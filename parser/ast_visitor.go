package parser

//OperationExpr
//TernaryExpr
//BinaryExpr
//AlterTableExpr
//AlterTable
//AlterTableAttachPartition
//AlterTableDetachPartition
//AlterTableDropPartition
//AlterTableFreezePartition
//AlterTableAddColumn
//AlterTableAddIndex
//AlterTableDropColumn
//AlterTableDropIndex
//AlterTableRemoveTTL
//AlterTableClearColumn
//AlterTableClearIndex
//AlterTableRenameColumn
//AlterTableModifyTTL
//AlterTableModifyColumn
//AlterTableReplacePartition
//RemovePropertyType
//TableIndex
//Ident
//UUID
//CreateDatabase
//CreateTable
//CreateMaterializedView
//CreateView
//CreateFunction
//RoleName
//SettingPair
//RoleSetting
//CreateRole
//AlterRole
//RoleRenamePair
//DestinationExpr
//ConstraintExpr
//NullLiteral
//NotNullLiteral
//NestedIdentifier
//ColumnIdentifier
//TableIdentifier
//TableSchemaExpr
//TableArgListExpr
//TableFunctionExpr
//OnClusterExpr
//DefaultExpr
//PartitionExpr
//PartitionByExpr
//PrimaryKeyExpr
//SampleByExpr
//TTLExpr
//TTLExprList
//OrderByExpr
//OrderByListExpr
//SettingsExpr
//SettingsExprList
//ParamExprList
//ArrayParamList
//ObjectParams
//FunctionExpr
//WindowFunctionExpr
//Column
//ScalarTypeExpr
//PropertyTypeExpr
//TypeWithParamsExpr
//ComplexTypeExpr
//NestedTypeExpr
//CompressionCodec
//NumberLiteral
//StringLiteral
//RatioExpr
//EnumValueExpr
//EnumValueExprList
//IntervalExpr
//EngineExpr
//ColumnTypeExpr
//ColumnArgList
//ColumnExprList
//WhenExpr
//CaseExpr
//CastExpr
//WithExpr
//TopExpr
//CreateLiveView
//WithTimeoutExpr
//TableExpr
//OnExpr
//UsingExpr
//JoinExpr
//JoinConstraintExpr
//FromExpr
//IsNullExpr
//IsNotNullExpr
//AliasExpr
//WhereExpr
//PrewhereExpr
//GroupByExpr
//HavingExpr
//LimitExpr
//LimitByExpr
//WindowConditionExpr
//WindowExpr
//WindowFrameExpr
//WindowFrameExtendExpr
//WindowFrameRangeExpr
//WindowFrameCurrentRow
//WindowFrameUnbounded
//WindowFrameNumber
//ArrayJoinExpr
//SelectQuery
//SubQueryExpr
//NotExpr
//NegateExpr
//GlobalInExpr
//ExtractExpr
//DropDatabase
//DropStmt
//DropUserOrRole
//UseExpr
//CTEExpr
//SetExpr
//FormatExpr
//OptimizeExpr
//DeduplicateExpr
//SystemExpr
//SystemFlushExpr
//SystemReloadExpr
//SystemSyncExpr
//SystemCtrlExpr
//SystemDropExpr
//TruncateTable
//SampleRatioExpr
//DeleteFromExpr
//ColumnNamesExpr
//ValuesExpr
//InsertExpr
//CheckExpr
//UnaryExpr
//RenameStmt
//TargetPair
//ExplainExpr
//PrivilegeExpr
//GrantPrivilegeExpr

