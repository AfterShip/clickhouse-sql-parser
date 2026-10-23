package parser

import "fmt"

func (p *Parser) consumeRowPolicyKeyword(keyword string) error {
	if !p.matchKeyword(keyword) {
		return &ParseError{Pos: p.Pos(), Got: p.current(), Keyword: keyword}
	}
	return p.lexer.consumeToken()
}

func (p *Parser) consumeRowPolicyPrefix() error {
	if p.matchKeyword(KeywordRow) {
		if err := p.consumeRowPolicyKeyword(KeywordRow); err != nil {
			return err
		}
	}
	return p.consumeRowPolicyKeyword(KeywordPolicy)
}

func (p *Parser) rowPolicyOnCluster() bool {
	state := p.lexer.saveState()
	defer p.lexer.restoreState(state)
	return p.tryConsumeKeywords(KeywordOn, KeywordCluster)
}

func rowPolicyIdentEnd(ident *Ident) Pos {
	if ident.QuoteType == BackTicks || ident.QuoteType == DoubleQuote {
		return ident.End() + 1
	}
	return ident.End()
}

func rowPolicyExprEnd(expr Expr) Pos {
	end := expr.End()
	Walk(expr, func(node Expr) bool {
		if ident, ok := node.(*Ident); ok && rowPolicyIdentEnd(ident) > end {
			end = rowPolicyIdentEnd(ident)
		}
		return true
	})
	return end
}

func rowPolicyTableEnd(table *TableIdentifier) Pos {
	return rowPolicyIdentEnd(table.Table)
}

func rowPolicyWildcard(ident *Ident) bool {
	return ident.Name == "*" && (ident.QuoteType == 0 || ident.QuoteType == Unquoted)
}

func (p *Parser) parseRowPolicyTable(allowGlobal bool) (*TableIdentifier, error) {
	first, err := p.parseIdentOrStar()
	if err != nil {
		return nil, err
	}
	if p.matchTokenKind(TokenKindMul) {
		return nil, fmt.Errorf("prefix wildcards are not supported in row policy targets")
	}
	if !p.matchTokenKind(TokenKindDot) {
		return &TableIdentifier{Table: first}, nil
	}
	if rowPolicyWildcard(first) && !allowGlobal {
		return nil, fmt.Errorf("row policy targets cannot use a wildcard database")
	}
	if _, err := p.tryConsumeTokenKind(TokenKindDot); err != nil {
		return nil, err
	}
	var table *Ident
	if p.matchTokenKind(TokenKindMul) {
		table, err = p.parseIdentOrStar()
	} else {
		table, err = p.parseAnyKeyword()
	}
	if err != nil {
		return nil, err
	}
	if p.matchTokenKind(TokenKindMul) {
		return nil, fmt.Errorf("prefix wildcards are not supported in row policy targets")
	}
	if rowPolicyWildcard(first) && !rowPolicyWildcard(table) {
		return nil, fmt.Errorf("global policy inspection requires *.*")
	}
	return &TableIdentifier{Database: first, Table: table}, nil
}

