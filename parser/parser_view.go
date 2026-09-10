package parser

import "fmt"

// parseCreateMaterializedView parses a CREATE MATERIALIZED VIEW statement.
//
// The syntax is as follows:
// CREATE MATERIALIZED VIEW [IF NOT EXISTS] [db.]table_name [ON CLUSTER cluster]
// REFRESH EVERY|AFTER interval [OFFSET interval]
// [RANDOMIZE FOR interval]
// [DEPENDS ON [db.]name [, [db.]name [, ...]]]
// [SETTINGS name = value [, name = value [, ...]]]
// [APPEND]
// [TO[db.]name] [(columns)] [ENGINE = engine]
// [EMPTY]
// [DEFINER = { user | CURRENT_USER }] [SQL SECURITY { DEFINER | NONE }]
// AS SELECT ...
// [COMMENT 'comment']
//
//nolint:funlen
func (p *Parser) parseCreateMaterializedView(pos Pos, orReplace bool) (*CreateMaterializedView, error) {
	if err := p.expectKeyword(KeywordMaterialized); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordView); err != nil {
		return nil, err
	}

	createMaterializedView := &CreateMaterializedView{CreatePos: pos, OrReplace: orReplace}

	// parse IF NOT EXISTS clause if exists
	var err error
	createMaterializedView.IfNotExists, err = p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}

	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	createMaterializedView.Name = tableIdentifier

	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	createMaterializedView.OnCluster = onCluster

	refreshExpr, err := p.tryParseRefreshExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	createMaterializedView.Refresh = refreshExpr

	if matched, consumeErr := p.tryConsumeKeywords(KeywordRandomize, KeywordFor); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		randomizeFor, err := p.parseInterval(false)
		if err != nil {
			return nil, err
		}
		createMaterializedView.RandomizeFor = randomizeFor
	}
	if matched, consumeErr := p.tryConsumeKeywords(KeywordDepends, KeywordOn); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		dependsOnTables := make([]*TableIdentifier, 0)
		table, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		dependsOnTables = append(dependsOnTables, table)
		for p.matchTokenKind(TokenKindComma) {
			if err := p.lexer.consumeToken(); err != nil {
				return nil, err
			}
			table, err := p.parseTableIdentifier(p.Pos())
			if err != nil {
				return nil, err
			}
			dependsOnTables = append(dependsOnTables, table)
		}
		createMaterializedView.DependsOn = dependsOnTables
	}
	settings, err := p.tryParseSettingsClause(p.Pos())
	if err != nil {
		return nil, err
	}
	createMaterializedView.Settings = settings
	createMaterializedView.HasAppend, err = p.tryConsumeKeywords(KeywordAppend)
	if err != nil {
		return nil, err
	}

	switch {
	case p.matchKeyword(KeywordTo):
		destination, err := p.parseDestinationClause(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.Destination = destination
		createMaterializedView.StatementEnd = destination.End()
		if p.matchTokenKind(TokenKindLParen) {
			tableSchema, err := p.parseTableSchemaClause(p.Pos())
			if err != nil {
				return nil, err
			}
			createMaterializedView.Destination.TableSchema = tableSchema
		}
	case p.matchTokenKind(TokenKindLParen):
		// Column list before ENGINE (e.g. SHOW CREATE TABLE output for RMVs with ENGINE = Memory)
		tableSchema, err := p.parseTableSchemaClause(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.TableSchema = tableSchema
		if !p.matchKeyword(KeywordEngine) {
			return nil, fmt.Errorf("unexpected token: %q, expected ENGINE after column list", p.currentTokenKind())
		}
		engineExpr, err := p.parseEngineExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.Engine = engineExpr
		createMaterializedView.StatementEnd = engineExpr.End()
	case p.matchKeyword(KeywordEngine):
		engineExpr, err := p.parseEngineExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.Engine = engineExpr
		createMaterializedView.StatementEnd = engineExpr.End()
	default:
		return nil, fmt.Errorf("unexpected token: %q, expected TO or ENGINE", p.currentTokenKind())
	}
	createMaterializedView.HasEmpty, err = p.tryConsumeKeywords(KeywordEmpty)
	if err != nil {
		return nil, err
	}

	// Parse DEFINER clause
	if matched, consumeErr := p.tryConsumeKeywords(KeywordDefiner); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		if err := p.expectTokenKind(TokenKindSingleEQ); err != nil {
			return nil, err
		}
		definer, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		createMaterializedView.Definer = definer
	}

	// Parse SQL SECURITY clause
	if matched, consumeErr := p.tryConsumeKeywords(KeywordSQL, KeywordSecurity); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		if !p.matchOneOfKeywords(KeywordDefiner, KeywordNone) {
			return nil, fmt.Errorf("expected DEFINER or NONE after SQL SECURITY, got %q", p.currentTokenKind())
		}
		createMaterializedView.SQLSecurity = p.current().String
		if err := p.lexer.consumeToken(); err != nil {
			return nil, err
		}
	}

	// Check for POPULATE before AS SELECT - only valid with ENGINE and no Destination
	if matched, consumeErr := p.tryConsumeKeywords(KeywordPopulate); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		if createMaterializedView.Destination != nil {
			return nil, fmt.Errorf("POPULATE is only allowed when using ENGINE, not with TO clause")
		}
		if createMaterializedView.Engine == nil {
			return nil, fmt.Errorf("POPULATE requires ENGINE to be specified")
		}
		createMaterializedView.Populate = true
		createMaterializedView.StatementEnd = p.Pos()
	}

	// COMMENT can appear either before or after AS SELECT.
	// ClickHouse 26.2+ outputs COMMENT before AS SELECT in SHOW CREATE TABLE,
	// while 25.8 outputs it after AS SELECT.
	comment, err := p.tryParseComment()
	if err != nil {
		return nil, err
	}
	createMaterializedView.Comment = comment

	if matched, consumeErr := p.tryConsumeKeywords(KeywordAs); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		subQuery, err := p.parseSubQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.SubQuery = subQuery
		createMaterializedView.StatementEnd = subQuery.End()
	}

	// Also try parsing COMMENT after AS SELECT (ClickHouse 25.x format)
	if createMaterializedView.Comment == nil {
		comment, err = p.tryParseComment()
		if err != nil {
			return nil, err
		}
		createMaterializedView.Comment = comment
	}
	return createMaterializedView, nil
}

