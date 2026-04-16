package parser

type OrderDirection string

const (
	OrderDirectionNone OrderDirection = ""
	OrderDirectionAsc  OrderDirection = "ASC"
	OrderDirectionDesc OrderDirection = "DESC"
)

type Expr interface {
	Pos() Pos
	End() Pos
	FormatSQL(*Formatter)
	Accept(visitor ASTVisitor) error
}

type DDL interface {
	Expr
	Type() string
}

type SelectItem struct {
	Expr Expr
	// Please refer: https://clickhouse.com/docs/en/sql-reference/statements/select#select-modifiers
	Modifiers []*FunctionExpr
	Alias     *Ident
}

func (s *SelectItem) Pos() Pos {
	return s.Expr.Pos()
}

func (s *SelectItem) End() Pos {
	if s.Alias != nil {
		return s.Alias.End()
	}
	if len(s.Modifiers) > 0 {
		return s.Modifiers[len(s.Modifiers)-1].End()
	}
	return s.Expr.End()
}

func (s *SelectItem) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Expr.Accept(visitor); err != nil {
		return err
	}
	for _, modifier := range s.Modifiers {
		if err := modifier.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Alias != nil {
		if err := s.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSelectItem(s)
}

type OperationExpr struct {
	OperationPos Pos
	Kind         TokenKind
}

func (o *OperationExpr) Pos() Pos {
	return o.OperationPos
}

func (o *OperationExpr) End() Pos {
	return o.OperationPos + Pos(len(o.Kind))
}

func (o *OperationExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	return visitor.VisitOperationExpr(o)
}

type TernaryOperation struct {
	Condition Expr
	TrueExpr  Expr
	FalseExpr Expr
}

func (t *TernaryOperation) Pos() Pos {
	return t.Condition.Pos()
}

func (t *TernaryOperation) End() Pos {
	return t.FalseExpr.End()
}

func (t *TernaryOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.TrueExpr.Accept(visitor); err != nil {
		return err
	}
	if err := t.FalseExpr.Accept(visitor); err != nil {
		return err
	}
	if err := t.Condition.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTernaryExpr(t)
}

type BinaryOperation struct {
	LeftExpr  Expr
	Operation TokenKind
	RightExpr Expr
	HasGlobal bool
	HasNot    bool
}

func (p *BinaryOperation) Pos() Pos {
	return p.LeftExpr.Pos()
}

func (p *BinaryOperation) End() Pos {
	return p.RightExpr.End()
}

func (p *BinaryOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if err := p.LeftExpr.Accept(visitor); err != nil {
		return err
	}
	if err := p.RightExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitBinaryExpr(p)
}

type IndexOperation struct {
	Object    Expr
	Operation TokenKind
	Index     Expr
}

func (i *IndexOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.Object.Accept(visitor); err != nil {
		return err
	}
	if err := i.Index.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIndexOperation(i)
}

func (i *IndexOperation) Pos() Pos {
	return i.Object.Pos()
}

func (i *IndexOperation) End() Pos {
	return i.Index.End()
}

type JoinTableExpr struct {
	Table        *TableExpr
	StatementEnd Pos
	SampleRatio  *SampleClause
	HasFinal     bool
}

func (j *JoinTableExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if err := j.Table.Accept(visitor); err != nil {
		return err
	}
	if j.SampleRatio != nil {
		return j.SampleRatio.Accept(visitor)
	}
	return visitor.VisitJoinTableExpr(j)
}

func (j *JoinTableExpr) Pos() Pos {
	return j.Table.Pos()
}

func (j *JoinTableExpr) End() Pos {
	return j.StatementEnd
}

type AlterTableClause interface {
	Expr
	AlterType() string
}

type AlterTable struct {
	AlterPos        Pos
	StatementEnd    Pos
	TableIdentifier *TableIdentifier
	OnCluster       *ClusterClause
	AlterExprs      []AlterTableClause
}

func (a *AlterTable) Pos() Pos {
	return a.AlterPos
}

func (a *AlterTable) End() Pos {
	return a.StatementEnd
}

func (a *AlterTable) Type() string {
	return "ALTER TABLE"
}

func (a *AlterTable) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.TableIdentifier.Accept(visitor); err != nil {
		return err
	}
	if a.OnCluster != nil {
		if err := a.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}

	for _, expr := range a.AlterExprs {
		if err := expr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTable(a)
}

type AlterTableAttachPartition struct {
	AttachPos Pos

	Partition *PartitionClause
	From      *TableIdentifier
}

func (a *AlterTableAttachPartition) Pos() Pos {
	return a.AttachPos
}

func (a *AlterTableAttachPartition) End() Pos {
	if a.From != nil {
		return a.From.End()
	}
	return a.Partition.End()
}

func (a *AlterTableAttachPartition) AlterType() string {
	return "ATTACH_PARTITION"
}

func (a *AlterTableAttachPartition) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Partition.Accept(visitor); err != nil {
		return err
	}
	if a.From != nil {
		if err := a.From.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableAttachPartition(a)
}

type AlterTableDetachPartition struct {
	DetachPos Pos
	Partition *PartitionClause
	Settings  *SettingsClause
}

func (a *AlterTableDetachPartition) Pos() Pos {
	return a.DetachPos
}

func (a *AlterTableDetachPartition) End() Pos {
	return a.Partition.End()
}

func (a *AlterTableDetachPartition) AlterType() string {
	return "DETACH_PARTITION"
}

func (a *AlterTableDetachPartition) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Partition.Accept(visitor); err != nil {
		return err
	}
	if a.Settings != nil {
		if err := a.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableDetachPartition(a)
}

type AlterTableDropPartition struct {
	DropPos     Pos
	HasDetached bool
	Partition   *PartitionClause
	Settings    *SettingsClause
}

func (a *AlterTableDropPartition) Pos() Pos {
	return a.DropPos
}

func (a *AlterTableDropPartition) End() Pos {
	if a.Settings != nil {
		a.Settings.End()
	}
	return a.Partition.End()
}

func (a *AlterTableDropPartition) AlterType() string {
	return "DROP_PARTITION"
}

func (a *AlterTableDropPartition) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Partition.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableDropPartition(a)
}

type AlterTableMaterializeProjection struct {
	MaterializedPos Pos
	StatementEnd    Pos
	IfExists        bool
	ProjectionName  *NestedIdentifier
	Partition       *PartitionClause
}

func (a *AlterTableMaterializeProjection) Pos() Pos {
	return a.MaterializedPos
}

func (a *AlterTableMaterializeProjection) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableMaterializeProjection) AlterType() string {
	return "MATERIALIZE_PROJECTION"
}

func (a *AlterTableMaterializeProjection) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.ProjectionName.Accept(visitor); err != nil {
		return err
	}
	if a.Partition != nil {
		if err := a.Partition.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableMaterializeProjection(a)
}

type AlterTableMaterializeIndex struct {
	MaterializedPos Pos
	StatementEnd    Pos
	IfExists        bool
	IndexName       *NestedIdentifier
	Partition       *PartitionClause
}

func (a *AlterTableMaterializeIndex) Pos() Pos {
	return a.MaterializedPos
}

func (a *AlterTableMaterializeIndex) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableMaterializeIndex) AlterType() string {
	return "MATERIALIZE_INDEX"
}

func (a *AlterTableMaterializeIndex) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.IndexName.Accept(visitor); err != nil {
		return err
	}
	if a.Partition != nil {
		if err := a.Partition.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableMaterializeIndex(a)
}

type AlterTableFreezePartition struct {
	FreezePos    Pos
	StatementEnd Pos
	Partition    *PartitionClause
}

func (a *AlterTableFreezePartition) Pos() Pos {
	return a.FreezePos
}

func (a *AlterTableFreezePartition) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableFreezePartition) AlterType() string {
	return "FREEZE_PARTITION"
}

func (a *AlterTableFreezePartition) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if a.Partition != nil {
		if err := a.Partition.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableFreezePartition(a)
}

type AlterTableAddColumn struct {
	AddPos       Pos
	StatementEnd Pos

	Column      *ColumnDef
	IfNotExists bool
	After       *NestedIdentifier
	Settings    *SettingsClause
}

func (a *AlterTableAddColumn) Pos() Pos {
	return a.AddPos
}

func (a *AlterTableAddColumn) End() Pos {
	if a.Settings != nil {
		return a.Settings.End()
	}
	return a.StatementEnd
}

func (a *AlterTableAddColumn) AlterType() string {
	return "ADD_COLUMN"
}

func (a *AlterTableAddColumn) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Column.Accept(visitor); err != nil {
		return err
	}
	if a.After != nil {
		if err := a.After.Accept(visitor); err != nil {
			return err
		}
	}
	if a.Settings != nil {
		if err := a.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableAddColumn(a)
}

type AlterTableAddIndex struct {
	AddPos       Pos
	StatementEnd Pos

	Index       *TableIndex
	IfNotExists bool
	After       *NestedIdentifier
}

func (a *AlterTableAddIndex) Pos() Pos {
	return a.AddPos
}

func (a *AlterTableAddIndex) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableAddIndex) AlterType() string {
	return "ADD_INDEX"
}

func (a *AlterTableAddIndex) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Index.Accept(visitor); err != nil {
		return err
	}
	if a.After != nil {
		if err := a.After.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableAddIndex(a)
}

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

type AlterTableAddProjection struct {
	AddPos       Pos
	StatementEnd Pos

	IfNotExists     bool
	TableProjection *TableProjection
	After           *NestedIdentifier
}

func (a *AlterTableAddProjection) Pos() Pos {
	return a.AddPos
}

func (a *AlterTableAddProjection) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableAddProjection) AlterType() string {
	return "ADD_PROJECTION"
}

func (a *AlterTableAddProjection) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.TableProjection.Accept(visitor); err != nil {
		return err
	}
	if a.After != nil {
		if err := a.After.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableAddProjection(a)
}

type AlterTableDropColumn struct {
	DropPos    Pos
	ColumnName *NestedIdentifier
	IfExists   bool
}

func (a *AlterTableDropColumn) Pos() Pos {
	return a.DropPos
}

func (a *AlterTableDropColumn) End() Pos {
	return a.ColumnName.End()
}

func (a *AlterTableDropColumn) AlterType() string {
	return "DROP_COLUMN"
}

func (a *AlterTableDropColumn) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.ColumnName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableDropColumn(a)
}

type AlterTableDropIndex struct {
	DropPos   Pos
	IndexName *NestedIdentifier
	IfExists  bool
}

func (a *AlterTableDropIndex) Pos() Pos {
	return a.DropPos
}

func (a *AlterTableDropIndex) End() Pos {
	return a.IndexName.End()
}

func (a *AlterTableDropIndex) AlterType() string {
	return "DROP_INDEX"
}

func (a *AlterTableDropIndex) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.IndexName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableDropIndex(a)
}

type AlterTableDropProjection struct {
	DropPos        Pos
	ProjectionName *NestedIdentifier
	IfExists       bool
}

func (a *AlterTableDropProjection) Pos() Pos {
	return a.DropPos
}

func (a *AlterTableDropProjection) End() Pos {
	return a.ProjectionName.End()
}

func (a *AlterTableDropProjection) AlterType() string {
	return "DROP_PROJECTION"
}

func (a *AlterTableDropProjection) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.ProjectionName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableDropProjection(a)
}

