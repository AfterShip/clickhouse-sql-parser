package parser

import (
	"fmt"
)

func (p *Parser) parseDDL(pos Pos) (DDL, error) {
	switch {
	case p.matchKeyword(KeywordCreate),
		p.matchKeyword(KeywordAttach):
		_ = p.lexer.consumeToken()
		switch {
		case p.matchKeyword(KeywordDatabase):
			return p.parseCreateDatabase(pos)
		case p.matchKeyword(KeywordTable),
			p.matchKeyword(KeywordTemporary):
			return p.parseCreateTable(pos)
		case p.matchKeyword(KeywordFunction):
			return p.parseCreateFunction(pos)
		case p.matchKeyword(KeywordMaterialized):
			return p.parseCreateMaterializedView(pos)
		case p.matchKeyword(KeywordLive):
			return p.parseCreateLiveView(pos)
		case p.matchKeyword(KeywordView):
			return p.parseCreateView(pos)
		case p.matchKeyword(KeywordRole):
			return p.parseCreateRole(pos)
		case p.matchKeyword(KeywordDictionary):
		case p.matchKeyword(KeywordFunction):
		case p.matchKeyword(KeywordRow):
		case p.matchKeyword(KeywordSettings):
		default:
			return nil, fmt.Errorf("expected keyword: DATABASE|TABLE|VIEW|DICTIONARY|FUNCTION|ROW|QUOTA|SETTINGS, but got %q",
				p.last().String)
		}
	case p.matchKeyword(KeywordAlter):
		_ = p.lexer.consumeToken()
		switch {
		case p.matchKeyword(KeywordRole):
			return p.parseAlterRole(pos)
		case p.matchKeyword(KeywordTable):
			return p.parseAlterTable(pos)
		default:
			return nil, fmt.Errorf("expected keyword: TABLE|ROLE, but got %q", p.last().String)
		}
	case p.matchKeyword(KeywordDrop),
		p.matchKeyword(KeywordDetach):
		_ = p.lexer.consumeToken()
		switch {
		case p.matchKeyword(KeywordDatabase):
			return p.parseDropDatabase(pos)
		case p.matchKeyword(KeywordTemporary),
			p.matchKeyword(KeywordView),
			p.matchKeyword(KeywordDictionary),
			p.matchKeyword(KeywordTable):
			return p.parseDropStmt(pos)
		case p.matchKeyword(KeywordUser),
			p.matchKeyword(KeywordRole):
			return p.parserDropUserOrRole(pos)
		default:
			return nil, fmt.Errorf("expected keyword: DATABASE|TABLE, but got %q", p.last().String)
		}
	case p.matchKeyword(KeywordTruncate):
		return p.parseTruncateTable(pos)
	case p.matchKeyword(KeywordRename):
		return p.parseRenameStmt(pos)
	}
	return nil, nil // nolint
}

func (p *Parser) parseCreateDatabase(pos Pos) (*CreateDatabase, error) {
	if err := p.consumeKeyword(KeywordDatabase); err != nil {
		return nil, err
	}

	// try to parse IF NOT EXISTS clause
	ifNotExists, err := p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}
	// parse database name
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	StatementEnd := name.End()
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	if onCluster != nil {
		StatementEnd = onCluster.End()
	}
	engineExpr, err := p.tryParseEngineExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if engineExpr != nil {
		StatementEnd = onCluster.End()
	}
	return &CreateDatabase{
		CreatePos:    pos,
		StatementEnd: StatementEnd,
		Name:         name,
		IfNotExists:  ifNotExists,
		OnCluster:    onCluster,
		Engine:       engineExpr,
	}, nil
}

