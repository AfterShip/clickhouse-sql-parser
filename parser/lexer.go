package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	TokenKindEOF          TokenKind = "<eof>"
	TokenKindIdent        TokenKind = "<ident>"
	TokenKindKeyword      TokenKind = "<keyword>"
	TokenKindInt          TokenKind = "<int>"
	TokenKindFloat        TokenKind = "<float>"
	TokenKindString       TokenKind = "<string>"
	TokenKindDot          TokenKind = "."
	TokenKindSingleEQ     TokenKind = "="
	TokenKindDoubleEQ     TokenKind = "=="
	TokenKindNE           TokenKind = "!="
	TokenKindLT           TokenKind = "<"
	TokenKindLE           TokenKind = "<="
	TokenKindGT           TokenKind = ">"
	TokenKindGE           TokenKind = ">="
	TokenKindQuestionMark TokenKind = "?"

	TokenKindPlus   TokenKind = "+"
	TokenKindMinus  TokenKind = "-"
	TokenKindMul    TokenKind = "*"
	TokenKindDiv    TokenKind = "/"
	TokenKindMod    TokenKind = "%"
	TokenKindConcat TokenKind = "||"

	TokenKindArrow TokenKind = "->"
	TokenKindDash  TokenKind = "::"

	TokenKindLParen   TokenKind = "("
	TokenKindRParen   TokenKind = ")"
	TokenKindLBrace   TokenKind = "{"
	TokenKindRBrace   TokenKind = "}"
	TokenKindLBracket TokenKind = "["
	TokenKindRBracket TokenKind = "]"

	TokenKindComma  TokenKind = ","
	TokenKindColon  TokenKind = ":"
	TokenKindAtSign TokenKind = "@"
)

const (
	Unquoted = iota + 1
	DoubleQuote
	BackTicks
	SingleQuote
)

type Pos int
type TokenKind string

type Token struct {
	Pos Pos
	End Pos

	Kind      TokenKind
	String    string
	Base      int // 10 or 16 on TokenKindInt
	QuoteType int
}

func (t *Token) ToString() string {
	if t.Kind == TokenKindKeyword {
		return strings.ToUpper(t.String)
	}
	return t.String
}

type lexerState struct {
	offset       int    // byte offset into input of the next unread character
	currentToken *Token // current lookahead token; nil at end of input
}

type Lexer struct {
	lexerState

	input string
}

func NewLexer(buf string) *Lexer {
	return &Lexer{input: buf}
}

func (l *Lexer) saveState() lexerState {
	return l.lexerState
}

func (l *Lexer) restoreState(state lexerState) {
	l.lexerState = state
}

func (l *Lexer) skipN(n int) {
	l.offset += n
}

func (l *Lexer) slice(i, j int) string {
	return l.input[l.offset+i : l.offset+j]
}

func (l *Lexer) peekN(n int) byte {
	return l.input[l.offset+n]
}

func (l *Lexer) peekOk(n int) bool {
	return l.offset+n < len(l.input)
}

func (l *Lexer) isKeyword(ident string) bool {
	return keywords.Contains(ident)
}

func (l *Lexer) consumeNumber() error {
	i := 0
	base := 10
	if l.peekN(0) == '+' || l.peekN(0) == '-' {
		// skip sign
		i++
	}
	if l.peekOk(i+1) && l.peekN(i) == '0' && l.peekN(i+1) == 'x' {
		i += 2
		base = 16
	}

	hasExp := false
	hasDot := false
	tokenKind := TokenKindInt
	// hasNumberPart is only set when a digit is actually consumed, so that a
	// bare prefix followed by a non-digit ("0x;") is rejected rather than
	// lexed as an empty hex literal
	hasNumberPart := false
	for l.peekOk(i) {
		c := l.peekN(i)
		switch {
		case base == 10 && IsDigit(c):
			hasNumberPart = true
			i++
			continue
		case base == 16 && IsHexDigit(c):
			hasNumberPart = true
			i++
			continue
		case c == '.': // float
			// a second dot ("1.2.3"), a dot after the exponent ("1e2.3") or
			// a dot in a hex literal ("0x1.8") cannot start a valid number tail
			if hasDot || hasExp || base == 16 {
				return errors.New("invalid number")
			}
			hasDot = true
			tokenKind = TokenKindFloat
			i++
			continue
		case base != 16 && (c == 'e' || c == 'E'):
			if hasExp {
				return errors.New("invalid number")
			}
			i++
			if l.peekOk(i) && (l.peekN(i) == '+' || l.peekN(i) == '-') {
				i++
			}
			if !l.peekOk(i) || !IsDigit(l.peekN(i)) {
				return errors.New("exponent part should contain at least one digit")
			}
			hasExp = true
			// scientific notation always denotes a floating-point value
			tokenKind = TokenKindFloat
			continue
		}
		break
	}
	if (l.peekOk(i) && IsIdentPart(l.peekN(i))) || !hasNumberPart {
		return errors.New("invalid number")
	}
	l.currentToken = &Token{
		Kind:   tokenKind,
		String: l.slice(0, i),
		Pos:    Pos(l.offset),
		End:    Pos(l.offset + i),
		Base:   base,
	}
	l.skipN(i)
	return nil
}

