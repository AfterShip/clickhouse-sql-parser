package parser

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

var runCompatible = flag.Bool("compatible", false, "run compatible test")

func TestParser_Compatible(t *testing.T) {
	if !*runCompatible {
		t.Skip("Compatible test runs only if -compatible is set")
	}
	dir := "./testdata/query/compatible/1_stateful"
	entries, err := os.ReadDir(dir)
	if err != nil {
		require.NoError(t, err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			fileBytes, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			require.NoError(t, err)
			parser := Parser{
				lexer: NewLexer(string(fileBytes)),
			}
			_, err = parser.ParseStmts()
			require.NoError(t, err)
		})
	}
}

func TestParser_ParseStatements(t *testing.T) {
	for _, dir := range []string{"./testdata/dml", "./testdata/ddl", "./testdata/query", "./testdata/basic"} {
		outputDir := dir + "/output"
		entries, err := os.ReadDir(dir)
		if err != nil {
			require.NoError(t, err)
		}
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			t.Run(entry.Name(), func(t *testing.T) {
				fileBytes, err := os.ReadFile(filepath.Join(dir, entry.Name()))
				require.NoError(t, err)
				parser := Parser{
					lexer: NewLexer(string(fileBytes)),
				}
				stmts, err := parser.ParseStmts()
				require.NoError(t, err)
				outputBytes, _ := json.MarshalIndent(stmts, "", "  ")
				g := goldie.New(t,
					goldie.WithNameSuffix(".golden.json"),
					goldie.WithDiffEngine(goldie.ColoredDiff),
					goldie.WithFixtureDir(outputDir))
				g.Assert(t, entry.Name(), outputBytes)
			})
		}
	}
}

func TestParser_Format(t *testing.T) {
	for _, dir := range []string{"./testdata/dml", "./testdata/ddl", "./testdata/query", "./testdata/basic"} {
		outputDir := dir + "/format"

		entries, err := os.ReadDir(dir)
		if err != nil {
			require.NoError(t, err)
		}
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			t.Run(entry.Name(), func(t *testing.T) {
				fileBytes, err := os.ReadFile(filepath.Join(dir, entry.Name()))
				require.NoError(t, err)
				parser := Parser{
					lexer: NewLexer(string(fileBytes)),
				}
				stmts, err := parser.ParseStmts()
				require.NoError(t, err)
				var builder strings.Builder
				builder.WriteString("-- Origin SQL:\n")
				builder.Write(fileBytes)
				builder.WriteString("\n\n-- Format SQL:\n")
				var formatSQLBuilder strings.Builder
				for _, stmt := range stmts {
					formatSQLBuilder.WriteString(Format(stmt))
					formatSQLBuilder.WriteByte(';')
					formatSQLBuilder.WriteByte('\n')
				}
				formatSQL := formatSQLBuilder.String()
				builder.WriteString(formatSQL)
				validFormatSQL(t, formatSQL)
				g := goldie.New(t,
					goldie.WithNameSuffix(""),
					goldie.WithDiffEngine(goldie.ColoredDiff),
					goldie.WithFixtureDir(outputDir))
				g.Assert(t, entry.Name(), []byte(builder.String()))
			})
		}
	}
}

func TestParser_FormatBeautify(t *testing.T) {
	for _, dir := range []string{"./testdata/dml", "./testdata/ddl", "./testdata/query", "./testdata/basic"} {
		outputDir := dir + "/format/beautify"

		entries, err := os.ReadDir(dir)
		if err != nil {
			require.NoError(t, err)
		}
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			t.Run(entry.Name(), func(t *testing.T) {
				fileBytes, err := os.ReadFile(filepath.Join(dir, entry.Name()))
				require.NoError(t, err)
				parser := Parser{
					lexer: NewLexer(string(fileBytes)),
				}
				stmts, err := parser.ParseStmts()
				require.NoError(t, err)
				var builder strings.Builder
				builder.WriteString("-- Origin SQL:\n")
				builder.Write(fileBytes)
				builder.WriteString("\n\n-- Beautify SQL:\n")
				for _, stmt := range stmts {
					formatter := NewFormatter()
					formatter.WithBeautify()
					formatter.WriteExpr(stmt)
					builder.WriteString(formatter.String())
					builder.WriteByte(';')
					builder.WriteByte('\n')
				}
				g := goldie.New(t,
					goldie.WithNameSuffix(""),
					goldie.WithDiffEngine(goldie.ColoredDiff),
					goldie.WithFixtureDir(outputDir))
				g.Assert(t, entry.Name(), []byte(builder.String()))
			})
		}
	}
}

