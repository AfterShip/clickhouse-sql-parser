package parser

import (
	"errors"
	"fmt"
	"strings"
)

type Parser struct {
	lexer *Lexer
	lines lineStarts // lazily built on the first error, for position lookup
}

// lineStarts returns the line-start offsets for the input, building them on
// first use. Errors are the cold path, so we avoid paying for this on success.
func (p *Parser) lineStarts() lineStarts {
	if p.lines == nil {
		p.lines = newLineStarts(p.lexer.input)
	}
	return p.lines
}

func NewParser(buffer string) *Parser {
	return &Parser{
		lexer: NewLexer(buffer),
	}
}

func (p *Parser) curTokenKind() TokenKind {
	if p.cur() == nil {
		return TokenKindEOF
	}
	return p.cur().Kind
}

// curTokenString returns the string of the current token, or "<EOF>" when the
// lexer has no current token. Use this on error paths where the current token
// may be nil (e.g. at end of input) to avoid a nil-pointer dereference.
func (p *Parser) curTokenString() string {
	if p.cur() == nil {
		return "<EOF>"
	}
	return p.cur().String
}

// cur returns the current lookahead token: the token the parser is looking at
// but has not consumed yet. It is nil at end of input.
func (p *Parser) cur() *Token {
	return p.lexer.current
}

func (p *Parser) End() Pos {
	if p.cur() == nil {
		return Pos(p.lexer.offset + 1)
	}
	return p.cur().End
}

func (p *Parser) Pos() Pos {
	last := p.cur()
	if last == nil {
		return Pos(p.lexer.offset)
	}
	return last.Pos
}

func (p *Parser) matchTokenKind(kind TokenKind) bool {
	return p.curTokenKind() == kind ||
		(kind == TokenKindIdent && p.curTokenKind() == TokenKindKeyword)
}

// expectTokenKind consumes the current token if it is the given kind.
func (p *Parser) expectTokenKind(kind TokenKind) error {
	if curToken := p.tryConsumeTokenKind(kind); curToken != nil {
		return nil
	}
	return &ParseError{
		Pos:      p.Pos(),
		Got:      p.cur(),
		Expected: []TokenKind{kind},
	}
}

func (p *Parser) tryConsumeTokenKind(kind TokenKind) *Token {
	if p.matchTokenKind(kind) {
		curToken := p.cur()
		_ = p.lexer.consumeToken()
		return curToken
	}
	return nil
}

func (p *Parser) matchKeyword(keyword string) bool {
	return p.matchTokenKind(TokenKindKeyword) && strings.EqualFold(p.cur().String, keyword)
}

func (p *Parser) matchOneOfKeywords(keywords ...string) bool {
	for _, keyword := range keywords {
		if p.matchKeyword(keyword) {
			return true
		}
	}
	return false
}

func (p *Parser) expectKeyword(keyword string) error {
	if !p.matchKeyword(keyword) {
		return &ParseError{
			Pos:     p.Pos(),
			Got:     p.cur(),
			Keyword: keyword,
		}
	}
	_ = p.lexer.consumeToken()
	return nil
}

func (p *Parser) tryConsumeKeywords(keywords ...string) bool {
	savedState := p.lexer.saveState()
	for _, keyword := range keywords {
		if !p.matchKeyword(keyword) {
			p.lexer.restoreState(savedState)
			return false
		}
		_ = p.lexer.consumeToken()
	}
	return true
}

func (p *Parser) tryParseIdent() *Ident {
	if p.curTokenKind() != TokenKindIdent {
		return nil
	}
	curToken := p.cur()
	_ = p.lexer.consumeToken()
	return &Ident{
		NamePos:   curToken.Pos,
		NameEnd:   curToken.End,
		Name:      curToken.String,
		QuoteType: curToken.QuoteType,
	}
}

func (p *Parser) parseIdent() (*Ident, error) {
	curToken := p.cur()
	if err := p.expectTokenKind(TokenKindIdent); err != nil {
		return nil, err
	}
	ident := &Ident{
		NamePos:   curToken.Pos,
		NameEnd:   curToken.End,
		Name:      curToken.String,
		QuoteType: curToken.QuoteType,
	}
	return ident, nil
}

