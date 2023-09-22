package parser

func (p *Parser) parseDropDatabase(pos Pos) (*DropDatabase, error) {
	if err := p.consumeKeyword(KeywordDatabase); err != nil {
		return nil, err
	}

	isExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	statementEnd := name.End()

	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	if onCluster != nil {
		statementEnd = onCluster.End()
	}

	return &DropDatabase{
		DropPos:      pos,
		Name:         name,
		IfExists:     isExists,
		OnCluster:    onCluster,
		StatementEnd: statementEnd,
	}, nil
}

func (p *Parser) parseDropTable(pos Pos) (*DropTable, error) {
	isTemporary := p.tryConsumeKeyword(KeywordTemporary) != nil

	if err := p.consumeKeyword(KeywordTable); err != nil {
		return nil, err
	}

	isExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := name.End()

	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	if onCluster != nil {
		statementEnd = onCluster.End()
	}

	noDelay, err := p.tryParseNoDelay()
	if err != nil {
		return nil, err
	}
	if noDelay {
		statementEnd = p.Pos()
	}

	return &DropTable{
		DropPos:      pos,
		Name:         name,
		IfExists:     isExists,
		OnCluster:    onCluster,
		IsTemporary:  isTemporary,
		NoDelay:      noDelay,
		StatementEnd: statementEnd,
	}, nil
}

func (p *Parser) tryParseNoDelay() (bool, error) {
	if p.tryConsumeKeyword(KeywordNo) == nil {
		return false, nil
	}

	if err := p.consumeKeyword(KeywordDelay); err != nil {
		return false, err
	}

	return true, nil
}
