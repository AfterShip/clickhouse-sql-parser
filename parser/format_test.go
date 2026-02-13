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

func TestFormatter_WithIdent_Chaining(t *testing.T) {
	// Test that WithIdent returns the formatter for chaining
	formatter := NewFormatter().WithIdent("    ")
	require.NotNil(t, formatter)
	require.Equal(t, "    ", formatter.ident)
}

func TestFormatter_ChainedMethods(t *testing.T) {
	// Test that methods can be chained together
	formatter := NewFormatter().WithBeautify().WithIdent("\t")
	require.NotNil(t, formatter)
	require.Equal(t, FormatModeBeautify, formatter.mode)
	require.Equal(t, "\t", formatter.ident)
}

func TestFormatter_WithIdent_CustomIndentation(t *testing.T) {
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
	formatter2 := NewFormatter().WithBeautify().WithIdent("    ")
	formatter2.WriteExpr(stmts[0])
	result2 := formatter2.String()

	// Test with tab indent
	formatter3 := NewFormatter().WithBeautify().WithIdent("\t")
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

func TestFormatter_DefaultIdent(t *testing.T) {
	// Test that default indent is 2 spaces
	formatter := NewFormatter()
	require.Equal(t, "  ", formatter.ident)
}
