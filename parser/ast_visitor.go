package parser

type ASTVisitor interface {
	VisitOperationExpr(expr *OperationExpr) (Expr, error)
	VisitTernaryExpr(expr *TernaryExpr) (Expr, error)
	VisitBinaryExpr(expr *BinaryExpr) (Expr, error)
	VisitAlterTable(expr *AlterTable) (Expr, error)
	VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) (Expr, error)
	VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) (Expr, error)
	VisitAlterTableDropPartition(expr *AlterTableDropPartition) (Expr, error)
	VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) (Expr, error)
	VisitAlterTableAddColumn(expr *AlterTableAddColumn) (Expr, error)
	VisitAlterTableAddIndex(expr *AlterTableAddIndex) (Expr, error)
	VisitAlterTableDropColumn(expr *AlterTableDropColumn) (Expr, error)
	VisitAlterTableDropIndex(expr *AlterTableDropIndex) (Expr, error)
	VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) (Expr, error)
	VisitAlterTableClearColumn(expr *AlterTableClearColumn) (Expr, error)
	VisitAlterTableClearIndex(expr *AlterTableClearIndex) (Expr, error)
	VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) (Expr, error)
	VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) (Expr, error)
	VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) (Expr, error)
	VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) (Expr, error)
	VisitRemovePropertyType(expr *RemovePropertyType) (Expr, error)
	VisitTableIndex(expr *TableIndex) (Expr, error)
	VisitIdent(expr *Ident) (Expr, error)
	VisitUUID(expr *UUID) (Expr, error)
	VisitCreateDatabase(expr *CreateDatabase) (Expr, error)
	VisitCreateTable(expr *CreateTable) (Expr, error)
	VisitCreateMaterializedView(expr *CreateMaterializedView) (Expr, error)
	VisitCreateView(expr *CreateView) (Expr, error)
	VisitCreateFunction(expr *CreateFunction) (Expr, error)
	VisitRoleName(expr *RoleName) (Expr, error)
	VisitSettingPair(expr *SettingPair) (Expr, error)
	VisitRoleSetting(expr *RoleSetting) (Expr, error)
	VisitCreateRole(expr *CreateRole) (Expr, error)
	VisitAlterRole(expr *AlterRole) (Expr, error)
	VisitRoleRenamePair(expr *RoleRenamePair) (Expr, error)
	VisitDestinationExpr(expr *DestinationExpr) (Expr, error)
	VisitConstraintExpr(expr *ConstraintExpr) (Expr, error)
	VisitNullLiteral(expr *NullLiteral) (Expr, error)
	VisitNotNullLiteral(expr *NotNullLiteral) (Expr, error)
	VisitNestedIdentifier(expr *NestedIdentifier) (Expr, error)
	VisitColumnIdentifier(expr *ColumnIdentifier) (Expr, error)
	VisitTableIdentifier(expr *TableIdentifier) (Expr, error)
	VisitTableSchemaExpr(expr *TableSchemaExpr) (Expr, error)
	VisitTableArgListExpr(expr *TableArgListExpr) (Expr, error)
	VisitTableFunctionExpr(expr *TableFunctionExpr) (Expr, error)
	VisitOnClusterExpr(expr *OnClusterExpr) (Expr, error)
	VisitDefaultExpr(expr *DefaultExpr) (Expr, error)
	VisitPartitionExpr(expr *PartitionExpr) (Expr, error)
	VisitPartitionByExpr(expr *PartitionByExpr) (Expr, error)
	VisitPrimaryKeyExpr(expr *PrimaryKeyExpr) (Expr, error)
	VisitSampleByExpr(expr *SampleByExpr) (Expr, error)
	VisitTTLExpr(expr *TTLExpr) (Expr, error)
	VisitTTLExprList(expr *TTLExprList) (Expr, error)
	VisitOrderByExpr(expr *OrderByExpr) (Expr, error)
	VisitOrderByListExpr(expr *OrderByListExpr) (Expr, error)
	VisitSettingsExpr(expr *SettingsExpr) (Expr, error)
	VisitSettingsExprList(expr *SettingsExprList) (Expr, error)
	VisitParamExprList(expr *ParamExprList) (Expr, error)
	VisitArrayParamList(expr *ArrayParamList) (Expr, error)
	VisitObjectParams(expr *ObjectParams) (Expr, error)
	VisitFunctionExpr(expr *FunctionExpr) (Expr, error)
	VisitWindowFunctionExpr(expr *WindowFunctionExpr) (Expr, error)
	VisitColumn(expr *Column) (Expr, error)
	VisitScalarTypeExpr(expr *ScalarTypeExpr) (Expr, error)
	VisitPropertyTypeExpr(expr *PropertyTypeExpr) (Expr, error)
	VisitTypeWithParamsExpr(expr *TypeWithParamsExpr) (Expr, error)
	VisitComplexTypeExpr(expr *ComplexTypeExpr) (Expr, error)
	VisitNestedTypeExpr(expr *NestedTypeExpr) (Expr, error)
	VisitCompressionCodec(expr *CompressionCodec) (Expr, error)
	VisitNumberLiteral(expr *NumberLiteral) (Expr, error)
	VisitStringLiteral(expr *StringLiteral) (Expr, error)
	VisitRatioExpr(expr *RatioExpr) (Expr, error)
	VisitEnumValueExpr(expr *EnumValueExpr) (Expr, error)
	VisitEnumValueExprList(expr *EnumValueExprList) (Expr, error)
	VisitIntervalExpr(expr *IntervalExpr) (Expr, error)
	VisitEngineExpr(expr *EngineExpr) (Expr, error)
	VisitColumnTypeExpr(expr *ColumnTypeExpr) (Expr, error)
	VisitColumnArgList(expr *ColumnArgList) (Expr, error)
	VisitColumnExprList(expr *ColumnExprList) (Expr, error)
	VisitWhenExpr(expr *WhenExpr) (Expr, error)
	VisitCaseExpr(expr *CaseExpr) (Expr, error)
	VisitCastExpr(expr *CastExpr) (Expr, error)
	VisitWithExpr(expr *WithExpr) (Expr, error)
	VisitTopExpr(expr *TopExpr) (Expr, error)
	VisitCreateLiveView(expr *CreateLiveView) (Expr, error)
	VisitWithTimeoutExpr(expr *WithTimeoutExpr) (Expr, error)
	VisitTableExpr(expr *TableExpr) (Expr, error)
	VisitOnExpr(expr *OnExpr) (Expr, error)
	VisitUsingExpr(expr *UsingExpr) (Expr, error)
	VisitJoinExpr(expr *JoinExpr) (Expr, error)
	VisitJoinConstraintExpr(expr *JoinConstraintExpr) (Expr, error)
	VisitFromExpr(expr *FromExpr) (Expr, error)
	VisitIsNullExpr(expr *IsNullExpr) (Expr, error)
	VisitIsNotNullExpr(expr *IsNotNullExpr) (Expr, error)
	VisitAliasExpr(expr *AliasExpr) (Expr, error)
	VisitWhereExpr(expr *WhereExpr) (Expr, error)
	VisitPrewhereExpr(expr *PrewhereExpr) (Expr, error)
	VisitGroupByExpr(expr *GroupByExpr) (Expr, error)
	VisitHavingExpr(expr *HavingExpr) (Expr, error)
	VisitLimitExpr(expr *LimitExpr) (Expr, error)
	VisitLimitByExpr(expr *LimitByExpr) (Expr, error)
	VisitWindowConditionExpr(expr *WindowConditionExpr) (Expr, error)
	VisitWindowExpr(expr *WindowExpr) (Expr, error)
	VisitWindowFrameExpr(expr *WindowFrameExpr) (Expr, error)
	VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) (Expr, error)
	VisitWindowFrameRangeExpr(expr *WindowFrameRangeExpr) (Expr, error)
	VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) (Expr, error)
	VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) (Expr, error)
	VisitWindowFrameNumber(expr *WindowFrameNumber) (Expr, error)
	VisitArrayJoinExpr(expr *ArrayJoinExpr) (Expr, error)
	VisitSelectQuery(expr *SelectQuery) (Expr, error)
	VisitSubQueryExpr(expr *SubQueryExpr) (Expr, error)
	VisitNotExpr(expr *NotExpr) (Expr, error)
	VisitNegateExpr(expr *NegateExpr) (Expr, error)
	VisitGlobalInExpr(expr *GlobalInExpr) (Expr, error)
	VisitExtractExpr(expr *ExtractExpr) (Expr, error)
	VisitDropDatabase(expr *DropDatabase) (Expr, error)
	VisitDropStmt(expr *DropStmt) (Expr, error)
	VisitDropUserOrRole(expr *DropUserOrRole) (Expr, error)
	VisitUseExpr(expr *UseExpr) (Expr, error)
	VisitCTEExpr(expr *CTEExpr) (Expr, error)
	VisitSetExpr(expr *SetExpr) (Expr, error)
	VisitFormatExpr(expr *FormatExpr) (Expr, error)
	VisitOptimizeExpr(expr *OptimizeExpr) (Expr, error)
	VisitDeduplicateExpr(expr *DeduplicateExpr) (Expr, error)
	VisitSystemExpr(expr *SystemExpr) (Expr, error)
	VisitSystemFlushExpr(expr *SystemFlushExpr) (Expr, error)
	VisitSystemReloadExpr(expr *SystemReloadExpr) (Expr, error)
	VisitSystemSyncExpr(expr *SystemSyncExpr) (Expr, error)
	VisitSystemCtrlExpr(expr *SystemCtrlExpr) (Expr, error)
	VisitSystemDropExpr(expr *SystemDropExpr) (Expr, error)
	VisitTruncateTable(expr *TruncateTable) (Expr, error)
	VisitSampleRatioExpr(expr *SampleRatioExpr) (Expr, error)
	VisitDeleteFromExpr(expr *DeleteFromExpr) (Expr, error)
	VisitColumnNamesExpr(expr *ColumnNamesExpr) (Expr, error)
	VisitValuesExpr(expr *ValuesExpr) (Expr, error)
	VisitInsertExpr(expr *InsertExpr) (Expr, error)
	VisitCheckExpr(expr *CheckExpr) (Expr, error)
	VisitUnaryExpr(expr *UnaryExpr) (Expr, error)
	VisitRenameStmt(expr *RenameStmt) (Expr, error)
	//VisitTargetPair(expr *TargetPair) (Expr, error)
	VisitExplainExpr(expr *ExplainExpr) (Expr, error)
	VisitPrivilegeExpr(expr *PrivilegeExpr) (Expr, error)
	VisitGrantPrivilegeExpr(expr *GrantPrivilegeExpr) (Expr, error)

	enter(expr Expr)
	leave(expr Expr)
}

