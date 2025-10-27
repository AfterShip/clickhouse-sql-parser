package parser

import (
	"errors"
	"fmt"
)

func (p *Parser) parseAlterTable(pos Pos) (*AlterTable, error) {
	alterTable := &AlterTable{
		AlterPos:   pos,
		AlterExprs: make([]AlterTableClause, 0),
	}
	if err := p.expectKeyword(KeywordTable); err != nil {
		return nil, err
	}

	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTable.TableIdentifier = tableIdentifier
	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTable.OnCluster = onCluster

	for !p.lexer.isEOF() {
		var alter AlterTableClause
		switch {
		case p.matchKeyword(KeywordAdd):
			alter, err = p.parseAlterTableAdd(p.Pos())
		case p.matchKeyword(KeywordDrop):
			alter, err = p.parseAlterTableDrop(p.Pos())
		case p.matchKeyword(KeywordAttach):
			alter, err = p.parseAlterTableAttachPartition(p.Pos())
		case p.matchKeyword(KeywordDetach):
			_ = p.lexer.consumeToken()
			alter, err = p.parseAlterTableDetachPartition(p.Pos())
		case p.matchKeyword(KeywordFreeze):
			alter, err = p.parseAlterTableFreezePartition(p.Pos())
		case p.matchKeyword(KeywordRemove):
			alter, err = p.parseAlterTableRemoveTTL(p.Pos())
		case p.matchKeyword(KeywordRename):
			alter, err = p.parseAlterTableRenameColumn(p.Pos())
		case p.matchKeyword(KeywordClear):
			alter, err = p.parseAlterTableClear(p.Pos())
		case p.matchKeyword(KeywordModify):
			alter, err = p.parseAlterTableModify(p.Pos())
		case p.matchKeyword(KeywordReplace):
			alter, err = p.parseAlterTableReplacePartition(p.Pos())
		case p.matchKeyword(KeywordMaterialize):
			alter, err = p.parseAlterTableMaterialize(p.Pos())
		case p.matchKeyword(KeywordReset):
			alter, err = p.parseAlterTableReset(p.Pos())
		case p.matchKeyword(KeywordDelete):
			alter, err = p.parseAlterTableDelete(p.Pos())
		case p.matchKeyword(KeywordUpdate):
			alter, err = p.parseAlterTableUpdate(p.Pos())
		default:
			return nil, errors.New("expected token: ADD|DROP|ATTACH|DETACH|FREEZE|REMOVE|CLEAR|MODIFY|REPLACE|MATERIALIZE|RESET|DELETE|UPDATE")
		}
		if err != nil {
			return nil, err
		}
		alterTable.AlterExprs = append(alterTable.AlterExprs, alter)
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
			break
		}
	}
	if len(alterTable.AlterExprs) == 0 {
		return nil, errors.New("expected token: ADD|DROP")
	}
	alterTable.StatementEnd = alterTable.AlterExprs[len(alterTable.AlterExprs)-1].End()

	return alterTable, nil
}

func (p *Parser) parseAlterTableAdd(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordAdd); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordColumn):
		return p.parseAlterTableAddColumn(pos)
	case p.matchKeyword(KeywordIndex):
		return p.parseAlterTableAddIndex(pos)
	case p.matchKeyword(KeywordProjection):
		return p.parseAlterTableAddProjection(pos)
	default:
		return nil, errors.New("expected token: COLUMN|INDEX|PROJECTION")
	}
}

func (p *Parser) parseAlterTableAddColumn(pos Pos) (*AlterTableAddColumn, error) {
	if err := p.expectKeyword(KeywordColumn); err != nil {
		return nil, err
	}

	ifNotExists, err := p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}

	column, err := p.parseTableColumnExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := column.End()

	after, err := p.tryParseAfterClause()
	if err != nil {
		return nil, err
	}
	if after != nil {
		statementEnd = after.End()
	}

	return &AlterTableAddColumn{
		AddPos:       pos,
		StatementEnd: statementEnd,
		Column:       column,
		IfNotExists:  ifNotExists,
		After:        after,
	}, nil
}

func (p *Parser) parseAlterTableAddIndex(pos Pos) (*AlterTableAddIndex, error) {
	indexPos := p.Pos()
	if err := p.expectKeyword(KeywordIndex); err != nil {
		return nil, err
	}

	ifNotExists, err := p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}
	index, err := p.parseTableIndex(indexPos)
	if err != nil {
		return nil, err
	}
	statementEnd := index.End()
	after, err := p.tryParseAfterClause()
	if err != nil {
		return nil, err
	}
	if after != nil {
		statementEnd = after.End()
	}
	return &AlterTableAddIndex{
		AddPos:       pos,
		StatementEnd: statementEnd,
		IfNotExists:  ifNotExists,
		Index:        index,
		After:        after,
	}, nil
}

