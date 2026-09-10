package parser

import (
	"fmt"
	"strings"
)

func (p *Parser) parseSetStmt(pos Pos) (*SetStmt, error) {
	if err := p.expectKeyword(KeywordSet); err != nil {
		return nil, err
	}
	settings, err := p.parseSettingsClause(p.Pos())
	if err != nil {
		return nil, err
	}
	return &SetStmt{
		SetPos:   pos,
		Settings: settings,
	}, nil
}

func (p *Parser) parseSettingsStmt(pos Pos) (*SetStmt, error) {
	if err := p.expectKeyword(KeywordSettings); err != nil {
		return nil, err
	}
	settings, err := p.parseSettingsClause(p.Pos())
	if err != nil {
		return nil, err
	}
	return &SetStmt{
		SetPos:   pos,
		Settings: settings,
	}, nil
}

func (p *Parser) parseSystemFlushExpr(pos Pos) (*SystemFlushExpr, error) {
	if err := p.expectKeyword(KeywordFlush); err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordLogs):
		curToken := p.current()
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		return &SystemFlushExpr{
			FlushPos:     pos,
			StatementEnd: curToken.End,
			Logs:         true,
		}, nil
	case p.matchKeyword(KeywordDistributed):
		if err := p.expectKeyword(KeywordDistributed); err != nil {
			return nil, err
		}
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
	if err := p.expectKeyword(KeywordReload); err != nil {
		return nil, err
	}

	var typ string
	var statementEnd Pos
	// Only RELOAD DICTIONARY takes a dictionary name.
	hasDictionaryName := false
	switch {
	case p.matchKeyword(KeywordDictionaries):
		typ = KeywordDictionaries
		statementEnd = p.current().End
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
	case p.matchKeyword(KeywordDictionary):
		typ = KeywordDictionary
		statementEnd = p.current().End
		hasDictionaryName = true
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
	case p.matchKeyword(KeywordEmbedded):
		if err := p.expectKeyword(KeywordEmbedded); err != nil {
			return nil, err
		}
		typ = "EMBEDDED DICTIONARIES"
		statementEnd = p.current().End
		if err := p.expectKeyword(KeywordDictionaries); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("expected DICTIONARIES|DICTIONARY|EMBEDDED")
	}

	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if onCluster != nil {
		statementEnd = onCluster.End()
	}

	var dictionary *TableIdentifier
	if hasDictionaryName {
		dictionary, err = p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		statementEnd = dictionary.End()

		// ClickHouse also accepts ON CLUSTER after the dictionary name.
		if onCluster == nil {
			onCluster, err = p.tryParseClusterClause(p.Pos())
			if err != nil {
				return nil, err
			}
			if onCluster != nil {
				statementEnd = onCluster.End()
			}
		}
	}

	return &SystemReloadExpr{
		ReloadPos:    pos,
		StatementEnd: statementEnd,
		OnCluster:    onCluster,
		Type:         typ,
		Dictionary:   dictionary,
	}, nil
}

func (p *Parser) parseSystemSyncExpr(pos Pos) (*SystemSyncExpr, error) {
	if err := p.expectKeyword(KeywordSync); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordReplica); err != nil {
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
	command := strings.ToUpper(p.current().String)
	if err := p.lexer.consumeToken(); err != nil {
		return nil, err
	}

	var typ string
	switch {
	case p.matchKeyword(KeywordDistributed):
		if err := p.expectKeyword(KeywordDistributed); err != nil {
			return nil, err
		}
		switch {
		case p.matchKeyword(KeywordSends):
			typ = "DISTRIBUTED SENDS"
		case p.matchKeyword(KeywordFetches):
			typ = "FETCHES"
		case p.matchKeyword(KeywordMerges):
			typ = "MERGES"
		case p.matchKeyword(KeywordTtl):
			typ = "TTL MERGES"
			if err := p.expectKeyword(KeywordMerges); err != nil {
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
	case p.matchKeyword(KeywordReplicated):
		if err := p.expectKeyword(KeywordReplicated); err != nil {
			return nil, err
		}
		curToken := p.current()
		if err := p.expectKeyword(KeywordSends); err != nil {
			return nil, err
		}
		typ = "REPLICATED SENDS"
		return &SystemCtrlExpr{
			CtrlPos:      pos,
			StatementEnd: curToken.End,
			Command:      command,
			Type:         typ,
		}, nil
	default:
		return nil, fmt.Errorf("expected DISTRIBUTED|REPLICATED")
	}
}

func (p *Parser) parseSystemDropExpr(pos Pos) (*SystemDropExpr, error) {
	if err := p.expectKeyword(KeywordDrop); err != nil {
		return nil, err
	}
	switch {
	case p.matchKeyword(KeywordDNS),
		p.matchKeyword(KeywordMark),
		p.matchKeyword(KeywordUncompressed),
		p.matchKeyword(KeywordFileSystem),
		p.matchKeyword(KeywordQuery):
		prefixToken := p.current()
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		curToken := p.current()
		if err := p.expectKeyword(KeywordCache); err != nil {
			return nil, err
		}
		return &SystemDropExpr{
			DropPos:      pos,
			StatementEnd: curToken.End,
			Type:         prefixToken.String + " CACHE",
		}, nil
	case p.matchKeyword(KeywordCompiled):
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordExpression); err != nil {
			return nil, err
		}
		curToken := p.current()
		if err := p.expectKeyword(KeywordCache); err != nil {
			return nil, err
		}
		return &SystemDropExpr{
			DropPos:      pos,
			StatementEnd: curToken.End,
			Type:         "COMPILED EXPRESSION CACHE",
		}, nil
	default:
		return nil, fmt.Errorf("expected DNS|MARK|REPLICA|DATABASE|UNCOMPRESSION|COMPILED|QUERY")
	}
}

func (p *Parser) tryParseDeduplicateClause(pos Pos) (*DeduplicateClause, error) {
	if !p.matchKeyword(KeywordDeduplicate) {
		return nil, nil
	}
	return p.parseDeduplicateClause(pos)
}

func (p *Parser) parseDeduplicateClause(pos Pos) (*DeduplicateClause, error) {
	if err := p.expectKeyword(KeywordDeduplicate); err != nil {
		return nil, err
	}
	if matched, consumeErr := p.tryConsumeKeywords(KeywordBy); consumeErr != nil {
		return nil, consumeErr
	} else if !matched {
		return &DeduplicateClause{
			DeduplicatePos: pos,
		}, nil
	}

	by, err := p.parseColumnExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	var except *ColumnExprList
	if matched, consumeErr := p.tryConsumeKeywords(KeywordExcept); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		except, err = p.parseColumnExprList(p.Pos())
		if err != nil {
			return nil, err
		}
	}
	return &DeduplicateClause{
		DeduplicatePos: pos,
		By:             by,
		Except:         except,
	}, nil
}

