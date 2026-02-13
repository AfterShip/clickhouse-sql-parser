package parser

import "strings"

const (
	whitespace byte = ' '
	newline    byte = '\n'
)

type FormatMode int

const (
	FormatModeCompact FormatMode = iota + 1
	FormatModeBeautify
)

// Formatter renders SQL.
type Formatter struct {
	builder strings.Builder

	mode        FormatMode
	indentLevel int
	lineStart   bool
	indent      string
}

func NewFormatter() *Formatter {
	return &Formatter{
		mode:      FormatModeCompact,
		lineStart: true,
		indent:    "  ",
	}
}

func (f *Formatter) WithBeautify() *Formatter {
	f.mode = FormatModeBeautify
	return f
}

// WithIndent sets the indentation string used when beautifying SQL.
// The indent parameter should not be empty to maintain proper formatting.
func (f *Formatter) WithIndent(indent string) *Formatter {
	f.indent = indent
	return f
}

func (f *Formatter) writeIndentIfNeeded() {
	if !f.lineStart {
		return
	}
	for i := 0; i < f.indentLevel; i++ {
		f.builder.WriteString(f.indent)
	}
	f.lineStart = false
}

func (f *Formatter) WriteString(s string) {
	for i := 0; i < len(s); i++ {
		f.WriteByte(s[i])
	}
}

func (f *Formatter) WriteByte(b byte) {
	if f.mode == FormatModeBeautify {
		if b == newline {
			f.builder.WriteByte(newline)
			f.lineStart = true
			return
		}
		f.writeIndentIfNeeded()
		f.builder.WriteByte(b)
	} else {
		f.builder.WriteByte(b)
	}
}

func (f *Formatter) WriteExpr(expr Expr) {
	if expr == nil {
		return
	}
	expr.FormatSQL(f)
}

func (f *Formatter) NewLine() {
	if f.mode != FormatModeBeautify {
		return
	}
	f.WriteByte(newline)
}

func (f *Formatter) Break() {
	if f.mode == FormatModeBeautify {
		f.NewLine()
		return
	}
	f.WriteByte(whitespace)
}

func (f *Formatter) Indent() {
	f.indentLevel++
}

func (f *Formatter) Dedent() {
	if f.indentLevel > 0 {
		f.indentLevel--
	}
}

func (f *Formatter) String() string {
	return f.builder.String()
}

// Format renders an expression into compact SQL.
func Format(expr Expr) string {
	formatter := NewFormatter()
	formatter.WriteExpr(expr)
	return formatter.String()
}

func (p *BinaryOperation) isLogicalOp() bool {
	switch p.Operation {
	case TokenKind(KeywordAnd), TokenKind(KeywordOr):
		return true
	default:
		return p.HasGlobal || p.HasNot
	}
}

func isLogicalBinaryOp(expr Expr) bool {
	if bin, ok := expr.(*BinaryOperation); ok {
		return bin.isLogicalOp()
	}
	return false
}

func (p *BinaryOperation) writeLogicalOperand(formatter *Formatter, expr Expr) {
	if isLogicalBinaryOp(expr) {
		formatter.WriteExpr(expr)
	} else {
		formatter.Indent()
		formatter.WriteExpr(expr)
		formatter.Dedent()
	}
}

func (p *BinaryOperation) FormatSQL(formatter *Formatter) {
	if p.isLogicalOp() && formatter.mode == FormatModeBeautify {
		p.writeLogicalOperand(formatter, p.LeftExpr)
		formatter.NewLine()
		if p.HasNot {
			formatter.WriteString("NOT ")
		} else if p.HasGlobal {
			formatter.WriteString("GLOBAL ")
		}
		formatter.WriteString(string(p.Operation))
		formatter.NewLine()
		p.writeLogicalOperand(formatter, p.RightExpr)
		return
	}
	formatter.WriteExpr(p.LeftExpr)
	if p.Operation != TokenKindDash {
		formatter.WriteByte(whitespace)
	}
	if p.HasNot {
		formatter.WriteString("NOT ")
	} else if p.HasGlobal {
		formatter.WriteString("GLOBAL ")
	}
	formatter.WriteString(string(p.Operation))
	if p.Operation != TokenKindDash {
		formatter.WriteByte(whitespace)
	}
	formatter.WriteExpr(p.RightExpr)
}

func (a *AliasExpr) FormatSQL(formatter *Formatter) {
	if _, isSelect := a.Expr.(*SelectQuery); isSelect {
		formatter.WriteByte('(')
		formatter.WriteExpr(a.Expr)
		formatter.WriteByte(')')
	} else {
		formatter.WriteExpr(a.Expr)
	}
	formatter.WriteString(" AS ")
	formatter.WriteExpr(a.Alias)
}

func (a *AlterRole) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ALTER ROLE ")
	if a.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	for i, roleRenamePair := range a.RoleRenamePairs {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(roleRenamePair)
	}
	if len(a.Settings) > 0 {
		formatter.WriteString(" SETTINGS ")
		for i, setting := range a.Settings {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(setting)
		}
	}
}

func (a *AlterTable) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ALTER TABLE ")
	formatter.WriteExpr(a.TableIdentifier)
	if a.OnCluster != nil {
		formatter.Break()
		formatter.WriteExpr(a.OnCluster)
	}
	for i, expr := range a.AlterExprs {
		formatter.Break()
		formatter.WriteExpr(expr)
		if i != len(a.AlterExprs)-1 {
			formatter.WriteString(",")
		}
	}
}

func (a *AlterTableAddColumn) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ADD COLUMN ")
	if a.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(a.Column)
	if a.After != nil {
		formatter.WriteString(" AFTER ")
		formatter.WriteExpr(a.After)
	}
	if a.Settings != nil {
		formatter.Break()
		formatter.WriteExpr(a.Settings)
	}
}

func (a *AlterTableAddIndex) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ADD ")
	formatter.WriteExpr(a.Index)
	if a.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	if a.After != nil {
		formatter.WriteString(" AFTER ")
		formatter.WriteExpr(a.After)
	}
}

func (a *AlterTableAddProjection) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ADD PROJECTION ")
	if a.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(a.TableProjection)
	if a.After != nil {
		formatter.WriteString(" AFTER ")
		formatter.WriteExpr(a.After)
	}
}

func (a *AlterTableAttachPartition) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ATTACH ")
	formatter.WriteExpr(a.Partition)
	if a.From != nil {
		formatter.WriteString(" FROM ")
		formatter.WriteExpr(a.From)
	}
}

func (a *AlterTableClearColumn) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CLEAR COLUMN ")
	if a.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(a.ColumnName)
	if a.PartitionExpr != nil {
		formatter.WriteString(" IN ")
		formatter.WriteExpr(a.PartitionExpr)
	}

}

func (a *AlterTableClearIndex) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CLEAR INDEX ")
	if a.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(a.IndexName)
	if a.PartitionExpr != nil {
		formatter.WriteString(" IN ")
		formatter.WriteExpr(a.PartitionExpr)
	}

}

func (a *AlterTableClearProjection) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CLEAR PROJECTION ")
	if a.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(a.ProjectionName)
	if a.PartitionExpr != nil {
		formatter.WriteString(" IN ")
		formatter.WriteExpr(a.PartitionExpr)
	}

}

func (a *AlterTableDelete) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DELETE WHERE ")
	formatter.WriteExpr(a.WhereClause)
}

func (a *AlterTableDetachPartition) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DETACH ")
	formatter.WriteExpr(a.Partition)
	if a.Settings != nil {
		formatter.Break()
		formatter.WriteExpr(a.Settings)
	}
}

func (a *AlterTableDropColumn) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP COLUMN ")
	if a.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(a.ColumnName)
}

func (a *AlterTableDropIndex) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP INDEX ")
	formatter.WriteExpr(a.IndexName)
	if a.IfExists {
		formatter.WriteString(" IF EXISTS")
	}
}

