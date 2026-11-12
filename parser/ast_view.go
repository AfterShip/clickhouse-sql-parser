package parser

type CreateMaterializedView struct {
	CreatePos    Pos // position of CREATE|ATTACH keyword
	StatementEnd Pos
	OrReplace    bool
	Name         *TableIdentifier
	IfNotExists  bool
	OnCluster    *ClusterClause
	Refresh      *RefreshExpr
	RandomizeFor *IntervalExpr
	DependsOn    []*TableIdentifier
	Settings     *SettingsClause
	HasAppend    bool
	Engine       *EngineExpr
	TableSchema  *TableSchemaClause // column list for ENGINE-based MVs (no TO clause)
	HasEmpty     bool
	Destination  *DestinationClause
	SubQuery     *SubQuery
	Populate     bool
	Comment      *StringLiteral
	Definer      *Ident
	SQLSecurity  string
}

func (c *CreateMaterializedView) Pos() Pos {
	return c.CreatePos
}

func (c *CreateMaterializedView) End() Pos {
	return c.StatementEnd
}

func (c *CreateMaterializedView) Type() string {
	return "MATERIALIZED_VIEW"
}

func (c *CreateMaterializedView) Accept(visitor ASTVisitor) error {
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
	if c.Refresh != nil {
		if err := c.Refresh.Accept(visitor); err != nil {
			return err
		}
	}
	if c.RandomizeFor != nil {
		if err := c.RandomizeFor.Accept(visitor); err != nil {
			return err
		}
	}
	if c.DependsOn != nil {
		for _, dep := range c.DependsOn {
			if err := dep.Accept(visitor); err != nil {
				return err
			}
		}
	}
	if c.Settings != nil {
		if err := c.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TableSchema != nil {
		if err := c.TableSchema.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Engine != nil {
		if err := c.Engine.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Destination != nil {
		// Destination.Accept visits its own TableSchema
		if err := c.Destination.Accept(visitor); err != nil {
			return err
		}
	}
	if c.SubQuery != nil {
		if err := c.SubQuery.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Definer != nil {
		if err := c.Definer.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Comment != nil {
		if err := c.Comment.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateMaterializedView(c)
}

type CreateView struct {
	CreatePos    Pos // position of CREATE|ATTACH keyword
	StatementEnd Pos
	OrReplace    bool
	Name         *TableIdentifier
	IfNotExists  bool
	UUID         *UUID
	OnCluster    *ClusterClause
	TableSchema  *TableSchemaClause
	Comment      *StringLiteral
	SubQuery     *SubQuery
}

func (c *CreateView) Pos() Pos {
	return c.CreatePos
}

func (c *CreateView) End() Pos {
	return c.StatementEnd
}

func (c *CreateView) Type() string {
	return "VIEW"
}

func (c *CreateView) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	if c.UUID != nil {
		if err := c.UUID.Accept(visitor); err != nil {
			return err
		}
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TableSchema != nil {
		if err := c.TableSchema.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Comment != nil {
		if err := c.Comment.Accept(visitor); err != nil {
			return err
		}
	}
	if c.SubQuery != nil {
		if err := c.SubQuery.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateView(c)
}

type RefreshExpr struct {
	RefreshPos Pos
	Frequency  string // EVERY|AFTER
	Interval   *IntervalExpr
	Offset     *IntervalExpr
}

func (r *RefreshExpr) Pos() Pos {
	return r.RefreshPos
}

func (r *RefreshExpr) End() Pos {
	if r.Offset != nil {
		return r.Offset.End()
	}
	return r.Interval.End()
}

func (r *RefreshExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if r.Interval != nil {
		if err := r.Interval.Accept(visitor); err != nil {
			return err
		}
	}
	if r.Offset != nil {
		if err := r.Offset.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRefreshExpr(r)
}

type CreateLiveView struct {
	CreatePos    Pos
	StatementEnd Pos
	Name         *TableIdentifier
	IfNotExists  bool
	UUID         *UUID
	OnCluster    *ClusterClause
	Destination  *DestinationClause
	TableSchema  *TableSchemaClause
	WithTimeout  *WithTimeoutClause
	SubQuery     *SubQuery
}

func (c *CreateLiveView) Type() string {
	return "LIVE_VIEW"
}

func (c *CreateLiveView) Pos() Pos {
	return c.CreatePos
}

func (c *CreateLiveView) End() Pos {
	return c.StatementEnd
}

func (c *CreateLiveView) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	if c.UUID != nil {
		if err := c.UUID.Accept(visitor); err != nil {
			return err
		}
	}
	if c.OnCluster != nil {
		if err := c.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Destination != nil {
		if err := c.Destination.Accept(visitor); err != nil {
			return err
		}
	}
	if c.TableSchema != nil {
		if err := c.TableSchema.Accept(visitor); err != nil {
			return err
		}
	}
	if c.WithTimeout != nil {
		if err := c.WithTimeout.Accept(visitor); err != nil {
			return err
		}
	}
	if c.SubQuery != nil {
		if err := c.SubQuery.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateLiveView(c)
}

type WithTimeoutClause struct {
	WithTimeoutPos Pos
	Expr           Expr
	Number         *NumberLiteral
}

func (w *WithTimeoutClause) Pos() Pos {
	return w.WithTimeoutPos
}

func (w *WithTimeoutClause) End() Pos {
	if w.Number != nil {
		return w.Number.End()
	}
	if w.Expr != nil {
		return w.Expr.End()
	}
	return w.WithTimeoutPos
}

func (w *WithTimeoutClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if w.Expr != nil {
		if err := w.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	if w.Number != nil {
		if err := w.Number.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitWithTimeoutExpr(w)
}