func (p *Parser) parseCreateTable(pos Pos) (*CreateTable, error) {
	createTable := &CreateTable{CreatePos: pos}

	createTable.HasTemporary = p.tryConsumeKeyword(KeywordTemporary) != nil

	if err := p.consumeKeyword(KeywordTable); err != nil {
		return nil, err
	}

	// parse IF NOT EXISTS clause if exists
	var err error
	createTable.IfNotExists, err = p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}

	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	createTable.Name = tableIdentifier

	// try parse UUID clause if exists
	uuid, err := p.tryParseUUID()
	if err != nil {
		return nil, err
	}
	createTable.UUID = uuid
	// parse ON CLUSTER clause if exists
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	createTable.OnCluster = onCluster

	tableSchema, err := p.parseTableSchemaExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	createTable.TableSchema = tableSchema

	engineExpr, err := p.tryParseEngineExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if engineExpr != nil {
		createTable.Engine = engineExpr
		createTable.StatementEnd = engineExpr.End()
	}

	if p.matchKeyword(KeywordAs) {
		subQuery, err := p.parseSubQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		createTable.SubQuery = subQuery
		createTable.StatementEnd = subQuery.End()
	}
	return createTable, nil
}

func (p *Parser) parseIdentOrFunction(_ Pos) (Expr, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	switch {
	case p.matchTokenKind("["):
		params, err := p.parseArrayParams(p.Pos())
		if err != nil {
			return nil, err
		}
		return &ObjectParams{
			Object: ident,
			Params: params,
		}, nil
	case p.matchTokenKind("("):
		params, err := p.parseFunctionParams(p.Pos())
		if err != nil {
			return nil, err
		}
		funcExpr := &FunctionExpr{
			Name:   ident,
			Params: params,
		}
		if overToken := p.tryConsumeKeyword(KeywordOver); overToken != nil {
			var overExpr Expr
			switch {
			case p.matchTokenKind(TokenIdent):
				overExpr, err = p.parseIdent()
			case p.matchTokenKind("("):
				overExpr, err = p.parseWindowCondition(p.Pos())
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("expected IDENT or (, but got %q", p.lastTokenKind())
			}

			if err != nil {
				return nil, err
			}
			return &WindowFunctionExpr{
				Function: funcExpr,
				OverPos:  overToken.Pos,
				OverExpr: overExpr,
			}, nil
		}
		return funcExpr, nil
	case p.tryConsumeTokenKind(".") != nil:
		switch {
		case p.matchTokenKind(TokenIdent):
			nextIdent, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			return &TableIdentifier{
				Database: ident,
				Table:    nextIdent,
			}, nil
		case p.matchTokenKind("*"):
			nextIdent, err := p.parseColumnStar(p.Pos())
			if err != nil {
				return nil, err
			}
			return &NestedIdentifier{
				Ident:    ident,
				DotIdent: nextIdent,
			}, nil
		}

	}
	return ident, nil
}

func (p *Parser) parseTableIdentifier(_ Pos) (*TableIdentifier, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	dotIdent, err := p.tryParseDotIdent()
	if err != nil {
		return nil, err
	}
	if dotIdent != nil {
		return &TableIdentifier{
			Database: ident,
			Table:    dotIdent,
		}, nil
	}
	return &TableIdentifier{
		Table: ident,
	}, nil
}

func (p *Parser) parseTableSchemaExpr(pos Pos) (*TableSchemaExpr, error) {
	switch {
	case p.matchTokenKind("("):
		// parse column definitions
		if _, err := p.consumeTokenKind("("); err != nil {
			return nil, err
		}

		columns, err := p.parseTableColumns()
		if err != nil {
			return nil, err
		}

		rightParenPos := p.Pos()
		if _, err := p.consumeTokenKind(")"); err != nil {
			return nil, err
		}
		return &TableSchemaExpr{
			SchemaPos: pos,
			SchemaEnd: rightParenPos,
			Columns:   columns,
		}, nil
	case p.tryConsumeKeyword(KeywordAs) != nil:
		switch {
		case p.matchTokenKind(TokenIdent):
			ident, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			switch {
			case p.matchTokenKind("."):
				// it's a database.table
				dotIdent, err := p.tryParseDotIdent()
				if err != nil {
					return nil, err
				}
				return &TableSchemaExpr{
					SchemaPos: pos,
					SchemaEnd: dotIdent.NameEnd,
					AliasTable: &TableIdentifier{
						Database: ident,
						Table:    dotIdent,
					},
				}, nil
			case p.matchTokenKind("("):
				// it's a table function
				argsExpr, err := p.parseTableArgList(pos)
				if err != nil {
					return nil, err
				}
				return &TableSchemaExpr{
					SchemaPos: pos,
					SchemaEnd: p.last().End,
					TableFunction: &TableFunctionExpr{
						Name: ident,
						Args: argsExpr,
					},
				}, nil
			default:
				return &TableSchemaExpr{
					SchemaPos: pos,
					SchemaEnd: p.last().End,
					AliasTable: &TableIdentifier{
						Table: ident,
					},
				}, nil
			}
		}
	}
	// no schema is ok for MATERIALIZED VIEW
	return nil, nil
}

