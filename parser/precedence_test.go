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

func TestSignedNumberAfterClosingBracketIsBinaryOperator(t *testing.T) {
	// A closing `)` or `]` ends an expression, so the following `+`/`-` is a
	// binary operator and not the sign of the next numeric literal.
	for _, sql := range []string{"SELECT (1)-1", "SELECT arr[1]-1", "SELECT f()-1"} {
		expr := parseSelectItemExpr(t, sql)
		op, ok := expr.(*BinaryOperation)
		require.True(t, ok, "%s: expected BinaryOperation at the top, got %T", sql, expr)
		require.Equal(t, TokenKindMinus, op.Operation, sql)
		right, ok := op.RightExpr.(*NumberLiteral)
		require.True(t, ok, "%s: right side should be the unsigned literal `1`, got %T", sql, op.RightExpr)
		require.Equal(t, "1", right.Literal, sql)
	}

	// A `+`/`-` after an opening bracket is still a sign, not an operator.
	expr := parseSelectItemExpr(t, "SELECT arr[-1]")
	require.Equal(t, "arr[-1]", Format(expr))
}

// keywordArgOf returns the sole argument of a call, unwrapping its ColumnExpr.
func keywordArgOf(t *testing.T, sql, name string) Expr {
	t.Helper()
	expr := parseSelectItemExpr(t, sql)
	fn, ok := expr.(*FunctionExpr)
	require.True(t, ok, "%s: expected FunctionExpr, got %T", sql, expr)
	require.Equal(t, name, fn.Name.Name, sql)
	require.Len(t, fn.Params.Items.Items, 1, sql)
	arg, ok := fn.Params.Items.Items[0].(*ColumnExpr)
	require.True(t, ok, "%s: argument should be wrapped in ColumnExpr, got %T", sql, fn.Params.Items.Items[0])
	return arg.Expr
}

func TestTrimKeywordArgs(t *testing.T) {
	from, ok := keywordArgOf(t, "SELECT trim(BOTH ' ' FROM ' x ')", "trim").(*BinaryOperation)
	require.True(t, ok)
	require.Equal(t, TokenKind(KeywordFrom), from.Operation)
	modifier, ok := from.LeftExpr.(*UnaryExpr)
	require.True(t, ok, "left side should be the BOTH modifier, got %T", from.LeftExpr)
	require.Equal(t, TokenKind(KeywordBoth), modifier.Kind)
	require.Equal(t, "SELECT trim(BOTH ' ' FROM ' x ')", Format(parseOneStmt(t, "SELECT trim(BOTH ' ' FROM ' x ')")))

	for _, sql := range []string{
		"SELECT trim(LEADING '0' FROM '00042')",
		"SELECT trim(TRAILING '/' FROM '/api/')",
	} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	// The modifier, the characters and FROM are all required.
	for _, sql := range []string{
		"SELECT trim(' ' FROM ' x ')",
		"SELECT trim(BOTH FROM ' x ')",
		"SELECT trim(BOTH ' ')",
	} {
		_, err := NewParser(sql).ParseStmts()
		require.Error(t, err, sql)
	}

	// BOTH is not reserved, so it stays usable as an ordinary name. ClickHouse
	// rejects trim(both) itself — its trim grammar always reads BOTH as the
	// modifier — so only the non-trim spellings are asserted.
	for _, sql := range []string{"SELECT both FROM t", "SELECT max(both) FROM t", "SELECT a AS both FROM t"} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	// Only `trim` has the keyword form; the explicit variants take commas alone,
	// as in ClickHouse.
	for _, sql := range []string{
		"SELECT trimBoth('xxhixx', 'x')",
		"SELECT trimLeft('xxhi', 'x')",
		"SELECT trimRight('hixx', 'x')",
	} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	for _, sql := range []string{
		"SELECT trimBoth(BOTH ' ' FROM ' x ')",
		"SELECT trimLeft(LEADING ' ' FROM ' x ')",
		"SELECT trimRight(TRAILING ' ' FROM ' x ')",
	} {
		_, err := NewParser(sql).ParseStmts()
		require.Error(t, err, sql)
	}

	// An argument alias works as it does for any other function.
	for _, sql := range []string{"SELECT trim('x' AS y)", "SELECT trim(BOTH ' ' FROM ' x ' AS y)"} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	// A trailing comma argument here is a ClickHouse arity error, not a syntax
	// error, so the parser accepts it.
	for _, sql := range []string{"SELECT trim(BOTH ' ' FROM ' x ', 'y')", "SELECT trim('  x  ', 'y')"} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}
}