func (p *Parser) parseOptimizeStmt(pos Pos) (*OptimizeStmt, error) {
	if err := p.expectKeyword(KeywordOptimize); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordTable); err != nil {
		return nil, err
	}

	table, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	statementEnd := table.End()

	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if onCluster != nil {
		statementEnd = onCluster.End()
	}

	partition, err := p.tryParsePartitionClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if partition != nil {
		statementEnd = partition.End()
	}

	hasFinal := false
	lastPos := p.Pos()
	if matched, consumeErr := p.tryConsumeKeywords(KeywordFinal); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		hasFinal = true
		statementEnd = lastPos
	}

	deduplicate, err := p.tryParseDeduplicateClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if deduplicate != nil {
		statementEnd = deduplicate.End()
	}

	return &OptimizeStmt{
		OptimizePos:  pos,
		StatementEnd: statementEnd,
		Table:        table,
		OnCluster:    onCluster,
		Partition:    partition,
		HasFinal:     hasFinal,
		Deduplicate:  deduplicate,
	}, nil
}

func (p *Parser) parseSystemStmt(pos Pos) (*SystemStmt, error) {
	if err := p.expectKeyword(KeywordSystem); err != nil {
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
	return &SystemStmt{
		SystemPos: pos,
		Expr:      expr,
	}, nil
}

func (p *Parser) parseCheckStmt(pos Pos) (*CheckStmt, error) {
	if err := p.expectKeyword(KeywordCheck); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordTable); err != nil {
		return nil, err
	}
	table, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	partition, err := p.tryParsePartitionClause(p.Pos())
	if err != nil {
		return nil, err
	}
	return &CheckStmt{
		CheckPos:  pos,
		Table:     table,
		Partition: partition,
	}, nil
}

func (p *Parser) parseRoleName(_ Pos) (*RoleName, error) {
	switch {
	case p.matchTokenKind(TokenKindIdent):
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		var scope *StringLiteral
		if token, consumeErr := p.tryConsumeTokenKind(TokenKindAtSign); consumeErr != nil {
			return nil, consumeErr
		} else if token != nil {
			scope, err = p.parseString(p.Pos())
			if err != nil {
				return nil, err
			}
		}
		onCluster, err := p.tryParseClusterClause(p.Pos())
		if err != nil {
			return nil, err
		}
		return &RoleName{
			Name:      name,
			Scope:     scope,
			OnCluster: onCluster,
		}, nil
	case p.matchTokenKind(TokenKindString):
		name, err := p.parseString(p.Pos())
		if err != nil {
			return nil, err
		}
		onCluster, err := p.tryParseClusterClause(p.Pos())
		if err != nil {
			return nil, err
		}
		return &RoleName{
			Name:      name,
			OnCluster: onCluster,
		}, nil
	default:
		return nil, fmt.Errorf("expected <ident> or <string>")
	}
}

func (p *Parser) tryParseRoleSettings(pos Pos) ([]*RoleSetting, error) {
	if matched, consumeErr := p.tryConsumeKeywords(KeywordSettings); consumeErr != nil {
		return nil, consumeErr
	} else if !matched {
		return nil, nil
	}
	return p.parseRoleSettings(pos)
}