func (p *Parser) parseTableColumns() ([]Expr, error) {
	columns := make([]Expr, 0)
	for !p.lexer.isEOF() {
		switch {
		case p.matchKeyword(KeywordIndex):
			indexPos := p.Pos()
			_ = p.lexer.consumeToken()
			index, err := p.parseTableIndex(indexPos)
			if err != nil {
				return nil, err
			}
			columns = append(columns, index)
		case p.matchKeyword(KeywordConstraint):
			constraintPos := p.Pos()
			_ = p.lexer.consumeToken()
			ident, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			if err := p.consumeKeyword(KeywordCheck); err != nil {
				return nil, err
			}
			expr, err := p.parseExpr(p.Pos())
			if err != nil {
				return nil, err
			}
			columns = append(columns, &ConstraintExpr{
				ConstraintPos: constraintPos,
				Constraint:    ident,
				Expr:          expr,
			})
		default:
			column, err := p.tryParseTableColumn(p.Pos())
			if err != nil {
				return nil, err
			}
			if column == nil {
				break
			}
			columns = append(columns, column)
		}
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	// end of column definitions
	return columns, nil
}

func (p *Parser) tryParseTableColumn(pos Pos) (*Column, error) {
	if !p.matchTokenKind(TokenIdent) {
		return nil, nil // nolint
	}
	return p.parseTableColumn(pos)
}

func (p *Parser) parseTableColumn(pos Pos) (*Column, error) {
	// Not a column definition, just return
	column := &Column{NamePos: pos}
	// parse column name
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	column.Name = name
	columnEnd := name.End()

	if p.matchTokenKind(TokenIdent) && !p.matchKeyword(KeywordRemove) {
		columnType, err := p.parseColumnType(p.Pos())
		if err != nil {
			return nil, err
		}
		column.Type = columnType
		columnEnd = columnType.End()
	}

	nullable := p.tryParseNull(p.Pos())
	if nullable != nil {
		columnEnd = nullable.End()
	}
	notNull, err := p.tryParseNotNull(p.Pos())
	if err != nil {
		return nil, err
	}
	if notNull != nil {
		columnEnd = notNull.End()
	}
	comment, err := p.tryParseColumnComment(p.Pos())
	if err != nil {
		return nil, err
	}
	if comment != nil {
		columnEnd = comment.End()
	}
	property, err := p.tryParseTableColumnPropertyExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if property != nil {
		columnEnd = property.End()
	}
	codec, err := p.tryParseCompressionCodecs(p.Pos())
	if err != nil {
		return nil, err
	}
	if codec != nil {
		columnEnd = codec.End()
	}

	column.ColumnEnd = columnEnd
	column.Comment = comment
	column.Codec = codec
	column.Nullable = nullable
	column.NotNull = notNull
	column.Property = property
	return column, nil
}

func (p *Parser) parseTableArgExpr(pos Pos) (Expr, error) {
	switch {
	case p.matchTokenKind(TokenIdent):
		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		switch {
		// nest identifier
		case p.matchTokenKind("."):
			dotIdent, err := p.tryParseDotIdent()
			if err != nil {
				return nil, err
			}
			return &NestedIdentifier{
				Ident:    ident,
				DotIdent: dotIdent,
			}, nil
		case p.matchTokenKind("("):
			argsExpr, err := p.parseTableArgList(pos)
			if err != nil {
				return nil, err
			}
			return &TableFunctionExpr{
				Name: ident,
				Args: argsExpr,
			}, nil
		default:
			return ident, nil
		}
	case p.matchTokenKind(TokenInt), p.matchTokenKind(TokenString), p.matchKeyword("NULL"):
		return p.parseLiteral(p.Pos())
	default:
		return nil, fmt.Errorf("unexpected token: %q, expected <Ident>, <literal>", p.last().String)
	}
}

func (p *Parser) parseTableArgList(pos Pos) (*TableArgListExpr, error) {
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}

	args := make([]Expr, 0)
	for !p.lexer.isEOF() {
		arg, err := p.parseTableArgExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}

	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}

	return &TableArgListExpr{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Args:          args,
	}, nil
}