func (p *Parser) tryParseRefreshExpr(pos Pos) (*RefreshExpr, error) {
	if matched, consumeErr := p.tryConsumeKeywords(KeywordRefresh); consumeErr != nil {
		return nil, consumeErr
	} else if !matched {
		return nil, nil // nolint
	}

	// REFRESH EVERY|AFTER interval
	refreshExpr := &RefreshExpr{RefreshPos: pos}
	if !p.matchOneOfKeywords(KeywordEvery, KeywordAfter) {
		return nil, fmt.Errorf("expected EVERY or AFTER, but got %q", p.currentTokenKind())
	}
	refreshExpr.Frequency = p.current().String
	if err := p.lexer.consumeToken(); err != nil {
		return nil, err
	}

	interval, err := p.parseInterval(false)
	if err != nil {
		return nil, err
	}
	refreshExpr.Interval = interval

	// [OFFSET interval]
	if matched, consumeErr := p.tryConsumeKeywords(KeywordOffset); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		offset, err := p.parseInterval(false)
		if err != nil {
			return nil, err
		}
		refreshExpr.Offset = offset
	}

	return refreshExpr, nil
}

// (ATTACH | CREATE) (OR REPLACE)? VIEW (IF NOT EXISTS)? tableIdentifier uuidClause? clusterClause? tableSchemaClause? subqueryClause
func (p *Parser) parseCreateView(pos Pos, orReplace bool) (*CreateView, error) {
	createView := &CreateView{CreatePos: pos, OrReplace: orReplace}
	if err := p.expectKeyword(KeywordView); err != nil {
		return nil, err
	}

	var err error
	createView.IfNotExists, err = p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}

	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	createView.Name = tableIdentifier

	uuid, err := p.tryParseUUID()
	if err != nil {
		return nil, err
	}
	createView.UUID = uuid

	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	createView.OnCluster = onCluster

	if p.matchTokenKind(TokenKindLParen) {
		tableSchema, err := p.parseTableSchemaClause(p.Pos())
		if err != nil {
			return nil, err
		}
		createView.TableSchema = tableSchema
	}

	// parse COMMENT clause if exists (before AS SELECT)
	comment, err := p.tryParseComment()
	if err != nil {
		return nil, err
	}
	createView.Comment = comment

	if matched, consumeErr := p.tryConsumeKeywords(KeywordAs); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		subQuery, err := p.parseSubQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		createView.SubQuery = subQuery
		createView.StatementEnd = subQuery.End()
	}

	return createView, nil
}

