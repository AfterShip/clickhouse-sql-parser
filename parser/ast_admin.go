package parser

type CreateFunction struct {
	CreatePos    Pos
	OrReplace    bool
	IfNotExists  bool
	FunctionName *Ident
	OnCluster    *ClusterClause
	Params       *ParamExprList
	Expr         Expr
}

func (c *CreateFunction) Type() string {
	return "FUNCTION"
}

func (c *CreateFunction) Pos() Pos {
	return c.CreatePos
}

func (c *CreateFunction) End() Pos {
	return c.Expr.End()
}

func (c *CreateFunction) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.FunctionName.Accept(visitor); err != nil {
		return err
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if err := c.Params.Accept(visitor); err != nil {
		return err
	}
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCreateFunction(c)
}

type CreateNamedCollection struct {
	CreatePos    Pos
	StatementEnd Pos
	Name         *Ident
	IfNotExists  bool
	OnCluster    *ClusterClause
	Params       []*NamedCollectionParam
}

func (c *CreateNamedCollection) Pos() Pos {
	return c.CreatePos
}

func (c *CreateNamedCollection) End() Pos {
	return c.StatementEnd
}

func (c *CreateNamedCollection) Type() string {
	return "NAMED COLLECTION"
}

func (c *CreateNamedCollection) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	for _, param := range c.Params {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateNamedCollection(c)
}

type NamedCollectionParam struct {
	ParamPos       Pos
	Name           *Ident
	Value          Expr
	Overridable    bool
	NotOverridable bool
}

func (n *NamedCollectionParam) Pos() Pos {
	return n.ParamPos
}

func (n *NamedCollectionParam) End() Pos {
	return n.Value.End()
}

func (n *NamedCollectionParam) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Name.Accept(visitor); err != nil {
		return err
	}
	if err := n.Value.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNamedCollectionParam(n)
}

type UseStmt struct {
	UsePos       Pos
	StatementEnd Pos
	Database     *Ident
}

func (u *UseStmt) Pos() Pos {
	return u.UsePos
}

func (u *UseStmt) End() Pos {
	return u.Database.End()
}

func (u *UseStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	if err := u.Database.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUseExpr(u)
}

type SetStmt struct {
	SetPos   Pos
	Settings *SettingsClause
}

func (s *SetStmt) Pos() Pos {
	return s.SetPos
}

func (s *SetStmt) End() Pos {
	return s.Settings.End()
}

func (s *SetStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Settings.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSetExpr(s)
}

type OptimizeStmt struct {
	OptimizePos  Pos
	StatementEnd Pos
	Table        *TableIdentifier
	OnCluster    *ClusterClause
	Partition    *PartitionClause
	HasFinal     bool
	Deduplicate  *DeduplicateClause
}

func (o *OptimizeStmt) Pos() Pos {
	return o.OptimizePos
}

func (o *OptimizeStmt) End() Pos {
	return o.StatementEnd
}

