package parser

import (
	"fmt"
	"strings"
)

const (
	PrecedenceUnknown = iota
	PrecedenceOr
	PrecedenceAnd
	PrecedenceQuery
	PrecedenceNot
	PrecedenceGlobal
	PrecedenceIs
	PrecedenceCompare
	PrecedenceBetweenLike
	precedenceIn
	PrecedenceAddSub
	PrecedenceMulDivMod
	PrecedenceBracket
	PrecedenceArrow
	PrecedenceDot
	PrecedenceDoubleColon
)

func (p *Parser) tryParseColumnComment(pos Pos) (*StringLiteral, error) {
	if p.tryConsumeKeyword(KeywordComment) == nil {
		return nil, nil // nolint
	}
	return p.parseString(pos)
}

func (p *Parser) getNextPrecedence() int {
	switch {
	case p.matchKeyword(KeywordOr):
		return PrecedenceOr
	case p.matchKeyword(KeywordAnd):
		return PrecedenceAnd
	case p.matchKeyword(KeywordIs):
		return PrecedenceIs
	case p.matchKeyword(KeywordNot):
		return PrecedenceNot
	case p.matchTokenKind(TokenKindDot):
		return PrecedenceDot
	case p.matchTokenKind(TokenKindDash):
		return PrecedenceDoubleColon
	case p.matchTokenKind(TokenKindSingleEQ), p.matchTokenKind(TokenKindLT), p.matchTokenKind(TokenKindLE),
		p.matchTokenKind(TokenKindGE), p.matchTokenKind(TokenKindGT), p.matchTokenKind(TokenKindDoubleEQ),
		p.matchTokenKind(TokenKindNE), p.matchTokenKind("<>"):
		return PrecedenceCompare
	case p.matchTokenKind(TokenKindPlus), p.matchTokenKind(TokenKindMinus):
		return PrecedenceAddSub
	case p.matchTokenKind(TokenKindMul), p.matchTokenKind(TokenKindDiv), p.matchTokenKind(TokenKindMod):
		return PrecedenceMulDivMod
	case p.matchTokenKind(TokenKindArrow):
		return PrecedenceArrow
	case p.matchTokenKind(TokenKindLParen), p.matchTokenKind(TokenKindLBracket):
		return PrecedenceBracket
	case p.matchTokenKind(TokenKindDash):
		return PrecedenceDoubleColon
	case p.matchTokenKind(TokenKindDot):
		return PrecedenceDot
	case p.matchKeyword(KeywordBetween), p.matchKeyword(KeywordLike), p.matchKeyword(KeywordIlike):
		return PrecedenceBetweenLike
	case p.matchKeyword(KeywordIn):
		return precedenceIn
	case p.matchKeyword(KeywordGlobal):
		return PrecedenceGlobal
	case p.matchTokenKind(TokenKindQuestionMark):
		return PrecedenceQuery
	default:
		return PrecedenceUnknown
	}
}