func TestSubstringKeywordArgs(t *testing.T) {
	// substring(s FROM 2 FOR 3) nests left-to-right: ((s FROM 2) FOR 3).
	forOp, ok := keywordArgOf(t, "SELECT substring('hello' FROM 2 FOR 3)", "substring").(*BinaryOperation)
	require.True(t, ok)
	require.Equal(t, TokenKind(KeywordFor), forOp.Operation)
	fromOp, ok := forOp.LeftExpr.(*BinaryOperation)
	require.True(t, ok, "left side should be the FROM operation, got %T", forOp.LeftExpr)
	require.Equal(t, TokenKind(KeywordFrom), fromOp.Operation)

	// FOR is optional, FROM is not, and the arguments are ordinary expressions.
	for _, sql := range []string{
		"SELECT substring('hello' FROM 2)",
		"SELECT substring(concat(a, b) FROM x + 1 FOR len(y))",
		"SELECT substring('hello', 2, 3)",
	} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	// A comma fills the same slot a separator keyword would, so substring mixes
	// the two freely.
	for _, sql := range []string{
		"SELECT substring('hello', 2 FOR 3)",
		"SELECT substring('hello' FROM 2, 3)",
	} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	// A fourth slot is a syntax error however the earlier separators were spelled.
	for _, sql := range []string{
		"SELECT substring('hello', 2 FOR 3, 4)",
		"SELECT substring('hello' FROM 2, 3, 4)",
		"SELECT substring('hello' FROM 2 FOR 3, 4)",
	} {
		_, err := NewParser(sql).ParseStmts()
		require.Error(t, err, sql)
	}

	// substringUTF8 has no keyword form, unlike overlayUTF8.
	for _, sql := range []string{
		"SELECT substring('hello' FOR 3)",
		"SELECT substringUTF8('hello' FROM 2 FOR 3)",
	} {
		_, err := NewParser(sql).ParseStmts()
		require.Error(t, err, sql)
	}

	// FROM keeps its clause meaning outside these argument lists.
	require.Equal(t, "SELECT substring FROM t", Format(parseOneStmt(t, "SELECT substring FROM t")))
}

func TestOverlayKeywordArgs(t *testing.T) {
	// overlay(s PLACING r FROM n FOR m) chains three separators in order.
	forOp, ok := keywordArgOf(t, "SELECT overlay('Hello world' PLACING 'SQL' FROM 7 FOR 5)", "overlay").(*BinaryOperation)
	require.True(t, ok)
	require.Equal(t, TokenKind(KeywordFor), forOp.Operation)
	fromOp, ok := forOp.LeftExpr.(*BinaryOperation)
	require.True(t, ok, "expected the FROM operation, got %T", forOp.LeftExpr)
	require.Equal(t, TokenKind(KeywordFrom), fromOp.Operation)
	placingOp, ok := fromOp.LeftExpr.(*BinaryOperation)
	require.True(t, ok, "expected the PLACING operation, got %T", fromOp.LeftExpr)
	require.Equal(t, TokenKind(KeywordPlacing), placingOp.Operation)

	for _, sql := range []string{
		"SELECT overlay('Hello world' PLACING 'SQL' FROM 7)",
		"SELECT overlayUTF8(s PLACING r FROM x + 1 FOR len(y)) FROM t",
		"SELECT overlay(a, b, c) FROM t",
	} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	// Separators only count in the declared order, and both PLACING and FROM
	// are required — ClickHouse rejects overlay(s PLACING r).
	for _, sql := range []string{
		"SELECT overlay(s FROM 2 PLACING r)",
		"SELECT overlay(s FOR 3)",
		"SELECT overlay('Hello world' PLACING 'SQL')",
	} {
		_, err := NewParser(sql).ParseStmts()
		require.Error(t, err, sql)
	}

	// Unlike substring, overlay cannot mix commas with its keyword form in
	// either direction.
	for _, sql := range []string{
		"SELECT overlay('Hello world' PLACING 'SQL' FROM 1, 2)",
		"SELECT overlay('Hello world' PLACING 'SQL' FROM 1 FOR 2, 3)",
		"SELECT overlay('Hello world', 'SQL' FROM 1)",
		"SELECT overlay('Hello world', 'SQL', 1 FOR 2)",
		"SELECT overlayUTF8('Hello world' PLACING 'SQL' FROM 1, 2)",
		"SELECT overlayUTF8('Hello world', 'SQL' FROM 1)",
	} {
		_, err := NewParser(sql).ParseStmts()
		require.Error(t, err, sql)
	}

	// The all-comma form keeps its own arity.
	for _, sql := range []string{
		"SELECT overlay('Hello world', 'SQL', 1)",
		"SELECT overlay('Hello world', 'SQL', 1, 2)",
	} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}

	// OVERLAY and PLACING are non-reserved, so they stay usable as names.
	for _, sql := range []string{"SELECT overlay, placing FROM t", "SELECT a AS overlay FROM t", "SELECT max(placing) FROM t"} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
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

