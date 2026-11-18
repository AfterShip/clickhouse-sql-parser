package parser

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
	// Settings holds the query-level SETTINGS clause. ClickHouse only
	// accepts it once, after the last alter clause.
	Settings *SettingsClause
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
	if a.Settings != nil {
		if err := a.Settings.Accept(visitor); err != nil {
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
	return visitor.VisitAlterTableDetachPartition(a)
}

type AlterTableDropPartition struct {
	DropPos     Pos
	HasDetached bool
	Partition   *PartitionClause
}

func (a *AlterTableDropPartition) Pos() Pos {
	return a.DropPos
}

func (a *AlterTableDropPartition) End() Pos {
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
}

func (a *AlterTableAddColumn) Pos() Pos {
	return a.AddPos
}

func (a *AlterTableAddColumn) End() Pos {
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

type AlterTableModifyOrderBy struct {
	ModifyPos    Pos
	StatementEnd Pos
	OrderBy      Expr
}

func (a *AlterTableModifyOrderBy) Pos() Pos {
	return a.ModifyPos
}

func (a *AlterTableModifyOrderBy) End() Pos {
	return a.StatementEnd
}

func (a *AlterTableModifyOrderBy) AlterType() string {
	return "MODIFY_ORDER_BY"
}

func (a *AlterTableModifyOrderBy) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.OrderBy.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAlterTableModifyOrderBy(a)
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
