package parser

import (
	"errors"
	"fmt"
)

func (p *Parser) parseAlterTable(pos Pos) (*AlterTable, error) {
	alterTable := &AlterTable{
		AlterPos:   pos,
		AlterExprs: make([]AlterTableExpr, 0),
	}
	if err := p.consumeKeyword(KeywordTable); err != nil {
		return nil, err
	}

	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTable.TableIdentifier = tableIdentifier
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTable.OnCluster = onCluster

	for !p.lexer.isEOF() {
		var alterExpr AlterTableExpr
		switch {
		case p.matchKeyword(KeywordAdd):
			alterExpr, err = p.parseAlterTableAdd(p.Pos())
		case p.matchKeyword(KeywordDrop):
			alterExpr, err = p.parseAlterTableDrop(p.Pos())
		case p.matchKeyword(KeywordAttach):
			alterExpr, err = p.parseAlterTableAttachPartition(p.Pos())
		case p.matchKeyword(KeywordDetach):
			_ = p.lexer.consumeToken()
			alterExpr, err = p.parseAlterTableDetachPartition(p.Pos())
		case p.matchKeyword(KeywordFreeze):
			alterExpr, err = p.parseAlterTableFreezePartition(p.Pos())
		case p.matchKeyword(KeywordRemove):
			alterExpr, err = p.parseAlterTableRemoveTTL(p.Pos())
		case p.matchKeyword(KeywordRename):
			alterExpr, err = p.parseAlterTableRenameColumn(p.Pos())
		case p.matchKeyword(KeywordClear):
			alterExpr, err = p.parseAlterTableClear(p.Pos())
		case p.matchKeyword(KeywordModify):
			alterExpr, err = p.parseAlterTableModify(p.Pos())
		case p.matchKeyword(KeywordReplace):
			alterExpr, err = p.parseAlterTableReplacePartition(p.Pos())

		default:
			return nil, errors.New("expected token: ADD|DROP|ATTACH|DETACH|FREEZE|REMOVE|CLEAR")
		}
		if err != nil {
			return nil, err
		}
		alterTable.AlterExprs = append(alterTable.AlterExprs, alterExpr)
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

func (p *Parser) parseAlterTableAdd(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordAdd); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordColumn):
		return p.parseAlterTableAddColumn(pos)
	case p.matchKeyword(KeywordIndex):
		return p.parseAlterTableAddIndex(pos)
	default:
		return nil, errors.New("expected token: COLUMN|INDEX")
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

	column, err := p.parseTableColumn(p.Pos())
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

func (p *Parser) parseAlterTableDrop(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordDrop); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordColumn):
		return p.parseAlterTableDropColumn(pos)
	case p.matchKeyword(KeywordIndex):
		return p.parseAlterTableDropIndex(pos)
	case p.matchKeyword(KeywordDetached):
		_ = p.lexer.consumeToken()
		return p.parseAlterTableDetachPartition(pos)
	case p.matchKeyword(KeywordPartition):
		return p.parseAlterTableDropPartition(pos)
	default:
		return nil, errors.New("expected keyword: COLUMN|INDEX|DETACH")
	}
}

// Syntax: ALTER TABLE DETACH partitionClause
func (p *Parser) parseAlterTableDetachPartition(pos Pos) (AlterTableExpr, error) {
	partitionPos := p.Pos()
	if err := p.consumeKeyword(KeywordPartition); err != nil {
		return nil, err
	}
	partitionExpr := &PartitionExpr{
		PartitionPos: partitionPos,
	}
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	partitionExpr.Expr = expr

	settings, err := p.tryParseSettingsExprList(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableDetachPartition{
		DetachPos: pos,
		Partition: partitionExpr,
		Settings:  settings,
	}, nil
}

func (p *Parser) tryParsePartitionExpr(pos Pos) (*PartitionExpr, error) {
	if !p.matchKeyword(KeywordPartition) {
		return nil, nil // nolint
	}
	return p.parsePartitionExpr(pos)
}

func (p *Parser) parsePartitionExpr(pos Pos) (*PartitionExpr, error) {
	if err := p.consumeKeyword(KeywordPartition); err != nil {
		return nil, err
	}

	partitionExpr := &PartitionExpr{
		PartitionPos: pos,
	}
	if p.tryConsumeKeyword(KeywordId) != nil {
		id, err := p.parseString(p.Pos())
		if err != nil {
			return nil, err
		}
		partitionExpr.ID = id
	} else if p.tryConsumeKeyword(KeywordAll) != nil {
		partitionExpr.All = true
	} else {
		expr, err := p.parseExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		partitionExpr.Expr = expr
	}
	return partitionExpr, nil
}

// Syntax: ALTER TABLE ATTACH partitionClause (FROM tableIdentifier)?
func (p *Parser) parseAlterTableAttachPartition(pos Pos) (AlterTableExpr, error) {
	alterTable := &AlterTableAttachPartition{AttachPos: pos}

	if err := p.consumeKeyword(KeywordAttach); err != nil {
		return nil, err
	}
	partitionExpr, err := p.parsePartitionExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	alterTable.Partition = partitionExpr
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

func (p *Parser) parseAlterTableDropColumn(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordColumn); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	name, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableDropColumn{
		DropPos:    pos,
		ColumnName: name,
		IfExists:   ifExists,
	}, nil
}

func (p *Parser) parseAlterTableDropIndex(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordIndex); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	name, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	return &AlterTableDropIndex{
		DropPos:   pos,
		IndexName: name,
		IfExists:  ifExists,
	}, nil
}

