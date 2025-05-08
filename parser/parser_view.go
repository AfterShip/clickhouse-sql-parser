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
func (p *Parser) parseCreateMaterializedView(pos Pos) (*CreateMaterializedView, error) {
	if err := p.expectKeyword(KeywordMaterialized); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(KeywordView); err != nil {
		return nil, err
	}

	createMaterializedView := &CreateMaterializedView{CreatePos: pos}

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

	if p.tryConsumeKeywords(KeywordRandomize, KeywordFor) {
		randomizeFor, err := p.parseInterval(false)
		if err != nil {
			return nil, err
		}
		createMaterializedView.RandomizeFor = randomizeFor
	}
	if p.tryConsumeKeywords(KeywordDepends, KeywordOn) {
		dependsOnTables := make([]*TableIdentifier, 0)
		table, err := p.parseTableIdentifier(p.Pos())
		if err != nil {
			return nil, err
		}
		dependsOnTables = append(dependsOnTables, table)
		for p.matchTokenKind(TokenKindComma) {
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
	createMaterializedView.HasAppend = p.tryConsumeKeywords(KeywordAppend)

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
	case p.matchKeyword(KeywordEngine):
		engineExpr, err := p.parseEngineExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.Engine = engineExpr
		createMaterializedView.StatementEnd = engineExpr.End()
	default:
		return nil, fmt.Errorf("unexpected token: %q, expected TO or ENGINE", p.lastTokenKind())
	}
	createMaterializedView.HasEmpty = p.tryConsumeKeywords(KeywordEmpty)

	// Parse DEFINER clause
	if p.tryConsumeKeywords(KeywordDefiner) {
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
	if p.tryConsumeKeywords(KeywordSQL, KeywordSecurity) {
		if !p.matchOneOfKeywords(KeywordDefiner, KeywordNone) {
			return nil, fmt.Errorf("expected DEFINER or NONE after SQL SECURITY, got %q", p.lastTokenKind())
		}
		createMaterializedView.SQLSecurity = p.last().String
		_ = p.lexer.consumeToken()
	}

	// Check for POPULATE before AS SELECT - only valid with ENGINE and no Destination
	if p.tryConsumeKeywords(KeywordPopulate) {
		if createMaterializedView.Destination != nil {
			return nil, fmt.Errorf("POPULATE is only allowed when using ENGINE, not with TO clause")
		}
		if createMaterializedView.Engine == nil {
			return nil, fmt.Errorf("POPULATE requires ENGINE to be specified")
		}
		createMaterializedView.Populate = true
		createMaterializedView.StatementEnd = p.Pos()
	}

	if p.tryConsumeKeywords(KeywordAs) {
		subQuery, err := p.parseSubQuery(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.SubQuery = subQuery
		createMaterializedView.StatementEnd = subQuery.End()
	}

	comment, err := p.tryParseComment()
	if err != nil {
		return nil, err
	}
	createMaterializedView.Comment = comment
	return createMaterializedView, nil
}

func (p *Parser) tryParseRefreshExpr(pos Pos) (*RefreshExpr, error) {
	if !p.tryConsumeKeywords(KeywordRefresh) {
		return nil, nil // nolint
	}

	// REFRESH EVERY|AFTER interval
	refreshExpr := &RefreshExpr{RefreshPos: pos}
	if !p.matchOneOfKeywords(KeywordEvery, KeywordAfter) {
		return nil, fmt.Errorf("expected EVERY or AFTER, but got %q", p.lastTokenKind())
	}
	refreshExpr.Frequency = p.last().String
	_ = p.lexer.consumeToken()

	interval, err := p.parseInterval(false)
	if err != nil {
		return nil, err
	}
	refreshExpr.Interval = interval

	// [OFFSET interval]
	if p.tryConsumeKeywords(KeywordOffset) {
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

	if p.tryConsumeKeywords(KeywordAs) {
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

	if p.tryConsumeKeywords(KeywordAs) {
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
	if !p.tryConsumeKeywords(KeywordWith) {
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