func (p *Parser) parseProjectionOrderBy(pos Pos) (*ProjectionOrderByClause, error) {
	if err := p.expectKeyword(KeywordOrder); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordBy); err != nil {
		return nil, err
	}
	columns, err := p.parseColumnExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	return &ProjectionOrderByClause{
		OrderByPos: pos,
		Columns:    columns,
	}, nil
}

func (p *Parser) parseProjectionSelect(pos Pos) (*ProjectionSelectStmt, error) {
	if err := p.expectTokenKind(TokenKindLParen); err != nil {
		return nil, err
	}
	with, err := p.tryParseWithClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordSelect); err != nil {
		return nil, err
	}
	columns, err := p.parseColumnExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	groupBy, err := p.tryParseGroupByClause(p.Pos())
	if err != nil {
		return nil, err
	}
	orderBy, err := p.parseProjectionOrderBy(p.Pos())
	if err != nil {
		return nil, err
	}

	lastToken := p.last()
	if err := p.expectTokenKind(TokenKindRParen); err != nil {
		return nil, err
	}
	return &ProjectionSelectStmt{
		LeftParenPos:  pos,
		RightParenPos: lastToken.Pos,
		With:          with,
		SelectColumns: columns,
		GroupBy:       groupBy,
		OrderBy:       orderBy,
	}, nil
}

func (p *Parser) parseTableProjection(pos Pos, includeProjectionKeyword bool) (*TableProjection, error) {
	if includeProjectionKeyword {
		if err := p.expectKeyword(KeywordProjection); err != nil {
			return nil, err
		}
	}
	identifier, err := p.ParseNestedIdentifier(pos)
	if err != nil {
		return nil, err
	}
	selectExpr, err := p.parseProjectionSelect(p.Pos())
	if err != nil {
		return nil, err
	}
	return &TableProjection{
		IncludeProjectionKeyword: includeProjectionKeyword,
		ProjectionPos:            pos,
		Identifier:               identifier,
		Select:                   selectExpr,
	}, nil
}

func (p *Parser) parseAlterTableAddProjection(pos Pos) (*AlterTableAddProjection, error) {
	if err := p.expectKeyword(KeywordProjection); err != nil {
		return nil, err
	}

	ifNotExists, err := p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}
	tableProjection, err := p.parseTableProjection(p.Pos(), false)
	if err != nil {
		return nil, err
	}
	statementEnd := tableProjection.End()
	after, err := p.tryParseAfterClause()
	if err != nil {
		return nil, err
	}
	if after != nil {
		statementEnd = after.End()
	}
	return &AlterTableAddProjection{
		AddPos:          pos,
		StatementEnd:    statementEnd,
		IfNotExists:     ifNotExists,
		TableProjection: tableProjection,
		After:           after,
	}, nil
}

func (p *Parser) parseTableIndex(pos Pos) (*TableIndex, error) {
	name, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	columnExpr, err := p.parseColumnsExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordType); err != nil {
		return nil, err
	}
	columnType, err := p.parseColumnType(p.Pos())
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordGranularity); err != nil {
		return nil, err
	}
	granularity, err := p.parseDecimal(p.Pos())
	if err != nil {
		return nil, err
	}

	return &TableIndex{
		IndexPos:    pos,
		Name:        name,
		ColumnExpr:  columnExpr,
		ColumnType:  columnType,
		Granularity: granularity,
	}, nil
}

func (p *Parser) parseAlterTableDrop(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordDrop); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordColumn), p.matchKeyword(KeywordIndex), p.matchKeyword(KeywordProjection):
		return p.parseAlterTableDropClause(pos)
	case p.matchKeyword(KeywordDetached), p.matchKeyword(KeywordPartition):
		return p.parseAlterTableDropPartition(pos)
	default:
		return nil, errors.New("expected keyword: COLUMN|INDEX|PROJECTION|DETACHED|PARTITION")
	}
}