func (p *Parser) parseIdentOrStar() (*Ident, error) {
	switch {
	case p.matchTokenKind(TokenKindIdent):
		return p.parseIdent()
	case p.matchTokenKind("*"):
		curToken := p.cur()
		_ = p.lexer.consumeToken()
		return &Ident{
			NamePos: curToken.Pos,
			NameEnd: curToken.End,
			Name:    curToken.String,
		}, nil
	default:
		return nil, fmt.Errorf("expected <ident> or *, but got %q", p.curTokenKind())
	}
}

func (p *Parser) parseIdentOrString() (*Ident, error) {
	switch {
	case p.matchTokenKind(TokenKindIdent):
		return p.parseIdent()
	case p.matchTokenKind(TokenKindString):
		curToken := p.cur()
		_ = p.lexer.consumeToken()
		return &Ident{
			NamePos:   curToken.Pos,
			NameEnd:   curToken.End,
			Name:      curToken.String,
			QuoteType: SingleQuote, // Treat string literals as single-quoted identifiers
		}, nil
	default:
		return nil, fmt.Errorf("expected <ident> or <string>, but got %q", p.curTokenKind())
	}
}

func (p *Parser) tryParseDotIdent(_ Pos) (*Ident, error) {
	if p.tryConsumeTokenKind(TokenKindDot) == nil {
		return nil, nil // nolint
	}
	return p.parseIdent()
}

func (p *Parser) tryParseDotIdentOrString(_ Pos) (*Ident, error) {
	if p.tryConsumeTokenKind(TokenKindDot) == nil {
		return nil, nil // nolint
	}
	return p.parseIdentOrString()
}

func (p *Parser) parseUUID() (*UUID, error) {
	if err := p.expectKeyword(KeywordUuid); err != nil {
		return nil, err
	}

	uuidString, err := p.parseString(p.Pos())
	if err != nil {
		return nil, err
	}
	return &UUID{
		Value: uuidString,
	}, nil
}

func (p *Parser) tryParseUUID() (*UUID, error) {
	if !p.matchKeyword(KeywordUuid) {
		return nil, nil // nolint
	}
	return p.parseUUID()
}

func (p *Parser) tryParseComment() (*StringLiteral, error) {
	if !p.tryConsumeKeywords(KeywordComment) {
		return nil, nil
	}
	return p.parseString(p.Pos())
}

func (p *Parser) tryParseIfExists() (bool, error) {
	if !p.tryConsumeKeywords(KeywordIf) {
		return false, nil
	}

	if err := p.expectKeyword(KeywordExists); err != nil {
		return false, err
	}
	return true, nil
}

func (p *Parser) tryParseIfNotExists() (bool, error) {
	if !p.tryConsumeKeywords(KeywordIf) {
		return false, nil
	}

	if err := p.expectKeyword(KeywordNot); err != nil {
		return false, err
	}

	if err := p.expectKeyword(KeywordExists); err != nil {
		return false, err
	}
	return true, nil
}

func (p *Parser) tryParseNull(pos Pos) *NullLiteral {
	if !p.tryConsumeKeywords(KeywordNull) {
		return nil
	}
	return &NullLiteral{NullPos: pos}
}

func (p *Parser) tryParseNotNull(pos Pos) (*NotNullLiteral, error) {
	if !p.tryConsumeKeywords(KeywordNot) {
		return nil, nil // nolint
	}
	notNull := &NotNullLiteral{NotPos: pos}

	nullPos := p.Pos()
	if err := p.expectKeyword(KeywordNull); err != nil {
		return notNull, err
	}
	notNull.NullLiteral = &NullLiteral{NullPos: nullPos}
	return notNull, nil
}

func (p *Parser) parseDecimal(pos Pos) (*NumberLiteral, error) {
	number, err := p.parseNumber(pos)
	if err != nil {
		return nil, err
	}
	if number.Base != 10 {
		return nil, fmt.Errorf("invalid decimal literal: %q", number.Literal)
	}
	return number, nil
}

