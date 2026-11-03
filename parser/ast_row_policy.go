package parser

// RowPolicyMode controls how a policy combines with other policies.
// The empty value means omitted, not explicitly PERMISSIVE.
type RowPolicyMode string

const rowPolicyType = "ROW POLICY"

const (
	RowPolicyPermissive  RowPolicyMode = "PERMISSIVE"
	RowPolicyRestrictive RowPolicyMode = "RESTRICTIVE"
)

// RowPolicyTarget binds a policy name to exactly one table or database wildcard.
// An unquoted Table.Table named * is a wildcard; a quoted * is a literal name.
type RowPolicyTarget struct {
	PolicyName *Ident
	Table      *TableIdentifier
}

func (r *RowPolicyTarget) Pos() Pos { return r.PolicyName.Pos() }
func (r *RowPolicyTarget) End() Pos { return rowPolicyTableEnd(r.Table) }
func (r *RowPolicyTarget) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.PolicyName.Accept(visitor); err != nil {
		return err
	}
	if err := r.Table.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitRowPolicyTarget(r)
}

func (r *RowPolicyTarget) FormatSQL(f *Formatter) {
	f.WriteExpr(r.PolicyName)
	f.WriteString(" ON ")
	f.WriteExpr(r.Table)
}

// RowPolicyCondition distinguishes an explicit NONE from an expression.
// A nil condition on a filter means that clause was omitted.
type RowPolicyCondition struct {
	ConditionPos Pos
	ConditionEnd Pos
	None         bool
	Expr         Expr
}

func (r *RowPolicyCondition) Pos() Pos { return r.ConditionPos }
func (r *RowPolicyCondition) End() Pos { return r.ConditionEnd }
func (r *RowPolicyCondition) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if r.Expr != nil {
		if err := r.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRowPolicyCondition(r)
}

func (r *RowPolicyCondition) FormatSQL(f *Formatter) {
	if r.None {
		f.WriteString("NONE")
		return
	}
	f.WriteExpr(r.Expr)
}

// RowPolicyFilter holds a row filter and its optional FOR clause.
// Empty Commands means FOR was omitted; explicit SELECT and ALL are preserved.
type RowPolicyFilter struct {
	FilterPos Pos
	FilterEnd Pos
	Commands  []string
	Using     *RowPolicyCondition
	WithCheck *RowPolicyCondition
}

func (r *RowPolicyFilter) Pos() Pos { return r.FilterPos }
func (r *RowPolicyFilter) End() Pos { return r.FilterEnd }
func (r *RowPolicyFilter) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if r.Using != nil {
		if err := r.Using.Accept(visitor); err != nil {
			return err
		}
	}
	if r.WithCheck != nil {
		if err := r.WithCheck.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRowPolicyFilter(r)
}

func (r *RowPolicyFilter) FormatSQL(f *Formatter) {
	hasFor := len(r.Commands) > 0
	if hasFor {
		f.WriteString("FOR ")
		for i, command := range r.Commands {
			if i > 0 {
				f.WriteString(", ")
			}
			f.WriteString(command)
		}
	}
	if r.Using != nil {
		if hasFor {
			f.WriteByte(whitespace)
		}
		f.WriteString("USING ")
		f.WriteExpr(r.Using)
	}
	if r.WithCheck != nil {
		if hasFor || r.Using != nil {
			f.WriteByte(whitespace)
		}
		f.WriteString("WITH CHECK ")
		f.WriteExpr(r.WithCheck)
	}
}

// RowPolicyPrincipal is a named user/role or the dynamic CURRENT_USER tag.
type RowPolicyPrincipal struct {
	PrincipalPos Pos
	PrincipalEnd Pos
	Name         *Ident
	Host         *Ident
	CurrentUser  bool
}

func (r *RowPolicyPrincipal) Pos() Pos { return r.PrincipalPos }
func (r *RowPolicyPrincipal) End() Pos { return r.PrincipalEnd }
func (r *RowPolicyPrincipal) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if r.Name != nil {
		if err := r.Name.Accept(visitor); err != nil {
			return err
		}
	}
	if r.Host != nil {
		if err := r.Host.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRowPolicyPrincipal(r)
}

func (r *RowPolicyPrincipal) FormatSQL(f *Formatter) {
	if r.CurrentUser {
		f.WriteString("CURRENT_USER")
		return
	}
	f.WriteExpr(r.Name)
	if r.Host != nil {
		f.WriteByte('@')
		f.WriteExpr(r.Host)
	}
}