func (a *AlterTableDropPartition) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP ")
	if a.HasDetached {
		formatter.WriteString("DETACHED ")
	}
	formatter.WriteExpr(a.Partition)
	if a.Settings != nil {
		formatter.Break()
		formatter.WriteExpr(a.Settings)
	}
}

func (a *AlterTableDropProjection) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP PROJECTION ")
	formatter.WriteExpr(a.ProjectionName)
	if a.IfExists {
		formatter.WriteString(" IF EXISTS")
	}
}

func (a *AlterTableFreezePartition) FormatSQL(formatter *Formatter) {
	formatter.WriteString("FREEZE")
	if a.Partition != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(a.Partition)
	}
}

func (a *AlterTableMaterializeIndex) FormatSQL(formatter *Formatter) {
	formatter.WriteString("MATERIALIZE INDEX")

	if a.IfExists {
		formatter.WriteString(" IF EXISTS")
	}
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(a.IndexName)
	if a.Partition != nil {
		formatter.WriteString(" IN ")
		formatter.WriteExpr(a.Partition)
	}
}

func (a *AlterTableMaterializeProjection) FormatSQL(formatter *Formatter) {
	formatter.WriteString("MATERIALIZE PROJECTION")

	if a.IfExists {
		formatter.WriteString(" IF EXISTS")
	}
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(a.ProjectionName)
	if a.Partition != nil {
		formatter.WriteString(" IN ")
		formatter.WriteExpr(a.Partition)
	}
}

func (a *AlterTableModifyColumn) FormatSQL(formatter *Formatter) {
	formatter.WriteString("MODIFY COLUMN ")
	if a.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(a.Column)
	if a.RemovePropertyType != nil {
		formatter.WriteExpr(a.RemovePropertyType)
	}
}

func (a *AlterTableModifyQuery) FormatSQL(formatter *Formatter) {
	formatter.WriteString("MODIFY QUERY ")
	formatter.WriteExpr(a.SelectExpr)
}

func (a *AlterTableModifySetting) FormatSQL(formatter *Formatter) {
	formatter.WriteString("MODIFY SETTING ")
	for i, setting := range a.Settings {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(setting)
	}
}

func (a *AlterTableModifyTTL) FormatSQL(formatter *Formatter) {
	formatter.WriteString("MODIFY ")
	formatter.WriteExpr(a.TTL)
}

func (a *AlterTableRemoveTTL) FormatSQL(formatter *Formatter) {
	formatter.WriteString("REMOVE TTL")
}

func (a *AlterTableRenameColumn) FormatSQL(formatter *Formatter) {
	formatter.WriteString("RENAME COLUMN ")
	if a.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(a.OldColumnName)
	formatter.WriteString(" TO ")
	formatter.WriteExpr(a.NewColumnName)
}

func (a *AlterTableReplacePartition) FormatSQL(formatter *Formatter) {
	formatter.WriteString("REPLACE ")
	formatter.WriteExpr(a.Partition)
	formatter.WriteString(" FROM ")
	formatter.WriteExpr(a.Table)
}

func (a *AlterTableResetSetting) FormatSQL(formatter *Formatter) {
	formatter.WriteString("RESET SETTING ")
	for i, setting := range a.Settings {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(setting)
	}
}

func (a *AlterTableUpdate) FormatSQL(formatter *Formatter) {
	formatter.WriteString("UPDATE ")
	for i, assignment := range a.Assignments {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(assignment)
	}
	if a.InPartition != nil {
		formatter.WriteString(" IN ")
		formatter.WriteExpr(a.InPartition)
	}
	formatter.WriteString(" WHERE ")
	formatter.WriteExpr(a.WhereClause)
}

func (a *ArrayParamList) FormatSQL(formatter *Formatter) {
	formatter.WriteString("[")
	for i, item := range a.Items.Items {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(item)
	}
	formatter.WriteString("]")
}

func (v *AssignmentValues) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('(')
	for i, value := range v.Values {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(value)
	}
	formatter.WriteByte(')')
}

func (a *AuthenticationClause) FormatSQL(formatter *Formatter) {
	if a.NotIdentified {
		formatter.WriteString("NOT IDENTIFIED")
		return
	}
	formatter.WriteString("IDENTIFIED")
	if a.AuthType != "" {
		formatter.WriteString(" WITH ")
		formatter.WriteString(a.AuthType)
	}
	if a.AuthValue != nil {
		formatter.WriteString(" BY ")
		formatter.WriteExpr(a.AuthValue)
	}
	if a.LdapServer != nil {
		formatter.WriteString(" WITH ldap SERVER ")
		formatter.WriteExpr(a.LdapServer)
	}
	if a.IsKerberos {
		formatter.WriteString(" WITH kerberos")
		if a.KerberosRealm != nil && a.KerberosRealm.Literal != "" {
			formatter.WriteString(" REALM ")
			formatter.WriteExpr(a.KerberosRealm)
		}
	}
}

func (f *BetweenClause) FormatSQL(formatter *Formatter) {
	if f.Expr != nil {
		formatter.WriteExpr(f.Expr)
		formatter.WriteString(" BETWEEN ")
	} else {
		formatter.WriteString("BETWEEN ")
	}
	formatter.WriteExpr(f.Between)
	formatter.WriteString(" AND ")
	formatter.WriteExpr(f.And)
}

func (b *BoolLiteral) FormatSQL(formatter *Formatter) {
	formatter.WriteString(b.Literal)
}

func (c *CTEStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(c.Expr)
	formatter.WriteString(" AS ")
	if _, isSelect := c.Alias.(*SelectQuery); isSelect {
		formatter.WriteByte('(')
		formatter.WriteExpr(c.Alias)
		formatter.WriteByte(')')
	} else {
		formatter.WriteExpr(c.Alias)
	}
}

func (c *CaseExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CASE")
	if c.Expr != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Expr)
	}
	formatter.Indent()
	for _, when := range c.Whens {
		formatter.Break()
		formatter.WriteExpr(when)
	}
	if c.Else != nil {
		formatter.Break()
		formatter.WriteString("ELSE ")
		formatter.WriteExpr(c.Else)
	}
	formatter.Dedent()
	formatter.Break()
	formatter.WriteString("END")
}

func (c *CastExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CAST(")
	formatter.WriteExpr(c.Expr)
	if c.Separator == "," {
		formatter.WriteString(", ")
	} else {
		formatter.WriteString(" AS ")
	}
	formatter.WriteExpr(c.AsType)
	formatter.WriteByte(')')
}

func (c *CheckStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CHECK TABLE ")
	formatter.WriteExpr(c.Table)
	if c.Partition != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Partition)
	}
}

func (o *ClusterClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ON CLUSTER ")
	formatter.WriteExpr(o.Expr)
}

func (c *ColumnArgList) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('(')
	for i, item := range c.Items {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(item)
	}
	formatter.WriteByte(')')
}

func (c *ColumnDef) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(c.Name)
	if c.Type != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Type)
	}
	if c.NotNull != nil {
		formatter.WriteString(" NOT NULL")
	} else if c.Nullable != nil {
		formatter.WriteString(" NULL")
	}
	if c.DefaultExpr != nil {
		formatter.WriteString(" DEFAULT ")
		formatter.WriteExpr(c.DefaultExpr)
	}
	if c.MaterializedExpr != nil {
		formatter.WriteString(" MATERIALIZED ")
		formatter.WriteExpr(c.MaterializedExpr)
	}
	if c.AliasExpr != nil {
		formatter.WriteString(" ALIAS ")
		formatter.WriteExpr(c.AliasExpr)
	}
	if c.Codec != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Codec)
	}
	if c.TTL != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.TTL)
	}
	if c.Comment != nil {
		formatter.WriteString(" COMMENT ")
		formatter.WriteExpr(c.Comment)
	}
}

func (c *ColumnExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(c.Expr)
	if c.Alias != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(c.Alias)
	}
}

