package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWalk_BasicTraversal(t *testing.T) {
	sql := `SELECT a, COUNT(b) FROM table1 WHERE id > 10 ORDER BY a;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	var nodeCount int
	Walk(stmts[0], func(node Expr) bool {
		nodeCount++
		return true
	})

	// Verify we visited multiple nodes
	require.Greater(t, nodeCount, 10, "Should have visited more than 10 nodes")
}

func TestWalk_JoinExpr(t *testing.T) {
	sql := `SELECT a, COUNT(b) FROM table1 JOIN table2 ON table1.id = table2.id;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	var onClauses int
	Walk(stmts[0], func(node Expr) bool {
		if _, ok := node.(*OnClause); ok {
			onClauses++
		}
		return true
	})

	require.Equal(t, 1, onClauses, "Should have visited exactly 1 OnClause")
}

func TestWalkWithBreak_EarlyTermination(t *testing.T) {
	sql := `SELECT a, COUNT(b) FROM table1 WHERE id > 10 ORDER BY a;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	var nodeCount int
	result := WalkWithBreak(stmts[0], func(node Expr) bool {
		nodeCount++
		// Stop after visiting 5 nodes
		return nodeCount < 5
	})

	require.False(t, result, "WalkWithBreak should return false when terminated early")
	require.Equal(t, 5, nodeCount, "Should have stopped at exactly 5 nodes")
}

func TestFind_FirstMatch(t *testing.T) {
	sql := `SELECT a, COUNT(b) FROM table1 WHERE id > 10;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	// Find the first FunctionExpr
	found, exists := Find(stmts[0], func(node Expr) bool {
		_, ok := node.(*FunctionExpr)
		return ok
	})

	require.True(t, exists, "Should find a FunctionExpr")
	require.NotNil(t, found)

	funcExpr, ok := found.(*FunctionExpr)
	require.True(t, ok, "Found node should be a FunctionExpr")
	require.Equal(t, "COUNT", funcExpr.Name.String(), "Should find the COUNT function")
}

func TestFindAll_MultipleMatches(t *testing.T) {
	sql := `SELECT a, COUNT(b), MAX(c) FROM table1 WHERE id > 10;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	// Find all FunctionExprs
	functions := FindAll(stmts[0], func(node Expr) bool {
		_, ok := node.(*FunctionExpr)
		return ok
	})

	require.Equal(t, 2, len(functions), "Should find 2 function expressions")

	funcNames := make([]string, len(functions))
	for i, fn := range functions {
		funcExpr := fn.(*FunctionExpr)
		funcNames[i] = funcExpr.Name.String()
	}

	require.Contains(t, funcNames, "COUNT")
	require.Contains(t, funcNames, "MAX")
}

func TestWalk_TableIdentifierRewriting(t *testing.T) {
	sql := `SELECT a, COUNT(b) FROM group_by_all GROUP BY CUBE(a) WITH CUBE WITH TOTALS ORDER BY a;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	// Rewrite table names
	Walk(stmts[0], func(node Expr) bool {
		if tableId, ok := node.(*TableIdentifier); ok {
			if tableId.Table.String() == "group_by_all" {
				tableId.Table = &Ident{Name: "hack"}
			}
		}
		return true
	})

	newSQL := stmts[0].String()
	require.Contains(t, newSQL, "hack", "Table name should be rewritten to 'hack'")
	require.NotContains(t, newSQL, "group_by_all", "Original table name should be gone")
}

func TestWalk_OrderByDirectionRewriting(t *testing.T) {
	sql := `SELECT a, COUNT(b) FROM table1 ORDER BY a ASC, b;`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	// Change all order directions to DESC
	Walk(stmts[0], func(node Expr) bool {
		if orderExpr, ok := node.(*OrderExpr); ok {
			orderExpr.Direction = OrderDirectionDesc
		}
		return true
	})

	newSQL := stmts[0].String()
	require.Contains(t, newSQL, string(OrderDirectionDesc), "Should contain DESC direction")
}