// Syntax: ALTER TABLE DETACH partitionClause
func (p *Parser) parseAlterTableDetachPartition(pos Pos) (AlterTableClause, error) {
	partitionPos := p.Pos()
	if err := p.expectKeyword(KeywordPartition); err != nil {
		return nil, err
	}
	partition := &PartitionClause{
		PartitionPos: partitionPos,
	}
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	partition.Expr = expr

	settings, err := p.tryParseSettingsClause(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableDetachPartition{
		DetachPos: pos,
		Partition: partition,
		Settings:  settings,
	}, nil
}

func (p *Parser) tryParsePartitionClause(pos Pos) (*PartitionClause, error) {
	if !p.matchKeyword(KeywordPartition) {
		return nil, nil // nolint
	}
	return p.parsePartitionClause(pos)
}

func (p *Parser) parsePartitionClause(pos Pos) (*PartitionClause, error) {
	if err := p.expectKeyword(KeywordPartition); err != nil {
		return nil, err
	}

	partition := &PartitionClause{
		PartitionPos: pos,
	}
	if p.tryConsumeKeywords(KeywordId) {
		id, err := p.parseString(p.Pos())
		if err != nil {
			return nil, err
		}
		partition.ID = id
	} else if p.tryConsumeKeywords(KeywordAll) {
		partition.All = true
	} else {
		expr, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		partition.Expr = expr
	}
	return partition, nil
}

// Syntax: ALTER TABLE ATTACH partitionClause (FROM tableIdentifier)?
func (p *Parser) parseAlterTableAttachPartition(pos Pos) (AlterTableClause, error) {
	alterTable := &AlterTableAttachPartition{AttachPos: pos}

	if err := p.expectKeyword(KeywordAttach); err != nil {
		return nil, err
	}
	partition, err := p.parsePartitionClause(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTable.Partition = partition
	// FROM [db.]table?
	if p.tryConsumeKeywords(KeywordFrom) {
		tableIdentifier, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		alterTable.From = tableIdentifier
	}
	return alterTable, nil
}

func (p *Parser) parseAlterTableDropClause(pos Pos) (AlterTableClause, error) {
	var kind string
	switch {
	case p.matchKeyword(KeywordColumn):
		kind = KeywordColumn
	case p.matchKeyword(KeywordIndex):
		kind = KeywordIndex
	case p.matchKeyword(KeywordProjection):
		kind = KeywordProjection
	default:
		return nil, fmt.Errorf("expected token: COLUMN|INDEX|PROJECTION, but got %s", p.lastTokenKind())
	}
	_ = p.lexer.consumeToken()

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	name, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	if kind == KeywordProjection {
		return &AlterTableDropProjection{
			DropPos:        pos,
			ProjectionName: name,
			IfExists:       ifExists,
		}, nil
	} else if kind == KeywordColumn {
		return &AlterTableDropColumn{
			DropPos:    pos,
			ColumnName: name,
			IfExists:   ifExists,
		}, nil
	} else {
		return &AlterTableDropIndex{
			DropPos:   pos,
			IndexName: name,
			IfExists:  ifExists,
		}, nil
	}
}

func (p *Parser) tryParseAfterClause() (*NestedIdentifier, error) {
	if !p.tryConsumeKeywords(KeywordAfter) {
		return nil, nil // nolint
	}

	return p.ParseNestedIdentifier(p.Pos())
}

// Syntax: ALTER TABLE DROP partitionClause
func (p *Parser) parseAlterTableDropPartition(pos Pos) (AlterTableClause, error) {
	var hasDetached bool
	if p.matchKeyword(KeywordDetached) {
		_ = p.lexer.consumeToken()
		hasDetached = true
	}
	partitionPos := p.Pos()
	if err := p.expectKeyword(KeywordPartition); err != nil {
		return nil, err
	}
	partition := &PartitionClause{
		PartitionPos: partitionPos,
	}
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	partition.Expr = expr

	settings, err := p.tryParseSettingsClause(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableDropPartition{
		DropPos:     pos,
		Partition:   partition,
		HasDetached: hasDetached,
		Settings:    settings,
	}, nil
}

func (p *Parser) parseAlterTableFreezePartition(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordFreeze); err != nil {
		return nil, err
	}
	alterTable := &AlterTableFreezePartition{
		FreezePos:    pos,
		StatementEnd: p.Pos(),
	}
	if p.matchKeyword(KeywordPartition) {
		partition, err := p.parsePartitionClause(p.Pos())
		if err != nil {
			return nil, err
		}
		alterTable.Partition = partition
		alterTable.StatementEnd = partition.End()
	}

	return alterTable, nil
}

func (p *Parser) parseAlterTableRemoveTTL(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordRemove); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordTtl); err != nil {
		return nil, err
	}

	return &AlterTableRemoveTTL{
		RemovePos:    pos,
		StatementEnd: p.Pos(),
	}, nil
}

