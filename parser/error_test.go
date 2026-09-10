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

func TestParseError_LexicalFailure(t *testing.T) {
	for _, prefix := range []string{"SELECT 1 ", "SELECT case ", "SELECT interval + ", "SELECT 1;\n", "SELECT 1 /* closed */\n"} {
		for _, suffix := range []struct {
			sql string
			msg string
		}{
			{"/*", "unclosed multi-line comment"},
			{"/* outer /* inner */", "unclosed multi-line comment"},
			{"$$unclosed", "invalid dollar-quoted string"},
			{"'unclosed", "invalid string"},
			{"`unclosed", "unclosed quoted identifier"},
			{"1e+", "exponent part should contain at least one digit"},
			{"中文", "unexpected character"},
		} {
			sql := prefix + suffix.sql
			t.Run(sql, func(t *testing.T) {
				stmts, err := NewParser(sql).ParseStmts()
				require.Error(t, err)
				require.Nil(t, stmts)

				var pe *ParseError
				require.ErrorAs(t, err, &pe)
				require.Contains(t, pe.Msg, suffix.msg)
				require.Equal(t, Pos(len(prefix)), pe.Pos)
				if strings.HasSuffix(prefix, "\n") {
					require.Equal(t, 2, pe.Line)
					require.Equal(t, 1, pe.Column)
				}
			})
		}
	}
}

func TestParser_TokenConsumptionError(t *testing.T) {
	for _, tc := range []struct {
		name  string
		sql   string
		parse func(*Parser) error
	}{
		{"required token", ") /*", func(p *Parser) error { return p.expectTokenKind(TokenKindRParen) }},
		{"optional token", ". /*", func(p *Parser) error {
			_, err := p.tryParseDotIdent(p.Pos())
			return err
		}},
		{"list separator", "a, /*", func(p *Parser) error {
			_, err := p.parseUserNames()
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.sql)
			require.NoError(t, p.lexer.consumeToken())
			err := tc.parse(p)
			var lexicalErr *lexerError
			require.ErrorAs(t, err, &lexicalErr)
			require.Equal(t, Pos(strings.Index(tc.sql, "/*")), lexicalErr.pos)
			require.EqualError(t, lexicalErr, "unclosed multi-line comment")
		})
	}
}

func TestParser_TryConsumeTokenKind(t *testing.T) {
	for _, tc := range []struct {
		name    string
		sql     string
		kind    TokenKind
		matched bool
		next    string
	}{
		{"mismatch", "a /*", TokenKindComma, false, "a"},
		{"advance", "a b", TokenKindIdent, true, "b"},
		{"last token", "a", TokenKindIdent, true, "<EOF>"},
		{"empty input", "", TokenKindIdent, false, "<EOF>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.sql)
			require.NoError(t, p.lexer.consumeToken())
			current := p.current()
			token, err := p.tryConsumeTokenKind(tc.kind)
			require.NoError(t, err)
			if tc.matched {
				require.Same(t, current, token)
			} else {
				require.Nil(t, token)
				require.Equal(t, current, p.current())
			}
			require.Equal(t, tc.next, p.currentTokenString())
		})
	}
}