func (l *Lexer) consumeIdent(_ Pos) error {
	token := &Token{}
	quoteType := Unquoted
	if l.peekOk(0) && (l.peekN(0) == '`' || l.peekN(0) == '"') {
		if l.peekOk(0) && l.peekN(0) == '`' {
			quoteType = BackTicks
		} else {
			quoteType = DoubleQuote
		}
		l.skipN(1)
	}

	i := 0
	if quoteType == Unquoted {
		if l.peekOk(i) && l.peekN(i) == '$' {
			i++
		}
		for l.peekOk(i) && IsIdentPart(l.peekN(i)) {
			i++
		}
	} else {
		for l.peekOk(i) && (quoteType == BackTicks && l.peekN(i) != '`' ||
			quoteType == DoubleQuote && l.peekN(i) != '"') {
			i++
		}
		if !l.peekOk(i) || (quoteType == BackTicks && l.peekN(i) != '`') ||
			(quoteType == DoubleQuote && l.peekN(i) != '"') {
			return fmt.Errorf("unclosed quoted identifier: %s", l.slice(0, i))
		}
	}
	slice := l.slice(0, i)
	if quoteType == Unquoted && l.isKeyword(strings.ToUpper(slice)) {
		token.Kind = TokenKindKeyword
	} else {
		token.Kind = TokenKindIdent
	}
	token.Pos = Pos(l.offset)
	token.End = Pos(l.offset + i)
	token.String = slice
	token.QuoteType = quoteType
	l.currentToken = token

	l.skipN(i)
	if quoteType != Unquoted {
		l.skipN(1)
	}
	return nil
}

func (l *Lexer) consumeSingleLineComment() {
	l.skipN(2)
	i := 0
	for l.peekOk(i) && l.peekN(i) != '\r' && l.peekN(i) != '\n' {
		i++
	}
	if l.peekOk(i) {
		// consume the newline too; at EOF there is none to consume
		i++
	}
	l.skipN(i)
}

func (l *Lexer) consumeMultiLineComment() error {
	l.skipN(2)
	i := 0
	for l.peekOk(i) {
		if l.peekOk(i+1) && l.peekN(i) == '*' && l.peekN(i+1) == '/' {
			l.skipN(i + 2)
			return nil
		}
		i++
	}
	l.skipN(i)
	return errors.New("unclosed multi-line comment")
}

func (l *Lexer) consumeString() error {
	i := 1
	endChar := byte('\'')
	for l.peekOk(i) {
		c := l.peekN(i)
		// backslash escape
		if c == '\\' {
			i++
			if l.peekOk(i) {
				i++
			}
			continue
		}
		// single quote
		if c == endChar {
			// double single quote ''
			if l.peekOk(i+1) && l.peekN(i+1) == endChar {
				i += 2
				continue
			}
			break
		}
		i++
	}
	if !l.peekOk(i) || l.peekN(i) != endChar {
		return errors.New("invalid string")
	}
	l.currentToken = &Token{
		Kind:   TokenKindString,
		String: l.slice(1, i),
		Pos:    Pos(l.offset + 1),
		End:    Pos(l.offset + i),
	}
	l.skipN(i + 1)
	return nil
}