func (c *ColumnExprList) FormatSQL(formatter *Formatter) {
	if c.HasDistinct {
		formatter.WriteString("DISTINCT ")
	}
	for i, item := range c.Items {
		formatter.WriteExpr(item)
		if i != len(c.Items)-1 {
			formatter.WriteString(", ")
		}
	}
}

func (c *ColumnNamesExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('(')
	for i, column := range c.ColumnNames {
		if i > 0 {
			formatter.WriteString(", ")
		}
		columnExpr := column
		formatter.WriteExpr(&columnExpr)
	}
	formatter.WriteByte(')')
}

func (c *ColumnTypeExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(c.Name)
}

func (c *ComplexType) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(c.Name)
	formatter.WriteByte('(')
	for i, param := range c.Params {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(param)
	}
	formatter.WriteByte(')')
}

func (c *CompressionCodec) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CODEC(")
	if c.Type != nil {
		formatter.WriteExpr(c.Type)
		if c.TypeLevel != nil {
			formatter.WriteByte('(')
			formatter.WriteExpr(c.TypeLevel)
			formatter.WriteByte(')')
		}
		formatter.WriteByte(',')
		formatter.WriteByte(whitespace)
	}
	if c.Name != nil {
		formatter.WriteExpr(c.Name)
		if c.Level != nil {
			formatter.WriteByte('(')
			formatter.WriteExpr(c.Level)
			formatter.WriteByte(')')
		}
	}
	formatter.WriteByte(')')
}

func (c *ConstraintClause) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(c.Constraint)
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(c.Expr)
}

func (c *CreateDatabase) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE DATABASE ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.Name)
	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}
	if c.Engine != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Engine)
	}
	if c.Comment != nil {
		formatter.WriteString(" COMMENT ")
		formatter.WriteExpr(c.Comment)
	}
}

func (c *CreateDictionary) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE ")
	if c.OrReplace {
		formatter.WriteString("OR REPLACE ")
	}
	formatter.WriteString("DICTIONARY ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.Name)

	if c.UUID != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.UUID)
	}

	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}

	if c.Schema != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Schema)
	}

	if c.Engine != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Engine)
	}

	if c.Comment != nil {
		formatter.WriteString(" COMMENT ")
		formatter.WriteExpr(c.Comment)
	}

}

func (c *CreateFunction) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE")
	if c.OrReplace {
		formatter.WriteString(" OR REPLACE")
	}
	formatter.WriteString(" FUNCTION ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.FunctionName)
	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}
	formatter.WriteString(" AS ")
	formatter.WriteExpr(c.Params)
	formatter.WriteString(" -> ")
	formatter.WriteExpr(c.Expr)
}

func (c *CreateLiveView) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE LIVE VIEW ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.Name)

	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}

	if c.WithTimeout != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.WithTimeout)
	}

	if c.Destination != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Destination)
	}

	if c.TableSchema != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.TableSchema)
	}

	if c.SubQuery != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(c.SubQuery)
	}

}

func (c *CreateMaterializedView) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE MATERIALIZED VIEW ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.Name)
	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}
	if c.Refresh != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Refresh)
	}
	if c.RandomizeFor != nil {
		formatter.WriteString(" RANDOMIZE FOR ")
		formatter.WriteExpr(c.RandomizeFor)
	}
	if c.DependsOn != nil {
		formatter.WriteString(" DEPENDS ON ")
		for i, dep := range c.DependsOn {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(dep)
		}
	}
	if c.Settings != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Settings)
	}
	if c.HasAppend {
		if c.Settings != nil {
			formatter.Break()
		} else {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteString("APPEND")
	}
	if c.Engine != nil {
		formatter.WriteExpr(c.Engine)
	}
	if c.Destination != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.Destination)
		if c.Destination.TableSchema != nil {
			formatter.WriteByte(whitespace)
			formatter.WriteExpr(c.Destination.TableSchema)
		}
	}
	if c.HasEmpty {
		formatter.WriteString(" EMPTY")
	}
	if c.Definer != nil {
		formatter.WriteString(" DEFINER = ")
		formatter.WriteExpr(c.Definer)
	}
	if c.SQLSecurity != "" {
		formatter.WriteString(" SQL SECURITY ")
		formatter.WriteString(c.SQLSecurity)
	}
	if c.Populate {
		formatter.WriteString(" POPULATE")
	}
	if c.SubQuery != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(c.SubQuery)
	}
	if c.Comment != nil {
		formatter.WriteString(" COMMENT ")
		formatter.WriteExpr(c.Comment)
	}
}

func (c *CreateNamedCollection) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE NAMED COLLECTION ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.Name)
	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}
	formatter.WriteString(" AS ")
	for i, param := range c.Params {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(param)
	}
}

func (c *CreateRole) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE ROLE ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	if c.OrReplace {
		formatter.WriteString("OR REPLACE ")
	}
	for i, roleName := range c.RoleNames {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(roleName)
	}
	if c.AccessStorageType != nil {
		formatter.WriteString(" IN ")
		formatter.WriteExpr(c.AccessStorageType)
	}
	if len(c.Settings) > 0 {
		formatter.WriteString(" SETTINGS ")
		for i, setting := range c.Settings {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(setting)
		}
	}
}

func (c *CreateTable) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE")
	if c.OrReplace {
		formatter.WriteString(" OR REPLACE")
	}
	if c.HasTemporary {
		formatter.WriteString(" TEMPORARY")
	}
	formatter.WriteString(" TABLE ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.Name)
	if c.UUID != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.UUID)
	}
	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}

	if c.TableSchema != nil {
		formatter.Break()
		formatter.WriteExpr(c.TableSchema)
	}
	if c.Engine != nil {
		formatter.WriteExpr(c.Engine)
	}
	if c.SubQuery != nil {
		formatter.Break()
		formatter.WriteString("AS ")
		formatter.WriteExpr(c.SubQuery)
	}
	if c.TableFunction != nil {
		formatter.Break()
		formatter.WriteString("AS ")
		formatter.WriteExpr(c.TableFunction)
	}
	if c.Comment != nil {
		formatter.Break()
		formatter.WriteString("COMMENT ")
		formatter.WriteExpr(c.Comment)
	}
}

func (c *CreateUser) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE USER ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	if c.OrReplace {
		formatter.WriteString("OR REPLACE ")
	}
	for i, userName := range c.UserNames {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(userName)
	}
	if c.Authentication != nil {
		formatter.Break()
		formatter.WriteExpr(c.Authentication)
	}
	if len(c.Hosts) > 0 {
		formatter.Break()
		for i, host := range c.Hosts {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(host)
		}
	}
	if c.DefaultRole != nil {
		formatter.Break()
		formatter.WriteExpr(c.DefaultRole)
	}
	if c.DefaultDatabase != nil {
		formatter.Break()
		formatter.WriteString("DEFAULT DATABASE ")
		formatter.WriteExpr(c.DefaultDatabase)
	} else if c.DefaultDbNone {
		formatter.Break()
		formatter.WriteString("DEFAULT DATABASE NONE")
	}
	if c.Grantees != nil {
		formatter.Break()
		formatter.WriteExpr(c.Grantees)
	}
	if len(c.Settings) > 0 {
		formatter.Break()
		formatter.WriteString("SETTINGS")
		formatter.Indent()
		for i, setting := range c.Settings {
			formatter.Break()
			formatter.WriteExpr(setting)
			if i < len(c.Settings)-1 {
				formatter.WriteString(",")
			}
		}
		formatter.Dedent()
	}
}

func (c *CreateView) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CREATE")
	if c.OrReplace {
		formatter.WriteString(" OR REPLACE")
	}
	formatter.WriteString(" VIEW ")
	if c.IfNotExists {
		formatter.WriteString("IF NOT EXISTS ")
	}
	formatter.WriteExpr(c.Name)
	if c.UUID != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.UUID)
	}

	if c.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.OnCluster)
	}

	if c.TableSchema != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(c.TableSchema)
	}

	if c.SubQuery != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(c.SubQuery)
	}
}

