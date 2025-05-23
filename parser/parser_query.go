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
		if p.matchKeyword(KeywordAll) || p.matchKeyword(KeywordAny) || p.matchKeyword(KeywordAsof) {
			modifiers = append(modifiers, p.last().String)
			_ = p.lexer.consumeToken()
		}
	case p.matchKeyword(KeywordLeft), p.matchKeyword(KeywordRight):
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
	}
	return modifiers
}

func (p *Parser) parseJoinTableExpr(_ Pos) (Expr, error) {
	switch {
	case p.matchTokenKind(TokenKindIdent), p.matchTokenKind(TokenKindLParen):
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

func (p *Parser) tryParseLimitClause(pos Pos) (*LimitClause, error) {
	if !p.matchKeyword(KeywordLimit) {
		return nil, nil
	}
	return p.parseLimitClause(pos)
}

func (p *Parser) parseLimitClause(pos Pos) (*LimitClause, error) {
	if err := p.expectKeyword(KeywordLimit); err != nil {
		return nil, err
	}

	limit, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	var offset Expr
	if p.tryConsumeKeywords(KeywordOffset) {
		offset, err = p.parseExpr(p.Pos())
	} else if p.tryConsumeTokenKind(TokenKindComma) != nil {
		offset = limit
		limit, err = p.parseExpr(p.Pos())
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
	}

	var expr Expr
	switch {
	case p.tryConsumeKeywords(KeywordBetween):
		betweenWindowFrame, err := p.parseWindowFrameClause(p.Pos())
		if err != nil {
			return nil, err
		}

		andPos := p.Pos()
		if err := p.expectKeyword(KeywordAnd); err != nil {
			return nil, err
		}

		andWindowFrame, err := p.parseWindowFrameClause(p.Pos())
		if err != nil {
			return nil, err
		}
		expr = &BetweenClause{
			Between: betweenWindowFrame,
			AndPos:  andPos,
			And:     andWindowFrame,
		}
	case p.matchKeyword(KeywordCurrent):
		currentPos := p.Pos()
		_ = p.lexer.consumeToken()
		rowEnd := p.End()
		if err := p.expectKeyword(KeywordRow); err != nil {
			return nil, err
		}
		expr = &WindowFrameCurrentRow{
			CurrentPos: currentPos,
			RowEnd:     rowEnd,
		}
	case p.matchKeyword(KeywordUnbounded):
		unboundedPos := p.Pos()
		_ = p.lexer.consumeToken()

		direction := ""
		switch {
		case p.matchKeyword(KeywordPreceding), p.matchKeyword(KeywordFollowing):
			direction = p.last().String
			_ = p.lexer.consumeToken()
		default:
			return nil, fmt.Errorf("expected PRECEDING or FOLLOWING, got %s", p.lastTokenKind())
		}
		expr = &WindowFrameUnbounded{
			UnboundedPos: unboundedPos,
			Direction:    direction,
		}
	case p.matchTokenKind(TokenKindInt):
		number, err := p.parseNumber(p.Pos())
		if err != nil {
			return nil, err
		}

		var unboundedEnd Pos
		direction := ""
		switch {
		case p.matchKeyword(KeywordPreceding), p.matchKeyword(KeywordFollowing):
			direction = p.last().String
			unboundedEnd = p.End()
			_ = p.lexer.consumeToken()
		default:
			return nil, fmt.Errorf("expected PRECEDING or FOLLOWING, got %s", p.lastTokenKind())
		}
		expr = &WindowFrameNumber{
			UnboundedEnd: unboundedEnd,
			Number:       number,
			Direction:    direction,
		}
	default:
		return nil, fmt.Errorf("expected BETWEEN, CURRENT, UNBOUNDED or integer, got %s", p.lastTokenKind())
	}
	return &WindowFrameClause{
		FramePos: pos,
		Type:     windowFrameType,
		Extend:   expr,
	}, nil
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

	windowName, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordAs); err != nil {
		return nil, err
	}

	condition, err := p.parseWindowCondition(p.Pos())
	if err != nil {
		return nil, err
	}

	return &WindowClause{
		WindowPos:  pos,
		Name:       windowName,
		WindowExpr: condition,
	}, nil
}

func (p *Parser) tryParseArrayJoinClause(pos Pos) (*ArrayJoinClause, error) {
	if !p.matchKeyword(KeywordLeft) && !p.matchKeyword(KeywordInner) && !p.matchKeyword(KeywordArray) {
		return nil, nil
	}
	return p.parseArrayJoinClause(pos)
}

func (p *Parser) parseArrayJoinClause(_ Pos) (*ArrayJoinClause, error) {
	var typ string
	switch {
	case p.matchKeyword(KeywordLeft), p.matchKeyword(KeywordInner):
		typ = p.last().String
		_ = p.lexer.consumeToken()
	}
	arrayPos := p.Pos()
	if err := p.expectKeyword(KeywordArray); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordJoin); err != nil {
		return nil, err
	}

	expr, err := p.parseColumnExprList(p.Pos())
	if err != nil {
		return nil, err
	}

	return &ArrayJoinClause{
		ArrayPos: arrayPos,
		Type:     typ,
		Expr:     expr,
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
	window, err := p.tryParseWindowClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if window != nil {
		selectQuery.Window = window
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
			unionDistinctExpr, err := p.parseSelectStmt(p.Pos())
			if err != nil {
				return nil, err
			}
			selectStmt.UnionDistinct = unionDistinctExpr
		default:
			return nil, fmt.Errorf("expected ALL or DISTINCT, got %s", p.lastTokenKind())
		}
	case p.tryConsumeKeywords(KeywordExcept):
		exceptExpr, err := p.parseSelectStmt(p.Pos())
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
	arrayJoin, err := p.tryParseArrayJoinClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if arrayJoin != nil {
		statementEnd = arrayJoin.End()
	}
	window, err := p.tryParseWindowClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if window != nil {
		statementEnd = window.End()
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
			limit, err = p.tryParseLimitClause(p.Pos())
			if err != nil {
				return nil, err
			}
			if limit != nil {
				statementEnd = limit.End()
			}
		case *LimitClause:
			limit = e
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
		SelectItems:  selectItems,
		From:         from,
		ArrayJoin:    arrayJoin,
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

func (p *Parser) tryParseColumnAliases() ([]*Ident, error) {
	if !p.matchTokenKind(TokenKindLParen) {
		return nil, nil
	}
	if err := p.expectTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}

	aliasList := make([]*Ident, 0)
	for {
		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		aliasList = append(aliasList, ident)
		if p.matchTokenKind(TokenKindRParen) {
			break
		}
		if err := p.expectTokenKind(TokenKindComma); err != nil {
			return nil, err
		}
	}
	if err := p.expectTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return aliasList, nil
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
