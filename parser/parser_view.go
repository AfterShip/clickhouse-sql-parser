package parser

import "fmt"

func (p *Parser) parseCreateMaterializedView(pos Pos) (*CreateMaterializedView, error) {
	if err := p.consumeKeyword(KeywordMaterialized); err != nil {
		return nil, err
	}
	if err := p.consumeKeyword(KeywordView); err != nil {
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

	// try parse UUID clause if exists
	uuid, err := p.tryParseUUID()
	if err != nil {
		return nil, err
	}
	createMaterializedView.UUID = uuid
	// parse ON CLUSTER clause if exists
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	createMaterializedView.OnCluster = onCluster

	tableSchema, err := p.parseTableSchemaExpr(p.Pos())
	if err != nil {
		return nil, err
	}
	createMaterializedView.TableSchema = tableSchema

	switch {
	case p.matchKeyword(KeywordTo):
		destinationExpr, err := p.parseDestinationExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		createMaterializedView.Destination = destinationExpr
		createMaterializedView.StatementEnd = destinationExpr.End()
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
		if p.matchKeyword(KeywordAs) {
			subQuery, err := p.parseSubQuery(p.Pos())
			if err != nil {
				return nil, err
			}
			createMaterializedView.SubQuery = subQuery
			createMaterializedView.StatementEnd = subQuery.End()
		}
	default:
		return nil, fmt.Errorf("unexpected token: %q, expected TO or ENGINE", p.lastTokenKind())
	}
	return createMaterializedView, nil
}

// (ATTACH | CREATE) (OR REPLACE)? VIEW (IF NOT EXISTS)? tableIdentifier uuidClause? clusterClause? tableSchemaClause? subqueryClause
func (p *Parser) parseCreateView(pos Pos) (*CreateView, error) {
	if err := p.consumeKeyword(KeywordView); err != nil {
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

	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	createView.OnCluster = onCluster

	if p.matchTokenKind("(") {
		tableSchema, err := p.parseTableSchemaExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		createView.TableSchema = tableSchema
	}

	subQueryExpr, err := p.parseSubQuery(p.Pos())
	if err != nil {
		return nil, err
	}
	createView.SubQuery = subQueryExpr
	createView.StatementEnd = subQueryExpr.End()

	return createView, nil
}

// # CreateLiveViewStmt
// (ATTACH | CREATE) LIVE VIEW (IF NOT EXISTS)? tableIdentifier uuidClause?
// clusterClause? (WITH TIMEOUT DECIMAL_LITERAL?)? destinationClause? tableSchemaClause? subqueryClause
func (p *Parser) parseCreateLiveView(pos Pos) (*CreateLiveView, error) {
	if err := p.consumeKeyword(KeywordLive); err != nil {
		return nil, err
	}

	if err := p.consumeKeyword(KeywordView); err != nil {
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
	onCluster, err := p.tryParseOnCluster(p.Pos())
	if err != nil {
		return nil, err
	}
	createLiveView.OnCluster = onCluster

	withTimeExpr, err := p.tryParseWithTimeout(p.Pos())
	if err != nil {
		return nil, err
	}
	createLiveView.WithTimeout = withTimeExpr

	if p.matchKeyword(KeywordTo) {
		destinationExpr, err := p.parseDestinationExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		createLiveView.Destination = destinationExpr
	}

	if p.matchTokenKind("(") {
		tableSchema, err := p.parseTableSchemaExpr(p.Pos())
		if err != nil {
			return nil, err
		}
		createLiveView.TableSchema = tableSchema
	}

	subQuery, err := p.parseSubQuery(p.Pos())
	if err != nil {
		return nil, err
	}
	createLiveView.SubQuery = subQuery
	createLiveView.StatementEnd = subQuery.End()

	return createLiveView, nil
}

func (p *Parser) tryParseWithTimeout(pos Pos) (*WithTimeoutExpr, error) {
	if p.tryConsumeKeyword(KeywordWith) == nil {
		return nil, nil // nolint
	}
	if err := p.consumeKeyword(KeywordTimeout); err != nil {
		return nil, err
	}

	withTimeoutExpr := &WithTimeoutExpr{WithTimeoutPos: pos}

	if p.matchTokenKind(TokenInt) {
		decimalNumber, err := p.parseDecimal(p.Pos())
		if err != nil {
			return nil, err
		}
		withTimeoutExpr.Number = decimalNumber
	}

	return withTimeoutExpr, nil
}