func (p *Parser) parseInfix(expr Expr, precedence int) (Expr, error) {
	switch {
	case p.matchTokenKind(TokenKindSingleEQ), p.matchTokenKind(TokenKindLT), p.matchTokenKind(TokenKindLE),
		p.matchTokenKind(TokenKindGE), p.matchTokenKind(TokenKindGT),
		p.matchTokenKind(TokenKindNE), p.matchTokenKind("<>"),
		p.matchTokenKind(TokenKindMinus), p.matchTokenKind(TokenKindPlus), p.matchTokenKind(TokenKindMul),
		p.matchTokenKind(TokenKindDiv), p.matchTokenKind(TokenKindMod),
		p.matchKeyword(KeywordIn), p.matchKeyword(KeywordLike),
		p.matchKeyword(KeywordIlike), p.matchKeyword(KeywordAnd), p.matchKeyword(KeywordOr),
		p.matchTokenKind(TokenKindArrow), p.matchTokenKind(TokenKindDoubleEQ):
		op := p.last().ToString()
		_ = p.lexer.consumeToken()
		rightExpr, err := p.parseSubExpr(p.Pos(), precedence)
		if err != nil {
			return nil, err
		}
		return &BinaryOperation{
			LeftExpr:  expr,
			Operation: TokenKind(op),
			RightExpr: rightExpr,
		}, nil
	case p.matchTokenKind(TokenKindDash):
		_ = p.lexer.consumeToken()

		if p.matchTokenKind(TokenKindIdent) && p.last().String == "Tuple" {
			name, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			if _, err = p.consumeTokenKind(TokenKindLParen); err != nil {
				return nil, err
			}
			// it's a tuple type definition after "::" operator
			rightExpr, err := p.parseNestedType(name, p.Pos())
			if err != nil {
				return nil, err
			}
			return &BinaryOperation{
				LeftExpr:  expr,
				Operation: TokenKindDash,
				RightExpr: rightExpr,
			}, nil
		}

		rightExpr, err := p.parseSubExpr(p.Pos(), precedence)
		if err != nil {
			return nil, err
		}
		return &BinaryOperation{
			LeftExpr:  expr,
			Operation: TokenKindDash,
			RightExpr: rightExpr,
		}, nil
	case p.matchKeyword(KeywordBetween):
		return p.parseBetweenClause(expr)
	case p.matchKeyword(KeywordGlobal):
		_ = p.lexer.consumeToken()
		if p.consumeKeyword(KeywordIn) != nil {
			return nil, fmt.Errorf("expected IN after GLOBAL, got %s", p.lastTokenKind())
		}
		rightExpr, err := p.parseSubExpr(p.Pos(), precedence)
		if err != nil {
			return nil, err
		}
		return &BinaryOperation{
			LeftExpr:  expr,
			Operation: "GLOBAL IN",
			RightExpr: rightExpr,
		}, nil
	case p.matchTokenKind(TokenKindDot):
		_ = p.lexer.consumeToken()
		// access column with dot notation
		var rightExpr Expr
		var err error
		if p.matchTokenKind(TokenKindIdent) {
			rightExpr, err = p.parseIdent()
		} else {
			rightExpr, err = p.parseDecimal(p.Pos())
		}
		if err != nil {
			return nil, err
		}
		return &IndexOperation{
			Object:    expr,
			Operation: TokenKindDot,
			Index:     rightExpr,
		}, nil
	case p.matchKeyword(KeywordNot):
		_ = p.lexer.consumeToken()
		switch {
		case p.matchKeyword(KeywordIn):
		case p.matchKeyword(KeywordLike):
		case p.matchKeyword(KeywordIlike):
		default:
			return nil, fmt.Errorf("expected IN, LIKE or ILIKE after NOT, got %s", p.lastTokenKind())
		}
		if p.matchKeyword(KeywordBetween) {
			return p.parseBetweenClause(expr)
		}
		op := p.last().ToString()
		_ = p.lexer.consumeToken()
		rightExpr, err := p.parseSubExpr(p.Pos(), precedence)
		if err != nil {
			return nil, err
		}
		return &BinaryOperation{
			LeftExpr:  expr,
			Operation: TokenKind("NOT " + op),
			RightExpr: rightExpr,
		}, nil
	case p.matchTokenKind(TokenKindLBracket):
		params, err := p.parseArrayParams(p.Pos())
		if err != nil {
			return nil, err
		}
		return &ObjectParams{
			Object: expr,
			Params: params,
		}, nil
	case p.matchTokenKind(TokenKindQuestionMark):
		return p.parseTernaryExpr(expr)
	case p.matchKeyword(KeywordIs):
		_ = p.lexer.consumeToken()
		isNotNull := p.tryConsumeKeyword(KeywordNot) != nil
		if err := p.consumeKeyword(KeywordNull); err != nil {
			return nil, err
		}
		if isNotNull {
			return &IsNotNullExpr{
				IsPos: p.Pos(),
				Expr:  expr,
			}, nil
		}
		return &IsNullExpr{
			IsPos: p.Pos(),
			Expr:  expr,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected token kind: %s", p.lastTokenKind())
	}
}

func (p *Parser) parseExpr(pos Pos) (Expr, error) {
	return p.parseSubExpr(pos, PrecedenceUnknown)
}

func (p *Parser) parseSubExpr(pos Pos, precedence int) (Expr, error) {
	expr, err := p.parseUnaryExpr(pos)
	if err != nil {
		return nil, err
	}
	for !p.lexer.isEOF() {
		nextPrecedence := p.getNextPrecedence()
		if nextPrecedence <= precedence {
			return expr, nil
		}
		// parse binary operation
		expr, err = p.parseInfix(expr, nextPrecedence)
		if err != nil {
			return nil, err
		}
	}
	return expr, nil
}

func (p *Parser) parseTernaryExpr(condition Expr) (*TernaryOperation, error) {
	if _, err := p.consumeTokenKind(TokenKindQuestionMark); err != nil {
		return nil, err
	}
	trueExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if _, err := p.consumeTokenKind(TokenKindColon); err != nil {
		return nil, err
	}
	falseExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &TernaryOperation{
		Condition: condition,
		TrueExpr:  trueExpr,
		FalseExpr: falseExpr,
	}, nil
}

func (p *Parser) parseColumnExtractExpr(pos Pos) (*ExtractExpr, error) {
	if err := p.consumeKeyword(KeywordExtract); err != nil {
		return nil, err
	}
	if _, err := p.consumeTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}

	// parse interval
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if !intervalType.Contains(strings.ToUpper(ident.Name)) {
		return nil, fmt.Errorf("unknown interval type: <%q>", ident.Name)
	}

	fromPos := p.Pos()
	if err := p.consumeKeyword(KeywordFrom); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &ExtractExpr{
		ExtractPos: pos,
		Interval:   ident,
		FromPos:    fromPos,
		FromExpr:   expr,
	}, nil
}