func (p *Parser) tryParseOnCluster(pos Pos) (*OnClusterExpr, error) {
	if p.tryConsumeKeyword(KeywordOn) == nil {
		return nil, nil // nolint
	}
	if err := p.consumeKeyword(KeywordCluster); err != nil {
		return nil, err
	}

	var expr Expr
	var err error
	switch {
	case p.matchTokenKind(TokenIdent):
		expr, err = p.parseIdent()
	case p.matchTokenKind(TokenString):
		expr, err = p.parseString(p.Pos())
	default:
		return nil, fmt.Errorf("expected <Ident> or <Literal>, but got %q", p.last().String)
	}
	if err != nil {
		return nil, err
	}
	return &OnClusterExpr{
		OnPos: pos,
		Expr:  expr,
	}, nil
}

func (p *Parser) tryParsePartitionByExpr(pos Pos) (*PartitionByExpr, error) {
	if p.tryConsumeKeyword(KeywordPartition) == nil {
		return nil, nil // nolint
	}

	if err := p.consumeKeyword(KeywordBy); err != nil {
		return nil, err
	}

	// parse partition key list
	columnExpr, err := p.parseColumnExprListWithRoundBracket(p.Pos())
	if err != nil {
		return nil, err
	}
	return &PartitionByExpr{
		PartitionPos: pos,
		Expr:         columnExpr,
	}, nil
}

func (p *Parser) tryParsePrimaryKeyExpr(pos Pos) (*PrimaryKeyExpr, error) {
	if p.tryConsumeKeyword(KeywordPrimary) == nil {
		return nil, nil // nolint
	}

	if err := p.consumeKeyword(KeywordKey); err != nil {
		return nil, err
	}

	// parse partition key list
	columnExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &PrimaryKeyExpr{
		PrimaryPos: pos,
		Expr:       columnExpr,
	}, nil
}

func (p *Parser) tryParseOrderByExprList(pos Pos) (*OrderByListExpr, error) {
	if p.tryConsumeKeyword(KeywordOrder) == nil {
		return nil, nil // nolint
	}

	if err := p.consumeKeyword(KeywordBy); err != nil {
		return nil, err
	}
	return p.parseOrderByExprList(pos)
}

func (p *Parser) parseOrderByExprList(pos Pos) (*OrderByListExpr, error) {
	orderByListExpr := &OrderByListExpr{OrderPos: pos, ListEnd: pos}
	items := make([]Expr, 0)
	for {
		expr, err := p.parseOrderByExpr(pos)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			break
		}
		items = append(items, expr)

		if p.lexer.isEOF() || p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	if len(items) > 0 {
		orderByListExpr.ListEnd = items[len(items)-1].End()
	}
	orderByListExpr.Items = items
	return orderByListExpr, nil
}

func (p *Parser) parseOrderByExpr(pos Pos) (*OrderByExpr, error) {
	// parse column expr
	columnExpr, err := p.parseExpr(pos)
	if err != nil {
		return nil, err
	}

	direction := OrderDirectionNone
	switch {
	case p.matchKeyword(KeywordAsc), p.matchKeyword(KeywordAscending):
		direction = OrderDirectionAsc
		_ = p.lexer.consumeToken()
	case p.matchKeyword(KeywordDesc), p.matchKeyword(KeywordDescending):
		direction = OrderDirectionDesc
		_ = p.lexer.consumeToken()
	}
	return &OrderByExpr{
		OrderPos:  pos,
		Expr:      columnExpr,
		Direction: direction,
	}, nil
}