func (p *Parser) parseRowPolicyNames(multiple bool) ([]*Ident, error) {
	name, err := p.parseIdentOrString()
	if err != nil {
		return nil, err
	}
	names := []*Ident{name}
	for multiple && p.matchTokenKind(TokenKindComma) {
		if _, err := p.tryConsumeTokenKind(TokenKindComma); err != nil {
			return nil, err
		}
		name, err = p.parseIdentOrString()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

func (p *Parser) parseRowPolicyTables(multiple bool) ([]*TableIdentifier, error) {
	table, err := p.parseRowPolicyTable(false)
	if err != nil {
		return nil, err
	}
	tables := []*TableIdentifier{table}
	for multiple && p.matchTokenKind(TokenKindComma) {
		state := p.lexer.saveState()
		if _, err := p.tryConsumeTokenKind(TokenKindComma); err != nil {
			return nil, err
		}
		candidate, err := p.parseRowPolicyTable(false)
		if err != nil || (candidate.Database == nil && p.matchKeyword(KeywordOn) && !p.rowPolicyOnCluster()) {
			p.lexer.restoreState(state)
			break
		}
		tables = append(tables, candidate)
	}
	return tables, nil
}

func bindRowPolicyTargets(names []*Ident, tables []*TableIdentifier) []*RowPolicyTarget {
	var targets []*RowPolicyTarget
	for _, name := range names {
		for _, table := range tables {
			policyCopy := *name
			tableCopy := *table
			columnCopy := *table.Table
			tableCopy.Table = &columnCopy
			if table.Database != nil {
				dbCopy := *table.Database
				tableCopy.Database = &dbCopy
			}
			targets = append(targets, &RowPolicyTarget{PolicyName: &policyCopy, Table: &tableCopy})
		}
	}
	return targets
}

func (p *Parser) parseRowPolicyTargetGroup(
	multiple, allowCluster bool,
) ([]*RowPolicyTarget, *ClusterClause, Pos, error) {
	names, err := p.parseRowPolicyNames(multiple)
	if err != nil {
		return nil, nil, 0, err
	}
	var cluster *ClusterClause
	if allowCluster && p.rowPolicyOnCluster() {
		cluster, _, err = p.parseRowPolicyCluster(nil)
		if err != nil {
			return nil, nil, 0, err
		}
	}
	if err := p.consumeRowPolicyKeyword(KeywordOn); err != nil {
		return nil, nil, 0, err
	}
	tables, err := p.parseRowPolicyTables(multiple && len(names) == 1)
	if err != nil {
		return nil, nil, 0, err
	}
	end := rowPolicyTableEnd(tables[len(tables)-1])
	if allowCluster && cluster == nil && p.rowPolicyOnCluster() {
		cluster, end, err = p.parseRowPolicyCluster(nil)
		if err != nil {
			return nil, nil, 0, err
		}
	}
	return bindRowPolicyTargets(names, tables), cluster, end, nil
}

func (p *Parser) parseRowPolicyTargets(allowCluster bool) ([]*RowPolicyTarget, *ClusterClause, Pos, error) {
	targets, cluster, end, err := p.parseRowPolicyTargetGroup(true, allowCluster)
	if err != nil {
		return nil, nil, 0, err
	}
	canAppendPair := len(targets) == 1
	for p.matchTokenKind(TokenKindComma) {
		if !canAppendPair || cluster != nil {
			return nil, nil, 0, fmt.Errorf("cannot combine multiple policy names/tables or clusters with another target group")
		}
		if _, err := p.tryConsumeTokenKind(TokenKindComma); err != nil {
			return nil, nil, 0, err
		}
		var pair []*RowPolicyTarget
		pair, cluster, end, err = p.parseRowPolicyTargetGroup(false, allowCluster)
		if err != nil {
			return nil, nil, 0, err
		}
		targets = append(targets, pair...)
	}
	return targets, cluster, end, nil
}

func (p *Parser) parseRowPolicyCondition() (*RowPolicyCondition, error) {
	condition := &RowPolicyCondition{ConditionPos: p.Pos()}
	if p.matchKeyword(KeywordNone) {
		condition.None = true
		condition.ConditionEnd = p.current().End
		if err := p.consumeRowPolicyKeyword(KeywordNone); err != nil {
			return nil, err
		}
		return condition, nil
	}
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	condition.Expr = expr
	condition.ConditionEnd = rowPolicyExprEnd(expr)
	return condition, nil
}

func (p *Parser) parseRowPolicyFilter() (*RowPolicyFilter, error) {
	filter := &RowPolicyFilter{FilterPos: p.Pos()}
	if p.matchKeyword(KeywordFor) {
		if err := p.consumeRowPolicyKeyword(KeywordFor); err != nil {
			return nil, err
		}
		for {
			if !p.matchOneOfKeywords(KeywordSelect, KeywordAll) {
				return nil, fmt.Errorf("row policy supports FOR SELECT or ALL, got %q", p.currentTokenString())
			}
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			if !p.matchTokenKind(TokenKindComma) {
				break
			}
			if _, err := p.tryConsumeTokenKind(TokenKindComma); err != nil {
				return nil, err
			}
		}
	}
	var err error
	if p.matchKeyword(KeywordUsing) {
		if err = p.consumeRowPolicyKeyword(KeywordUsing); err != nil {
			return nil, err
		}
		filter.Using, err = p.parseRowPolicyCondition()
		if err != nil {
			return nil, err
		}
		filter.FilterEnd = filter.Using.End()
	}
	if p.matchKeyword(KeywordWith) {
		if err = p.consumeRowPolicyKeyword(KeywordWith); err != nil {
			return nil, err
		}
		if err = p.consumeRowPolicyKeyword(KeywordCheck); err != nil {
			return nil, err
		}
		filter.WithCheck, err = p.parseRowPolicyCondition()
		if err != nil {
			return nil, err
		}
		filter.FilterEnd = filter.WithCheck.End()
	}
	if filter.Using == nil && filter.WithCheck == nil {
		return nil, fmt.Errorf("row policy filter requires USING or WITH CHECK")
	}
	return filter, nil
}

func (p *Parser) parseRowPolicyPrincipal() (*RowPolicyPrincipal, error) {
	principal := &RowPolicyPrincipal{PrincipalPos: p.Pos()}
	if p.matchOneOfKeywords(KeywordCurrentUser, KeywordCurrentUserAlt) {
		principal.CurrentUser = true
		principal.PrincipalEnd = p.current().End
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
		if p.matchTokenKind(TokenKindLParen) {
			if _, err := p.tryConsumeTokenKind(TokenKindLParen); err != nil {
				return nil, err
			}
			if !p.matchTokenKind(TokenKindRParen) {
				return nil, fmt.Errorf("CURRENT_USER takes no arguments")
			}
			principal.PrincipalEnd = p.current().End
			if _, err := p.tryConsumeTokenKind(TokenKindRParen); err != nil {
				return nil, err
			}
		}
		return principal, nil
	}
	var err error
	principal.Name, err = p.parseIdentOrString()
	if err != nil {
		return nil, err
	}
	principal.PrincipalEnd = rowPolicyIdentEnd(principal.Name)
	if principal.Name.QuoteType != SingleQuote && p.matchTokenKind(TokenKindAtSign) {
		if _, err := p.tryConsumeTokenKind(TokenKindAtSign); err != nil {
			return nil, err
		}
		principal.Host, err = p.parseIdentOrString()
		if err != nil {
			return nil, err
		}
		principal.PrincipalEnd = rowPolicyIdentEnd(principal.Host)
	}
	return principal, nil
}

func (p *Parser) parseRowPolicyPrincipalList(allowAll bool) ([]*RowPolicyPrincipal, bool, Pos, error) {
	var principals []*RowPolicyPrincipal
	var all bool
	var end Pos
	for {
		switch {
		case p.matchKeyword(KeywordNone):
			end = p.current().End
			if err := p.consumeRowPolicyKeyword(KeywordNone); err != nil {
				return nil, false, 0, err
			}
		case allowAll && p.matchKeyword(KeywordAll):
			all = true
			end = p.current().End
			if err := p.consumeRowPolicyKeyword(KeywordAll); err != nil {
				return nil, false, 0, err
			}
		default:
			principal, err := p.parseRowPolicyPrincipal()
			if err != nil {
				return nil, false, 0, err
			}
			principals = append(principals, principal)
			end = principal.End()
		}
		if !p.matchTokenKind(TokenKindComma) {
			break
		}
		if _, err := p.tryConsumeTokenKind(TokenKindComma); err != nil {
			return nil, false, 0, err
		}
	}
	if all {
		principals = nil
	}
	return principals, all, end, nil
}

func (p *Parser) parseRowPolicySubjects() (*RowPolicySubjects, error) {
	subjects := &RowPolicySubjects{ToPos: p.Pos()}
	if err := p.consumeRowPolicyKeyword(KeywordTo); err != nil {
		return nil, err
	}
	var err error
	subjects.Principals, subjects.All, subjects.ToEnd, err = p.parseRowPolicyPrincipalList(true)
	if err != nil {
		return nil, err
	}
	if p.matchKeyword(KeywordExcept) {
		if err := p.consumeRowPolicyKeyword(KeywordExcept); err != nil {
			return nil, err
		}
		subjects.Except, _, subjects.ToEnd, err = p.parseRowPolicyPrincipalList(false)
		if err != nil {
			return nil, err
		}
	}
	return subjects, nil
}

func (p *Parser) parseRowPolicyRename(allowed bool) (*Ident, Pos, error) {
	if !allowed {
		return nil, 0, fmt.Errorf("RENAME TO requires a single ALTER ROW POLICY target")
	}
	if err := p.consumeRowPolicyKeyword(KeywordRename); err != nil {
		return nil, 0, err
	}
	if err := p.consumeRowPolicyKeyword(KeywordTo); err != nil {
		return nil, 0, err
	}
	name, err := p.parseIdentOrString()
	if err != nil {
		return nil, 0, err
	}
	return name, rowPolicyIdentEnd(name), nil
}

func (p *Parser) parseRowPolicyMode(existing RowPolicyMode) (RowPolicyMode, Pos, error) {
	if existing != "" {
		return "", 0, fmt.Errorf("duplicate row policy mode")
	}
	if err := p.consumeRowPolicyKeyword(KeywordAs); err != nil {
		return "", 0, err
	}
	var mode RowPolicyMode
	switch {
	case p.matchKeyword(KeywordPermissive):
		mode = RowPolicyPermissive
	case p.matchKeyword(KeywordRestrictive):
		mode = RowPolicyRestrictive
	default:
		return "", 0, fmt.Errorf("expected PERMISSIVE or RESTRICTIVE, got %q", p.currentTokenString())
	}
	end := p.current().End
	if err := p.lexer.consumeToken(); err != nil {
		return "", 0, err
	}
	return mode, end, nil
}

func (p *Parser) parseRowPolicyStorage(keyword string, existing *Ident) (*Ident, Pos, error) {
	if existing != nil {
		return nil, 0, fmt.Errorf("duplicate row policy access storage")
	}
	if err := p.consumeRowPolicyKeyword(keyword); err != nil {
		return nil, 0, err
	}
	storage, err := p.parseIdentOrString()
	if err != nil {
		return nil, 0, err
	}
	return storage, rowPolicyIdentEnd(storage), nil
}

func (p *Parser) parseRowPolicyCluster(existing *ClusterClause) (*ClusterClause, Pos, error) {
	if existing != nil {
		return nil, 0, fmt.Errorf("duplicate ON CLUSTER in row policy")
	}
	cluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, 0, err
	}
	if cluster == nil {
		return nil, 0, fmt.Errorf("expected ON CLUSTER")
	}
	return cluster, rowPolicyExprEnd(cluster), nil
}

func (p *Parser) parseRowPolicyTail(options *RowPolicyOptions, end Pos) (Pos, error) {
	if p.matchKeyword(KeywordTo) {
		subjects, err := p.parseRowPolicySubjects()
		if err != nil {
			return 0, err
		}
		options.To, end = subjects, subjects.End()
	}
	if p.rowPolicyOnCluster() {
		var err error
		options.OnCluster, end, err = p.parseRowPolicyCluster(options.OnCluster)
		if err != nil {
			return 0, err
		}
	}
	return end, nil
}

func (p *Parser) parseRowPolicyOptions(options *RowPolicyOptions, alter bool, end Pos) (*Ident, Pos, error) {
	var rename *Ident
	for {
		var err error
		switch {
		case p.matchKeyword(KeywordRename):
			rename, end, err = p.parseRowPolicyRename(alter && len(options.Targets) == 1 && rename == nil)
		case p.matchKeyword(KeywordAs):
			options.Mode, end, err = p.parseRowPolicyMode(options.Mode)
		case p.matchOneOfKeywords(KeywordFor, KeywordUsing, KeywordWith):
			var filter *RowPolicyFilter
			filter, err = p.parseRowPolicyFilter()
			if err == nil {
				options.Filters = append(options.Filters, filter)
				end = filter.End()
				if p.matchTokenKind(TokenKindComma) {
					_, err = p.tryConsumeTokenKind(TokenKindComma)
					if err == nil && !p.matchOneOfKeywords(KeywordFor, KeywordUsing, KeywordWith) {
						err = fmt.Errorf("expected another row policy filter after comma")
					}
				}
			}
		case p.matchKeyword(KeywordIn):
			options.AccessStorageType, end, err = p.parseRowPolicyStorage(KeywordIn, options.AccessStorageType)
		case p.rowPolicyOnCluster():
			options.OnCluster, end, err = p.parseRowPolicyCluster(options.OnCluster)
		default:
			end, err = p.parseRowPolicyTail(options, end)
			return rename, end, err
		}
		if err != nil {
			return nil, 0, err
		}
	}
}

func (p *Parser) parseCreateRowPolicy(pos Pos) (*CreateRowPolicy, error) {
	if err := p.consumeRowPolicyPrefix(); err != nil {
		return nil, err
	}
	statement := &CreateRowPolicy{CreatePos: pos}
	var err error
	statement.IfNotExists, err = p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}
	if !statement.IfNotExists && p.matchKeyword(KeywordOr) {
		if err := p.consumeRowPolicyKeyword(KeywordOr); err != nil {
			return nil, err
		}
		if err := p.consumeRowPolicyKeyword(KeywordReplace); err != nil {
			return nil, err
		}
		statement.OrReplace = true
	}
	statement.Targets, statement.OnCluster, statement.StatementEnd, err = p.parseRowPolicyTargets(true)
	if err != nil {
		return nil, err
	}
	_, statement.StatementEnd, err = p.parseRowPolicyOptions(&statement.RowPolicyOptions, false, statement.StatementEnd)
	return statement, err
}

