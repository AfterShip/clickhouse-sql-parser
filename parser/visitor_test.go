package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestVisitor_Identical(t *testing.T) {
	visitor := DefaultASTVisitor{}

	for _, dir := range []string{"./testdata/dml", "./testdata/ddl", "./testdata/query", "./testdata/basic"} {
		outputDir := dir + "/format"

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
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
				for _, stmt := range stmts {
					err := stmt.Accept(&visitor)
					require.NoError(t, err)

					builder.WriteString(stmt.String(0))
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

type simpleRewriteVisitor struct {
	DefaultASTVisitor
}

func (v *simpleRewriteVisitor) VisitTableIdentifier(expr *TableIdentifier) error {
	if expr.Table.String(0) == "group_by_all" {
		expr.Table = &Ident{Name: "hack"}
	}
	return nil
}

func (v *simpleRewriteVisitor) VisitOrderByExpr(expr *OrderExpr) error {
	expr.Direction = OrderDirectionDesc
	return nil
}

func TestVisitor_SimpleRewrite(t *testing.T) {
	visitor := simpleRewriteVisitor{}

	sql := `SELECT a, COUNT(b) FROM group_by_all GROUP BY CUBE(a) WITH CUBE WITH TOTALS ORDER BY a;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	require.Equal(t, 1, len(stmts))
	stmt := stmts[0]

	err = stmt.Accept(&visitor)
	require.NoError(t, err)
	newSql := stmt.String(0)

	require.NotSame(t, sql, newSql)
	require.True(t, strings.Contains(newSql, "hack"))
	require.True(t, strings.Contains(newSql, string(OrderDirectionDesc)))
}

type nestedRewriteVisitor struct {
	DefaultASTVisitor
	stack []Expr
}

func (v *nestedRewriteVisitor) VisitTableIdentifier(expr *TableIdentifier) error {
	expr.Table = &Ident{Name: fmt.Sprintf("table%d", len(v.stack))}
	return nil
}

func (v *nestedRewriteVisitor) enter(expr Expr) {
	if s, ok := expr.(*SelectQuery); ok {
		v.stack = append(v.stack, s)
	}
}

func (v *nestedRewriteVisitor) leave(expr Expr) {
	if _, ok := expr.(*SelectQuery); ok {
		v.stack = v.stack[1:]
	}
}

func TestVisitor_NestRewrite(t *testing.T) {
	visitor := nestedRewriteVisitor{}

	sql := `SELECT replica_name FROM system.ha_replicas UNION DISTINCT SELECT replica_name FROM system.ha_unique_replicas format JSON`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	require.Equal(t, 1, len(stmts))
	stmt := stmts[0]

	err = stmt.Accept(&visitor)
	require.NoError(t, err)
	newSql := stmt.String(0)

	require.NotSame(t, sql, newSql)
	require.Less(t, strings.Index(newSql, "table1"), strings.Index(newSql, "table2"))
}