func (p *Parser) parseNumber(pos Pos) (*NumberLiteral, error) {
	var err error

	curToken := p.cur()
	switch {
	case p.matchTokenKind(TokenKindInt):
		err = p.expectTokenKind(TokenKindInt)
	case p.matchTokenKind(TokenKindFloat):
		err = p.expectTokenKind(TokenKindFloat)
	case p.matchTokenKind(TokenKindDot):
		_ = p.lexer.consumeToken()
		curToken = p.cur()
		if err := p.expectTokenKind(TokenKindInt); err != nil {
			return nil, err
		}
		if curToken.Base != 10 {
			return nil, fmt.Errorf("invalid decimal literal: %q", curToken.String)
		}
		curToken.String = "." + curToken.String
		curToken.Kind = TokenKindFloat
	default:
		return nil, fmt.Errorf("expected <int> or <float>, but got %q", p.curTokenKind())
	}
	if err != nil {
		return nil, err
	}
	number := &NumberLiteral{
		NumPos:  pos,
		NumEnd:  curToken.End,
		Literal: curToken.String,
		Base:    curToken.Base,
	}
	return number, nil
}

func (p *Parser) parseString(pos Pos) (*StringLiteral, error) {
	curToken := p.cur()
	if err := p.expectTokenKind(TokenKindString); err != nil {
		return nil, err
	}

	str := &StringLiteral{
		LiteralPos: pos,
		LiteralEnd: curToken.End,
		Literal:    curToken.String,
	}
	return str, nil
}

func (p *Parser) parseLiteral(pos Pos) (Literal, error) {
	switch {
	case p.matchTokenKind(TokenKindInt), p.matchTokenKind(TokenKindFloat):
		return p.parseNumber(pos)
	case p.matchTokenKind(TokenKindString):
		return p.parseString(pos)
	case p.matchTokenKind(TokenKindIdent):
		return p.parseIdent()
	case p.matchKeyword(KeywordNull):
		// accept the NULL keyword
		return &NullLiteral{NullPos: pos}, nil
	default:
		return nil, fmt.Errorf("expected <int>, <string>, <ident> or keyword <NULL>, but got %q", p.curTokenKind())
	}
}

func (p *Parser) ParseNestedIdentifier(pos Pos) (*NestedIdentifier, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	dotIdent, err := p.tryParseDotIdent(p.Pos())
	if err != nil {
		return nil, err
	}
	if dotIdent != nil {
		return &NestedIdentifier{
			Ident:    ident,
			DotIdent: dotIdent,
		}, nil
	}
	return &NestedIdentifier{
		Ident: ident,
	}, nil
}

func (p *Parser) tryParseFormat(pos Pos) (*FormatClause, error) {
	if !p.matchKeyword(KeywordFormat) {
		return nil, nil // nolint
	}
	return p.parseFormat(pos)
}

func (p *Parser) parseFormat(pos Pos) (*FormatClause, error) {
	if err := p.expectKeyword(KeywordFormat); err != nil {
		return nil, err
	}
	formatIdent, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &FormatClause{
		FormatPos: pos,
		Format:    formatIdent,
	}, nil
}

// wrapError finalizes a parse error: it ensures the error is a *ParseError with
// line/column resolved and the input attached for caret rendering. Errors that
// already originate as *ParseError (from the expect* helpers) keep their
// captured position and expected-token information; the long tail of
// fmt.Errorf sites is wrapped here with the current position.
func (p *Parser) wrapError(err error) error {
	if err == nil {
		return nil
	}

	var pe *ParseError
	if !errors.As(err, &pe) {
		pe = &ParseError{
			Pos: p.Pos(),
			Got: p.cur(),
			Msg: err.Error(),
		}
	}
	if pe.Line == 0 {
		pe.Line, pe.Column = p.lineStarts().position(int(pe.Pos))
	}
	pe.input = p.lexer.input
	pe.starts = p.lineStarts()
	return pe
}

func (p *Parser) parseRatioExpr(pos Pos) (*RatioExpr, error) {
	numerator, err := p.parseNumber(pos)
	if err != nil {
		return nil, err
	}

	var denominator *NumberLiteral
	if p.tryConsumeTokenKind(TokenKindDiv) != nil {
		denominator, err = p.parseNumber(pos)
		if err != nil {
			return nil, err
		}
	}
	return &RatioExpr{
		Numerator:   numerator,
		Denominator: denominator,
	}, nil
}