func (l *Lexer) skipComments() error {
	for !l.isEOF() {
		l.skipSpace()
		if !l.peekOk(0) {
			return nil
		}
		switch l.peekN(0) {
		case '-':
			if l.peekOk(1) && l.peekN(1) == '-' {
				l.consumeSingleLineComment()
				continue
			}
			return nil
		case '/': // multi-line comment
			if l.peekOk(1) && l.peekN(1) == '*' {
				if err := l.consumeMultiLineComment(); err != nil {
					return err
				}
				continue
			}
			return nil
		case '\r', '\n':
			// skip \r\n or \n\r
			l.skipN(1)
		default:
			return nil
		}
	}
	return nil
}

func (l *Lexer) peekToken() (*Token, error) {
	savedState := l.saveState()
	if err := l.consumeToken(); err != nil {
		return nil, err
	}
	token := l.currentToken

	l.restoreState(savedState)
	return token, nil
}

func (l *Lexer) hasPrecedenceToken(last *Token) bool {
	return last != nil && (last.Kind == TokenKindIdent ||
		last.Kind == TokenKindKeyword ||
		last.Kind == TokenKindInt ||
		last.Kind == TokenKindFloat ||
		last.Kind == TokenKindString)
}

func (l *Lexer) consumeToken() error {
	// replace the current token; keep the previous one to disambiguate unary +/-
	prevToken := l.currentToken
	l.currentToken = nil
	if err := l.skipComments(); err != nil {
		return err
	}
	l.skipSpace()
	if l.isEOF() {
		return nil
	}
	switch l.peekN(0) {
	case '>', '<', '!', '=', '|':
		if l.peekN(0) == '|' && l.peekOk(1) && l.peekN(1) == '|' || // ||
			l.peekN(0) == '<' && l.peekOk(1) && l.peekN(1) == '>' || // <>
			l.peekN(0) == '=' && l.peekOk(1) && l.peekN(1) == '=' || // ==
			l.peekN(0) != '|' && l.peekOk(1) && l.peekN(1) == '=' { // |=
			l.currentToken = &Token{
				String: l.slice(0, 2),
				Kind:   TokenKind(l.slice(0, 2)),
				Pos:    Pos(l.offset),
				End:    Pos(l.offset + 2),
			}
			l.skipN(2)
			return nil
		}

	case '+', '-':
		// hasPrecedenceToken is used to distinguish between unary and binary operators
		if !l.hasPrecedenceToken(prevToken) && l.peekOk(1) && IsDigit(l.peekN(1)) {
			return l.consumeNumber()
		} else if l.peekOk(1) && l.peekN(1) == '>' {
			l.currentToken = &Token{
				String: l.slice(0, 2),
				Kind:   TokenKindArrow,
				Pos:    Pos(l.offset),
				End:    Pos(l.offset + 2),
			}
			l.skipN(2)
			return nil
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.consumeNumber()
	case '`', '$', '"':
		return l.consumeIdent(Pos(l.offset))
	case '\'':
		return l.consumeString()
	case ':':
		if l.peekOk(1) && l.peekN(1) == ':' {
			l.currentToken = &Token{
				String: l.slice(0, 2),
				Kind:   TokenKindDash,
				Pos:    Pos(l.offset),
				End:    Pos(l.offset + 2),
			}
			l.skipN(2)
			return nil
		}
	case '.':
		l.currentToken = &Token{
			String: l.slice(0, 1),
			Kind:   TokenKindDot,
			Pos:    Pos(l.offset),
			End:    Pos(l.offset + 1),
		}
		l.skipN(1)
		return nil
	}

	if IsIdentStart(l.peekN(0)) {
		return l.consumeIdent(Pos(l.offset))
	}

	// Non-ASCII bytes can only appear inside quoted identifiers or string
	// literals, which are handled above. Report the whole rune instead of
	// emitting a one-byte token that splits the UTF-8 sequence.
	if l.peekN(0) >= utf8.RuneSelf {
		r, _ := utf8.DecodeRuneInString(l.input[l.offset:])
		return fmt.Errorf("unexpected character %q", r)
	}

	token := &Token{}
	token.Pos = Pos(l.offset)
	token.End = Pos(l.offset + 1)
	token.String = l.input[l.offset : l.offset+1]
	token.Kind = TokenKind(token.String)
	l.skipN(1)
	l.currentToken = token
	return nil
}

func (l *Lexer) isEOF() bool {
	return l.offset >= len(l.input)
}

func (l *Lexer) skipSpace() {
	for !l.isEOF() {
		r, size := utf8.DecodeRuneInString(l.input[l.offset:])
		if !unicode.IsSpace(r) {
			break
		}
		l.offset += size
	}
}
