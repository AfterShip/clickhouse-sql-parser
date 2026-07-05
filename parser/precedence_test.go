package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// parseSelectItemExpr parses a single-statement SELECT and returns the first
// projection expression, so tests can assert the tree structure directly.
func parseSelectItemExpr(t *testing.T, sql string) Expr {
	t.Helper()
	stmts, err := NewParser(sql).ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	selectQuery, ok := stmts[0].(*SelectQuery)
	require.True(t, ok, "expected *SelectQuery, got %T", stmts[0])
	require.NotEmpty(t, selectQuery.SelectItems)
	return selectQuery.SelectItems[0].Expr
}

func TestLambdaBodyBindsLoosest(t *testing.T) {
	// `x -> x + 1` is a lambda whose body is the whole `x + 1`, not `(x -> x) + 1`.
	expr := parseSelectItemExpr(t, "SELECT arrayMap(x -> x + 1, arr)")
	fn, ok := expr.(*FunctionExpr)
	require.True(t, ok)
	firstArg, ok := fn.Params.Items.Items[0].(*ColumnExpr)
	require.True(t, ok)
	lambda, ok := firstArg.Expr.(*BinaryOperation)
	require.True(t, ok)
	require.Equal(t, TokenKindArrow, lambda.Operation)
	body, ok := lambda.RightExpr.(*BinaryOperation)
	require.True(t, ok, "lambda body should be the binary operation `x + 1`, got %T", lambda.RightExpr)
	require.Equal(t, TokenKind("+"), body.Operation)
}

func TestLambdaIsRightAssociative(t *testing.T) {
	expr := parseSelectItemExpr(t, "SELECT x -> y -> x + y")
	outer, ok := expr.(*BinaryOperation)
	require.True(t, ok)
	require.Equal(t, TokenKindArrow, outer.Operation)
	_, ok = outer.LeftExpr.(*Ident)
	require.True(t, ok, "outer lambda parameter should be `x`, got %T", outer.LeftExpr)
	inner, ok := outer.RightExpr.(*BinaryOperation)
	require.True(t, ok)
	require.Equal(t, TokenKindArrow, inner.Operation)
}

func TestNotBindsLooserThanComparison(t *testing.T) {
	// `NOT a = b` negates the whole comparison: NOT (a = b).
	expr := parseSelectItemExpr(t, "SELECT NOT a = b")
	not, ok := expr.(*UnaryExpr)
	require.True(t, ok, "expected UnaryExpr at the top, got %T", expr)
	require.Equal(t, TokenKind("NOT"), not.Kind)
	cmp, ok := not.Expr.(*BinaryOperation)
	require.True(t, ok, "NOT operand should be the comparison, got %T", not.Expr)
	require.Equal(t, TokenKind("="), cmp.Operation)
}

func TestDoubleNot(t *testing.T) {
	expr := parseSelectItemExpr(t, "SELECT NOT NOT a")
	outer, ok := expr.(*UnaryExpr)
	require.True(t, ok)
	require.Equal(t, TokenKind("NOT"), outer.Kind)
	inner, ok := outer.Expr.(*UnaryExpr)
	require.True(t, ok, "inner expression should be the second NOT, got %T", outer.Expr)
	require.Equal(t, TokenKind("NOT"), inner.Kind)
	_, ok = inner.Expr.(*Ident)
	require.True(t, ok)
}

func TestNotStopsAtAnd(t *testing.T) {
	// NOT binds tighter than AND: `NOT a AND b` is `(NOT a) AND b`.
	expr := parseSelectItemExpr(t, "SELECT NOT a AND b")
	and, ok := expr.(*BinaryOperation)
	require.True(t, ok, "expected AND at the top, got %T", expr)
	require.Equal(t, TokenKind("AND"), and.Operation)
	_, ok = and.LeftExpr.(*UnaryExpr)
	require.True(t, ok, "left side of AND should be `NOT a`, got %T", and.LeftExpr)
}

func TestTernaryBindsBelowOr(t *testing.T) {
	// `a OR b ? 1 : 2` groups as `(a OR b) ? 1 : 2`.
	expr := parseSelectItemExpr(t, "SELECT a OR b ? 1 : 2")
	ternary, ok := expr.(*TernaryOperation)
	require.True(t, ok, "expected TernaryOperation at the top, got %T", expr)
	cond, ok := ternary.Condition.(*BinaryOperation)
	require.True(t, ok, "ternary condition should be `a OR b`, got %T", ternary.Condition)
	require.Equal(t, TokenKind("OR"), cond.Operation)
}

func TestNotBetween(t *testing.T) {
	expr := parseSelectItemExpr(t, "SELECT x NOT BETWEEN 1 AND 2")
	between, ok := expr.(*BetweenClause)
	require.True(t, ok, "expected BetweenClause, got %T", expr)
	require.True(t, between.Not)
	require.Equal(t, "x NOT BETWEEN 1 AND 2", Format(expr))

	expr = parseSelectItemExpr(t, "SELECT x BETWEEN 1 AND 2")
	between, ok = expr.(*BetweenClause)
	require.True(t, ok)
	require.False(t, between.Not)
}

func TestNotInGroupsLikeIn(t *testing.T) {
	// NOT IN must bind with IN's precedence: `a = b NOT IN (1)` groups as
	// `a = (b NOT IN (1))`, exactly like `a = b IN (1)`.
	for _, sql := range []string{"SELECT a = b IN (1)", "SELECT a = b NOT IN (1)"} {
		expr := parseSelectItemExpr(t, sql)
		eq, ok := expr.(*BinaryOperation)
		require.True(t, ok, "%s: expected `=` at the top, got %T", sql, expr)
		require.Equal(t, TokenKind("="), eq.Operation, sql)
		in, ok := eq.RightExpr.(*BinaryOperation)
		require.True(t, ok, "%s: right side of `=` should be the IN operation, got %T", sql, eq.RightExpr)
		require.Contains(t, string(in.Operation), "IN", sql)
	}
}

func TestIntersect(t *testing.T) {
	stmts, err := NewParser("SELECT 1 INTERSECT SELECT 2").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	selectQuery, ok := stmts[0].(*SelectQuery)
	require.True(t, ok)
	require.NotNil(t, selectQuery.Intersect)
	require.Equal(t, "SELECT 1 INTERSECT SELECT 2", Format(stmts[0]))
}