func (p *Parser) parseRoleSetting(_ Pos) (*RoleSetting, error) {
	pairs := make([]*SettingPair, 0)
	for p.matchTokenKind(TokenKindIdent) {
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
		case p.matchTokenKind(TokenKindSingleEQ),
			p.matchTokenKind(TokenKindInt),
			p.matchTokenKind(TokenKindFloat),
			p.matchTokenKind(TokenKindString):
			var op TokenKind
			if token, consumeErr := p.tryConsumeTokenKind(TokenKindSingleEQ); consumeErr != nil {
				return nil, consumeErr
			} else if token != nil {
				op = token.Kind
			}
			value, err := p.parseLiteral(p.Pos())
			if err != nil {
				return nil, err
			}
			// docs: https://clickhouse.com/docs/en/sql-reference/statements/alter/role
			// the operator "=" was required if the variable name is NOT in
			// ["MIN", "MAX", "PROFILE"] and value is existed.
			if value != nil && name.Name != "MIN" && name.Name != "MAX" && name.Name != "PROFILE" && op != TokenKindSingleEQ {
				return nil, fmt.Errorf("expected operator = or no value, but got %s", op)
			}
			pairs = append(pairs, &SettingPair{
				Name:      name,
				Operation: op,
				Value:     value,
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

func (p *Parser) parseRoleSettings(_ Pos) ([]*RoleSetting, error) {
	settings := make([]*RoleSetting, 0)
	for {
		setting, err := p.parseRoleSetting(p.Pos())
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
		if token, consumeErr := p.tryConsumeTokenKind(TokenKindComma); consumeErr != nil {
			return nil, consumeErr
		} else if token == nil {
			break
		}
	}
	return settings, nil
}

func (p *Parser) parseCreateRole(pos Pos) (*CreateRole, error) {
	if err := p.expectKeyword(KeywordRole); err != nil {
		return nil, err
	}

	ifNotExists := false
	orReplace := false
	switch {
	case p.matchKeyword(KeywordIf):
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordNot); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordExists); err != nil {
			return nil, err
		}
		ifNotExists = true
	case p.matchKeyword(KeywordOr):
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordReplace); err != nil {
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
	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		roleName, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		roleNames = append(roleNames, roleName)
	}
	statementEnd := roleNames[len(roleNames)-1].End()

	var accessStorageType *Ident
	if matched, consumeErr := p.tryConsumeKeywords(KeywordIn); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
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

func (p *Parser) parseAuthenticationClause(pos Pos) (*AuthenticationClause, error) {
	auth := &AuthenticationClause{AuthPos: pos}

	if matched, consumeErr := p.tryConsumeKeywords(KeywordNot); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		if err := p.expectKeyword(KeywordIdentified); err != nil {
			return nil, err
		}
		auth.NotIdentified = true
		auth.AuthEnd = p.current().End
		return auth, nil
	}

	if err := p.expectKeyword(KeywordIdentified); err != nil {
		return nil, err
	}
	auth.AuthEnd = p.current().End

	if matched, consumeErr := p.tryConsumeKeywords(KeywordWith); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		if p.matchKeyword(KeywordLdap) {
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			if err := p.expectKeyword(KeywordServer); err != nil {
				return nil, err
			}
			server, err := p.parseString(p.Pos())
			if err != nil {
				return nil, err
			}
			auth.LdapServer = server
			auth.AuthEnd = server.End()
		} else if p.matchKeyword(KeywordKerberos) {
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			auth.IsKerberos = true
			auth.AuthEnd = p.current().End
			if matched, consumeErr := p.tryConsumeKeywords(KeywordRealm); consumeErr != nil {
				return nil, consumeErr
			} else if matched {
				realm, err := p.parseString(p.Pos())
				if err != nil {
					return nil, err
				}
				auth.KerberosRealm = realm
				auth.AuthEnd = realm.End()
			}
		} else if p.matchTokenKind(TokenKindIdent) {
			// Auth types like no_password, plaintext_password, etc.
			authType := p.current().String
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			auth.AuthType = authType
			auth.AuthEnd = p.current().End

			if matched, consumeErr := p.tryConsumeKeywords(KeywordBy); consumeErr != nil {
				return nil, consumeErr
			} else if matched {
				value, err := p.parseString(p.Pos())
				if err != nil {
					return nil, err
				}
				auth.AuthValue = value
				auth.AuthEnd = value.End()
			}
		}
	}

	return auth, nil
}

func (p *Parser) parseHostClause(pos Pos) (*HostClause, error) {
	if err := p.expectKeyword(KeywordHost); err != nil {
		return nil, err
	}

	host := &HostClause{HostPos: pos}

	switch {
	case p.matchOneOfKeywords(KeywordLocal, KeywordAny, KeywordNone):
		hostType := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		host.HostType = hostType
		host.HostEnd = p.current().End
	case p.matchOneOfKeywords(KeywordName, KeywordRegexp, KeywordIp, KeywordLike):
		hostType := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		host.HostType = hostType
		value, err := p.parseString(p.Pos())
		if err != nil {
			return nil, err
		}
		host.HostValue = value
		host.HostEnd = value.End()
	default:
		return nil, fmt.Errorf("expected LOCAL|NAME|REGEXP|IP|LIKE|ANY|NONE")
	}

	return host, nil
}

func (p *Parser) parseDefaultRoleClause(pos Pos) (*DefaultRoleClause, error) {
	if err := p.expectKeyword(KeywordDefault); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordRole); err != nil {
		return nil, err
	}

	defaultRole := &DefaultRoleClause{DefaultPos: pos}

	if matched, consumeErr := p.tryConsumeKeywords(KeywordNone); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		defaultRole.None = true
		defaultRole.DefaultEnd = p.current().End
		return defaultRole, nil
	}

	roles := make([]*RoleName, 0)
	role, err := p.parseRoleName(p.Pos())
	if err != nil {
		return nil, err
	}
	roles = append(roles, role)

	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		role, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	defaultRole.Roles = roles
	defaultRole.DefaultEnd = roles[len(roles)-1].End()
	return defaultRole, nil
}