type VisitFunc func(expr Expr) (Expr, error)
type EnterLeaveFunc func(expr Expr)

func DefaultVisitFunc(expr Expr) (Expr, error) {
	return expr, nil
}

type defaultASTVisitor struct {
	Visit VisitFunc
	Enter EnterLeaveFunc
	Leave EnterLeaveFunc
}

func NewDefaultASTVisitor(visitFunc VisitFunc, enterFunc EnterLeaveFunc, leaveFunc EnterLeaveFunc) ASTVisitor {
	if visitFunc == nil {
		visitFunc = DefaultVisitFunc
	}
	return &defaultASTVisitor{
		Visit: visitFunc,
		Enter: enterFunc,
		Leave: leaveFunc,
	}
}

func (v *defaultASTVisitor) VisitOperationExpr(expr *OperationExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTernaryExpr(expr *TernaryExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitBinaryExpr(expr *BinaryExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTable(expr *AlterTable) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropPartition(expr *AlterTableDropPartition) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAddColumn(expr *AlterTableAddColumn) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAddIndex(expr *AlterTableAddIndex) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropColumn(expr *AlterTableDropColumn) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropIndex(expr *AlterTableDropIndex) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableClearColumn(expr *AlterTableClearColumn) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableClearIndex(expr *AlterTableClearIndex) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRemovePropertyType(expr *RemovePropertyType) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableIndex(expr *TableIndex) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIdent(expr *Ident) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUUID(expr *UUID) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateDatabase(expr *CreateDatabase) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateTable(expr *CreateTable) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateMaterializedView(expr *CreateMaterializedView) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateView(expr *CreateView) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateFunction(expr *CreateFunction) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleName(expr *RoleName) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingPair(expr *SettingPair) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleSetting(expr *RoleSetting) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateRole(expr *CreateRole) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterRole(expr *AlterRole) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleRenamePair(expr *RoleRenamePair) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDestinationExpr(expr *DestinationExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitConstraintExpr(expr *ConstraintExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNullLiteral(expr *NullLiteral) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNotNullLiteral(expr *NotNullLiteral) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNestedIdentifier(expr *NestedIdentifier) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnIdentifier(expr *ColumnIdentifier) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableIdentifier(expr *TableIdentifier) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableSchemaExpr(expr *TableSchemaExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableArgListExpr(expr *TableArgListExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableFunctionExpr(expr *TableFunctionExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOnClusterExpr(expr *OnClusterExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDefaultExpr(expr *DefaultExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPartitionExpr(expr *PartitionExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPartitionByExpr(expr *PartitionByExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrimaryKeyExpr(expr *PrimaryKeyExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSampleByExpr(expr *SampleByExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTTLExpr(expr *TTLExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTTLExprList(expr *TTLExprList) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOrderByExpr(expr *OrderByExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOrderByListExpr(expr *OrderByListExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingsExpr(expr *SettingsExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingsExprList(expr *SettingsExprList) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitParamExprList(expr *ParamExprList) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitArrayParamList(expr *ArrayParamList) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitObjectParams(expr *ObjectParams) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFunctionExpr(expr *FunctionExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFunctionExpr(expr *WindowFunctionExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumn(expr *Column) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitScalarTypeExpr(expr *ScalarTypeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPropertyTypeExpr(expr *PropertyTypeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTypeWithParamsExpr(expr *TypeWithParamsExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitComplexTypeExpr(expr *ComplexTypeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNestedTypeExpr(expr *NestedTypeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCompressionCodec(expr *CompressionCodec) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNumberLiteral(expr *NumberLiteral) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitStringLiteral(expr *StringLiteral) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRatioExpr(expr *RatioExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEnumValueExpr(expr *EnumValueExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEnumValueExprList(expr *EnumValueExprList) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIntervalExpr(expr *IntervalExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEngineExpr(expr *EngineExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnTypeExpr(expr *ColumnTypeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnArgList(expr *ColumnArgList) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnExprList(expr *ColumnExprList) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWhenExpr(expr *WhenExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCaseExpr(expr *CaseExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCastExpr(expr *CastExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWithExpr(expr *WithExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTopExpr(expr *TopExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateLiveView(expr *CreateLiveView) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWithTimeoutExpr(expr *WithTimeoutExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableExpr(expr *TableExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOnExpr(expr *OnExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUsingExpr(expr *UsingExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitJoinExpr(expr *JoinExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitJoinConstraintExpr(expr *JoinConstraintExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFromExpr(expr *FromExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIsNullExpr(expr *IsNullExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIsNotNullExpr(expr *IsNotNullExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAliasExpr(expr *AliasExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWhereExpr(expr *WhereExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrewhereExpr(expr *PrewhereExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGroupByExpr(expr *GroupByExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitHavingExpr(expr *HavingExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitLimitExpr(expr *LimitExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitLimitByExpr(expr *LimitByExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowConditionExpr(expr *WindowConditionExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowExpr(expr *WindowExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameExpr(expr *WindowFrameExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameRangeExpr(expr *WindowFrameRangeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameNumber(expr *WindowFrameNumber) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitArrayJoinExpr(expr *ArrayJoinExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSelectQuery(expr *SelectQuery) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSubQueryExpr(expr *SubQueryExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNotExpr(expr *NotExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNegateExpr(expr *NegateExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGlobalInExpr(expr *GlobalInExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitExtractExpr(expr *ExtractExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropDatabase(expr *DropDatabase) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropStmt(expr *DropStmt) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropUserOrRole(expr *DropUserOrRole) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUseExpr(expr *UseExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCTEExpr(expr *CTEExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSetExpr(expr *SetExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFormatExpr(expr *FormatExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOptimizeExpr(expr *OptimizeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDeduplicateExpr(expr *DeduplicateExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemExpr(expr *SystemExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemFlushExpr(expr *SystemFlushExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemReloadExpr(expr *SystemReloadExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemSyncExpr(expr *SystemSyncExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemCtrlExpr(expr *SystemCtrlExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemDropExpr(expr *SystemDropExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTruncateTable(expr *TruncateTable) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSampleRatioExpr(expr *SampleRatioExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDeleteFromExpr(expr *DeleteFromExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnNamesExpr(expr *ColumnNamesExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitValuesExpr(expr *ValuesExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitInsertExpr(expr *InsertExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCheckExpr(expr *CheckExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUnaryExpr(expr *UnaryExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRenameStmt(expr *RenameStmt) (Expr, error) {
	return v.Visit(expr)
}

//func (v *defaultASTVisitor) VisitTargetPair(expr *TargetPair) (Expr, error) {
//	return v.Visit(expr)
//}

func (v *defaultASTVisitor) VisitExplainExpr(expr *ExplainExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrivilegeExpr(expr *PrivilegeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGrantPrivilegeExpr(expr *GrantPrivilegeExpr) (Expr, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) enter(expr Expr) {
	if v.Enter != nil {
		v.Enter(expr)
	}
}

func (v *defaultASTVisitor) leave(expr Expr) {
	if v.Leave != nil {
		v.Leave(expr)
	}
}
