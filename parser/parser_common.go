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

func (p *Parser) lastTokenKind() TokenKind {
	if p.last() == nil {
		return TokenKindEOF
	}
	return p.last().Kind
}

// lastTokenString returns the string of the last token, or "<EOF>" when the
// lexer has no current token. Use this on error paths where the last token may
// be nil (e.g. at end of input) to avoid a nil-pointer dereference.
func (p *Parser) lastTokenString() string {
	if p.last() == nil {
		return "<EOF>"
	}
	return p.last().String
}

func (p *Parser) last() *Token {
	return p.lexer.lastToken
}

func (p *Parser) End() Pos {
	if p.last() == nil {
		return Pos(p.lexer.current + 1)
	}
	return p.last().End
}

func (p *Parser) Pos() Pos {
	last := p.last()
	if last == nil {
		return Pos(p.lexer.current)
	}
	return last.Pos
}

// matchTokenKind reports whether the current token matches the given kind.
// A non-reserved keyword also matches TokenKindIdent: most ClickHouse keywords
// (DATE, KEY, FIRST, ...) are valid identifiers anywhere an identifier is
// expected. Reserved keywords (see reservedKeywords) do not match, so a
// missing identifier before e.g. FROM or WHERE fails fast instead of silently
// swallowing the clause keyword as a name.
func (p *Parser) matchTokenKind(kind TokenKind) bool {
	if p.lastTokenKind() == kind {
		return true
	}
	return kind == TokenKindIdent &&
		p.lastTokenKind() == TokenKindKeyword &&
		!reservedKeywords.Contains(strings.ToUpper(p.last().String))
}

// expectTokenKind consumes the last token if it is the given kind.
func (p *Parser) expectTokenKind(kind TokenKind) error {
	if lastToken := p.tryConsumeTokenKind(kind); lastToken != nil {
		return nil
	}
	return &ParseError{
		Pos:      p.Pos(),
		Got:      p.last(),
		Expected: []TokenKind{kind},
	}
}

func (p *Parser) tryConsumeTokenKind(kind TokenKind) *Token {
	if p.matchTokenKind(kind) {
		lastToken := p.last()
		_ = p.lexer.consumeToken()
		return lastToken
	}
	return nil
}

func (p *Parser) matchKeyword(keyword string) bool {
	return p.matchTokenKind(TokenKindKeyword) && strings.EqualFold(p.last().String, keyword)
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
			Got:     p.last(),
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
	if p.lastTokenKind() != TokenKindIdent {
		return nil
	}
	lastToken := p.last()
	_ = p.lexer.consumeToken()
	return &Ident{
		NamePos:   lastToken.Pos,
		NameEnd:   lastToken.End,
		Name:      lastToken.String,
		QuoteType: lastToken.QuoteType,
	}
}

// parseIdentAnyKeyword parses the current token as an identifier, accepting
// any keyword token — reserved or not — as the name. Use it only in positions
// where context has already proven the token is a name and not the start of a
// clause or expression: after AS, after a dot in a qualified name, or a select
// item the lookahead disambiguated.
func (p *Parser) parseIdentAnyKeyword() (*Ident, error) {
	last := p.last()
	if last == nil || (last.Kind != TokenKindIdent && last.Kind != TokenKindKeyword) {
		return nil, &ParseError{
			Pos:      p.Pos(),
			Got:      last,
			Expected: []TokenKind{TokenKindIdent},
		}
	}
	_ = p.lexer.consumeToken()
	return &Ident{
		NamePos:   last.Pos,
		NameEnd:   last.End,
		Name:      last.String,
		QuoteType: last.QuoteType,
	}, nil
}

func (p *Parser) parseIdent() (*Ident, error) {
	lastToken := p.last()
	if err := p.expectTokenKind(TokenKindIdent); err != nil {
		return nil, err
	}
	ident := &Ident{
		NamePos:   lastToken.Pos,
		NameEnd:   lastToken.End,
		Name:      lastToken.String,
		QuoteType: lastToken.QuoteType,
	}
	return ident, nil
}

func (p *Parser) parseIdentOrStar() (*Ident, error) {
	switch {
	case p.matchTokenKind(TokenKindIdent):
		return p.parseIdent()
	case p.matchTokenKind("*"):
		lastToken := p.last()
		_ = p.lexer.consumeToken()
		return &Ident{
			NamePos: lastToken.Pos,
			NameEnd: lastToken.End,
			Name:    lastToken.String,
		}, nil
	default:
		return nil, fmt.Errorf("expected <ident> or *, but got %q", p.lastTokenKind())
	}
}

func (p *Parser) parseIdentOrString() (*Ident, error) {
	switch {
	case p.matchTokenKind(TokenKindIdent):
		return p.parseIdent()
	case p.matchTokenKind(TokenKindString):
		lastToken := p.last()
		_ = p.lexer.consumeToken()
		return &Ident{
			NamePos:   lastToken.Pos,
			NameEnd:   lastToken.End,
			Name:      lastToken.String,
			QuoteType: SingleQuote, // Treat string literals as single-quoted identifiers
		}, nil
	default:
		return nil, fmt.Errorf("expected <ident> or <string>, but got %q", p.lastTokenKind())
	}
}

func (p *Parser) tryParseDotIdent(_ Pos) (*Ident, error) {
	if p.tryConsumeTokenKind(TokenKindDot) == nil {
		return nil, nil // nolint
	}
	// After a dot the token can only be a member name, so even reserved
	// keywords are accepted (e.g. `db.from`, `t.limit`).
	return p.parseIdentAnyKeyword()
}

func (p *Parser) tryParseDotIdentOrString(_ Pos) (*Ident, error) {
	if p.tryConsumeTokenKind(TokenKindDot) == nil {
		return nil, nil // nolint
	}
	// After a dot the token can only be a member name, so even reserved
	// keywords are accepted (e.g. `db.from`).
	if p.lastTokenKind() == TokenKindKeyword {
		return p.parseIdentAnyKeyword()
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

	lastToken := p.last()
	switch {
	case p.matchTokenKind(TokenKindInt):
		err = p.expectTokenKind(TokenKindInt)
	case p.matchTokenKind(TokenKindFloat):
		err = p.expectTokenKind(TokenKindFloat)
	case p.matchTokenKind(TokenKindDot):
		_ = p.lexer.consumeToken()
		lastToken = p.last()
		if err := p.expectTokenKind(TokenKindInt); err != nil {
			return nil, err
		}
		if lastToken.Base != 10 {
			return nil, fmt.Errorf("invalid decimal literal: %q", lastToken.String)
		}
		lastToken.String = "." + lastToken.String
		lastToken.Kind = TokenKindFloat
	default:
		return nil, fmt.Errorf("expected <int> or <float>, but got %q", p.lastTokenKind())
	}
	if err != nil {
		return nil, err
	}
	number := &NumberLiteral{
		NumPos:  pos,
		NumEnd:  lastToken.End,
		Literal: lastToken.String,
		Base:    lastToken.Base,
	}
	return number, nil
}

func (p *Parser) parseString(pos Pos) (*StringLiteral, error) {
	lastToken := p.last()
	if err := p.expectTokenKind(TokenKindString); err != nil {
		return nil, err
	}

	str := &StringLiteral{
		LiteralPos: pos,
		LiteralEnd: lastToken.End,
		Literal:    lastToken.String,
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
		return nil, fmt.Errorf("expected <int>, <string>, <ident> or keyword <NULL>, but got %q", p.lastTokenKind())
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
			Got: p.last(),
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
