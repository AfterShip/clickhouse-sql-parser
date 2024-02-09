package parser

type ASTVisitor interface {
	VisitOperationExpr(expr *OperationExpr) error
	VisitTernaryExpr(expr *TernaryExpr) error
	VisitBinaryExpr(expr *BinaryExpr) error
	VisitAlterTable(expr *AlterTable) error
	VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) error
	VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) error
	VisitAlterTableDropPartition(expr *AlterTableDropPartition) error
	VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) error
	VisitAlterTableAddColumn(expr *AlterTableAddColumn) error
	VisitAlterTableAddIndex(expr *AlterTableAddIndex) error
	VisitAlterTableDropColumn(expr *AlterTableDropColumn) error
	VisitAlterTableDropIndex(expr *AlterTableDropIndex) error
	VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) error
	VisitAlterTableClearColumn(expr *AlterTableClearColumn) error
	VisitAlterTableClearIndex(expr *AlterTableClearIndex) error
	VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) error
	VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) error
	VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) error
	VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) error
	VisitRemovePropertyType(expr *RemovePropertyType) error
	VisitTableIndex(expr *TableIndex) error
	VisitIdent(expr *Ident) error
	VisitUUID(expr *UUID) error
	VisitCreateDatabase(expr *CreateDatabase) error
	VisitCreateTable(expr *CreateTable) error
	VisitCreateMaterializedView(expr *CreateMaterializedView) error
	VisitCreateView(expr *CreateView) error
	VisitCreateFunction(expr *CreateFunction) error
	VisitRoleName(expr *RoleName) error
	VisitSettingPair(expr *SettingPair) error
	VisitRoleSetting(expr *RoleSetting) error
	VisitCreateRole(expr *CreateRole) error
	VisitAlterRole(expr *AlterRole) error
	VisitRoleRenamePair(expr *RoleRenamePair) error
	VisitDestinationExpr(expr *DestinationExpr) error
	VisitConstraintExpr(expr *ConstraintExpr) error
	VisitNullLiteral(expr *NullLiteral) error
	VisitNotNullLiteral(expr *NotNullLiteral) error
	VisitNestedIdentifier(expr *NestedIdentifier) error
	VisitColumnIdentifier(expr *ColumnIdentifier) error
	VisitTableIdentifier(expr *TableIdentifier) error
	VisitTableSchemaExpr(expr *TableSchemaExpr) error
	VisitTableArgListExpr(expr *TableArgListExpr) error
	VisitTableFunctionExpr(expr *TableFunctionExpr) error
	VisitOnClusterExpr(expr *OnClusterExpr) error
	VisitDefaultExpr(expr *DefaultExpr) error
	VisitPartitionExpr(expr *PartitionExpr) error
	VisitPartitionByExpr(expr *PartitionByExpr) error
	VisitPrimaryKeyExpr(expr *PrimaryKeyExpr) error
	VisitSampleByExpr(expr *SampleByExpr) error
	VisitTTLExpr(expr *TTLExpr) error
	VisitTTLExprList(expr *TTLExprList) error
	VisitOrderByExpr(expr *OrderByExpr) error
	VisitOrderByListExpr(expr *OrderByListExpr) error
	VisitSettingsExpr(expr *SettingsExpr) error
	VisitSettingsExprList(expr *SettingsExprList) error
	VisitParamExprList(expr *ParamExprList) error
	VisitArrayParamList(expr *ArrayParamList) error
	VisitObjectParams(expr *ObjectParams) error
	VisitFunctionExpr(expr *FunctionExpr) error
	VisitWindowFunctionExpr(expr *WindowFunctionExpr) error
	VisitColumn(expr *Column) error
	VisitScalarTypeExpr(expr *ScalarTypeExpr) error
	VisitPropertyTypeExpr(expr *PropertyTypeExpr) error
	VisitTypeWithParamsExpr(expr *TypeWithParamsExpr) error
	VisitComplexTypeExpr(expr *ComplexTypeExpr) error
	VisitNestedTypeExpr(expr *NestedTypeExpr) error
	VisitCompressionCodec(expr *CompressionCodec) error
	VisitNumberLiteral(expr *NumberLiteral) error
	VisitStringLiteral(expr *StringLiteral) error
	VisitRatioExpr(expr *RatioExpr) error
	VisitEnumValueExpr(expr *EnumValueExpr) error
	VisitEnumValueExprList(expr *EnumValueExprList) error
	VisitIntervalExpr(expr *IntervalExpr) error
	VisitEngineExpr(expr *EngineExpr) error
	VisitColumnTypeExpr(expr *ColumnTypeExpr) error
	VisitColumnArgList(expr *ColumnArgList) error
	VisitColumnExprList(expr *ColumnExprList) error
	VisitWhenExpr(expr *WhenExpr) error
	VisitCaseExpr(expr *CaseExpr) error
	VisitCastExpr(expr *CastExpr) error
	VisitWithExpr(expr *WithExpr) error
	VisitTopExpr(expr *TopExpr) error
	VisitCreateLiveView(expr *CreateLiveView) error
	VisitWithTimeoutExpr(expr *WithTimeoutExpr) error
	VisitTableExpr(expr *TableExpr) error
	VisitOnExpr(expr *OnExpr) error
	VisitUsingExpr(expr *UsingExpr) error
	VisitJoinExpr(expr *JoinExpr) error
	VisitJoinConstraintExpr(expr *JoinConstraintExpr) error
	VisitFromExpr(expr *FromExpr) error
	VisitIsNullExpr(expr *IsNullExpr) error
	VisitIsNotNullExpr(expr *IsNotNullExpr) error
	VisitAliasExpr(expr *AliasExpr) error
	VisitWhereExpr(expr *WhereExpr) error
	VisitPrewhereExpr(expr *PrewhereExpr) error
	VisitGroupByExpr(expr *GroupByExpr) error
	VisitHavingExpr(expr *HavingExpr) error
	VisitLimitExpr(expr *LimitExpr) error
	VisitLimitByExpr(expr *LimitByExpr) error
	VisitWindowConditionExpr(expr *WindowConditionExpr) error
	VisitWindowExpr(expr *WindowExpr) error
	VisitWindowFrameExpr(expr *WindowFrameExpr) error
	VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) error
	VisitWindowFrameRangeExpr(expr *WindowFrameRangeExpr) error
	VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) error
	VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) error
	VisitWindowFrameNumber(expr *WindowFrameNumber) error
	VisitArrayJoinExpr(expr *ArrayJoinExpr) error
	VisitSelectQuery(expr *SelectQuery) error
	VisitSubQueryExpr(expr *SubQueryExpr) error
	VisitNotExpr(expr *NotExpr) error
	VisitNegateExpr(expr *NegateExpr) error
	VisitGlobalInExpr(expr *GlobalInExpr) error
	VisitExtractExpr(expr *ExtractExpr) error
	VisitDropDatabase(expr *DropDatabase) error
	VisitDropStmt(expr *DropStmt) error
	VisitDropUserOrRole(expr *DropUserOrRole) error
	VisitUseExpr(expr *UseExpr) error
	VisitCTEExpr(expr *CTEExpr) error
	VisitSetExpr(expr *SetExpr) error
	VisitFormatExpr(expr *FormatExpr) error
	VisitOptimizeExpr(expr *OptimizeExpr) error
	VisitDeduplicateExpr(expr *DeduplicateExpr) error
	VisitSystemExpr(expr *SystemExpr) error
	VisitSystemFlushExpr(expr *SystemFlushExpr) error
	VisitSystemReloadExpr(expr *SystemReloadExpr) error
	VisitSystemSyncExpr(expr *SystemSyncExpr) error
	VisitSystemCtrlExpr(expr *SystemCtrlExpr) error
	VisitSystemDropExpr(expr *SystemDropExpr) error
	VisitTruncateTable(expr *TruncateTable) error
	VisitSampleRatioExpr(expr *SampleRatioExpr) error
	VisitDeleteFromExpr(expr *DeleteFromExpr) error
	VisitColumnNamesExpr(expr *ColumnNamesExpr) error
	VisitValuesExpr(expr *ValuesExpr) error
	VisitInsertExpr(expr *InsertExpr) error
	VisitCheckExpr(expr *CheckExpr) error
	VisitUnaryExpr(expr *UnaryExpr) error
	VisitRenameStmt(expr *RenameStmt) error
	//VisitTargetPair(expr *TargetPair) (error)
	VisitExplainExpr(expr *ExplainExpr) error
	VisitPrivilegeExpr(expr *PrivilegeExpr) error
	VisitGrantPrivilegeExpr(expr *GrantPrivilegeExpr) error

	enter(expr Expr)
	leave(expr Expr)
}

