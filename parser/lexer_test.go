package parser

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumeComment(t *testing.T) {
	comments := []string{
		"-- hello world",
		"-- hello world\n",
		"-- hello world\r\n",
		"-- hello world\r",
		"/* hello world */",
		"/* hello world */\n",
		"/* hello world */\r\n",
		"/* hello world */\r",
		"/* hello world */ /* hello world */",
		"/* hello world */ /* hello world */\n",
		"/* hello world */ /* hello world */\r\n",
		"/* hello world */ /* hello world */\r",
	}
	for _, c := range comments {
		lexer := NewLexer(c)
		err := lexer.consumeToken()
		require.NoError(t, err)
	}

}

// TestConsumeUnterminatedComment guards against an infinite loop (a DoS hang)
// when a block comment is never closed. consumeMultiLineComment previously
// looped on isEOF() while only advancing a local index, so l.offset never
// reached EOF and the lexer spun forever. The test runs in a goroutine with a
// timeout so a regression fails fast instead of hanging the whole test binary.
func TestConsumeUnterminatedComment(t *testing.T) {
	inputs := []string{
		"/*",
		"/* unterminated",
		"/* unterminated *",
		"SELECT 1 /* unterminated",
	}
	for _, c := range inputs {
		c := c
		done := make(chan struct{})
		go func() {
			defer close(done)
			lexer := NewLexer(c)
			for !lexer.isEOF() {
				if err := lexer.consumeToken(); err != nil {
					break
				}
			}
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatalf("lexer did not terminate on unterminated comment: %q", c)
		}
	}
}

func TestConsumeString(t *testing.T) {
	t.Run("Simple strings", func(t *testing.T) {
		strs := []string{
			"'hello world'",
			"'123'",
		}
		for _, s := range strs {
			lexer := NewLexer(s)
			err := lexer.consumeToken()
			require.NoError(t, err)
			require.Equal(t, TokenKindString, lexer.currentToken.Kind)
			require.Equal(t, strings.Trim(s, "'"), lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})

	t.Run("Strings with backslash-escaped quotes", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{`'hello\'world'`, `hello\'world`},
			{`'test\''`, `test\'`},
			{`'\'abc\''`, `\'abc\'`},
		}
		for _, tc := range testCases {
			lexer := NewLexer(tc.input)
			err := lexer.consumeToken()
			require.NoError(t, err, "Failed to parse: %s", tc.input)
			require.Equal(t, TokenKindString, lexer.currentToken.Kind)
			require.Equal(t, tc.expected, lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})

	t.Run("Strings with double single quotes", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{`'hello''world'`, `hello''world`},
			{`'test''123'`, `test''123`},
			{`'abc''def''ghi'`, `abc''def''ghi`},
		}
		for _, tc := range testCases {
			lexer := NewLexer(tc.input)
			err := lexer.consumeToken()
			require.NoError(t, err, "Failed to parse: %s", tc.input)
			require.Equal(t, TokenKindString, lexer.currentToken.Kind)
			require.Equal(t, tc.expected, lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})

	t.Run("Strings with backslash-escaped backslashes", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{`'a\\b'`, `a\\b`},
			{`'test\\123'`, `test\\123`},
		}
		for _, tc := range testCases {
			lexer := NewLexer(tc.input)
			err := lexer.consumeToken()
			require.NoError(t, err, "Failed to parse: %s", tc.input)
			require.Equal(t, TokenKindString, lexer.currentToken.Kind)
			require.Equal(t, tc.expected, lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})
}

func TestConsumeNumber(t *testing.T) {
	t.Run("Integer number", func(t *testing.T) {
		integers := []string{
			"123",
			"123e+10",
			"123e-10",
			"123e10",
			"123E10",
			"123E+10",
			"123E-10",
		}
		for _, i := range integers {
			lexer := NewLexer(i)
			err := lexer.consumeToken()
			require.NoError(t, err)
			require.Equal(t, TokenKindInt, lexer.currentToken.Kind)
			require.Equal(t, 10, lexer.currentToken.Base)
			require.Equal(t, i, lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})

	t.Run("Hexadecimal number", func(t *testing.T) {
		numbers := []string{
			"0x123",
			"0x1",
		}
		for _, n := range numbers {
			lexer := NewLexer(n)
			err := lexer.consumeToken()
			require.NoError(t, err)
			require.Equal(t, TokenKindInt, lexer.currentToken.Kind)
			require.Equal(t, 16, lexer.currentToken.Base)
			require.Equal(t, n, lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})

	t.Run("Invalid number", func(t *testing.T) {
		invalidNumbers := []string{
			"123e",
			"123e+",
			"123e-",
			"123e",
			"123E",
			"123E+",
			"123E-",
			"0x",
			"0xg",
		}
		for _, n := range invalidNumbers {
			lexer := NewLexer(n)
			err := lexer.consumeToken()
			require.Error(t, err)
		}
	})

	t.Run("Float number", func(t *testing.T) {
		floats := []string{
			"123.456",
			"123.456e+10",
			"123.456e-10",
			"123.456e10",
			"123.456E10",
			"123.456E+10",
			"123.456E-10",
		}
		for _, f := range floats {
			lexer := NewLexer(f)
			err := lexer.consumeToken()
			require.NoError(t, err)
			require.Equal(t, TokenKindFloat, lexer.currentToken.Kind)
			require.Equal(t, f, lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})

	t.Run("Invalid float number", func(t *testing.T) {
		invalidFloats := []string{
			"123.456b",
			"123.456e",
			"123.456e+",
			"123.456e-",
			"123.456e+10e",
			"123.456e-10e",
			"123.456e10e",
			"123.456E10e",
			"123.456E+10e",
			"123.456E-10e",
			"123.456e+10e+10",
		}
		for _, f := range invalidFloats {
			lexer := NewLexer(f)
			err := lexer.consumeToken()
			assert.Error(t, err)
		}
	})

	t.Run("Name", func(t *testing.T) {
		idents := []string{
			"`CASE`",
			"`TEST`",
			"`WHEN`",
			"hello",
			"hello_world",
			"hello123",
			"hello_123",
			"hello_123_world",
			"hello_123_world_456",
			"hello_123_world_456_789",
			"hello_123_world_456_789_abc",
			"hello_123_world_456_789_abc_def",
			"hello_123_world_456_789_abc_def_ghi",
			"hello_123_world_456_789_abc_def_ghi_jkl",
			"hello_123_world_456_789_abc_def_ghi_jkl_mno",
			"hello_123_world_456_789_abc_def_ghi_jkl_mno_pqr",
		}
		for _, i := range idents {
			lexer := NewLexer(i)
			err := lexer.consumeToken()
			require.NoError(t, err)
			require.Equal(t, TokenKindIdent, lexer.currentToken.Kind)
			require.Equal(t, strings.Trim(i, "`"), lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})

	t.Run("Keyword", func(t *testing.T) {
		for _, k := range keywords.Members() {
			lexer := NewLexer(k)
			err := lexer.consumeToken()
			require.NoError(t, err)
			require.Equal(t, TokenKindKeyword, lexer.currentToken.Kind)
			require.Equal(t, k, lexer.currentToken.String)
			require.True(t, lexer.isEOF())
		}
	})
}
