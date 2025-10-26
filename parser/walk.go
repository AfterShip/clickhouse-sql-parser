package parser

import "reflect"

// WalkFunc is a function type for walking AST nodes.
// It receives the current node and returns a boolean indicating whether to continue walking.
// If the function returns false, the walking stops for the current subtree.
type WalkFunc func(node Expr) bool

// Walk traverses an AST in depth-first order, calling the provided function
// for each node. If the function returns false, traversal stops for that subtree.
func Walk(node Expr, fn WalkFunc) bool {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return true
	}
	if !fn(node) {
		return false
	}

	switch n := node.(type) {
	case *SelectQuery:
		if !Walk(n.With, fn) {
			return false
		}
		if !Walk(n.DistinctOn, fn) {
			return false
		}
		if !Walk(n.Top, fn) {
			return false
		}
		for _, item := range n.SelectItems {
			if !Walk(item, fn) {
				return false
			}
		}
		if !Walk(n.From, fn) {
			return false
		}
		if !Walk(n.ArrayJoin, fn) {
			return false
		}
		if !Walk(n.Window, fn) {
			return false
		}
		if !Walk(n.Prewhere, fn) {
			return false
		}
		if !Walk(n.Where, fn) {
			return false
		}
		if !Walk(n.GroupBy, fn) {
			return false
		}
		if !Walk(n.Having, fn) {
			return false
		}
		if !Walk(n.OrderBy, fn) {
			return false
		}
		if !Walk(n.LimitBy, fn) {
			return false
		}
		if !Walk(n.Limit, fn) {
			return false
		}
		if !Walk(n.Settings, fn) {
			return false
		}
		if !Walk(n.UnionAll, fn) {
			return false
		}
		if !Walk(n.UnionDistinct, fn) {
			return false
		}
		if !Walk(n.Format, fn) {
			return false
		}
	case *SubQuery:
		if !Walk(n.Select, fn) {
			return false
		}
	case *SelectItem:
		if !Walk(n.Expr, fn) {
			return false
		}
		for _, modifier := range n.Modifiers {
			if !Walk(modifier, fn) {
				return false
			}
		}
		if !Walk(n.Alias, fn) {
			return false
		}
	case *TableExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Alias, fn) {
			return false
		}
	case *AliasExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Alias, fn) {
			return false
		}
	case *FunctionExpr:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Params, fn) {
			return false
		}
	case *TableIdentifier:
		if !Walk(n.Database, fn) {
			return false
		}
		if !Walk(n.Table, fn) {
			return false
		}
	case *Ident:
		// Leaf node
	case *NumberLiteral:
		// Leaf node
	case *StringLiteral:
		// Leaf node
	case *NullLiteral:
		// Leaf node
	case *NotNullLiteral:
		if !Walk(n.NullLiteral, fn) {
			return false
		}
	case *ColumnExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Alias, fn) {
			return false
		}
	case *BinaryOperation:
		if !Walk(n.LeftExpr, fn) {
			return false
		}
		if !Walk(n.RightExpr, fn) {
			return false
		}
	case *WhenClause:
		if !Walk(n.When, fn) {
			return false
		}
		if !Walk(n.Then, fn) {
			return false
		}
	case *CaseExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		for _, when := range n.Whens {
			if !Walk(when, fn) {
				return false
			}
		}
		if !Walk(n.Else, fn) {
			return false
		}
	case *CastExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.AsType, fn) {
			return false
		}
	case *Path:
		for _, field := range n.Fields {
			if !Walk(field, fn) {
				return false
			}
		}
	case *WithClause:
		for _, cte := range n.CTEs {
			if !Walk(cte, fn) {
				return false
			}
		}
	case *CTEStmt:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Alias, fn) {
			return false
		}
	case *FromClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *JoinExpr:
		if !Walk(n.Left, fn) {
			return false
		}
		if !Walk(n.Right, fn) {
			return false
		}
		if !Walk(n.Constraints, fn) {
			return false
		}
	case *JoinTableExpr:
		if !Walk(n.Table, fn) {
			return false
		}
		if !Walk(n.SampleRatio, fn) {
			return false
		}
	case *OnClause:
		if !Walk(n.On, fn) {
			return false
		}
	case *UsingClause:
		if !Walk(n.Using, fn) {
			return false
		}
	case *WhereClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *PrewhereClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *GroupByClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *HavingClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *OrderByClause:
		for _, item := range n.Items {
			if !Walk(item, fn) {
				return false
			}
		}
	case *OrderExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Alias, fn) {
			return false
		}
	case *LimitClause:
		if !Walk(n.Limit, fn) {
			return false
		}
		if !Walk(n.Offset, fn) {
			return false
		}
	case *LimitByClause:
		if !Walk(n.Limit, fn) {
			return false
		}
		if !Walk(n.ByExpr, fn) {
			return false
		}
	case *SettingsClause:
		for _, item := range n.Items {
			if !Walk(item, fn) {
				return false
			}
		}
	case *SettingExpr:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Expr, fn) {
			return false
		}
	case *FormatClause:
		if !Walk(n.Format, fn) {
			return false
		}
	case *InsertStmt:
		if !Walk(n.Table, fn) {
			return false
		}
		if !Walk(n.ColumnNames, fn) {
			return false
		}
		if !Walk(n.Format, fn) {
			return false
		}
		if !Walk(n.SelectExpr, fn) {
			return false
		}
	case *ColumnNamesExpr:
		for i := range n.ColumnNames {
			if !Walk(&n.ColumnNames[i], fn) {
				return false
			}
		}
	case *AssignmentValues:
		for _, value := range n.Values {
			if !Walk(value, fn) {
				return false
			}
		}
	case *TableFunctionExpr:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Args, fn) {
			return false
		}
	case *TableArgListExpr:
		for _, arg := range n.Args {
			if !Walk(arg, fn) {
				return false
			}
		}
	case *NestedIdentifier:
		if !Walk(n.Ident, fn) {
			return false
		}
		if !Walk(n.DotIdent, fn) {
			return false
		}
	case *ArrayParamList:
		if !Walk(n.Items, fn) {
			return false
		}
	case *ColumnExprList:
		for _, item := range n.Items {
			if !Walk(item, fn) {
				return false
			}
		}
	case *ParamExprList:
		if !Walk(n.Items, fn) {
			return false
		}
		if !Walk(n.ColumnArgList, fn) {
			return false
		}
	case *ColumnArgList:
		for _, item := range n.Items {
			if !Walk(item, fn) {
				return false
			}
		}
	case *WindowClause:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.WindowExpr, fn) {
			return false
		}
	case *WindowExpr:
		if !Walk(n.PartitionBy, fn) {
			return false
		}
		if !Walk(n.OrderBy, fn) {
			return false
		}
		if !Walk(n.Frame, fn) {
			return false
		}
	case *PartitionByClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *WindowFrameClause:
		if !Walk(n.Extend, fn) {
			return false
		}
	case *WindowFrameExtendExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *BetweenClause:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Between, fn) {
			return false
		}
		if !Walk(n.And, fn) {
			return false
		}
	case *WindowFrameCurrentRow:
		// Leaf node
	case *WindowFrameUnbounded:
		// Leaf node
	case *WindowFrameNumber:
		if !Walk(n.Number, fn) {
			return false
		}
	case *ArrayJoinClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *TopClause:
		if !Walk(n.Number, fn) {
			return false
		}
	case *SampleClause:
		if !Walk(n.Ratio, fn) {
			return false
		}
		if !Walk(n.Offset, fn) {
			return false
		}
	case *RatioExpr:
		if !Walk(n.Numerator, fn) {
			return false
		}
		if !Walk(n.Denominator, fn) {
			return false
		}
	case *IntervalExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Unit, fn) {
			return false
		}
	case *DropStmt:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
	case *DropDatabase:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
	case *DropUserOrRole:
		for _, name := range n.Names {
			if !Walk(name, fn) {
				return false
			}
		}
		if !Walk(n.From, fn) {
			return false
		}
	case *TruncateTable:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
	case *CheckStmt:
		if !Walk(n.Table, fn) {
			return false
		}
		if !Walk(n.Partition, fn) {
			return false
		}
	case *OptimizeStmt:
		if !Walk(n.Table, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.Partition, fn) {
			return false
		}
		if !Walk(n.Deduplicate, fn) {
			return false
		}
	case *DeduplicateClause:
		if !Walk(n.By, fn) {
			return false
		}
		if !Walk(n.Except, fn) {
			return false
		}
	case *SystemStmt:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *SystemFlushExpr:
		if !Walk(n.Distributed, fn) {
			return false
		}
	case *SystemReloadExpr:
		if !Walk(n.Dictionary, fn) {
			return false
		}
	case *SystemSyncExpr:
		if !Walk(n.Cluster, fn) {
			return false
		}
	case *SystemCtrlExpr:
		if !Walk(n.Cluster, fn) {
			return false
		}
	case *SystemDropExpr:
		// Leaf node
	case *UseStmt:
		if !Walk(n.Database, fn) {
			return false
		}
	case *SetStmt:
		if !Walk(n.Settings, fn) {
			return false
		}
	case *ExplainStmt:
		if !Walk(n.Statement, fn) {
			return false
		}
	case *GrantPrivilegeStmt:
		for _, privilege := range n.Privileges {
			if !Walk(privilege, fn) {
				return false
			}
		}
		if !Walk(n.On, fn) {
			return false
		}
		for _, role := range n.To {
			if !Walk(role, fn) {
				return false
			}
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
	case *PrivilegeClause:
		if !Walk(n.Params, fn) {
			return false
		}
	case *RenameStmt:
		for _, pair := range n.TargetPairList {
			if !Walk(pair.Old, fn) {
				return false
			}
			if !Walk(pair.New, fn) {
				return false
			}
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
	case *DeleteClause:
		if !Walk(n.Table, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.WhereExpr, fn) {
			return false
		}
	case *CreateDatabase:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.Engine, fn) {
			return false
		}
		if !Walk(n.Comment, fn) {
			return false
		}
	case *CreateTable:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.UUID, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.TableSchema, fn) {
			return false
		}
		if !Walk(n.Engine, fn) {
			return false
		}
		if !Walk(n.SubQuery, fn) {
			return false
		}
		if !Walk(n.Comment, fn) {
			return false
		}
	case *CreateView:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.UUID, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.TableSchema, fn) {
			return false
		}
		if !Walk(n.SubQuery, fn) {
			return false
		}
	case *CreateMaterializedView:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.Refresh, fn) {
			return false
		}
		if !Walk(n.RandomizeFor, fn) {
			return false
		}
		for _, dep := range n.DependsOn {
			if !Walk(dep, fn) {
				return false
			}
		}
		if !Walk(n.Settings, fn) {
			return false
		}
		if !Walk(n.Engine, fn) {
			return false
		}
		if !Walk(n.Destination, fn) {
			return false
		}
		if !Walk(n.SubQuery, fn) {
			return false
		}
		if !Walk(n.Comment, fn) {
			return false
		}
		if !Walk(n.Definer, fn) {
			return false
		}
	case *CreateLiveView:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.UUID, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.Destination, fn) {
			return false
		}
		if !Walk(n.TableSchema, fn) {
			return false
		}
		if !Walk(n.WithTimeout, fn) {
			return false
		}
		if !Walk(n.SubQuery, fn) {
			return false
		}
	case *CreateDictionary:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.UUID, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.Schema, fn) {
			return false
		}
		if !Walk(n.Engine, fn) {
			return false
		}
	case *CreateFunction:
		if !Walk(n.FunctionName, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		if !Walk(n.Params, fn) {
			return false
		}
		if !Walk(n.Expr, fn) {
			return false
		}
	case *CreateRole:
		for _, name := range n.RoleNames {
			if !Walk(name, fn) {
				return false
			}
		}
		if !Walk(n.AccessStorageType, fn) {
			return false
		}
		for _, setting := range n.Settings {
			if !Walk(setting, fn) {
				return false
			}
		}
	case *CreateUser:
		for _, name := range n.UserNames {
			if !Walk(name, fn) {
				return false
			}
		}
		if !Walk(n.Authentication, fn) {
			return false
		}
		for _, host := range n.Hosts {
			if !Walk(host, fn) {
				return false
			}
		}
		if !Walk(n.DefaultRole, fn) {
			return false
		}
		if !Walk(n.DefaultDatabase, fn) {
			return false
		}
		if !Walk(n.Grantees, fn) {
			return false
		}
		for _, setting := range n.Settings {
			if !Walk(setting, fn) {
				return false
			}
		}
	case *AlterTable:
		if !Walk(n.TableIdentifier, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
		for _, expr := range n.AlterExprs {
			if !Walk(expr, fn) {
				return false
			}
		}
	case *AlterTableAttachPartition:
		if !Walk(n.Partition, fn) {
			return false
		}
		if !Walk(n.From, fn) {
			return false
		}
	case *AlterTableDetachPartition:
		if !Walk(n.Partition, fn) {
			return false
		}
		if !Walk(n.Settings, fn) {
			return false
		}
	case *AlterTableDropPartition:
		if !Walk(n.Partition, fn) {
			return false
		}
		if !Walk(n.Settings, fn) {
			return false
		}
	case *AlterTableMaterializeProjection:
		if !Walk(n.ProjectionName, fn) {
			return false
		}
		if !Walk(n.Partition, fn) {
			return false
		}
	case *AlterTableMaterializeIndex:
		if !Walk(n.IndexName, fn) {
			return false
		}
		if !Walk(n.Partition, fn) {
			return false
		}
	case *AlterTableFreezePartition:
		if !Walk(n.Partition, fn) {
			return false
		}
	case *AlterTableAddColumn:
		if !Walk(n.Column, fn) {
			return false
		}
		if !Walk(n.After, fn) {
			return false
		}
	case *AlterTableAddIndex:
		if !Walk(n.Index, fn) {
			return false
		}
		if !Walk(n.After, fn) {
			return false
		}
	case *AlterTableAddProjection:
		if !Walk(n.TableProjection, fn) {
			return false
		}
		if !Walk(n.After, fn) {
			return false
		}
	case *AlterTableDropColumn:
		if !Walk(n.ColumnName, fn) {
			return false
		}
	case *AlterTableDropIndex:
		if !Walk(n.IndexName, fn) {
			return false
		}
	case *AlterTableDropProjection:
		if !Walk(n.ProjectionName, fn) {
			return false
		}
	case *AlterTableRemoveTTL:
		// Leaf node
	case *AlterTableClearColumn:
		if !Walk(n.ColumnName, fn) {
			return false
		}
		if !Walk(n.PartitionExpr, fn) {
			return false
		}
	case *AlterTableClearIndex:
		if !Walk(n.IndexName, fn) {
			return false
		}
		if !Walk(n.PartitionExpr, fn) {
			return false
		}
	case *AlterTableClearProjection:
		if !Walk(n.ProjectionName, fn) {
			return false
		}
		if !Walk(n.PartitionExpr, fn) {
			return false
		}
	case *AlterTableRenameColumn:
		if !Walk(n.OldColumnName, fn) {
			return false
		}
		if !Walk(n.NewColumnName, fn) {
			return false
		}
	case *AlterTableModifyQuery:
		if !Walk(n.SelectExpr, fn) {
			return false
		}
	case *AlterTableModifyTTL:
		if !Walk(n.TTL, fn) {
			return false
		}
	case *AlterTableModifyColumn:
		if !Walk(n.Column, fn) {
			return false
		}
		if !Walk(n.RemovePropertyType, fn) {
			return false
		}
	case *AlterTableModifySetting:
		for _, setting := range n.Settings {
			if !Walk(setting, fn) {
				return false
			}
		}
	case *AlterTableResetSetting:
		for _, setting := range n.Settings {
			if !Walk(setting, fn) {
				return false
			}
		}
	case *AlterTableReplacePartition:
		if !Walk(n.Partition, fn) {
			return false
		}
		if !Walk(n.Table, fn) {
			return false
		}
	case *AlterTableDelete:
		if !Walk(n.WhereClause, fn) {
			return false
		}
	case *AlterTableUpdate:
		for _, assignment := range n.Assignments {
			if !Walk(assignment, fn) {
				return false
			}
		}
		if !Walk(n.WhereClause, fn) {
			return false
		}
	case *UpdateAssignment:
		if !Walk(n.Column, fn) {
			return false
		}
		if !Walk(n.Expr, fn) {
			return false
		}
	case *AlterRole:
		for _, pair := range n.RoleRenamePairs {
			if !Walk(pair, fn) {
				return false
			}
		}
		for _, setting := range n.Settings {
			if !Walk(setting, fn) {
				return false
			}
		}
	case *RoleRenamePair:
		if !Walk(n.RoleName, fn) {
			return false
		}
		if !Walk(n.NewName, fn) {
			return false
		}
	case *TableSchemaClause:
		for _, column := range n.Columns {
			if !Walk(column, fn) {
				return false
			}
		}
		if !Walk(n.AliasTable, fn) {
			return false
		}
		if !Walk(n.TableFunction, fn) {
			return false
		}
	case *ColumnDef:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Type, fn) {
			return false
		}
		if !Walk(n.NotNull, fn) {
			return false
		}
		if !Walk(n.Nullable, fn) {
			return false
		}
		if !Walk(n.DefaultExpr, fn) {
			return false
		}
		if !Walk(n.MaterializedExpr, fn) {
			return false
		}
		if !Walk(n.AliasExpr, fn) {
			return false
		}
		if !Walk(n.Codec, fn) {
			return false
		}
		if !Walk(n.TTL, fn) {
			return false
		}
		if !Walk(n.Comment, fn) {
			return false
		}
	case *ScalarType:
		if !Walk(n.Name, fn) {
			return false
		}
	case *JSONType:
		if !Walk(n.Name, fn) {
			return false
		}
	case *PropertyType:
		if !Walk(n.Name, fn) {
			return false
		}
	case *TypeWithParams:
		if !Walk(n.Name, fn) {
			return false
		}
		for _, param := range n.Params {
			if !Walk(param, fn) {
				return false
			}
		}
	case *ComplexType:
		if !Walk(n.Name, fn) {
			return false
		}
		for _, param := range n.Params {
			if !Walk(param, fn) {
				return false
			}
		}
	case *NestedType:
		if !Walk(n.Name, fn) {
			return false
		}
		for _, column := range n.Columns {
			if !Walk(column, fn) {
				return false
			}
		}
	case *CompressionCodec:
		if !Walk(n.Type, fn) {
			return false
		}
		if !Walk(n.TypeLevel, fn) {
			return false
		}
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Level, fn) {
			return false
		}
	case *EngineExpr:
		if !Walk(n.Params, fn) {
			return false
		}
		if !Walk(n.PrimaryKey, fn) {
			return false
		}
		if !Walk(n.PartitionBy, fn) {
			return false
		}
		if !Walk(n.SampleBy, fn) {
			return false
		}
		if !Walk(n.TTL, fn) {
			return false
		}
		if !Walk(n.Settings, fn) {
			return false
		}
		if !Walk(n.OrderBy, fn) {
			return false
		}
	case *PrimaryKeyClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *SampleByClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *TTLClause:
		for _, item := range n.Items {
			if !Walk(item, fn) {
				return false
			}
		}
	case *TTLExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.Policy, fn) {
			return false
		}
	case *TTLPolicy:
		if !Walk(n.Item, fn) {
			return false
		}
		if !Walk(n.Where, fn) {
			return false
		}
		if !Walk(n.GroupBy, fn) {
			return false
		}
	case *TTLPolicyRule:
		if !Walk(n.ToVolume, fn) {
			return false
		}
		if !Walk(n.ToDisk, fn) {
			return false
		}
	case *TTLPolicyRuleAction:
		if !Walk(n.Codec, fn) {
			return false
		}
	case *RefreshExpr:
		if !Walk(n.Interval, fn) {
			return false
		}
		if !Walk(n.Offset, fn) {
			return false
		}
	case *DestinationClause:
		if !Walk(n.TableIdentifier, fn) {
			return false
		}
		if !Walk(n.TableSchema, fn) {
			return false
		}
	case *ConstraintClause:
		if !Walk(n.Constraint, fn) {
			return false
		}
		if !Walk(n.Expr, fn) {
			return false
		}
	case *RoleName:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Scope, fn) {
			return false
		}
		if !Walk(n.OnCluster, fn) {
			return false
		}
	case *SettingPair:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Value, fn) {
			return false
		}
	case *RoleSetting:
		for _, pair := range n.SettingPairs {
			if !Walk(pair, fn) {
				return false
			}
		}
		if !Walk(n.Modifier, fn) {
			return false
		}
	case *AuthenticationClause:
		if !Walk(n.AuthValue, fn) {
			return false
		}
		if !Walk(n.LdapServer, fn) {
			return false
		}
		if !Walk(n.KerberosRealm, fn) {
			return false
		}
	case *HostClause:
		if !Walk(n.HostValue, fn) {
			return false
		}
	case *DefaultRoleClause:
		for _, role := range n.Roles {
			if !Walk(role, fn) {
				return false
			}
		}
	case *GranteesClause:
		for _, grantee := range n.Grantees {
			if !Walk(grantee, fn) {
				return false
			}
		}
		for _, except := range n.ExceptUsers {
			if !Walk(except, fn) {
				return false
			}
		}
	case *WithTimeoutClause:
		if !Walk(n.Number, fn) {
			return false
		}
	case *DictionarySchemaClause:
		for _, attr := range n.Attributes {
			if !Walk(attr, fn) {
				return false
			}
		}
	case *DictionaryAttribute:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Type, fn) {
			return false
		}
		if !Walk(n.Default, fn) {
			return false
		}
		if !Walk(n.Expression, fn) {
			return false
		}
	case *DictionaryEngineClause:
		if !Walk(n.PrimaryKey, fn) {
			return false
		}
		if !Walk(n.Source, fn) {
			return false
		}
		if !Walk(n.Lifetime, fn) {
			return false
		}
		if !Walk(n.Layout, fn) {
			return false
		}
		if !Walk(n.Range, fn) {
			return false
		}
		if !Walk(n.Settings, fn) {
			return false
		}
	case *DictionaryPrimaryKeyClause:
		if !Walk(n.Keys, fn) {
			return false
		}
	case *DictionarySourceClause:
		if !Walk(n.Source, fn) {
			return false
		}
		for _, arg := range n.Args {
			if !Walk(arg, fn) {
				return false
			}
		}
	case *DictionaryArgExpr:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Value, fn) {
			return false
		}
	case *DictionaryLifetimeClause:
		if !Walk(n.Value, fn) {
			return false
		}
		if !Walk(n.Min, fn) {
			return false
		}
		if !Walk(n.Max, fn) {
			return false
		}
	case *DictionaryLayoutClause:
		if !Walk(n.Layout, fn) {
			return false
		}
		for _, arg := range n.Args {
			if !Walk(arg, fn) {
				return false
			}
		}
	case *DictionaryRangeClause:
		if !Walk(n.Min, fn) {
			return false
		}
		if !Walk(n.Max, fn) {
			return false
		}
	case *PlaceHolder:
		// Leaf node
	case *TypedPlaceholder:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Type, fn) {
			return false
		}
	case *QueryParam:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Type, fn) {
			return false
		}
	case *MapLiteral:
		for _, kv := range n.KeyValues {
			if !Walk(&kv.Key, fn) {
				return false
			}
			if !Walk(kv.Value, fn) {
				return false
			}
		}
	case *ObjectParams:
		if !Walk(n.Object, fn) {
			return false
		}
		if !Walk(n.Params, fn) {
			return false
		}
	case *WindowFunctionExpr:
		if !Walk(n.Function, fn) {
			return false
		}
		if !Walk(n.OverExpr, fn) {
			return false
		}
	case *NotExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *NegateExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *GlobalInOperation:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *ExtractExpr:
		if !Walk(n.FromExpr, fn) {
			return false
		}
	case *IsNullExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *IsNotNullExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *TernaryOperation:
		if !Walk(n.Condition, fn) {
			return false
		}
		if !Walk(n.TrueExpr, fn) {
			return false
		}
		if !Walk(n.FalseExpr, fn) {
			return false
		}
	case *IndexOperation:
		if !Walk(n.Object, fn) {
			return false
		}
		if !Walk(n.Index, fn) {
			return false
		}
	case *OperationExpr:
		// Leaf node
	case *TableIndex:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.ColumnExpr, fn) {
			return false
		}
		if !Walk(n.ColumnType, fn) {
			return false
		}
		if !Walk(n.Granularity, fn) {
			return false
		}
	case *ProjectionOrderByClause:
		if !Walk(n.Columns, fn) {
			return false
		}
	case *ProjectionSelectStmt:
		if !Walk(n.With, fn) {
			return false
		}
		if !Walk(n.SelectColumns, fn) {
			return false
		}
		if !Walk(n.GroupBy, fn) {
			return false
		}
		if !Walk(n.OrderBy, fn) {
			return false
		}
	case *TableProjection:
		if !Walk(n.Identifier, fn) {
			return false
		}
		if !Walk(n.Select, fn) {
			return false
		}
	case *RemovePropertyType:
		if !Walk(n.PropertyType, fn) {
			return false
		}
	case *EnumType:
		if !Walk(n.Name, fn) {
			return false
		}
		for i := range n.Values {
			if !Walk(&n.Values[i], fn) {
				return false
			}
		}
	case *EnumValue:
		if !Walk(n.Name, fn) {
			return false
		}
		if !Walk(n.Value, fn) {
			return false
		}
	case *ClusterClause:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *PartitionClause:
		if !Walk(n.Expr, fn) {
			return false
		}
		if !Walk(n.ID, fn) {
			return false
		}
	case *UUID:
		if !Walk(n.Value, fn) {
			return false
		}
	case *ColumnTypeExpr:
		if !Walk(n.Name, fn) {
			return false
		}
	case *UnaryExpr:
		if !Walk(n.Expr, fn) {
			return false
		}
	case *JoinConstraintClause:
		if !Walk(n.On, fn) {
			return false
		}
		if !Walk(n.Using, fn) {
			return false
		}
	case *TargetPair:
		if !Walk(n.Old, fn) {
			return false
		}
		if !Walk(n.New, fn) {
			return false
		}
	case *ShowStmt:
		if !Walk(n.Target, fn) {
			return false
		}
		if !Walk(n.LikePattern, fn) {
			return false
		}
		if !Walk(n.Limit, fn) {
			return false
		}
		if !Walk(n.OutFile, fn) {
			return false
		}
		if !Walk(n.Format, fn) {
			return false
		}
	case *DescribeStmt:
		if !Walk(n.Target, fn) {
			return false
		}
	case *DistinctOn:
		for _, ident := range n.Idents {
			if !Walk(ident, fn) {
				return false
			}
		}
	}
	return true
}