func (d *DeduplicateClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString(" DEDUPLICATE")
	if d.By != nil {
		formatter.WriteString(" BY ")
		formatter.WriteExpr(d.By)
	}
	if d.Except != nil {
		formatter.WriteString(" EXCEPT ")
		formatter.WriteExpr(d.Except)
	}
}

func (d *DefaultRoleClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DEFAULT ROLE ")
	if d.None {
		formatter.WriteString("NONE")
	} else {
		for i, role := range d.Roles {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(role)
		}
	}
}

func (d *DeleteClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DELETE FROM ")
	formatter.WriteExpr(d.Table)
	if d.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(d.OnCluster)
	}
	if d.WhereExpr != nil {
		formatter.WriteString(" WHERE ")
		formatter.WriteExpr(d.WhereExpr)
	}
}

func (d *DescribeStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DESCRIBE ")
	if d.DescribeType != "" {
		formatter.WriteString(d.DescribeType)
		formatter.WriteByte(whitespace)
	}
	formatter.WriteExpr(d.Target)
}

func (d *DestinationClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("TO ")
	formatter.WriteExpr(d.TableIdentifier)
}

func (d *DictionaryArgExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(d.Name)
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(d.Value)
}

func (d *DictionaryAttribute) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(d.Name)
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(d.Type)

	if d.Default != nil {
		formatter.WriteString(" DEFAULT ")
		formatter.WriteExpr(d.Default)
	}

	if d.Expression != nil {
		formatter.WriteString(" EXPRESSION ")
		formatter.WriteExpr(d.Expression)
	}

	if d.Hierarchical {
		formatter.WriteString(" HIERARCHICAL")
	}

	if d.Injective {
		formatter.WriteString(" INJECTIVE")
	}

	if d.IsObjectId {
		formatter.WriteString(" IS_OBJECT_ID")
	}

}

func (d *DictionaryEngineClause) FormatSQL(formatter *Formatter) {
	paddingSpace := false
	if d.PrimaryKey != nil {
		formatter.WriteExpr(d.PrimaryKey)
		paddingSpace = true
	}
	if d.Source != nil {
		if paddingSpace {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(d.Source)
		paddingSpace = true
	}
	if d.Lifetime != nil {
		if paddingSpace {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(d.Lifetime)
		paddingSpace = true
	}
	if d.Layout != nil {
		if paddingSpace {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(d.Layout)
		paddingSpace = true
	}
	if d.Range != nil {
		if paddingSpace {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(d.Range)
		paddingSpace = true
	}
	if d.Settings != nil {
		if paddingSpace {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteString("SETTINGS(")
		for i, item := range d.Settings.Items {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(item)
		}
		formatter.WriteString(")")
	}
}

func (d *DictionaryLayoutClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("LAYOUT(")
	formatter.WriteExpr(d.Layout)
	formatter.WriteString("(")
	for i, arg := range d.Args {
		if i > 0 {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(arg)
	}
	formatter.WriteString("))")
}

func (d *DictionaryLifetimeClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("LIFETIME(")
	if d.Value != nil {
		formatter.WriteExpr(d.Value)
	} else if d.Min != nil && d.Max != nil {
		formatter.WriteString("MIN ")
		formatter.WriteExpr(d.Min)
		formatter.WriteString(" MAX ")
		formatter.WriteExpr(d.Max)
	}
	formatter.WriteString(")")
}

func (d *DictionaryPrimaryKeyClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("PRIMARY KEY ")
	formatter.WriteExpr(d.Keys)
}

func (d *DictionaryRangeClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("RANGE(")
	formatter.WriteString("MIN ")
	if d.Min != nil {
		formatter.WriteExpr(d.Min)
	}
	if d.Max != nil {
		formatter.WriteString(" MAX ")
		formatter.WriteExpr(d.Max)
	}
	formatter.WriteString(")")
}

func (d *DictionarySchemaClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("(")
	for i, attr := range d.Attributes {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(attr)
	}
	formatter.WriteString(")")
}

func (d *DictionarySourceClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SOURCE(")
	formatter.WriteExpr(d.Source)
	formatter.WriteString("(")
	for i, arg := range d.Args {
		if i > 0 {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(arg)
	}
	formatter.WriteString("))")
}

func (s *DistinctOn) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ON (")
	for i, ident := range s.Idents {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(ident)
	}
	formatter.WriteByte(')')
}

func (d *DropDatabase) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP DATABASE ")
	if d.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(d.Name)
	if d.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(d.OnCluster)
	}
}

func (d *DropStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP ")
	if d.IsTemporary {
		formatter.WriteString("TEMPORARY ")
	}
	formatter.WriteString(d.DropTarget + " ")
	if d.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(d.Name)
	if d.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(d.OnCluster)
	}
	if len(d.Modifier) != 0 {
		formatter.WriteString(" " + d.Modifier)
	}
}

func (d *DropUserOrRole) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP " + d.Target + " ")
	if d.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	for i, name := range d.Names {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(name)
	}
	if len(d.Modifier) != 0 {
		formatter.WriteString(" " + d.Modifier)
	}
	if d.From != nil {
		formatter.WriteString(" FROM ")
		formatter.WriteExpr(d.From)
	}
}

func (e *EngineExpr) FormatSQL(formatter *Formatter) {
	formatter.Break()
	formatter.WriteString("ENGINE = ")
	formatter.WriteString(e.Name)
	if e.Params != nil {
		formatter.WriteExpr(e.Params)
	}
	if e.OrderBy != nil {
		formatter.Break()
		formatter.WriteExpr(e.OrderBy)
	}
	if e.PartitionBy != nil {
		formatter.Break()
		formatter.WriteExpr(e.PartitionBy)
	}
	if e.PrimaryKey != nil {
		formatter.Break()
		formatter.WriteExpr(e.PrimaryKey)
	}
	if e.SampleBy != nil {
		formatter.Break()
		formatter.WriteExpr(e.SampleBy)
	}
	if e.TTL != nil {
		formatter.Break()
		formatter.WriteExpr(e.TTL)
	}
	if e.Settings != nil {
		formatter.Break()
		formatter.WriteExpr(e.Settings)
	}
}

func (e *EnumType) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(e.Name)
	formatter.WriteByte('(')
	for i, enum := range e.Values {
		if i > 0 {
			formatter.WriteString(", ")
		}
		enumExpr := enum
		formatter.WriteExpr(&enumExpr)
	}
	formatter.WriteByte(')')
}

func (e *EnumValue) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(e.Name)
	formatter.WriteByte('=')
	formatter.WriteExpr(e.Value)
}

func (e *ExplainStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("EXPLAIN ")
	formatter.WriteString(e.Type)
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(e.Statement)
}

func (e *ExtractExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("EXTRACT(")
	for i, param := range e.Parameters {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(param)
	}
	formatter.WriteByte(')')
}

func (f *Fill) FormatSQL(formatter *Formatter) {
	formatter.WriteString("WITH FILL")
	if f.From != nil {
		formatter.WriteString(" FROM ")
		formatter.WriteExpr(f.From)
	}
	if f.To != nil {
		formatter.WriteString(" TO ")
		formatter.WriteExpr(f.To)
	}
	if f.Step != nil {
		formatter.WriteString(" STEP ")
		formatter.WriteExpr(f.Step)
	}
	if f.Staleness != nil {
		formatter.WriteString(" STALENESS ")
		formatter.WriteExpr(f.Staleness)
	}
}

func (f *FormatClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("FORMAT ")
	formatter.WriteExpr(f.Format)
}

func (f *FromClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("FROM")
	formatter.Indent()
	formatter.Break()
	formatter.WriteExpr(f.Expr)
	formatter.Dedent()
}

func (f *FunctionExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(f.Name)
	formatter.WriteExpr(f.Params)
}

func (g *GlobalInOperation) FormatSQL(formatter *Formatter) {
	formatter.WriteString("GLOBAL ")
	formatter.WriteExpr(g.Expr)
}

func (g *GrantPrivilegeStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("GRANT ")
	if g.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(g.OnCluster)
	}
	for i, privilege := range g.Privileges {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(privilege)
	}
	formatter.WriteString(" ON ")
	formatter.WriteExpr(g.On)
	formatter.WriteString(" TO ")
	for i, role := range g.To {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(role)
	}
	for _, option := range g.WithOptions {
		formatter.WriteString(" WITH " + option + " OPTION")
	}

}

func (g *GranteesClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("GRANTEES ")
	if g.Any {
		formatter.WriteString("ANY")
	} else if g.None {
		formatter.WriteString("NONE")
	} else {
		for i, grantee := range g.Grantees {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(grantee)
		}
	}
	if len(g.ExceptUsers) > 0 {
		formatter.Break()
		formatter.WriteString("EXCEPT ")
		for i, except := range g.ExceptUsers {
			if i > 0 {
				formatter.WriteString(", ")
			}
			formatter.WriteExpr(except)
		}
	}
}

func (g *GroupByClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("GROUP BY")

	formatter.Indent()
	defer formatter.Dedent()
	if g.AggregateType != "" {
		formatter.Break()
		formatter.WriteString(g.AggregateType)
	}
	if g.Expr != nil {
		if g.AggregateType == "" {
			formatter.Break()
		}
		formatter.WriteExpr(g.Expr)
	}
	if g.WithCube {
		formatter.Break()
		formatter.WriteString("WITH CUBE")
	}
	if g.WithRollup {
		formatter.Break()
		formatter.WriteString("WITH ROLLUP")
	}
	if g.WithTotals {
		formatter.Break()
		formatter.WriteString("WITH TOTALS")
	}
}

func (h *HavingClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("HAVING ")
	formatter.WriteExpr(h.Expr)
}

func (h *HostClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("HOST ")
	formatter.WriteString(h.HostType)
	if h.HostValue != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(h.HostValue)
	}
}