func (p *Parser) tryParseTTLExprList(pos Pos) (*TTLExprList, error) {
	if p.tryConsumeKeyword(KeywordTtl) == nil {
		return nil, nil // nolint
	}
	ttlExprList := &TTLExprList{TTLPos: pos, ListEnd: pos}
	// accept the TTL keyword
	items, err := p.parseTTLExprList(pos)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		ttlExprList.ListEnd = items[len(items)-1].End()
	}
	ttlExprList.Items = items
	return ttlExprList, nil
}

func (p *Parser) parseTTLExprList(pos Pos) ([]*TTLExpr, error) {
	items := make([]*TTLExpr, 0)
	expr, err := p.parseTTLExpr(pos)
	if err != nil {
		return nil, err
	}
	items = append(items, expr)
	for !p.lexer.isEOF() && p.tryConsumeTokenKind(",") != nil {
		expr, err = p.parseTTLExpr(pos)
		if err != nil {
			return nil, err
		}
		items = append(items, expr)
	}
	return items, nil
}

func (p *Parser) parseTTLExpr(pos Pos) (*TTLExpr, error) {
	columnExpr, err := p.parseExpr(pos)
	if err != nil {
		return nil, err
	}
	switch {
	case p.matchKeyword(KeywordDelete):
		_ = p.lexer.consumeToken()
	case p.matchKeyword(KeywordTo):
		_ = p.lexer.consumeToken()
		if p.tryConsumeKeyword(KeywordDisk) != nil || p.tryConsumeKeyword(KeywordVolume) != nil {
			_, err := p.parseString(p.Pos())
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("expected keyword <DISK> or <VOLUME>, but got %q", p.last().String)
		}
	}
	return &TTLExpr{
		TTLPos: pos,
		Expr:   columnExpr,
	}, nil
}

func (p *Parser) tryParseSampleByExpr(pos Pos) (*SampleByExpr, error) {
	if p.tryConsumeKeyword(KeywordSample) == nil {
		return nil, nil // nolint
	}

	if err := p.consumeKeyword(KeywordBy); err != nil {
		return nil, err
	}

	// parse sample by expr
	columnExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &SampleByExpr{
		SamplePos: pos,
		Expr:      columnExpr,
	}, nil
}

func (p *Parser) tryParseSettingsExprList(pos Pos) (*SettingsExprList, error) {
	if p.tryConsumeKeyword(KeywordSettings) == nil {
		return nil, nil // nolint
	}
	return p.parseSettingsExprList(pos)
}

func (p *Parser) parseSettingsExprList(pos Pos) (*SettingsExprList, error) {
	settingsExprList := &SettingsExprList{SettingsPos: pos, ListEnd: pos}
	items := make([]*SettingsExpr, 0)
	expr, err := p.parseSettingsExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	items = append(items, expr)
	for p.tryConsumeTokenKind(",") != nil {
		expr, err = p.parseSettingsExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		items = append(items, expr)
	}
	if len(items) > 0 {
		settingsExprList.ListEnd = items[len(items)-1].End()
	}
	settingsExprList.Items = items
	return settingsExprList, nil
}

func (p *Parser) parseSettingsExpr(pos Pos) (*SettingsExpr, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	if _, err := p.consumeTokenKind("="); err != nil {
		return nil, err
	}

	var expr Expr
	switch {
	case p.matchTokenKind(TokenInt):
		number, err := p.parseNumber(p.Pos())
		if err != nil {
			return nil, err
		}
		expr = number
	case p.matchTokenKind(TokenString):
		str, err := p.parseString(p.Pos())
		expr = str
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected token: %q, expected <number> or <string>", p.last().String)
	}

	return &SettingsExpr{
		SettingsPos: pos,
		Name:        ident,
		Expr:        expr,
	}, nil
}

func (p *Parser) parseDefaultExpr(pos Pos) (Expr, error) {
	if err := p.consumeKeyword(KeywordDefault); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr(pos)
	if err != nil {
		return nil, err
	}
	return &DefaultExpr{
		DefaultPos: pos,
		Expr:       expr,
	}, nil
}

func (p *Parser) parseDestinationExpr(pos Pos) (*DestinationExpr, error) {
	if err := p.consumeKeyword(KeywordTo); err != nil {
		return nil, err
	}

	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	return &DestinationExpr{
		ToPos:           pos,
		TableIdentifier: tableIdentifier,
	}, nil
}