func (p *Parser) parseUnaryExpr(pos Pos) (Expr, error) {
	kind := p.last()
	switch {
	case p.matchTokenKind(TokenKindPlus),
		p.matchTokenKind(TokenKindMinus),
		p.matchKeyword(KeywordNot):
		_ = p.lexer.consumeToken()
	default:
		return p.parseColumnExpr(pos)
	}

	expr, err := p.parseColumnExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	return &UnaryExpr{
		UnaryPos: pos,
		Kind:     kind.Kind,
		Expr:     expr,
	}, nil

}

func (p *Parser) peekTokenKind(kind TokenKind) bool {
	if p.lexer.isEOF() {
		return false
	}
	token, err := p.lexer.peekToken()
	if err != nil || token == nil {
		return false
	}
	return token.Kind == kind
}

func (p *Parser) parseColumnExpr(pos Pos) (Expr, error) { //nolint:funlen
	switch {
	case p.matchKeyword(KeywordInterval):
		return p.parseColumnExprInterval(pos)
	case p.matchKeyword(KeywordDate), p.matchKeyword(KeywordTimestamp):
		nextToken, err := p.lexer.peekToken()
		if err != nil {
			return nil, err
		}
		if nextToken != nil && nextToken.Kind == TokenKindString {
			return p.parseString(p.Pos())
		}
		return p.parseIdentOrFunction(pos)
	case p.matchKeyword(KeywordCast):
		return p.parseColumnCastExpr(pos)
	case p.matchKeyword(KeywordCase):
		return p.parseColumnCaseExpr(pos)
	case p.matchKeyword(KeywordExtract):
		return p.parseColumnExtractExpr(pos)
	case p.matchTokenKind(TokenKindIdent):
		return p.parseIdentOrFunction(pos)
	case p.matchTokenKind(TokenKindString): // string literal
		return p.parseString(pos)
	case p.matchTokenKind(TokenKindInt),
		p.matchTokenKind(TokenKindFloat): // number literal
		return p.parseNumber(pos)
	case p.matchTokenKind(TokenKindLParen):
		if peek, _ := p.lexer.peekToken(); peek != nil {
			if peek.Kind == TokenKindKeyword && strings.EqualFold(peek.String, KeywordSelect) {
				return p.parseSubQuery(pos)
			}
		}
		return p.parseFunctionParams(p.Pos())
	case p.matchTokenKind("*"):
		return p.parseColumnStar(p.Pos())
	case p.matchTokenKind(TokenKindLBracket):
		return p.parseArrayParams(p.Pos())
	case p.matchTokenKind(TokenKindLBrace):
		// The map literal string also starts with '{', so we need to check the next token
		// to determine if it is a map literal or a query param.
		if p.peekTokenKind(TokenKindIdent) {
			return p.parseQueryParam(p.Pos())
		}
		return p.parseMapLiteral(p.Pos())
	case p.matchTokenKind(TokenKindDot):
		return p.parseNumber(p.Pos())
	case p.matchTokenKind(TokenKindQuestionMark):
		// Placeholder `?`
		_ = p.lexer.consumeToken()
		return &PlaceHolder{
			PlaceholderPos: pos,
			PlaceHolderEnd: pos,
			Type:           string(TokenKindQuestionMark),
		}, nil
	default:
		return nil, fmt.Errorf("unexpected token kind: %s", p.lastTokenKind())
	}
}