func (i *Ident) FormatSQL(formatter *Formatter) {
	switch i.QuoteType {
	case BackTicks:
		formatter.WriteByte('`')
		formatter.WriteString(i.Name)
		formatter.WriteByte('`')
	case DoubleQuote:
		formatter.WriteByte('"')
		formatter.WriteString(i.Name)
		formatter.WriteByte('"')
	case SingleQuote:
		formatter.WriteByte('\'')
		formatter.WriteString(i.Name)
		formatter.WriteByte('\'')
	default:
		formatter.WriteString(i.Name)
	}
}

func (i *IndexOperation) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(i.Object)
	formatter.WriteString(string(i.Operation))
	formatter.WriteExpr(i.Index)
}

func (i *InsertStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("INSERT INTO ")
	if i.HasTableKeyword {
		formatter.WriteString("TABLE ")
	}
	formatter.WriteExpr(i.Table)
	if i.ColumnNames != nil {
		formatter.Break()
		formatter.Indent()
		formatter.WriteExpr(i.ColumnNames)
		formatter.Dedent()
	}
	if i.Format != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(i.Format)
	}

	if i.SelectExpr != nil {
		formatter.Break()
		formatter.WriteExpr(i.SelectExpr)
	} else if len(i.Values) > 0 {
		formatter.Break()
		formatter.WriteString("VALUES")
		formatter.Indent()
		for j, value := range i.Values {
			formatter.Break()
			formatter.WriteExpr(value)
			if j != len(i.Values)-1 {
				formatter.WriteByte(',')
			}
		}
		formatter.Dedent()
	}
}

func (i *InterpolateClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("INTERPOLATE")
	if len(i.Items) > 0 {
		formatter.WriteString(" (")
		for idx, item := range i.Items {
			formatter.WriteExpr(item)
			if idx != len(i.Items)-1 {
				formatter.WriteString(", ")
			}
		}
		formatter.WriteByte(')')
	}
}

func (i *InterpolateItem) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(i.Column)
	if i.Expr != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(i.Expr)
	}
}

func (i *IntervalExpr) FormatSQL(formatter *Formatter) {
	if i.IntervalPos != 0 {
		formatter.WriteString("INTERVAL ")
	}
	formatter.WriteExpr(i.Expr)
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(i.Unit)
}

func (i *IntervalFrom) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(i.Interval)
	formatter.WriteString(" FROM ")
	formatter.WriteExpr(i.FromExpr)
}

func (n *IsNotNullExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(n.Expr)
	formatter.WriteString(" IS NOT NULL")
}

func (n *IsNullExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(n.Expr)
	formatter.WriteString(" IS NULL")
}

func (j *JSONPath) FormatSQL(formatter *Formatter) {
	for i, ident := range j.Idents {
		if i > 0 {
			formatter.WriteByte('.')
		}
		formatter.WriteExpr(ident)
	}
}

func (j *JSONOption) FormatSQL(formatter *Formatter) {
	wroteAny := false
	if j.SkipPath != nil {
		formatter.WriteString("SKIP ")
		j.SkipPath.FormatSQL(formatter)
		wroteAny = true
	}
	if j.SkipRegex != nil {
		formatter.WriteString(" SKIP REGEXP ")
		formatter.WriteExpr(j.SkipRegex)
		wroteAny = true
	}
	if j.MaxDynamicPaths != nil {
		formatter.WriteString("max_dynamic_paths")
		formatter.WriteByte('=')
		formatter.WriteExpr(j.MaxDynamicPaths)
		wroteAny = true
	}
	if j.MaxDynamicTypes != nil {
		formatter.WriteString("max_dynamic_types")
		formatter.WriteByte('=')
		formatter.WriteExpr(j.MaxDynamicTypes)
		wroteAny = true
	}
	if j.Column != nil && j.Column.Path != nil && j.Column.Type != nil {
		// Add a leading space if there is already content.
		if wroteAny {
			formatter.WriteByte(whitespace)
		}
		j.Column.Path.FormatSQL(formatter)
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(j.Column.Type)
	}
}

func (j *JSONOptions) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('(')
	// Ensure stable, readable ordering:
	// 1) numeric options (max_dynamic_*), 2) type-hint items, 3) skip options (SKIP, SKIP REGEXP)
	// Preserve original relative order within each group.
	numericOptionItems := make([]*JSONOption, 0, len(j.Items))
	columnItems := make([]*JSONOption, 0, len(j.Items))
	skipOptionItems := make([]*JSONOption, 0, len(j.Items))
	for _, item := range j.Items {
		if item.MaxDynamicPaths != nil || item.MaxDynamicTypes != nil {
			numericOptionItems = append(numericOptionItems, item)
			continue
		}
		if item.Column != nil {
			columnItems = append(columnItems, item)
			continue
		}
		if item.SkipPath != nil || item.SkipRegex != nil {
			skipOptionItems = append(skipOptionItems, item)
			continue
		}
		// Fallback: treat as numeric option to avoid dropping unknown future fields.
		numericOptionItems = append(numericOptionItems, item)
	}

	wroteItem := false
	writeItems := func(items []*JSONOption) {
		for _, item := range items {
			if wroteItem {
				formatter.WriteString(", ")
			}
			item.FormatSQL(formatter)
			wroteItem = true
		}
	}

	writeItems(numericOptionItems)
	writeItems(columnItems)
	writeItems(skipOptionItems)
	formatter.WriteByte(')')
}

func (j *JSONType) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(j.Name)
	if j.Options != nil {
		j.Options.FormatSQL(formatter)
	}
}

func (j *JoinConstraintClause) FormatSQL(formatter *Formatter) {
	if j.On != nil {
		formatter.WriteString("ON ")
		formatter.WriteExpr(j.On)
	} else {
		formatter.WriteString("USING ")
		formatter.WriteExpr(j.Using)
	}
}

