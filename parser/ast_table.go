package parser

type ProjectionOrderByClause struct {
	OrderByPos Pos
	Columns    *ColumnExprList
}

func (p *ProjectionOrderByClause) Pos() Pos {
	return p.OrderByPos
}

func (p *ProjectionOrderByClause) End() Pos {
	return p.Columns.End()
}

func (p *ProjectionOrderByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if p.Columns != nil {
		if err := p.Columns.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitProjectionOrderBy(p)
}

type ProjectionSelectStmt struct {
	LeftParenPos  Pos
	RightParenPos Pos
	With          *WithClause
	SelectColumns *ColumnExprList
	GroupBy       *GroupByClause
	OrderBy       *ProjectionOrderByClause
}

func (p *ProjectionSelectStmt) Pos() Pos {
	return p.LeftParenPos

}

func (p *ProjectionSelectStmt) End() Pos {
	return p.RightParenPos
}

func (p *ProjectionSelectStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if p.With != nil {
		if err := p.With.Accept(visitor); err != nil {
			return err
		}
	}
	if err := p.SelectColumns.Accept(visitor); err != nil {
		return err
	}
	if p.GroupBy != nil {
		if err := p.GroupBy.Accept(visitor); err != nil {
			return err
		}
	}
	if p.OrderBy != nil {
		if err := p.OrderBy.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitProjectionSelect(p)
}

type TableProjection struct {
	IncludeProjectionKeyword bool
	ProjectionPos            Pos
	Identifier               *NestedIdentifier
	Select                   *ProjectionSelectStmt
}

func (t *TableProjection) Pos() Pos {
	return t.ProjectionPos
}

func (t *TableProjection) End() Pos {
	return t.Select.End()
}

func (t *TableProjection) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Identifier.Accept(visitor); err != nil {
		return err
	}
	if err := t.Select.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTableProjection(t)
}

type TableIndex struct {
	IndexPos Pos

	Name        *NestedIdentifier
	ColumnExpr  *ColumnExpr
	ColumnType  Expr
	Granularity *NumberLiteral
}

func (a *TableIndex) Pos() Pos {
	return a.IndexPos
}

func (a *TableIndex) End() Pos {
	return a.Granularity.End()
}

func (a *TableIndex) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Name.Accept(visitor); err != nil {
		return err
	}
	if err := a.ColumnExpr.Accept(visitor); err != nil {
		return err
	}
	if err := a.ColumnType.Accept(visitor); err != nil {
		return err
	}
	if err := a.Granularity.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTableIndex(a)
}

type CreateDatabase struct {
	CreatePos    Pos // position of CREATE keyword
	StatementEnd Pos
	Name         Expr
	IfNotExists  bool // true if 'IF NOT EXISTS' is specified
	OnCluster    *ClusterClause
	Engine       *EngineExpr
	Comment      *StringLiteral
}

func (c *CreateDatabase) Pos() Pos {
	return c.CreatePos
}

func (c *CreateDatabase) End() Pos {
	return c.StatementEnd
}

func (c *CreateDatabase) Type() string {
	return "DATABASE"
}

