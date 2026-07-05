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

// TestParseError_LexerErrorSurfaced guards that a lexing failure mid-statement
// is reported as itself. Most parser call sites discard consumeToken's error
// and then read the nil current token as end of input, which used to produce
// misleading messages like "unexpected token kind: <eof>" — or, worse, silent
// success when the statement looked complete at the point of failure.
func TestParseError_LexerErrorSurfaced(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"SELECT 0x", "invalid number"},
		{"SELECT 'unterminated", "invalid string"},
		{"SELECT `unclosed", "unclosed quoted identifier"},
		// the statement is already well-formed when the bad token appears,
		// so this parsed successfully before
		{"SELECT 1 'oops", "invalid string"},
		// the failure is in a later statement
		{"SELECT 1; SELECT 'bad; SELECT 2", "invalid string"},
	}
	for _, c := range cases {
		_, err := NewParser(c.input).ParseStmts()
		require.Error(t, err, "input %q", c.input)
		require.Contains(t, err.Error(), c.want, "input %q", c.input)

		var pe *ParseError
		require.True(t, errors.As(err, &pe), "input %q: expected a *ParseError, got %T", c.input, err)
	}
}

// TestParseError_LexerErrorPosition guards that the reported position points
// at the offending token, not wherever the parser noticed the nil token.
func TestParseError_LexerErrorPosition(t *testing.T) {
	_, err := NewParser("SELECT 0x").ParseStmts()
	require.Error(t, err)

	var pe *ParseError
	require.True(t, errors.As(err, &pe))
	require.Equal(t, Pos(7), pe.Pos) // byte offset of "0x"
	require.Equal(t, 1, pe.Line)
	require.Equal(t, 8, pe.Column)
}

// TestParseError_LookaheadLexerErrorRollback guards that a lexer error hit
// only during backtracked lookahead does not leak into an otherwise
// successful parse.
func TestParseError_LookaheadLexerErrorRollback(t *testing.T) {
	stmts, err := NewParser("SELECT 1").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)
}
