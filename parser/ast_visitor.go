package parser

type ASTVisitor interface {
	VisitOperationExpr(expr *OperationExpr) error
	VisitTernaryExpr(expr *TernaryOperation) error
	VisitBinaryExpr(expr *BinaryOperation) error
	VisitIndexOperation(expr *IndexOperation) error
	VisitAlterTable(expr *AlterTable) error
	VisitAlterTableAttachPartition(expr *AlterTableAttachPartition) error
	VisitAlterTableDetachPartition(expr *AlterTableDetachPartition) error
	VisitAlterTableDropPartition(expr *AlterTableDropPartition) error
	VisitAlterTableFreezePartition(expr *AlterTableFreezePartition) error
	VisitAlterTableAddColumn(expr *AlterTableAddColumn) error
	VisitAlterTableAddIndex(expr *AlterTableAddIndex) error
	VisitAlterTableAddProjection(expr *AlterTableAddProjection) error
	VisitTableProjection(expr *TableProjection) error
	VisitProjectionOrderBy(expr *ProjectionOrderByClause) error
	VisitProjectionSelect(expr *ProjectionSelectStmt) error
	VisitAlterTableDropColumn(expr *AlterTableDropColumn) error
	VisitAlterTableDropIndex(expr *AlterTableDropIndex) error
	VisitAlterTableDropProjection(expr *AlterTableDropProjection) error
	VisitAlterTableRemoveTTL(expr *AlterTableRemoveTTL) error
	VisitAlterTableClearColumn(expr *AlterTableClearColumn) error
	VisitAlterTableClearIndex(expr *AlterTableClearIndex) error
	VisitAlterTableClearProjection(expr *AlterTableClearProjection) error
	VisitAlterTableMaterializeIndex(expr *AlterTableMaterializeIndex) error
	VisitAlterTableMaterializeProjection(expr *AlterTableMaterializeProjection) error
	VisitAlterTableRenameColumn(expr *AlterTableRenameColumn) error
	VisitAlterTableModifyTTL(expr *AlterTableModifyTTL) error
	VisitAlterTableModifyQuery(expr *AlterTableModifyQuery) error
	VisitAlterTableModifyColumn(expr *AlterTableModifyColumn) error
	VisitAlterTableModifySetting(expr *AlterTableModifySetting) error
	VisitAlterTableResetSetting(expr *AlterTableResetSetting) error
	VisitAlterTableReplacePartition(expr *AlterTableReplacePartition) error
	VisitAlterTableDelete(expr *AlterTableDelete) error
	VisitAlterTableUpdate(expr *AlterTableUpdate) error
	VisitUpdateAssignment(expr *UpdateAssignment) error
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
	VisitCreateUser(expr *CreateUser) error
	VisitAuthenticationClause(expr *AuthenticationClause) error
	VisitHostClause(expr *HostClause) error
	VisitDefaultRoleClause(expr *DefaultRoleClause) error
	VisitGranteesClause(expr *GranteesClause) error
	VisitAlterRole(expr *AlterRole) error
	VisitRoleRenamePair(expr *RoleRenamePair) error
	VisitDestinationExpr(expr *DestinationClause) error
	VisitConstraintExpr(expr *ConstraintClause) error
	VisitNullLiteral(expr *NullLiteral) error
	VisitNotNullLiteral(expr *NotNullLiteral) error
	VisitPath(expr *Path) error
	VisitNestedIdentifier(expr *NestedIdentifier) error
	VisitTableIdentifier(expr *TableIdentifier) error
	VisitTableSchemaExpr(expr *TableSchemaClause) error
	VisitTableArgListExpr(expr *TableArgListExpr) error
	VisitTableFunctionExpr(expr *TableFunctionExpr) error
	VisitOnClusterExpr(expr *ClusterClause) error
	VisitPartitionExpr(expr *PartitionClause) error
	VisitPartitionByExpr(expr *PartitionByClause) error
	VisitPrimaryKeyExpr(expr *PrimaryKeyClause) error
	VisitSampleByExpr(expr *SampleByClause) error
	VisitTTLExpr(expr *TTLExpr) error
	VisitTTLExprList(expr *TTLClause) error
	VisitTTLPolicy(expr *TTLPolicy) error
	VisitTTLPolicyRule(expr *TTLPolicyRule) error
	VisitTTLPolicyItemAction(expr *TTLPolicyRuleAction) error
	VisitRefreshExpr(expr *RefreshExpr) error
	VisitOrderByExpr(expr *OrderExpr) error
	VisitOrderByListExpr(expr *OrderByClause) error
	VisitSettingsExpr(expr *SettingExpr) error
	VisitSettingsExprList(expr *SettingsClause) error
	VisitParamExprList(expr *ParamExprList) error
	VisitMapLiteral(expr *MapLiteral) error
	VisitArrayParamList(expr *ArrayParamList) error
	VisitQueryParam(expr *QueryParam) error
	VisitObjectParams(expr *ObjectParams) error
	VisitFunctionExpr(expr *FunctionExpr) error
	VisitWindowFunctionExpr(expr *WindowFunctionExpr) error
	VisitColumnDef(expr *ColumnDef) error
	VisitColumnExpr(expr *ColumnExpr) error
	VisitTypedPlaceholder(expr *TypedPlaceholder) error
	VisitScalarType(expr *ScalarType) error
	VisitJSONType(expr *JSONType) error
	VisitPropertyType(expr *PropertyType) error
	VisitTypeWithParams(expr *TypeWithParams) error
	VisitComplexType(expr *ComplexType) error
	VisitNestedType(expr *NestedType) error
	VisitCompressionCodec(expr *CompressionCodec) error
	VisitNumberLiteral(expr *NumberLiteral) error
	VisitStringLiteral(expr *StringLiteral) error
	VisitRatioExpr(expr *RatioExpr) error
	VisitEnumValue(expr *EnumValue) error
	VisitEnumType(expr *EnumType) error
	VisitIntervalExpr(expr *IntervalExpr) error
	VisitEngineExpr(expr *EngineExpr) error
	VisitColumnTypeExpr(expr *ColumnTypeExpr) error
	VisitColumnArgList(expr *ColumnArgList) error
	VisitColumnExprList(expr *ColumnExprList) error
	VisitWhenExpr(expr *WhenClause) error
	VisitCaseExpr(expr *CaseExpr) error
	VisitCastExpr(expr *CastExpr) error
	VisitWithExpr(expr *WithClause) error
	VisitTopExpr(expr *TopClause) error
	VisitCreateLiveView(expr *CreateLiveView) error
	VisitCreateDictionary(expr *CreateDictionary) error
	VisitDictionarySchemaClause(expr *DictionarySchemaClause) error
	VisitDictionaryAttribute(expr *DictionaryAttribute) error
	VisitDictionaryEngineClause(expr *DictionaryEngineClause) error
	VisitDictionaryPrimaryKeyClause(expr *DictionaryPrimaryKeyClause) error
	VisitDictionarySourceClause(expr *DictionarySourceClause) error
	VisitDictionaryArgExpr(expr *DictionaryArgExpr) error
	VisitDictionaryLifetimeClause(expr *DictionaryLifetimeClause) error
	VisitDictionaryLayoutClause(expr *DictionaryLayoutClause) error
	VisitDictionaryRangeClause(expr *DictionaryRangeClause) error
	VisitWithTimeoutExpr(expr *WithTimeoutClause) error
	VisitTableExpr(expr *TableExpr) error
	VisitOnExpr(expr *OnClause) error
	VisitUsingExpr(expr *UsingClause) error
	VisitJoinExpr(expr *JoinExpr) error
	VisitJoinConstraintExpr(expr *JoinConstraintClause) error
	VisitJoinTableExpr(expr *JoinTableExpr) error
	VisitFromExpr(expr *FromClause) error
	VisitIsNullExpr(expr *IsNullExpr) error
	VisitIsNotNullExpr(expr *IsNotNullExpr) error
	VisitAliasExpr(expr *AliasExpr) error
	VisitWhereExpr(expr *WhereClause) error
	VisitPrewhereExpr(expr *PrewhereClause) error
	VisitGroupByExpr(expr *GroupByClause) error
	VisitHavingExpr(expr *HavingClause) error
	VisitLimitExpr(expr *LimitClause) error
	VisitLimitByExpr(expr *LimitByClause) error
	VisitWindowConditionExpr(expr *WindowExpr) error
	VisitWindowExpr(expr *WindowClause) error
	VisitWindowFrameExpr(expr *WindowFrameClause) error
	VisitWindowFrameExtendExpr(expr *WindowFrameExtendExpr) error
	VisitBetweenClause(expr *BetweenClause) error
	VisitWindowFrameCurrentRow(expr *WindowFrameCurrentRow) error
	VisitWindowFrameUnbounded(expr *WindowFrameUnbounded) error
	VisitWindowFrameNumber(expr *WindowFrameNumber) error
	VisitArrayJoinExpr(expr *ArrayJoinClause) error
	VisitSelectQuery(expr *SelectQuery) error
	VisitSubQueryExpr(expr *SubQuery) error
	VisitNotExpr(expr *NotExpr) error
	VisitNegateExpr(expr *NegateExpr) error
	VisitGlobalInExpr(expr *GlobalInOperation) error
	VisitExtractExpr(expr *ExtractExpr) error
	VisitDropDatabase(expr *DropDatabase) error
	VisitDropStmt(expr *DropStmt) error
	VisitDropUserOrRole(expr *DropUserOrRole) error
	VisitUseExpr(expr *UseStmt) error
	VisitCTEExpr(expr *CTEStmt) error
	VisitSetExpr(expr *SetStmt) error
	VisitFormatExpr(expr *FormatClause) error
	VisitOptimizeExpr(expr *OptimizeStmt) error
	VisitDeduplicateExpr(expr *DeduplicateClause) error
	VisitSystemExpr(expr *SystemStmt) error
	VisitSystemFlushExpr(expr *SystemFlushExpr) error
	VisitSystemReloadExpr(expr *SystemReloadExpr) error
	VisitSystemSyncExpr(expr *SystemSyncExpr) error
	VisitSystemCtrlExpr(expr *SystemCtrlExpr) error
	VisitSystemDropExpr(expr *SystemDropExpr) error
	VisitTruncateTable(expr *TruncateTable) error
	VisitSampleRatioExpr(expr *SampleClause) error
	VisitPlaceHolderExpr(expr *PlaceHolder) error
	VisitDeleteFromExpr(expr *DeleteClause) error
	VisitColumnNamesExpr(expr *ColumnNamesExpr) error
	VisitValuesExpr(expr *AssignmentValues) error
	VisitInsertExpr(expr *InsertStmt) error
	VisitCheckExpr(expr *CheckStmt) error
	VisitUnaryExpr(expr *UnaryExpr) error
	VisitRenameStmt(expr *RenameStmt) error
	VisitExplainExpr(expr *ExplainStmt) error
	VisitPrivilegeExpr(expr *PrivilegeClause) error
	VisitGrantPrivilegeExpr(expr *GrantPrivilegeStmt) error
	VisitShowExpr(expr *ShowStmt) error
	VisitDescribeExpr(expr *DescribeStmt) error
	VisitSelectItem(expr *SelectItem) error
	VisitTargetPairExpr(expr *TargetPair) error
	VisitDistinctOn(expr *DistinctOn) error

	Enter(expr Expr)
	Leave(expr Expr)
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

func (v *DefaultASTVisitor) VisitTernaryExpr(expr *TernaryOperation) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitBinaryExpr(expr *BinaryOperation) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitIndexOperation(expr *IndexOperation) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitJoinTableExpr(expr *JoinTableExpr) error {
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

func (v *DefaultASTVisitor) VisitAlterTableAddProjection(expr *AlterTableAddProjection) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitProjectionOrderBy(expr *ProjectionOrderByClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitProjectionSelect(expr *ProjectionSelectStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTableProjection(expr *TableProjection) error {
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

func (v *DefaultASTVisitor) VisitAlterTableDropProjection(expr *AlterTableDropProjection) error {
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

func (v *DefaultASTVisitor) VisitAlterTableClearProjection(expr *AlterTableClearProjection) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableMaterializeProjection(expr *AlterTableMaterializeProjection) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableMaterializeIndex(expr *AlterTableMaterializeIndex) error {
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

func (v *DefaultASTVisitor) VisitAlterTableModifyQuery(expr *AlterTableModifyQuery) error {
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

func (v *DefaultASTVisitor) VisitAlterTableModifySetting(expr *AlterTableModifySetting) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableResetSetting(expr *AlterTableResetSetting) error {
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

func (v *DefaultASTVisitor) VisitAlterTableDelete(expr *AlterTableDelete) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAlterTableUpdate(expr *AlterTableUpdate) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitUpdateAssignment(expr *UpdateAssignment) error {
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

func (v *DefaultASTVisitor) VisitCreateUser(expr *CreateUser) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitAuthenticationClause(expr *AuthenticationClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitHostClause(expr *HostClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDefaultRoleClause(expr *DefaultRoleClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitGranteesClause(expr *GranteesClause) error {
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

func (v *DefaultASTVisitor) VisitDestinationExpr(expr *DestinationClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitConstraintExpr(expr *ConstraintClause) error {
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

func (v *DefaultASTVisitor) VisitPath(expr *Path) error {
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

func (v *DefaultASTVisitor) VisitTableSchemaExpr(expr *TableSchemaClause) error {
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

func (v *DefaultASTVisitor) VisitOnClusterExpr(expr *ClusterClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPartitionExpr(expr *PartitionClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPartitionByExpr(expr *PartitionByClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPrimaryKeyExpr(expr *PrimaryKeyClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSampleByExpr(expr *SampleByClause) error {
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

func (v *DefaultASTVisitor) VisitTTLExprList(expr *TTLClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTTLPolicy(expr *TTLPolicy) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTTLPolicyRule(expr *TTLPolicyRule) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTTLPolicyItemAction(expr *TTLPolicyRuleAction) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitRefreshExpr(expr *RefreshExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOrderByExpr(expr *OrderExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOrderByListExpr(expr *OrderByClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSettingsExpr(expr *SettingExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSettingsExprList(expr *SettingsClause) error {
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

func (v *DefaultASTVisitor) VisitQueryParam(expr *QueryParam) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitMapLiteral(expr *MapLiteral) error {
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

func (v *DefaultASTVisitor) VisitColumnDef(expr *ColumnDef) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitColumnExpr(expr *ColumnExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTypedPlaceholder(expr *TypedPlaceholder) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitScalarType(expr *ScalarType) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitJSONType(expr *JSONType) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPropertyType(expr *PropertyType) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTypeWithParams(expr *TypeWithParams) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitComplexType(expr *ComplexType) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitNestedType(expr *NestedType) error {
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

func (v *DefaultASTVisitor) VisitEnumValue(expr *EnumValue) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitEnumType(expr *EnumType) error {
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

func (v *DefaultASTVisitor) VisitWhenExpr(expr *WhenClause) error {
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

func (v *DefaultASTVisitor) VisitWithExpr(expr *WithClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTopExpr(expr *TopClause) error {
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

func (v *DefaultASTVisitor) VisitCreateDictionary(expr *CreateDictionary) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionarySchemaClause(expr *DictionarySchemaClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionaryAttribute(expr *DictionaryAttribute) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionaryEngineClause(expr *DictionaryEngineClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionaryPrimaryKeyClause(expr *DictionaryPrimaryKeyClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionarySourceClause(expr *DictionarySourceClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionaryArgExpr(expr *DictionaryArgExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionaryLifetimeClause(expr *DictionaryLifetimeClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionaryLayoutClause(expr *DictionaryLayoutClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDictionaryRangeClause(expr *DictionaryRangeClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWithTimeoutExpr(expr *WithTimeoutClause) error {
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

func (v *DefaultASTVisitor) VisitOnExpr(expr *OnClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitUsingExpr(expr *UsingClause) error {
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

func (v *DefaultASTVisitor) VisitJoinConstraintExpr(expr *JoinConstraintClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitFromExpr(expr *FromClause) error {
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

func (v *DefaultASTVisitor) VisitWhereExpr(expr *WhereClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPrewhereExpr(expr *PrewhereClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitGroupByExpr(expr *GroupByClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitHavingExpr(expr *HavingClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitLimitExpr(expr *LimitClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitLimitByExpr(expr *LimitByClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowConditionExpr(expr *WindowExpr) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowExpr(expr *WindowClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitWindowFrameExpr(expr *WindowFrameClause) error {
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

func (v *DefaultASTVisitor) VisitBetweenClause(expr *BetweenClause) error {
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

func (v *DefaultASTVisitor) VisitArrayJoinExpr(expr *ArrayJoinClause) error {
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

func (v *DefaultASTVisitor) VisitSubQueryExpr(expr *SubQuery) error {
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

func (v *DefaultASTVisitor) VisitGlobalInExpr(expr *GlobalInOperation) error {
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

func (v *DefaultASTVisitor) VisitUseExpr(expr *UseStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCTEExpr(expr *CTEStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSetExpr(expr *SetStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitFormatExpr(expr *FormatClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitOptimizeExpr(expr *OptimizeStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDeduplicateExpr(expr *DeduplicateClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSystemExpr(expr *SystemStmt) error {
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

func (v *DefaultASTVisitor) VisitSampleRatioExpr(expr *SampleClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPlaceHolderExpr(expr *PlaceHolder) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDeleteFromExpr(expr *DeleteClause) error {
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

func (v *DefaultASTVisitor) VisitValuesExpr(expr *AssignmentValues) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitInsertExpr(expr *InsertStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitCheckExpr(expr *CheckStmt) error {
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

func (v *DefaultASTVisitor) VisitExplainExpr(expr *ExplainStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitPrivilegeExpr(expr *PrivilegeClause) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitGrantPrivilegeExpr(expr *GrantPrivilegeStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitShowExpr(expr *ShowStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDescribeExpr(expr *DescribeStmt) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitSelectItem(expr *SelectItem) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitTargetPairExpr(expr *TargetPair) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) VisitDistinctOn(expr *DistinctOn) error {
	if v.Visit != nil {
		return v.Visit(expr)
	}
	return nil
}

func (v *DefaultASTVisitor) Enter(expr Expr) {}

func (v *DefaultASTVisitor) Leave(expr Expr) {}
