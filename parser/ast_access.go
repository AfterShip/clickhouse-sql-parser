package parser

type RoleName struct {
	Name      Expr
	Scope     *StringLiteral
	OnCluster *ClusterClause
}

func (r *RoleName) Pos() Pos {
	return r.Name.Pos()
}

func (r *RoleName) End() Pos {
	if r.Scope != nil {
		return r.Scope.End()
	}
	if r.OnCluster != nil {
		return r.OnCluster.End()
	}
	return r.Name.End()
}

func (r *RoleName) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.Name.Accept(visitor); err != nil {
		return err
	}
	if r.Scope != nil {
		if err := r.Scope.Accept(visitor); err != nil {
			return err
		}
	}
	if r.OnCluster != nil {
		if err := r.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRoleName(r)
}

type SettingPair struct {
	Name      *Ident
	Operation TokenKind
	Value     Expr
}

func (s *SettingPair) Pos() Pos {
	return s.Name.NamePos
}

func (s *SettingPair) End() Pos {
	return s.Value.End()
}

func (s *SettingPair) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	if s.Value != nil {
		if err := s.Value.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSettingPair(s)
}

type RoleSetting struct {
	SettingPairs []*SettingPair
	Modifier     *Ident
}

func (r *RoleSetting) Pos() Pos {
	if len(r.SettingPairs) > 0 {
		return r.SettingPairs[0].Pos()
	}
	return r.Modifier.NamePos
}

func (r *RoleSetting) End() Pos {
	if r.Modifier != nil {
		return r.Modifier.NameEnd
	}
	return r.SettingPairs[len(r.SettingPairs)-1].End()
}

func (r *RoleSetting) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	for _, settingPair := range r.SettingPairs {
		if err := settingPair.Accept(visitor); err != nil {
			return err
		}
	}
	if r.Modifier != nil {
		if err := r.Modifier.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRoleSetting(r)
}

type CreateRole struct {
	CreatePos         Pos
	StatementEnd      Pos
	IfNotExists       bool
	OrReplace         bool
	RoleNames         []*RoleName
	AccessStorageType *Ident
	Settings          []*RoleSetting
}

func (c *CreateRole) Pos() Pos {
	return c.CreatePos
}

func (c *CreateRole) End() Pos {
	return c.StatementEnd
}

func (c *CreateRole) Type() string {
	return "ROLE"
}

func (c *CreateRole) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, roleName := range c.RoleNames {
		if err := roleName.Accept(visitor); err != nil {
			return err
		}
	}
	if c.AccessStorageType != nil {
		if err := c.AccessStorageType.Accept(visitor); err != nil {
			return err
		}
	}
	for _, setting := range c.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateRole(c)
}

type AuthenticationClause struct {
	AuthPos       Pos
	AuthEnd       Pos
	NotIdentified bool
	AuthType      string // "no_password", "plaintext_password", "sha256_password", etc.
	AuthValue     *StringLiteral
	LdapServer    *StringLiteral
	KerberosRealm *StringLiteral
	IsKerberos    bool
}

func (a *AuthenticationClause) Pos() Pos {
	return a.AuthPos
}

func (a *AuthenticationClause) End() Pos {
	return a.AuthEnd
}