// validFormatSQL Verify that the format sql can be re-parsed with consistent results
func validFormatSQL(t *testing.T, sql string) {
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	var builder strings.Builder
	for _, stmt := range stmts {
		builder.WriteString(Format(stmt))
		builder.WriteByte(';')
		builder.WriteByte('\n')
	}
	require.Equal(t, sql, builder.String())
}

func TestParser_InvalidSyntax(t *testing.T) {
	invalidSQLs := []string{
		"SELECT * FROM",
		// WITH FILL error cases
		"SELECT n FROM t ORDER BY n WITH",                             // WITH without FILL
		"SELECT n FROM t ORDER BY n FILL",                             // FILL without WITH
		"SELECT n FROM t ORDER BY n WITH FILL FROM",                   // FROM without value
		"SELECT n FROM t ORDER BY n WITH FILL TO",                     // TO without value
		"SELECT n FROM t ORDER BY n WITH FILL STEP",                   // STEP without value
		"SELECT n FROM t ORDER BY n WITH FILL STALENESS",              // STALENESS without value
		"SELECT n FROM t ORDER BY n WITH FILL INTERPOLATE (x",         // Missing closing paren
		"SELECT n FROM t ORDER BY n WITH FILL INTERPOLATE x AS x + 1", // Missing parens around column list
		"ALTER TABLE foo_mv MODIFY QUERY AS SELECT * FROM baz",        // MODIFY QUERY followed by an invalid query
		// A sorting key is a single expression, so a key over several columns
		// is written as a tuple. A comma starts the next engine clause in
		// CREATE TABLE and the next clause in ALTER TABLE.
		"CREATE TABLE t (a DateTime, b String) ENGINE = MergeTree ORDER BY a, b",
		"ALTER TABLE t MODIFY ORDER BY a, b",
		// A broken column value must report an error, not stop the parser
		"CREATE TABLE t (x String DEFAULT CAST(a +, 'String'))",
		"CREATE TABLE t (x String MATERIALIZED a +)",
		"CREATE TABLE t (x String ALIAS a +)",
		// A TTL GROUP BY action only accepts a plain expression list; the
		// query-level modifiers are syntax errors for ClickHouse in a TTL.
		// (GROUP BY ALL and CUBE/ROLLUP(...) read as ordinary key
		// expressions there and are only rejected semantically.)
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id WITH TOTALS",
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id WITH ROLLUP",
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id WITH CUBE",
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY GROUPING SETS (id)",
		// OR REPLACE is only accepted for MATERIALIZED VIEW under CREATE,
		// never under ATTACH
		"ATTACH OR REPLACE MATERIALIZED VIEW mv TO dest AS SELECT * FROM src",
		// A TTL GROUP BY key list is greedy: after the keys the comma
		// continues the list, so a trailing TTL action after a
		// comma-separated key is rejected by ClickHouse (unlike after a
		// complete SET assignment, where the comma starts the next TTL
		// rule).
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id, created + INTERVAL 2 DAY DELETE",
		// A TTL WHERE clause belongs only to the DELETE action; ClickHouse
		// rejects it after GROUP BY, RECOMPRESS, and TO DISK/VOLUME.
		"CREATE TABLE t (id UInt64, created DateTime, x UInt64) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id SET x = sum(x) WHERE id > 0",
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id WHERE id > 0",
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY RECOMPRESS CODEC(ZSTD(1)) WHERE id > 0",
		"CREATE TABLE t (id UInt64, created DateTime) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY TO VOLUME 'v1' WHERE id > 0",
		// Invalid ARRAY JOIN types (only ARRAY JOIN, LEFT ARRAY JOIN, and INNER ARRAY JOIN are valid)
		"SELECT * FROM t RIGHT ARRAY JOIN arr AS a", // RIGHT ARRAY JOIN not supported
		"SELECT * FROM t FULL ARRAY JOIN arr AS a",  // FULL ARRAY JOIN not supported
		// ARRAY JOIN reads a column list, so it has no locality — in any spelling,
		// since modifiers keep their source case
		"SELECT * FROM t GLOBAL ARRAY JOIN arr AS a",
		"SELECT * FROM t LOCAL LEFT ARRAY JOIN arr AS a",
		"SELECT * FROM t GLOBAL array JOIN arr AS a",
		"SELECT * FROM t LOCAL left array JOIN arr AS a",
		// A join takes at most one locality, and it precedes the join type
		"SELECT * FROM t1 GLOBAL LOCAL JOIN t2 ON t1.a = t2.a",
		"SELECT * FROM t1 LEFT GLOBAL JOIN t2 ON t1.a = t2.a",
		// A locality without a join is not an alias, so it can't be dropped
		"SELECT * FROM t1 GLOBAL",
		"SELECT * FROM t1 LOCAL",
		"SELECT * FROM t global",
		"SELECT * FROM t1 GLOBAL JOIN t2 ON t1.a = t2.a GLOBAL",
		"SELECT * FROM t1 GLOBAL WHERE a = 1",
		"SELECT * FROM t1 LOCAL SETTINGS max_threads = 1",
		// GLOBAL is a keyword everywhere, never an implicit alias
		"SELECT a GLOBAL FROM t",
		"SELECT a GLOBAL, b FROM t",
		// GLOBAL is an operator only before IN/NOT IN, and LOCAL is never one
		"SELECT * FROM t WHERE a GLOBAL LIKE 'x'",
		"SELECT * FROM t WHERE a GLOBAL NOT LIKE 'x'",
		"SELECT * FROM t WHERE a GLOBAL BETWEEN 1 AND 2",
		"SELECT * FROM t WHERE a NOT GLOBAL IN (SELECT 1)",
		"SELECT * FROM t WHERE a GLOBAL GLOBAL IN (SELECT 1)",
		"SELECT * FROM t WHERE a LOCAL IN (SELECT 1)",
		"SELECT * FROM t WHERE a LOCAL NOT IN (SELECT 1)",
		"00e1d",    // invalid number that leaves curToken nil
		"CREATE--", // trailing comment pushes p.Pos() past end of input (wrapError out-of-range)
		// Inputs that previously caused a nil-pointer dereference while
		// formatting the error message at EOF (p.current() is nil).
		"ALTER ",
		"SELECT*FROM A(0A",
		"SET A=",
		// An operator as the last token of the input is missing its right
		// operand. The infix loop used to stop before it, so the operator
		// either reached the statement parser as unexpected trailing input or,
		// when it was a keyword, was read as an implicit alias.
		"SELECT a +",
		"SELECT a GLOBAL",
		"SELECT a REGEXP",
		"SELECT * FROM t WHERE a AND",
		// A parenthesized select must still be closed and UNION still needs
		// ALL or DISTINCT
		"(SELECT 1",
		"(SELECT 1) UNION SELECT 2",
		// ClickHouse rejects a set operator once SETTINGS is bound to a
		// parenthesized group
		"(SELECT 1) SETTINGS max_threads=1 UNION ALL SELECT 2",
	}
	for _, sql := range invalidSQLs {
		parser := NewParser(sql)
		_, err := parser.ParseStmts()
		require.Error(t, err, "Expected error for SQL: %s", sql)
	}
}