func TestWalk_NestedQueryDepthTracking(t *testing.T) {
	sql := `SELECT replica_name FROM system.ha_replicas UNION DISTINCT SELECT replica_name FROM system.ha_unique_replicas`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	var tableNames []string

	Walk(stmts[0], func(node Expr) bool {
		// Track nesting depth
		if tableID, ok := node.(*JoinTableExpr); ok {
			tableName := tableID.Table.String()
			tableNames = append(tableNames, tableName+"@depth")
		}
		return true
	})

	require.Len(t, tableNames, 2, "Should find 2 table identifiers")
}

func TestWalk_SimpleNodeCounting(t *testing.T) {
	sql := `SELECT a FROM table1`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	var nodeCount int
	Walk(stmts[0], func(node Expr) bool {
		nodeCount++
		return true
	})

	require.Greater(t, nodeCount, 0, "Walk should visit nodes")
	require.Greater(t, nodeCount, 3, "Should visit at least SELECT, column, table nodes")
}

func TestFind_NoMatch(t *testing.T) {
	sql := `SELECT a FROM table1`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	// Try to find a non-existent node type
	found, exists := Find(stmts[0], func(node Expr) bool {
		// Look for AlterTable in a SELECT statement (should not exist)
		_, ok := node.(*AlterTable)
		return ok
	})

	require.False(t, exists, "Should not find AlterTable in SELECT statement")
	require.Nil(t, found, "Found node should be nil when not found")
}

func TestFindAll_EmptyResult(t *testing.T) {
	sql := `SELECT a FROM table1`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)

	// Try to find non-existent node types
	results := FindAll(stmts[0], func(node Expr) bool {
		// Look for AlterTable in a SELECT statement (should not exist)
		_, ok := node.(*AlterTable)
		return ok
	})

	require.Empty(t, results, "Should return empty slice when no matches found")
}

func TestWalk_ShowStmtNewFields(t *testing.T) {
	// Test that Walk properly traverses all new fields in ShowStmt
	sql := `SHOW DATABASES LIKE 'prod%' LIMIT 5 INTO OUTFILE '/tmp/prod_dbs.txt' FORMAT JSON`
	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	_, ok := stmts[0].(*ShowStmt)
	require.True(t, ok, "Statement should be ShowStmt")

	// Collect all nodes during walk
	var foundNodes []Expr
	Walk(stmts[0], func(node Expr) bool {
		foundNodes = append(foundNodes, node)
		return true
	})

	// Should find the ShowStmt itself plus all its expression fields
	require.Greater(t, len(foundNodes), 4, "Should visit at least ShowStmt + 4 expression fields")

	// Find specific types of expressions that should be walked
	var stringLiterals []*StringLiteral
	var numberLiterals []*NumberLiteral

	for _, node := range foundNodes {
		switch n := node.(type) {
		case *StringLiteral:
			stringLiterals = append(stringLiterals, n)
		case *NumberLiteral:
			numberLiterals = append(numberLiterals, n)
		}
	}

	// Should find exactly 3 string literals: LIKE pattern, OUTFILE path, FORMAT type
	require.Equal(t, 3, len(stringLiterals), "Should find exactly 3 StringLiteral nodes")

	// Should find exactly 1 number literal: LIMIT value
	require.Equal(t, 1, len(numberLiterals), "Should find exactly 1 NumberLiteral node")

	// Verify the specific values
	stringValues := make([]string, len(stringLiterals))
	for i, sl := range stringLiterals {
		stringValues[i] = sl.Literal
	}
	require.Contains(t, stringValues, "prod%", "Should contain LIKE pattern")
	require.Contains(t, stringValues, "/tmp/prod_dbs.txt", "Should contain OUTFILE path")
	require.Contains(t, stringValues, "JSON", "Should contain FORMAT type")

	require.Equal(t, "5", numberLiterals[0].Literal, "Should contain LIMIT value")
}
