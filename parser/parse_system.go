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

func (p *Parser) parseSystemFlushExpr(pos Pos) (*SystemFlushExpr, error) {
	if err := p.expectKeyword(KeywordFlush); err != nil {
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
	case p.tryConsumeKeywords(KeywordDistributed):
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

	switch {
	case p.matchKeyword(KeywordDictionaries):
		lastToken := p.last()
		_ = p.lexer.consumeToken()
		return &SystemReloadExpr{
			ReloadPos:    pos,
			StatementEnd: lastToken.End,
			Type:         KeywordDictionaries,
		}, nil
	case p.tryConsumeKeywords(KeywordDictionary):
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
	case p.tryConsumeKeywords(KeywordEmbedded):
		lastToken := p.last()
		if err := p.expectKeyword(KeywordDictionaries); err != nil {
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
	command := strings.ToUpper(p.last().String)
	_ = p.lexer.consumeToken()

	var typ string
	switch {
	case p.tryConsumeKeywords(KeywordDistributed):
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
	case p.tryConsumeKeywords(KeywordReplicated):
		lastToken := p.last()
		if err := p.expectKeyword(KeywordSends); err != nil {
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
	if err := p.expectKeyword(KeywordDrop); err != nil {
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
		if err := p.expectKeyword(KeywordCache); err != nil {
			return nil, err
		}
		return &SystemDropExpr{
			DropPos:      pos,
			StatementEnd: lastToken.End,
			Type:         prefixToken.String + " CACHE",
		}, nil
	case p.matchKeyword(KeywordCompiled):
		_ = p.lexer.consumeToken()
		if err := p.expectKeyword(KeywordExpression); err != nil {
			return nil, err
		}
		lastToken := p.last()
		if err := p.expectKeyword(KeywordCache); err != nil {
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
	if !p.tryConsumeKeywords(KeywordBy) {
		return &DeduplicateClause{
			DeduplicatePos: pos,
		}, nil
	}

	by, err := p.parseColumnExprList(p.Pos())
	if err != nil {
		return nil, err
	}
	var except *ColumnExprList
	if p.tryConsumeKeywords(KeywordExcept) {
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
	if p.tryConsumeKeywords(KeywordFinal) {
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
		if p.tryConsumeTokenKind(TokenKindAtSign) != nil {
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
	if !p.tryConsumeKeywords(KeywordSettings) {
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
			if token := p.tryConsumeTokenKind(TokenKindSingleEQ); token != nil {
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
		if p.tryConsumeTokenKind(TokenKindComma) == nil {
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
		_ = p.lexer.consumeToken()
		if err := p.expectKeyword(KeywordNot); err != nil {
			return nil, err
		}
		if err := p.expectKeyword(KeywordExists); err != nil {
			return nil, err
		}
		ifNotExists = true
	case p.matchKeyword(KeywordOr):
		_ = p.lexer.consumeToken()
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
	for p.tryConsumeTokenKind(TokenKindComma) != nil {
		roleName, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		roleNames = append(roleNames, roleName)
	}
	statementEnd := roleNames[len(roleNames)-1].End()

	var accessStorageType *Ident
	if p.tryConsumeKeywords(KeywordIn) {
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

	if p.tryConsumeKeywords(KeywordNot) {
		if err := p.expectKeyword(KeywordIdentified); err != nil {
			return nil, err
		}
		auth.NotIdentified = true
		auth.AuthEnd = p.last().End
		return auth, nil
	}

	if err := p.expectKeyword(KeywordIdentified); err != nil {
		return nil, err
	}
	auth.AuthEnd = p.last().End

	if p.tryConsumeKeywords(KeywordWith) {
		if p.matchKeyword(KeywordLdap) {
			_ = p.lexer.consumeToken()
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
			_ = p.lexer.consumeToken()
			auth.IsKerberos = true
			auth.AuthEnd = p.last().End
			if p.tryConsumeKeywords(KeywordRealm) {
				realm, err := p.parseString(p.Pos())
				if err != nil {
					return nil, err
				}
				auth.KerberosRealm = realm
				auth.AuthEnd = realm.End()
			}
		} else if p.matchTokenKind(TokenKindIdent) {
			// Auth types like no_password, plaintext_password, etc.
			authType := p.last().String
			_ = p.lexer.consumeToken()
			auth.AuthType = authType
			auth.AuthEnd = p.last().End

			if p.tryConsumeKeywords(KeywordBy) {
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
		hostType := p.last().String
		_ = p.lexer.consumeToken()
		host.HostType = hostType
		host.HostEnd = p.last().End
	case p.matchOneOfKeywords(KeywordName, KeywordRegexp, KeywordIp, KeywordLike):
		hostType := p.last().String
		_ = p.lexer.consumeToken()
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

	if p.tryConsumeKeywords(KeywordNone) {
		defaultRole.None = true
		defaultRole.DefaultEnd = p.last().End
		return defaultRole, nil
	}

	roles := make([]*RoleName, 0)
	role, err := p.parseRoleName(p.Pos())
	if err != nil {
		return nil, err
	}
	roles = append(roles, role)

	for p.tryConsumeTokenKind(TokenKindComma) != nil {
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

	if p.tryConsumeKeywords(KeywordAny) {
		grantees.Any = true
		grantees.GranteesEnd = p.last().End
	} else if p.tryConsumeKeywords(KeywordNone) {
		grantees.None = true
		grantees.GranteesEnd = p.last().End
	} else {
		// Parse list of grantees
		granteeList := make([]*RoleName, 0)
		grantee, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		granteeList = append(granteeList, grantee)

		for p.tryConsumeTokenKind(TokenKindComma) != nil {
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
	if p.tryConsumeKeywords(KeywordExcept) {
		exceptList := make([]*RoleName, 0)
		except, err := p.parseRoleName(p.Pos())
		if err != nil {
			return nil, err
		}
		exceptList = append(exceptList, except)

		for p.tryConsumeTokenKind(TokenKindComma) != nil {
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
		_ = p.lexer.consumeToken()
		if err := p.expectKeyword(KeywordNot); err != nil {
			return err
		}
		if err := p.expectKeyword(KeywordExists); err != nil {
			return err
		}
		createUser.IfNotExists = true
	case p.matchKeyword(KeywordOr):
		_ = p.lexer.consumeToken()
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

	for p.tryConsumeTokenKind(TokenKindComma) != nil {
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

	for p.tryConsumeTokenKind(TokenKindComma) != nil {
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
		_ = p.lexer.consumeToken() // consume DEFAULT
		_ = p.lexer.consumeToken() // consume DATABASE
		if p.tryConsumeKeywords(KeywordNone) {
			createUser.DefaultDbNone = true
			createUser.StatementEnd = p.last().End
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
			_ = p.lexer.consumeToken() // consume SETTINGS keyword
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
		target = p.last().String
		_ = p.lexer.consumeToken()
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
	for p.tryConsumeTokenKind(TokenKindComma) != nil {
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
	if p.tryConsumeKeywords(KeywordFrom) {
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
	keyword := p.last().String
	_ = p.lexer.consumeToken()

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
	case p.tryConsumeKeywords(KeywordIndex):
		keywords = append(keywords, KeywordIndex)
	case p.matchOneOfKeywords(KeywordUpdate, KeywordDelete, KeywordUser, KeywordRole, KeywordQuota):
		keyword := p.last().String
		_ = p.lexer.consumeToken()
		keywords = append(keywords, keyword)
	case p.matchOneOfKeywords(KeywordAdd, KeywordDrop, KeywordModify, KeywordClear, KeywordComment, KeywordRename, KeywordMaterialized):
		keyword := p.last().String
		_ = p.lexer.consumeToken()
		keywords = append(keywords, keyword)
		switch {
		case p.tryConsumeKeywords(KeywordColumn):
			keywords = append(keywords, KeywordColumn)
		case p.tryConsumeKeywords(KeywordIndex):
			keywords = append(keywords, KeywordIndex)
			keywords = append(keywords, KeywordConstraint)
		case p.tryConsumeKeywords(KeywordTtl):
			keywords = append(keywords, KeywordTtl)
		default:
			return nil, fmt.Errorf("expected COLUMN|INDEX")
		}
	case p.tryConsumeKeywords(KeywordOrder):
		if err := p.expectKeyword(KeywordBy); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordOrder, KeywordBy)
	case p.tryConsumeKeywords(KeywordSample):
		if err := p.expectKeyword(KeywordBy); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordSample, KeywordBy)
	case p.tryConsumeKeywords(KeywordSettings):
		keywords = append(keywords, KeywordSettings)
	case p.tryConsumeKeywords(KeywordView):
		keywords = append(keywords, KeywordView)
		switch {
		case p.tryConsumeKeywords(KeywordModify):
			keywords = append(keywords, KeywordModify)
		case p.tryConsumeKeywords(KeywordRefresh):
			keywords = append(keywords, KeywordRefresh)
		default:
			return nil, fmt.Errorf("expected MODIFY|REFRESH")
		}
	case p.matchOneOfKeywords(KeywordMove, KeywordFreeze):
		keyword := p.last().String
		_ = p.lexer.consumeToken()
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
		keyword := p.last().String
		_ = p.lexer.consumeToken()
		keywords = append(keywords, keyword)
	case p.tryConsumeKeywords(KeywordTemporary):
		if err := p.expectKeyword(KeywordTable); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordTemporary, KeywordTable)
	case p.tryConsumeKeywords(KeywordRows):
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
		keyword := p.last().String
		_ = p.lexer.consumeToken()
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
		keyword := p.last().String
		_ = p.lexer.consumeToken()
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
		keyword := p.last().String
		_ = p.lexer.consumeToken()
		keywords = append(keywords, keyword)
	case p.tryConsumeKeywords(KeywordDrop):
		keywords = append(keywords, KeywordDrop)
		switch {
		case p.tryConsumeKeywords(KeywordCache):
			keywords = append(keywords, KeywordCache)
		case p.matchOneOfKeywords(KeywordMark, KeywordDNS, KeywordUncompressed):
			keyword := p.last().String
			_ = p.lexer.consumeToken()
			keywords = append(keywords, keyword)
			if err := p.expectKeyword(KeywordCache); err != nil {
				return nil, err
			}
			keywords = append(keywords, KeywordCache)
		default:
			return nil, fmt.Errorf("expected CACHE|MARK|DNS|UNCOMPRESSED")
		}
	case p.tryConsumeKeywords(KeywordReload):
		keywords = append(keywords, KeywordReload)
		switch {
		case p.matchOneOfKeywords(KeywordDictionary, KeywordFunction, KeywordFunctions, KeywordConfig):
			keyword := p.last().String
			_ = p.lexer.consumeToken()
			keywords = append(keywords, keyword)
		default:
			return nil, fmt.Errorf("expected DICTIONARY|FUNCTION|FUNCTIONS|CONFIG")
		}
	case p.tryConsumeKeywords(KeywordFlush):
		keywords = append(keywords, KeywordFlush)
		switch {
		case p.matchOneOfKeywords(KeywordLogs, KeywordDistributed):
			keyword := p.last().String
			_ = p.lexer.consumeToken()
			keywords = append(keywords, keyword)
		default:
			return nil, fmt.Errorf("expected LOGS|DISTRIBUTED")
		}
	case p.tryConsumeKeywords(KeywordTtl):
		keywords = append(keywords, KeywordTtl)
		if err := p.expectKeyword(KeywordMerges); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordMerges)
	case p.matchOneOfKeywords(KeywordSync, KeywordRestart):
		keyword := p.last().String
		_ = p.lexer.consumeToken()
		keywords = append(keywords, keyword)
		if err := p.expectKeyword(KeywordReplica); err != nil {
			return nil, err
		}
		keywords = append(keywords, KeywordReplica)
	case p.tryConsumeKeywords(KeywordReplication):
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
		if p.last().String == "dictGet" {
			_ = p.lexer.consumeToken()
			return &PrivilegeClause{
				PrivilegePos: pos,
				Keywords:     []string{"dictGet"},
			}, nil
		}
	}
	switch {
	case p.matchOneOfKeywords(KeywordSelect, KeywordInsert):
		return p.parsePrivilegeSelectOrInsert(pos)
	case p.tryConsumeKeywords(KeywordAlter):
		return p.parsePrivilegeAlter(pos)
	case p.tryConsumeKeywords(KeywordCreate):
		return p.parsePrivilegeCreate(pos)
	case p.tryConsumeKeywords(KeywordDrop):
		return p.parsePrivilegeDrop(pos)
	case p.tryConsumeKeywords(KeywordShow):
		return p.parsePrivilegeShow(pos)
	case p.matchKeyword(KeywordAll), p.matchTokenKind(KeywordNone):
		_ = p.lexer.consumeToken()
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{KeywordAll},
		}, nil
	case p.tryConsumeKeywords(KeywordKill):
		if err := p.expectKeyword(KeywordQuery); err != nil {
			return nil, err
		}
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{KeywordKill, KeywordQuery},
		}, nil
	case p.tryConsumeKeywords(KeywordSystem):
		return p.parsePrivilegeSystem(pos)
	case p.tryConsumeKeywords(KeywordAdmin):
		if err := p.expectKeyword(KeywordOption); err != nil {
			return nil, err
		}
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{KeywordAdmin, KeywordOption},
		}, nil
	case p.matchOneOfKeywords(KeywordOptimize, KeywordTruncate):
		keyword := p.last().String
		_ = p.lexer.consumeToken()
		return &PrivilegeClause{
			PrivilegePos: pos,
			Keywords:     []string{keyword},
		}, nil
	case p.tryConsumeKeywords(KeywordRole):
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
	for p.tryConsumeTokenKind(TokenKindComma) != nil {
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
	ident, err := p.parseIdent()
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

	if p.tryConsumeTokenKind(TokenKindDot) == nil {
		return &TableIdentifier{
			Table: ident,
		}, nil
	}
	dotIdent, err := p.parseIdentOrStar()
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
	for p.tryConsumeTokenKind(TokenKindComma) != nil {
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
	for p.tryConsumeTokenKind(TokenKindComma) != nil {
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
	if p.tryConsumeKeywords(KeywordRename) {
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
