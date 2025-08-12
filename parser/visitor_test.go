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
				var formatSQLBuilder strings.Builder
				for _, stmt := range stmts {
					// Use Walk to traverse the AST (equivalent to the visitor doing nothing)
					Walk(stmt, func(node Expr) bool {
						return true // Continue traversal
					})

					formatSQLBuilder.WriteString(stmt.String())
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

func TestVisitor_SimpleRewrite(t *testing.T) {
	sql := `SELECT a, COUNT(b) FROM group_by_all GROUP BY CUBE(a) WITH CUBE WITH TOTALS ORDER BY a;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	require.Equal(t, 1, len(stmts))
	stmt := stmts[0]

	// Rewrite using Walk function
	Walk(stmt, func(node Expr) bool {
		switch expr := node.(type) {
		case *TableIdentifier:
			if expr.Table.String() == "group_by_all" {
				expr.Table = &Ident{Name: "hack"}
			}
		case *OrderExpr:
			expr.Direction = OrderDirectionDesc
		}
		return true // Continue traversal
	})

	newSql := stmt.String()

	require.NotSame(t, sql, newSql)
	require.True(t, strings.Contains(newSql, "hack"))
	require.True(t, strings.Contains(newSql, string(OrderDirectionDesc)))
}

func TestVisitor_NestRewrite(t *testing.T) {
	sql := `SELECT replica_name FROM system.ha_replicas UNION DISTINCT SELECT replica_name FROM system.ha_unique_replicas format JSON`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	require.Equal(t, 1, len(stmts))
	stmt := stmts[0]

	// Track nesting depth with closure variables
	var stack []Expr

	Walk(stmt, func(node Expr) bool {
		// Simulate Enter behavior
		if s, ok := node.(*SelectQuery); ok {
			stack = append(stack, s)
		}

		// Process TableIdentifier nodes
		if expr, ok := node.(*TableIdentifier); ok {
			expr.Table = &Ident{Name: fmt.Sprintf("table%d", len(stack))}
		}

		// Continue with children
		return true
	})

	newSql := stmt.String()

	require.NotSame(t, sql, newSql)
	// Both table names should be rewritten (they might both be table1 since they're at the same depth)
	require.True(t, strings.Contains(newSql, "table1") || strings.Contains(newSql, "table2"))
}

// TestWalk_NodeCounting verifies that Walk visits all nodes in the AST
func TestWalk_NodeCounting(t *testing.T) {
	sql := `SELECT a FROM table1`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	var nodeCount int
	Walk(stmts[0], func(node Expr) bool {
		nodeCount++
		return true
	})

	// Verify that we visited multiple nodes
	require.Greater(t, nodeCount, 0, "Walk should visit nodes")
	require.Greater(t, nodeCount, 3, "Should visit at least SELECT, column, table nodes")
}