type VisitFunc func(expr Expr) error
type EnterLeaveFunc func(expr Expr)

func DefaultVisitFunc(expr Expr) error {
	return nil
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

func (v *defaultASTVisitor) VisitOperationExpr(expr *OperationExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTernaryExpr(expr *TernaryExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitBinaryExpr(expr *BinaryExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTable(expr *AlterTable) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropPartition(expr *AlterTableDropPartition) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAddColumn(expr *AlterTableAddColumn) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAddIndex(expr *AlterTableAddIndex) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropColumn(expr *AlterTableDropColumn) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropIndex(expr *AlterTableDropIndex) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableClearColumn(expr *AlterTableClearColumn) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableClearIndex(expr *AlterTableClearIndex) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRemovePropertyType(expr *RemovePropertyType) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableIndex(expr *TableIndex) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIdent(expr *Ident) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUUID(expr *UUID) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateDatabase(expr *CreateDatabase) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateTable(expr *CreateTable) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateMaterializedView(expr *CreateMaterializedView) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateView(expr *CreateView) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateFunction(expr *CreateFunction) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleName(expr *RoleName) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingPair(expr *SettingPair) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleSetting(expr *RoleSetting) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateRole(expr *CreateRole) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterRole(expr *AlterRole) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleRenamePair(expr *RoleRenamePair) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDestinationExpr(expr *DestinationExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitConstraintExpr(expr *ConstraintExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNullLiteral(expr *NullLiteral) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNotNullLiteral(expr *NotNullLiteral) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNestedIdentifier(expr *NestedIdentifier) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnIdentifier(expr *ColumnIdentifier) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableIdentifier(expr *TableIdentifier) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableSchemaExpr(expr *TableSchemaExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableArgListExpr(expr *TableArgListExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableFunctionExpr(expr *TableFunctionExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOnClusterExpr(expr *OnClusterExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDefaultExpr(expr *DefaultExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPartitionExpr(expr *PartitionExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPartitionByExpr(expr *PartitionByExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrimaryKeyExpr(expr *PrimaryKeyExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSampleByExpr(expr *SampleByExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTTLExpr(expr *TTLExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTTLExprList(expr *TTLExprList) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOrderByExpr(expr *OrderByExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOrderByListExpr(expr *OrderByListExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingsExpr(expr *SettingsExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingsExprList(expr *SettingsExprList) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitParamExprList(expr *ParamExprList) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitArrayParamList(expr *ArrayParamList) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitObjectParams(expr *ObjectParams) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFunctionExpr(expr *FunctionExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFunctionExpr(expr *WindowFunctionExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumn(expr *Column) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitScalarTypeExpr(expr *ScalarTypeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPropertyTypeExpr(expr *PropertyTypeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTypeWithParamsExpr(expr *TypeWithParamsExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitComplexTypeExpr(expr *ComplexTypeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNestedTypeExpr(expr *NestedTypeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCompressionCodec(expr *CompressionCodec) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNumberLiteral(expr *NumberLiteral) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitStringLiteral(expr *StringLiteral) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRatioExpr(expr *RatioExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEnumValueExpr(expr *EnumValueExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEnumValueExprList(expr *EnumValueExprList) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIntervalExpr(expr *IntervalExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEngineExpr(expr *EngineExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnTypeExpr(expr *ColumnTypeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnArgList(expr *ColumnArgList) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnExprList(expr *ColumnExprList) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWhenExpr(expr *WhenExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCaseExpr(expr *CaseExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCastExpr(expr *CastExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWithExpr(expr *WithExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTopExpr(expr *TopExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateLiveView(expr *CreateLiveView) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWithTimeoutExpr(expr *WithTimeoutExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableExpr(expr *TableExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOnExpr(expr *OnExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUsingExpr(expr *UsingExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitJoinExpr(expr *JoinExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitJoinConstraintExpr(expr *JoinConstraintExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFromExpr(expr *FromExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIsNullExpr(expr *IsNullExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIsNotNullExpr(expr *IsNotNullExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAliasExpr(expr *AliasExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWhereExpr(expr *WhereExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrewhereExpr(expr *PrewhereExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGroupByExpr(expr *GroupByExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitHavingExpr(expr *HavingExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitLimitExpr(expr *LimitExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitLimitByExpr(expr *LimitByExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowConditionExpr(expr *WindowConditionExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowExpr(expr *WindowExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameExpr(expr *WindowFrameExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameRangeExpr(expr *WindowFrameRangeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameNumber(expr *WindowFrameNumber) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitArrayJoinExpr(expr *ArrayJoinExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSelectQuery(expr *SelectQuery) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSubQueryExpr(expr *SubQueryExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNotExpr(expr *NotExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNegateExpr(expr *NegateExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGlobalInExpr(expr *GlobalInExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitExtractExpr(expr *ExtractExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropDatabase(expr *DropDatabase) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropStmt(expr *DropStmt) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropUserOrRole(expr *DropUserOrRole) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUseExpr(expr *UseExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCTEExpr(expr *CTEExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSetExpr(expr *SetExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFormatExpr(expr *FormatExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOptimizeExpr(expr *OptimizeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDeduplicateExpr(expr *DeduplicateExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemExpr(expr *SystemExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemFlushExpr(expr *SystemFlushExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemReloadExpr(expr *SystemReloadExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemSyncExpr(expr *SystemSyncExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemCtrlExpr(expr *SystemCtrlExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemDropExpr(expr *SystemDropExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTruncateTable(expr *TruncateTable) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSampleRatioExpr(expr *SampleRatioExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDeleteFromExpr(expr *DeleteFromExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnNamesExpr(expr *ColumnNamesExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitValuesExpr(expr *ValuesExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitInsertExpr(expr *InsertExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCheckExpr(expr *CheckExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUnaryExpr(expr *UnaryExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRenameStmt(expr *RenameStmt) error {
	return v.Visit(expr)
}

//func (v *defaultASTVisitor) VisitTargetPair(expr *TargetPair) (error) {
//	return v.Visit(expr)
//}

func (v *defaultASTVisitor) VisitExplainExpr(expr *ExplainExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrivilegeExpr(expr *PrivilegeExpr) error {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGrantPrivilegeExpr(expr *GrantPrivilegeExpr) error {
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