type ASTVisitor interface {
	VisitOperationExpr(expr *OperationExpr) (interface{}, error)
	VisitTernaryExpr(expr *TernaryExpr) (interface{}, error)
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitAlterTable(expr *AlterTable) (interface{}, error)
	VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) (interface{}, error)
	VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) (interface{}, error)
	VisitAlterTableDropPartition(expr *AlterTableDropPartition) (interface{}, error)
	VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) (interface{}, error)
	VisitAlterTableAddColumn(expr *AlterTableAddColumn) (interface{}, error)
	VisitAlterTableAddIndex(expr *AlterTableAddIndex) (interface{}, error)
	VisitAlterTableDropColumn(expr *AlterTableDropColumn) (interface{}, error)
	VisitAlterTableDropIndex(expr *AlterTableDropIndex) (interface{}, error)
	VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) (interface{}, error)
	VisitAlterTableClearColumn(expr *AlterTableClearColumn) (interface{}, error)
	VisitAlterTableClearIndex(expr *AlterTableClearIndex) (interface{}, error)
	VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) (interface{}, error)
	VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) (interface{}, error)
	VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) (interface{}, error)
	VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) (interface{}, error)
	VisitRemovePropertyType(expr *RemovePropertyType) (interface{}, error)
	VisitTableIndex(expr *TableIndex) (interface{}, error)
	VisitIdent(expr *Ident) (interface{}, error)
	VisitUUID(expr *UUID) (interface{}, error)
	VisitCreateDatabase(expr *CreateDatabase) (interface{}, error)
	VisitCreateTable(expr *CreateTable) (interface{}, error)
	VisitCreateMaterializedView(expr *CreateMaterializedView) (interface{}, error)
	VisitCreateView(expr *CreateView) (interface{}, error)
	VisitCreateFunction(expr *CreateFunction) (interface{}, error)
	VisitRoleName(expr *RoleName) (interface{}, error)
	VisitSettingPair(expr *SettingPair) (interface{}, error)
	VisitRoleSetting(expr *RoleSetting) (interface{}, error)
	VisitCreateRole(expr *CreateRole) (interface{}, error)
	VisitAlterRole(expr *AlterRole) (interface{}, error)
	VisitRoleRenamePair(expr *RoleRenamePair) (interface{}, error)
	VisitDestinationExpr(expr *DestinationExpr) (interface{}, error)
	VisitConstraintExpr(expr *ConstraintExpr) (interface{}, error)
	VisitNullLiteral(expr *NullLiteral) (interface{}, error)
	VisitNotNullLiteral(expr *NotNullLiteral) (interface{}, error)
	VisitNestedIdentifier(expr *NestedIdentifier) (interface{}, error)
	VisitColumnIdentifier(expr *ColumnIdentifier) (interface{}, error)
	VisitTableIdentifier(expr *TableIdentifier) (interface{}, error)
	VisitTableSchemaExpr(expr *TableSchemaExpr) (interface{}, error)
	VisitTableArgListExpr(expr *TableArgListExpr) (interface{}, error)
	VisitTableFunctionExpr(expr *TableFunctionExpr) (interface{}, error)
	VisitOnClusterExpr(expr *OnClusterExpr) (interface{}, error)
	VisitDefaultExpr(expr *DefaultExpr) (interface{}, error)
	VisitPartitionExpr(expr *PartitionExpr) (interface{}, error)
	VisitPartitionByExpr(expr *PartitionByExpr) (interface{}, error)
	VisitPrimaryKeyExpr(expr *PrimaryKeyExpr) (interface{}, error)
	VisitSampleByExpr(expr *SampleByExpr) (interface{}, error)
	VisitTTLExpr(expr *TTLExpr) (interface{}, error)
	VisitTTLExprList(expr *TTLExprList) (interface{}, error)
	VisitOrderByExpr(expr *OrderByExpr) (interface{}, error)
	VisitOrderByListExpr(expr *OrderByListExpr) (interface{}, error)
	VisitSettingsExpr(expr *SettingsExpr) (interface{}, error)
	VisitSettingsExprList(expr *SettingsExprList) (interface{}, error)
	VisitParamExprList(expr *ParamExprList) (interface{}, error)
	VisitArrayParamList(expr *ArrayParamList) (interface{}, error)
	VisitObjectParams(expr *ObjectParams) (interface{}, error)
	VisitFunctionExpr(expr *FunctionExpr) (interface{}, error)
	VisitWindowFunctionExpr(expr *WindowFunctionExpr) (interface{}, error)
	VisitColumn(expr *Column) (interface{}, error)
	VisitScalarTypeExpr(expr *ScalarTypeExpr) (interface{}, error)
	VisitPropertyTypeExpr(expr *PropertyTypeExpr) (interface{}, error)
	VisitTypeWithParamsExpr(expr *TypeWithParamsExpr) (interface{}, error)
	VisitComplexTypeExpr(expr *ComplexTypeExpr) (interface{}, error)
	VisitNestedTypeExpr(expr *NestedTypeExpr) (interface{}, error)
	VisitCompressionCodec(expr *CompressionCodec) (interface{}, error)
	VisitNumberLiteral(expr *NumberLiteral) (interface{}, error)
	VisitStringLiteral(expr *StringLiteral) (interface{}, error)
	VisitRatioExpr(expr *RatioExpr) (interface{}, error)
	VisitEnumValueExpr(expr *EnumValueExpr) (interface{}, error)
	VisitEnumValueExprList(expr *EnumValueExprList) (interface{}, error)
	VisitIntervalExpr(expr *IntervalExpr) (interface{}, error)
	VisitEngineExpr(expr *EngineExpr) (interface{}, error)
	VisitColumnTypeExpr(expr *ColumnTypeExpr) (interface{}, error)
	VisitColumnArgList(expr *ColumnArgList) (interface{}, error)
	VisitColumnExprList(expr *ColumnExprList) (interface{}, error)
	VisitWhenExpr(expr *WhenExpr) (interface{}, error)
	VisitCaseExpr(expr *CaseExpr) (interface{}, error)
	VisitCastExpr(expr *CastExpr) (interface{}, error)
	VisitWithExpr(expr *WithExpr) (interface{}, error)
	VisitTopExpr(expr *TopExpr) (interface{}, error)
	VisitCreateLiveView(expr *CreateLiveView) (interface{}, error)
	VisitWithTimeoutExpr(expr *WithTimeoutExpr) (interface{}, error)
	VisitTableExpr(expr *TableExpr) (interface{}, error)
	VisitOnExpr(expr *OnExpr) (interface{}, error)
	VisitUsingExpr(expr *UsingExpr) (interface{}, error)
	VisitJoinExpr(expr *JoinExpr) (interface{}, error)
	VisitJoinConstraintExpr(expr *JoinConstraintExpr) (interface{}, error)
	VisitFromExpr(expr *FromExpr) (interface{}, error)
	VisitIsNullExpr(expr *IsNullExpr) (interface{}, error)
	VisitIsNotNullExpr(expr *IsNotNullExpr) (interface{}, error)
	VisitAliasExpr(expr *AliasExpr) (interface{}, error)
	VisitWhereExpr(expr *WhereExpr) (interface{}, error)
	VisitPrewhereExpr(expr *PrewhereExpr) (interface{}, error)
	VisitGroupByExpr(expr *GroupByExpr) (interface{}, error)
	VisitHavingExpr(expr *HavingExpr) (interface{}, error)
	VisitLimitExpr(expr *LimitExpr) (interface{}, error)
	VisitLimitByExpr(expr *LimitByExpr) (interface{}, error)
	VisitWindowConditionExpr(expr *WindowConditionExpr) (interface{}, error)
	VisitWindowExpr(expr *WindowExpr) (interface{}, error)
	VisitWindowFrameExpr(expr *WindowFrameExpr) (interface{}, error)
	VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) (interface{}, error)
	VisitWindowFrameRangeExpr(expr *WindowFrameRangeExpr) (interface{}, error)
	VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) (interface{}, error)
	VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) (interface{}, error)
	VisitWindowFrameNumber(expr *WindowFrameNumber) (interface{}, error)
	VisitArrayJoinExpr(expr *ArrayJoinExpr) (interface{}, error)
	VisitSelectQuery(expr *SelectQuery) (interface{}, error)
	VisitSubQueryExpr(expr *SubQueryExpr) (interface{}, error)
	VisitNotExpr(expr *NotExpr) (interface{}, error)
	VisitNegateExpr(expr *NegateExpr) (interface{}, error)
	VisitGlobalInExpr(expr *GlobalInExpr) (interface{}, error)
	VisitExtractExpr(expr *ExtractExpr) (interface{}, error)
	VisitDropDatabase(expr *DropDatabase) (interface{}, error)
	VisitDropStmt(expr *DropStmt) (interface{}, error)
	VisitDropUserOrRole(expr *DropUserOrRole) (interface{}, error)
	VisitUseExpr(expr *UseExpr) (interface{}, error)
	VisitCTEExpr(expr *CTEExpr) (interface{}, error)
	VisitSetExpr(expr *SetExpr) (interface{}, error)
	VisitFormatExpr(expr *FormatExpr) (interface{}, error)
	VisitOptimizeExpr(expr *OptimizeExpr) (interface{}, error)
	VisitDeduplicateExpr(expr *DeduplicateExpr) (interface{}, error)
	VisitSystemExpr(expr *SystemExpr) (interface{}, error)
	VisitSystemFlushExpr(expr *SystemFlushExpr) (interface{}, error)
	VisitSystemReloadExpr(expr *SystemReloadExpr) (interface{}, error)
	VisitSystemSyncExpr(expr *SystemSyncExpr) (interface{}, error)
	VisitSystemCtrlExpr(expr *SystemCtrlExpr) (interface{}, error)
	VisitSystemDropExpr(expr *SystemDropExpr) (interface{}, error)
	VisitTruncateTable(expr *TruncateTable) (interface{}, error)
	VisitSampleRatioExpr(expr *SampleRatioExpr) (interface{}, error)
	VisitDeleteFromExpr(expr *DeleteFromExpr) (interface{}, error)
	VisitColumnNamesExpr(expr *ColumnNamesExpr) (interface{}, error)
	VisitValuesExpr(expr *ValuesExpr) (interface{}, error)
	VisitInsertExpr(expr *InsertExpr) (interface{}, error)
	VisitCheckExpr(expr *CheckExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitRenameStmt(expr *RenameStmt) (interface{}, error)
	//VisitTargetPair(expr *TargetPair) (interface{}, error)
	VisitExplainExpr(expr *ExplainExpr) (interface{}, error)
	VisitPrivilegeExpr(expr *PrivilegeExpr) (interface{}, error)
	VisitGrantPrivilegeExpr(expr *GrantPrivilegeExpr) (interface{}, error)

	enter(expr Expr)
	leave(expr Expr)
}