type AlterTableRemoveTTL struct {
	RemovePos    Pos
	StatementEnd Pos
}

func (a *AlterTableRemoveTTL) Pos() Pos {
	return a.RemovePos
}

func (a *AlterTableRemoveTTL) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableRemoveTTL) AlterType() string {
	return "REMOVE_TTL"
}

func (a *AlterTableRemoveTTL) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	return visitor.VisitAlterTableRemoveTTL(a)
}

type AlterTableClearColumn struct {
	ClearPos     Pos
	StatementEnd Pos

	IfExists      bool
	ColumnName    *NestedIdentifier
	PartitionExpr *PartitionClause
}

func (a *AlterTableClearColumn) Pos() Pos {
	return a.ClearPos
}

func (a *AlterTableClearColumn) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableClearColumn) AlterType() string {
	return "CLEAR_COLUMN"
}

func (a *AlterTableClearColumn) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.ColumnName.Accept(visitor); err != nil {
		return err
	}
	if a.PartitionExpr != nil {
		if err := a.PartitionExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableClearColumn(a)
}

type AlterTableClearIndex struct {
	ClearPos     Pos
	StatementEnd Pos

	IfExists      bool
	IndexName     *NestedIdentifier
	PartitionExpr *PartitionClause
}

func (a *AlterTableClearIndex) Pos() Pos {
	return a.ClearPos
}

func (a *AlterTableClearIndex) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableClearIndex) AlterType() string {
	return "CLEAR_INDEX"
}

func (a *AlterTableClearIndex) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.IndexName.Accept(visitor); err != nil {
		return err
	}
	if a.PartitionExpr != nil {
		if err := a.PartitionExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableClearIndex(a)
}

type AlterTableClearProjection struct {
	ClearPos     Pos
	StatementEnd Pos

	IfExists       bool
	ProjectionName *NestedIdentifier
	PartitionExpr  *PartitionClause
}

func (a *AlterTableClearProjection) Pos() Pos {
	return a.ClearPos
}

func (a *AlterTableClearProjection) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableClearProjection) AlterType() string {
	return "CLEAR_PROJECTION"
}

func (a *AlterTableClearProjection) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.ProjectionName.Accept(visitor); err != nil {
		return err
	}
	if a.PartitionExpr != nil {
		if err := a.PartitionExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableClearProjection(a)
}

type AlterTableRenameColumn struct {
	RenamePos Pos

	IfExists      bool
	OldColumnName *NestedIdentifier
	NewColumnName *NestedIdentifier
}

func (a *AlterTableRenameColumn) Pos() Pos {
	return a.RenamePos
}

func (a *AlterTableRenameColumn) End() Pos {
	return a.NewColumnName.End()
}

func (a *AlterTableRenameColumn) AlterType() string {
	return "RENAME_COLUMN"
}

func (a *AlterTableRenameColumn) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.OldColumnName.Accept(visitor); err != nil {
		return err
	}
	if err := a.NewColumnName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableRenameColumn(a)
}

type AlterTableModifyQuery struct {
	ModifyPos    Pos
	StatementEnd Pos
	SelectExpr   *SelectQuery
}

func (a *AlterTableModifyQuery) Pos() Pos {
	return a.ModifyPos
}

func (a *AlterTableModifyQuery) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableModifyQuery) AlterType() string {
	return "MODIFY_QUERY"
}

func (a *AlterTableModifyQuery) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.SelectExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableModifyQuery(a)
}

type AlterTableModifyTTL struct {
	ModifyPos    Pos
	StatementEnd Pos
	TTL          *TTLClause
}

func (a *AlterTableModifyTTL) Pos() Pos {
	return a.ModifyPos
}

func (a *AlterTableModifyTTL) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableModifyTTL) AlterType() string {
	return "MODIFY_TTL"
}

func (a *AlterTableModifyTTL) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.TTL.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableModifyTTL(a)
}

type AlterTableModifyColumn struct {
	ModifyPos    Pos
	StatementEnd Pos

	IfExists           bool
	Column             *ColumnDef
	RemovePropertyType *RemovePropertyType
}

func (a *AlterTableModifyColumn) Pos() Pos {
	return a.ModifyPos
}

func (a *AlterTableModifyColumn) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableModifyColumn) AlterType() string {
	return "MODIFY_COLUMN"
}

func (a *AlterTableModifyColumn) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Column.Accept(visitor); err != nil {
		return err
	}
	if a.RemovePropertyType != nil {
		if err := a.RemovePropertyType.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableModifyColumn(a)
}

type AlterTableModifySetting struct {
	ModifyPos    Pos
	StatementEnd Pos
	Settings     []*SettingExpr
}

func (a *AlterTableModifySetting) Pos() Pos {
	return a.ModifyPos
}

func (a *AlterTableModifySetting) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableModifySetting) AlterType() string {
	return "MODIFY_SETTING"
}

func (a *AlterTableModifySetting) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	for _, setting := range a.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableModifySetting(a)
}

type AlterTableResetSetting struct {
	ResetPos     Pos
	StatementEnd Pos
	Settings     []*Ident
}

func (a *AlterTableResetSetting) Pos() Pos {
	return a.ResetPos
}

func (a *AlterTableResetSetting) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableResetSetting) AlterType() string {
	return "RESET_SETTING"
}

func (a *AlterTableResetSetting) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	for _, setting := range a.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterTableResetSetting(a)
}

type AlterTableReplacePartition struct {
	ReplacePos Pos
	Partition  *PartitionClause
	Table      *TableIdentifier
}

func (a *AlterTableReplacePartition) Pos() Pos {
	return a.ReplacePos
}

func (a *AlterTableReplacePartition) End() Pos {
	return a.Table.End()
}

func (a *AlterTableReplacePartition) AlterType() string {
	return "REPLACE_PARTITION"
}

func (a *AlterTableReplacePartition) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Partition.Accept(visitor); err != nil {
		return err
	}
	if err := a.Table.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableReplacePartition(a)
}

type AlterTableDelete struct {
	DeletePos    Pos
	StatementEnd Pos
	WhereClause  Expr
}

func (a *AlterTableDelete) Pos() Pos {
	return a.DeletePos
}

func (a *AlterTableDelete) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableDelete) AlterType() string {
	return "DELETE"
}

func (a *AlterTableDelete) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.WhereClause.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableDelete(a)
}

type AlterTableUpdate struct {
	UpdatePos    Pos
	StatementEnd Pos
	Assignments  []*UpdateAssignment
	InPartition  *PartitionClause
	WhereClause  Expr
}

func (a *AlterTableUpdate) Pos() Pos {
	return a.UpdatePos
}

func (a *AlterTableUpdate) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableUpdate) AlterType() string {
	return "UPDATE"
}

func (a *AlterTableUpdate) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	for _, assignment := range a.Assignments {
		if err := assignment.Accept(visitor); err != nil {
			return err
		}
	}
	if a.InPartition != nil {
		if err := a.InPartition.Accept(visitor); err != nil {
			return err
		}
	}
	if err := a.WhereClause.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableUpdate(a)
}

type UpdateAssignment struct {
	AssignmentPos Pos
	Column        *NestedIdentifier
	Expr          Expr
}

func (u *UpdateAssignment) Pos() Pos {
	return u.AssignmentPos
}

func (u *UpdateAssignment) End() Pos {
	return u.Expr.End()
}

func (u *UpdateAssignment) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	if err := u.Column.Accept(visitor); err != nil {
		return err
	}
	if err := u.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUpdateAssignment(u)
}

type RemovePropertyType struct {
	RemovePos Pos

	PropertyType Expr
}

func (a *RemovePropertyType) Pos() Pos {
	return a.RemovePos
}

func (a *RemovePropertyType) End() Pos {
	return a.PropertyType.End()
}

func (a *RemovePropertyType) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.PropertyType.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitRemovePropertyType(a)
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

type Ident struct {
	Name      string
	QuoteType int
	NamePos   Pos
	NameEnd   Pos
}

func (i *Ident) Pos() Pos {
	return i.NamePos
}

func (i *Ident) End() Pos {
	return i.NameEnd
}

func (i *Ident) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	return visitor.VisitIdent(i)
}

type UUID struct {
	Value *StringLiteral
}

func (u *UUID) Pos() Pos {
	return u.Value.LiteralPos
}

func (u *UUID) End() Pos {
	return u.Value.LiteralEnd
}

func (u *UUID) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	return visitor.VisitUUID(u)
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
	return visitor.VisitCreateTable(c)
}

type CreateMaterializedView struct {
	CreatePos    Pos // position of CREATE|ATTACH keyword
	StatementEnd Pos
	Name         *TableIdentifier
	IfNotExists  bool
	OnCluster    *ClusterClause
	Refresh      *RefreshExpr
	RandomizeFor *IntervalExpr
	DependsOn    []*TableIdentifier
	Settings     *SettingsClause
	HasAppend    bool
	Engine       *EngineExpr
	HasEmpty     bool
	Destination  *DestinationClause
	SubQuery     *SubQuery
	Populate     bool
	Comment      *StringLiteral
	Definer      *Ident
	SQLSecurity  string
}

func (c *CreateMaterializedView) Pos() Pos {
	return c.CreatePos
}

func (c *CreateMaterializedView) End() Pos {
	return c.StatementEnd
}

func (c *CreateMaterializedView) Type() string {
	return "MATERIALIZED_VIEW"
}