func (p *Parser) parseColumnCastExpr(pos Pos) (Expr, error) {
	if err := p.consumeKeyword(KeywordCast); err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}

	columnExpr, err := p.parseColumnExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	var separator string
	asPos := p.Pos()
	switch {
	// CAST(x, T) and CAST(x AS T) are equivalent
	case p.matchKeyword(KeywordAs), p.matchTokenKind(","):
		separator = p.last().String
		_ = p.lexer.consumeToken()
	default:
		return nil, fmt.Errorf("expected AS or , but got %s", p.lastTokenKind())
	}

	var asColumnType Expr
	// CAST(1 AS 'Float') or CAST(1 AS Float) are equivalent
	if p.matchTokenKind(TokenKindString) {
		asColumnType, err = p.parseString(p.Pos())
	} else {
		asColumnType, err = p.parseColumnType(p.Pos())
	}
	if err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}

	return &CastExpr{
		CastPos:   pos,
		AsPos:     asPos,
		Separator: separator,
		Expr:      columnExpr,
		AsType:    asColumnType,
	}, nil
}

func (p *Parser) parseColumnExprListWithLParen(pos Pos) (*ColumnExprList, error) {
	return p.parseColumnExprListWithTerm(TokenKindRParen, pos)
}

func (p *Parser) parseColumnExprListWithSquareBracket(pos Pos) (*ColumnExprList, error) {
	return p.parseColumnExprListWithTerm(TokenKindRBracket, pos)
}

func (p *Parser) parseColumnExprList(pos Pos) (*ColumnExprList, error) {
	return p.parseColumnExprListWithTerm("", pos)
}

func (p *Parser) parseColumnExprListWithTerm(term TokenKind, pos Pos) (*ColumnExprList, error) {
	columnExprList := &ColumnExprList{
		ListPos: pos,
		ListEnd: pos,
	}
	columnExprList.HasDistinct = p.tryConsumeKeyword(KeywordDistinct) != nil
	columnList := make([]Expr, 0)
	for !p.lexer.isEOF() || p.last() != nil {
		if term != "" && p.matchTokenKind(term) {
			break
		}
		columnExpr, err := p.parseColumnsExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		if columnExpr == nil {
			break
		}
		columnList = append(columnList, columnExpr)
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}
	columnExprList.Items = columnList
	if len(columnList) > 0 {
		columnExprList.ListEnd = columnList[len(columnList)-1].End()
	}
	return columnExprList, nil
}

func (p *Parser) parseSelectItems() ([]*SelectItem, error) {
	selectItems := make([]*SelectItem, 0)
	for !p.lexer.isEOF() || p.last() != nil {
		selectItem, err := p.parseSelectItem()
		if err != nil {
			return nil, err
		}
		if selectItem == nil {
			break
		}
		selectItems = append(selectItems, selectItem)
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}
	return selectItems, nil
}

// Syntax: INTERVAL expr interval
func (p *Parser) parseColumnExprInterval(pos Pos) (Expr, error) {
	if err := p.consumeKeyword(KeywordInterval); err != nil {
		return nil, err
	}

	// store the column expr if it needs
	columnExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	// parse interval
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if !intervalType.Contains(strings.ToUpper(ident.Name)) {
		return nil, fmt.Errorf("unknown interval type: <%q>", ident.Name)
	}
	return &IntervalExpr{
		IntervalPos: pos,
		Expr:        columnExpr,
		Unit:        ident,
	}, nil
}

func (p *Parser) parseFunctionExpr(_ Pos) (*FunctionExpr, error) {
	// parse function name
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	// parse function params
	params, err := p.parseFunctionParams(p.Pos())
	if err != nil {
		return nil, err
	}
	return &FunctionExpr{
		Name:   name,
		Params: params,
	}, nil
}