func (p *Parser) tryParseEngineExpr(pos Pos) (*EngineExpr, error) {
	if !p.matchKeyword(KeywordEngine) {
		return nil, nil // nolint
	}
	return p.parseEngineExpr(pos)
}

func (p *Parser) parseEngineExpr(pos Pos) (*EngineExpr, error) {
	if err := p.consumeKeyword(KeywordEngine); err != nil {
		return nil, err
	}
	_ = p.tryConsumeTokenKind("=")

	engineExpr := &EngineExpr{EnginePos: pos}
	var engineEnd Pos
	switch {
	case p.matchKeyword(KeywordNull):
		engineExpr.Name = KeywordNull
		engineEnd = p.last().End
		_ = p.lexer.consumeToken()
	case p.matchTokenKind(TokenIdent):
		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		engineExpr.Name = ident.Name
		engineEnd = ident.End()
		if p.matchTokenKind("(") {
			params, err := p.parseFunctionParams(p.Pos())
			if err != nil {
				return nil, err
			}
			engineExpr.Params = params
			engineExpr.EngineEnd = params.End()
		}
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.lastTokenKind())
	}

	for !p.lexer.isEOF() {
		switch {
		case p.matchKeyword(KeywordOrder):
			orderByExprList, err := p.tryParseOrderByExprList(p.Pos())
			if err != nil {
				return nil, err
			}
			engineExpr.OrderByListExpr = orderByExprList
			engineEnd = orderByExprList.End()
		case p.matchKeyword(KeywordPartition):
			partitionByExpr, err := p.tryParsePartitionByExpr(p.Pos())
			if err != nil {
				return nil, err
			}
			engineExpr.PartitionBy = partitionByExpr
			engineEnd = partitionByExpr.End()
		case p.matchKeyword(KeywordPrimary):
			primaryKeyExpr, err := p.tryParsePrimaryKeyExpr(p.Pos())
			if err != nil {
				return nil, err
			}
			engineExpr.PrimaryKey = primaryKeyExpr
			engineEnd = primaryKeyExpr.End()
		case p.matchKeyword(KeywordSample):
			sampleByExpr, err := p.tryParseSampleByExpr(p.Pos())
			if err != nil {
				return nil, err
			}
			engineExpr.SampleBy = sampleByExpr
			engineEnd = sampleByExpr.End()
		case p.matchKeyword(KeywordTtl):
			ttlExprList, err := p.tryParseTTLExprList(p.Pos())
			if err != nil {
				return nil, err
			}
			engineExpr.TTLExprList = ttlExprList
			engineEnd = ttlExprList.End()
		case p.matchKeyword(KeywordSettings):
			settingsExprList, err := p.tryParseSettingsExprList(p.Pos())
			if err != nil {
				return nil, err
			}
			engineExpr.SettingsExprList = settingsExprList
			engineEnd = settingsExprList.End()
		default:
			engineExpr.EngineEnd = engineEnd
			return engineExpr, nil
		}
	}
	engineExpr.EngineEnd = engineEnd
	return engineExpr, nil
}