// RowPolicySubjects represents the TO clause, including inclusions/exclusions.
// A non-nil empty set is TO NONE; nil on an ALTER means keep the existing set.
type RowPolicySubjects struct {
	ToPos      Pos
	ToEnd      Pos
	All        bool
	Principals []*RowPolicyPrincipal
	Except     []*RowPolicyPrincipal
}

func (r *RowPolicySubjects) Pos() Pos { return r.ToPos }
func (r *RowPolicySubjects) End() Pos { return r.ToEnd }
func (r *RowPolicySubjects) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	for _, principal := range r.Principals {
		if err := principal.Accept(visitor); err != nil {
			return err
		}
	}
	for _, principal := range r.Except {
		if err := principal.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRowPolicySubjects(r)
}

func (r *RowPolicySubjects) FormatSQL(f *Formatter) {
	f.WriteString("TO ")
	if r.All {
		f.WriteString("ALL")
	} else if len(r.Principals) == 0 {
		f.WriteString("NONE")
	} else {
		formatRowPolicyPrincipals(f, r.Principals)
	}
	if len(r.Except) > 0 {
		f.WriteString(" EXCEPT ")
		formatRowPolicyPrincipals(f, r.Except)
	}
}

func formatRowPolicyPrincipals(f *Formatter, principals []*RowPolicyPrincipal) {
	for i, principal := range principals {
		if i > 0 {
			f.WriteString(", ")
		}
		f.WriteExpr(principal)
	}
}

func formatRowPolicyTargets(f *Formatter, targets []*RowPolicyTarget) {
	for i, target := range targets {
		if i > 0 {
			f.WriteString(", ")
		}
		f.WriteExpr(target)
	}
}

// RowPolicyOptions is the shared desired definition for CREATE and ALTER.
// Omitted fields stay omitted so ALTER does not reset unrelated attributes.
type RowPolicyOptions struct {
	Targets           []*RowPolicyTarget
	OnCluster         *ClusterClause
	AccessStorageType *Ident
	Mode              RowPolicyMode
	Filters           []*RowPolicyFilter
	To                *RowPolicySubjects
}

func (r *RowPolicyOptions) format(f *Formatter, rename *Ident) {
	formatRowPolicyTargets(f, r.Targets)
	if r.AccessStorageType != nil {
		f.WriteString(" IN ")
		f.WriteExpr(r.AccessStorageType)
	}
	if r.OnCluster != nil {
		f.WriteByte(whitespace)
		f.WriteExpr(r.OnCluster)
	}
	if rename != nil {
		f.WriteString(" RENAME TO ")
		f.WriteExpr(rename)
	}
	if r.Mode != "" {
		f.WriteString(" AS ")
		f.WriteString(string(r.Mode))
	}
	for i, filter := range r.Filters {
		if i > 0 {
			f.WriteByte(',')
		}
		f.WriteByte(whitespace)
		f.WriteExpr(filter)
	}
	if r.To != nil {
		f.WriteByte(whitespace)
		f.WriteExpr(r.To)
	}
}

func (r *RowPolicyOptions) accept(visitor ASTVisitor, rename *Ident) error {
	for _, target := range r.Targets {
		if err := target.Accept(visitor); err != nil {
			return err
		}
	}
	if r.AccessStorageType != nil {
		if err := r.AccessStorageType.Accept(visitor); err != nil {
			return err
		}
	}
	if r.OnCluster != nil {
		if err := r.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if rename != nil {
		if err := rename.Accept(visitor); err != nil {
			return err
		}
	}
	for _, filter := range r.Filters {
		if err := filter.Accept(visitor); err != nil {
			return err
		}
	}
	if r.To != nil {
		return r.To.Accept(visitor)
	}
	return nil
}

// CreateRowPolicy creates one or more explicitly bound policies.
type CreateRowPolicy struct {
	CreatePos    Pos
	StatementEnd Pos
	IfNotExists  bool
	OrReplace    bool
	RowPolicyOptions
}

func (r *CreateRowPolicy) Pos() Pos     { return r.CreatePos }
func (r *CreateRowPolicy) End() Pos     { return r.StatementEnd }
func (r *CreateRowPolicy) Type() string { return rowPolicyType }
func (r *CreateRowPolicy) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.RowPolicyOptions.accept(visitor, nil); err != nil {
		return err
	}
	return visitor.VisitCreateRowPolicy(r)
}

