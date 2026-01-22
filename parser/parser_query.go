package parser

import (
	"errors"
	"fmt"
)

func (p *Parser) tryParseWithClause(pos Pos) (*WithClause, error) {
	if !p.matchKeyword(KeywordWith) {
		return nil, nil
	}
	return p.parseWithClause(pos)
}

func (p *Parser) parseWithClause(pos Pos) (*WithClause, error) {
	if err := p.expectKeyword(KeywordWith); err != nil {
		return nil, err
	}

	cteExpr, err := p.parseCTEStmt(p.Pos())
	if err != nil {
		return nil, err
	}
	ctes := []*CTEStmt{cteExpr}
	for p.tryConsumeTokenKind(TokenKindComma) != nil {
		cteExpr, err := p.parseCTEStmt(p.Pos())
		if err != nil {
			return nil, err
		}
		ctes = append(ctes, cteExpr)
	}

	return &WithClause{
		WithPos: pos,
		CTEs:    ctes,
		EndPos:  ctes[len(ctes)-1].End(),
	}, nil
}

func (p *Parser) tryParseTopClause(pos Pos) (*TopClause, error) {
	if !p.matchKeyword(KeywordTop) {
		return nil, nil
	}
	return p.parseTopClause(pos)
}

func (p *Parser) parseTopClause(pos Pos) (*TopClause, error) {
	if err := p.expectKeyword(KeywordTop); err != nil {
		return nil, err
	}

	number, err := p.parseNumber(p.Pos())
	if err != nil {
		return nil, err
	}
	topEnd := number.End()

	withTies := false
	if p.tryConsumeKeywords(KeywordWith) {
		topEnd = p.End()
		if err := p.expectKeyword(KeywordTies); err != nil {
			return nil, err
		}
		withTies = true
	}
	return &TopClause{
		TopPos:   pos,
		TopEnd:   topEnd,
		Number:   number,
		WithTies: withTies,
	}, nil
}

func (p *Parser) tryParseDistinctOn(pos Pos) (*DistinctOn, error) {
	if !p.matchKeyword(KeywordOn) {
		return nil, nil
	}
	return p.parseDistinctOn(pos)
}

func (p *Parser) parseDistinctOn(pos Pos) (*DistinctOn, error) {
	if err := p.expectKeyword(KeywordOn); err != nil {
		return nil, err
	}

	if err := p.expectTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}

	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	idents := []*Ident{ident}

	for p.matchTokenKind(TokenKindComma) {
		_ = p.lexer.consumeToken()

		ident, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
		idents = append(idents, ident)
	}

	if err := p.expectTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}

	return &DistinctOn{
		Idents:        idents,
		DistinctOnPos: pos,
		DistinctOnEnd: p.Pos(),
	}, nil
}

func (p *Parser) tryParseFromClause(pos Pos) (*FromClause, error) {
	if !p.matchKeyword(KeywordFrom) {
		return nil, nil
	}
	return p.parseFromClause(pos)
}

func (p *Parser) parseFromClause(pos Pos) (*FromClause, error) {
	if err := p.expectKeyword(KeywordFrom); err != nil {
		return nil, err
	}

	expr, err := p.parseJoinExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &FromClause{
		FromPos: pos,
		Expr:    expr,
	}, nil
}

