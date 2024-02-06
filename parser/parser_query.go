package parser

import (
	"errors"
	"fmt"
)

func (p *Parser) tryParseWithExpr(pos Pos) (*WithExpr, error) {
	if !p.matchKeyword(KeywordWith) {
		return nil, nil
	}
	return p.parseWithExpr(pos)
}

func (p *Parser) parseWithExpr(pos Pos) (*WithExpr, error) {
	if err := p.consumeKeyword(KeywordWith); err != nil {
		return nil, err
	}

	cteExpr, err := p.parseCTEExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	ctes := []*CTEExpr{cteExpr}
	for p.tryConsumeTokenKind(",") != nil {
		cteExpr, err := p.parseCTEExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		ctes = append(ctes, cteExpr)
	}

	return &WithExpr{
		WithPos: pos,
		CTEs:    ctes,
		EndPos:  ctes[len(ctes)-1].End(),
	}, nil
}

func (p *Parser) tryParseTopExpr(pos Pos) (*TopExpr, error) {
	if !p.matchKeyword(KeywordTop) {
		return nil, nil
	}
	return p.parseTopExpr(pos)
}

func (p *Parser) parseTopExpr(pos Pos) (*TopExpr, error) {
	if err := p.consumeKeyword(KeywordTop); err != nil {
		return nil, err
	}

	number, err := p.parseNumber(p.Pos())
	if err != nil {
		return nil, err
	}
	topEnd := number.End()

	withTies := false
	if p.tryConsumeKeyword(KeywordWith) != nil {
		topEnd = p.last().End
		if err := p.consumeKeyword(KeywordTies); err != nil {
			return nil, err
		}
		withTies = true
	}
	return &TopExpr{
		TopPos:   pos,
		TopEnd:   topEnd,
		Number:   number,
		WithTies: withTies,
	}, nil
}

func (p *Parser) tryParseFromExpr(pos Pos) (*FromExpr, error) {
	if !p.matchKeyword(KeywordFrom) {
		return nil, nil
	}
	return p.parseFromExpr(pos)
}

func (p *Parser) parseFromExpr(pos Pos) (*FromExpr, error) {
	if err := p.consumeKeyword(KeywordFrom); err != nil {
		return nil, err
	}

	expr, err := p.parseJoinExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &FromExpr{
		FromPos: pos,
		Expr:    expr,
	}, nil
}