func (p *Parser) tryParseAfterClause() (*NestedIdentifier, error) {
	if p.tryConsumeKeyword(KeywordAfter) == nil {
		return nil, nil // nolint
	}

	return p.ParseNestedIdentifier(p.Pos())
}

// Syntax: ALTER TABLE DROP partitionClause
func (p *Parser) parseAlterTableDropPartition(pos Pos) (AlterTableExpr, error) {
	partitionPos := p.Pos()
	if err := p.consumeKeyword(KeywordPartition); err != nil {
		return nil, err
	}
	partitionExpr := &PartitionExpr{
		PartitionPos: partitionPos,
	}
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	partitionExpr.Expr = expr

	return &AlterTableDropPartition{
		DropPos:   pos,
		Partition: partitionExpr,
	}, nil
}

func (p *Parser) parseAlterTableFreezePartition(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordFreeze); err != nil {
		return nil, err
	}
	alterTable := &AlterTableFreezePartition{
		FreezePos:    pos,
		StatementEnd: p.Pos(),
	}
	if p.matchKeyword(KeywordPartition) {
		partitionExpr, err := p.parsePartitionExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		alterTable.Partition = partitionExpr
		alterTable.StatementEnd = partitionExpr.End()
	}

	return alterTable, nil
}

func (p *Parser) parseAlterTableRemoveTTL(pos Pos) (AlterTableExpr, error) {
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

func (p *Parser) parseAlterTableClear(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordClear); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordColumn):
		return p.parseAlterTableClearColumn(pos)
	case p.matchKeyword(KeywordIndex):
		return p.parseAlterTableClearIndex(pos)

	default:
		return nil, errors.New("expected token: COLUMN|INDEX|PROJECTION")
	}
}

// Syntax: ALTER TABLE CLEAR COLUMN (IF EXISTS)? nestedIdentifier (IN partitionClause)?
func (p *Parser) parseAlterTableClearColumn(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordColumn); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	columnName, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := columnName.End()

	var partitionExpr *PartitionExpr
	if p.tryConsumeKeyword(KeywordIn) != nil {
		partitionExpr, err = p.tryParsePartitionExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		if partitionExpr != nil {
			statementEnd = partitionExpr.End()
		}
	}

	return &AlterTableClearColumn{
		ClearPos:      pos,
		StatementEnd:  statementEnd,
		IfExists:      ifExists,
		ColumnName:    columnName,
		PartitionExpr: partitionExpr,
	}, nil
}

// Syntax: ALTER TABLE CLEAR INDEX (IF EXISTS)? nestedIdentifier (IN partitionClause)?
func (p *Parser) parseAlterTableClearIndex(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordIndex); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	indexName, err := p.ParseNestedIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := indexName.End()

	var partitionExpr *PartitionExpr
	if p.tryConsumeKeyword(KeywordIn) != nil {
		partitionExpr, err = p.tryParsePartitionExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		if partitionExpr != nil {
			statementEnd = partitionExpr.End()
		}
	}

	return &AlterTableClearIndex{
		ClearPos:      pos,
		StatementEnd:  statementEnd,
		IfExists:      ifExists,
		IndexName:     indexName,
		PartitionExpr: partitionExpr,
	}, nil
}

// Syntax: ALTER TABLE RENAME COLUMN (IF EXISTS)? nestedIdentifier TO nestedIdentifier
func (p *Parser) parseAlterTableRenameColumn(pos Pos) (AlterTableExpr, error) {
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

func (p *Parser) parseAlterTableModify(pos Pos) (AlterTableExpr, error) {
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
func (p *Parser) parseAlterTableModifyColumn(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordColumn); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	// at least parse out column name
	column, err := p.parseTableColumn(p.Pos())
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

func (p *Parser) parseAlterTableReplacePartition(pos Pos) (AlterTableExpr, error) {
	if err := p.consumeKeyword(KeywordReplace); err != nil {
		return nil, err
	}

	partitionExpr, err := p.parsePartitionExpr(p.Pos())
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
		Partition:  partitionExpr,
		Table:      table,
	}, nil
}