func (p *Parser) tryParseJoinConstraints(pos Pos) (Expr, error) {
	switch {
	case p.tryConsumeKeywords(KeywordOn):
		columnExprList, err := p.parseColumnExprList(p.Pos())
		if err != nil {
			return nil, err
		}
		return &OnClause{
			OnPos: pos,
			On:    columnExprList,
		}, nil
	case p.tryConsumeKeywords(KeywordUsing):
		hasParen := p.tryConsumeTokenKind(TokenKindLParen) != nil
		columnExprList, err := p.parseColumnExprListWithLParen(p.Pos())
		if err != nil {
			return nil, err
		}
		if hasParen {
			if err := p.expectTokenKind(TokenKindRParen); err != nil {
				return nil, err
			}
		}
		return &UsingClause{
			UsingPos: pos,
			Using:    columnExprList,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseJoinOp(_ Pos) []string {
	var modifiers []string
	switch {
	case p.tryConsumeKeywords(KeywordCross): // cross join
		modifiers = append(modifiers, KeywordCross)
	case p.matchKeyword(KeywordAny), p.matchKeyword(KeywordAll):
		modifiers = append(modifiers, p.last().String)
		_ = p.lexer.consumeToken()
		if p.matchKeyword(KeywordFull) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
		if p.matchKeyword(KeywordLeft) || p.matchKeyword(KeywordRight) || p.matchKeyword(KeywordInner) || p.matchKeyword(KeywordOuter) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
	case p.matchKeyword(KeywordSemi), p.matchKeyword(KeywordAsof):
		modifiers = append(modifiers, p.last().String)
		_ = p.lexer.consumeToken()
		if p.matchKeyword(KeywordLeft) || p.matchKeyword(KeywordRight) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
		if p.matchKeyword(KeywordOuter) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
	case p.matchKeyword(KeywordInner):
		modifiers = append(modifiers, p.last().String)
		_ = p.lexer.consumeToken()
		if p.matchKeyword(KeywordAll) || p.matchKeyword(KeywordAny) || p.matchKeyword(KeywordAsof) || p.matchKeyword(KeywordArray) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
	case p.matchKeyword(KeywordLeft):
		modifiers = append(modifiers, p.last().String)
		_ = p.lexer.consumeToken()
		if p.matchKeyword(KeywordOuter) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
		if p.matchKeyword(KeywordSemi) || p.matchKeyword(KeywordAnti) ||
			p.matchKeyword(KeywordAny) || p.matchKeyword(KeywordAll) ||
			p.matchKeyword(KeywordAsof) || p.matchKeyword(KeywordArray) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
	case p.matchKeyword(KeywordRight):
		modifiers = append(modifiers, p.last().String)
		_ = p.lexer.consumeToken()
		if p.matchKeyword(KeywordOuter) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
		if p.matchKeyword(KeywordSemi) || p.matchKeyword(KeywordAnti) ||
			p.matchKeyword(KeywordAny) || p.matchKeyword(KeywordAll) ||
			p.matchKeyword(KeywordAsof) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
	case p.matchKeyword(KeywordFull):
		modifiers = append(modifiers, p.last().String)
		_ = p.lexer.consumeToken()
		if p.matchKeyword(KeywordOuter) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
		if p.matchKeyword(KeywordAll) || p.matchKeyword(KeywordAny) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
	case p.matchKeyword(KeywordArray):
		modifiers = append(modifiers, p.last().String)
		_ = p.lexer.consumeToken()
	}
	return modifiers
}

func (p *Parser) parseJoinTableExpr(_ Pos) (Expr, error) {
	switch {
	case p.matchTokenKind(TokenKindIdent), p.matchTokenKind(TokenKindString), p.matchTokenKind(TokenKindLParen):
		tableExpr, err := p.parseTableExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		statementEnd := tableExpr.End()

		hasFinal := p.matchKeyword(KeywordFinal)
		if hasFinal {
			statementEnd = p.End()
			_ = p.lexer.consumeToken()
		}

		sampleRatio, err := p.tryParseSampleClause(p.Pos())
		if err != nil {
			return nil, err
		}
		if sampleRatio != nil {
			statementEnd = sampleRatio.End()
		}
		return &JoinTableExpr{
			Table:        tableExpr,
			SampleRatio:  sampleRatio,
			HasFinal:     hasFinal,
			StatementEnd: statementEnd,
		}, nil
	default:
		return nil, fmt.Errorf("expected table name or subquery, got %s", fmt.Sprintf("%v", p.lastTokenKind()))
	}
}

func (p *Parser) parseJoinRightExpr(pos Pos) (expr Expr, err error) {
	var rightExpr Expr
	var modifiers []string
	switch {
	case p.tryConsumeKeywords(KeywordGlobal):
	case p.tryConsumeKeywords(KeywordLocal):
	case p.tryConsumeTokenKind(TokenKindComma) != nil:
		return p.parseJoinExpr(p.Pos())
	default:
		modifiers = p.parseJoinOp(p.Pos())
	}

	if len(modifiers) != 0 && !p.matchKeyword(KeywordJoin) {
		return nil, fmt.Errorf("expected JOIN, got %s", p.lastTokenKind())
	}
	if !p.tryConsumeKeywords(KeywordJoin) {
		return nil, nil
	}

	modifiers = append(modifiers, KeywordJoin)

	// Check if this is an ARRAY JOIN
	isArrayJoin := false
	for _, mod := range modifiers {
		if mod == KeywordArray {
			isArrayJoin = true
			break
		}
	}

	if isArrayJoin {
		// For ARRAY JOIN, parse column expression list instead of table expression
		expr, err = p.parseColumnExprList(p.Pos())
		if err != nil {
			return nil, err
		}

		// ARRAY JOIN doesn't have constraints (ON/USING)
		// try parse next join
		rightExpr, err = p.parseJoinRightExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		return &JoinExpr{
			JoinPos:     pos,
			Left:        expr,
			Right:       rightExpr,
			Modifiers:   modifiers,
			Constraints: nil,
		}, nil
	}

	expr, err = p.parseJoinTableExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	constrains, err := p.tryParseJoinConstraints(p.Pos())
	if err != nil {
		return nil, err
	}

	// try parse next join
	rightExpr, err = p.parseJoinRightExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &JoinExpr{
		JoinPos:     pos,
		Left:        expr,
		Right:       rightExpr,
		Modifiers:   modifiers,
		Constraints: constrains,
	}, nil
}

func (p *Parser) parseJoinExpr(pos Pos) (expr Expr, err error) {
	if expr, err = p.parseJoinTableExpr(p.Pos()); err != nil {
		return nil, err
	}
	rightExpr, err := p.parseJoinRightExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if rightExpr == nil {
		return expr, nil
	}
	return &JoinExpr{
		JoinPos: pos,
		Left:    expr,
		Right:   rightExpr,
	}, nil
}

func (p *Parser) parseTableExpr(pos Pos) (*TableExpr, error) {
	var expr Expr
	var err error
	switch {
	case p.matchTokenKind(TokenKindString), p.matchTokenKind(TokenKindIdent):
		// table name
		tableIdentifier, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		// it's a table name
		if tableIdentifier.Database != nil || !p.matchTokenKind(TokenKindLParen) { // database.table
			expr = tableIdentifier
		} else {
			// table function expr
			tableArgs, err := p.parseTableArgList(p.Pos())
			if err != nil {
				return nil, err
			}
			expr = &TableFunctionExpr{
				Name: tableIdentifier.Table,
				Args: tableArgs,
			}
		}
	case p.matchTokenKind(TokenKindLParen):
		expr, err = p.parseSubQuery(p.Pos())
	default:
		return nil, errors.New("expect table name or subquery")
	}
	if err != nil {
		return nil, err
	}

	tableEnd := expr.End()
	if p.tryConsumeKeywords(KeywordAs) {
		alias, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		expr = &AliasExpr{
			Expr:     expr,
			AliasPos: alias.Pos(),
			Alias:    alias,
		}
		tableEnd = expr.End()
	} else if p.matchTokenKind(TokenKindIdent) && p.lastTokenKind() != TokenKindKeyword {
		alias, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		expr = &AliasExpr{
			Expr:     expr,
			AliasPos: alias.Pos(),
			Alias:    alias,
		}
		tableEnd = expr.End()
	}

	isFinalExist := false
	if p.tryConsumeKeywords(KeywordFinal) {
		switch expr.(type) {
		case *TableFunctionExpr:
			return nil, errors.New("table function doesn't support FINAL")
		case *SelectQuery:
			return nil, errors.New("subquery doesn't support FINAL")
		}
		isFinalExist = true
		tableEnd = expr.End()
	}

	return &TableExpr{
		TablePos: pos,
		TableEnd: tableEnd,
		Expr:     expr,
		HasFinal: isFinalExist,
	}, nil
}

func (p *Parser) tryParsePrewhereClause(pos Pos) (*PrewhereClause, error) {
	if !p.matchKeyword(KeywordPrewhere) {
		return nil, nil
	}
	return p.parsePrewhereClause(pos)
}
func (p *Parser) parsePrewhereClause(pos Pos) (*PrewhereClause, error) {
	if err := p.expectKeyword(KeywordPrewhere); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &PrewhereClause{
		PrewherePos: pos,
		Expr:        expr,
	}, nil
}

func (p *Parser) tryParseWhereClause(pos Pos) (*WhereClause, error) {
	if !p.matchKeyword(KeywordWhere) {
		return nil, nil
	}
	return p.parseWhereClause(pos)
}

func (p *Parser) parseWhereClause(pos Pos) (*WhereClause, error) {
	if err := p.expectKeyword(KeywordWhere); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &WhereClause{
		WherePos: pos,
		Expr:     expr,
	}, nil
}

func (p *Parser) tryParseGroupByClause(pos Pos) (*GroupByClause, error) {
	if !p.matchKeyword(KeywordGroup) {
		return nil, nil
	}
	return p.parseGroupByClause(pos)
}

// syntax: groupByClause? (WITH (CUBE | ROLLUP))? (WITH TOTALS)?
func (p *Parser) parseGroupByClause(pos Pos) (*GroupByClause, error) {
	if err := p.expectKeyword(KeywordGroup); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordBy); err != nil {
		return nil, err
	}

	var expr Expr
	var err error
	aggregateType := ""
	switch {
	case p.matchKeyword(KeywordCube) || p.matchKeyword(KeywordRollup):
		aggregateType = p.last().String
		_ = p.lexer.consumeToken()
		expr, err = p.parseFunctionParams(p.Pos())
	case p.tryConsumeKeywords(KeywordGrouping, KeywordSets):
		aggregateType = "GROUPING SETS"
		expr, err = p.parseFunctionParams(p.Pos())
	case p.tryConsumeKeywords(KeywordAll):
		aggregateType = "ALL"
	default:
		expr, err = p.parseColumnExprListWithLParen(p.Pos())
	}
	if err != nil {
		return nil, err
	}
	groupBy := &GroupByClause{
		GroupByPos:    pos,
		AggregateType: aggregateType,
		Expr:          expr,
	}

	// parse WITH CUBE, ROLLUP, TOTALS
	for p.tryConsumeKeywords(KeywordWith) {
		switch {
		case p.tryConsumeKeywords(KeywordCube):
			groupBy.WithCube = true
		case p.tryConsumeKeywords(KeywordRollup):
			groupBy.WithRollup = true
		case p.tryConsumeKeywords(KeywordTotals):
			groupBy.WithTotals = true
		default:
			return nil, fmt.Errorf("expected CUBE, ROLLUP or TOTALS, got %s", p.lastTokenKind())
		}
	}
	groupBy.GroupByEnd = p.Pos()

	return groupBy, nil
}

func (p *Parser) tryParseLimitAfterLimitByClause(pos Pos) (*LimitClause, error) {
	if !p.matchKeyword(KeywordLimit) {
		return nil, nil
	}

	return p.parseLimitClause(pos)
}

func (p *Parser) tryParseLimitClause(pos Pos) (*LimitClause, error) {
	if !p.matchKeyword(KeywordLimit) && !p.matchKeyword(KeywordOffset) {
		return nil, nil
	}

	return p.parseLimitClause(pos)
}

func (p *Parser) parseLimitClause(pos Pos) (*LimitClause, error) {
	var limit Expr
	var offset Expr
	var err error
	if p.tryConsumeKeywords(KeywordLimit) {
		limit, err = p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}

		if p.tryConsumeKeywords(KeywordOffset) {
			offset, err = p.parseExpr(p.Pos())
		} else if p.tryConsumeTokenKind(TokenKindComma) != nil {
			offset = limit
			limit, err = p.parseExpr(p.Pos())
		}
	} else if p.tryConsumeKeywords(KeywordOffset) {
		offset, err = p.parseExpr(p.Pos())
	}

	if err != nil {
		return nil, err
	}

	return &LimitClause{
		LimitPos: pos,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (p *Parser) tryParseLimitByClause(pos Pos) (Expr, error) {
	if !p.matchKeyword(KeywordLimit) {
		return nil, nil
	}
	return p.parseLimitByClause(pos)
}

func (p *Parser) parseBetweenClause(expr Expr) (*BetweenClause, error) {
	if err := p.expectKeyword(KeywordBetween); err != nil {
		return nil, err
	}

	betweenExpr, err := p.parseSubExpr(p.Pos(), PrecedenceBetweenLike)
	if err != nil {
		return nil, err
	}

	andPos := p.Pos()
	if err := p.expectKeyword(KeywordAnd); err != nil {
		return nil, err
	}

	andExpr, err := p.parseSubExpr(p.Pos(), PrecedenceBetweenLike)
	if err != nil {
		return nil, err
	}

	return &BetweenClause{
		Expr:    expr,
		Between: betweenExpr,
		AndPos:  andPos,
		And:     andExpr,
	}, nil
}

func (p *Parser) parseLimitByClause(pos Pos) (Expr, error) {
	limit, err := p.parseLimitClause(pos)
	if err != nil {
		return nil, err
	}

	var by *ColumnExprList
	if !p.tryConsumeKeywords(KeywordBy) {
		return limit, nil
	}
	if by, err = p.parseColumnExprListWithLParen(p.Pos()); err != nil {
		return nil, err
	}
	return &LimitByClause{
		Limit:  limit,
		ByExpr: by,
	}, nil
}

func (p *Parser) tryParseWindowFrameClause(pos Pos) (*WindowFrameClause, error) {
	if !p.matchKeyword(KeywordRows) && !p.matchKeyword(KeywordRange) {
		return nil, nil
	}
	return p.parseWindowFrameClause(pos)
}

func (p *Parser) parseWindowFrameClause(pos Pos) (*WindowFrameClause, error) {
	var windowFrameType string
	if p.matchKeyword(KeywordRows) || p.matchKeyword(KeywordRange) {
		windowFrameType = p.last().String
		_ = p.lexer.consumeToken()
	} else {
		return nil, fmt.Errorf("expected ROWS or RANGE for window frame")
	}

	var expr Expr
	if p.tryConsumeKeywords(KeywordBetween) {
		left, err := p.parseFrameExtent()
		if err != nil {
			return nil, err
		}
		andPos := p.Pos()
		if err := p.expectKeyword(KeywordAnd); err != nil {
			return nil, err
		}
		right, err := p.parseFrameExtent()
		if err != nil {
			return nil, err
		}
		expr = &BetweenClause{
			Between: left,
			AndPos:  andPos,
			And:     right,
		}
	} else {
		// single extent
		extent, err := p.parseFrameExtent()
		if err != nil {
			return nil, err
		}
		expr = extent
	}

	return &WindowFrameClause{
		FramePos: pos,
		Type:     windowFrameType,
		Extend:   expr,
	}, nil
}

// parseFrameExtent parses a single frame extent
func (p *Parser) parseFrameExtent() (Expr, error) {
	switch {
	case p.matchKeyword(KeywordCurrent):
		return p.parseFrameCurrentRow()
	case p.matchKeyword(KeywordUnbounded):
		return p.parseFrameUnbounded()
	case p.matchTokenKind(TokenKindInt):
		return p.parseFrameNumber()
	case p.matchTokenKind(TokenKindLBrace):
		return p.parseFrameParam()
	case p.matchKeyword(KeywordInterval):
		return p.parseFrameInterval()
	default:
		return nil, fmt.Errorf("expected UNBOUNDED, CURRENT ROW, integer, parameter, or interval")
	}
}

func (p *Parser) parseFrameCurrentRow() (Expr, error) {
	currentPos := p.Pos()
	_ = p.lexer.consumeToken()
	if err := p.expectKeyword(KeywordRow); err != nil {
		return nil, err
	}
	rowEnd := p.End()
	return &WindowFrameCurrentRow{
		CurrentPos: currentPos,
		RowEnd:     rowEnd,
	}, nil
}

func (p *Parser) parseFrameUnbounded() (Expr, error) {
	unboundedPos := p.Pos()
	_ = p.lexer.consumeToken()

	direction, err := p.parseFrameDirection()
	if err != nil {
		return nil, err
	}
	return &WindowFrameUnbounded{
		UnboundedPos: unboundedPos,
		Direction:    direction,
	}, nil
}

func (p *Parser) parseFrameNumber() (Expr, error) {
	number, err := p.parseNumber(p.Pos())
	if err != nil {
		return nil, err
	}

	direction, endPos, err := p.parseFrameDirectionWithEnd()
	if err != nil {
		return nil, err
	}
	return &WindowFrameNumber{
		EndPos:    endPos,
		Number:    number,
		Direction: direction,
	}, nil
}

func (p *Parser) parseFrameParam() (Expr, error) {
	queryParam, err := p.parseQueryParam(p.Pos())
	if err != nil {
		return nil, err
	}

	direction, endPos, err := p.parseFrameDirectionWithEnd()
	if err != nil {
		return nil, err
	}
	return &WindowFrameParam{
		Param:     queryParam,
		EndPos:    endPos,
		Direction: direction,
	}, nil
}

func (p *Parser) parseFrameInterval() (Expr, error) {
	intervalExpr, err := p.parseInterval(true)
	if err != nil {
		return nil, err
	}

	direction, endPos, err := p.parseFrameDirectionWithEnd()
	if err != nil {
		return nil, err
	}
	return &WindowFrameExtendExpr{
		Expr:      intervalExpr,
		Direction: direction,
		EndPos:    endPos,
	}, nil
}

func (p *Parser) parseFrameDirection() (string, error) {
	switch {
	case p.matchKeyword(KeywordPreceding), p.matchKeyword(KeywordFollowing):
		direction := p.last().String
		_ = p.lexer.consumeToken()
		return direction, nil
	default:
		return "", fmt.Errorf("expected PRECEDING or FOLLOWING, got %s", p.lastTokenKind())
	}
}

func (p *Parser) parseFrameDirectionWithEnd() (string, Pos, error) {
	if !p.matchKeyword(KeywordPreceding) && !p.matchKeyword(KeywordFollowing) {
		return "", 0, fmt.Errorf("expected PRECEDING or FOLLOWING, got %s", p.lastTokenKind())
	}
	endPos := p.End()
	direction := p.last().String
	_ = p.lexer.consumeToken()
	return direction, endPos, nil
}

func (p *Parser) tryParseWindowClause(pos Pos) (*WindowClause, error) {
	if !p.matchKeyword(KeywordWindow) {
		return nil, nil
	}
	return p.parseWindowClause(pos)
}

func (p *Parser) parseWindowCondition(pos Pos) (*WindowExpr, error) {
	if err := p.expectTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}
	partitionBy, err := p.tryParsePartitionByClause(pos)
	if err != nil {
		return nil, err
	}
	orderBy, err := p.tryParseOrderByClause(p.Pos())
	if err != nil {
		return nil, err
	}
	frame, err := p.tryParseWindowFrameClause(p.Pos())
	if err != nil {
		return nil, err
	}
	rightParenPos := p.Pos()
	if err := p.expectTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &WindowExpr{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		PartitionBy:   partitionBy,
		OrderBy:       orderBy,
		Frame:         frame,
	}, nil
}

func (p *Parser) parseWindowClause(pos Pos) (*WindowClause, error) {
	if err := p.expectKeyword(KeywordWindow); err != nil {
		return nil, err
	}

	windows := make([]*WindowDefinition, 0, 1)
	for {
		windowName, err := p.parseIdent()
		if err != nil {
			return nil, err
		}

		asPos := p.Pos()
		if err := p.expectKeyword(KeywordAs); err != nil {
			return nil, err
		}

		condition, err := p.parseWindowCondition(p.Pos())
		if err != nil {
			return nil, err
		}

		windows = append(windows, &WindowDefinition{
			Name:  windowName,
			AsPos: asPos,
			Expr:  condition,
		})

		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}

	var endPos Pos
	if len(windows) > 0 {
		endPos = windows[len(windows)-1].End()
	}

	return &WindowClause{
		WindowPos: pos,
		EndPos:    endPos,
		Windows:   windows,
	}, nil
}


func (p *Parser) tryParseHavingClause(pos Pos) (*HavingClause, error) {
	if !p.matchKeyword(KeywordHaving) {
		return nil, nil
	}
	return p.parseHavingClause(pos)
}

func (p *Parser) parseHavingClause(pos Pos) (*HavingClause, error) {
	if err := p.expectKeyword(KeywordHaving); err != nil {
		return nil, err
	}

	expr, err := p.parseColumnsExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	return &HavingClause{
		HavingPos: pos,
		Expr:      expr,
	}, nil
}

func (p *Parser) parseSubQuery(_ Pos) (*SubQuery, error) {

	hasParen := p.tryConsumeTokenKind(TokenKindLParen) != nil

	selectQuery, err := p.parseSelectQuery(p.Pos())
	if err != nil {
		return nil, err
	}
	if hasParen {
		if err := p.expectTokenKind(TokenKindRParen); err != nil {
			return nil, err
		}
	}

	return &SubQuery{
		HasParen: hasParen,
		Select:   selectQuery,
	}, nil
}

func (p *Parser) parseSelectQuery(_ Pos) (*SelectQuery, error) {
	if !p.matchKeyword(KeywordSelect) && !p.matchKeyword(KeywordWith) && !p.matchTokenKind(TokenKindLParen) {
		return nil, fmt.Errorf("expected SELECT, WITH or (, got %s", p.lastTokenKind())
	}

	hasParen := p.tryConsumeTokenKind(TokenKindLParen) != nil
	selectStmt, err := p.parseSelectStmt(p.Pos())
	if err != nil {
		return nil, err
	}
	switch {
	case p.tryConsumeKeywords(KeywordUnion):
		switch {
		case p.tryConsumeKeywords(KeywordAll):
			unionAllExpr, err := p.parseSelectQuery(p.Pos())
			if err != nil {
				return nil, err
			}
			selectStmt.UnionAll = unionAllExpr
		case p.tryConsumeKeywords(KeywordDistinct):
			unionDistinctExpr, err := p.parseSelectQuery(p.Pos())
			if err != nil {
				return nil, err
			}
			selectStmt.UnionDistinct = unionDistinctExpr
		default:
			return nil, fmt.Errorf("expected ALL or DISTINCT, got %s", p.lastTokenKind())
		}
	case p.tryConsumeKeywords(KeywordExcept):
		exceptExpr, err := p.parseSelectQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		selectStmt.Except = exceptExpr
	}
	if hasParen {
		if err := p.expectTokenKind(TokenKindRParen); err != nil {
			return nil, err
		}
	}
	return selectStmt, nil
}

func (p *Parser) parseSelectStmt(pos Pos) (*SelectQuery, error) { // nolint: funlen
	withClause, err := p.tryParseWithClause(pos)
	if err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordSelect); err != nil {
		return nil, err
	}
	// DISTINCT?
	hasDistinct := p.tryConsumeKeywords(KeywordDistinct)
	distinctOn, err := p.tryParseDistinctOn(p.Pos())
	if err != nil {
		return nil, err
	}

	top, err := p.tryParseTopClause(p.Pos())
	if err != nil {
		return nil, err
	}
	selectItems, err := p.parseSelectItems()
	if err != nil {
		return nil, err
	}

	statementEnd := pos
	if len(selectItems) > 0 {
		statementEnd = selectItems[len(selectItems)-1].End()
	}
	from, err := p.tryParseFromClause(p.Pos())
	if err != nil {
		return nil, err
	}

	if from != nil {
		statementEnd = from.End()
	}
	prewhere, err := p.tryParsePrewhereClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if prewhere != nil {
		statementEnd = prewhere.End()
	}
	where, err := p.tryParseWhereClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if where != nil {
		statementEnd = where.End()
	}
	groupBy, err := p.tryParseGroupByClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if groupBy != nil {
		statementEnd = groupBy.End()
	}
	withTotal := false
	lastPos := p.Pos()
	if p.tryConsumeKeywords(KeywordWith) {
		if err := p.expectKeyword(KeywordTotals); err != nil {
			return nil, err
		}
		withTotal = true
		statementEnd = lastPos
	}
	having, err := p.tryParseHavingClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if having != nil {
		statementEnd = having.End()
	}
	window, err := p.tryParseWindowClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if window != nil {
		statementEnd = window.End()
	}
	orderBy, err := p.tryParseOrderByClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if orderBy != nil {
		statementEnd = orderBy.End()
	}

	var limitBy *LimitByClause
	var limit *LimitClause
	parsedLimitBy, err := p.tryParseLimitByClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if parsedLimitBy != nil {
		statementEnd = parsedLimitBy.End()
		switch e := parsedLimitBy.(type) {
		case *LimitByClause:
			limitBy = e
			limit, err = p.tryParseLimitAfterLimitByClause(p.Pos())
			if err != nil {
				return nil, err
			}
			if limit != nil {
				statementEnd = limit.End()
			}
		case *LimitClause:
			limit = e
		}
	} else {
		limit, err = p.tryParseLimitClause(p.Pos())
		if err != nil {
			return nil, err
		}
		if limit != nil {
			statementEnd = limit.End()
		}
	}

	settings, err := p.tryParseSettingsClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if settings != nil {
		statementEnd = settings.End()
	}

	format, err := p.tryParseFormat(p.Pos())
	if err != nil {
		return nil, err
	}
	if format != nil {
		statementEnd = format.End()
	}

	return &SelectQuery{
		With:         withClause,
		SelectPos:    pos,
		StatementEnd: statementEnd,
		Top:          top,
		HasDistinct:  hasDistinct,
		DistinctOn:   distinctOn,
		SelectItems:  selectItems,
		From:         from,
		Window:       window,
		Prewhere:     prewhere,
		Where:        where,
		GroupBy:      groupBy,
		Having:       having,
		OrderBy:      orderBy,
		LimitBy:      limitBy,
		Limit:        limit,
		Settings:     settings,
		Format:       format,
		WithTotal:    withTotal,
	}, nil
}

func (p *Parser) parseCTEStmt(pos Pos) (*CTEStmt, error) {
	expr, err := p.parseExpr(pos)
	if err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordAs); err != nil {
		return nil, err
	}
	if p.matchTokenKind(TokenKindLParen) {
		selectQuery, err := p.parseSelectQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		return &CTEStmt{
			CTEPos: pos,
			Expr:   expr,
			Alias:  selectQuery,
		}, nil
	}
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	return &CTEStmt{
		CTEPos: pos,
		Expr:   expr,
		Alias:  name,
	}, nil
}

