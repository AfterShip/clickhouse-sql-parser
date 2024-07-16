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
	if err := p.consumeKeyword(KeywordTable); err != nil {
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
		default:
			return nil, errors.New("expected token: ADD|DROP|ATTACH|DETACH|FREEZE|REMOVE|CLEAR")
		}
		if err != nil {
			return nil, err
		}
		alterTable.AlterExprs = append(alterTable.AlterExprs, alter)
		if p.tryConsumeTokenKind(",") == nil {
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
	if err := p.consumeKeyword(KeywordAdd); err != nil {
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
	if err := p.consumeKeyword(KeywordColumn); err != nil {
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
	if err := p.consumeKeyword(KeywordIndex); err != nil {
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
	if err := p.consumeKeyword(KeywordOrder); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordBy); err != nil {
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
	if _, err := p.consumeTokenKind("("); err != nil {
		return nil, err
	}
	with, err := p.tryParseWithClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordSelect); err != nil {
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
	rightParen, err := p.consumeTokenKind(")")
	if err != nil {
		return nil, err
	}
	return &ProjectionSelectStmt{
		LeftParenPos:  pos,
		RightParenPos: rightParen.Pos,
		With:          with,
		SelectColumns: columns,
		GroupBy:       groupBy,
		OrderBy:       orderBy,
	}, nil
}

func (p *Parser) parseTableProjection(pos Pos) (*TableProjection, error) {
	identifier, err := p.ParseNestedIdentifier(pos)
	if err != nil {
		return nil, err
	}
	selectExpr, err := p.parseProjectionSelect(p.Pos())
	if err != nil {
		return nil, err
	}
	return &TableProjection{
		ProjectionPos: pos,
		Identifier:    identifier,
		Select:        selectExpr,
	}, nil
}

func (p *Parser) parseAlterTableAddProjection(pos Pos) (*AlterTableAddProjection, error) {
	if err := p.consumeKeyword(KeywordProjection); err != nil {
		return nil, err
	}

	ifNotExists, err := p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}
	tableProjection, err := p.parseTableProjection(p.Pos())
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

	if err := p.consumeKeyword(KeywordType); err != nil {
		return nil, err
	}
	columnType, err := p.parseColumnType(p.Pos())
	if err != nil {
		return nil, err
	}

	if err := p.consumeKeyword(KeywordGranularity); err != nil {
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
	if err := p.consumeKeyword(KeywordDrop); err != nil {
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
	if err := p.consumeKeyword(KeywordPartition); err != nil {
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
	if err := p.consumeKeyword(KeywordPartition); err != nil {
		return nil, err
	}

	partition := &PartitionClause{
		PartitionPos: pos,
	}
	if p.tryConsumeKeyword(KeywordId) != nil {
		id, err := p.parseString(p.Pos())
		if err != nil {
			return nil, err
		}
		partition.ID = id
	} else if p.tryConsumeKeyword(KeywordAll) != nil {
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

	if err := p.consumeKeyword(KeywordAttach); err != nil {
		return nil, err
	}
	partition, err := p.parsePartitionClause(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTable.Partition = partition
	// FROM [db.]table?
	if p.tryConsumeKeyword(KeywordFrom) != nil {
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
	if p.tryConsumeKeyword(KeywordAfter) == nil {
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
	if err := p.consumeKeyword(KeywordPartition); err != nil {
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
	if err := p.consumeKeyword(KeywordFreeze); err != nil {
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
	if err := p.consumeKeyword(KeywordRemove); err != nil {
		return nil, err
	}

	if err := p.consumeKeyword(KeywordTtl); err != nil {
		return nil, err
	}

	return &AlterTableRemoveTTL{
		RemovePos:    pos,
		StatementEnd: p.Pos(),
	}, nil
}

func (p *Parser) parseAlterTableClear(pos Pos) (AlterTableClause, error) {
	if err := p.consumeKeyword(KeywordClear); err != nil {
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
	if p.tryConsumeKeyword(KeywordIn) != nil {
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
	if err := p.consumeKeyword(KeywordRename); err != nil {
		return nil, err
	}

	if err := p.consumeKeyword(KeywordColumn); err != nil {
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

	if err = p.consumeKeyword(KeywordTo); err != nil {
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
	if err := p.consumeKeyword(KeywordModify); err != nil {
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
	default:
		return nil, fmt.Errorf("expected keyword: COLUMN, but got %q",
			p.last().String)
	}

}

// syntax: MODIFY COLUMN (IF EXISTS)? tableColumnDfnt
func (p *Parser) parseAlterTableModifyColumn(pos Pos) (AlterTableClause, error) {
	if err := p.consumeKeyword(KeywordColumn); err != nil {
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

	if err := p.consumeKeyword(KeywordRemove); err != nil {
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
	if err := p.consumeKeyword(KeywordReplace); err != nil {
		return nil, err
	}

	partition, err := p.parsePartitionClause(p.Pos())
	if err != nil {
		return nil, err
	}

	if err = p.consumeKeyword(KeywordFrom); err != nil {
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
	if err := p.consumeKeyword(KeywordMaterialize); err != nil {
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
	if p.tryConsumeKeyword(KeywordIn) != nil {
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
