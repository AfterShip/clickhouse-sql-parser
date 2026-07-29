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
		// Invalid ARRAY JOIN types (only ARRAY JOIN, LEFT ARRAY JOIN, and INNER ARRAY JOIN are valid)
		"SELECT * FROM t RIGHT ARRAY JOIN arr AS a", // RIGHT ARRAY JOIN not supported
		"SELECT * FROM t FULL ARRAY JOIN arr AS a",  // FULL ARRAY JOIN not supported
		// ARRAY JOIN reads a column list, so it has no GLOBAL/LOCAL locality.
		// Join modifiers keep their source spelling, so the lowercase spellings
		// have to be rejected too.
		"SELECT * FROM t GLOBAL ARRAY JOIN arr AS a",
		"SELECT * FROM t LOCAL LEFT ARRAY JOIN arr AS a",
		"SELECT * FROM t GLOBAL array JOIN arr AS a",
		"SELECT * FROM t LOCAL left array JOIN arr AS a",
		// A join takes at most one locality, and it precedes the join type
		"SELECT * FROM t1 GLOBAL LOCAL JOIN t2 ON t1.a = t2.a",
		"SELECT * FROM t1 LEFT GLOBAL JOIN t2 ON t1.a = t2.a",
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
	}
	for _, sql := range invalidSQLs {
		parser := NewParser(sql)
		_, err := parser.ParseStmts()
		require.Error(t, err, "Expected error for SQL: %s", sql)
	}
}
