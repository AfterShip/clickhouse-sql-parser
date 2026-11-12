package parser

type CreateDictionary struct {
	CreatePos    Pos
	StatementEnd Pos
	OrReplace    bool
	Name         *TableIdentifier
	IfNotExists  bool
	UUID         *UUID
	OnCluster    *ClusterClause
	Schema       *DictionarySchemaClause
	Engine       *DictionaryEngineClause
	Comment      *StringLiteral
}

func (c *CreateDictionary) Type() string {
	return "DICTIONARY"
}

func (c *CreateDictionary) Pos() Pos {
	return c.CreatePos
}

func (c *CreateDictionary) End() Pos {
	return c.StatementEnd
}

func (c *CreateDictionary) Accept(visitor ASTVisitor) error {
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
	if c.Schema != nil {
		if err := c.Schema.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Engine != nil {
		if err := c.Engine.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Comment != nil {
		if err := c.Comment.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCreateDictionary(c)
}

type DictionarySchemaClause struct {
	SchemaPos  Pos
	Attributes []*DictionaryAttribute
	RParenPos  Pos
}

func (d *DictionarySchemaClause) Pos() Pos {
	return d.SchemaPos
}

func (d *DictionarySchemaClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionarySchemaClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	for _, attr := range d.Attributes {
		if err := attr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionarySchemaClause(d)
}

type DictionaryAttribute struct {
	NamePos      Pos
	AttrEnd      Pos // end of the last property, or of the type when there is none
	Name         *Ident
	Type         ColumnType
	Default      Literal
	Expression   Expr
	Hierarchical bool
	Injective    bool
	IsObjectId   bool
}

func (d *DictionaryAttribute) Pos() Pos {
	return d.NamePos
}

func (d *DictionaryAttribute) End() Pos {
	return d.AttrEnd
}

func (d *DictionaryAttribute) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Name.Accept(visitor); err != nil {
		return err
	}
	if err := d.Type.Accept(visitor); err != nil {
		return err
	}
	if d.Default != nil {
		if err := d.Default.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Expression != nil {
		if err := d.Expression.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryAttribute(d)
}

type DictionaryEngineClause struct {
	EnginePos  Pos
	PrimaryKey *DictionaryPrimaryKeyClause
	Source     *DictionarySourceClause
	Lifetime   *DictionaryLifetimeClause
	Layout     *DictionaryLayoutClause
	Range      *DictionaryRangeClause
	Settings   *SettingsClause
}

func (d *DictionaryEngineClause) Pos() Pos {
	return d.EnginePos
}

func (d *DictionaryEngineClause) End() Pos {
	if d.Settings != nil {
		return d.Settings.End()
	}
	if d.Range != nil {
		return d.Range.End()
	}
	if d.Layout != nil {
		return d.Layout.End()
	}
	if d.Lifetime != nil {
		return d.Lifetime.End()
	}
	if d.Source != nil {
		return d.Source.End()
	}
	if d.PrimaryKey != nil {
		return d.PrimaryKey.End()
	}
	return d.EnginePos
}

func (d *DictionaryEngineClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if d.PrimaryKey != nil {
		if err := d.PrimaryKey.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Source != nil {
		if err := d.Source.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Lifetime != nil {
		if err := d.Lifetime.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Layout != nil {
		if err := d.Layout.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Range != nil {
		if err := d.Range.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Settings != nil {
		if err := d.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryEngineClause(d)
}

type DictionaryPrimaryKeyClause struct {
	PrimaryKeyPos Pos
	Keys          *ColumnExprList
	RParenPos     Pos
}

func (d *DictionaryPrimaryKeyClause) Pos() Pos {
	return d.PrimaryKeyPos
}

func (d *DictionaryPrimaryKeyClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryPrimaryKeyClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Keys.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitDictionaryPrimaryKeyClause(d)
}

type DictionarySourceClause struct {
	SourcePos Pos
	Source    *Ident
	Args      []*DictionaryArgExpr
	RParenPos Pos
}

func (d *DictionarySourceClause) Pos() Pos {
	return d.SourcePos
}

func (d *DictionarySourceClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionarySourceClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Source.Accept(visitor); err != nil {
		return err
	}
	for _, arg := range d.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionarySourceClause(d)
}

type DictionaryArgExpr struct {
	ArgPos Pos
	Name   *Ident
	Value  Expr // can be Ident with optional parentheses or literal, nil for a nested argument list
	// Args holds a nested argument list such as CREDENTIALS(USER 'u' PASSWORD 'p').
	Args []*DictionaryArgExpr
	// LParenPos tells a nested list apart from a value: it is only set when the
	// argument is written with parentheses, which may hold no argument at all.
	LParenPos Pos
	RParenPos Pos
}

func (d *DictionaryArgExpr) Pos() Pos {
	return d.ArgPos
}

func (d *DictionaryArgExpr) End() Pos {
	if d.Value != nil {
		return d.Value.End()
	}
	return d.RParenPos + 1
}

func (d *DictionaryArgExpr) hasArgList() bool {
	return d.LParenPos != 0
}

func (d *DictionaryArgExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Name.Accept(visitor); err != nil {
		return err
	}
	if d.Value != nil {
		if err := d.Value.Accept(visitor); err != nil {
			return err
		}
	}
	for _, arg := range d.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryArgExpr(d)
}

type DictionaryLifetimeClause struct {
	LifetimePos Pos
	Min         *NumberLiteral
	Max         *NumberLiteral
	Value       *NumberLiteral // for simple LIFETIME(value) form
	RParenPos   Pos
}

func (d *DictionaryLifetimeClause) Pos() Pos {
	return d.LifetimePos
}

func (d *DictionaryLifetimeClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryLifetimeClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if d.Value != nil {
		if err := d.Value.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Min != nil {
		if err := d.Min.Accept(visitor); err != nil {
			return err
		}
	}
	if d.Max != nil {
		if err := d.Max.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryLifetimeClause(d)
}

type DictionaryLayoutClause struct {
	LayoutPos Pos
	Layout    *Ident
	Args      []*DictionaryArgExpr
	RParenPos Pos
}

func (d *DictionaryLayoutClause) Pos() Pos {
	return d.LayoutPos
}

func (d *DictionaryLayoutClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryLayoutClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Layout.Accept(visitor); err != nil {
		return err
	}
	for _, arg := range d.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDictionaryLayoutClause(d)
}

type DictionaryRangeClause struct {
	RangePos  Pos
	Min       *Ident
	Max       *Ident
	RParenPos Pos
}

func (d *DictionaryRangeClause) Pos() Pos {
	return d.RangePos
}

func (d *DictionaryRangeClause) End() Pos {
	return d.RParenPos + 1
}

func (d *DictionaryRangeClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(d)
	defer visitor.Leave(d)
	if err := d.Min.Accept(visitor); err != nil {
		return err
	}
	if err := d.Max.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitDictionaryRangeClause(d)
}
