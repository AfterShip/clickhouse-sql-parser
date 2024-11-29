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
	case p.matchTokenKind(opTypeCast):
		return PrecedenceDoubleColon
	case p.matchTokenKind(opTypeEQ), p.matchTokenKind(opTypeLT), p.matchTokenKind(opTypeLE),
		p.matchTokenKind(opTypeGE), p.matchTokenKind(opTypeGT), p.matchTokenKind(opTypeDoubleEQ),
		p.matchTokenKind(opTypeNE), p.matchTokenKind("<>"):
		return PrecedenceCompare
	case p.matchTokenKind(opTypePlus), p.matchTokenKind(opTypeMinus):
		return PrecedenceAddSub
	case p.matchTokenKind(opTypeMul), p.matchTokenKind(opTypeDiv), p.matchTokenKind(opTypeMod):
		return PrecedenceMulDivMod
	case p.matchTokenKind(opTypeArrow):
		return PrecedenceArrow
	case p.matchTokenKind("("), p.matchTokenKind("["):
		return PrecedenceBracket
	case p.matchTokenKind(opTypeCast):
		return PrecedenceDoubleColon
	case p.matchKeyword(KeywordBetween), p.matchKeyword(KeywordLike), p.matchKeyword(KeywordIlike):
		return PrecedenceBetweenLike
	case p.matchKeyword(KeywordIn):
		return precedenceIn
	case p.matchKeyword(KeywordGlobal):
		return PrecedenceGlobal
	case p.matchTokenKind(opTypeQuery):
		return PrecedenceQuery
	default:
		return PrecedenceUnknown
	}
}