func (c *CreateDatabase) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if c.Name != nil {
		if err := c.Name.Accept(visitor); err != nil {
			return err
		}
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Engine != nil {
		if err := c.Engine.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Comment != nil {
		if err := c.Comment.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateDatabase(c)
}

type CreateTable struct {
	CreatePos     Pos // position of CREATE|ATTACH keyword
	StatementEnd  Pos
	OrReplace     bool
	Name          *TableIdentifier
	IfNotExists   bool
	UUID          *UUID
	OnCluster     *ClusterClause
	TableSchema   *TableSchemaClause
	Engine        *EngineExpr
	SubQuery      *SubQuery
	TableFunction *TableFunctionExpr
	HasTemporary  bool
	Comment       *StringLiteral
}

func (c *CreateTable) Pos() Pos {
	return c.CreatePos
}

func (c *CreateTable) End() Pos {
	return c.StatementEnd
}

func (c *CreateTable) Type() string {
	return "CREATE TABLE"
}

func (c *CreateTable) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	if c.UUID != nil {
		if err := c.UUID.Accept(visitor); err != nil {
			return err
		}
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TableSchema != nil {
		if err := c.TableSchema.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Engine != nil {
		if err := c.Engine.Accept(visitor); err != nil {
			return err
		}
	}
	if c.SubQuery != nil {
		if err := c.SubQuery.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TableFunction != nil {
		if err := c.TableFunction.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Comment != nil {
		if err := c.Comment.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateTable(c)
}

type DestinationClause struct {
	ToPos           Pos
	TableIdentifier *TableIdentifier
	TableSchema     *TableSchemaClause
}

func (d *DestinationClause) Pos() Pos {
	return d.ToPos
}

func (d *DestinationClause) End() Pos {
	if d.TableSchema != nil {
		return d.TableSchema.End()
	}
	return d.TableIdentifier.End()
}

func (d *DestinationClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.TableIdentifier.Accept(visitor); err != nil {
		return err
	}
	if d.TableSchema != nil {
		if err := d.TableSchema.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDestinationExpr(d)
}

type ConstraintClause struct {
	ConstraintPos Pos
	Constraint    *Ident
	Expr          Expr
}

func (c *ConstraintClause) Pos() Pos {
	return c.ConstraintPos
}

func (c *ConstraintClause) End() Pos {
	return c.Expr.End()
}

func (c *ConstraintClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Constraint.Accept(visitor); err != nil {
		return err
	}
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitConstraintExpr(c)
}

type TableSchemaClause struct {
	SchemaPos     Pos
	SchemaEnd     Pos
	Columns       []Expr
	AliasTable    *TableIdentifier
	TableFunction *TableFunctionExpr
}

func (t *TableSchemaClause) Pos() Pos {
	return t.SchemaPos
}

func (t *TableSchemaClause) End() Pos {
	return t.SchemaEnd
}

func (t *TableSchemaClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	for _, column := range t.Columns {
		if err := column.Accept(visitor); err != nil {
			return err
		}
	}
	if t.AliasTable != nil {
		if err := t.AliasTable.Accept(visitor); err != nil {
			return err
		}
	}
	if t.TableFunction != nil {
		if err := t.TableFunction.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTableSchemaExpr(t)
}

type ClusterClause struct {
	OnPos Pos
	Expr  Expr
}

func (o *ClusterClause) Pos() Pos {
	return o.OnPos
}

func (o *ClusterClause) End() Pos {
	return o.Expr.End()
}

func (o *ClusterClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitOnClusterExpr(o)
}

type PartitionClause struct {
	PartitionPos Pos
	Expr         Expr
	ID           *StringLiteral
	All          bool
}

func (p *PartitionClause) Pos() Pos {
	return p.PartitionPos
}

func (p *PartitionClause) End() Pos {
	if p.ID != nil {
		return p.ID.LiteralEnd
	}
	return p.Expr.End()
}

func (p *PartitionClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if p.Expr != nil {
		if err := p.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	if p.ID != nil {
		if err := p.ID.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPartitionExpr(p)
}

type PartitionByClause struct {
	PartitionPos Pos
	Expr         Expr
}

func (p *PartitionByClause) Pos() Pos {
	return p.PartitionPos
}

func (p *PartitionByClause) End() Pos {
	return p.Expr.End()
}

func (p *PartitionByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if err := p.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPartitionByExpr(p)
}

type PrimaryKeyClause struct {
	PrimaryPos Pos
	Expr       Expr
}

func (p *PrimaryKeyClause) Pos() Pos {
	return p.PrimaryPos
}

func (p *PrimaryKeyClause) End() Pos {
	return p.Expr.End()
}

func (p *PrimaryKeyClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if err := p.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPrimaryKeyExpr(p)
}

type SampleByClause struct {
	SamplePos Pos
	Expr      Expr
}

func (s *SampleByClause) Pos() Pos {
	return s.SamplePos
}

func (s *SampleByClause) End() Pos {
	return s.Expr.End()
}

func (s *SampleByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSampleByExpr(s)
}

type TTLPolicyRuleAction struct {
	ActionPos Pos
	ActionEnd Pos
	Action    string
	Codec     *CompressionCodec
}

func (t *TTLPolicyRuleAction) Pos() Pos {
	return t.ActionPos
}

func (t *TTLPolicyRuleAction) End() Pos {
	if t.Codec != nil {
		return t.Codec.End()
	}
	return t.ActionEnd
}

func (t *TTLPolicyRuleAction) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if t.Codec != nil {
		if err := t.Codec.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTTLPolicyItemAction(t)
}

type TTLPolicyRule struct {
	RulePos  Pos
	ToVolume *StringLiteral
	ToDisk   *StringLiteral
	Action   *TTLPolicyRuleAction
	GroupBy  *GroupByClause
	Set      []*UpdateAssignment
}

func (t *TTLPolicyRule) Pos() Pos {
	return t.RulePos
}

func (t *TTLPolicyRule) End() Pos {
	if len(t.Set) > 0 {
		return t.Set[len(t.Set)-1].End()
	}
	if t.GroupBy != nil {
		return t.GroupBy.End()
	}
	if t.Action != nil {
		return t.Action.End()
	}
	if t.ToDisk != nil {
		return t.ToDisk.LiteralEnd
	}
	return t.ToVolume.LiteralEnd
}

func (t *TTLPolicyRule) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if t.ToVolume != nil {
		if err := t.ToVolume.Accept(visitor); err != nil {
			return err
		}
	}
	if t.ToDisk != nil {
		if err := t.ToDisk.Accept(visitor); err != nil {
			return err
		}
	}
	if t.Action != nil {
		if err := t.Action.Accept(visitor); err != nil {
			return err
		}
	}
	if t.GroupBy != nil {
		if err := t.GroupBy.Accept(visitor); err != nil {
			return err
		}
	}
	for _, set := range t.Set {
		if err := set.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTTLPolicyRule(t)
}

type TTLPolicy struct {
	Item  *TTLPolicyRule
	Where *WhereClause
}

func (t *TTLPolicy) Pos() Pos {
	if t.Item != nil {
		return t.Item.Pos()
	}
	return t.Where.Pos()
}

func (t *TTLPolicy) End() Pos {
	if t.Where != nil {
		return t.Where.End()
	}
	return t.Item.End()
}

func (t *TTLPolicy) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if t.Item != nil {
		if err := t.Item.Accept(visitor); err != nil {
			return err
		}
	}
	if t.Where != nil {
		if err := t.Where.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTTLPolicy(t)
}

type TTLExpr struct {
	TTLPos Pos
	Expr   Expr
	Policy *TTLPolicy
}

func (t *TTLExpr) Pos() Pos {
	return t.TTLPos
}

func (t *TTLExpr) End() Pos {
	if t.Policy != nil {
		return t.Policy.End()
	}
	return t.Expr.End()
}

func (t *TTLExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Expr.Accept(visitor); err != nil {
		return err
	}
	if t.Policy != nil {
		if err := t.Policy.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTTLExpr(t)
}

type TTLClause struct {
	TTLPos  Pos
	ListEnd Pos
	Items   []*TTLExpr
}

func (t *TTLClause) Pos() Pos {
	return t.TTLPos
}

func (t *TTLClause) End() Pos {
	return t.ListEnd
}

func (t *TTLClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	for _, item := range t.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTTLExprList(t)
}

type ColumnDef struct {
	NamePos   Pos
	ColumnEnd Pos
	Name      *NestedIdentifier
	Type      ColumnType
	NotNull   *NotNullLiteral
	Nullable  *NullLiteral

	DefaultExpr      Expr
	MaterializedExpr Expr
	AliasExpr        Expr

	Codec *CompressionCodec
	TTL   *TTLClause

	Comment          *StringLiteral
	CompressionCodec *Ident
}

func (c *ColumnDef) Pos() Pos {
	return c.Name.Pos()
}

func (c *ColumnDef) End() Pos {
	return c.ColumnEnd
}

func (c *ColumnDef) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	if c.Type != nil {
		if err := c.Type.Accept(visitor); err != nil {
			return err
		}
	}
	if c.NotNull != nil {
		if err := c.NotNull.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Nullable != nil {
		if err := c.Nullable.Accept(visitor); err != nil {
			return err
		}
	}
	if c.DefaultExpr != nil {
		if err := c.DefaultExpr.Accept(visitor); err != nil {
			return err
		}
	}
	if c.MaterializedExpr != nil {
		if err := c.MaterializedExpr.Accept(visitor); err != nil {
			return err
		}
	}
	if c.AliasExpr != nil {
		if err := c.AliasExpr.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Codec != nil {
		if err := c.Codec.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TTL != nil {
		if err := c.TTL.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Comment != nil {
		if err := c.Comment.Accept(visitor); err != nil {
			return err
		}
	}
	if c.CompressionCodec != nil {
		if err := c.CompressionCodec.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnDef(c)
}

type CompressionCodec struct {
	CodecPos      Pos
	RightParenPos Pos
	Type          *Ident
	TypeLevel     *NumberLiteral
	Name          *Ident
	Level         *NumberLiteral // compression level
}

func (c *CompressionCodec) Pos() Pos {
	return c.CodecPos
}

func (c *CompressionCodec) End() Pos {
	return c.RightParenPos
}

func (c *CompressionCodec) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if c.Type != nil {
		if err := c.Type.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TypeLevel != nil {
		if err := c.TypeLevel.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Name != nil {
		if err := c.Name.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Level != nil {
		if err := c.Level.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCompressionCodec(c)
}

// TODO(@git-hulk): split into EngineClause and EngineExpr
type EngineExpr struct {
	EnginePos   Pos
	EngineEnd   Pos
	Name        string
	Params      *ParamExprList
	PrimaryKey  *PrimaryKeyClause
	PartitionBy *PartitionByClause
	SampleBy    *SampleByClause
	TTL         *TTLClause
	Settings    *SettingsClause
	OrderBy     *OrderByClause
}

func (e *EngineExpr) Pos() Pos {
	return e.EnginePos
}

func (e *EngineExpr) End() Pos {
	return e.EngineEnd
}

func (e *EngineExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	if e.Params != nil {
		if err := e.Params.Accept(visitor); err != nil {
			return err
		}
	}
	if e.PrimaryKey != nil {
		if err := e.PrimaryKey.Accept(visitor); err != nil {
			return err
		}
	}
	if e.PartitionBy != nil {
		if err := e.PartitionBy.Accept(visitor); err != nil {
			return err
		}
	}
	if e.SampleBy != nil {
		if err := e.SampleBy.Accept(visitor); err != nil {
			return err
		}
	}
	if e.TTL != nil {
		if err := e.TTL.Accept(visitor); err != nil {
			return err
		}
	}
	if e.Settings != nil {
		if err := e.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	if e.OrderBy != nil {
		if err := e.OrderBy.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitEngineExpr(e)
}

type DropDatabase struct {
	DropPos      Pos
	StatementEnd Pos
	Name         *Ident
	IfExists     bool
	OnCluster    *ClusterClause
	Modifier     string
}

func (d *DropDatabase) Pos() Pos {
	return d.DropPos
}

func (d *DropDatabase) End() Pos {
	return d.StatementEnd
}

func (d *DropDatabase) Type() string {
	return "DATABASE"
}

func (d *DropDatabase) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Name.Accept(visitor); err != nil {
		return err
	}
	if d.OnCluster != nil {
		if err := d.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDropDatabase(d)
}

type DropStmt struct {
	DropPos      Pos
	StatementEnd Pos

	DropTarget  string
	Name        *TableIdentifier
	IfExists    bool
	OnCluster   *ClusterClause
	IsTemporary bool
	Modifier    string
}

func (d *DropStmt) Pos() Pos {
	return d.DropPos
}

func (d *DropStmt) End() Pos {
	return d.StatementEnd
}

func (d *DropStmt) Type() string {
	return "DROP " + d.DropTarget
}

func (d *DropStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Name.Accept(visitor); err != nil {
		return err
	}
	if d.OnCluster != nil {
		if err := d.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDropStmt(d)

}

type TruncateTable struct {
	TruncatePos  Pos
	StatementEnd Pos
	IsTemporary  bool
	IfExists     bool
	Name         *TableIdentifier
	OnCluster    *ClusterClause
}

func (t *TruncateTable) Pos() Pos {
	return t.TruncatePos
}

func (t *TruncateTable) End() Pos {
	return t.StatementEnd
}

func (t *TruncateTable) Type() string {
	return "TRUNCATE TABLE"
}

func (t *TruncateTable) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Name.Accept(visitor); err != nil {
		return err
	}
	if t.OnCluster != nil {
		if err := t.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTruncateTable(t)
}

type RenameStmt struct {
	RenamePos    Pos
	StatementEnd Pos

	RenameTarget   string
	TargetPairList []*TargetPair
	OnCluster      *ClusterClause
}

func (r *RenameStmt) Pos() Pos {
	return r.RenamePos
}

func (r *RenameStmt) End() Pos {
	return r.StatementEnd
}

func (r *RenameStmt) Type() string {
	return "RENAME " + r.RenameTarget
}

func (r *RenameStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	for _, pair := range r.TargetPairList {
		if err := pair.Old.Accept(visitor); err != nil {
			return err
		}
		if err := pair.New.Accept(visitor); err != nil {
			return err
		}
	}
	if r.OnCluster != nil {
		if err := r.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRenameStmt(r)
}

type TargetPair struct {
	Old *TableIdentifier
	New *TableIdentifier
}

func (t *TargetPair) Pos() Pos {
	return t.Old.Pos()
}

func (t *TargetPair) End() Pos {
	return t.New.End()
}

func (t *TargetPair) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Old.Accept(visitor); err != nil {
		return err
	}
	if err := t.New.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTargetPairExpr(t)
}