func (p *Parser) tryParseJoinConstraints(pos Pos) (Expr, error) {
	switch {
	case p.tryConsumeKeyword(KeywordOn) != nil:
		columnExprList, err := p.parseColumnExprList(p.Pos())
		if err != nil {
			return nil, err
		}
		return &OnExpr{
			OnPos: pos,
			On:    columnExprList,
		}, nil
	case p.tryConsumeKeyword(KeywordUsing) != nil:
		hasParen := p.tryConsumeTokenKind("(") != nil
		columnExprList, err := p.parseColumnExprListWithRoundBracket(p.Pos())
		if err != nil {
			return nil, err
		}
		if hasParen {
			if _, err := p.consumeTokenKind(")"); err != nil {
				return nil, err
			}
		}
		return &UsingExpr{
			UsingPos: pos,
			Using:    columnExprList,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseJoinOp(_ Pos) (Expr, []string, error) {
	var modifiers []string
	switch {
	case p.tryConsumeKeyword(KeywordCross) != nil: // cross join
		modifiers = append(modifiers, KeywordCross)
	case p.tryConsumeTokenKind(",") != nil:
		expr, err := p.parseJoinExpr(p.Pos())
		if err != nil {
			return nil, nil, err
		}
		return expr, nil, nil
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

	if p.tryConsumeKeyword(KeywordJoin) != nil {
		modifiers = append(modifiers, KeywordJoin)
		expr, err := p.parseJoinExpr(p.Pos())
		if err != nil {
			return nil, nil, err
		}
		return expr, modifiers, nil
	}
	return nil, modifiers, nil
}

func (p *Parser) parseJoinExpr(pos Pos) (expr Expr, err error) {
	var sampleRatio *SampleRatioExpr
	switch {
	case p.matchTokenKind(TokenString), p.matchTokenKind(TokenIdent), p.matchTokenKind("("):
		expr, err = p.parseTableExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		_ = p.tryConsumeKeyword(KeywordFinal)

		sampleRatio, err = p.tryParseSampleRationExpr(p.Pos())
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("expected table name or subquery, got %s", p.last().Kind)
	}

	// TODO: store global/local in AST
	if p.matchKeyword(KeywordGlobal) || p.matchKeyword(KeywordLocal) {
		_ = p.lexer.consumeToken()
	}
	rightExpr, modifiers, err := p.parseJoinOp(p.Pos())
	if err != nil {
		return nil, err
	}
	// rightExpr is nil means no join op
	if rightExpr == nil {
		return
	}
	constrains, err := p.tryParseJoinConstraints(p.Pos())
	if err != nil {
		return nil, err
	}
	return &JoinExpr{
		JoinPos:     pos,
		Left:        expr,
		Right:       rightExpr,
		Modifiers:   modifiers,
		SampleRatio: sampleRatio,
		Constraints: constrains,
	}, nil
}

func (p *Parser) parseTableExpr(pos Pos) (*TableExpr, error) {
	var expr Expr
	var err error
	switch {
	case p.matchTokenKind(TokenString):
		expr, err = p.parseString(p.Pos())
	case p.matchTokenKind(TokenIdent):
		// table name
		tableIdentifier, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		// it's a table name
		if tableIdentifier.Database != nil || !p.matchTokenKind("(") { // database.table
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
	case p.matchTokenKind("("):
		expr, err = p.parseSelectQuery(p.Pos())
	default:
		return nil, errors.New("expect table name or subquery")
	}
	if err != nil {
		return nil, err
	}

	tableEnd := expr.End()
	if asToken := p.tryConsumeKeyword(KeywordAs); asToken != nil {
		alias, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		expr = &AliasExpr{
			Expr:     expr,
			AliasPos: asToken.Pos,
			Alias:    alias,
		}
		tableEnd = expr.End()
	}

	isFinalExist := false
	if asToken := p.tryConsumeKeyword(KeywordFinal); asToken != nil {
		switch expr.(type) {
		case *TableFunctionExpr:
			return nil, errors.New("tablefunction doesn't support FINAL")
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

func (p *Parser) tryParsePrewhereExpr(pos Pos) (*PrewhereExpr, error) {
	if !p.matchKeyword(KeywordPrewhere) {
		return nil, nil
	}
	return p.parsePrewhereExpr(pos)
}
func (p *Parser) parsePrewhereExpr(pos Pos) (*PrewhereExpr, error) {
	if err := p.consumeKeyword(KeywordPrewhere); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &PrewhereExpr{
		PrewherePos: pos,
		Expr:        expr,
	}, nil
}

func (p *Parser) tryParseWhereExpr(pos Pos) (*WhereExpr, error) {
	if !p.matchKeyword(KeywordWhere) {
		return nil, nil
	}
	return p.parseWhereExpr(pos)
}

func (p *Parser) parseWhereExpr(pos Pos) (*WhereExpr, error) {
	if err := p.consumeKeyword(KeywordWhere); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &WhereExpr{
		WherePos: pos,
		Expr:     expr,
	}, nil
}

func (p *Parser) tryParseGroupByExpr(pos Pos) (*GroupByExpr, error) {
	if !p.matchKeyword(KeywordGroup) {
		return nil, nil
	}
	return p.parseGroupByExpr(pos)
}

// syntax: groupByClause? (WITH (CUBE | ROLLUP))? (WITH TOTALS)?
func (p *Parser) parseGroupByExpr(pos Pos) (*GroupByExpr, error) {
	if err := p.consumeKeyword(KeywordGroup); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordBy); err != nil {
		return nil, err
	}

	var expr Expr
	var err error
	aggregateType := ""
	if p.matchKeyword(KeywordCube) || p.matchKeyword(KeywordRollup) {
		aggregateType = p.last().String
		_ = p.lexer.consumeToken()
		expr, err = p.parseFunctionParams(p.Pos())
	} else {
		expr, err = p.parseColumnExprListWithRoundBracket(p.Pos())
	}
	if err != nil {
		return nil, err
	}

	groupByExpr := &GroupByExpr{
		GroupByPos:    pos,
		AggregateType: aggregateType,
		Expr:          expr,
	}

	// parse WITH CUBE, ROLLUP, TOTALS
	for p.tryConsumeKeyword(KeywordWith) != nil {
		switch {
		case p.tryConsumeKeyword(KeywordCube) != nil:
			groupByExpr.WithCube = true
		case p.tryConsumeKeyword(KeywordRollup) != nil:
			groupByExpr.WithRollup = true
		case p.tryConsumeKeyword(KeywordTotals) != nil:
			groupByExpr.WithTotals = true
		default:
			return nil, fmt.Errorf("expected CUBE, ROLLUP or TOTALS, got %s", p.lastTokenKind())
		}
	}

	return groupByExpr, nil
}

func (p *Parser) tryParseLimitExpr(pos Pos) (*LimitExpr, error) {
	if !p.matchKeyword(KeywordLimit) {
		return nil, nil
	}
	return p.parseLimitExpr(pos)
}

func (p *Parser) parseLimitExpr(pos Pos) (*LimitExpr, error) {
	if err := p.consumeKeyword(KeywordLimit); err != nil {
		return nil, err
	}

	limit, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	var offset Expr
	if p.tryConsumeKeyword(KeywordOffset) != nil {
		offset, err = p.parseExpr(p.Pos())
	} else if p.tryConsumeTokenKind(",") != nil {
		offset = limit
		limit, err = p.parseExpr(p.Pos())
	}
	if err != nil {
		return nil, err
	}

	return &LimitExpr{
		LimitPos: pos,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (p *Parser) tryParseLimitByExpr(pos Pos) (Expr, error) {
	if !p.matchKeyword(KeywordLimit) {
		return nil, nil
	}
	return p.parseLimitByExpr(pos)
}

func (p *Parser) parseLimitByExpr(pos Pos) (Expr, error) {
	limitExpr, err := p.parseLimitExpr(pos)
	if err != nil {
		return nil, err
	}

	var by *ColumnExprList
	if p.tryConsumeKeyword(KeywordBy) == nil {
		return limitExpr, nil
	}
	if by, err = p.parseColumnExprListWithRoundBracket(p.Pos()); err != nil {
		return nil, err
	}
	return &LimitByExpr{
		Limit:  limitExpr,
		ByExpr: by,
	}, nil
}

func (p *Parser) tryParseWindowFrameExpr(pos Pos) (*WindowFrameExpr, error) {
	if !p.matchKeyword(KeywordRows) && !p.matchKeyword(KeywordRange) {
		return nil, nil
	}
	return p.parseWindowFrameExpr(pos)
}

func (p *Parser) parseWindowFrameExpr(pos Pos) (*WindowFrameExpr, error) {
	var windowFrameType string
	if p.matchKeyword(KeywordRows) || p.matchKeyword(KeywordRange) {
		windowFrameType = p.last().String
		_ = p.lexer.consumeToken()
	}

	var expr Expr
	switch {
	case p.tryConsumeKeyword(KeywordBetween) != nil:
		betweenExpr, err := p.parseWindowFrameExpr(p.Pos())
		if err != nil {
			return nil, err
		}

		andPos := p.Pos()
		if err := p.consumeKeyword(KeywordAnd); err != nil {
			return nil, err
		}

		andExpr, err := p.parseWindowFrameExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		expr = &WindowFrameRangeExpr{
			BetweenPos:  pos,
			BetweenExpr: betweenExpr,
			AndPos:      andPos,
			AndExpr:     andExpr,
		}
	case p.matchKeyword(KeywordCurrent):
		currentPos := p.Pos()
		_ = p.lexer.consumeToken()
		rowEnd := p.last().End
		if err := p.consumeKeyword(KeywordRow); err != nil {
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
	case p.matchTokenKind(TokenInt):
		number, err := p.parseNumber(p.Pos())
		if err != nil {
			return nil, err
		}

		var unboundedEnd Pos
		direction := ""
		switch {
		case p.matchKeyword(KeywordPreceding), p.matchKeyword(KeywordFollowing):
			direction = p.last().String
			unboundedEnd = p.last().End
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
	return &WindowFrameExpr{
		FramePos: pos,
		Type:     windowFrameType,
		Extend:   expr,
	}, nil
}

func (p *Parser) tryParseWindowExpr(pos Pos) (*WindowExpr, error) {
	if !p.matchKeyword(KeywordWindow) {
		return nil, nil
	}
	return p.parseWindowExpr(pos)
}

func (p *Parser) parseWindowCondition(pos Pos) (*WindowConditionExpr, error) {
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}
	partitionBy, err := p.tryParsePartitionByExpr(pos)
	if err != nil {
		return nil, err
	}
	orderBy, err := p.tryParseOrderByExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	frame, err := p.tryParseWindowFrameExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return &WindowConditionExpr{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		PartitionBy:   partitionBy,
		OrderBy:       orderBy,
		Frame:         frame,
	}, nil
}

func (p *Parser) parseWindowExpr(pos Pos) (*WindowExpr, error) {
	if err := p.consumeKeyword(KeywordWindow); err != nil {
		return nil, err
	}

	windowName, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	if err := p.consumeKeyword(KeywordAs); err != nil {
		return nil, err
	}

	condition, err := p.parseWindowCondition(p.Pos())
	if err != nil {
		return nil, err
	}

	return &WindowExpr{
		WindowPos:           pos,
		Name:                windowName,
		WindowConditionExpr: condition,
	}, nil
}

func (p *Parser) tryParseArrayJoin(pos Pos) (*ArrayJoinExpr, error) {
	if !p.matchKeyword(KeywordLeft) && !p.matchKeyword(KeywordInner) && !p.matchKeyword(KeywordArray) {
		return nil, nil
	}
	return p.parseArrayJoin(pos)
}

func (p *Parser) parseArrayJoin(_ Pos) (*ArrayJoinExpr, error) {
	var typ string
	switch {
	case p.matchKeyword(KeywordLeft), p.matchKeyword(KeywordInner):
		typ = p.last().String
		_ = p.lexer.consumeToken()
	}
	arrayPos := p.Pos()
	if err := p.consumeKeyword(KeywordArray); err != nil {
		return nil, err
	}

	if err := p.consumeKeyword(KeywordJoin); err != nil {
		return nil, err
	}

	expr, err := p.parseColumnExprList(p.Pos())
	if err != nil {
		return nil, err
	}

	return &ArrayJoinExpr{
		ArrayPos: arrayPos,
		Type:     typ,
		Expr:     expr,
	}, nil
}

func (p *Parser) tryParseHavingExpr(pos Pos) (*HavingExpr, error) {
	if !p.matchKeyword(KeywordHaving) {
		return nil, nil
	}
	return p.parseHavingExpr(pos)
}

func (p *Parser) parseHavingExpr(pos Pos) (*HavingExpr, error) {
	if err := p.consumeKeyword(KeywordHaving); err != nil {
		return nil, err
	}

	expr, err := p.parseColumnsExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	return &HavingExpr{
		HavingPos: pos,
		Expr:      expr,
	}, nil
}

func (p *Parser) parseSubQuery(pos Pos) (*SubQueryExpr, error) {
	if err := p.consumeKeyword(KeywordAs); err != nil {
		return nil, err
	}

	selectExprList, err := p.parseSelectQuery(p.Pos())
	if err != nil {
		return nil, err
	}

	return &SubQueryExpr{
		AsPos:  pos,
		Select: selectExprList,
	}, nil
}

func (p *Parser) parseSelectQuery(_ Pos) (*SelectQuery, error) {
	if !p.matchKeyword(KeywordSelect) && !p.matchKeyword(KeywordWith) && !p.matchTokenKind("(") {
		return nil, fmt.Errorf("expected SELECT, WITH or (, got %s", p.lastTokenKind())
	}

	hasParen := p.tryConsumeTokenKind("(") != nil
	selectExpr, err := p.parseSelectStatement(p.Pos())
	if err != nil {
		return nil, err
	}
	switch {
	case p.tryConsumeKeyword(KeywordUnion) != nil:
		switch {
		case p.tryConsumeKeyword(KeywordAll) != nil:
			unionAllExpr, err := p.parseSelectStatement(p.Pos())
			if err != nil {
				return nil, err
			}
			selectExpr.UnionAll = unionAllExpr
		case p.tryConsumeKeyword(KeywordDistinct) != nil:
			unionDistinctExpr, err := p.parseSelectStatement(p.Pos())
			if err != nil {
				return nil, err
			}
			selectExpr.UnionDistinct = unionDistinctExpr
		default:
			return nil, fmt.Errorf("expected ALL or DISTINCT, got %s", p.lastTokenKind())
		}
	case p.tryConsumeKeyword(KeywordExcept) != nil:
		exceptExpr, err := p.parseSelectStatement(p.Pos())
		if err != nil {
			return nil, err
		}
		selectExpr.Except = exceptExpr
	}
	if hasParen {
		if _, err := p.consumeTokenKind(")"); err != nil {
			return nil, err
		}
	}
	return selectExpr, nil
}

func (p *Parser) parseSelectStatement(pos Pos) (*SelectQuery, error) { // nolint: funlen
	withExpr, err := p.tryParseWithExpr(pos)
	if err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordSelect); err != nil {
		return nil, err
	}
	// DISTINCT?
	_ = p.tryConsumeKeyword(KeywordDistinct)

	topExpr, err := p.tryParseTopExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	selectColumns, err := p.parseColumnExprListWithRoundBracket(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := selectColumns.End()
	fromExpr, err := p.tryParseFromExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	if fromExpr != nil {
		statementEnd = fromExpr.End()
	}
	arrayJoinExpr, err := p.tryParseArrayJoin(p.Pos())
	if err != nil {
		return nil, err
	}
	if arrayJoinExpr != nil {
		statementEnd = arrayJoinExpr.End()
	}
	windowExpr, err := p.tryParseWindowExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if windowExpr != nil {
		statementEnd = windowExpr.End()
	}
	prewhereExpr, err := p.tryParsePrewhereExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if prewhereExpr != nil {
		statementEnd = prewhereExpr.End()
	}
	whereExpr, err := p.tryParseWhereExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if whereExpr != nil {
		statementEnd = whereExpr.End()
	}
	groupByExpr, err := p.tryParseGroupByExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if groupByExpr != nil {
		statementEnd = groupByExpr.End()
	}

	withTotal := false
	lastPos := p.Pos()
	if p.tryConsumeKeyword(KeywordWith) != nil {
		if err := p.consumeKeyword(KeywordTotals); err != nil {
			return nil, err
		}
		withTotal = true
		statementEnd = lastPos
	}
	havingExpr, err := p.tryParseHavingExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if havingExpr != nil {
		statementEnd = havingExpr.End()
	}
	orderByExpr, err := p.tryParseOrderByExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	if orderByExpr != nil {
		statementEnd = orderByExpr.End()
	}

	var limitByExpr *LimitByExpr
	var limitExpr *LimitExpr
	parsedLimitBy, err := p.tryParseLimitByExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if parsedLimitBy != nil {
		statementEnd = parsedLimitBy.End()
		switch e := parsedLimitBy.(type) {
		case *LimitByExpr:
			limitByExpr = e
			limitExpr, err = p.tryParseLimitExpr(p.Pos())
			if err != nil {
				return nil, err
			}
			if limitExpr != nil {
				statementEnd = limitExpr.End()
			}
		case *LimitExpr:
			limitExpr = e
		}
	}

	settingsExpr, err := p.tryParseSettingsExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	if settingsExpr != nil {
		statementEnd = settingsExpr.End()
	}

	return &SelectQuery{
		With:          withExpr,
		SelectPos:     pos,
		StatementEnd:  statementEnd,
		Top:           topExpr,
		SelectColumns: selectColumns,
		From:          fromExpr,
		ArrayJoin:     arrayJoinExpr,
		Window:        windowExpr,
		Prewhere:      prewhereExpr,
		Where:         whereExpr,
		GroupBy:       groupByExpr,
		Having:        havingExpr,
		OrderBy:       orderByExpr,
		LimitBy:       limitByExpr,
		Limit:         limitExpr,
		Settings:      settingsExpr,
		WithTotal:     withTotal,
	}, nil
}

func (p *Parser) parseCTEExpr(pos Pos) (*CTEExpr, error) {
	expr, err := p.parseOrExpr(pos)
	if err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordAs); err != nil {
		return nil, err
	}
	if p.matchTokenKind("(") {
		selectQuery, err := p.parseSelectQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		return &CTEExpr{
			CTEPos: pos,
			Expr:   expr,
			Alias:  selectQuery,
		}, nil
	}
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	return &CTEExpr{
		CTEPos: pos,
		Expr:   expr,
		Alias:  name,
	}, nil
}

func (p *Parser) tryParseColumnAliases() ([]*Ident, error) {
	if !p.matchTokenKind("(") {
		return nil, nil
	}
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}

	aliasList := make([]*Ident, 0)
	for {
		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		aliasList = append(aliasList, ident)
		if p.matchTokenKind(")") {
			break
		}
		if _, err := p.consumeTokenKind(","); err != nil {
			return nil, err
		}
	}
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return aliasList, nil
}

func (p *Parser) tryParseSampleRationExpr(pos Pos) (*SampleRatioExpr, error) {
	if !p.matchKeyword(KeywordSample) {
		return nil, nil
	}
	return p.parseSampleRationExpr(pos)
}

func (p *Parser) parseSampleRationExpr(pos Pos) (*SampleRatioExpr, error) {
	if err := p.consumeKeyword(KeywordSample); err != nil {
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

	return &SampleRatioExpr{
		SamplePos: pos,
		Ratio:     ratio,
		Offset:    offset,
	}, nil
}

func (p *Parser) parseExplainExpr(pos Pos) (*ExplainExpr, error) {
	if err := p.consumeKeyword(KeywordExplain); err != nil {
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
	expr, err := p.parseSelectQuery(p.Pos())
	if err != nil {
		return nil, err
	}
	return &ExplainExpr{
		ExplainPos: pos,
		Type:       explainType,
		Statement:  expr,
	}, nil
}