func (p *Parser) parseGranteesClause(pos Pos) (*GranteesClause, error) {
	if err := p.expectKeyword(KeywordGrantees); err != nil {
		return nil, err
	}

	grantees := &GranteesClause{GranteesPos: pos}

	if matched, consumeErr := p.tryConsumeKeywords(KeywordAny); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		grantees.Any = true
		grantees.GranteesEnd = p.current().End
	} else if matched, consumeErr := p.tryConsumeKeywords(KeywordNone); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		grantees.None = true
		grantees.GranteesEnd = p.current().End
	} else {
		// Parse list of grantees
		granteeList := make([]*RoleName, 0)
		grantee, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		granteeList = append(granteeList, grantee)

		for {
			token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
			if consumeErr != nil {
				return nil, consumeErr
			}
			if token == nil {
				break
			}
			grantee, err := p.parseRoleName(p.Pos())
			if err != nil {
				return nil, err
			}
			granteeList = append(granteeList, grantee)
		}

		grantees.Grantees = granteeList
		grantees.GranteesEnd = granteeList[len(granteeList)-1].End()
	}

	// Check for EXCEPT clause
	if matched, consumeErr := p.tryConsumeKeywords(KeywordExcept); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		exceptList := make([]*RoleName, 0)
		except, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		exceptList = append(exceptList, except)

		for {
			token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
			if consumeErr != nil {
				return nil, consumeErr
			}
			if token == nil {
				break
			}
			except, err := p.parseRoleName(p.Pos())
			if err != nil {
				return nil, err
			}
			exceptList = append(exceptList, except)
		}

		grantees.ExceptUsers = exceptList
		grantees.GranteesEnd = exceptList[len(exceptList)-1].End()
	}

	return grantees, nil
}

func (p *Parser) parseCreateUserModifiers(createUser *CreateUser) error {
	switch {
	case p.matchKeyword(KeywordIf):
		if err := p.lexer.consumeToken(); err != nil {
			return err
		}
		if err := p.expectKeyword(KeywordNot); err != nil {
			return err
		}
		if err := p.expectKeyword(KeywordExists); err != nil {
			return err
		}
		createUser.IfNotExists = true
	case p.matchKeyword(KeywordOr):
		if err := p.lexer.consumeToken(); err != nil {
			return err
		}
		if err := p.expectKeyword(KeywordReplace); err != nil {
			return err
		}
		createUser.OrReplace = true
	}
	return nil
}

func (p *Parser) parseUserNames() ([]*RoleName, error) {
	userNames := make([]*RoleName, 0)
	userName, err := p.parseRoleName(p.Pos())
	if err != nil {
		return nil, err
	}
	userNames = append(userNames, userName)

	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		userName, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		userNames = append(userNames, userName)
	}
	return userNames, nil
}