func (p *Parser) parseInfix(expr Expr, precedence int) (Expr, error) {
	switch {
	case p.matchTokenKind(opTypeEQ), p.matchTokenKind(opTypeLT), p.matchTokenKind(opTypeLE),
		p.matchTokenKind(opTypeGE), p.matchTokenKind(opTypeGT),
		p.matchTokenKind(opTypeNE), p.matchTokenKind("<>"),
		p.matchTokenKind(opTypeMinus), p.matchTokenKind(opTypePlus), p.matchTokenKind(opTypeMul),
		p.matchTokenKind(opTypeDiv), p.matchTokenKind(opTypeMod),
		p.matchKeyword(KeywordIn), p.matchKeyword(KeywordLike),
		p.matchKeyword(KeywordIlike), p.matchKeyword(KeywordAnd), p.matchKeyword(KeywordOr),
		p.matchTokenKind(opTypeCast), p.matchTokenKind(opTypeArrow), p.matchTokenKind(opTypeDoubleEQ):
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
	case p.matchTokenKind("["):
		params, err := p.parseArrayParams(p.Pos())
		if err != nil {
			return nil, err
		}
		return &ObjectParams{
			Object: expr,
			Params: params,
		}, nil
	case p.matchTokenKind(opTypeQuery):
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
	if _, err := p.consumeTokenKind("?"); err != nil {
		return nil, err
	}
	trueExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if _, err := p.consumeTokenKind(":"); err != nil {
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
	if _, err := p.consumeTokenKind("("); err != nil {
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
	if err := p.consumeKeyword(")"); err != nil {
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
	case p.matchTokenKind(opTypePlus),
		p.matchTokenKind(opTypeMinus),
		p.matchKeyword(KeywordNot):
		_ = p.lexer.consumeToken()
	default:
		return p.parseColumnExpr(pos)
	}

	var expr Expr
	var err error
	switch {
	case p.matchTokenKind(TokenIdent),
		p.matchTokenKind("("):
		expr, err = p.parseExpr(p.Pos())
	default:
		expr, err = p.parseColumnExpr(p.Pos())
	}
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
		if nextToken != nil && nextToken.Kind == TokenString {
			return p.parseString(p.Pos())
		}
		return p.parseIdentOrFunction(pos)
	case p.matchKeyword(KeywordCast):
		return p.parseColumnCastExpr(pos)
	case p.matchKeyword(KeywordCase):
		return p.parseColumnCaseExpr(pos)
	case p.matchKeyword(KeywordExtract):
		return p.parseColumnExtractExpr(pos)
	case p.matchTokenKind(TokenIdent):
		return p.parseIdentOrFunction(pos)
	case p.matchTokenKind(TokenString): // string literal
		return p.parseString(pos)
	case p.matchTokenKind(TokenInt),
		p.matchTokenKind(TokenFloat): // number literal
		return p.parseNumber(pos)
	case p.matchTokenKind("("):
		if peek, _ := p.lexer.peekToken(); peek != nil {
			if peek.Kind == TokenKeyword && strings.EqualFold(peek.String, KeywordSelect) {
				return p.parseSubQuery(pos)
			}
		}
		return p.parseFunctionParams(p.Pos())
	case p.matchTokenKind("*"):
		return p.parseColumnStar(p.Pos())
	case p.matchTokenKind("["):
		return p.parseArrayParams(p.Pos())
	case p.matchTokenKind("{"):
		// The map literal string also starts with '{', so we need to check the next token
		// to determine if it is a map literal or a query param.
		if p.peekTokenKind(TokenIdent) {
			return p.parseQueryParam(p.Pos())
		}
		return p.parseMapLiteral(p.Pos())
	case p.matchTokenKind(opTypeQuery):
		// Placeholder `?`
		_ = p.lexer.consumeToken()
		return &PlaceHolder{
			PlaceholderPos: pos,
			PlaceHolderEnd: pos,
			Type:           opTypeQuery,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected token kind: %s", p.lastTokenKind())
	}
}

func (p *Parser) parseColumnCastExpr(pos Pos) (Expr, error) {
	if err := p.consumeKeyword(KeywordCast); err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind("("); err != nil {
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
	asColumnType, err := p.parseColumnType(p.Pos())
	if err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind(")"); err != nil {
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

func (p *Parser) parseColumnExprListWithRoundBracket(pos Pos) (*ColumnExprList, error) {
	return p.parseColumnExprListWithTerm(")", pos)
}

func (p *Parser) parseColumnExprListWithSquareBracket(pos Pos) (*ColumnExprList, error) {
	return p.parseColumnExprListWithTerm("]", pos)
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
		if p.tryConsumeTokenKind(",") == nil {
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
		if p.tryConsumeTokenKind(",") == nil {
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
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}
	distinct := false
	if p.tryConsumeKeyword(KeywordDistinct) != nil {
		distinct = true
	}
	var items []Expr
	for !p.lexer.isEOF() && !p.matchTokenKind(")") {
		item, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
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
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}
	params, err := p.parseColumnExprListWithRoundBracket(p.Pos())
	if err != nil {
		return nil, err
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
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
	if p.matchTokenKind("(") {
		columnArgList, err := p.parseColumnArgList(p.Pos())
		if err != nil {
			return nil, err
		}
		paramExprList.ColumnArgList = columnArgList
	}
	return paramExprList, nil
}

func (p *Parser) parseMapLiteral(pos Pos) (*MapLiteral, error) {
	if _, err := p.consumeTokenKind("{"); err != nil {
		return nil, err
	}

	keyValues := make([]KeyValue, 0)
	for !p.lexer.isEOF() && !p.matchTokenKind("}") {
		key, err := p.parseString(p.Pos())
		if err != nil {
			return nil, err
		}
		if _, err := p.consumeTokenKind(":"); err != nil {
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
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	rightBracePos := p.Pos()
	if _, err := p.consumeTokenKind("}"); err != nil {
		return nil, err
	}
	return &MapLiteral{
		LBracePos: pos,
		RBracePos: rightBracePos,
		KeyValues: keyValues,
	}, nil
}

func (p *Parser) parseQueryParam(pos Pos) (*QueryParam, error) {
	if _, err := p.consumeTokenKind("{"); err != nil {
		return nil, err
	}

	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.consumeTokenKind(":"); err != nil {
		return nil, err
	}
	columnType, err := p.parseColumnType(p.Pos())
	if err != nil {
		return nil, err
	}
	rightBracePos := p.Pos()
	if _, err := p.consumeTokenKind("}"); err != nil {
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
	if _, err := p.consumeTokenKind("["); err != nil {
		return nil, err
	}
	params, err := p.parseColumnExprListWithSquareBracket(p.Pos())
	if err != nil {
		return nil, err
	}
	rightBracketPos := p.Pos()
	if _, err := p.consumeTokenKind("]"); err != nil {
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

func (p *Parser) parseColumnType(_ Pos) (Expr, error) { // nolint:funlen
	if p.matchTokenKind(TokenString) {
		return p.parseString(p.Pos())
	}
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if p.tryConsumeTokenKind("(") != nil {
		switch {
		case p.matchTokenKind(TokenIdent):
			if ident.Name == "Nested" {
				return p.parseNestedType(ident, p.Pos())
			}
			return p.parseComplexType(ident, p.Pos())
		case p.matchTokenKind(TokenString):
			if peekToken, err := p.lexer.peekToken(); err == nil && peekToken.Kind == opTypeEQ {
				// enum values
				return p.parseEnumExpr(ident,p.Pos())
			}
			// like Datetime('Asia/Dubai')
			return p.parseColumnTypeWithParams(ident, p.Pos())
		case p.matchTokenKind(TokenInt), p.matchTokenKind(TokenFloat):
			// fixed size
			return p.parseColumnTypeWithParams(ident, p.Pos())
		default:
			return nil, fmt.Errorf("unexpected token kind: %v", p.lastTokenKind())
		}
	}
	return &ScalarTypeExpr{Name: ident}, nil
}

func (p *Parser) parseColumnPropertyType(_ Pos) (Expr, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &PropertyTypeExpr{
		Name: ident,
	}, nil
}

func (p *Parser) parseComplexType(name *Ident, pos Pos) (Expr, error) {
	subTypes := make([]Expr, 0)
	for !p.lexer.isEOF() && !p.matchTokenKind(")") {
		subExpr, err := p.parseColumnType(p.Pos())
		if err != nil {
			return nil, err
		}
		subTypes = append(subTypes, subExpr)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return &ComplexTypeExpr{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Name:          name,
		Params:        subTypes,
	}, nil
}

func (p *Parser) parseEnumExpr(name *Ident, pos Pos) (*EnumValueList, error) {
	enumValueList := &EnumValueList{
		Name: name,
		ListPos: pos,
		Enums:   make([]EnumValue, 0),
	}
	for !p.lexer.isEOF() && !p.matchTokenKind(")") {
		enumValueExpr, err := p.parseEnumValueExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		if enumValueExpr == nil {
			break
		}
		enumValueList.Enums = append(enumValueList.Enums, *enumValueExpr)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	if len(enumValueList.Enums) > 0 {
		enumValueList.ListEnd = enumValueList.Enums[len(enumValueList.Enums)-1].Value.NumEnd
	}
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return enumValueList, nil
}

func (p *Parser) parseColumnTypeWithParams(name *Ident, pos Pos) (*TypeWithParamsExpr, error) {
	params := make([]Literal, 0)
	param, err := p.parseLiteral(p.Pos())
	if err != nil {
		return nil, err
	}
	params = append(params, param)
	for !p.lexer.isEOF() && p.tryConsumeTokenKind(",") != nil {
		size, err := p.parseLiteral(p.Pos())
		if err != nil {
			return nil, err
		}
		params = append(params, size)
	}

	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return &TypeWithParamsExpr{
		Name:          name,
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Params:        params,
	}, nil
}

func (p *Parser) parseNestedType(name *Ident, pos Pos) (*NestedTypeExpr, error) {
	columns, err := p.parseTableColumns()
	if err != nil {
		return nil, err
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return &NestedTypeExpr{
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

	if _, err := p.consumeTokenKind("("); err != nil {
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
		if _, err := p.consumeTokenKind(","); err != nil {
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
	if _, err := p.consumeTokenKind(")"); err != nil {
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

	if _, err := p.consumeTokenKind(opTypeEQ); err != nil {
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
	if p.tryConsumeTokenKind("(") == nil {
		return nil, nil // nolint
	}

	num, err := p.parseNumber(pos)
	if err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return num, nil
}
