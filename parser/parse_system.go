package parser

import (
	"fmt"
	"strings"
)

func (p *Parser) parseSetExpr(pos Pos) (*SetExpr, error) {
	if err := p.consumeKeyword(KeywordSet); err != nil {
		return nil, err
	}
	settings, err := p.parseSettingsExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	return &SetExpr{
		SetPos:   pos,
		Settings: settings,
	}, nil
}

func (p *Parser) parseSystemFlushExpr(pos Pos) (*SystemFlushExpr, error) {
	if err := p.consumeKeyword(KeywordFlush); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordLogs):
		lastToken := p.last()
		_ = p.lexer.consumeToken()
		return &SystemFlushExpr{
			FlushPos:     pos,
			StatementEnd: lastToken.End,
			Logs:         true,
		}, nil
	case p.tryConsumeKeyword(KeywordDistributed) != nil:
		distributed, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		return &SystemFlushExpr{
			FlushPos:     pos,
			StatementEnd: distributed.End(),
			Distributed:  distributed,
		}, nil
	default:
		return nil, fmt.Errorf("expected LOGS|DISTRIBUTED")
	}
}

func (p *Parser) parseSystemReloadExpr(pos Pos) (*SystemReloadExpr, error) {
	if err := p.consumeKeyword(KeywordReload); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordDictionaries):
		lastToken := p.last()
		_ = p.lexer.consumeToken()
		return &SystemReloadExpr{
			ReloadPos:    pos,
			StatementEnd: lastToken.End,
			Type:         KeywordDictionaries,
		}, nil
	case p.tryConsumeKeyword(KeywordDictionary) != nil:
		dictionary, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		return &SystemReloadExpr{
			ReloadPos:    pos,
			StatementEnd: dictionary.End(),
			Type:         KeywordDictionary,
			Dictionary:   dictionary,
		}, nil
	case p.tryConsumeKeyword("EMBEDDED") != nil:
		lastToken := p.last()
		if err := p.consumeKeyword(KeywordDictionaries); err != nil {
			return nil, err
		}
		return &SystemReloadExpr{
			ReloadPos:    pos,
			StatementEnd: lastToken.End,
			Type:         "EMBEDDED DICTIONARIES",
		}, nil
	default:
		return nil, fmt.Errorf("expected DICTIONARIES|CONFIG")
	}
}

func (p *Parser) parseSystemSyncExpr(pos Pos) (*SystemSyncExpr, error) {
	if err := p.consumeKeyword(KeywordSync); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordReplica); err != nil {
		return nil, err
	}
	cluster, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	return &SystemSyncExpr{
		SyncPos: pos,
		Cluster: cluster,
	}, nil
}

func (p *Parser) parseSystemCtrlExpr(pos Pos) (*SystemCtrlExpr, error) {
	if !p.matchKeyword(KeywordStart) && !p.matchKeyword(KeywordStop) {
		return nil, fmt.Errorf("expected START|STOP")
	}
	command := strings.ToUpper(p.last().String)
	_ = p.lexer.consumeToken()

	var typ string
	switch {
	case p.tryConsumeKeyword(KeywordDistributed) != nil:
		switch {
		case p.matchKeyword(KeywordSends):
			typ = "DISTRIBUTED SENDS"
		case p.matchKeyword(KeywordFetches):
			typ = "FETCHES"
		case p.matchKeyword(KeywordMerges):
			typ = "MERGES"
		case p.matchKeyword(KeywordTtl):
			typ = "TTL MERGES"
			if err := p.consumeKeyword(KeywordMerges); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("expected SENDS|FETCHES|MERGES|TTL")
		}
		cluster, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		return &SystemCtrlExpr{
			CtrlPos:      pos,
			StatementEnd: cluster.End(),
			Command:      command,
			Type:         typ,
			Cluster:      cluster,
		}, nil
	case p.tryConsumeKeyword(KeywordReplicated) != nil:
		lastToken := p.last()
		if err := p.consumeKeyword(KeywordSends); err != nil {
			return nil, err
		}
		typ = "REPLICATED SENDS"
		return &SystemCtrlExpr{
			CtrlPos:      pos,
			StatementEnd: lastToken.End,
			Command:      command,
			Type:         typ,
		}, nil
	default:
		return nil, fmt.Errorf("expected DISTRIBUTED|REPLICATED")
	}
}