func (p *Parser) parseHostClauses() ([]*HostClause, error) {
	hosts := make([]*HostClause, 0)
	host, err := p.parseHostClause(p.Pos())
	if err != nil {
		return nil, err
	}
	hosts = append(hosts, host)

	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		host, err := p.parseHostClause(p.Pos())
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func (p *Parser) parseDefaultClause(createUser *CreateUser) (bool, error) {
	nextToken, err := p.lexer.peekToken()
	if err != nil {
		return false, err
	}

	if nextToken.String == KeywordRole {
		defaultRole, err := p.parseDefaultRoleClause(p.Pos())
		if err != nil {
			return false, err
		}
		createUser.DefaultRole = defaultRole
		createUser.StatementEnd = defaultRole.End()
		return true, nil
	} else if nextToken.String == KeywordDatabase {
		if err := p.lexer.consumeToken(); err != nil { // consume DEFAULT
			return false, err
		}
		if err := p.lexer.consumeToken(); err != nil { // consume DATABASE
			return false, err
		}
		if matched, consumeErr := p.tryConsumeKeywords(KeywordNone); consumeErr != nil {
			return false, consumeErr
		} else if matched {
			createUser.DefaultDbNone = true
			createUser.StatementEnd = p.current().End
		} else {
			db, err := p.parseIdent()
			if err != nil {
				return false, err
			}
			createUser.DefaultDatabase = db
			createUser.StatementEnd = db.End()
		}
		return true, nil
	}
	return false, nil
}

func (p *Parser) parseOptionalClauses(createUser *CreateUser) error {
	continueParsing := true
	for continueParsing {
		switch {
		case p.matchOneOfKeywords(KeywordNot, KeywordIdentified):
			auth, err := p.parseAuthenticationClause(p.Pos())
			if err != nil {
				return err
			}
			createUser.Authentication = auth
			createUser.StatementEnd = auth.End()

		case p.matchKeyword(KeywordValid):
			if err := p.lexer.consumeToken(); err != nil { // consume VALID keyword
				return err
			}
			if err := p.expectKeyword(KeywordUntil); err != nil {
				return err
			}
			validUntil, err := p.parseString(p.Pos())
			if err != nil {
				return err
			}
			createUser.ValidUntil = validUntil
			createUser.StatementEnd = validUntil.End()

		case p.matchKeyword(KeywordHost):
			hosts, err := p.parseHostClauses()
			if err != nil {
				return err
			}
			createUser.Hosts = hosts
			createUser.StatementEnd = hosts[len(hosts)-1].End()

		case p.matchKeyword(KeywordDefault):
			parsed, err := p.parseDefaultClause(createUser)
			if err != nil {
				return err
			}
			if !parsed {
				continueParsing = false
			}

		case p.matchKeyword(KeywordGrantees):
			grantees, err := p.parseGranteesClause(p.Pos())
			if err != nil {
				return err
			}
			createUser.Grantees = grantees
			createUser.StatementEnd = grantees.End()

		case p.matchKeyword(KeywordSettings):
			if err := p.lexer.consumeToken(); err != nil { // consume SETTINGS keyword
				return err
			}
			settings, err := p.parseRoleSettings(p.Pos())
			if err != nil {
				return err
			}
			createUser.Settings = settings
			if len(settings) > 0 {
				createUser.StatementEnd = settings[len(settings)-1].End()
			}

		default:
			continueParsing = false
		}
	}
	return nil
}

func (p *Parser) parseCreateUser(pos Pos) (*CreateUser, error) {
	if err := p.expectKeyword(KeywordUser); err != nil {
		return nil, err
	}

	createUser := &CreateUser{CreatePos: pos}

	// Handle IF NOT EXISTS or OR REPLACE
	if err := p.parseCreateUserModifiers(createUser); err != nil {
		return nil, err
	}

	// Parse user names
	userNames, err := p.parseUserNames()
	if err != nil {
		return nil, err
	}
	createUser.UserNames = userNames
	createUser.StatementEnd = userNames[len(userNames)-1].End()

	// Parse optional clauses
	if err := p.parseOptionalClauses(createUser); err != nil {
		return nil, err
	}

	return createUser, nil
}

func (p *Parser) parserDropUserOrRole(pos Pos) (*DropUserOrRole, error) {
	var target string
	switch {
	case p.matchOneOfKeywords(KeywordUser, KeywordRole):
		target = p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("expected USER|ROLE")
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	names := make([]*RoleName, 0)
	name, err := p.parseRoleName(p.Pos())
	if err != nil {
		return nil, err
	}
	names = append(names, name)
	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		name, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	statementEnd := names[len(names)-1].End()

	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	if onCluster != nil {
		statementEnd = onCluster.End()
	}

	var from *Ident
	if matched, consumeErr := p.tryConsumeKeywords(KeywordFrom); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		from, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
	}

	modifier, err := p.tryParseModifier()
	if err != nil {
		return nil, err
	}

	return &DropUserOrRole{
		DropPos:      pos,
		StatementEnd: statementEnd,
		Target:       target,
		IfExists:     ifExists,
		Names:        names,
		From:         from,
		Modifier:     modifier,
	}, nil
}

func (p *Parser) parsePrivilegeSelectOrInsert(pos Pos) (*PrivilegeClause, error) {
	keyword := p.current().String
	if err := p.lexer.consumeToken(); err != nil {
		return nil, err
	}

	var err error
	var params *ParamExprList
	if p.matchTokenKind(TokenKindLParen) {
		params, err = p.parseFunctionParams(p.Pos())
		if err != nil {
			return nil, err
		}
	}
	return &PrivilegeClause{
		PrivilegePos: pos,
		Keywords:     []string{keyword},
		Params:       params,
	}, nil
}

func (p *Parser) parsePrivilegeAlter(pos Pos) (*PrivilegeClause, error) {
	keywords := []string{KeywordAlter}
	switch {
	case p.matchKeyword(KeywordIndex):
		if err := p.expectKeyword(KeywordIndex); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordIndex)
	case p.matchOneOfKeywords(KeywordUpdate, KeywordDelete, KeywordUser, KeywordRole, KeywordQuota):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	case p.matchOneOfKeywords(KeywordAdd, KeywordDrop, KeywordModify, KeywordClear, KeywordComment, KeywordRename, KeywordMaterialized):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
		switch {
		case p.matchKeyword(KeywordColumn):
			if err := p.expectKeyword(KeywordColumn); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordColumn)
		case p.matchKeyword(KeywordIndex):
			if err := p.expectKeyword(KeywordIndex); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordIndex)
			keywords = append(keywords, KeywordConstraint)
		case p.matchKeyword(KeywordTtl):
			if err := p.expectKeyword(KeywordTtl); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordTtl)
		default:
			return nil, fmt.Errorf("expected COLUMN|INDEX")
		}
	case p.matchKeyword(KeywordOrder):
		if err := p.expectKeyword(KeywordOrder); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordBy); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordOrder, KeywordBy)
	case p.matchKeyword(KeywordSample):
		if err := p.expectKeyword(KeywordSample); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordBy); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordSample, KeywordBy)
	case p.matchKeyword(KeywordSettings):
		if err := p.expectKeyword(KeywordSettings); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordSettings)
	case p.matchKeyword(KeywordView):
		if err := p.expectKeyword(KeywordView); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordView)
		switch {
		case p.matchKeyword(KeywordModify):
			if err := p.expectKeyword(KeywordModify); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordModify)
		case p.matchKeyword(KeywordRefresh):
			if err := p.expectKeyword(KeywordRefresh); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordRefresh)
		default:
			return nil, fmt.Errorf("expected MODIFY|REFRESH")
		}
	case p.matchOneOfKeywords(KeywordMove, KeywordFreeze):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
		if err := p.expectKeyword(KeywordPartition); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordPartition)
	default:
		return nil, fmt.Errorf("expected UPDATE|DELETE|ADD|DROP|MODIFY|CLEAR|COMMENT|RENAME|MATERIALIZED|ORDER|SAMPLE|SETTINGS|VIEW|MOVE|FREEZE")
	}
	return &PrivilegeClause{
		PrivilegePos: pos,
		Keywords:     keywords,
	}, nil
}

