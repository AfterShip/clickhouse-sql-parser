package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStmt(t *testing.T) {
	t.Run("single statement", func(t *testing.T) {
		stmt, err := ParseStmt("SELECT a FROM t")
		require.NoError(t, err)
		require.IsType(t, &SelectQuery{}, stmt)
	})

	t.Run("single statement with trailing semicolon", func(t *testing.T) {
		stmt, err := ParseStmt("SELECT a FROM t;")
		require.NoError(t, err)
		require.IsType(t, &SelectQuery{}, stmt)
	})

	t.Run("empty input", func(t *testing.T) {
		_, err := ParseStmt("")
		require.ErrorContains(t, err, "no statement")
	})

	t.Run("multiple statements", func(t *testing.T) {
		_, err := ParseStmt("SELECT 1; SELECT 2")
		require.ErrorContains(t, err, "exactly one statement")
	})

	t.Run("syntax error", func(t *testing.T) {
		_, err := ParseStmt("SELECT a FROM WHERE")
		require.Error(t, err)
	})
}

func TestParseExpr(t *testing.T) {
	t.Run("column reference", func(t *testing.T) {
		expr, err := ParseExpr("a")
		require.NoError(t, err)
		require.IsType(t, &Ident{}, expr)
	})

	t.Run("function call", func(t *testing.T) {
		expr, err := ParseExpr("toDate(created_at) + 1")
		require.NoError(t, err)
		require.Equal(t, "toDate(created_at) + 1", Format(expr))
	})

	t.Run("case expression", func(t *testing.T) {
		expr, err := ParseExpr("CASE WHEN a > 1 THEN 'x' ELSE 'y' END")
		require.NoError(t, err)
		require.IsType(t, &CaseExpr{}, expr)
	})

	t.Run("empty input", func(t *testing.T) {
		_, err := ParseExpr("")
		require.Error(t, err)
	})

	t.Run("trailing tokens", func(t *testing.T) {
		_, err := ParseExpr("a + 1 b")
		require.ErrorContains(t, err, "unexpected token after expression")
	})

	t.Run("syntax error", func(t *testing.T) {
		_, err := ParseExpr("f(")
		require.Error(t, err)
	})
}

func TestFormatBeautify(t *testing.T) {
	stmt, err := ParseStmt("SELECT a, b FROM t WHERE a = 1")
	require.NoError(t, err)

	beautified := FormatBeautify(stmt)
	require.Contains(t, beautified, "\n")

	// Must match the long-form formatter API it wraps.
	formatter := NewFormatter().WithBeautify()
	formatter.WriteExpr(stmt)
	require.Equal(t, formatter.String(), beautified)

	// Beautified SQL must still parse to the same compact form.
	reparsed, err := ParseStmt(beautified)
	require.NoError(t, err)
	require.Equal(t, Format(stmt), Format(reparsed))

	require.Equal(t, "", FormatBeautify(nil))
	require.False(t, strings.Contains(Format(stmt), "\n"))
}
