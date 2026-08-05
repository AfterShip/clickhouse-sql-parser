package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReservedKeywordsAreKeywords(t *testing.T) {
	for _, kw := range reservedKeywords.Members() {
		require.True(t, keywords.Contains(kw),
			"reserved keyword %q is missing from the keywords set", kw)
	}
}

// TestReservedKeywordRejectedAsIdentifier asserts that a reserved keyword no
// longer silently fills an identifier slot: a missing name before a clause
// keyword must fail at the keyword instead of swallowing it (the bug class
// behind #268/#269: e.g. `SELECT a FROM WHERE b = 1` used to parse FROM as a
// bare alias of `a`).
func TestReservedKeywordRejectedAsIdentifier(t *testing.T) {
	cases := []string{
		"SELECT a FROM WHERE b = 1",                  // FROM must not become an alias of `a`
		"SELECT a, b FROM GROUP BY a",                // FROM must not become an alias of `b`
		"INSERT INTO SELECT 1",                       // SELECT must not become a table name
		"SELECT a AS FROM t",                         // explicit alias FROM consumes the clause keyword
		"SELECT a FROM t JOIN ON a = b",              // ON must not become a table name
		"CREATE TABLE t (from String) ENGINE=Memory", // reserved keyword as column name needs quoting
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := NewParser(sql).ParseStmts()
			require.Error(t, err)
		})
	}
}

// TestNonReservedKeywordAsIdentifier asserts that non-reserved keywords keep
// working as identifiers anywhere an identifier is expected.
func TestNonReservedKeywordAsIdentifier(t *testing.T) {
	cases := []string{
		"SELECT key FROM t",
		"SELECT date, first, last, timestamp FROM t",
		"CREATE TABLE t (key String, date Date) ENGINE=Memory",
		"SELECT t.key FROM t",
		"SELECT if(a, 1, 2), any(b), left(c, 1) FROM t",
		"SELECT * FROM t WHERE key = 1",
		"SELECT max(end) FROM t",
		"SELECT * FROM t WHERE end > start",
		"INSERT INTO t (end) VALUES (1)",
		"CREATE TABLE t (end DateTime) ENGINE=Memory",
		"SELECT CASE WHEN end > start THEN end ELSE start END FROM t",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
		})
	}
}

// TestCaseExprRequiresEnd asserts that END, which is non-reserved so that
// columns can be named `end`, still terminates a CASE expression: a CASE
// without it must fail instead of running into the next clause.
func TestCaseExprRequiresEnd(t *testing.T) {
	cases := []string{
		"SELECT CASE WHEN a THEN 1 ELSE 2 FROM t",
		"SELECT CASE WHEN a THEN 1 FROM t",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := NewParser(sql).ParseStmts()
			require.Error(t, err)
		})
	}
}

// TestReservedKeywordInDisambiguatedPositions asserts that reserved keywords
// are still accepted as names where context proves they cannot start a clause:
// after AS, after a dot in a qualified name, lookahead-disambiguated select
// items, query parameters, window names, and GRANT options.
func TestReservedKeywordInDisambiguatedPositions(t *testing.T) {
	cases := []string{
		"SELECT 1 AS from",
		"SELECT 1 AS interval, 2 AS from, 3 AS limit",
		"SELECT * FROM t AS from",
		"SELECT a FROM db.from",
		"SELECT t.from FROM t",
		"SELECT a, limit FROM t",
		"SELECT case;",
		"SELECT limit",
		"SELECT a FROM t WHERE ts < {end:UInt32}",
		"SELECT sum(x) OVER (order) FROM t WINDOW order AS (PARTITION BY team)",
		"SELECT sum(x) OVER order FROM t WINDOW order AS (PARTITION BY team)",
		"GRANT SELECT(x) ON db.table TO john WITH GRANT OPTION",
		"GRANT SELECT ON db.from TO john",
		"CREATE TABLE t (j JSON(max_dynamic_paths=1, SKIP a.from)) ENGINE=Memory",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
		})
	}
}

func TestExpressionContinuationDoesNotWeakenClauseParsing(t *testing.T) {
	cases := []string{
		"SELECT a FROM WHERE b = 1",
		"SELECT a, b FROM GROUP BY a",
		"SELECT a AS FROM t",
		"SELECT a FROM t JOIN ON a = b",
		"SELECT CASE WHEN a THEN 1 ELSE 2 FROM t",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := NewParser(sql).ParseStmts()
			require.Error(t, err)
		})
	}
}

func TestClauseKeywordsStillParseAsClauses(t *testing.T) {
	cases := []string{
		"SELECT a FROM t LIMIT 10",
		"SELECT a FROM t LIMIT 10 OFFSET 5",
		"SELECT a FROM t LIMIT 1 BY a",
		"SELECT a FROM t ORDER BY a DESC LIMIT 5",
		"SELECT DISTINCT a FROM t",
		"SELECT a FROM t UNION ALL SELECT b FROM s",
		"SELECT a FROM t WHERE a BETWEEN 1 AND 2",
		"SELECT INTERVAL 1 DAY",
		"SELECT a + INTERVAL 1 DAY FROM t",
		"SELECT * EXCEPT (a) FROM t",
		"SELECT a FROM t FORMAT JSON",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
		})
	}
}

// TestReservedOperatorKeywordsAreCallable covers the regressions from the
// review of #275: operator keywords double as ordinary ClickHouse function
// names and must stay callable when followed by '('.
func TestReservedOperatorKeywordsAreCallable(t *testing.T) {
	inputs := []string{
		"SELECT and(a, b) FROM t",
		"SELECT or(a, b) FROM t",
		"SELECT in(1, [1])",
		"SELECT like(s, '%a%') FROM t",
		"SELECT ilike(s, '%a%') FROM t",
	}
	for _, sql := range inputs {
		t.Run(sql, func(t *testing.T) {
			stmts, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
			require.Len(t, stmts, 1)
		})
	}
}

// TestReservedKeywordAliasesAfterAs covers the review regressions of #275:
// AS proves the next token is an alias name, so even reserved keywords are
// accepted in expression lists, ORDER BY, and non-parenthesized CTEs.
func TestReservedKeywordAliasesAfterAs(t *testing.T) {
	inputs := []string{
		"SELECT (1 AS from)",
		"SELECT sum(x AS from) FROM t",
		"SELECT a FROM t ORDER BY x AS from",
		"WITH 1 AS from SELECT from",
	}
	for _, sql := range inputs {
		t.Run(sql, func(t *testing.T) {
			stmts, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
			require.Len(t, stmts, 1)
		})
	}
}