type VisitFunc func(expr Expr) (interface{}, error)
type EnterLeaveFunc func(expr Expr)

func DefaultVisitFunc(expr Expr) (interface{}, error) {
	return nil, nil
}

type defaultASTVisitor struct {
	Visit VisitFunc
	Enter EnterLeaveFunc
	Leave EnterLeaveFunc
}

func NewDefaultASTVisitor(visitFunc VisitFunc, enterFunc EnterLeaveFunc, leaveFunc EnterLeaveFunc) ASTVisitor {
	return &defaultASTVisitor{
		Visit: DefaultVisitFunc,
		Enter: nil,
		Leave: nil,
	}
}

func (v *defaultASTVisitor) VisitOperationExpr(expr *OperationExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTernaryExpr(expr *TernaryExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTable(expr *AlterTable) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropPartition(expr *AlterTableDropPartition) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAddColumn(expr *AlterTableAddColumn) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableAddIndex(expr *AlterTableAddIndex) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropColumn(expr *AlterTableDropColumn) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableDropIndex(expr *AlterTableDropIndex) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableClearColumn(expr *AlterTableClearColumn) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableClearIndex(expr *AlterTableClearIndex) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRemovePropertyType(expr *RemovePropertyType) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableIndex(expr *TableIndex) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIdent(expr *Ident) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUUID(expr *UUID) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateDatabase(expr *CreateDatabase) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateTable(expr *CreateTable) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateMaterializedView(expr *CreateMaterializedView) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateView(expr *CreateView) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateFunction(expr *CreateFunction) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleName(expr *RoleName) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingPair(expr *SettingPair) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleSetting(expr *RoleSetting) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateRole(expr *CreateRole) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAlterRole(expr *AlterRole) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRoleRenamePair(expr *RoleRenamePair) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDestinationExpr(expr *DestinationExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitConstraintExpr(expr *ConstraintExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNullLiteral(expr *NullLiteral) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNotNullLiteral(expr *NotNullLiteral) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNestedIdentifier(expr *NestedIdentifier) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnIdentifier(expr *ColumnIdentifier) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableIdentifier(expr *TableIdentifier) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableSchemaExpr(expr *TableSchemaExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableArgListExpr(expr *TableArgListExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableFunctionExpr(expr *TableFunctionExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOnClusterExpr(expr *OnClusterExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDefaultExpr(expr *DefaultExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPartitionExpr(expr *PartitionExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPartitionByExpr(expr *PartitionByExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrimaryKeyExpr(expr *PrimaryKeyExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSampleByExpr(expr *SampleByExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTTLExpr(expr *TTLExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTTLExprList(expr *TTLExprList) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOrderByExpr(expr *OrderByExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOrderByListExpr(expr *OrderByListExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingsExpr(expr *SettingsExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSettingsExprList(expr *SettingsExprList) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitParamExprList(expr *ParamExprList) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitArrayParamList(expr *ArrayParamList) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitObjectParams(expr *ObjectParams) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFunctionExpr(expr *FunctionExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFunctionExpr(expr *WindowFunctionExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumn(expr *Column) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitScalarTypeExpr(expr *ScalarTypeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPropertyTypeExpr(expr *PropertyTypeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTypeWithParamsExpr(expr *TypeWithParamsExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitComplexTypeExpr(expr *ComplexTypeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNestedTypeExpr(expr *NestedTypeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCompressionCodec(expr *CompressionCodec) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNumberLiteral(expr *NumberLiteral) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitStringLiteral(expr *StringLiteral) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRatioExpr(expr *RatioExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEnumValueExpr(expr *EnumValueExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEnumValueExprList(expr *EnumValueExprList) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIntervalExpr(expr *IntervalExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitEngineExpr(expr *EngineExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnTypeExpr(expr *ColumnTypeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnArgList(expr *ColumnArgList) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnExprList(expr *ColumnExprList) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWhenExpr(expr *WhenExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCaseExpr(expr *CaseExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCastExpr(expr *CastExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWithExpr(expr *WithExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTopExpr(expr *TopExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCreateLiveView(expr *CreateLiveView) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWithTimeoutExpr(expr *WithTimeoutExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTableExpr(expr *TableExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOnExpr(expr *OnExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUsingExpr(expr *UsingExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitJoinExpr(expr *JoinExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitJoinConstraintExpr(expr *JoinConstraintExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFromExpr(expr *FromExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIsNullExpr(expr *IsNullExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitIsNotNullExpr(expr *IsNotNullExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitAliasExpr(expr *AliasExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWhereExpr(expr *WhereExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrewhereExpr(expr *PrewhereExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGroupByExpr(expr *GroupByExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitHavingExpr(expr *HavingExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitLimitExpr(expr *LimitExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitLimitByExpr(expr *LimitByExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowConditionExpr(expr *WindowConditionExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowExpr(expr *WindowExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameExpr(expr *WindowFrameExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameRangeExpr(expr *WindowFrameRangeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitWindowFrameNumber(expr *WindowFrameNumber) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitArrayJoinExpr(expr *ArrayJoinExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSelectQuery(expr *SelectQuery) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSubQueryExpr(expr *SubQueryExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNotExpr(expr *NotExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitNegateExpr(expr *NegateExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGlobalInExpr(expr *GlobalInExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitExtractExpr(expr *ExtractExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropDatabase(expr *DropDatabase) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropStmt(expr *DropStmt) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDropUserOrRole(expr *DropUserOrRole) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUseExpr(expr *UseExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCTEExpr(expr *CTEExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSetExpr(expr *SetExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitFormatExpr(expr *FormatExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitOptimizeExpr(expr *OptimizeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDeduplicateExpr(expr *DeduplicateExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemExpr(expr *SystemExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemFlushExpr(expr *SystemFlushExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemReloadExpr(expr *SystemReloadExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemSyncExpr(expr *SystemSyncExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemCtrlExpr(expr *SystemCtrlExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSystemDropExpr(expr *SystemDropExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitTruncateTable(expr *TruncateTable) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitSampleRatioExpr(expr *SampleRatioExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitDeleteFromExpr(expr *DeleteFromExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitColumnNamesExpr(expr *ColumnNamesExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitValuesExpr(expr *ValuesExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitInsertExpr(expr *InsertExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitCheckExpr(expr *CheckExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitRenameStmt(expr *RenameStmt) (interface{}, error) {
	return v.Visit(expr)
}

//func (v *defaultASTVisitor) VisitTargetPair(expr *TargetPair) (interface{}, error) {
//	return v.Visit(expr)
//}

func (v *defaultASTVisitor) VisitExplainExpr(expr *ExplainExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitPrivilegeExpr(expr *PrivilegeExpr) (interface{}, error) {
	return v.Visit(expr)
}

func (v *defaultASTVisitor) VisitGrantPrivilegeExpr(expr *GrantPrivilegeExpr) (interface{}, error) {
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