func (p *Parser) parseColumnArgList(pos Pos) (*ColumnArgList, error) {
	if _, err := p.consumeTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}
	distinct := false
	if p.tryConsumeKeyword(KeywordDistinct) != nil {
		distinct = true
	}
	var items []Expr
	for !p.lexer.isEOF() && !p.matchTokenKind(TokenKindRParen) {
		item, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &ColumnArgList{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Distinct:      distinct,
		Items:         items,
	}, nil
}

func (p *Parser) parseFunctionParams(pos Pos) (*ParamExprList, error) {
	if _, err := p.consumeTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}
	params, err := p.parseColumnExprListWithLParen(p.Pos())
	if err != nil {
		return nil, err
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	paramExprList := &ParamExprList{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Items:         params,
	}

	// For some aggregate functions might support parametric arguments:
	// e.g. QUANTILE(0.5)(x) or QUANTILE(0.5, 0.9)(x).
	// So we need to have a check if there is another argument list with detecting the left bracket.
	if p.matchTokenKind(TokenKindLParen) {
		columnArgList, err := p.parseColumnArgList(p.Pos())
		if err != nil {
			return nil, err
		}
		paramExprList.ColumnArgList = columnArgList
	}
	return paramExprList, nil
}

func (p *Parser) parseMapLiteral(pos Pos) (*MapLiteral, error) {
	if _, err := p.consumeTokenKind(TokenKindLBrace); err != nil {
		return nil, err
	}

	keyValues := make([]KeyValue, 0)
	for !p.lexer.isEOF() && !p.matchTokenKind(TokenKindRBrace) {
		key, err := p.parseString(p.Pos())
		if err != nil {
			return nil, err
		}
		if _, err := p.consumeTokenKind(TokenKindColon); err != nil {
			return nil, err
		}
		value, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		keyValues = append(keyValues, KeyValue{
			Key:   *key,
			Value: value,
		})
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}
	rightBracePos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRBrace); err != nil {
		return nil, err
	}
	return &MapLiteral{
		LBracePos: pos,
		RBracePos: rightBracePos,
		KeyValues: keyValues,
	}, nil
}

func (p *Parser) parseQueryParam(pos Pos) (*QueryParam, error) {
	if _, err := p.consumeTokenKind(TokenKindLBrace); err != nil {
		return nil, err
	}

	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.consumeTokenKind(TokenKindColon); err != nil {
		return nil, err
	}
	columnType, err := p.parseColumnType(p.Pos())
	if err != nil {
		return nil, err
	}
	rightBracePos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRBrace); err != nil {
		return nil, err
	}
	return &QueryParam{
		LBracePos: pos,
		RBracePos: rightBracePos,
		Name:      ident,
		Type:      columnType,
	}, nil
}

func (p *Parser) parseArrayParams(pos Pos) (*ArrayParamList, error) {
	if _, err := p.consumeTokenKind(TokenKindLBracket); err != nil {
		return nil, err
	}
	params, err := p.parseColumnExprListWithSquareBracket(p.Pos())
	if err != nil {
		return nil, err
	}
	rightBracketPos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRBracket); err != nil {
		return nil, err
	}
	return &ArrayParamList{
		LeftBracketPos:  pos,
		RightBracketPos: rightBracketPos,
		Items:           params,
	}, nil
}

func (p *Parser) parseColumnsExpr(pos Pos) (*ColumnExpr, error) {
	expr, err := p.parseExpr(pos)
	if err != nil {
		return nil, err
	}

	var alias *Ident
	if p.tryConsumeKeyword(KeywordAs) != nil {
		alias, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
	}
	return &ColumnExpr{
		Expr:  expr,
		Alias: alias,
	}, nil
}

func (p *Parser) parseSelectItem() (*SelectItem, error) {
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	modifiers := make([]*FunctionExpr, 0)
	for {
		if p.matchKeyword(KeywordExcept) || p.matchKeyword(KeywordApply) || p.matchKeyword(KeywordReplace) {
			modifier, err := p.parseFunctionExpr(p.Pos())
			if err != nil {
				return nil, err
			}
			modifiers = append(modifiers, modifier)
		} else {
			break
		}
	}

	var alias *Ident
	if p.tryConsumeKeyword(KeywordAs) != nil {
		alias, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
	}
	return &SelectItem{
		Expr:      expr,
		Modifiers: modifiers,
		Alias:     alias,
	}, nil
}