func (p *Parser) parseAlterRowPolicy(pos Pos) (*AlterRowPolicy, error) {
	if err := p.consumeRowPolicyPrefix(); err != nil {
		return nil, err
	}
	statement := &AlterRowPolicy{AlterPos: pos}
	var err error
	statement.IfExists, err = p.tryParseIfExists()
	if err != nil {
		return nil, err
	}
	statement.Targets, statement.OnCluster, statement.StatementEnd, err = p.parseRowPolicyTargets(true)
	if err != nil {
		return nil, err
	}
	statement.RenameTo, statement.StatementEnd, err = p.parseRowPolicyOptions(
		&statement.RowPolicyOptions, true, statement.StatementEnd,
	)
	return statement, err
}

func (p *Parser) parseDropRowPolicy(pos Pos) (*DropRowPolicy, error) {
	if err := p.consumeRowPolicyPrefix(); err != nil {
		return nil, err
	}
	statement := &DropRowPolicy{DropPos: pos}
	var err error
	statement.IfExists, err = p.tryParseIfExists()
	if err != nil {
		return nil, err
	}
	statement.Targets, statement.OnCluster, statement.StatementEnd, err = p.parseRowPolicyTargets(true)
	if err != nil {
		return nil, err
	}
	if p.matchKeyword(KeywordFrom) {
		statement.AccessStorageType, statement.StatementEnd, err = p.parseRowPolicyStorage(KeywordFrom, nil)
		if err != nil {
			return nil, err
		}
	}
	if p.rowPolicyOnCluster() {
		statement.OnCluster, statement.StatementEnd, err = p.parseRowPolicyCluster(statement.OnCluster)
		if err != nil {
			return nil, err
		}
	}
	return statement, nil
}

