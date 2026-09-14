package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatter_WithBeautify_Chaining(t *testing.T) {
	// Test that WithBeautify returns the formatter for chaining
	formatter := NewFormatter().WithBeautify()
	require.NotNil(t, formatter)
	require.Equal(t, FormatModeBeautify, formatter.mode)
}

func TestFormatter_WithIndent_Chaining(t *testing.T) {
	// Test that WithIndent returns the formatter for chaining
	formatter := NewFormatter().WithIndent("    ")
	require.NotNil(t, formatter)
	require.Equal(t, "    ", formatter.indent)
}

func TestFormatter_ChainedMethods(t *testing.T) {
	// Test that methods can be chained together
	formatter := NewFormatter().WithBeautify().WithIndent("\t")
	require.NotNil(t, formatter)
	require.Equal(t, FormatModeBeautify, formatter.mode)
	require.Equal(t, "\t", formatter.indent)
}

func TestFormatter_WithIndent_CustomIndentation(t *testing.T) {
	// Test actual formatting with custom indent using parsed SQL
	sql := "SELECT col1, col2 FROM table1 WHERE col1 > 10"

	parser := NewParser(sql)
	stmts, err := parser.ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	// Test with default 2-space indent
	formatter1 := NewFormatter().WithBeautify()
	formatter1.WriteExpr(stmts[0])
	result1 := formatter1.String()

	// Test with 4-space indent
	formatter2 := NewFormatter().WithBeautify().WithIndent("    ")
	formatter2.WriteExpr(stmts[0])
	result2 := formatter2.String()

	// Test with tab indent
	formatter3 := NewFormatter().WithBeautify().WithIndent("\t")
	formatter3.WriteExpr(stmts[0])
	result3 := formatter3.String()

	// Verify all results are different (due to different indentation)
	require.NotEqual(t, result1, result2)
	require.NotEqual(t, result1, result3)
	require.NotEqual(t, result2, result3)

	// Verify they all contain the basic SQL keywords
	require.Contains(t, result1, "SELECT")
	require.Contains(t, result2, "SELECT")
	require.Contains(t, result3, "SELECT")
	require.Contains(t, result1, "FROM")
	require.Contains(t, result2, "FROM")
	require.Contains(t, result3, "FROM")
}

func TestFormatter_DefaultIndent(t *testing.T) {
	// Test that default indent is 2 spaces
	formatter := NewFormatter()
	require.Equal(t, "  ", formatter.indent)
}