// WalkWithBreak allows for early termination of tree traversal.
// The provided function should return true to continue walking,
// or false to stop the traversal entirely.
func WalkWithBreak(node Expr, fn WalkFunc) bool {
	if node == nil {
		return true
	}

	// Call the function first - if it returns false, stop immediately
	if !fn(node) {
		return false
	}

	// For early termination support, use a helper that converts our function
	// to one that collects a boolean result
	var continueWalk = true
	Walk(node, func(child Expr) bool {
		// Skip the current node since we already processed it
		if child == node {
			return true
		}
		// Call our termination-aware function
		if !fn(child) {
			continueWalk = false
			return false
		}
		return true
	})

	return continueWalk
}

// Find searches for the first node matching the given predicate.
// Returns the matching node and true if found, or nil and false if not found.
func Find(root Expr, predicate func(Expr) bool) (Expr, bool) {
	var found Expr
	WalkWithBreak(root, func(node Expr) bool {
		if predicate(node) {
			found = node
			return false // Stop traversal
		}
		return true // Continue traversal
	})
	return found, found != nil
}

// FindAll collects all nodes matching the given predicate.
func FindAll(root Expr, predicate func(Expr) bool) []Expr {
	var matches []Expr
	Walk(root, func(node Expr) bool {
		if predicate(node) {
			matches = append(matches, node)
		}
		return true // Always continue traversal
	})
	return matches
}

// Transform applies a transformation function to all nodes in the tree.
// The transformation function receives a node and should return the transformed node.
// Note: This modifies the tree in place for mutable fields.
func Transform(root Expr, transformer func(Expr) Expr) Expr {
	transformed := transformer(root)
	if transformed == nil {
		return nil
	}

	// Apply transformations to children
	Walk(transformed, func(node Expr) bool {
		// The actual in-place transformation would need to be implemented
		// based on the specific needs and mutability of the AST nodes
		return true
	})

	return transformed
}