func (j *JoinExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(j.Left)
	if j.Right != nil {
		writeJoinSQL(formatter, j.Right)
	}
}

func writeJoinSQL(formatter *Formatter, expr Expr) {
	joinExpr, ok := expr.(*JoinExpr)
	if !ok {
		formatter.WriteByte(',')
		formatter.WriteExpr(expr)
		return
	}

	if len(joinExpr.Modifiers) == 0 {
		formatter.WriteByte(',')
		formatter.WriteExpr(joinExpr.Left)
	} else {
		formatter.Break()
		formatter.WriteString(strings.Join(joinExpr.Modifiers, " "))
		formatter.Indent()
		formatter.Break()
		formatter.WriteExpr(joinExpr.Left)
		if joinExpr.Constraints != nil {
			formatter.WriteByte(whitespace)
			formatter.WriteExpr(joinExpr.Constraints)
		}
		formatter.Dedent()
	}
	if joinExpr.Right != nil {
		writeJoinSQL(formatter, joinExpr.Right)
	}
}

func (j *JoinTableExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(j.Table)
	if j.SampleRatio != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(j.SampleRatio)
	}
	if j.HasFinal {
		formatter.WriteString(" FINAL")
	}
}

func (l *LimitByClause) FormatSQL(formatter *Formatter) {
	if l.Limit != nil {
		formatter.WriteExpr(l.Limit)
	}
	if l.ByExpr != nil {
		formatter.WriteString(" BY ")
		formatter.WriteExpr(l.ByExpr)
	}
}

func (l *LimitClause) FormatSQL(formatter *Formatter) {
	if l.Limit != nil {
		formatter.WriteString("LIMIT ")
		formatter.WriteExpr(l.Limit)
		if l.Offset != nil {
			formatter.WriteByte(whitespace)
		}
	}
	if l.Offset != nil {
		formatter.WriteString("OFFSET ")
		formatter.WriteExpr(l.Offset)
	}
}

func (m *MapLiteral) FormatSQL(formatter *Formatter) {
	formatter.WriteString("{")

	for i, value := range m.KeyValues {
		if i > 0 {
			formatter.WriteString(", ")
		}
		key := value.Key
		formatter.WriteExpr(&key)
		formatter.WriteString(": ")
		formatter.WriteExpr(value.Value)
	}
	formatter.WriteString("}")
}

func (n *NamedCollectionParam) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(n.Name)
	formatter.WriteString(" = ")
	formatter.WriteExpr(n.Value)
	if n.NotOverridable {
		formatter.WriteString(" NOT OVERRIDABLE")
	} else if n.Overridable {
		formatter.WriteString(" OVERRIDABLE")
	}
}

func (n *NamedParameterExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(n.Name)
	formatter.WriteByte('=')
	formatter.WriteExpr(n.Value)
}

func (n *NegateExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('-')
	formatter.WriteExpr(n.Expr)
}

func (n *NestedIdentifier) FormatSQL(formatter *Formatter) {
	if n.DotIdent != nil {
		formatter.WriteExpr(n.Ident)
		formatter.WriteByte('.')
		formatter.WriteExpr(n.DotIdent)
	} else {
		formatter.WriteExpr(n.Ident)
	}
}

func (n *NestedType) FormatSQL(formatter *Formatter) {
	// on the same level as the column type
	formatter.WriteExpr(n.Name)
	formatter.WriteByte('(')
	for i, column := range n.Columns {
		formatter.WriteExpr(column)
		if i != len(n.Columns)-1 {
			formatter.WriteString(", ")
		}
	}
	// right paren needs to be on the same level as the column
	formatter.WriteByte(')')
}

func (n *NotExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("NOT")
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(n.Expr)
}

func (n *NotNullLiteral) FormatSQL(formatter *Formatter) {
	formatter.WriteString("NOT NULL")
}

func (n *NullLiteral) FormatSQL(formatter *Formatter) {
	formatter.WriteString("NULL")
}

func (n *NumberLiteral) FormatSQL(formatter *Formatter) {
	formatter.WriteString(n.Literal)
}

func (o *ObjectParams) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(o.Object)
	formatter.WriteExpr(o.Params)
}

func (o *OnClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ON ")
	formatter.WriteExpr(o.On)
}

func (o *OperationExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString(strings.ToUpper(string(o.Kind)))
}

func (o *OptimizeStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("OPTIMIZE TABLE ")
	formatter.WriteExpr(o.Table)
	if o.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(o.OnCluster)
	}
	if o.Partition != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(o.Partition)
	}
	if o.HasFinal {
		formatter.WriteString(" FINAL")
	}
	if o.Deduplicate != nil {
		formatter.WriteExpr(o.Deduplicate)
	}
}

func (o *OrderByClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ORDER BY")
	formatter.Indent()
	for i, item := range o.Items {
		if i == 0 {
			formatter.Break()
		} else {
			formatter.WriteByte(',')
			formatter.Break()
		}
		formatter.WriteExpr(item)
	}
	if o.Interpolate != nil {
		formatter.Break()
		formatter.WriteExpr(o.Interpolate)
	}
	formatter.Dedent()
}

func (o *OrderExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(o.Expr)
	if o.Alias != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(o.Alias)
	}
	if o.Direction != OrderDirectionNone {
		formatter.WriteByte(whitespace)
		formatter.WriteString(string(o.Direction))
	}
	if o.Fill != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(o.Fill)
	}
}

func (f *ParamExprList) FormatSQL(formatter *Formatter) {
	formatter.WriteString("(")
	formatter.WriteExpr(f.Items)
	formatter.WriteString(")")
	if f.ColumnArgList != nil {
		formatter.WriteExpr(f.ColumnArgList)
	}
}

func (p *PartitionByClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("PARTITION BY ")
	formatter.WriteExpr(p.Expr)
}

func (p *PartitionClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("PARTITION ")
	if p.ID != nil {
		formatter.WriteExpr(p.ID)
	} else if p.All {
		formatter.WriteString("ALL")
	} else {
		formatter.WriteExpr(p.Expr)
	}
}

func (p *Path) FormatSQL(formatter *Formatter) {
	for i, ident := range p.Fields {
		if i > 0 {
			formatter.WriteByte('.')
		}
		formatter.WriteExpr(ident)
	}
}

func (p *PlaceHolder) FormatSQL(formatter *Formatter) {
	formatter.WriteString(p.Type)
}

func (w *PrewhereClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("PREWHERE ")
	formatter.WriteExpr(w.Expr)
}

func (p *PrimaryKeyClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("PRIMARY KEY ")
	formatter.WriteExpr(p.Expr)
}

func (p *PrivilegeClause) FormatSQL(formatter *Formatter) {
	for i, keyword := range p.Keywords {
		if i > 0 {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteString(keyword)
	}
	if p.Params != nil {
		formatter.WriteExpr(p.Params)
	}
}

func (p *ProjectionOrderByClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("ORDER BY ")
	formatter.WriteExpr(p.Columns)
}

func (p *ProjectionSelectStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("(")
	if p.With != nil {
		formatter.WriteExpr(p.With)
		formatter.WriteByte(whitespace)
	}
	formatter.WriteString("SELECT ")
	formatter.WriteExpr(p.SelectColumns)
	if p.GroupBy != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(p.GroupBy)
	}
	if p.OrderBy != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(p.OrderBy)
	}
	formatter.WriteString(")")
}

func (c *PropertyType) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(c.Name)
}

func (q *QueryParam) FormatSQL(formatter *Formatter) {
	formatter.WriteString("{")
	formatter.WriteExpr(q.Name)
	formatter.WriteString(": ")
	formatter.WriteExpr(q.Type)
	formatter.WriteString("}")
}

func (r *RatioExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(r.Numerator)
	if r.Denominator != nil {
		formatter.WriteString("/")
		formatter.WriteExpr(r.Denominator)
	}
}

