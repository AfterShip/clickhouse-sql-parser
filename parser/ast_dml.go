package parser

type UpdateAssignment struct {
	AssignmentPos Pos
	Column        *NestedIdentifier
	Expr          Expr
}

func (u *UpdateAssignment) Pos() Pos {
	return u.AssignmentPos
}

func (u *UpdateAssignment) End() Pos {
	return u.Expr.End()
}

func (u *UpdateAssignment) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	if err := u.Column.Accept(visitor); err != nil {
		return err
	}
	if err := u.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUpdateAssignment(u)
}

type DeleteClause struct {
	DeletePos Pos
	Table     *TableIdentifier
	OnCluster *ClusterClause
	WhereExpr Expr
}

func (d *DeleteClause) Pos() Pos {
	return d.DeletePos
}

func (d *DeleteClause) End() Pos {
	return d.WhereExpr.End()
}

func (d *DeleteClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Table.Accept(visitor); err != nil {
		return err
	}
	if d.OnCluster != nil {
		if err := d.OnCluster.Accept(visitor); err != nil {
			return err
		}
	}
	if d.WhereExpr != nil {
		if err := d.WhereExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDeleteFromExpr(d)
}

type ColumnNamesExpr struct {
	LeftParenPos  Pos
	RightParenPos Pos
	ColumnNames   []NestedIdentifier
}

func (c *ColumnNamesExpr) Pos() Pos {
	return c.LeftParenPos
}

func (c *ColumnNamesExpr) End() Pos {
	return c.RightParenPos
}

func (c *ColumnNamesExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for i := range c.ColumnNames {
		if err := c.ColumnNames[i].Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnNamesExpr(c)
}

type AssignmentValues struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Values        []Expr
}

func (v *AssignmentValues) Pos() Pos {
	return v.LeftParenPos
}

func (v *AssignmentValues) End() Pos {
	return v.RightParenPos
}

func (v *AssignmentValues) Accept(visitor ASTVisitor) error {
	visitor.Enter(v)
	defer visitor.Leave(v)
	for _, value := range v.Values {
		if err := value.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitValuesExpr(v)
}

type InsertStmt struct {
	InsertPos       Pos
	Format          *FormatClause
	HasTableKeyword bool
	Table           Expr
	ColumnNames     *ColumnNamesExpr
	Values          []*AssignmentValues
	SelectExpr      *SelectQuery
}

func (i *InsertStmt) Pos() Pos {
	return i.InsertPos
}

func (i *InsertStmt) End() Pos {
	if i.SelectExpr != nil {
		return i.SelectExpr.End()
	}
	if len(i.Values) > 0 {
		return i.Values[len(i.Values)-1].End()
	}
	// `INSERT INTO t FORMAT CSV` carries neither VALUES nor a SELECT — the
	// data arrives out of band after the statement
	if i.Format != nil {
		return i.Format.End()
	}
	if i.ColumnNames != nil {
		return i.ColumnNames.End()
	}
	return i.Table.End()
}

func (i *InsertStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.Table.Accept(visitor); err != nil {
		return err
	}
	if i.ColumnNames != nil {
		if err := i.ColumnNames.Accept(visitor); err != nil {
			return err
		}
	}
	if i.Format != nil {
		if err := i.Format.Accept(visitor); err != nil {
			return err
		}
	}
	for _, value := range i.Values {
		if err := value.Accept(visitor); err != nil {
			return err
		}
	}
	if i.SelectExpr != nil {
		if err := i.SelectExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitInsertExpr(i)
}