func (p *Parser) parseShowWithPolicies(pos Pos) (Expr, error) {
	state := p.lexer.saveState()
	if err := p.consumeRowPolicyKeyword(KeywordShow); err != nil {
		return nil, err
	}
	create := p.matchKeyword(KeywordCreate)
	if create {
		if err := p.consumeRowPolicyKeyword(KeywordCreate); err != nil {
			return nil, err
		}
	}
	if p.matchOneOfKeywords(KeywordRow, KeywordPolicy, KeywordPolicies) {
		return p.parseShowRowPolicies(pos, create)
	}
	p.lexer.restoreState(state)
	return p.parseShowStmt(pos)
}

func (p *Parser) parseShowRowPolicies(pos Pos, create bool) (*ShowRowPolicies, error) {
	statement := &ShowRowPolicies{ShowPos: pos, Create: create}
	if p.matchKeyword(KeywordRow) {
		if err := p.consumeRowPolicyKeyword(KeywordRow); err != nil {
			return nil, err
		}
	}
	plural := p.matchKeyword(KeywordPolicies)
	if !plural && (!create || !p.matchKeyword(KeywordPolicy)) {
		return nil, fmt.Errorf("expected ROW POLICIES or CREATE ROW POLICY")
	}
	statement.StatementEnd = p.current().End
	if err := p.lexer.consumeToken(); err != nil {
		return nil, err
	}
	if p.matchKeyword(KeywordOn) {
		if err := p.consumeRowPolicyKeyword(KeywordOn); err != nil {
			return nil, err
		}
		var err error
		statement.Table, err = p.parseRowPolicyTable(true)
		if err != nil {
			return nil, err
		}
		statement.StatementEnd = rowPolicyTableEnd(statement.Table)
		return statement, nil
	}
	if create && !p.matchTokenKind(TokenKindEOF, ";") && p.current() != nil {
		state := p.lexer.saveState()
		targets, _, end, err := p.parseRowPolicyTargets(false)
		if err == nil {
			statement.Targets, statement.StatementEnd = targets, end
			return statement, nil
		}
		p.lexer.restoreState(state)
		statement.Name, err = p.parseIdentOrString()
		if err != nil {
			return nil, err
		}
		statement.StatementEnd = rowPolicyIdentEnd(statement.Name)
		return statement, nil
	}
	if !plural {
		return nil, fmt.Errorf("SHOW CREATE ROW POLICY requires a policy or table selector")
	}
	return statement, nil
}
