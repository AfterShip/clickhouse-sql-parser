package parser

func (p *Parser) parseDropDatabase(pos Pos) (*DropDatabase, error) {
	if err := p.expectKeyword(KeywordDatabase); err != nil {
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

	onCluster, err := p.tryParseClusterClause(p.Pos())
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

func (p *Parser) parseDropStmt(pos Pos) (*DropStmt, error) {
	var isTemporary bool
	dropTarget := KeywordTable
	switch {
	case p.tryConsumeKeywords(KeywordDictionary):
		dropTarget = KeywordDictionary
	case p.tryConsumeKeywords(KeywordView):
		dropTarget = KeywordView
	default:
		isTemporary = p.tryConsumeKeywords(KeywordTemporary)
		if err := p.expectKeyword(KeywordTable); err != nil {
			return nil, err
		}
	}

	isExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}

	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}

	modifier, err := p.tryParseModifier()
	if err != nil {
		return nil, err
	}

	return &DropStmt{
		DropPos:      pos,
		DropTarget:   dropTarget,
		Name:         name,
		IfExists:     isExists,
		OnCluster:    onCluster,
		IsTemporary:  isTemporary,
		Modifier:     modifier,
		StatementEnd: p.Pos(),
	}, nil
}

func (p *Parser) tryParseModifier() (string, error) {
	switch {
	case p.tryConsumeKeywords(KeywordSync):
		return "SYNC", nil
	case p.tryConsumeKeywords(KeywordNo):
		if err := p.expectKeyword(KeywordDelay); err != nil {
			return "", err
		}
		return "NO DELAY", nil
	}
	return "", nil
}