func TestFormatter_ExpressionGrouping(t *testing.T) {
	a, b, c := &Ident{Name: "a"}, &Ident{Name: "b"}, &Ident{Name: "c"}
	binary := func(op TokenKind, left, right Expr) *BinaryOperation {
		return &BinaryOperation{Operation: op, LeftExpr: left, RightExpr: right}
	}
	minus := func(expr Expr) *UnaryExpr { return &UnaryExpr{Kind: TokenKindMinus, Expr: expr} }
	not := func(expr Expr) *UnaryExpr { return &UnaryExpr{Kind: TokenKind(KeywordNot), Expr: expr} }
	ternary := &TernaryOperation{Condition: a, TrueExpr: b, FalseExpr: c}
	index := &NumberLiteral{Literal: "1", Base: 10}
	object, ok := parseSelectItemExpr(t, "SELECT arr[1]").(*ObjectParams)
	require.True(t, ok)
	array := object.Params
	for _, tc := range []struct {
		name string
		expr Expr
		sql  string
	}{
		{"sum under product", binary("*", binary("+", a, b), c), "(a + b) * c"},
		{"right subtraction", binary("-", a, binary("-", b, c)), "a - (b - c)"},
		{"right addition", binary("+", a, binary("+", b, c)), "a + (b + c)"},
		{"right division", binary("/", a, binary("*", b, c)), "a / (b * c)"},
		{"left associative", binary("-", binary("-", a, b), c), "a - b - c"},
		{"logical grouping", binary("AND", binary("OR", a, b), c), "(a OR b) AND c"},
		{"right logical grouping", binary("OR", a, binary("OR", b, c)), "a OR (b OR c)"},
		{"not logical", not(binary("AND", a, b)), "NOT (a AND b)"},
		{"not operand", binary("IN", not(a), b), "(NOT a) IN b"},
		{"comparison membership", binary("=", a, binary("IN", b, c)), "a = (b IN c)"},
		{"negative product", minus(binary("*", a, b)), "- (a * b)"},
		{"negative object", &IndexOperation{Object: minus(a), Operation: ".", Index: index}, "(- a).1"},
		{"numeric object", &IndexOperation{Object: index, Operation: ".", Index: index}, "(1).1"},
		{"sum object", &ObjectParams{Object: binary("+", a, b), Params: array}, "(a + b)[1]"},
		{"cast sum", binary("::", binary("+", a, b), &Ident{Name: "Int64"}), "(a + b)::Int64"},
		{"cast negative", binary("::", minus(a), &Ident{Name: "Int64"}), "(- a)::Int64"},
		{"null logical", &IsNullExpr{Expr: binary("OR", a, b)}, "(a OR b) IS NULL"},
		{"between operand", &BetweenClause{Expr: binary("OR", a, b), Between: b, And: c}, "(a OR b) BETWEEN b AND c"},
		{"between bound", &BetweenClause{Expr: a, Between: binary("AND", b, c), And: c}, "a BETWEEN (b AND c) AND c"},
		{"ternary condition", &TernaryOperation{Condition: ternary, TrueExpr: a, FalseExpr: b}, "(a ? b : c) ? a : b"},
		{"ternary branch", &TernaryOperation{Condition: a, TrueExpr: b, FalseExpr: ternary}, "a ? b : (a ? b : c)"},
		{"ternary logical branch", &TernaryOperation{Condition: a, TrueExpr: b, FalseExpr: binary("OR", b, c)}, "a ? b : (b OR c)"},
		{"ternary under sum", binary("+", ternary, c), "(a ? b : c) + c"},
		{"lambda branch", &TernaryOperation{Condition: a, TrueExpr: binary("->", b, c), FalseExpr: c}, "a ? (b -> c) : c"},
		{"nested lambda", binary("->", a, binary("->", b, c)), "a -> (b -> c)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.sql, Format(tc.expr))
			for _, beautify := range []bool{false, true} {
				formatter := NewFormatter()
				if beautify {
					formatter.WithBeautify()
				}
				formatter.WriteExpr(tc.expr)
				reparsed := parseSelectItemExpr(t, "SELECT "+formatter.String())
				require.Equal(t, expressionShape(tc.expr), expressionShape(reparsed))
			}
		})
	}
}

// Ignore positions and explicit grouping wrappers, while retaining every
// operator and operand. Formatting/reparsing must preserve this tree.
func expressionShape(expr Expr) string {
	switch e := expr.(type) {
	case *ColumnExpr:
		return expressionShape(e.Expr)
	case *ParamExprList:
		if len(e.Items.Items) == 1 {
			return expressionShape(e.Items.Items[0])
		}
	case *BinaryOperation:
		return fmt.Sprintf("(%s %t %t %s %s)", e.Operation, e.HasGlobal, e.HasNot,
			expressionShape(e.LeftExpr), expressionShape(e.RightExpr))
	case *UnaryExpr:
		return fmt.Sprintf("(%s %s)", e.Kind, expressionShape(e.Expr))
	case *TernaryOperation:
		return fmt.Sprintf("(?: %s %s %s)", expressionShape(e.Condition), expressionShape(e.TrueExpr), expressionShape(e.FalseExpr))
	case *IndexOperation:
		return fmt.Sprintf("(%s %s %s)", e.Operation, expressionShape(e.Object), expressionShape(e.Index))
	case *ObjectParams:
		return fmt.Sprintf("([] %s %s)", expressionShape(e.Object), Format(e.Params))
	case *BetweenClause:
		return fmt.Sprintf("(BETWEEN %t %s %s %s)", e.Not, expressionShape(e.Expr), expressionShape(e.Between), expressionShape(e.And))
	case *IsNullExpr:
		return "(IS NULL " + expressionShape(e.Expr) + ")"
	case *IsNotNullExpr:
		return "(IS NOT NULL " + expressionShape(e.Expr) + ")"
	}
	return fmt.Sprintf("%T:%s", expr, Format(expr))
}

func TestFormatter_EditedExpressionGrouping(t *testing.T) {
	expr, ok := parseSelectItemExpr(t, "SELECT a * b").(*BinaryOperation)
	require.True(t, ok)
	expr.LeftExpr = parseSelectItemExpr(t, "SELECT c + d")
	require.Equal(t, "(c + d) * b", Format(expr))
	reparsed := parseSelectItemExpr(t, "SELECT "+Format(expr))
	require.Equal(t, expressionShape(expr), expressionShape(reparsed))
}