func (p *Parser) parseAlterTableClear(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordClear); err != nil {
		return nil, err
	}
	return p.parseAlterTableClearClause(pos)
}

// Syntax: ALTER TABLE CLEAR COLUMN|INDEX|PROJECTION (IF EXISTS)? nestedIdentifier (IN partitionClause)?
func (p *Parser) parseAlterTableClearClause(pos Pos) (AlterTableClause, error) {
	var kind string
	switch {
	case p.matchKeyword(KeywordColumn):
		kind = KeywordColumn
	case p.matchKeyword(KeywordIndex):
		kind = KeywordIndex
	case p.matchKeyword(KeywordProjection):
		kind = KeywordProjection
	default:
		return nil, fmt.Errorf("expected keyword: COLUMN|INDEX|PROJECTION, but got %q", p.lastTokenKind())
	}
	_ = p.lexer.consumeToken()

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	name, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := name.End()

	var partition *PartitionClause
	if p.tryConsumeKeywords(KeywordIn) {
		partition, err = p.tryParsePartitionClause(p.Pos())
		if err != nil {
			return nil, err
		}
		if partition != nil {
			statementEnd = partition.End()
		}
	}

	if kind == KeywordProjection {
		return &AlterTableClearProjection{
			ClearPos:       pos,
			StatementEnd:   statementEnd,
			IfExists:       ifExists,
			ProjectionName: name,
			PartitionExpr:  partition,
		}, nil
	} else if kind == KeywordColumn {
		return &AlterTableClearColumn{
			ClearPos:      pos,
			StatementEnd:  statementEnd,
			IfExists:      ifExists,
			ColumnName:    name,
			PartitionExpr: partition,
		}, nil
	} else {
		return &AlterTableClearIndex{
			ClearPos:      pos,
			StatementEnd:  statementEnd,
			IfExists:      ifExists,
			IndexName:     name,
			PartitionExpr: partition,
		}, nil
	}
}

// Syntax: ALTER TABLE RENAME COLUMN (IF EXISTS)? nestedIdentifier TO nestedIdentifier
func (p *Parser) parseAlterTableRenameColumn(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordRename); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordColumn); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	oldColumnName, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	if err = p.expectKeyword(KeywordTo); err != nil {
		return nil, err
	}

	newColumnName, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableRenameColumn{
		RenamePos:     pos,
		IfExists:      ifExists,
		OldColumnName: oldColumnName,
		NewColumnName: newColumnName,
	}, nil
}

func (p *Parser) parseAlterTableModify(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordModify); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordColumn):
		return p.parseAlterTableModifyColumn(pos)
	case p.matchKeyword(KeywordTtl):
		_ = p.lexer.consumeToken()
		ttlExpr, err := p.parseTTLExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		return &AlterTableModifyTTL{
			ModifyPos:    pos,
			StatementEnd: ttlExpr.End(),
			TTL:          ttlExpr,
		}, nil
	case p.matchKeyword(KeywordQuery):
		_ = p.lexer.consumeToken()
		selectQuery, _ := p.parseSelectQuery(pos)
		return &AlterTableModifyQuery{
			ModifyPos:    pos,
			StatementEnd: selectQuery.End(),
			SelectExpr:   selectQuery,
		}, nil
	case p.matchKeyword(KeywordSetting):
		_ = p.lexer.consumeToken() // consume "SETTING"
		settings, err := p.parseSettingsList(p.Pos())
		if err != nil {
			return nil, err
		}
		// settings must not be empty
		statementEnd := settings[len(settings)-1].End()
		return &AlterTableModifySetting{
			ModifyPos:    pos,
			StatementEnd: statementEnd,
			Settings:     settings,
		}, nil
	default:
		return nil, fmt.Errorf("expected keyword: COLUMN|TTL|QUERY|SETTING, but got %q",
			p.last().String)
	}

}

// syntax: MODIFY COLUMN (IF EXISTS)? tableColumnDfnt
func (p *Parser) parseAlterTableModifyColumn(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordColumn); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	// at least parse out column name
	column, err := p.parseTableColumnExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	alterTableModifyColumn := &AlterTableModifyColumn{
		ModifyPos:    pos,
		StatementEnd: column.End(),
		IfExists:     ifExists,
		Column:       column,
	}

	// syntax: MODIFY COLUMN (IF EXISTS)? nestedIdentifier REMOVE tableColumnPropertyType
	removePropertyType, err := p.tryParseRemovePropertyTypeExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTableModifyColumn.RemovePropertyType = removePropertyType

	return alterTableModifyColumn, nil
}