func (p *Parser) parseColumnCaseExpr(pos Pos) (*CaseExpr, error) {
	// CASE expr
	caseExpr := &CaseExpr{CasePos: pos}
	if err := p.consumeKeyword(KeywordCase); err != nil {
		return nil, err
	}

	// case expr is optional
	if !p.matchKeyword(KeywordWhen) {
		expr, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		caseExpr.Expr = expr
	}

	// WHEN expr THEN expr
	whenClauses := make([]*WhenClause, 0)
	for p.matchKeyword(KeywordWhen) {
		whenPos := p.Pos()
		_ = p.lexer.consumeToken()
		whenCondition, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}

		thenPos := p.Pos()
		if err := p.consumeKeyword(KeywordThen); err != nil {
			return nil, err
		}
		thenCondition, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}

		whenClauses = append(whenClauses, &WhenClause{
			WhenPos: whenPos,
			ThenPos: thenPos,
			When:    whenCondition,
			Then:    thenCondition,
		})
	}
	caseExpr.Whens = whenClauses

	// ELSE expr
	if elseToken := p.tryConsumeKeyword(KeywordElse); elseToken != nil {
		elseExpr, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		caseExpr.ElsePos = elseToken.Pos
		caseExpr.Else = elseExpr
	}

	if err := p.consumeKeyword(KeywordEnd); err != nil {
		return nil, err
	}

	return caseExpr, nil
}

func (p *Parser) parseColumnType(_ Pos) (ColumnType, error) { // nolint:funlen
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if p.tryConsumeTokenKind(TokenKindLParen) != nil {
		switch {
		case p.matchTokenKind(TokenKindIdent):
			switch ident.Name {
			case "Nested":
				return p.parseNestedType(ident, p.Pos())
			case "JSON":
				return p.parseJSONType(ident, p.Pos())
			default:
				return p.parseComplexType(ident, p.Pos())
			}
		case p.matchTokenKind(TokenKindString):
			if peekToken, err := p.lexer.peekToken(); err == nil && peekToken.Kind == TokenKindSingleEQ {
				// enum values
				return p.parseEnumType(ident, p.Pos())
			}
			// like Datetime('Asia/Dubai')
			return p.parseColumnTypeWithParams(ident, p.Pos())
		case p.matchTokenKind(TokenKindInt), p.matchTokenKind(TokenKindFloat):
			// fixed size
			return p.parseColumnTypeWithParams(ident, p.Pos())
		default:
			return nil, fmt.Errorf("unexpected token kind: %v", p.lastTokenKind())
		}
	}
	return &ScalarType{Name: ident}, nil
}

func (p *Parser) parseColumnPropertyType(_ Pos) (Expr, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &PropertyType{
		Name: ident,
	}, nil
}