func (a *AuthenticationClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if a.AuthValue != nil {
		if err := a.AuthValue.Accept(visitor); err != nil {
			return err
		}
	}
	if a.LdapServer != nil {
		if err := a.LdapServer.Accept(visitor); err != nil {
			return err
		}
	}
	if a.KerberosRealm != nil {
		if err := a.KerberosRealm.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAuthenticationClause(a)
}

type HostClause struct {
	HostPos   Pos
	HostEnd   Pos
	HostType  string // "LOCAL", "NAME", "REGEXP", "IP", "LIKE", "ANY", "NONE"
	HostValue *StringLiteral
}

func (h *HostClause) Pos() Pos {
	return h.HostPos
}

func (h *HostClause) End() Pos {
	return h.HostEnd
}

func (h *HostClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(h)
	defer visitor.Leave(h)
	if h.HostValue != nil {
		if err := h.HostValue.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitHostClause(h)
}

type DefaultRoleClause struct {
	DefaultPos Pos
	DefaultEnd Pos
	Roles      []*RoleName
	None       bool
}

func (d *DefaultRoleClause) Pos() Pos {
	return d.DefaultPos
}

func (d *DefaultRoleClause) End() Pos {
	return d.DefaultEnd
}

func (d *DefaultRoleClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	for _, role := range d.Roles {
		if err := role.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDefaultRoleClause(d)
}

type GranteesClause struct {
	GranteesPos Pos
	GranteesEnd Pos
	Grantees    []*RoleName
	ExceptUsers []*RoleName
	Any         bool
	None        bool
}

func (g *GranteesClause) Pos() Pos {
	return g.GranteesPos
}

func (g *GranteesClause) End() Pos {
	return g.GranteesEnd
}

func (g *GranteesClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	for _, grantee := range g.Grantees {
		if err := grantee.Accept(visitor); err != nil {
			return err
		}
	}
	for _, except := range g.ExceptUsers {
		if err := except.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitGranteesClause(g)
}

type CreateUser struct {
	CreatePos       Pos
	StatementEnd    Pos
	IfNotExists     bool
	OrReplace       bool
	UserNames       []*RoleName
	Authentication  *AuthenticationClause
	ValidUntil      *StringLiteral
	Hosts           []*HostClause
	DefaultRole     *DefaultRoleClause
	DefaultDatabase *Ident
	DefaultDbNone   bool
	Grantees        *GranteesClause
	Settings        []*RoleSetting
}

func (c *CreateUser) Pos() Pos {
	return c.CreatePos
}

func (c *CreateUser) End() Pos {
	return c.StatementEnd
}

func (c *CreateUser) Type() string {
	return "USER"
}

func (c *CreateUser) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, userName := range c.UserNames {
		if err := userName.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Authentication != nil {
		if err := c.Authentication.Accept(visitor); err != nil {
			return err
		}
	}
	if c.ValidUntil != nil {
		if err := c.ValidUntil.Accept(visitor); err != nil {
			return err
		}
	}
	for _, host := range c.Hosts {
		if err := host.Accept(visitor); err != nil {
			return err
		}
	}
	if c.DefaultRole != nil {
		if err := c.DefaultRole.Accept(visitor); err != nil {
			return err
		}
	}
	if c.DefaultDatabase != nil {
		if err := c.DefaultDatabase.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Grantees != nil {
		if err := c.Grantees.Accept(visitor); err != nil {
			return err
		}
	}
	for _, setting := range c.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateUser(c)
}

type AlterRole struct {
	AlterPos        Pos
	StatementEnd    Pos
	IfExists        bool
	RoleRenamePairs []*RoleRenamePair
	Settings        []*RoleSetting
}

func (a *AlterRole) Pos() Pos {
	return a.AlterPos
}

func (a *AlterRole) End() Pos {
	return a.StatementEnd
}

func (a *AlterRole) Type() string {
	return "ROLE"
}

func (a *AlterRole) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	for _, roleRenamePair := range a.RoleRenamePairs {
		if err := roleRenamePair.Accept(visitor); err != nil {
			return err
		}
	}
	for _, setting := range a.Settings {
		if err := setting.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitAlterRole(a)
}

type RoleRenamePair struct {
	RoleName     *RoleName
	NewName      Expr
	StatementEnd Pos
}

func (r *RoleRenamePair) Pos() Pos {
	return r.RoleName.Pos()
}

func (r *RoleRenamePair) End() Pos {
	return r.StatementEnd
}

func (r *RoleRenamePair) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.RoleName.Accept(visitor); err != nil {
		return err
	}
	if r.NewName != nil {
		if err := r.NewName.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRoleRenamePair(r)
}

type DropUserOrRole struct {
	DropPos      Pos
	Target       string
	StatementEnd Pos
	Names        []*RoleName
	IfExists     bool
	Modifier     string
	From         *Ident
}

func (d *DropUserOrRole) Pos() Pos {
	return d.DropPos
}

func (d *DropUserOrRole) End() Pos {
	return d.StatementEnd
}

func (d *DropUserOrRole) Type() string {
	return d.Target
}

func (d *DropUserOrRole) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	for _, name := range d.Names {
		if err := name.Accept(visitor); err != nil {
			return err
		}
	}
	if d.From != nil {
		if err := d.From.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDropUserOrRole(d)
}

type PrivilegeClause struct {
	PrivilegePos Pos
	PrivilegeEnd Pos
	Keywords     []string
	Params       *ParamExprList
}

func (p *PrivilegeClause) Pos() Pos {
	return p.PrivilegePos
}

func (p *PrivilegeClause) End() Pos {
	return p.PrivilegeEnd
}

func (p *PrivilegeClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if p.Params != nil {
		if err := p.Params.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPrivilegeExpr(p)
}

type GrantPrivilegeStmt struct {
	GrantPos     Pos
	StatementEnd Pos
	OnCluster    *ClusterClause
	Privileges   []*PrivilegeClause
	On           *TableIdentifier
	To           []*Ident
	WithOptions  []string
}

func (g *GrantPrivilegeStmt) Pos() Pos {
	return g.GrantPos
}

func (g *GrantPrivilegeStmt) End() Pos {
	return g.StatementEnd
}

func (g *GrantPrivilegeStmt) Type() string {
	return "GRANT PRIVILEGE"
}

func (g *GrantPrivilegeStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	if g.OnCluster != nil {
		if err := g.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	for _, privilege := range g.Privileges {
		if err := privilege.Accept(visitor); err != nil {
			return err
		}
	}
	if err := g.On.Accept(visitor); err != nil {
		return err
	}
	for _, role := range g.To {
		if err := role.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitGrantPrivilegeExpr(g)
}
