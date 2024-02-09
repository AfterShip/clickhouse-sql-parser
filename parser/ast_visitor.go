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
}

type VisitFunc func(expr Expr) error

type DefaultASTVisitor struct {
	Visit VisitFunc
}

func (v *DefaultASTVisitor) VisitOperationExpr(expr *OperationExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTernaryExpr(expr *TernaryExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitBinaryExpr(expr *BinaryExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTable(expr *AlterTable) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableDropPartition(expr *AlterTableDropPartition) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableAddColumn(expr *AlterTableAddColumn) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableAddIndex(expr *AlterTableAddIndex) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableDropColumn(expr *AlterTableDropColumn) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableDropIndex(expr *AlterTableDropIndex) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableClearColumn(expr *AlterTableClearColumn) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableClearIndex(expr *AlterTableClearIndex) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitRemovePropertyType(expr *RemovePropertyType) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTableIndex(expr *TableIndex) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitIdent(expr *Ident) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitUUID(expr *UUID) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCreateDatabase(expr *CreateDatabase) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCreateTable(expr *CreateTable) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCreateMaterializedView(expr *CreateMaterializedView) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCreateView(expr *CreateView) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCreateFunction(expr *CreateFunction) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitRoleName(expr *RoleName) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSettingPair(expr *SettingPair) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitRoleSetting(expr *RoleSetting) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCreateRole(expr *CreateRole) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterRole(expr *AlterRole) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitRoleRenamePair(expr *RoleRenamePair) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDestinationExpr(expr *DestinationExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitConstraintExpr(expr *ConstraintExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNullLiteral(expr *NullLiteral) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNotNullLiteral(expr *NotNullLiteral) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNestedIdentifier(expr *NestedIdentifier) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitColumnIdentifier(expr *ColumnIdentifier) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTableIdentifier(expr *TableIdentifier) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTableSchemaExpr(expr *TableSchemaExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTableArgListExpr(expr *TableArgListExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTableFunctionExpr(expr *TableFunctionExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOnClusterExpr(expr *OnClusterExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDefaultExpr(expr *DefaultExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPartitionExpr(expr *PartitionExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPartitionByExpr(expr *PartitionByExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPrimaryKeyExpr(expr *PrimaryKeyExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSampleByExpr(expr *SampleByExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTTLExpr(expr *TTLExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTTLExprList(expr *TTLExprList) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOrderByExpr(expr *OrderByExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOrderByListExpr(expr *OrderByListExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSettingsExpr(expr *SettingsExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSettingsExprList(expr *SettingsExprList) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitParamExprList(expr *ParamExprList) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitArrayParamList(expr *ArrayParamList) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitObjectParams(expr *ObjectParams) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitFunctionExpr(expr *FunctionExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFunctionExpr(expr *WindowFunctionExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitColumn(expr *Column) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitScalarTypeExpr(expr *ScalarTypeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPropertyTypeExpr(expr *PropertyTypeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTypeWithParamsExpr(expr *TypeWithParamsExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitComplexTypeExpr(expr *ComplexTypeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNestedTypeExpr(expr *NestedTypeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCompressionCodec(expr *CompressionCodec) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNumberLiteral(expr *NumberLiteral) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitStringLiteral(expr *StringLiteral) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitRatioExpr(expr *RatioExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitEnumValueExpr(expr *EnumValueExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitEnumValueExprList(expr *EnumValueExprList) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitIntervalExpr(expr *IntervalExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitEngineExpr(expr *EngineExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitColumnTypeExpr(expr *ColumnTypeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitColumnArgList(expr *ColumnArgList) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitColumnExprList(expr *ColumnExprList) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWhenExpr(expr *WhenExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCaseExpr(expr *CaseExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCastExpr(expr *CastExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWithExpr(expr *WithExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTopExpr(expr *TopExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCreateLiveView(expr *CreateLiveView) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWithTimeoutExpr(expr *WithTimeoutExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTableExpr(expr *TableExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOnExpr(expr *OnExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitUsingExpr(expr *UsingExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitJoinExpr(expr *JoinExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitJoinConstraintExpr(expr *JoinConstraintExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitFromExpr(expr *FromExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitIsNullExpr(expr *IsNullExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitIsNotNullExpr(expr *IsNotNullExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAliasExpr(expr *AliasExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWhereExpr(expr *WhereExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPrewhereExpr(expr *PrewhereExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitGroupByExpr(expr *GroupByExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitHavingExpr(expr *HavingExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitLimitExpr(expr *LimitExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitLimitByExpr(expr *LimitByExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowConditionExpr(expr *WindowConditionExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowExpr(expr *WindowExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFrameExpr(expr *WindowFrameExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFrameRangeExpr(expr *WindowFrameRangeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFrameNumber(expr *WindowFrameNumber) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitArrayJoinExpr(expr *ArrayJoinExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSelectQuery(expr *SelectQuery) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSubQueryExpr(expr *SubQueryExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNotExpr(expr *NotExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNegateExpr(expr *NegateExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitGlobalInExpr(expr *GlobalInExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitExtractExpr(expr *ExtractExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDropDatabase(expr *DropDatabase) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDropStmt(expr *DropStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDropUserOrRole(expr *DropUserOrRole) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitUseExpr(expr *UseExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCTEExpr(expr *CTEExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSetExpr(expr *SetExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitFormatExpr(expr *FormatExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOptimizeExpr(expr *OptimizeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDeduplicateExpr(expr *DeduplicateExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSystemExpr(expr *SystemExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSystemFlushExpr(expr *SystemFlushExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSystemReloadExpr(expr *SystemReloadExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSystemSyncExpr(expr *SystemSyncExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSystemCtrlExpr(expr *SystemCtrlExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSystemDropExpr(expr *SystemDropExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTruncateTable(expr *TruncateTable) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSampleRatioExpr(expr *SampleRatioExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDeleteFromExpr(expr *DeleteFromExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitColumnNamesExpr(expr *ColumnNamesExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitValuesExpr(expr *ValuesExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitInsertExpr(expr *InsertExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCheckExpr(expr *CheckExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitUnaryExpr(expr *UnaryExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitRenameStmt(expr *RenameStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitExplainExpr(expr *ExplainExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPrivilegeExpr(expr *PrivilegeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitGrantPrivilegeExpr(expr *GrantPrivilegeExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}