func (p *Parser) parseComplexType(name *Ident, pos Pos) (*ComplexType, error) {
	subTypes := make([]ColumnType, 0)
	for !p.lexer.isEOF() && !p.matchTokenKind(TokenKindRParen) {
		subExpr, err := p.parseColumnType(p.Pos())
		if err != nil {
			return nil, err
		}
		subTypes = append(subTypes, subExpr)
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &ComplexType{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Name:          name,
		Params:        subTypes,
	}, nil
}

func (p *Parser) parseEnumType(name *Ident, pos Pos) (*EnumType, error) {
	enumType := &EnumType{
		Name:    name,
		ListPos: pos,
		Values:  make([]EnumValue, 0),
	}
	for !p.lexer.isEOF() && !p.matchTokenKind(TokenKindRParen) {
		enumValue, err := p.parseEnumValueExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		if enumValue == nil {
			break
		}
		enumType.Values = append(enumType.Values, *enumValue)
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}
	if len(enumType.Values) > 0 {
		enumType.ListEnd = enumType.Values[len(enumType.Values)-1].Value.NumEnd
	}
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return enumType, nil
}

func (p *Parser) parseColumnTypeWithParams(name *Ident, pos Pos) (*TypeWithParams, error) {
	params := make([]Literal, 0)
	param, err := p.parseLiteral(p.Pos())
	if err != nil {
		return nil, err
	}
	params = append(params, param)
	for !p.lexer.isEOF() && p.tryConsumeTokenKind(TokenKindComma) != nil {
		size, err := p.parseLiteral(p.Pos())
		if err != nil {
			return nil, err
		}
		params = append(params, size)
	}

	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &TypeWithParams{
		Name:          name,
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Params:        params,
	}, nil
}

func (p *Parser) parseJSONPath() (*JSONPath, error) {
	idents := make([]*Ident, 0)
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	idents = append(idents, ident)

	for !p.lexer.isEOF() && p.tryConsumeTokenKind(TokenKindDot) != nil {
		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		idents = append(idents, ident)
	}
	return &JSONPath{
		Idents: idents,
	}, nil
}

func (p *Parser) parseJSONOption() (*JSONOption, error) {
	switch {
	case p.tryConsumeKeyword(KeywordSkip) != nil:
		if p.tryConsumeKeyword(KeywordRegexp) != nil {
			regex, err := p.parseString(p.Pos())
			if err != nil {
				return nil, err
			}
			return &JSONOption{
				SkipRegex: regex,
			}, nil
		}
		jsonPath, err := p.parseJSONPath()
		if err != nil {
			return nil, err
		}
		return &JSONOption{
			SkipPath: jsonPath,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected token kind: %s", p.lastTokenKind())
	}
}

func (p *Parser) parseJSONType(name *Ident, pos Pos) (*JSONType, error) {
	if p.matchTokenKind(TokenKindLParen) {
		return &JSONType{Name: name}, nil
	}

	options := make([]*JSONOption, 0)
	for !p.lexer.isEOF() && !p.matchTokenKind(TokenKindRParen) {
		option, err := p.parseJSONOption()
		if err != nil {
			return nil, err
		}
		options = append(options, option)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}

	rparenPos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &JSONType{
		Name: name,
		Options: &JSONOptions{
			LParen: pos,
			RParen: rparenPos,
			Items:  options,
		},
	}, nil
}

func (p *Parser) parseNestedType(name *Ident, pos Pos) (*NestedType, error) {
	columns, err := p.parseTableColumns()
	if err != nil {
		return nil, err
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &NestedType{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Name:          name,
		Columns:       columns,
	}, nil
}

func (p *Parser) tryParseCompressionCodecs(pos Pos) (*CompressionCodec, error) {
	if p.tryConsumeKeyword(KeywordCodec) == nil {
		return nil, nil // nolint
	}

	if _, err := p.consumeTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}

	// parse codec name
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	// parse DELTA if  CODEC(Delta, ZSTD(1))
	// or CODEC(Delta(9), ZSTD(1)) or CODEC(T64, ZSTD(1))
	var codecType *Ident
	var typeLevel *NumberLiteral
	switch strings.ToUpper(name.Name) {
	case "DELTA", "DOUBLEDELTA", "T64", "GORILLA":
		codecType = name
		// try parse delta level
		typeLevel, err = p.tryParseCompressionLevel(p.Pos())
		if err != nil {
			return nil, err
		}
		// consume comma
		if _, err := p.consumeTokenKind(TokenKindComma); err != nil {
			return nil, err
		}
		name, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
	}

	var level *NumberLiteral
	// TODO: check if the codec name is valid
	switch strings.ToUpper(name.Name) {
	case "ZSTD", "LZ4HC", "LH4":
		level, err = p.tryParseCompressionLevel(p.Pos())
		if err != nil {
			return nil, err
		}
	}

	rightParenPos := p.last().End
	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}

	return &CompressionCodec{
		CodecPos:      pos,
		RightParenPos: rightParenPos,
		Type:          codecType,
		TypeLevel:     typeLevel,
		Name:          name,
		Level:         level,
	}, nil
}

func (p *Parser) parseEnumValueExpr(pos Pos) (*EnumValue, error) {
	name, err := p.parseString(pos)
	if err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind(TokenKindSingleEQ); err != nil {
		return nil, err
	}

	value, err := p.parseNumber(p.Pos())
	if err != nil {
		return nil, err
	}
	return &EnumValue{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) parseColumnStar(pos Pos) (*Ident, error) {
	if _, err := p.consumeTokenKind("*"); err != nil {
		return nil, err
	}
	return &Ident{
		NamePos: pos,
		NameEnd: pos,
		Name:    "*",
	}, nil
}

func (p *Parser) tryParseCompressionLevel(pos Pos) (*NumberLiteral, error) {
	if p.tryConsumeTokenKind(TokenKindLParen) == nil {
		return nil, nil // nolint
	}

	num, err := p.parseNumber(pos)
	if err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return num, nil
}