func TestGlobalJoinLocalityPrefixesTheJoinType(t *testing.T) {
	// the locality prefixes the join type, and both have to reach the modifiers
	for _, tc := range []struct {
		sql       string
		modifiers []string
	}{
		{"SELECT * FROM t1 GLOBAL JOIN t2 ON t1.a = t2.a", []string{"GLOBAL", "JOIN"}},
		{"SELECT * FROM t1 GLOBAL INNER JOIN t2 ON t1.a = t2.a", []string{"GLOBAL", "INNER", "JOIN"}},
		{"SELECT * FROM t1 GLOBAL LEFT OUTER JOIN t2 ON t1.a = t2.a", []string{"GLOBAL", "LEFT", "OUTER", "JOIN"}},
		{"SELECT * FROM t1 GLOBAL ANY LEFT JOIN t2 ON t1.a = t2.a", []string{"GLOBAL", "ANY", "LEFT", "JOIN"}},
		{"SELECT * FROM t1 GLOBAL CROSS JOIN t2", []string{"GLOBAL", "CROSS", "JOIN"}},
		{"SELECT * FROM t1 AS x LOCAL FULL JOIN t2 USING a", []string{"LOCAL", "FULL", "JOIN"}},
	} {
		from := parseOneStmt(t, tc.sql).(*SelectQuery).From.Expr
		join, ok := from.(*JoinExpr)
		require.True(t, ok, "%s: expected a *JoinExpr in FROM, got %T", tc.sql, from)

		right, ok := join.Right.(*JoinExpr)
		require.True(t, ok, "%s: expected the join to carry a right side, got %T", tc.sql, join.Right)
		require.Equal(t, tc.modifiers, right.Modifiers, tc.sql)
		require.Equal(t, tc.sql, Format(parseOneStmt(t, tc.sql)), tc.sql)
	}
}

func TestLocalityWithoutAJoinKeepsItsKeyword(t *testing.T) {
	// a GLOBAL with no join used to be consumed and dropped from Format(); the
	// error now has to name it rather than the end of the statement
	_, err := NewParser("SELECT * FROM t1 GLOBAL").ParseStmts()
	require.ErrorContains(t, err, "GLOBAL")

	// an explicit alias still names a table, as it does on the server
	for _, sql := range []string{"SELECT * FROM t AS global", "SELECT * FROM t AS local"} {
		require.Equal(t, sql, Format(parseOneStmt(t, sql)), sql)
	}
}

func TestGlobalAfterJoinConstraintStartsANewJoin(t *testing.T) {
	// a GLOBAL after an ON/USING clause starts the next join, not a GLOBAL IN
	sql := "SELECT * FROM t1 GLOBAL JOIN t2 ON t1.a = t2.a GLOBAL LEFT JOIN t3 ON t1.a = t3.a"
	stmt := parseOneStmt(t, sql)
	join := stmt.(*SelectQuery).From.Expr.(*JoinExpr)

	first, ok := join.Right.(*JoinExpr)
	require.True(t, ok)
	require.Equal(t, []string{"GLOBAL", "JOIN"}, first.Modifiers)

	second, ok := first.Right.(*JoinExpr)
	require.True(t, ok, "the second GLOBAL should start another join, got %T", first.Right)
	require.Equal(t, []string{"GLOBAL", "LEFT", "JOIN"}, second.Modifiers)
	require.Equal(t, sql, Format(stmt))
}

func TestGlobalInAndGlobalNotIn(t *testing.T) {
	for _, tc := range []struct {
		sql string
		op  TokenKind
	}{
		{"SELECT a GLOBAL IN (SELECT 1)", "GLOBAL IN"},
		{"SELECT a GLOBAL NOT IN (SELECT 1)", "GLOBAL NOT IN"},
	} {
		expr := parseSelectItemExpr(t, tc.sql)
		in, ok := expr.(*BinaryOperation)
		require.True(t, ok, "%s: expected BinaryOperation at the top, got %T", tc.sql, expr)
		require.Equal(t, tc.op, in.Operation, tc.sql)
	}
}

func TestGlobalInGroupsLikeIn(t *testing.T) {
	// GLOBAL IN binds like IN: `a = b GLOBAL IN (1)` is `a = (b GLOBAL IN (1))`
	for _, sql := range []string{"SELECT a = b GLOBAL IN (1)", "SELECT a = b GLOBAL NOT IN (1)"} {
		expr := parseSelectItemExpr(t, sql)
		eq, ok := expr.(*BinaryOperation)
		require.True(t, ok, "%s: expected `=` at the top, got %T", sql, expr)
		require.Equal(t, TokenKind("="), eq.Operation, sql)

		in, ok := eq.RightExpr.(*BinaryOperation)
		require.True(t, ok, "%s: right side of `=` should be the GLOBAL IN operation, got %T", sql, eq.RightExpr)
		require.Contains(t, string(in.Operation), "GLOBAL", sql)
	}
}