func (p *Parser) parsePrivilegeCreate(pos Pos) (*PrivilegeClause, error) {
	keywords := []string{KeywordCreate}
	switch {
	case p.matchOneOfKeywords(KeywordDatabase, KeywordDictionary, KeywordTable, KeywordFunction, KeywordView, KeywordUser, KeywordRole, KeywordQuota):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	case p.matchKeyword(KeywordTemporary):
		if err := p.expectKeyword(KeywordTemporary); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordTable); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordTemporary, KeywordTable)
	case p.matchKeyword(KeywordRows):
		if err := p.expectKeyword(KeywordRows); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordPolicy); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordRows, KeywordPolicy)
	default:
		return nil, fmt.Errorf("expected DATABASE|DICTIONARY|TABLE|FUNCTION|VIEW|USER|ROLE|ROWS")
	}
	return &PrivilegeClause{
		PrivilegePos: pos,
		Keywords:     keywords,
	}, nil
}

func (p *Parser) parsePrivilegeDrop(pos Pos) (*PrivilegeClause, error) {
	keywords := []string{KeywordDrop}
	switch {
	case p.matchOneOfKeywords(KeywordDatabase, KeywordDictionary, KeywordUser, KeywordRole, KeywordQuota, KeywordTable, KeywordFunction, KeywordView):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	default:
		return nil, fmt.Errorf("expected DATABASE|DICTIONARY|TABLE|FUNCTION|VIEW")
	}
	return &PrivilegeClause{
		PrivilegePos: pos,
		Keywords:     keywords,
	}, nil
}

func (p *Parser) parsePrivilegeShow(pos Pos) (*PrivilegeClause, error) {
	keywords := []string{KeywordShow}
	switch {
	case p.matchOneOfKeywords(KeywordDatabases, KeywordDictionaries, KeywordTables, KeywordColumns):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	default:
		return nil, fmt.Errorf("expected DATABASES|DICTIONARIES|TABLES|COLUMNS")
	}
	return &PrivilegeClause{
		PrivilegePos: pos,
		Keywords:     keywords,
	}, nil
}