func (p *Parser) parseSystemDropExpr(pos Pos) (*SystemDropExpr, error) {
	if err := p.consumeKeyword(KeywordDrop); err != nil {
		return nil, err
	}
	switch {
	case p.matchKeyword(KeywordDNS),
		p.matchKeyword(KeywordMark),
		p.matchKeyword(KeywordUncompressed),
		p.matchKeyword(KeywordFileSystem),
		p.matchKeyword(KeywordQuery):
		prefixToken := p.last()
		_ = p.lexer.consumeToken()
		lastToken := p.last()
		if err := p.consumeKeyword(KeywordCache); err != nil {
			return nil, err
		}
		return &SystemDropExpr{
			DropPos:      pos,
			StatementEnd: lastToken.End,
			Type:         prefixToken.String + " CACHE",
		}, nil
	case p.matchKeyword(KeywordCompiled):
		_ = p.lexer.consumeToken()
		if err := p.consumeKeyword(KeywordExpression); err != nil {
			return nil, err
		}
		lastToken := p.last()
		if err := p.consumeKeyword(KeywordCache); err != nil {
			return nil, err
		}
		return &SystemDropExpr{
			DropPos:      pos,
			StatementEnd: lastToken.End,
			Type:         "COMPILED EXPRESSION CACHE",
		}, nil
	default:
		return nil, fmt.Errorf("expected DNS|MARK|REPLICA|DATABASE|UNCOMPRESSION|COMPILED|QUERY")
	}
}

func (p *Parser) tryParseDeduplicateExpr(pos Pos) (*DeduplicateExpr, error) {
	if !p.matchKeyword(KeywordDeduplicate) {
		return nil, nil
	}
	return p.parseDeduplicateExpr(pos)
}

func (p *Parser) parseDeduplicateExpr(pos Pos) (*DeduplicateExpr, error) {
	if err := p.consumeKeyword(KeywordDeduplicate); err != nil {
		return nil, err
	}
	if p.tryConsumeKeyword(KeywordBy) == nil {
		return &DeduplicateExpr{
			DeduplicatePos: pos,
		}, nil
	}

	by, err := p.parseColumnExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	var except *ColumnExprList
	if p.tryConsumeKeyword(KeywordExcept) != nil {
		except, err = p.parseColumnExprList(p.Pos())
		if err != nil {
			return nil, err
		}
	}
	return &DeduplicateExpr{
		DeduplicatePos: pos,
		By:             by,
		Except:         except,
	}, nil
}

func (p *Parser) parseOptimizeExpr(pos Pos) (*OptimizeExpr, error) {
	if err := p.consumeKeyword(KeywordOptimize); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordTable); err != nil {
		return nil, err
	}

	table, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	statmentEnd := table.End()

	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	if onCluster != nil {
		statmentEnd = onCluster.End()
	}

	partitionExpr, err := p.tryParsePartitionExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if partitionExpr != nil {
		statmentEnd = partitionExpr.End()
	}

	hasFinal := false
	lastPos := p.Pos()
	if p.tryConsumeKeyword(KeywordFinal) != nil {
		hasFinal = true
		statmentEnd = lastPos
	}

	deduplicate, err := p.tryParseDeduplicateExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	if deduplicate != nil {
		statmentEnd = deduplicate.End()
	}

	return &OptimizeExpr{
		OptimizePos:  pos,
		StatementEnd: statmentEnd,
		Table:        table,
		OnCluster:    onCluster,
		Partition:    partitionExpr,
		HasFinal:     hasFinal,
		Deduplicate:  deduplicate,
	}, nil
}

func (p *Parser) parseSystemExpr(pos Pos) (*SystemExpr, error) {
	if err := p.consumeKeyword(KeywordSystem); err != nil {
		return nil, err
	}

	var err error
	var expr Expr
	switch {
	case p.matchKeyword(KeywordFlush):
		expr, err = p.parseSystemFlushExpr(p.Pos())
	case p.matchKeyword(KeywordReload):
		expr, err = p.parseSystemReloadExpr(p.Pos())
	case p.matchKeyword(KeywordSync):
		expr, err = p.parseSystemSyncExpr(p.Pos())
	case p.matchKeyword(KeywordStart), p.matchKeyword(KeywordStop):
		expr, err = p.parseSystemCtrlExpr(p.Pos())
	case p.matchKeyword(KeywordDrop):
		expr, err = p.parseSystemDropExpr(p.Pos())
	default:
		return nil, fmt.Errorf("expected FLUSH|RELOAD|SYNC|START|STOP")
	}
	if err != nil {
		return nil, err
	}
	return &SystemExpr{
		SystemPos: pos,
		Expr:      expr,
	}, nil
}