func (p *Parser) parseStatement(pos Pos) (Expr, error) {
	var err error
	var expr Expr
	switch {
	case p.matchKeyword(KeywordCreate),
		p.matchKeyword(KeywordAttach),
		p.matchKeyword(KeywordAlter),
		p.matchKeyword(KeywordDrop),
		p.matchKeyword(KeywordDetach),
		p.matchKeyword(KeywordTruncate),
		p.matchKeyword(KeywordRename):
		expr, err = p.parseDDL(pos)
	case p.matchKeyword(KeywordSelect), p.matchKeyword(KeywordWith):
		expr, err = p.parseSelectQuery(pos)
	case p.matchKeyword(KeywordDelete):
		expr, err = p.parseDeleteFrom(pos)
	case p.matchKeyword(KeywordInsert):
		expr, err = p.parseInsertExpr(p.Pos())
	case p.matchKeyword(KeywordUse):
		expr, err = p.parseUseStatement(pos)
	case p.matchKeyword(KeywordSet):
		expr, err = p.parseSetExpr(pos)
	case p.matchKeyword(KeywordSystem):
		expr, err = p.parseSystemExpr(pos)
	case p.matchKeyword(KeywordOptimize):
		expr, err = p.parseOptimizeExpr(pos)
	case p.matchKeyword(KeywordCheck):
		expr, err = p.parseCheckExpr(pos)
	case p.matchKeyword(KeywordExplain):
		expr, err = p.parseExplainExpr(pos)
	case p.matchKeyword(KeywordGrant):
		expr, err = p.parseGrantPrivilege(pos)
	default:
		return nil, fmt.Errorf("unexpected token: %q", p.last().String)
	}
	if err != nil {
		return nil, err
	}
	_, err = p.tryParseFormatExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	// Statement can be terminated by ';' or EOF
	if p.last() != nil && !p.matchTokenKind(";") {
		return nil, fmt.Errorf("<EOF> or ';' was expected, but got: %q", p.last().String)
	}
	return expr, nil
}

func (p *Parser) ParseStatements() ([]Expr, error) {
	var statements []Expr
	for {
		_ = p.lexer.consumeToken()
		if p.lexer.isEOF() {
			break
		}
		if p.matchTokenKind(";") {
			continue
		}
		statement, err := p.parseStatement(p.Pos())
		if err != nil {
			return nil, p.wrapError(err)
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

func (p *Parser) parseUseStatement(pos Pos) (*UseExpr, error) {
	if err := p.consumeKeyword(KeywordUse); err != nil {
		return nil, err
	}

	database, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	return &UseExpr{
		UsePos:       pos,
		Database:     database,
		StatementEnd: database.End(),
	}, nil
}

// syntax: TRUNCATE TEMPORARY? TABLE (IF EXISTS)? tableIdentifier clusterClause?;
func (p *Parser) parseTruncateTable(pos Pos) (*TruncateTable, error) {
	if err := p.consumeKeyword(KeywordTruncate); err != nil {
		return nil, err
	}

	isTemporary := p.tryConsumeKeyword(KeywordTemporary) != nil

	if err := p.consumeKeyword(KeywordTable); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	tableName, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}

	truncateTable := &TruncateTable{
		TruncatePos:  pos,
		IsTemporary:  isTemporary,
		IfExists:     ifExists,
		Name:         tableName,
		OnCluster:    onCluster,
		StatementEnd: tableName.End(),
	}

	if onCluster != nil {
		truncateTable.StatementEnd = onCluster.End()
	}

	return truncateTable, nil
}

func (p *Parser) parseDeleteFrom(pos Pos) (*DeleteFromExpr, error) {
	if err := p.consumeKeyword(KeywordDelete); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordFrom); err != nil {
		return nil, err
	}
	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}

	if err := p.consumeKeyword(KeywordWhere); err != nil {
		return nil, err
	}
	whereExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	return &DeleteFromExpr{
		DeletePos: pos,
		Table:     tableIdentifier,
		OnCluster: onCluster,
		WhereExpr: whereExpr,
	}, nil
}

func (p *Parser) parseColumnNamesExpr(pos Pos) (*ColumnNamesExpr, error) {
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}

	var columnNames []NestedIdentifier
	for !p.lexer.isEOF() && p.tryConsumeTokenKind(")") == nil {
		name, err := p.ParseNestedIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
		columnNames = append(columnNames, *name)
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}
	return &ColumnNamesExpr{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		ColumnNames:   columnNames,
	}, nil
}

func (p *Parser) parseValuesExpr(pos Pos) (*ValuesExpr, error) {
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}

	var value Expr
	var err error
	values := make([]Expr, 0)
	for !p.lexer.isEOF() && p.tryConsumeTokenKind(")") == nil {
		switch {
		case p.matchTokenKind("("):
			value, err = p.parseValuesExpr(p.Pos())
		default:
			value, err = p.parseExpr(p.Pos())
		}
		if err != nil {
			return nil, err
		}
		values = append(values, value)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	rightParenPos := p.Pos()
	if _, err := p.consumeTokenKind(")"); err != nil {
		return nil, err
	}

	return &ValuesExpr{
		LeftParenPos:  pos,
		RightParenPos: rightParenPos,
		Values:        values,
	}, nil
}

