package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestVisitor_Identical(t *testing.T) {
	visitor := NewDefaultASTVisitor(nil, nil, nil)

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
				stmts, err := parser.ParseStatements()
				require.NoError(t, err)
				var builder strings.Builder
				builder.WriteString("-- Origin SQL:\n")
				builder.Write(fileBytes)
				builder.WriteString("\n\n-- Format SQL:\n")
				for _, stmt := range stmts {
					e, err := stmt.Accept(visitor)
					require.NoError(t, err)

					builder.WriteString(e.String(0))
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

type testRewriteVisitor struct {
	ASTVisitor
}

func (v *testRewriteVisitor) VisitTableIdentifier(expr *TableIdentifier) (Expr, error) {
	if expr.Table.Name == "group_by_all" {
		expr.Table.Name = "hack"
	}
	return expr, nil
}

func TestVisitor_Rewrite(t *testing.T) {
	visitor := testRewriteVisitor{
		ASTVisitor: NewDefaultASTVisitor(nil, nil, nil),
	}

	sql := `SELECT a, COUNT(b) FROM group_by_all GROUP BY CUBE(a) WITH CUBE WITH TOTALS ORDER BY a;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStatements()
	require.NoError(t, err)

	require.Equal(t, len(stmts), 1)
	stmt := stmts[0]

	newStmt, err := stmt.Accept(&visitor)
	require.NoError(t, err)
	newSql := newStmt.String(0)

	require.NotSame(t, sql, newSql)
	require.True(t, strings.Contains(newSql, "hack"))
}