func (r *RefreshExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("REFRESH ")
	formatter.WriteString(r.Frequency)
	if r.Interval != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(r.Interval)
	}
	if r.Offset != nil {
		formatter.WriteString(" OFFSET ")
		formatter.WriteExpr(r.Offset)
	}
}

func (a *RemovePropertyType) FormatSQL(formatter *Formatter) {
	formatter.WriteString(" REMOVE ")
	formatter.WriteExpr(a.PropertyType)
}

func (r *RenameStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("RENAME " + r.RenameTarget + " ")
	for i, pair := range r.TargetPairList {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(pair.Old)
		formatter.WriteString(" TO ")
		formatter.WriteExpr(pair.New)
	}
	if r.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(r.OnCluster)
	}
}

func (r *RoleName) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(r.Name)
	if r.Scope != nil {
		formatter.WriteString("@")
		formatter.WriteExpr(r.Scope)
	}
	if r.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(r.OnCluster)
	}
}

func (r *RoleRenamePair) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(r.RoleName)
	if r.NewName != nil {
		formatter.WriteString(" RENAME TO ")
		formatter.WriteExpr(r.NewName)
	}
}

func (r *RoleSetting) FormatSQL(formatter *Formatter) {
	for i, settingPair := range r.SettingPairs {
		if i > 0 {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(settingPair)
	}
	if r.Modifier != nil {
		if len(r.SettingPairs) > 0 {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(r.Modifier)
	}
}

func (s *SampleByClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SAMPLE BY ")
	formatter.WriteExpr(s.Expr)
}

func (s *SampleClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SAMPLE ")
	formatter.WriteExpr(s.Ratio)
	if s.Offset != nil {
		formatter.WriteString(" OFFSET ")
		formatter.WriteExpr(s.Offset)
	}
}

func (s *ScalarType) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(s.Name)
}

func (s *SelectItem) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(s.Expr)
	for _, modifier := range s.Modifiers {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(modifier)
	}
	if s.Alias != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(s.Alias)
	}
}

func (s *SelectQuery) FormatSQL(formatter *Formatter) {
	if s.With != nil {
		formatter.WriteString("WITH")
		formatter.Indent()
		for i, cte := range s.With.CTEs {
			if i == 0 {
				formatter.Break()
			} else {
				formatter.WriteByte(',')
				formatter.Break()
			}
			formatter.WriteExpr(cte)
		}
		formatter.Dedent()
		formatter.Break()
	}
	formatter.WriteString("SELECT")
	if s.HasDistinct {
		formatter.WriteString(" DISTINCT")
		if s.DistinctOn != nil {
			formatter.WriteByte(whitespace)
			formatter.WriteExpr(s.DistinctOn)
		}
	}
	if s.Top != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(s.Top)
	}
	formatter.Indent()
	for i, selectItem := range s.SelectItems {
		if i == 0 {
			formatter.Break()
		} else {
			formatter.WriteByte(',')
			formatter.Break()
		}
		formatter.WriteExpr(selectItem)
	}
	formatter.Dedent()
	if s.From != nil {
		formatter.Break()
		formatter.WriteExpr(s.From)
	}
	if s.Window != nil {
		formatter.Break()
		formatter.WriteExpr(s.Window)
	}
	if s.Prewhere != nil {
		formatter.Break()
		formatter.WriteExpr(s.Prewhere)
	}
	if s.Where != nil {
		formatter.Break()
		formatter.WriteExpr(s.Where)
	}
	if s.GroupBy != nil {
		formatter.Break()
		formatter.WriteExpr(s.GroupBy)
	}
	if s.Having != nil {
		formatter.Break()
		formatter.WriteExpr(s.Having)
	}
	if s.OrderBy != nil {
		formatter.Break()
		formatter.WriteExpr(s.OrderBy)
	}
	if s.LimitBy != nil {
		formatter.Break()
		formatter.WriteExpr(s.LimitBy)
	}
	if s.Limit != nil {
		formatter.Break()
		formatter.WriteExpr(s.Limit)
	}
	if s.Settings != nil {
		formatter.Break()
		formatter.WriteExpr(s.Settings)
	}
	if s.Format != nil {
		formatter.Break()
		formatter.WriteExpr(s.Format)
	}
	if s.UnionAll != nil {
		formatter.Break()
		formatter.WriteString("UNION ALL")
		formatter.Break()
		formatter.WriteExpr(s.UnionAll)
	} else if s.UnionDistinct != nil {
		formatter.Break()
		formatter.WriteString("UNION DISTINCT")
		formatter.Break()
		formatter.WriteExpr(s.UnionDistinct)
	} else if s.Except != nil {
		formatter.Break()
		formatter.WriteString("EXCEPT")
		formatter.Break()
		formatter.WriteExpr(s.Except)
	}
}

func (s *SetStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SET ")
	for i, item := range s.Settings.Items {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(item)
	}
}

func (s *SettingExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(s.Name)
	formatter.WriteByte('=')
	formatter.WriteExpr(s.Expr)
}

func (s *SettingPair) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(s.Name)
	if s.Value != nil {
		if s.Operation == TokenKindSingleEQ {
			formatter.WriteString(string(s.Operation))
		} else {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(s.Value)
	}
}

func (s *SettingsClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SETTINGS")
	formatter.Indent()
	for i, item := range s.Items {
		if i == 0 {
			formatter.Break()
		} else {
			formatter.WriteByte(',')
			formatter.Break()
		}
		formatter.WriteExpr(item)
	}
	formatter.Dedent()
}

func (s *ShowStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SHOW ")
	formatter.WriteString(s.ShowType)
	if s.Target != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(s.Target)
	}

	// Add optional clauses for SHOW DATABASES
	if s.LikeType != "" && s.LikePattern != nil {
		if s.NotLike {
			formatter.WriteString(" NOT ")
		} else {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteString(s.LikeType)
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(s.LikePattern)
	}

	if s.Limit != nil {
		formatter.WriteString(" LIMIT ")
		formatter.WriteExpr(s.Limit)
	}

	if s.OutFile != nil {
		formatter.WriteString(" INTO OUTFILE ")
		formatter.WriteExpr(s.OutFile)
	}

	if s.Format != nil {
		formatter.WriteString(" FORMAT ")
		formatter.WriteExpr(s.Format)
	}

}

func (s *StringLiteral) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('\'')
	formatter.WriteString(s.Literal)
	formatter.WriteByte('\'')
}

func (s *SubQuery) FormatSQL(formatter *Formatter) {
	if s.HasParen {
		formatter.WriteByte('(')
	}
	formatter.WriteExpr(s.Select)
	if s.HasParen {
		formatter.WriteByte(')')
	}
}

func (s *SystemCtrlExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString(s.Command)
	formatter.WriteByte(whitespace)
	formatter.WriteString(s.Type)
	if s.Cluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(s.Cluster)
	}
}

func (s *SystemDropExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("DROP ")
	formatter.WriteString(s.Type)
}

func (s *SystemFlushExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("FLUSH ")
	if s.Logs {
		formatter.WriteString("LOGS")
	} else {
		formatter.WriteExpr(s.Distributed)
	}
}

func (s *SystemReloadExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("RELOAD ")
	formatter.WriteString(s.Type)
	if s.Dictionary != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(s.Dictionary)
	}
}

func (s *SystemStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SYSTEM")
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(s.Expr)
}

func (s *SystemSyncExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString("SYNC ")
	formatter.WriteExpr(s.Cluster)
}

func (t *TTLClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("TTL ")
	for i, item := range t.Items {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(item)
	}
}

func (t *TTLExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(t.Expr)
	if t.Policy != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(t.Policy)
	}
}

func (t *TTLPolicy) FormatSQL(formatter *Formatter) {

	if t.Item != nil {
		formatter.WriteExpr(t.Item)
	}
	if t.Where != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(t.Where)
	}
	if t.GroupBy != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(t.GroupBy)
	}
}