func (p *Parser) parsePrivilegeSystem(pos Pos) (*PrivilegeClause, error) {
	keywords := []string{KeywordShow}
	switch {
	case p.matchOneOfKeywords(KeywordShutdown, KeywordMerges, KeywordFetches, KeywordSends, KeywordMoves, KeywordCluster):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	case p.matchKeyword(KeywordDrop):
		if err := p.expectKeyword(KeywordDrop); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordDrop)
		switch {
		case p.matchKeyword(KeywordCache):
			if err := p.expectKeyword(KeywordCache); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordCache)
		case p.matchOneOfKeywords(KeywordMark, KeywordDNS, KeywordUncompressed):
			keyword := p.current().String
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			keywords = append(keywords, keyword)
			if err := p.expectKeyword(KeywordCache); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordCache)
		default:
			return nil, fmt.Errorf("expected CACHE|MARK|DNS|UNCOMPRESSED")
		}
	case p.matchKeyword(KeywordReload):
		if err := p.expectKeyword(KeywordReload); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordReload)
		switch {
		case p.matchOneOfKeywords(KeywordDictionary, KeywordFunction, KeywordFunctions, KeywordConfig):
			keyword := p.current().String
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			keywords = append(keywords, keyword)
		default:
			return nil, fmt.Errorf("expected DICTIONARY|FUNCTION|FUNCTIONS|CONFIG")
		}
	case p.matchKeyword(KeywordFlush):
		if err := p.expectKeyword(KeywordFlush); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordFlush)
		switch {
		case p.matchOneOfKeywords(KeywordLogs, KeywordDistributed):
			keyword := p.current().String
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			keywords = append(keywords, keyword)
		default:
			return nil, fmt.Errorf("expected LOGS|DISTRIBUTED")
		}
	case p.matchKeyword(KeywordTtl):
		if err := p.expectKeyword(KeywordTtl); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordTtl)
		if err := p.expectKeyword(KeywordMerges); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordMerges)
	case p.matchOneOfKeywords(KeywordSync, KeywordRestart):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
		if err := p.expectKeyword(KeywordReplica); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordReplica)
	case p.matchKeyword(KeywordReplication):
		if err := p.expectKeyword(KeywordReplication); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordReplication)
		if err := p.expectKeyword(KeywordQueues); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordQueues)
	default:
		return nil, fmt.Errorf("expected QUEUES|SHUTDOWN|MERGES|FETCHES|SENDS|MOVES|CLUSTER|DROP|RELOAD|FLUSH|TTL|SYNC|RESTART|REPLICATION")
	}
	return &PrivilegeClause{
		PrivilegePos: pos,
		Keywords:     keywords,
	}, nil
}

func (p *Parser) parsePrivilegeClause(pos Pos) (*PrivilegeClause, error) {
	if p.matchTokenKind(TokenKindIdent) {
		if p.current().String == "dictGet" {
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			return &PrivilegeClause{
				PrivilegePos: pos,
				Keywords:     []string{"dictGet"},
			}, nil
		}
	}
	switch {
	case p.matchOneOfKeywords(KeywordSelect, KeywordInsert):
		return p.parsePrivilegeSelectOrInsert(pos)
	case p.matchKeyword(KeywordAlter):
		if err := p.expectKeyword(KeywordAlter); err != nil {
			return nil, err
		}
		return p.parsePrivilegeAlter(pos)
	case p.matchKeyword(KeywordCreate):
		if err := p.expectKeyword(KeywordCreate); err != nil {
			return nil, err
		}
		return p.parsePrivilegeCreate(pos)
	case p.matchKeyword(KeywordDrop):
		if err := p.expectKeyword(KeywordDrop); err != nil {
			return nil, err
		}
		return p.parsePrivilegeDrop(pos)
	case p.matchKeyword(KeywordShow):
		if err := p.expectKeyword(KeywordShow); err != nil {
			return nil, err
		}
		return p.parsePrivilegeShow(pos)
	case p.matchKeyword(KeywordAll), p.matchTokenKind(KeywordNone):
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{KeywordAll},
		}, nil
	case p.matchKeyword(KeywordKill):
		if err := p.expectKeyword(KeywordKill); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordQuery); err != nil {
			return nil, err
		}
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{KeywordKill, KeywordQuery},
		}, nil
	case p.matchKeyword(KeywordSystem):
		if err := p.expectKeyword(KeywordSystem); err != nil {
			return nil, err
		}
		return p.parsePrivilegeSystem(pos)
	case p.matchKeyword(KeywordAdmin):
		if err := p.expectKeyword(KeywordAdmin); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordOption); err != nil {
			return nil, err
		}
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{KeywordAdmin, KeywordOption},
		}, nil
	case p.matchOneOfKeywords(KeywordOptimize, KeywordTruncate):
		keyword := p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{keyword},
		}, nil
	case p.matchKeyword(KeywordRole):
		if err := p.expectKeyword(KeywordRole); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordAdmin); err != nil {
			return nil, err
		}
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{KeywordRole, KeywordAdmin},
		}, nil
	}
	return nil, fmt.Errorf("expected SELECT|INSERT|ALTER|CREATE|DROP|SHOW|KILL|SYSTEM|OPTIMIZE|TRUNCATE")
}