func (p *Parser) tryParseRemovePropertyTypeExpr(pos Pos) (*RemovePropertyType, error) {
	if !p.matchKeyword(KeywordRemove) {
		return nil, nil
	}

	if err := p.expectKeyword(KeywordRemove); err != nil {
		return nil, err
	}

	columnPropertyType, err := p.parseColumnPropertyType(p.Pos())
	if err != nil {
		return nil, err
	}

	return &RemovePropertyType{
		RemovePos:    pos,
		PropertyType: columnPropertyType,
	}, nil
}

func (p *Parser) parseAlterTableReplacePartition(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordReplace); err != nil {
		return nil, err
	}

	partition, err := p.parsePartitionClause(p.Pos())
	if err != nil {
		return nil, err
	}

	if err = p.expectKeyword(KeywordFrom); err != nil {
		return nil, err
	}

	table, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableReplacePartition{
		ReplacePos: pos,
		Partition:  partition,
		Table:      table,
	}, nil
}

func (p *Parser) parseAlterTableMaterialize(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordMaterialize); err != nil {
		return nil, err
	}
	var kind string
	switch {
	case p.matchKeyword(KeywordIndex):
		kind = KeywordIndex
	case p.matchKeyword(KeywordProjection):
		kind = KeywordProjection
	default:
		return nil, fmt.Errorf("expected keyword: INDEX|PROJECTION, but got %q", p.lastTokenKind())
	}
	_ = p.lexer.consumeToken()

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}
	name, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := name.End()
	var partition *PartitionClause
	if p.tryConsumeKeywords(KeywordIn) {
		partition, err = p.tryParsePartitionClause(p.Pos())
		if err != nil {
			return nil, err
		}
		statementEnd = partition.End()
	}
	if kind == KeywordIndex {
		return &AlterTableMaterializeIndex{
			MaterializedPos: pos,
			StatementEnd:    statementEnd,
			IfExists:        ifExists,
			IndexName:       name,
			Partition:       partition,
		}, nil
	}
	return &AlterTableMaterializeProjection{
		MaterializedPos: pos,
		StatementEnd:    statementEnd,
		IfExists:        ifExists,
		ProjectionName:  name,
		Partition:       partition,
	}, nil
}

func (p *Parser) parseAlterTableReset(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordReset); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordSetting); err != nil {
		return nil, err
	}

	// Parse comma-separated setting names inline
	var settings []*Ident
	setting, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	settings = append(settings, setting)

	for p.tryConsumeTokenKind(TokenKindComma) != nil {
		setting, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	statementEnd := settings[len(settings)-1].End()

	return &AlterTableResetSetting{
		ResetPos:     pos,
		StatementEnd: statementEnd,
		Settings:     settings,
	}, nil
}

// Syntax: ALTER TABLE DELETE WHERE condition
func (p *Parser) parseAlterTableDelete(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordDelete); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordWhere); err != nil {
		return nil, err
	}

	whereExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableDelete{
		DeletePos:    pos,
		StatementEnd: whereExpr.End(),
		WhereClause:  whereExpr,
	}, nil
}

// Syntax: ALTER TABLE UPDATE column1 = expr1 [, column2 = expr2, ...] WHERE condition
func (p *Parser) parseAlterTableUpdate(pos Pos) (AlterTableClause, error) {
	if err := p.expectKeyword(KeywordUpdate); err != nil {
		return nil, err
	}

	// Parse at least one assignment
	assignments := make([]*UpdateAssignment, 0)
	assignment, err := p.parseUpdateAssignment(p.Pos())
	if err != nil {
		return nil, err
	}
	assignments = append(assignments, assignment)

	// Parse additional comma-separated assignments
	for p.tryConsumeTokenKind(TokenKindComma) != nil {
		assignment, err = p.parseUpdateAssignment(p.Pos())
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	if err := p.expectKeyword(KeywordWhere); err != nil {
		return nil, err
	}

	whereExpr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableUpdate{
		UpdatePos:    pos,
		StatementEnd: whereExpr.End(),
		Assignments:  assignments,
		WhereClause:  whereExpr,
	}, nil
}

// Parse column = expression assignment
func (p *Parser) parseUpdateAssignment(pos Pos) (*UpdateAssignment, error) {
	column, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	if err := p.expectTokenKind(TokenKindSingleEQ); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}

	return &UpdateAssignment{
		AssignmentPos: pos,
		Column:        column,
		Expr:          expr,
	}, nil
}