func (p *Parser) parseInsertExpr(pos Pos) (*InsertExpr, error) {
	if err := p.consumeKeyword(KeywordInsert); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordInto); err != nil {
		return nil, err
	}
	_ = p.tryConsumeKeyword(KeywordTable)

	var table Expr
	var err error
	if p.tryConsumeKeyword(KeywordFunction) != nil {
		table, err = p.parseFunctionExpr(p.Pos())
	} else {
		table, err = p.parseTableIdentifier(p.Pos())
	}
	if err != nil {
		return nil, err
	}

	insertExpr := &InsertExpr{
		InsertPos: pos,
		Table:     table,
	}

	for i := 0; i < 1; i++ {
		switch {
		case p.matchKeyword(KeywordFormat):
			insertExpr.Format, err = p.parseFormatExpr(p.Pos())
		case p.matchKeyword(KeywordValues):
			// consume VALUES keyword
			_ = p.lexer.consumeToken()
		case p.matchKeyword(KeywordSelect):
			insertExpr.SelectExpr, err = p.parseSelectQuery(p.Pos())
			if err != nil {
				return nil, err
			}
			return insertExpr, nil
		default:
			if insertExpr.ColumnNames == nil {
				// columns
				insertExpr.ColumnNames, err = p.parseColumnNamesExpr(p.Pos())
				// need another pass to parse keywords
				i--
			}
		}
	}

	if err != nil {
		return nil, err
	}

	values := make([]*ValuesExpr, 0)
	for !p.lexer.isEOF() {
		value, err := p.parseValuesExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		values = append(values, value)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	insertExpr.Values = values

	return insertExpr, nil
}

func (p *Parser) parseRenameStmt(pos Pos) (*RenameStmt, error) {
	if err := p.consumeKeyword(KeywordRename); err != nil {
		return nil, err
	}

	renameTarget := KeywordTable
	switch {
	case p.tryConsumeKeyword(KeywordDictionary) != nil:
		renameTarget = KeywordDictionary
	case p.tryConsumeKeyword(KeywordDatabase) != nil:
		renameTarget = KeywordDatabase
	default:
		if err := p.consumeKeyword(KeywordTable); err != nil {
			return nil, err
		}
	}

	targetPair, err := p.parseTargetPair(p.Pos())
	if err != nil {
		return nil, err
	}
	tablePairList := []*TargetPair{targetPair}
	for p.tryConsumeTokenKind(",") != nil {
		tablePair, err := p.parseTargetPair(p.Pos())
		if err != nil {
			return nil, err
		}
		tablePairList = append(tablePairList, tablePair)
	}

	renameStmt := &RenameStmt{
		RenamePos:    pos,
		StatementEnd: tablePairList[len(tablePairList)-1].End(),

		RenameTarget:   renameTarget,
		TargetPairList: tablePairList,
	}

	onClusterExpr, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	if onClusterExpr != nil {
		renameStmt.OnCluster = onClusterExpr
		renameStmt.StatementEnd = onClusterExpr.End()
	}

	return renameStmt, nil
}

func (p *Parser) parseTargetPair(_ Pos) (*TargetPair, error) {
	oldTable, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	if err = p.consumeKeyword(KeywordTo); err != nil {
		return nil, err
	}
	newTable, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	return &TargetPair{
		Old: oldTable,
		New: newTable,
	}, nil
}

func (p *Parser) parseCreateFunction(pos Pos) (*CreateFunction, error) {
	if err := p.consumeKeyword(KeywordFunction); err != nil {
		return nil, err
	}
	ifNotExists, err := p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}
	functionName, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordAs); err != nil {
		return nil, err
	}
	params, err := p.parseFunctionParams(p.Pos())
	if err != nil {
		return nil, err
	}
	if _, err := p.consumeTokenKind(TokenArrow); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &CreateFunction{
		CreatePos:    pos,
		IfNotExists:  ifNotExists,
		FunctionName: functionName,
		OnCluster:    onCluster,
		Params:       params,
		Expr:         expr,
	}, nil
}
