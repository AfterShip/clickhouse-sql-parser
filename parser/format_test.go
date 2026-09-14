package parser

import (
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

func TestFormatter_ExplicitExpressionParentheses(t *testing.T) {
	for _, sql := range []string{
		"a + b * c",
		"(a + b) * c",
		"a - b - c",
		"a - (b - c)",
		"((a))",
		"(a OR b) AND c",
		"NOT (a AND b)",
		"-tuple(1, 2).1",
		"-arr[1]",
		"-x::Int64",
		"+x",
		"- -x",
		"-(-x)",
		"- -1",
		"- -1::Int64",
		"- - -x",
		"(-a).1",
		"(1).1",
		"a::Int64[1]",
		"(a + b)::Int64",
		"a = b IN (1)",
		"a = (b IN (1))",
		"a BETWEEN (b AND c) AND c",
		"a ? b : (a ? b : c)",
		"a -> (b -> c)",
	} {
		t.Run(sql, func(t *testing.T) {
			expr := parseSelectItemExpr(t, "SELECT "+sql)
			require.Equal(t, sql, Format(expr))
			formatter := NewFormatter().WithBeautify()
			formatter.WriteExpr(expr)
			reparsed := parseSelectItemExpr(t, "SELECT "+formatter.String())
			require.Equal(t, sql, Format(reparsed))
		})
	}
}

func TestFormatter_EditedExpressionRequiresExplicitParentheses(t *testing.T) {
	expr, ok := parseSelectItemExpr(t, "SELECT a * b").(*BinaryOperation)
	require.True(t, ok)
	expr.LeftExpr = parseSelectItemExpr(t, "SELECT c + d")
	require.Equal(t, "c + d * b", Format(expr))
	expr.LeftExpr = parseSelectItemExpr(t, "SELECT (c + d)")
	require.Equal(t, "(c + d) * b", Format(expr))
}