func (p *Parser) parseCheckExpr(pos Pos) (*CheckExpr, error) {
	if err := p.consumeKeyword(KeywordCheck); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordTable); err != nil {
		return nil, err
	}
	table, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	partition, err := p.tryParsePartitionExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	return &CheckExpr{
		CheckPos:  pos,
		Table:     table,
		Partition: partition,
	}, nil
}

func (p *Parser) parseRoleName(pos Pos) (*RoleName, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	return &RoleName{
		Name:      name,
		OnCluster: onCluster,
	}, nil
}

func (p *Parser) tryParseRoleSettings(pos Pos) ([]*RoleSetting, error) {
	if p.tryConsumeKeyword(KeywordSettings) == nil {
		return nil, nil
	}
	return p.parseRoleSettings(pos)
}

func (p *Parser) parseRoleSetting(pos Pos) (*RoleSetting, error) {
	pairs := make([]*SettingPair, 0)
	for p.matchTokenKind(TokenIdent) {
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		switch name.Name {
		case "NONE", "READABLE", "WRITABLE", "CONST", "CHANGEABLE_IN_READONLY":
			return &RoleSetting{
				Modifier:     name,
				SettingPairs: pairs,
			}, nil
		}
		switch {
		case p.matchTokenKind("="),
			p.matchTokenKind(TokenInt),
			p.matchTokenKind(TokenFloat),
			p.matchTokenKind(TokenString):
			_ = p.tryConsumeTokenKind("=")
			value, err := p.parseLiteral(p.Pos())
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, &SettingPair{
				Name:  name,
				Value: value,
			})
		default:
			pairs = append(pairs, &SettingPair{
				Name: name,
			})
		}

	}
	return &RoleSetting{
		SettingPairs: pairs,
	}, nil
}

func (p *Parser) parseRoleSettings(pos Pos) ([]*RoleSetting, error) {
	settings := make([]*RoleSetting, 0)
	for {
		setting, err := p.parseRoleSetting(p.Pos())
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
		if p.tryConsumeTokenKind(",") == nil {
			break
		}
	}
	return settings, nil
}

func (p *Parser) parseCreateRole(pos Pos) (*CreateRole, error) {
	if err := p.consumeKeyword(KeywordRole); err != nil {
		return nil, err
	}

	ifNotExists := false
	orReplace := false
	switch {
	case p.matchKeyword(KeywordIf):
		_ = p.lexer.consumeToken()
		if err := p.consumeKeyword(KeywordNot); err != nil {
			return nil, err
		}
		if err := p.consumeKeyword(KeywordExists); err != nil {
			return nil, err
		}
		ifNotExists = true
	case p.matchKeyword(KeywordOr):
		_ = p.lexer.consumeToken()
		if err := p.consumeKeyword(KeywordReplace); err != nil {
			return nil, err
		}
		orReplace = true
	}

	roleNames := make([]*RoleName, 0)
	roleName, err := p.parseRoleName(p.Pos())
	if err != nil {
		return nil, err
	}
	roleNames = append(roleNames, roleName)
	for p.tryConsumeTokenKind(",") != nil {
		roleName, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		roleNames = append(roleNames, roleName)
	}
	statementEnd := roleNames[len(roleNames)-1].End()

	var accessStorageType *Ident
	if p.tryConsumeKeyword(KeywordIn) != nil {
		accessStorageType, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
		statementEnd = accessStorageType.NameEnd
	}

	settings, err := p.tryParseRoleSettings(p.Pos())
	if err != nil {
		return nil, err
	}
	if settings != nil {
		statementEnd = settings[len(settings)-1].End()
	}

	return &CreateRole{
		CreatePos:         pos,
		StatementEnd:      statementEnd,
		IfNotExists:       ifNotExists,
		OrReplace:         orReplace,
		RoleNames:         roleNames,
		AccessStorageType: accessStorageType,
		Settings:          settings,
	}, nil
}