func (r *CreateRowPolicy) FormatSQL(f *Formatter) {
	f.WriteString("CREATE ROW POLICY ")
	if r.IfNotExists {
		f.WriteString("IF NOT EXISTS ")
	} else if r.OrReplace {
		f.WriteString("OR REPLACE ")
	}
	r.RowPolicyOptions.format(f, nil)
}

// AlterRowPolicy changes only the specified policy attributes.
type AlterRowPolicy struct {
	AlterPos     Pos
	StatementEnd Pos
	IfExists     bool
	RenameTo     *Ident
	RowPolicyOptions
}

func (r *AlterRowPolicy) Pos() Pos     { return r.AlterPos }
func (r *AlterRowPolicy) End() Pos     { return r.StatementEnd }
func (r *AlterRowPolicy) Type() string { return rowPolicyType }
func (r *AlterRowPolicy) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.RowPolicyOptions.accept(visitor, r.RenameTo); err != nil {
		return err
	}
	return visitor.VisitAlterRowPolicy(r)
}

func (r *AlterRowPolicy) FormatSQL(f *Formatter) {
	f.WriteString("ALTER ROW POLICY ")
	if r.IfExists {
		f.WriteString("IF EXISTS ")
	}
	r.RowPolicyOptions.format(f, r.RenameTo)
}

// DropRowPolicy drops exactly the listed policy/table identities.
type DropRowPolicy struct {
	DropPos           Pos
	StatementEnd      Pos
	IfExists          bool
	Targets           []*RowPolicyTarget
	OnCluster         *ClusterClause
	AccessStorageType *Ident
}

func (r *DropRowPolicy) Pos() Pos     { return r.DropPos }
func (r *DropRowPolicy) End() Pos     { return r.StatementEnd }
func (r *DropRowPolicy) Type() string { return rowPolicyType }
func (r *DropRowPolicy) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	for _, target := range r.Targets {
		if err := target.Accept(visitor); err != nil {
			return err
		}
	}
	if r.AccessStorageType != nil {
		if err := r.AccessStorageType.Accept(visitor); err != nil {
			return err
		}
	}
	if r.OnCluster != nil {
		if err := r.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDropRowPolicy(r)
}

func (r *DropRowPolicy) FormatSQL(f *Formatter) {
	f.WriteString("DROP ROW POLICY ")
	if r.IfExists {
		f.WriteString("IF EXISTS ")
	}
	formatRowPolicyTargets(f, r.Targets)
	if r.AccessStorageType != nil {
		f.WriteString(" FROM ")
		f.WriteExpr(r.AccessStorageType)
	}
	if r.OnCluster != nil {
		f.WriteByte(whitespace)
		f.WriteExpr(r.OnCluster)
	}
}

// ShowRowPolicies lists policies or their CREATE definitions. A nil selector
// lists all policies; Name, Table and Targets are mutually exclusive selectors.
type ShowRowPolicies struct {
	ShowPos      Pos
	StatementEnd Pos
	Create       bool
	Name         *Ident
	Table        *TableIdentifier
	Targets      []*RowPolicyTarget
}

func (r *ShowRowPolicies) Pos() Pos { return r.ShowPos }
func (r *ShowRowPolicies) End() Pos { return r.StatementEnd }
func (r *ShowRowPolicies) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if r.Name != nil {
		if err := r.Name.Accept(visitor); err != nil {
			return err
		}
	}
	if r.Table != nil {
		if err := r.Table.Accept(visitor); err != nil {
			return err
		}
	}
	for _, target := range r.Targets {
		if err := target.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitShowRowPolicies(r)
}

func (r *ShowRowPolicies) FormatSQL(f *Formatter) {
	f.WriteString("SHOW ")
	if r.Create {
		f.WriteString("CREATE ")
	}
	if len(r.Targets) > 0 {
		f.WriteString("ROW POLICY ")
		formatRowPolicyTargets(f, r.Targets)
		return
	}
	f.WriteString("ROW POLICIES")
	if r.Name != nil {
		f.WriteByte(whitespace)
		f.WriteExpr(r.Name)
	}
	if r.Table != nil {
		f.WriteString(" ON ")
		f.WriteExpr(r.Table)
	}
}
