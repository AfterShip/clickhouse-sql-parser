package parser

func (r *RowPolicyTarget) FormatSQL(f *Formatter) {
	f.WriteExpr(r.PolicyName)
	f.WriteString(" ON ")
	f.WriteExpr(r.Table)
}

func (r *RowPolicyCondition) FormatSQL(f *Formatter) {
	if r.None {
		f.WriteString("NONE")
		return
	}
	f.WriteExpr(r.Expr)
}

func (r *RowPolicyFilter) FormatSQL(f *Formatter) {
	f.WriteString("FOR SELECT")
	if r.Using != nil {
		f.WriteString(" USING ")
		f.WriteExpr(r.Using)
	}
	if r.WithCheck != nil {
		f.WriteString(" WITH CHECK ")
		f.WriteExpr(r.WithCheck)
	}
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

func (r *CreateRowPolicy) FormatSQL(f *Formatter) {
	f.WriteString("CREATE ROW POLICY ")
	if r.IfNotExists {
		f.WriteString("IF NOT EXISTS ")
	} else if r.OrReplace {
		f.WriteString("OR REPLACE ")
	}
	r.RowPolicyOptions.format(f, nil)
}

func (r *AlterRowPolicy) FormatSQL(f *Formatter) {
	f.WriteString("ALTER ROW POLICY ")
	if r.IfExists {
		f.WriteString("IF EXISTS ")
	}
	r.RowPolicyOptions.format(f, r.RenameTo)
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