func (c *CreateMaterializedView) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Refresh != nil {
		if err := c.Refresh.Accept(visitor); err != nil {
			return err
		}
	}
	if c.RandomizeFor != nil {
		if err := c.RandomizeFor.Accept(visitor); err != nil {
			return err
		}
	}
	if c.DependsOn != nil {
		for _, dep := range c.DependsOn {
			if err := dep.Accept(visitor); err != nil {
				return err
			}
		}
	}
	if c.Settings != nil {
		if err := c.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Engine != nil {
		if err := c.Engine.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Destination != nil {
		if err := c.Destination.Accept(visitor); err != nil {
			return err
		}
		if c.Destination.TableSchema != nil {
			if err := c.Destination.TableSchema.Accept(visitor); err != nil {
				return err
			}
		}
	}
	if c.SubQuery != nil {
		if err := c.SubQuery.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Definer != nil {
		if err := c.Definer.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Comment != nil {
		if err := c.Comment.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateMaterializedView(c)
}

type CreateView struct {
	CreatePos    Pos // position of CREATE|ATTACH keyword
	StatementEnd Pos
	OrReplace    bool
	Name         *TableIdentifier
	IfNotExists  bool
	UUID         *UUID
	OnCluster    *ClusterClause
	TableSchema  *TableSchemaClause
	SubQuery     *SubQuery
}

func (c *CreateView) Pos() Pos {
	return c.CreatePos
}

func (c *CreateView) End() Pos {
	return c.StatementEnd
}

func (c *CreateView) Type() string {
	return "VIEW"
}

func (c *CreateView) Accept(visitor ASTVisitor) error {
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
	if c.SubQuery != nil {
		if err := c.SubQuery.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateView(c)
}

type CreateFunction struct {
	CreatePos    Pos
	OrReplace    bool
	IfNotExists  bool
	FunctionName *Ident
	OnCluster    *ClusterClause
	Params       *ParamExprList
	Expr         Expr
}

func (c *CreateFunction) Type() string {
	return "FUNCTION"
}

func (c *CreateFunction) Pos() Pos {
	return c.CreatePos
}

func (c *CreateFunction) End() Pos {
	return c.Expr.End()
}

func (c *CreateFunction) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.FunctionName.Accept(visitor); err != nil {
		return err
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if err := c.Params.Accept(visitor); err != nil {
		return err
	}
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCreateFunction(c)
}

type RoleName struct {
	Name      Expr
	Scope     *StringLiteral
	OnCluster *ClusterClause
}

func (r *RoleName) Pos() Pos {
	return r.Name.Pos()
}

func (r *RoleName) End() Pos {
	if r.Scope != nil {
		return r.Scope.End()
	}
	if r.OnCluster != nil {
		return r.OnCluster.End()
	}
	return r.Name.End()
}

func (r *RoleName) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.Name.Accept(visitor); err != nil {
		return err
	}
	if r.Scope != nil {
		if err := r.Scope.Accept(visitor); err != nil {
			return err
		}
	}
	if r.OnCluster != nil {
		if err := r.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRoleName(r)
}

type SettingPair struct {
	Name      *Ident
	Operation TokenKind
	Value     Expr
}

func (s *SettingPair) Pos() Pos {
	return s.Name.NamePos
}

func (s *SettingPair) End() Pos {
	return s.Value.End()
}

func (s *SettingPair) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	if s.Value != nil {
		if err := s.Value.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSettingPair(s)
}

type RoleSetting struct {
	SettingPairs []*SettingPair
	Modifier     *Ident
}

func (r *RoleSetting) Pos() Pos {
	if len(r.SettingPairs) > 0 {
		return r.SettingPairs[0].Pos()
	}
	return r.Modifier.NamePos
}

func (r *RoleSetting) End() Pos {
	if r.Modifier != nil {
		return r.Modifier.NameEnd
	}
	return r.SettingPairs[len(r.SettingPairs)-1].End()
}

func (r *RoleSetting) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	for _, settingPair := range r.SettingPairs {
		if err := settingPair.Accept(visitor); err != nil {
			return err
		}
	}
	if r.Modifier != nil {
		if err := r.Modifier.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRoleSetting(r)
}

type CreateRole struct {
	CreatePos         Pos
	StatementEnd      Pos
	IfNotExists       bool
	OrReplace         bool
	RoleNames         []*RoleName
	AccessStorageType *Ident
	Settings          []*RoleSetting
}

func (c *CreateRole) Pos() Pos {
	return c.CreatePos
}

func (c *CreateRole) End() Pos {
	return c.StatementEnd
}

func (c *CreateRole) Type() string {
	return "ROLE"
}

func (c *CreateRole) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, roleName := range c.RoleNames {
		if err := roleName.Accept(visitor); err != nil {
			return err
		}
	}
	if c.AccessStorageType != nil {
		if err := c.AccessStorageType.Accept(visitor); err != nil {
			return err
		}
	}
	for _, setting := range c.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateRole(c)
}

type AuthenticationClause struct {
	AuthPos       Pos
	AuthEnd       Pos
	NotIdentified bool
	AuthType      string // "no_password", "plaintext_password", "sha256_password", etc.
	AuthValue     *StringLiteral
	LdapServer    *StringLiteral
	KerberosRealm *StringLiteral
	IsKerberos    bool
}

func (a *AuthenticationClause) Pos() Pos {
	return a.AuthPos
}

func (a *AuthenticationClause) End() Pos {
	return a.AuthEnd
}

func (a *AuthenticationClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if a.AuthValue != nil {
		if err := a.AuthValue.Accept(visitor); err != nil {
			return err
		}
	}
	if a.LdapServer != nil {
		if err := a.LdapServer.Accept(visitor); err != nil {
			return err
		}
	}
	if a.KerberosRealm != nil {
		if err := a.KerberosRealm.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAuthenticationClause(a)
}

type HostClause struct {
	HostPos   Pos
	HostEnd   Pos
	HostType  string // "LOCAL", "NAME", "REGEXP", "IP", "LIKE", "ANY", "NONE"
	HostValue *StringLiteral
}

func (h *HostClause) Pos() Pos {
	return h.HostPos
}

func (h *HostClause) End() Pos {
	return h.HostEnd
}

func (h *HostClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(h)
	defer visitor.Leave(h)
	if h.HostValue != nil {
		if err := h.HostValue.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitHostClause(h)
}

type DefaultRoleClause struct {
	DefaultPos Pos
	DefaultEnd Pos
	Roles      []*RoleName
	None       bool
}

func (d *DefaultRoleClause) Pos() Pos {
	return d.DefaultPos
}

func (d *DefaultRoleClause) End() Pos {
	return d.DefaultEnd
}

func (d *DefaultRoleClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	for _, role := range d.Roles {
		if err := role.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDefaultRoleClause(d)
}

type GranteesClause struct {
	GranteesPos Pos
	GranteesEnd Pos
	Grantees    []*RoleName
	ExceptUsers []*RoleName
	Any         bool
	None        bool
}

func (g *GranteesClause) Pos() Pos {
	return g.GranteesPos
}

func (g *GranteesClause) End() Pos {
	return g.GranteesEnd
}

func (g *GranteesClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	for _, grantee := range g.Grantees {
		if err := grantee.Accept(visitor); err != nil {
			return err
		}
	}
	for _, except := range g.ExceptUsers {
		if err := except.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitGranteesClause(g)
}

type CreateUser struct {
	CreatePos       Pos
	StatementEnd    Pos
	IfNotExists     bool
	OrReplace       bool
	UserNames       []*RoleName
	Authentication  *AuthenticationClause
	ValidUntil      *StringLiteral
	Hosts           []*HostClause
	DefaultRole     *DefaultRoleClause
	DefaultDatabase *Ident
	DefaultDbNone   bool
	Grantees        *GranteesClause
	Settings        []*RoleSetting
}

func (c *CreateUser) Pos() Pos {
	return c.CreatePos
}

func (c *CreateUser) End() Pos {
	return c.StatementEnd
}

func (c *CreateUser) Type() string {
	return "USER"
}

func (c *CreateUser) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, userName := range c.UserNames {
		if err := userName.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Authentication != nil {
		if err := c.Authentication.Accept(visitor); err != nil {
			return err
		}
	}
	if c.ValidUntil != nil {
		if err := c.ValidUntil.Accept(visitor); err != nil {
			return err
		}
	}
	for _, host := range c.Hosts {
		if err := host.Accept(visitor); err != nil {
			return err
		}
	}
	if c.DefaultRole != nil {
		if err := c.DefaultRole.Accept(visitor); err != nil {
			return err
		}
	}
	if c.DefaultDatabase != nil {
		if err := c.DefaultDatabase.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Grantees != nil {
		if err := c.Grantees.Accept(visitor); err != nil {
			return err
		}
	}
	for _, setting := range c.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateUser(c)
}

type AlterRole struct {
	AlterPos        Pos
	StatementEnd    Pos
	IfExists        bool
	RoleRenamePairs []*RoleRenamePair
	Settings        []*RoleSetting
}

func (a *AlterRole) Pos() Pos {
	return a.AlterPos
}

func (a *AlterRole) End() Pos {
	return a.StatementEnd
}

func (a *AlterRole) Type() string {
	return "ROLE"
}

func (a *AlterRole) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	for _, roleRenamePair := range a.RoleRenamePairs {
		if err := roleRenamePair.Accept(visitor); err != nil {
			return err
		}
	}
	for _, setting := range a.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterRole(a)
}

type RoleRenamePair struct {
	RoleName     *RoleName
	NewName      Expr
	StatementEnd Pos
}

func (r *RoleRenamePair) Pos() Pos {
	return r.RoleName.Pos()
}

func (r *RoleRenamePair) End() Pos {
	return r.StatementEnd
}

func (r *RoleRenamePair) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.RoleName.Accept(visitor); err != nil {
		return err
	}
	if r.NewName != nil {
		if err := r.NewName.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRoleRenamePair(r)
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
	return d.TableIdentifier.End()
}

func (d *DestinationClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.TableIdentifier.Accept(visitor); err != nil {
		return err
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

type NullLiteral struct {
	NullPos Pos
}

func (n *NullLiteral) Pos() Pos {
	return n.NullPos
}

func (n *NullLiteral) End() Pos {
	return n.NullPos + 4
}

func (n *NullLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	return visitor.VisitNullLiteral(n)
}

type NotNullLiteral struct {
	NotPos      Pos
	NullLiteral *NullLiteral
}

func (n *NotNullLiteral) Pos() Pos {
	return n.NotPos
}

func (n *NotNullLiteral) End() Pos {
	return n.NullLiteral.End()
}

func (n *NotNullLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.NullLiteral.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNotNullLiteral(n)
}

type NestedIdentifier struct {
	Ident    *Ident
	DotIdent *Ident
}

func (n *NestedIdentifier) Pos() Pos {
	return n.Ident.Pos()
}

func (n *NestedIdentifier) End() Pos {
	if n.DotIdent != nil {
		return n.DotIdent.End()
	}
	return n.Ident.End()
}

func (n *NestedIdentifier) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Ident.Accept(visitor); err != nil {
		return err
	}
	if n.DotIdent != nil {
		if err := n.DotIdent.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitNestedIdentifier(n)
}

type Path struct {
	Fields []*Ident
}

func (p *Path) Pos() Pos {
	if len(p.Fields) > 0 {
		return p.Fields[0].Pos()
	}
	return 0
}

func (p *Path) End() Pos {
	if len(p.Fields) > 0 {
		return p.Fields[len(p.Fields)-1].End()
	}
	return 0
}

func (p *Path) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	for _, ident := range p.Fields {
		if err := ident.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPath(p)
}

type TableIdentifier struct {
	Database *Ident
	Table    *Ident
}

func (t *TableIdentifier) Pos() Pos {
	if t.Database != nil {
		return t.Database.Pos()
	}
	return t.Table.Pos()
}

func (t *TableIdentifier) End() Pos {
	return t.Table.End()
}

func (t *TableIdentifier) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if t.Database != nil {
		if err := t.Database.Accept(visitor); err != nil {
			return err
		}
	}
	if err := t.Table.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTableIdentifier(t)
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

type TableArgListExpr struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Args          []Expr
}

func (t *TableArgListExpr) Pos() Pos {
	return t.LeftParenPos
}

func (t *TableArgListExpr) End() Pos {
	return t.RightParenPos
}

func (t *TableArgListExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	for _, arg := range t.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTableArgListExpr(t)
}

type TableFunctionExpr struct {
	Name Expr
	Args *TableArgListExpr
}

func (t *TableFunctionExpr) Pos() Pos {
	return t.Name.Pos()
}

func (t *TableFunctionExpr) End() Pos {
	return t.Args.End()
}

func (t *TableFunctionExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Name.Accept(visitor); err != nil {
		return err
	}
	if err := t.Args.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTableFunctionExpr(t)
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

type RefreshExpr struct {
	RefreshPos Pos
	Frequency  string // EVERY|AFTER
	Interval   *IntervalExpr
	Offset     *IntervalExpr
}

func (r *RefreshExpr) Pos() Pos {
	return r.RefreshPos
}

func (r *RefreshExpr) End() Pos {
	if r.Offset != nil {
		return r.Offset.End()
	}
	return r.Interval.End()
}

func (r *RefreshExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if r.Interval != nil {
		if err := r.Interval.Accept(visitor); err != nil {
			return err
		}
	}
	if r.Offset != nil {
		if err := r.Offset.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRefreshExpr(r)
}

type TTLPolicyRule struct {
	RulePos  Pos
	ToVolume *StringLiteral
	ToDisk   *StringLiteral
	Action   *TTLPolicyRuleAction
}

func (t *TTLPolicyRule) Pos() Pos {
	return t.RulePos
}

func (t *TTLPolicyRule) End() Pos {
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
	return visitor.VisitTTLPolicyRule(t)
}

type TTLPolicy struct {
	Item    *TTLPolicyRule
	Where   *WhereClause
	GroupBy *GroupByClause
}

func (t *TTLPolicy) Pos() Pos {
	if t.Item != nil {
		return t.Item.Pos()
	}
	if t.Where != nil {
		return t.Where.Pos()
	}
	return t.GroupBy.Pos()
}

func (t *TTLPolicy) End() Pos {
	if t.GroupBy != nil {
		return t.GroupBy.End()
	}
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
	if t.GroupBy != nil {
		if err := t.GroupBy.Accept(visitor); err != nil {
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

type Fill struct {
	FillPos   Pos
	From      Expr // optional
	To        Expr // optional
	Step      Expr // optional
	Staleness Expr // optional
}

func (f *Fill) Pos() Pos {
	return f.FillPos
}

func (f *Fill) End() Pos {
	if f.Staleness != nil {
		return f.Staleness.End()
	}
	if f.Step != nil {
		return f.Step.End()
	}
	if f.To != nil {
		return f.To.End()
	}
	if f.From != nil {
		return f.From.End()
	}
	return f.FillPos + Pos(len("FILL"))
}

func (f *Fill) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if f.From != nil {
		if err := f.From.Accept(visitor); err != nil {
			return err
		}
	}
	if f.To != nil {
		if err := f.To.Accept(visitor); err != nil {
			return err
		}
	}
	if f.Step != nil {
		if err := f.Step.Accept(visitor); err != nil {
			return err
		}
	}
	if f.Staleness != nil {
		if err := f.Staleness.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitFill(f)
}

type OrderExpr struct {
	OrderPos  Pos
	Expr      Expr
	Alias     *Ident
	Direction OrderDirection
	Fill      *Fill // optional WITH FILL clause
}

func (o *OrderExpr) Pos() Pos {
	return o.OrderPos
}

func (o *OrderExpr) End() Pos {
	if o.Fill != nil {
		return o.Fill.End()
	}
	if o.Alias != nil {
		return o.Alias.End()
	}
	return o.Expr.End()
}

func (o *OrderExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.Expr.Accept(visitor); err != nil {
		return err
	}
	if o.Alias != nil {
		if err := o.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Fill != nil {
		if err := o.Fill.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitOrderByExpr(o)
}

type InterpolateItem struct {
	Column *Ident
	Expr   Expr // optional AS expression
}

func (i *InterpolateItem) Pos() Pos {
	return i.Column.Pos()
}

func (i *InterpolateItem) End() Pos {
	if i.Expr != nil {
		return i.Expr.End()
	}
	return i.Column.End()
}

func (i *InterpolateItem) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.Column.Accept(visitor); err != nil {
		return err
	}
	if i.Expr != nil {
		if err := i.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitInterpolateItem(i)
}

type InterpolateClause struct {
	InterpolatePos Pos
	ListEnd        Pos
	Items          []*InterpolateItem // can be nil for INTERPOLATE without columns
}

func (i *InterpolateClause) Pos() Pos {
	return i.InterpolatePos
}

func (i *InterpolateClause) End() Pos {
	return i.ListEnd
}

func (i *InterpolateClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	for _, item := range i.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitInterpolateClause(i)
}

type OrderByClause struct {
	OrderPos    Pos
	ListEnd     Pos
	Items       []Expr
	Interpolate *InterpolateClause // optional INTERPOLATE clause
}

func (o *OrderByClause) Pos() Pos {
	return o.OrderPos
}

func (o *OrderByClause) End() Pos {
	if o.Interpolate != nil {
		return o.Interpolate.End()
	}
	return o.ListEnd
}

func (o *OrderByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	for _, item := range o.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Interpolate != nil {
		if err := o.Interpolate.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitOrderByListExpr(o)
}

type SettingExpr struct {
	SettingsPos Pos
	Name        *Ident
	Expr        Expr
}

func (s *SettingExpr) Pos() Pos {
	return s.SettingsPos
}

func (s *SettingExpr) End() Pos {
	return s.Expr.End()
}

func (s *SettingExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	if err := s.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSettingsExpr(s)
}

type SettingsClause struct {
	SettingsPos Pos
	ListEnd     Pos
	Items       []*SettingExpr
}

func (s *SettingsClause) Pos() Pos {
	return s.SettingsPos
}

func (s *SettingsClause) End() Pos {
	return s.ListEnd
}

func (s *SettingsClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	for _, item := range s.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSettingsExprList(s)
}

type ParamExprList struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Items         *ColumnExprList
	ColumnArgList *ColumnArgList
}

func (f *ParamExprList) Pos() Pos {
	return f.LeftParenPos
}

func (f *ParamExprList) End() Pos {
	return f.RightParenPos
}

func (f *ParamExprList) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Items.Accept(visitor); err != nil {
		return err
	}
	if f.ColumnArgList != nil {
		if err := f.ColumnArgList.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitParamExprList(f)
}

type KeyValue struct {
	Key   StringLiteral
	Value Expr
}

type MapLiteral struct {
	LBracePos Pos
	RBracePos Pos
	KeyValues []KeyValue
}

func (m *MapLiteral) Pos() Pos {
	return m.LBracePos
}

func (m *MapLiteral) End() Pos {
	return m.RBracePos
}

func (m *MapLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(m)
	defer visitor.Leave(m)
	for _, kv := range m.KeyValues {
		if err := kv.Key.Accept(visitor); err != nil {
			return err
		}
		if err := kv.Value.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitMapLiteral(m)
}

type NamedParameterExpr struct {
	NamePos Pos
	Name    *Ident
	Value   Expr
}

func (n *NamedParameterExpr) Pos() Pos {
	return n.NamePos
}

func (n *NamedParameterExpr) End() Pos {
	return n.Value.End()
}

func (n *NamedParameterExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Name.Accept(visitor); err != nil {
		return err
	}
	if err := n.Value.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNamedParameterExpr(n)
}

type QueryParam struct {
	LBracePos Pos
	RBracePos Pos
	Name      *Ident
	Type      ColumnType
}

func (q *QueryParam) Pos() Pos {
	return q.LBracePos
}

func (q *QueryParam) End() Pos {
	return q.RBracePos
}

func (q *QueryParam) Accept(visitor ASTVisitor) error {
	visitor.Enter(q)
	defer visitor.Leave(q)
	if err := q.Name.Accept(visitor); err != nil {
		return err
	}
	if err := q.Type.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitQueryParam(q)
}

type ArrayParamList struct {
	LeftBracketPos  Pos
	RightBracketPos Pos
	Items           *ColumnExprList
}

func (a *ArrayParamList) Pos() Pos {
	return a.LeftBracketPos
}

func (a *ArrayParamList) End() Pos {
	return a.RightBracketPos
}

func (a *ArrayParamList) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Items.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitArrayParamList(a)
}

type ObjectParams struct {
	Object Expr
	Params *ArrayParamList
}

func (o *ObjectParams) Pos() Pos {
	return o.Object.Pos()
}

func (o *ObjectParams) End() Pos {
	return o.Params.End()
}

func (o *ObjectParams) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.Object.Accept(visitor); err != nil {
		return err
	}
	if err := o.Params.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitObjectParams(o)
}

type FunctionExpr struct {
	Name   *Ident
	Params *ParamExprList
}

func (f *FunctionExpr) Pos() Pos {
	return f.Name.NamePos
}

func (f *FunctionExpr) End() Pos {
	return f.Params.RightParenPos
}

func (f *FunctionExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Name.Accept(visitor); err != nil {
		return err
	}
	if err := f.Params.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitFunctionExpr(f)
}

type WindowFunctionExpr struct {
	Function *FunctionExpr
	OverPos  Pos
	OverExpr Expr
}

func (w *WindowFunctionExpr) Pos() Pos {
	return w.Function.Pos()
}

func (w *WindowFunctionExpr) End() Pos {
	return w.OverExpr.End()
}

func (w *WindowFunctionExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.Function.Accept(visitor); err != nil {
		return err
	}
	if err := w.OverExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFunctionExpr(w)
}

type TypedPlaceholder struct {
	LeftBracePos  Pos
	RightBracePos Pos
	Name          *Ident
	Type          ColumnType
}

func (t *TypedPlaceholder) Pos() Pos {
	return t.LeftBracePos
}

func (t *TypedPlaceholder) End() Pos {
	return t.RightBracePos
}

func (t *TypedPlaceholder) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Name.Accept(visitor); err != nil {
		return err
	}
	if err := t.Type.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTypedPlaceholder(t)
}

type ColumnExpr struct {
	Expr  Expr
	Alias *Ident
}

func (c *ColumnExpr) Pos() Pos {
	return c.Expr.Pos()
}

func (c *ColumnExpr) End() Pos {
	if c.Alias != nil {
		return c.Alias.NameEnd
	}
	return c.Expr.End()
}

func (c *ColumnExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	if c.Alias != nil {
		if err := c.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnExpr(c)
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
	return visitor.VisitColumnDef(c)
}

type ColumnType interface {
	Expr
	Type() string
}

type ScalarType struct {
	Name *Ident
}

func (s *ScalarType) Pos() Pos {
	return s.Name.NamePos
}

func (s *ScalarType) End() Pos {
	return s.Name.NameEnd
}

func (s *ScalarType) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitScalarType(s)
}

func (s *ScalarType) Type() string {
	return s.Name.Name
}

type JSONPath struct {
	Idents []*Ident
}

type JSONTypeHint struct {
	Path *JSONPath
	Type ColumnType
}

type JSONOption struct {
	SkipPath        *JSONPath
	SkipRegex       *StringLiteral
	MaxDynamicPaths *NumberLiteral
	MaxDynamicTypes *NumberLiteral
	// Type hint for specific JSON subcolumn path, e.g., "message String" or "a.b UInt64"
	Column *JSONTypeHint
}

type JSONOptions struct {
	LParen Pos
	RParen Pos
	Items  []*JSONOption
}

func (j *JSONOptions) Pos() Pos {
	return j.LParen
}

func (j *JSONOptions) End() Pos {
	return j.RParen
}

type JSONType struct {
	Name    *Ident
	Options *JSONOptions
}

func (j *JSONType) Pos() Pos {
	return j.Name.NamePos
}

func (j *JSONType) End() Pos {
	if j.Options != nil {
		return j.Options.RParen
	}
	return j.Name.NameEnd
}

func (j *JSONType) Type() string {
	return j.Name.Name
}

func (j *JSONType) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if err := j.Name.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitJSONType(j)
}

type PropertyType struct {
	Name *Ident
}

func (c *PropertyType) Pos() Pos {
	return c.Name.NamePos
}

func (c *PropertyType) End() Pos {
	return c.Name.NameEnd
}

func (c *PropertyType) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPropertyType(c)
}

func (c *PropertyType) Type() string {
	return c.Name.Name
}

type TypeWithParams struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Name          *Ident
	Params        []Literal
}

func (s *TypeWithParams) Pos() Pos {
	return s.Name.NamePos
}

func (s *TypeWithParams) End() Pos {
	return s.RightParenPos
}

func (s *TypeWithParams) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	for _, param := range s.Params {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTypeWithParams(s)
}

func (s *TypeWithParams) Type() string {
	return s.Name.Name
}

type ComplexType struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Name          *Ident
	Params        []ColumnType
}

func (c *ComplexType) Pos() Pos {
	return c.Name.NamePos
}

func (c *ComplexType) End() Pos {
	return c.RightParenPos
}

func (c *ComplexType) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	for _, param := range c.Params {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitComplexType(c)
}

func (c *ComplexType) Type() string {
	return c.Name.Name
}

type NestedType struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Name          *Ident
	Columns       []Expr
}

func (n *NestedType) Pos() Pos {
	return n.Name.NamePos
}

func (n *NestedType) End() Pos {
	return n.RightParenPos
}

func (n *NestedType) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Name.Accept(visitor); err != nil {
		return err
	}
	for _, column := range n.Columns {
		if err := column.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitNestedType(n)
}

func (n *NestedType) Type() string {
	return n.Name.Name
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

type Literal interface {
	Expr
}

type NumberLiteral struct {
	NumPos  Pos
	NumEnd  Pos
	Literal string
	Base    int
}

func (n *NumberLiteral) Pos() Pos {
	return n.NumPos
}

func (n *NumberLiteral) End() Pos {
	return n.NumEnd
}

func (n *NumberLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	return visitor.VisitNumberLiteral(n)
}

type StringLiteral struct {
	LiteralPos Pos
	LiteralEnd Pos
	Literal    string
}

func (s *StringLiteral) Pos() Pos {
	return s.LiteralPos
}

func (s *StringLiteral) End() Pos {
	return s.LiteralEnd
}

func (s *StringLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	return visitor.VisitStringLiteral(s)
}

type BoolLiteral struct {
	LiteralPos Pos
	LiteralEnd Pos
	Literal    string
}

func (b *BoolLiteral) Pos() Pos {
	return b.LiteralPos
}

func (b *BoolLiteral) End() Pos {
	return b.LiteralEnd
}

func (b *BoolLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(b)
	defer visitor.Leave(b)
	return visitor.VisitBoolLiteral(b)
}

type PlaceHolder struct {
	PlaceholderPos Pos
	PlaceHolderEnd Pos
	Type           string
}

func (p *PlaceHolder) Pos() Pos {
	return p.PlaceholderPos
}

func (p *PlaceHolder) End() Pos {
	return p.PlaceHolderEnd
}

func (p *PlaceHolder) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	return visitor.VisitPlaceHolderExpr(p)
}

type RatioExpr struct {
	Numerator *NumberLiteral
	// numberLiteral (SLASH numberLiteral)?
	Denominator *NumberLiteral
}

func (r *RatioExpr) Pos() Pos {
	return r.Numerator.NumPos
}

func (r *RatioExpr) End() Pos {
	if r.Denominator != nil {
		return r.Denominator.NumEnd
	}
	return r.Numerator.NumEnd
}

func (r *RatioExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.Numerator.Accept(visitor); err != nil {
		return err
	}
	if r.Denominator != nil {
		if err := r.Denominator.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRatioExpr(r)
}

type EnumValue struct {
	Name  *StringLiteral
	Value *NumberLiteral
}

func (e *EnumValue) Pos() Pos {
	return e.Name.Pos()
}

func (e *EnumValue) End() Pos {
	return e.Value.End()
}

func (e *EnumValue) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	if err := e.Name.Accept(visitor); err != nil {
		return err
	}
	if err := e.Value.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitEnumValue(e)
}

type EnumType struct {
	Name    *Ident
	ListPos Pos
	ListEnd Pos
	Values  []EnumValue
}

func (e *EnumType) Pos() Pos {
	return e.ListPos
}

func (e *EnumType) End() Pos {
	return e.ListEnd
}

func (e *EnumType) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	if err := e.Name.Accept(visitor); err != nil {
		return err
	}
	for i := range e.Values {
		if err := e.Values[i].Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitEnumType(e)
}

func (e *EnumType) Type() string {
	return e.Name.Name
}

type IntervalExpr struct {
	// INTERVAL keyword position which might be omitted(IntervalPos = 0)
	IntervalPos Pos

	Expr Expr
	Unit *Ident
}

func (i *IntervalExpr) Pos() Pos {
	if i.IntervalPos != 0 {
		return i.IntervalPos
	}
	return i.Expr.Pos()
}

func (i *IntervalExpr) End() Pos {
	return i.Unit.End()
}

func (i *IntervalExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := i.Unit.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIntervalExpr(i)
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

type ColumnTypeExpr struct {
	Name *Ident
}

func (c *ColumnTypeExpr) Pos() Pos {
	return c.Name.NamePos
}

func (c *ColumnTypeExpr) End() Pos {
	return c.Name.NameEnd
}

func (c *ColumnTypeExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitColumnTypeExpr(c)
}

type ColumnArgList struct {
	Distinct      bool
	LeftParenPos  Pos
	RightParenPos Pos
	Items         []Expr
}

func (c *ColumnArgList) Pos() Pos {
	return c.LeftParenPos
}

func (c *ColumnArgList) End() Pos {
	return c.RightParenPos
}

func (c *ColumnArgList) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, item := range c.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnArgList(c)
}

type ColumnExprList struct {
	ListPos     Pos
	ListEnd     Pos
	HasDistinct bool
	Items       []Expr
}

func (c *ColumnExprList) Pos() Pos {
	return c.ListPos
}

func (c *ColumnExprList) End() Pos {
	return c.ListEnd
}

func (c *ColumnExprList) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, item := range c.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnExprList(c)
}

type WhenClause struct {
	WhenPos Pos
	ThenPos Pos
	When    Expr
	Then    Expr
	ElsePos Pos
	Else    Expr
}

func (w *WhenClause) Pos() Pos {
	return w.WhenPos
}

func (w *WhenClause) End() Pos {
	if w.Else != nil {
		return w.Else.End()
	}
	return w.Then.End()
}

func (w *WhenClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.When.Accept(visitor); err != nil {
		return err
	}
	if err := w.Then.Accept(visitor); err != nil {
		return err
	}
	if w.Else != nil {
		if err := w.Else.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitWhenExpr(w)
}

type CaseExpr struct {
	CasePos Pos
	EndPos  Pos
	Expr    Expr // optional
	Whens   []*WhenClause
	ElsePos Pos
	Else    Expr
}

func (c *CaseExpr) Pos() Pos {
	return c.CasePos
}

func (c *CaseExpr) End() Pos {
	return c.EndPos
}

func (c *CaseExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if c.Expr != nil {
		if err := c.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	for _, when := range c.Whens {
		if err := when.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Else != nil {
		if err := c.Else.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCaseExpr(c)
}

type CastExpr struct {
	CastPos   Pos
	Expr      Expr
	Separator string
	AsPos     Pos
	AsType    Expr
}

func (c *CastExpr) Pos() Pos {
	return c.CastPos
}

func (c *CastExpr) End() Pos {
	return c.AsType.End()
}

func (c *CastExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := c.AsType.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCastExpr(c)
}

type WithClause struct {
	WithPos Pos
	EndPos  Pos
	CTEs    []*CTEStmt
}

func (w *WithClause) Pos() Pos {
	return w.WithPos
}

func (w *WithClause) End() Pos {
	return w.EndPos
}

func (w *WithClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	for _, cte := range w.CTEs {
		if err := cte.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitWithExpr(w)
}

type TopClause struct {
	TopPos   Pos
	TopEnd   Pos
	Number   *NumberLiteral
	WithTies bool
}

func (t *TopClause) Pos() Pos {
	return t.TopPos
}

func (t *TopClause) End() Pos {
	return t.TopEnd
}

func (t *TopClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Number.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTopExpr(t)
}

type CreateLiveView struct {
	CreatePos    Pos
	StatementEnd Pos
	Name         *TableIdentifier
	IfNotExists  bool
	UUID         *UUID
	OnCluster    *ClusterClause
	Destination  *DestinationClause
	TableSchema  *TableSchemaClause
	WithTimeout  *WithTimeoutClause
	SubQuery     *SubQuery
}

func (c *CreateLiveView) Type() string {
	return "LIVE_VIEW"
}

func (c *CreateLiveView) Pos() Pos {
	return c.CreatePos
}

func (c *CreateLiveView) End() Pos {
	return c.StatementEnd
}

func (c *CreateLiveView) Accept(visitor ASTVisitor) error {
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
	if c.Destination != nil {
		if err := c.Destination.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TableSchema != nil {
		if err := c.TableSchema.Accept(visitor); err != nil {
			return err
		}
	}
	if c.WithTimeout != nil {
		if err := c.WithTimeout.Accept(visitor); err != nil {
			return err
		}
	}
	if c.SubQuery != nil {
		if err := c.SubQuery.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateLiveView(c)
}

type CreateDictionary struct {
	CreatePos    Pos
	StatementEnd Pos
	OrReplace    bool
	Name         *TableIdentifier
	IfNotExists  bool
	UUID         *UUID
	OnCluster    *ClusterClause
	Schema       *DictionarySchemaClause
	Engine       *DictionaryEngineClause
	Comment      *StringLiteral
}

func (c *CreateDictionary) Type() string {
	return "DICTIONARY"
}

func (c *CreateDictionary) Pos() Pos {
	return c.CreatePos
}

func (c *CreateDictionary) End() Pos {
	return c.StatementEnd
}

func (c *CreateDictionary) Accept(visitor ASTVisitor) error {
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
	if c.Schema != nil {
		if err := c.Schema.Accept(visitor); err != nil {
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
	return visitor.VisitCreateDictionary(c)
}

type CreateNamedCollection struct {
	CreatePos    Pos
	StatementEnd Pos
	Name         *Ident
	IfNotExists  bool
	OnCluster    *ClusterClause
	Params       []*NamedCollectionParam
}

func (c *CreateNamedCollection) Pos() Pos {
	return c.CreatePos
}

func (c *CreateNamedCollection) End() Pos {
	return c.StatementEnd
}

func (c *CreateNamedCollection) Type() string {
	return "NAMED COLLECTION"
}

func (c *CreateNamedCollection) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	for _, param := range c.Params {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateNamedCollection(c)
}

type NamedCollectionParam struct {
	ParamPos       Pos
	Name           *Ident
	Value          Expr
	Overridable    bool
	NotOverridable bool
}

func (n *NamedCollectionParam) Pos() Pos {
	return n.ParamPos
}

func (n *NamedCollectionParam) End() Pos {
	return n.Value.End()
}

func (n *NamedCollectionParam) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Name.Accept(visitor); err != nil {
		return err
	}
	if err := n.Value.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNamedCollectionParam(n)
}

type DictionarySchemaClause struct {
	SchemaPos  Pos
	Attributes []*DictionaryAttribute
	RParenPos  Pos
}

func (d *DictionarySchemaClause) Pos() Pos {
	return d.SchemaPos
}

func (d *DictionarySchemaClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionarySchemaClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	for _, attr := range d.Attributes {
		if err := attr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionarySchemaClause(d)
}

type DictionaryAttribute struct {
	NamePos      Pos
	Name         *Ident
	Type         ColumnType
	Default      Literal
	Expression   Expr
	Hierarchical bool
	Injective    bool
	IsObjectId   bool
}

func (d *DictionaryAttribute) Pos() Pos {
	return d.NamePos
}

func (d *DictionaryAttribute) End() Pos {
	if d.IsObjectId {
		return d.NamePos + Pos(len("IS_OBJECT_ID"))
	}
	if d.Injective {
		return d.NamePos + Pos(len("INJECTIVE"))
	}
	if d.Hierarchical {
		return d.NamePos + Pos(len("HIERARCHICAL"))
	}
	if d.Expression != nil {
		return d.Expression.End()
	}
	if d.Default != nil {
		return d.Default.End()
	}
	return d.Type.End()
}

func (d *DictionaryAttribute) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Name.Accept(visitor); err != nil {
		return err
	}
	if err := d.Type.Accept(visitor); err != nil {
		return err
	}
	if d.Default != nil {
		if err := d.Default.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Expression != nil {
		if err := d.Expression.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryAttribute(d)
}

type DictionaryEngineClause struct {
	EnginePos  Pos
	PrimaryKey *DictionaryPrimaryKeyClause
	Source     *DictionarySourceClause
	Lifetime   *DictionaryLifetimeClause
	Layout     *DictionaryLayoutClause
	Range      *DictionaryRangeClause
	Settings   *SettingsClause
}

func (d *DictionaryEngineClause) Pos() Pos {
	return d.EnginePos
}

func (d *DictionaryEngineClause) End() Pos {
	if d.Settings != nil {
		return d.Settings.End()
	}
	if d.Range != nil {
		return d.Range.End()
	}
	if d.Layout != nil {
		return d.Layout.End()
	}
	if d.Lifetime != nil {
		return d.Lifetime.End()
	}
	if d.Source != nil {
		return d.Source.End()
	}
	if d.PrimaryKey != nil {
		return d.PrimaryKey.End()
	}
	return d.EnginePos
}

func (d *DictionaryEngineClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if d.PrimaryKey != nil {
		if err := d.PrimaryKey.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Source != nil {
		if err := d.Source.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Lifetime != nil {
		if err := d.Lifetime.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Layout != nil {
		if err := d.Layout.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Range != nil {
		if err := d.Range.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Settings != nil {
		if err := d.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryEngineClause(d)
}

type DictionaryPrimaryKeyClause struct {
	PrimaryKeyPos Pos
	Keys          *ColumnExprList
	RParenPos     Pos
}

func (d *DictionaryPrimaryKeyClause) Pos() Pos {
	return d.PrimaryKeyPos
}

func (d *DictionaryPrimaryKeyClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryPrimaryKeyClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Keys.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitDictionaryPrimaryKeyClause(d)
}

type DictionarySourceClause struct {
	SourcePos Pos
	Source    *Ident
	Args      []*DictionaryArgExpr
	RParenPos Pos
}

func (d *DictionarySourceClause) Pos() Pos {
	return d.SourcePos
}

func (d *DictionarySourceClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionarySourceClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Source.Accept(visitor); err != nil {
		return err
	}
	for _, arg := range d.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionarySourceClause(d)
}

type DictionaryArgExpr struct {
	ArgPos Pos
	Name   *Ident
	Value  Expr // can be Ident with optional parentheses or literal
}

func (d *DictionaryArgExpr) Pos() Pos {
	return d.ArgPos
}

func (d *DictionaryArgExpr) End() Pos {
	return d.Value.End()
}

func (d *DictionaryArgExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Name.Accept(visitor); err != nil {
		return err
	}
	if err := d.Value.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitDictionaryArgExpr(d)
}

type DictionaryLifetimeClause struct {
	LifetimePos Pos
	Min         *NumberLiteral
	Max         *NumberLiteral
	Value       *NumberLiteral // for simple LIFETIME(value) form
	RParenPos   Pos
}

func (d *DictionaryLifetimeClause) Pos() Pos {
	return d.LifetimePos
}

func (d *DictionaryLifetimeClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryLifetimeClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if d.Value != nil {
		if err := d.Value.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Min != nil {
		if err := d.Min.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Max != nil {
		if err := d.Max.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryLifetimeClause(d)
}

type DictionaryLayoutClause struct {
	LayoutPos Pos
	Layout    *Ident
	Args      []*DictionaryArgExpr
	RParenPos Pos
}

func (d *DictionaryLayoutClause) Pos() Pos {
	return d.LayoutPos
}

func (d *DictionaryLayoutClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryLayoutClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Layout.Accept(visitor); err != nil {
		return err
	}
	for _, arg := range d.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryLayoutClause(d)
}

type DictionaryRangeClause struct {
	RangePos  Pos
	Min       *Ident
	Max       *Ident
	RParenPos Pos
}

func (d *DictionaryRangeClause) Pos() Pos {
	return d.RangePos
}

func (d *DictionaryRangeClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryRangeClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Min.Accept(visitor); err != nil {
		return err
	}
	if err := d.Max.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitDictionaryRangeClause(d)
}

type WithTimeoutClause struct {
	WithTimeoutPos Pos
	Expr           Expr
	Number         *NumberLiteral
}

func (w *WithTimeoutClause) Pos() Pos {
	return w.WithTimeoutPos
}

func (w *WithTimeoutClause) End() Pos {
	return w.Number.End()
}

func (w *WithTimeoutClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.Number.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWithTimeoutExpr(w)
}

type TableExpr struct {
	TablePos Pos
	TableEnd Pos
	Alias    *AliasExpr
	Expr     Expr
	HasFinal bool
}

func (t *TableExpr) Pos() Pos {
	return t.TablePos
}

func (t *TableExpr) End() Pos {
	return t.TableEnd
}

func (t *TableExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Expr.Accept(visitor); err != nil {
		return err
	}
	if t.Alias != nil {
		if err := t.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTableExpr(t)
}

type OnClause struct {
	OnPos Pos
	On    *ColumnExprList
}

func (o *OnClause) Pos() Pos {
	return o.OnPos
}

func (o *OnClause) End() Pos {
	return o.On.End()
}

func (o *OnClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.On.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitOnExpr(o)
}

type UsingClause struct {
	UsingPos Pos
	Using    *ColumnExprList
}

func (u *UsingClause) Pos() Pos {
	return u.UsingPos
}

func (u *UsingClause) End() Pos {
	return u.Using.End()
}

func (u *UsingClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	if err := u.Using.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUsingExpr(u)
}

type JoinExpr struct {
	JoinPos     Pos
	Left        Expr
	Right       Expr
	Modifiers   []string
	Constraints Expr
}

func (j *JoinExpr) Pos() Pos {
	return j.JoinPos
}

func (j *JoinExpr) End() Pos {
	// Return the rightmost position
	if j.Right != nil {
		return j.Right.End()
	}
	if j.Constraints != nil {
		return j.Constraints.End()
	}
	return j.Left.End()
}

func (j *JoinExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if err := j.Left.Accept(visitor); err != nil {
		return err
	}
	if j.Right != nil {
		if err := j.Right.Accept(visitor); err != nil {
			return err
		}
	}
	if j.Constraints != nil {
		if err := j.Constraints.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitJoinExpr(j)
}

type JoinConstraintClause struct {
	ConstraintPos Pos
	On            *ColumnExprList
	Using         *ColumnExprList
}

func (j *JoinConstraintClause) Pos() Pos {
	return j.ConstraintPos
}

func (j *JoinConstraintClause) End() Pos {
	if j.On != nil {
		return j.On.End()
	}
	return j.Using.End()
}

func (j *JoinConstraintClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if j.On != nil {
		if err := j.On.Accept(visitor); err != nil {
			return err
		}
	}
	if j.Using != nil {
		if err := j.Using.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitJoinConstraintExpr(j)
}

type FromClause struct {
	FromPos Pos
	Expr    Expr
}

func (f *FromClause) Pos() Pos {
	return f.FromPos
}

func (f *FromClause) End() Pos {
	return f.Expr.End()
}

func (f *FromClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitFromExpr(f)
}

type IsNullExpr struct {
	IsPos Pos
	Expr  Expr
}

func (n *IsNullExpr) Pos() Pos {
	return n.IsPos
}

func (n *IsNullExpr) End() Pos {
	return n.Expr.End()
}

func (n *IsNullExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIsNullExpr(n)
}

type IsNotNullExpr struct {
	IsPos Pos
	Expr  Expr
}

func (n *IsNotNullExpr) Pos() Pos {
	return n.Expr.Pos()
}

func (n *IsNotNullExpr) End() Pos {
	return n.Expr.End()
}

func (n *IsNotNullExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIsNotNullExpr(n)
}

type AliasExpr struct {
	Expr     Expr
	AliasPos Pos
	Alias    Expr
}

func (a *AliasExpr) Pos() Pos {
	return a.AliasPos
}

func (a *AliasExpr) End() Pos {
	return a.Alias.End()
}

func (a *AliasExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := a.Alias.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAliasExpr(a)
}

type WhereClause struct {
	WherePos Pos
	Expr     Expr
}

func (w *WhereClause) Pos() Pos {
	return w.WherePos
}

func (w *WhereClause) End() Pos {
	return w.Expr.End()
}

func (w *WhereClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWhereExpr(w)
}

type PrewhereClause struct {
	PrewherePos Pos
	Expr        Expr
}

func (w *PrewhereClause) Pos() Pos {
	return w.PrewherePos
}

func (w *PrewhereClause) End() Pos {
	return w.Expr.End()
}

func (w *PrewhereClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPrewhereExpr(w)
}

type GroupByClause struct {
	GroupByPos    Pos
	GroupByEnd    Pos
	AggregateType string
	Expr          Expr
	WithCube      bool
	WithRollup    bool
	WithTotals    bool
}

func (g *GroupByClause) Pos() Pos {
	return g.GroupByPos
}

func (g *GroupByClause) End() Pos {
	return g.GroupByEnd
}

func (g *GroupByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	if g.Expr != nil {
		if err := g.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitGroupByExpr(g)
}

type HavingClause struct {
	HavingPos Pos
	Expr      Expr
}

func (h *HavingClause) Pos() Pos {
	return h.HavingPos
}

func (h *HavingClause) End() Pos {
	return h.Expr.End()
}

func (h *HavingClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(h)
	defer visitor.Leave(h)
	if err := h.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitHavingExpr(h)
}

type LimitClause struct {
	LimitPos Pos
	Limit    Expr
	Offset   Expr
}

func (l *LimitClause) Pos() Pos {
	return l.LimitPos
}

func (l *LimitClause) End() Pos {
	if l.Offset != nil {
		return l.Offset.End()
	}
	return l.Limit.End()
}

func (l *LimitClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(l)
	defer visitor.Leave(l)
	if l.Limit != nil {
		if err := l.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if l.Offset != nil {
		if err := l.Offset.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitLimitExpr(l)
}

type LimitByClause struct {
	Limit  *LimitClause
	ByExpr *ColumnExprList
}

func (l *LimitByClause) Pos() Pos {
	return l.Limit.Pos()
}

func (l *LimitByClause) End() Pos {
	if l.ByExpr != nil {
		return l.ByExpr.End()
	}
	if l.Limit != nil {
		return l.Limit.End()
	}
	return l.Limit.End()
}

func (l *LimitByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(l)
	defer visitor.Leave(l)
	if l.Limit != nil {
		if err := l.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if l.ByExpr != nil {
		if err := l.ByExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitLimitByExpr(l)
}

type WindowExpr struct {
	LeftParenPos  Pos
	RightParenPos Pos
	PartitionBy   *PartitionByClause
	OrderBy       *OrderByClause
	Frame         *WindowFrameClause
}

func (w *WindowExpr) Pos() Pos {
	return w.LeftParenPos
}

func (w *WindowExpr) End() Pos {
	return w.RightParenPos
}

func (w *WindowExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if w.PartitionBy != nil {
		if err := w.PartitionBy.Accept(visitor); err != nil {
			return err
		}
	}
	if w.OrderBy != nil {
		if err := w.OrderBy.Accept(visitor); err != nil {
			return err
		}
	}
	if w.Frame != nil {
		if err := w.Frame.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitWindowConditionExpr(w)
}

type WindowDefinition struct {
	Name  *Ident
	AsPos Pos
	Expr  *WindowExpr
}

func (w *WindowDefinition) Pos() Pos {
	if w == nil || w.Name == nil {
		return 0
	}
	return w.Name.Pos()
}

func (w *WindowDefinition) End() Pos {
	if w == nil || w.Expr == nil {
		return 0
	}
	return w.Expr.End()
}

type WindowClause struct {
	WindowPos Pos
	EndPos    Pos
	Windows   []*WindowDefinition
}

func (w *WindowClause) Pos() Pos {
	return w.WindowPos
}

func (w *WindowClause) End() Pos {
	if w.EndPos != 0 {
		return w.EndPos
	}
	if len(w.Windows) == 0 {
		return w.WindowPos
	}
	return w.Windows[len(w.Windows)-1].End()
}

func (w *WindowClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	for _, window := range w.Windows {
		if window == nil {
			continue
		}
		if window.Name != nil {
			if err := window.Name.Accept(visitor); err != nil {
				return err
			}
		}
		if window.Expr != nil {
			if err := window.Expr.Accept(visitor); err != nil {
				return err
			}
		}
	}
	return visitor.VisitWindowExpr(w)
}

type WindowFrameClause struct {
	FramePos Pos
	Type     string
	Extend   Expr
}

func (f *WindowFrameClause) Pos() Pos {
	return f.FramePos
}

func (f *WindowFrameClause) End() Pos {
	return f.Extend.End()
}

func (f *WindowFrameClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Extend.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameExpr(f)
}

type WindowFrameExtendExpr struct {
	Expr      Expr
	Direction string
	EndPos    Pos
}

func (f *WindowFrameExtendExpr) Pos() Pos {
	return f.Expr.Pos()
}

func (f *WindowFrameExtendExpr) End() Pos {
	if f.EndPos != 0 {
		return f.EndPos
	}
	return f.Expr.End()
}

func (f *WindowFrameExtendExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameExtendExpr(f)
}

type BetweenClause struct {
	Expr    Expr
	Between Expr
	AndPos  Pos
	And     Expr
}

func (f *BetweenClause) Pos() Pos {
	if f.Expr != nil {
		return f.Expr.Pos()
	}
	return f.Between.Pos()
}

func (f *BetweenClause) End() Pos {
	return f.And.End()
}

func (f *BetweenClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if f.Expr != nil {
		if err := f.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	if err := f.Between.Accept(visitor); err != nil {
		return err
	}
	if err := f.And.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitBetweenClause(f)
}

type WindowFrameCurrentRow struct {
	CurrentPos Pos
	RowEnd     Pos
}

func (f *WindowFrameCurrentRow) Pos() Pos {
	return f.CurrentPos
}

func (f *WindowFrameCurrentRow) End() Pos {
	return f.RowEnd
}

func (f *WindowFrameCurrentRow) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	return visitor.VisitWindowFrameCurrentRow(f)
}

type WindowFrameUnbounded struct {
	UnboundedPos Pos
	UnboundedEnd Pos
	Direction    string
}

func (f *WindowFrameUnbounded) Pos() Pos {
	return f.UnboundedPos
}

func (f *WindowFrameUnbounded) End() Pos {
	return f.UnboundedEnd
}

func (f *WindowFrameUnbounded) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	return visitor.VisitWindowFrameUnbounded(f)
}

type WindowFrameNumber struct {
	Number    *NumberLiteral
	EndPos    Pos
	Direction string
}

func (f *WindowFrameNumber) Pos() Pos {
	return f.Number.Pos()
}

func (f *WindowFrameNumber) End() Pos {
	return f.EndPos
}

func (f *WindowFrameNumber) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Number.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameNumber(f)
}

type WindowFrameParam struct {
	Param     *QueryParam
	EndPos    Pos
	Direction string
}

func (f *WindowFrameParam) Pos() Pos {
	return f.Param.Pos()
}

func (f *WindowFrameParam) End() Pos {
	return f.EndPos
}

func (f *WindowFrameParam) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Param.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameParam(f)
}

type SelectQuery struct {
	SelectPos     Pos
	StatementEnd  Pos
	With          *WithClause
	Top           *TopClause
	HasDistinct   bool
	DistinctOn    *DistinctOn
	SelectItems   []*SelectItem
	From          *FromClause
	Window        *WindowClause
	Prewhere      *PrewhereClause
	Where         *WhereClause
	GroupBy       *GroupByClause
	WithTotal     bool
	Having        *HavingClause
	OrderBy       *OrderByClause
	LimitBy       *LimitByClause
	Limit         *LimitClause
	Settings      *SettingsClause
	Format        *FormatClause
	UnionAll      *SelectQuery
	UnionDistinct *SelectQuery
	Except        *SelectQuery
}

func (s *SelectQuery) Pos() Pos {
	return s.SelectPos
}

func (s *SelectQuery) End() Pos {
	return s.StatementEnd
}

func (s *SelectQuery) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.With != nil {
		if err := s.With.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Top != nil {
		if err := s.Top.Accept(visitor); err != nil {
			return err
		}
	}
	if s.SelectItems != nil {
		for _, item := range s.SelectItems {
			if err := item.Accept(visitor); err != nil {
				return err
			}
		}
	}
	if s.From != nil {
		if err := s.From.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Window != nil {
		if err := s.Window.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Prewhere != nil {
		if err := s.Prewhere.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Where != nil {
		if err := s.Where.Accept(visitor); err != nil {
			return err
		}
	}
	if s.GroupBy != nil {
		if err := s.GroupBy.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Having != nil {
		if err := s.Having.Accept(visitor); err != nil {
			return err
		}
	}
	if s.OrderBy != nil {
		if err := s.OrderBy.Accept(visitor); err != nil {
			return err
		}
	}
	if s.LimitBy != nil {
		if err := s.LimitBy.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Limit != nil {
		if err := s.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Settings != nil {
		if err := s.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Format != nil {
		if err := s.Format.Accept(visitor); err != nil {
			return err
		}
	}
	if s.UnionAll != nil {
		if err := s.UnionAll.Accept(visitor); err != nil {
			return err
		}
	}
	if s.UnionDistinct != nil {
		if err := s.UnionDistinct.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Except != nil {
		if err := s.Except.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSelectQuery(s)
}

type DistinctOn struct {
	Idents        []*Ident
	DistinctOnPos Pos
	DistinctOnEnd Pos
}

func (s *DistinctOn) Pos() Pos {
	return s.DistinctOnPos
}

func (s *DistinctOn) End() Pos {
	return s.DistinctOnEnd
}

func (s *DistinctOn) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	for _, ident := range s.Idents {
		if err := ident.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDistinctOn(s)
}

type SubQuery struct {
	HasParen bool
	Select   *SelectQuery
}

func (s *SubQuery) Pos() Pos {
	return s.Select.Pos()
}

func (s *SubQuery) End() Pos {
	return s.Select.End()
}

func (s *SubQuery) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Select != nil {
		if err := s.Select.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSubQueryExpr(s)
}

type NotExpr struct {
	NotPos Pos
	Expr   Expr
}

func (n *NotExpr) Pos() Pos {
	return n.NotPos
}

func (n *NotExpr) End() Pos {
	return n.Expr.End()
}

func (n *NotExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNotExpr(n)
}

type NegateExpr struct {
	NegatePos Pos
	Expr      Expr
}

func (n *NegateExpr) Pos() Pos {
	return n.NegatePos
}

func (n *NegateExpr) End() Pos {
	return n.Expr.End()
}

func (n *NegateExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNegateExpr(n)
}

type GlobalInOperation struct {
	GlobalPos Pos
	Expr      Expr
}

func (g *GlobalInOperation) Pos() Pos {
	return g.GlobalPos
}

func (g *GlobalInOperation) End() Pos {
	return g.Expr.End()
}

func (g *GlobalInOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	if err := g.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitGlobalInExpr(g)
}

type IntervalFrom struct {
	Interval *Ident
	FromPos  Pos
	FromExpr Expr
}

func (i *IntervalFrom) Pos() Pos {
	return i.Interval.NamePos
}

func (i *IntervalFrom) End() Pos {
	return i.FromExpr.End()
}

func (i *IntervalFrom) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.FromExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIntervalFrom(i)
}

type ExtractExpr struct {
	ExtractPos Pos
	ExtractEnd Pos
	Parameters []Expr
}

func (e *ExtractExpr) Pos() Pos {
	return e.ExtractPos
}

func (e *ExtractExpr) End() Pos {
	return e.ExtractEnd
}

func (e *ExtractExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	for _, param := range e.Parameters {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitExtractExpr(e)
}

type DropDatabase struct {
	DropPos      Pos
	StatementEnd Pos
	Name         *Ident
	IfExists     bool
	OnCluster    *ClusterClause
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

type DropUserOrRole struct {
	DropPos      Pos
	Target       string
	StatementEnd Pos
	Names        []*RoleName
	IfExists     bool
	Modifier     string
	From         *Ident
}

func (d *DropUserOrRole) Pos() Pos {
	return d.DropPos
}

func (d *DropUserOrRole) End() Pos {
	return d.StatementEnd
}

func (d *DropUserOrRole) Type() string {
	return d.Target
}

func (d *DropUserOrRole) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	for _, name := range d.Names {
		if err := name.Accept(visitor); err != nil {
			return err
		}
	}
	if d.From != nil {
		if err := d.From.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDropUserOrRole(d)
}

type UseStmt struct {
	UsePos       Pos
	StatementEnd Pos
	Database     *Ident
}

func (u *UseStmt) Pos() Pos {
	return u.UsePos
}

func (u *UseStmt) End() Pos {
	return u.Database.End()
}

func (u *UseStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	if err := u.Database.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUseExpr(u)
}

type CTEStmt struct {
	CTEPos Pos
	Expr   Expr
	Alias  Expr
}

func (c *CTEStmt) Pos() Pos {
	return c.CTEPos
}

func (c *CTEStmt) End() Pos {
	return c.Expr.End()
}

func (c *CTEStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := c.Alias.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCTEExpr(c)
}

type SetStmt struct {
	SetPos   Pos
	Settings *SettingsClause
}

func (s *SetStmt) Pos() Pos {
	return s.SetPos
}

func (s *SetStmt) End() Pos {
	return s.Settings.End()
}

func (s *SetStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Settings.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSetExpr(s)
}

type FormatClause struct {
	FormatPos Pos
	Format    *Ident
}

func (f *FormatClause) Pos() Pos {
	return f.FormatPos
}

func (f *FormatClause) End() Pos {
	return f.Format.End()
}

func (f *FormatClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Format.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitFormatExpr(f)
}

type OptimizeStmt struct {
	OptimizePos  Pos
	StatementEnd Pos
	Table        *TableIdentifier
	OnCluster    *ClusterClause
	Partition    *PartitionClause
	HasFinal     bool
	Deduplicate  *DeduplicateClause
}

func (o *OptimizeStmt) Pos() Pos {
	return o.OptimizePos
}

func (o *OptimizeStmt) End() Pos {
	return o.StatementEnd
}

func (o *OptimizeStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.Table.Accept(visitor); err != nil {
		return err
	}
	if o.OnCluster != nil {
		if err := o.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Partition != nil {
		if err := o.Partition.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Deduplicate != nil {
		if err := o.Deduplicate.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitOptimizeExpr(o)
}

type DeduplicateClause struct {
	DeduplicatePos Pos
	By             *ColumnExprList
	Except         *ColumnExprList
}

func (d *DeduplicateClause) Pos() Pos {
	return d.DeduplicatePos
}

func (d *DeduplicateClause) End() Pos {
	if d.By != nil {
		return d.By.End()
	} else if d.Except != nil {
		return d.Except.End()
	}
	return d.DeduplicatePos + Pos(len(KeywordDeduplicate))
}

func (d *DeduplicateClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if d.By != nil {
		if err := d.By.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Except != nil {
		if err := d.Except.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDeduplicateExpr(d)
}

type SystemStmt struct {
	SystemPos Pos
	Expr      Expr
}

func (s *SystemStmt) Pos() Pos {
	return s.SystemPos
}

func (s *SystemStmt) End() Pos {
	return s.Expr.End()
}

func (s *SystemStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSystemExpr(s)
}

type SystemFlushExpr struct {
	FlushPos     Pos
	StatementEnd Pos
	Logs         bool
	Distributed  *TableIdentifier
}

func (s *SystemFlushExpr) Pos() Pos {
	return s.FlushPos
}

func (s *SystemFlushExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemFlushExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Distributed != nil {
		if err := s.Distributed.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSystemFlushExpr(s)
}

type SystemReloadExpr struct {
	ReloadPos    Pos
	StatementEnd Pos
	Dictionary   *TableIdentifier
	Type         string
}

func (s *SystemReloadExpr) Pos() Pos {
	return s.ReloadPos
}

func (s *SystemReloadExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemReloadExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Dictionary != nil {
		if err := s.Dictionary.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSystemReloadExpr(s)
}

type SystemSyncExpr struct {
	SyncPos Pos
	Cluster *TableIdentifier
}

func (s *SystemSyncExpr) Pos() Pos {
	return s.SyncPos
}

func (s *SystemSyncExpr) End() Pos {
	return s.Cluster.End()
}

func (s *SystemSyncExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Cluster.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSystemSyncExpr(s)
}

type SystemCtrlExpr struct {
	CtrlPos      Pos
	StatementEnd Pos
	Command      string // START, STOP
	Type         string // REPLICATED, DISTRIBUTED
	Cluster      *TableIdentifier
}

func (s *SystemCtrlExpr) Pos() Pos {
	return s.CtrlPos
}

func (s *SystemCtrlExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemCtrlExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Cluster != nil {
		if err := s.Cluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSystemCtrlExpr(s)
}

type SystemDropExpr struct {
	DropPos      Pos
	StatementEnd Pos
	Type         string
}

func (s *SystemDropExpr) Pos() Pos {
	return s.DropPos
}

func (s *SystemDropExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemDropExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	return visitor.VisitSystemDropExpr(s)
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

type SampleClause struct {
	SamplePos Pos
	Ratio     *RatioExpr
	Offset    *RatioExpr
}

func (s *SampleClause) Pos() Pos {
	return s.SamplePos
}

func (s *SampleClause) End() Pos {
	if s.Offset != nil {
		return s.Offset.End()
	}
	return s.Ratio.End()
}

func (s *SampleClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Ratio.Accept(visitor); err != nil {
		return err
	}
	if s.Offset != nil {
		if err := s.Offset.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSampleRatioExpr(s)
}

type DeleteClause struct {
	DeletePos Pos
	Table     *TableIdentifier
	OnCluster *ClusterClause
	WhereExpr Expr
}

func (d *DeleteClause) Pos() Pos {
	return d.DeletePos
}

func (d *DeleteClause) End() Pos {
	return d.WhereExpr.End()
}

func (d *DeleteClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Table.Accept(visitor); err != nil {
		return err
	}
	if d.OnCluster != nil {
		if err := d.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if d.WhereExpr != nil {
		if err := d.WhereExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDeleteFromExpr(d)
}

type ColumnNamesExpr struct {
	LeftParenPos  Pos
	RightParenPos Pos
	ColumnNames   []NestedIdentifier
}

func (c *ColumnNamesExpr) Pos() Pos {
	return c.LeftParenPos
}

func (c *ColumnNamesExpr) End() Pos {
	return c.RightParenPos
}

func (c *ColumnNamesExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for i := range c.ColumnNames {
		if err := c.ColumnNames[i].Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnNamesExpr(c)
}

type AssignmentValues struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Values        []Expr
}

func (v *AssignmentValues) Pos() Pos {
	return v.LeftParenPos
}

func (v *AssignmentValues) End() Pos {
	return v.RightParenPos
}

func (v *AssignmentValues) Accept(visitor ASTVisitor) error {
	visitor.Enter(v)
	defer visitor.Leave(v)
	for _, value := range v.Values {
		if err := value.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitValuesExpr(v)
}

type InsertStmt struct {
	InsertPos       Pos
	Format          *FormatClause
	HasTableKeyword bool
	Table           Expr
	ColumnNames     *ColumnNamesExpr
	Values          []*AssignmentValues
	SelectExpr      *SelectQuery
}

func (i *InsertStmt) Pos() Pos {
	return i.InsertPos
}

func (i *InsertStmt) End() Pos {
	if i.SelectExpr != nil {
		return i.SelectExpr.End()
	}
	return i.Values[len(i.Values)-1].End()
}

func (i *InsertStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if i.Format != nil {
		if err := i.Format.Accept(visitor); err != nil {
			return err
		}
	}
	if err := i.Table.Accept(visitor); err != nil {
		return err
	}
	if i.ColumnNames != nil {
		if err := i.ColumnNames.Accept(visitor); err != nil {
			return err
		}
	}
	for _, value := range i.Values {
		if err := value.Accept(visitor); err != nil {
			return err
		}
	}
	if i.SelectExpr != nil {
		if err := i.SelectExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitInsertExpr(i)
}

type CheckStmt struct {
	CheckPos  Pos
	Table     *TableIdentifier
	Partition *PartitionClause
}

func (c *CheckStmt) Pos() Pos {
	return c.CheckPos
}

func (c *CheckStmt) End() Pos {
	return c.Partition.End()
}

func (c *CheckStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Table.Accept(visitor); err != nil {
		return err
	}
	if c.Partition != nil {
		if err := c.Partition.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCheckExpr(c)
}

type UnaryExpr struct {
	UnaryPos Pos
	Kind     TokenKind
	Expr     Expr
}

func (n *UnaryExpr) Pos() Pos {
	return n.UnaryPos
}

func (n *UnaryExpr) End() Pos {
	return n.Expr.End()
}

func (n *UnaryExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUnaryExpr(n)
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

type ExplainStmt struct {
	ExplainPos Pos
	Type       string
	Statement  Expr
}

func (e *ExplainStmt) Pos() Pos {
	return e.ExplainPos
}

func (e *ExplainStmt) End() Pos {
	return e.Statement.End()
}

func (e *ExplainStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	if err := e.Statement.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitExplainExpr(e)
}

type PrivilegeClause struct {
	PrivilegePos Pos
	PrivilegeEnd Pos
	Keywords     []string
	Params       *ParamExprList
}

func (p *PrivilegeClause) Pos() Pos {
	return p.PrivilegePos
}

func (p *PrivilegeClause) End() Pos {
	return p.PrivilegeEnd
}

func (p *PrivilegeClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if p.Params != nil {
		if err := p.Params.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPrivilegeExpr(p)
}

type GrantPrivilegeStmt struct {
	GrantPos     Pos
	StatementEnd Pos
	OnCluster    *ClusterClause
	Privileges   []*PrivilegeClause
	On           *TableIdentifier
	To           []*Ident
	WithOptions  []string
}

func (g *GrantPrivilegeStmt) Pos() Pos {
	return g.GrantPos
}

func (g *GrantPrivilegeStmt) End() Pos {
	return g.StatementEnd
}

func (g *GrantPrivilegeStmt) Type() string {
	return "GRANT PRIVILEGE"
}

func (g *GrantPrivilegeStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	if g.OnCluster != nil {
		if err := g.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	for _, privilege := range g.Privileges {
		if err := privilege.Accept(visitor); err != nil {
			return err
		}
	}
	if err := g.On.Accept(visitor); err != nil {
		return err
	}
	for _, role := range g.To {
		if err := role.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitGrantPrivilegeExpr(g)
}

type ShowStmt struct {
	ShowPos      Pos
	StatementEnd Pos
	ShowType     string           // e.g., "CREATE TABLE", "DATABASES", "TABLES"
	Target       *TableIdentifier // for SHOW CREATE TABLE table_name

	// Optional clauses for SHOW DATABASES
	NotLike     bool           // true if NOT LIKE/ILIKE
	LikeType    string         // "LIKE" or "ILIKE", empty if not used
	LikePattern Expr           // pattern expression for LIKE/ILIKE
	Limit       Expr           // limit expression
	OutFile     *StringLiteral // filename for INTO OUTFILE
	Format      *StringLiteral // format specification
}

func (s *ShowStmt) Pos() Pos {
	return s.ShowPos
}

func (s *ShowStmt) End() Pos {
	// Find the rightmost element to determine the end position
	if s.Format != nil {
		return s.Format.End()
	}
	if s.OutFile != nil {
		return s.OutFile.End()
	}
	if s.Limit != nil {
		return s.Limit.End()
	}
	if s.LikePattern != nil {
		return s.LikePattern.End()
	}
	if s.Target != nil {
		return s.Target.End()
	}
	return s.StatementEnd
}

func (s *ShowStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Target != nil {
		if err := s.Target.Accept(visitor); err != nil {
			return err
		}
	}
	if s.LikePattern != nil {
		if err := s.LikePattern.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Limit != nil {
		if err := s.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if s.OutFile != nil {
		if err := s.OutFile.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Format != nil {
		if err := s.Format.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitShowExpr(s)
}

type DescribeStmt struct {
	DescribePos  Pos
	StatementEnd Pos
	DescribeType string // e.g., "TABLE", empty if not used
	Target       *TableIdentifier
}

func (d *DescribeStmt) Pos() Pos {
	return d.DescribePos
}

func (d *DescribeStmt) End() Pos {
	return d.Target.End()
}

func (d *DescribeStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Target.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitDescribeExpr(d)
}
