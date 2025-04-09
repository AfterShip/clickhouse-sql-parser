package parser

import "fmt"

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

	// parse ON CLUSTER clause if exists
	onCluster, err := p.tryParseClusterClause(p.Pos())
	if err != nil {
		return nil, err
	}
	createMaterializedView.OnCluster = onCluster

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
		if populate := p.tryConsumeKeyword(KeywordPopulate); populate != nil {
			createMaterializedView.Populate = true
			createMaterializedView.StatementEnd = populate.End
		}
	default:
		return nil, fmt.Errorf("unexpected token: %q, expected TO or ENGINE", p.lastTokenKind())
	}
	if p.tryConsumeKeyword(KeywordAs) != nil {
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

// (ATTACH | CREATE) (OR REPLACE)? VIEW (IF NOT EXISTS)? tableIdentifier uuidClause? clusterClause? tableSchemaClause? subqueryClause
func (p *Parser) parseCreateView(pos Pos) (*CreateView, error) {
	if err := p.expectKeyword(KeywordView); err != nil {
		return nil, err
	}

	createView := &CreateView{CreatePos: pos}
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

	if p.tryConsumeKeyword(KeywordAs) != nil {
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

	if p.tryConsumeKeyword(KeywordAs) != nil {
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
	if p.tryConsumeKeyword(KeywordWith) == nil {
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
