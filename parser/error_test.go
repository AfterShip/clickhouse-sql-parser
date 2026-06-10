package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseError_Structured(t *testing.T) {
	// The error is on the second line, so line/column must reflect the
	// multi-line offset rather than a flat byte count.
	_, err := NewParser("SELECT 1\nFROM 123").ParseStmts()
	require.Error(t, err)

	var pe *ParseError
	require.True(t, errors.As(err, &pe), "expected a *ParseError, got %T", err)
	require.Equal(t, 2, pe.Line)
	require.GreaterOrEqual(t, pe.Column, 1)
	require.Contains(t, pe.Error(), "line 2:")
	// The rendered message includes the offending source line and a caret.
	require.Contains(t, pe.Error(), "FROM 123")
	require.Contains(t, pe.Error(), "^")
}

func TestParseError_ExpectedKeyword(t *testing.T) {
	// "IF" must be followed by "EXISTS" / "NOT EXISTS"; the failure carries the
	// expected keyword structurally.
	_, err := NewParser("DROP TABLE IF foo").ParseStmts()
	require.Error(t, err)

	var pe *ParseError
	require.True(t, errors.As(err, &pe))
	require.Equal(t, "EXISTS", pe.Keyword)
	require.Equal(t, 1, pe.Line)
}

func TestParseError_ExpectedTokenKind(t *testing.T) {
	// An unclosed function-call paren flows through expectTokenKind, so the
	// failure carries the expected token kind structurally.
	_, err := NewParser("SELECT count(a").ParseStmts()
	require.Error(t, err)

	var pe *ParseError
	require.True(t, errors.As(err, &pe))
	require.Equal(t, []TokenKind{TokenKindRParen}, pe.Expected)
	require.True(t, strings.HasPrefix(pe.Error(), "line "))
}
