package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	TokenEOF     TokenKind = "<eof>"
	TokenIdent   TokenKind = "<ident>"
	TokenKeyword TokenKind = "<keyword>"
	TokenInt     TokenKind = "<int>"
	TokenFloat   TokenKind = "<float>"
	TokenString  TokenKind = "<string>"
	TokenDot               = "."
)

const (
	Unquoted = iota + 1
	DoubleQuote
	BackTicks
)

type Pos int
type TokenKind string

type Token struct {
	Pos Pos
	End Pos

	Kind      TokenKind
	String    string
	Base      int // 10 or 16 on TokenInt
	QuoteType int
}

func (t *Token) ToString() string {
	if t.Kind == TokenKeyword {
		return strings.ToUpper(t.String)
	}
	return t.String
}

type Lexer struct {
	input     string
	current   int
	lastToken *Token
}

func NewLexer(buf string) *Lexer {
	return &Lexer{input: buf}
}

func (l *Lexer) skipN(n int) {
	l.current += n
}

func (l *Lexer) slice(i, j int) string {
	return l.input[l.current+i : l.current+j]
}

func (l *Lexer) peekN(n int) byte {
	return l.input[l.current+n]
}

func (l *Lexer) peekOk(n int) bool {
	return l.current+n < len(l.input)
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
	if l.peekN(0) == '0' && l.peekOk(1) && l.peekN(1) == 'x' {
		i += 2
		base = 16
	}

	hasExp := false
	tokenKind := TokenInt
	hasNumberPart := false
	for l.peekOk(i) {
		hasNumberPart = true
		c := l.peekN(i)
		switch {
		case base == 10 && IsDigit(c):
			i++
			continue
		case base == 16 && IsHexDigit(c):
			i++
			continue
		case c == '.': // float
			tokenKind = TokenFloat
			i++
			continue
		case base != 16 && (c == 'e' || c == 'E' || c == 'p' || c == 'P'):
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
			continue
		}
		break
	}
	if (l.peekOk(i) && IsIdentPart(l.peekN(i))) || !hasNumberPart {
		return errors.New("invalid number")
	}
	l.lastToken = &Token{
		Kind:   tokenKind,
		String: l.slice(0, i),
		Pos:    Pos(l.current),
		End:    Pos(l.current + i),
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
		if l.peekN(i) == '$' {
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
		token.Kind = TokenKeyword
	} else {
		token.Kind = TokenIdent
	}
	token.Pos = Pos(l.current)
	token.End = Pos(l.current + i)
	token.String = slice
	token.QuoteType = quoteType
	l.lastToken = token

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
	l.skipN(i + 1)
}

func (l *Lexer) consumeMultiLineComment() {
	l.skipN(2)
	i := 0
	for !l.isEOF() {
		if l.peekOk(i+1) && l.peekN(i) == '*' && l.peekN(i+1) == '/' {
			i += 2
			break
		}
		i++
	}
	l.skipN(i)
}

func (l *Lexer) consumeString() error {
	i := 1
	endChar := byte('\'')
	for l.peekOk(i) && l.peekN(i) != endChar {
		i++
	}
	if !l.peekOk(i) {
		return errors.New("invalid string")
	}
	l.lastToken = &Token{
		Kind:   TokenString,
		String: l.slice(1, i),
		Pos:    Pos(l.current + 1),
		End:    Pos(l.current + i),
	}
	l.skipN(i + 1)
	return nil
}

func (l *Lexer) skipComments() {
	for !l.isEOF() {
		l.skipSpace()
		if !l.peekOk(0) {
			return
		}
		switch l.peekN(0) {
		case '-':
			if l.peekOk(1) && l.peekN(1) == '-' {
				l.consumeSingleLineComment()
				continue
			}
			return
		case '/': // multi-line comment
			if l.peekOk(1) && l.peekN(1) == '*' {
				l.consumeMultiLineComment()
				continue
			}
			return
		case '\r', '\n':
			// skip \r\n or \n\r
			l.skipN(1)
		default:
			return
		}
	}
}

func (l *Lexer) peekToken() (*Token, error) {
	saveToken := l.lastToken
	saveCurrent := l.current
	if err := l.consumeToken(); err != nil {
		return nil, err
	}
	token := l.lastToken

	l.lastToken = saveToken
	l.current = saveCurrent
	return token, nil
}

func (l *Lexer) hasPrecedenceToken(last *Token) bool {
	return last != nil && (last.Kind == TokenIdent ||
		last.Kind == TokenKeyword ||
		last.Kind == TokenInt ||
		last.Kind == TokenFloat ||
		last.Kind == TokenString)
}

func (l *Lexer) consumeToken() error {
	// clear last token
	lastToken := l.lastToken
	l.lastToken = nil
	l.skipComments()
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
			l.lastToken = &Token{
				String: l.slice(0, 2),
				Kind:   TokenKind(l.slice(0, 2)),
				Pos:    Pos(l.current),
				End:    Pos(l.current + 2),
			}
			l.skipN(2)
			return nil
		}

	case '+', '-':
		// hasPrecedenceToken is used to distinguish between unary and binary operators
		if !l.hasPrecedenceToken(lastToken) && l.peekOk(1) && IsDigit(l.peekN(1)) {
			return l.consumeNumber()
		} else if l.peekOk(1) && l.peekN(1) == '>' {
			l.lastToken = &Token{
				String: l.slice(0, 2),
				Kind:   opTypeArrow,
				Pos:    Pos(l.current),
				End:    Pos(l.current + 2),
			}
			l.skipN(2)
			return nil
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.consumeNumber()
	case '`', '$', '"':
		return l.consumeIdent(Pos(l.current))
	case '\'':
		return l.consumeString()
	case ':':
		if l.peekOk(1) && l.peekN(1) == ':' {
			l.lastToken = &Token{
				String: l.slice(0, 2),
				Kind:   opTypeCast,
				Pos:    Pos(l.current),
				End:    Pos(l.current + 2),
			}
			l.skipN(2)
			return nil
		}
	case '.':
		l.lastToken = &Token{
			String: l.slice(0, 1),
			Kind:   TokenDot,
			Pos:    Pos(l.current),
			End:    Pos(l.current + 1),
		}
		l.skipN(1)
		return nil
	}

	if IsIdentStart(l.peekN(0)) {
		return l.consumeIdent(Pos(l.current))
	}

	token := &Token{}
	token.Pos = Pos(l.current)
	token.End = Pos(l.current + 1)
	token.String = l.input[l.current : l.current+1]
	token.Kind = TokenKind(token.String)
	l.skipN(1)
	l.lastToken = token
	return nil
}

func (l *Lexer) isEOF() bool {
	return l.current >= len(l.input)
}

func (l *Lexer) skipSpace() {
	for !l.isEOF() {
		r, size := utf8.DecodeRuneInString(l.input[l.current:])
		if !unicode.IsSpace(r) {
			break
		}
		l.current += size
	}
}