func TestParser_ParenthesizedSetOperationOperands(t *testing.T) {
	// A parenthesized operand becomes a group node, so the operator after
	// ')' binds to the whole group instead of leaking into its chain.
	stmts, err := NewParser("(SELECT 1 UNION DISTINCT SELECT 2) UNION ALL SELECT 3").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	group, ok := stmts[0].(*SelectQuery)
	require.True(t, ok)
	require.NotNil(t, group.InnerQuery)
	require.NotNil(t, group.InnerQuery.UnionDistinct)
	require.NotNil(t, group.UnionAll)
	require.Nil(t, group.InnerQuery.UnionDistinct.UnionAll)

	stmts, err = NewParser("SELECT a FROM ((SELECT 1 AS a) UNION ALL (SELECT 2 AS a))").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	outer, ok := stmts[0].(*SelectQuery)
	require.True(t, ok)
	joinTable, ok := outer.From.Expr.(*JoinTableExpr)
	require.True(t, ok)
	subQuery, ok := joinTable.Table.Expr.(*SubQuery)
	require.True(t, ok)
	require.NotNil(t, subQuery.Select.InnerQuery)
	require.NotNil(t, subQuery.Select.UnionAll)
	require.NotNil(t, subQuery.Select.UnionAll.InnerQuery)

	// Grouping survives the round trip: ClickHouse gives INTERSECT higher
	// precedence than UNION, so dropping the parens would change semantics.
	sql := "(SELECT 1 UNION ALL SELECT 2) INTERSECT SELECT 2"
	stmts, err = NewParser(sql).ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	require.Equal(t, sql, Format(stmts[0]))

	// SETTINGS and FORMAT after ')' bind to the group.
	stmts, err = NewParser("(SELECT 1) SETTINGS max_threads=1 FORMAT JSONEachRow").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	group, ok = stmts[0].(*SelectQuery)
	require.True(t, ok)
	require.NotNil(t, group.InnerQuery)
	require.NotNil(t, group.Settings)
	require.NotNil(t, group.Format)
}