func (p *Parser) tryParseSampleClause(pos Pos) (*SampleClause, error) {
	if !p.matchKeyword(KeywordSample) {
		return nil, nil
	}
	return p.parseSampleClause(pos)
}

func (p *Parser) parseSampleClause(pos Pos) (*SampleClause, error) {
	if err := p.expectKeyword(KeywordSample); err != nil {
		return nil, err
	}
	ratio, err := p.parseRatioExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	var offset *RatioExpr
	if p.matchKeyword(KeywordOffset) {
		_ = p.lexer.consumeToken()
		offset, err = p.parseRatioExpr(p.Pos())
		if err != nil {
			return nil, err
		}
	}

	return &SampleClause{
		SamplePos: pos,
		Ratio:     ratio,
		Offset:    offset,
	}, nil
}

func (p *Parser) parseExplainStmt(pos Pos) (*ExplainStmt, error) {
	if err := p.expectKeyword(KeywordExplain); err != nil {
		return nil, err
	}

	var explainType string
	switch {
	case p.matchKeyword(KeywordSyntax),
		p.matchKeyword(KeywordPipeline),
		p.matchKeyword(KeywordEstimate),
		p.matchKeyword(KeywordAst):
		explainType = p.last().String
		_ = p.lexer.consumeToken()
	default:
		return nil, fmt.Errorf("expected SYNTAX, PIPELINE, ESTIMATE or AST, got %s", p.lastTokenKind())
	}
	stmt, err := p.parseSelectQuery(p.Pos())
	if err != nil {
		return nil, err
	}
	return &ExplainStmt{
		ExplainPos: pos,
		Type:       explainType,
		Statement:  stmt,
	}, nil
}