func (p *Parser) parsePrivilegeRoles(_ Pos) ([]*Ident, error) {
	roles := make([]*Ident, 0)
	role, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	roles = append(roles, role)
	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		role, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (p *Parser) parseGrantOptions(_ Pos) ([]string, error) {
	options := make([]string, 0)
	for p.matchKeyword(KeywordWith) {
		option, err := p.parseGrantOption(p.Pos())
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}
	return options, nil
}

func (p *Parser) parseGrantOption(_ Pos) (string, error) {
	if err := p.expectKeyword(KeywordWith); err != nil {
		return "", err
	}
	// Between WITH and OPTION the token can only be the option name, which may
	// be a reserved keyword (e.g. `WITH GRANT OPTION`).
	ident, err := p.parseAnyKeyword()
	if err != nil {
		return "", err
	}
	if err := p.expectKeyword(KeywordOption); err != nil {
		return "", err
	}
	return ident.Name, nil
}

func (p *Parser) parseGrantSource(_ Pos) (*TableIdentifier, error) {
	ident, err := p.parseIdentOrStar()
	if err != nil {
		return nil, err
	}

	if token, consumeErr := p.tryConsumeTokenKind(TokenKindDot); consumeErr != nil {
		return nil, consumeErr
	} else if token == nil {
		return &TableIdentifier{
			Table: ident,
		}, nil
	}
	var dotIdent *Ident
	if p.matchTokenKind("*") {
		dotIdent, err = p.parseIdentOrStar()
	} else {
		dotIdent, err = p.parseAnyKeyword()
	}
	if err != nil {
		return nil, err
	}
	return &TableIdentifier{
		Database: ident,
		Table:    dotIdent,
	}, nil
}

func (p *Parser) parseGrantPrivilegeStmt(pos Pos) (*GrantPrivilegeStmt, error) {
	if err := p.expectKeyword(KeywordGrant); err != nil {
		return nil, err
	}
	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	var privileges []*PrivilegeClause
	privilege, err := p.parsePrivilegeClause(p.Pos())
	if err != nil {
		return nil, err
	}
	privileges = append(privileges, privilege)
	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		privilege, err := p.parsePrivilegeClause(p.Pos())
		if err != nil {
			return nil, err
		}
		privileges = append(privileges, privilege)
	}
	statementEnd := privileges[len(privileges)-1].End()

	if err := p.expectKeyword(KeywordOn); err != nil {
		return nil, err
	}
	on, err := p.parseGrantSource(p.Pos())
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordTo); err != nil {
		return nil, err
	}
	toRoles, err := p.parsePrivilegeRoles(p.Pos())
	if err != nil {
		return nil, err
	}
	if len(toRoles) != 0 {
		statementEnd = toRoles[len(toRoles)-1].NameEnd
	}
	options, err := p.parseGrantOptions(p.Pos())
	if err != nil {
		return nil, err
	}
	if len(options) != 0 {
		statementEnd = p.End()
	}

	return &GrantPrivilegeStmt{
		GrantPos:     pos,
		StatementEnd: statementEnd,
		OnCluster:    onCluster,
		Privileges:   privileges,
		On:           on,
		To:           toRoles,
		WithOptions:  options,
	}, nil
}

func (p *Parser) parseAlterRole(pos Pos) (*AlterRole, error) {
	if err := p.expectKeyword(KeywordRole); err != nil {
		return nil, err
	}

	ifExists, err := p.tryParseIfExists()
	if err != nil {
		return nil, err
	}

	roleRenamePairs := make([]*RoleRenamePair, 0)
	roleRenamePair, err := p.parseRoleRenamePair(p.Pos())
	if err != nil {
		return nil, err
	}
	roleRenamePairs = append(roleRenamePairs, roleRenamePair)
	for {
		token, consumeErr := p.tryConsumeTokenKind(TokenKindComma)
		if consumeErr != nil {
			return nil, consumeErr
		}
		if token == nil {
			break
		}
		roleRenamePair, err := p.parseRoleRenamePair(p.Pos())
		if err != nil {
			return nil, err
		}
		roleRenamePairs = append(roleRenamePairs, roleRenamePair)
	}
	statementEnd := roleRenamePairs[len(roleRenamePairs)-1].End()

	settings, err := p.tryParseRoleSettings(p.Pos())
	if err != nil {
		return nil, err
	}
	if settings != nil {
		statementEnd = settings[len(settings)-1].End()
	}

	return &AlterRole{
		AlterPos:        pos,
		StatementEnd:    statementEnd,
		IfExists:        ifExists,
		RoleRenamePairs: roleRenamePairs,
		Settings:        settings,
	}, nil
}

func (p *Parser) parseRoleRenamePair(_ Pos) (*RoleRenamePair, error) {
	roleName, err := p.parseRoleName(p.Pos())
	if err != nil {
		return nil, err
	}
	roleRenamePair := &RoleRenamePair{
		RoleName:     roleName,
		StatementEnd: roleName.End(),
	}
	if matched, consumeErr := p.tryConsumeKeywords(KeywordRename); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		if err := p.expectKeyword(KeywordTo); err != nil {
			return nil, err
		}
		newName, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		roleRenamePair.NewName = newName
		roleRenamePair.StatementEnd = newName.NameEnd
	}
	return roleRenamePair, nil
}