func (t *TTLPolicyRule) FormatSQL(formatter *Formatter) {
	if t.ToVolume != nil {
		formatter.WriteString("TO VOLUME ")
		formatter.WriteExpr(t.ToVolume)
	} else if t.ToDisk != nil {
		formatter.WriteString("TO DISK ")
		formatter.WriteExpr(t.ToDisk)
	} else if t.Action != nil {
		formatter.WriteExpr(t.Action)
	}
}

func (t *TTLPolicyRuleAction) FormatSQL(formatter *Formatter) {
	formatter.WriteString(t.Action)
	if t.Codec != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(t.Codec)
	}
}

func (t *TableArgListExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('(')
	for i, arg := range t.Args {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(arg)
	}
	formatter.WriteByte(')')
}

func (t *TableExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(t.Expr)
	if t.Alias != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(t.Alias)
	}
	if t.HasFinal {
		formatter.WriteString(" FINAL")
	}
}

func (t *TableFunctionExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(t.Name)
	formatter.WriteExpr(t.Args)
}

func (t *TableIdentifier) FormatSQL(formatter *Formatter) {
	if t.Database != nil {
		formatter.WriteExpr(t.Database)
		formatter.WriteByte('.')
	}
	formatter.WriteExpr(t.Table)
}

func (a *TableIndex) FormatSQL(formatter *Formatter) {
	formatter.WriteString("INDEX")
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(a.Name)
	// Add space only if column expression doesn't start with '('
	columnExprStr := Format(a.ColumnExpr)
	if len(columnExprStr) > 0 && columnExprStr[0] != '(' {
		formatter.WriteByte(whitespace)
	}
	formatter.WriteString(columnExprStr)
	formatter.WriteByte(whitespace)
	formatter.WriteString("TYPE")
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(a.ColumnType)
	formatter.WriteByte(whitespace)
	formatter.WriteString("GRANULARITY")
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(a.Granularity)
}

func (t *TableProjection) FormatSQL(formatter *Formatter) {
	if t.IncludeProjectionKeyword {
		formatter.WriteString("PROJECTION ")
	}
	formatter.WriteExpr(t.Identifier)
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(t.Select)
}

func (t *TableSchemaClause) FormatSQL(formatter *Formatter) {
	if len(t.Columns) > 0 {
		formatter.WriteByte('(')
		formatter.Indent()
		for i, column := range t.Columns {
			if i == 0 {
				formatter.NewLine()
			} else {
				formatter.WriteByte(',')
				formatter.Break()
			}
			formatter.WriteExpr(column)
		}
		formatter.Dedent()
		formatter.NewLine()
		formatter.WriteByte(')')
	}
	if t.AliasTable != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(t.AliasTable)
	}
	if t.TableFunction != nil {
		formatter.WriteString(" AS ")
		formatter.WriteExpr(t.TableFunction)
	}
}

func (t *TargetPair) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(t.Old)
	formatter.WriteString(" TO ")
	formatter.WriteExpr(t.New)
}

func (t *TernaryOperation) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(t.Condition)
	formatter.WriteString(" ? ")
	formatter.WriteExpr(t.TrueExpr)
	formatter.WriteString(" : ")
	formatter.WriteExpr(t.FalseExpr)
}

func (t *TopClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("TOP ")
	formatter.WriteString(t.Number.Literal)
	if t.WithTies {
		formatter.WriteString(" WITH TIES")
	}
}

func (t *TruncateTable) FormatSQL(formatter *Formatter) {
	formatter.WriteString("TRUNCATE ")
	if t.IsTemporary {
		formatter.WriteString("TEMPORARY ")
	}
	formatter.WriteString("TABLE ")
	if t.IfExists {
		formatter.WriteString("IF EXISTS ")
	}
	formatter.WriteExpr(t.Name)
	if t.OnCluster != nil {
		formatter.WriteByte(whitespace)
		formatter.WriteExpr(t.OnCluster)
	}
}

func (s *TypeWithParams) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(s.Name)
	formatter.WriteByte('(')
	for i, size := range s.Params {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(size)
	}
	formatter.WriteByte(')')
}

func (t *TypedPlaceholder) FormatSQL(formatter *Formatter) {
	formatter.WriteString("{")
	formatter.WriteExpr(t.Name)
	formatter.WriteByte(':')
	formatter.WriteExpr(t.Type)
	formatter.WriteString("}")
}

func (u *UUID) FormatSQL(formatter *Formatter) {
	formatter.WriteString("UUID ")
	formatter.WriteExpr(u.Value)
}

func (n *UnaryExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteString(string(n.Kind))
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(n.Expr)
}

func (u *UpdateAssignment) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(u.Column)
	formatter.WriteString(" = ")
	formatter.WriteExpr(u.Expr)
}

func (u *UseStmt) FormatSQL(formatter *Formatter) {
	formatter.WriteString("USE ")
	formatter.WriteExpr(u.Database)
}

func (u *UsingClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("USING ")
	formatter.WriteExpr(u.Using)
}

func (w *WhenClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("WHEN ")
	formatter.WriteExpr(w.When)
	formatter.WriteString(" THEN ")
	formatter.WriteExpr(w.Then)
	if w.Else != nil {
		formatter.WriteString(" ELSE ")
		formatter.WriteExpr(w.Else)
	}
}

func (w *WhereClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("WHERE")
	if isLogicalBinaryOp(w.Expr) {
		formatter.Break()
		formatter.WriteExpr(w.Expr)
	} else {
		formatter.Indent()
		formatter.Break()
		formatter.WriteExpr(w.Expr)
		formatter.Dedent()
	}
}

func (w *WindowDefinition) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(w.Name)
	formatter.WriteString(" AS ")
	formatter.WriteExpr(w.Expr)
}

func (w *WindowClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("WINDOW ")
	for i, window := range w.Windows {
		window.FormatSQL(formatter)
		if i != len(w.Windows)-1 {
			formatter.WriteString(", ")
		}
	}
}

func (w *WindowExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteByte('(')
	hasPart := false
	if w.PartitionBy != nil {
		formatter.WriteExpr(w.PartitionBy)
		hasPart = true
	}
	if w.OrderBy != nil {
		if hasPart {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(w.OrderBy)
		hasPart = true
	}
	if w.Frame != nil {
		if hasPart {
			formatter.WriteByte(whitespace)
		}
		formatter.WriteExpr(w.Frame)
	}
	formatter.WriteByte(')')
}

func (f *WindowFrameClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString(f.Type)
	formatter.WriteByte(whitespace)
	formatter.WriteExpr(f.Extend)
}

func (f *WindowFrameCurrentRow) FormatSQL(formatter *Formatter) {
	formatter.WriteString("CURRENT ROW")
}

func (f *WindowFrameExtendExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(f.Expr)
	if f.Direction != "" {
		formatter.WriteByte(whitespace)
		formatter.WriteString(f.Direction)
	}
}

func (f *WindowFrameNumber) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(f.Number)
	formatter.WriteByte(whitespace)
	formatter.WriteString(f.Direction)
}

func (f *WindowFrameParam) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(f.Param)
	formatter.WriteByte(whitespace)
	formatter.WriteString(f.Direction)
}

func (f *WindowFrameUnbounded) FormatSQL(formatter *Formatter) {
	formatter.WriteString("UNBOUNDED ")
	formatter.WriteString(f.Direction)
}

func (w *WindowFunctionExpr) FormatSQL(formatter *Formatter) {
	formatter.WriteExpr(w.Function)
	formatter.WriteString(" OVER ")
	formatter.WriteExpr(w.OverExpr)
}

func (w *WithClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("WITH ")
	for i, cte := range w.CTEs {
		if i > 0 {
			formatter.WriteString(", ")
		}
		formatter.WriteExpr(cte)
	}
}

func (w *WithTimeoutClause) FormatSQL(formatter *Formatter) {
	formatter.WriteString("WITH TIMEOUT ")
	formatter.WriteExpr(w.Number)
}