func (o *OptimizeStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.Table.Accept(visitor); err != nil {
		return err
	}
	if o.OnCluster != nil {
		if err := o.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Partition != nil {
		if err := o.Partition.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Deduplicate != nil {
		if err := o.Deduplicate.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitOptimizeExpr(o)
}

type DeduplicateClause struct {
	DeduplicatePos Pos
	By             *ColumnExprList
	Except         *ColumnExprList
}

func (d *DeduplicateClause) Pos() Pos {
	return d.DeduplicatePos
}

func (d *DeduplicateClause) End() Pos {
	if d.By != nil {
		return d.By.End()
	} else if d.Except != nil {
		return d.Except.End()
	}
	return d.DeduplicatePos + Pos(len(KeywordDeduplicate))
}

func (d *DeduplicateClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if d.By != nil {
		if err := d.By.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Except != nil {
		if err := d.Except.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDeduplicateExpr(d)
}

type SystemStmt struct {
	SystemPos Pos
	Expr      Expr
}

func (s *SystemStmt) Pos() Pos {
	return s.SystemPos
}

func (s *SystemStmt) End() Pos {
	return s.Expr.End()
}

func (s *SystemStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSystemExpr(s)
}

type SystemFlushExpr struct {
	FlushPos     Pos
	StatementEnd Pos
	Logs         bool
	Distributed  *TableIdentifier
}

func (s *SystemFlushExpr) Pos() Pos {
	return s.FlushPos
}

func (s *SystemFlushExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemFlushExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Distributed != nil {
		if err := s.Distributed.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSystemFlushExpr(s)
}

type SystemReloadExpr struct {
	ReloadPos    Pos
	StatementEnd Pos
	OnCluster    *ClusterClause
	Dictionary   *TableIdentifier
	Type         string
}

func (s *SystemReloadExpr) Pos() Pos {
	return s.ReloadPos
}

func (s *SystemReloadExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemReloadExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.OnCluster != nil {
		if err := s.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Dictionary != nil {
		if err := s.Dictionary.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSystemReloadExpr(s)
}

type SystemSyncExpr struct {
	SyncPos Pos
	Cluster *TableIdentifier
}

func (s *SystemSyncExpr) Pos() Pos {
	return s.SyncPos
}

func (s *SystemSyncExpr) End() Pos {
	return s.Cluster.End()
}

func (s *SystemSyncExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Cluster.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSystemSyncExpr(s)
}

type SystemCtrlExpr struct {
	CtrlPos      Pos
	StatementEnd Pos
	Command      string // START, STOP
	Type         string // REPLICATED, DISTRIBUTED
	Cluster      *TableIdentifier
}

func (s *SystemCtrlExpr) Pos() Pos {
	return s.CtrlPos
}

func (s *SystemCtrlExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemCtrlExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Cluster != nil {
		if err := s.Cluster.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSystemCtrlExpr(s)
}

type SystemDropExpr struct {
	DropPos      Pos
	StatementEnd Pos
	Type         string
}

func (s *SystemDropExpr) Pos() Pos {
	return s.DropPos
}

func (s *SystemDropExpr) End() Pos {
	return s.StatementEnd
}

func (s *SystemDropExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	return visitor.VisitSystemDropExpr(s)
}

type CheckStmt struct {
	CheckPos  Pos
	Table     *TableIdentifier
	Partition *PartitionClause
}

func (c *CheckStmt) Pos() Pos {
	return c.CheckPos
}

func (c *CheckStmt) End() Pos {
	return c.Partition.End()
}

func (c *CheckStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Table.Accept(visitor); err != nil {
		return err
	}
	if c.Partition != nil {
		if err := c.Partition.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCheckExpr(c)
}

type ExplainStmt struct {
	ExplainPos Pos
	Type       string
	Statement  Expr
}

func (e *ExplainStmt) Pos() Pos {
	return e.ExplainPos
}

func (e *ExplainStmt) End() Pos {
	return e.Statement.End()
}

func (e *ExplainStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	if err := e.Statement.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitExplainExpr(e)
}

type ShowStmt struct {
	ShowPos      Pos
	StatementEnd Pos
	ShowType     string           // e.g., "CREATE TABLE", "DATABASES", "TABLES"
	Target       *TableIdentifier // for SHOW CREATE TABLE table_name

	// Optional clauses for SHOW DATABASES
	NotLike     bool           // true if NOT LIKE/ILIKE
	LikeType    string         // "LIKE" or "ILIKE", empty if not used
	LikePattern Expr           // pattern expression for LIKE/ILIKE
	Limit       Expr           // limit expression
	OutFile     *StringLiteral // filename for INTO OUTFILE
	Format      *StringLiteral // format specification
}

func (s *ShowStmt) Pos() Pos {
	return s.ShowPos
}

func (s *ShowStmt) End() Pos {
	// Find the rightmost element to determine the end position
	if s.Format != nil {
		return s.Format.End()
	}
	if s.OutFile != nil {
		return s.OutFile.End()
	}
	if s.Limit != nil {
		return s.Limit.End()
	}
	if s.LikePattern != nil {
		return s.LikePattern.End()
	}
	if s.Target != nil {
		return s.Target.End()
	}
	return s.StatementEnd
}

func (s *ShowStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Target != nil {
		if err := s.Target.Accept(visitor); err != nil {
			return err
		}
	}
	if s.LikePattern != nil {
		if err := s.LikePattern.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Limit != nil {
		if err := s.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if s.OutFile != nil {
		if err := s.OutFile.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Format != nil {
		if err := s.Format.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitShowExpr(s)
}

type DescribeStmt struct {
	DescribePos  Pos
	StatementEnd Pos
	DescribeType string // e.g., "TABLE", empty if not used
	Target       *TableIdentifier
}

func (d *DescribeStmt) Pos() Pos {
	return d.DescribePos
}

func (d *DescribeStmt) End() Pos {
	return d.Target.End()
}

func (d *DescribeStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Target.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitDescribeExpr(d)
}