// # CreateLiveViewStmt
// (ATTACH | CREATE) LIVE VIEW (IF NOT EXISTS)? tableIdentifier uuidClause?
// clusterClause? (WITH TIMEOUT DECIMAL_LITERAL?)? destinationClause? tableSchemaClause? subqueryClause
func (p *Parser) parseCreateLiveView(pos Pos) (*CreateLiveView, error) {
	if err := p.expectKeyword(KeywordLive); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(KeywordView); err != nil {
		return nil, err
	}

	createLiveView := &CreateLiveView{CreatePos: pos}
	// parse IF NOT EXISTS clause if exists
	var err error
	createLiveView.IfNotExists, err = p.tryParseIfNotExists()
	if err != nil {
		return nil, err
	}

	tableIdentifier, err := p.parseTableIdentifier(p.Pos())
	if err != nil {
		return nil, err
	}
	createLiveView.Name = tableIdentifier

	// try parse UUID clause if exists
	uuid, err := p.tryParseUUID()
	if err != nil {
		return nil, err
	}
	createLiveView.UUID = uuid
	// parse ON CLUSTER clause if exists
	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	createLiveView.OnCluster = onCluster

	withTimeout, err := p.tryParseWithTimeout(p.Pos())
	if err != nil {
		return nil, err
	}
	createLiveView.WithTimeout = withTimeout

	if p.matchKeyword(KeywordTo) {
		destination, err := p.parseDestinationClause(p.Pos())
		if err != nil {
			return nil, err
		}
		createLiveView.Destination = destination
	}

	if p.matchTokenKind(TokenKindLParen) {
		tableSchema, err := p.parseTableSchemaClause(p.Pos())
		if err != nil {
			return nil, err
		}
		createLiveView.TableSchema = tableSchema
	}

	if matched, consumeErr := p.tryConsumeKeywords(KeywordAs); consumeErr != nil {
		return nil, consumeErr
	} else if matched {
		subQuery, err := p.parseSubQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		createLiveView.SubQuery = subQuery
		createLiveView.StatementEnd = subQuery.End()
	}

	return createLiveView, nil
}

func (p *Parser) tryParseWithTimeout(pos Pos) (*WithTimeoutClause, error) {
	if matched, consumeErr := p.tryConsumeKeywords(KeywordWith); consumeErr != nil {
		return nil, consumeErr
	} else if !matched {
		return nil, nil // nolint
	}
	if err := p.expectKeyword(KeywordTimeout); err != nil {
		return nil, err
	}

	withTimeout := &WithTimeoutClause{WithTimeoutPos: pos}

	if p.matchTokenKind(TokenKindInt) {
		decimalNumber, err := p.parseDecimal(p.Pos())
		if err != nil {
			return nil, err
		}
		withTimeout.Number = decimalNumber
	}

	return withTimeout, nil
}
