package parser

type ColumnType interface {
	Expr
	Type() string
}

type ScalarType struct {
	Name *Ident
}

func (s *ScalarType) Pos() Pos {
	return s.Name.NamePos
}

func (s *ScalarType) End() Pos {
	return s.Name.NameEnd
}

func (s *ScalarType) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitScalarType(s)
}

func (s *ScalarType) Type() string {
	return s.Name.Name
}

type JSONPath struct {
	Idents []*Ident
}

func (j *JSONPath) accept(visitor ASTVisitor) error {
	if j == nil {
		return nil
	}
	for _, ident := range j.Idents {
		if ident == nil {
			continue
		}
		if err := ident.Accept(visitor); err != nil {
			return err
		}
	}
	return nil
}

type JSONTypeHint struct {
	Path *JSONPath
	Type ColumnType
}

func (j *JSONTypeHint) accept(visitor ASTVisitor) error {
	if j == nil {
		return nil
	}
	if err := j.Path.accept(visitor); err != nil {
		return err
	}
	if j.Type != nil {
		if err := j.Type.Accept(visitor); err != nil {
			return err
		}
	}
	return nil
}

type JSONOption struct {
	SkipPath        *JSONPath
	SkipRegex       *StringLiteral
	MaxDynamicPaths *NumberLiteral
	MaxDynamicTypes *NumberLiteral
	// Type hint for specific JSON subcolumn path, e.g., "message String" or "a.b UInt64"
	Column *JSONTypeHint
}

func (j *JSONOption) accept(visitor ASTVisitor) error {
	if j == nil {
		return nil
	}
	if err := j.SkipPath.accept(visitor); err != nil {
		return err
	}
	if j.SkipRegex != nil {
		if err := j.SkipRegex.Accept(visitor); err != nil {
			return err
		}
	}
	if j.MaxDynamicPaths != nil {
		if err := j.MaxDynamicPaths.Accept(visitor); err != nil {
			return err
		}
	}
	if j.MaxDynamicTypes != nil {
		if err := j.MaxDynamicTypes.Accept(visitor); err != nil {
			return err
		}
	}
	if err := j.Column.accept(visitor); err != nil {
		return err
	}
	return nil
}

type JSONOptions struct {
	LParen Pos
	RParen Pos
	Items  []*JSONOption
}

func (j *JSONOptions) Pos() Pos {
	return j.LParen
}

func (j *JSONOptions) End() Pos {
	return j.RParen
}

func (j *JSONOptions) accept(visitor ASTVisitor) error {
	if j == nil {
		return nil
	}
	for _, item := range j.Items {
		if err := item.accept(visitor); err != nil {
			return err
		}
	}
	return nil
}

type JSONType struct {
	Name    *Ident
	Options *JSONOptions
}

func (j *JSONType) Pos() Pos {
	return j.Name.NamePos
}

func (j *JSONType) End() Pos {
	if j.Options != nil {
		return j.Options.RParen
	}
	return j.Name.NameEnd
}

func (j *JSONType) Type() string {
	return j.Name.Name
}

func (j *JSONType) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if err := j.Name.Accept(visitor); err != nil {
		return err
	}
	if err := j.Options.accept(visitor); err != nil {
		return err
	}
	return visitor.VisitJSONType(j)
}

type PropertyType struct {
	Name *Ident
}

func (c *PropertyType) Pos() Pos {
	return c.Name.NamePos
}

func (c *PropertyType) End() Pos {
	return c.Name.NameEnd
}

func (c *PropertyType) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPropertyType(c)
}

func (c *PropertyType) Type() string {
	return c.Name.Name
}

type TypeWithParams struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Name          *Ident
	Params        []Literal
}

func (s *TypeWithParams) Pos() Pos {
	return s.Name.NamePos
}

func (s *TypeWithParams) End() Pos {
	return s.RightParenPos
}

func (s *TypeWithParams) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	for _, param := range s.Params {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTypeWithParams(s)
}

func (s *TypeWithParams) Type() string {
	return s.Name.Name
}

type TypeWithNamedParams struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Name          *Ident
	Params        []*NamedParameterExpr
}

func (s *TypeWithNamedParams) Pos() Pos {
	return s.Name.NamePos
}

func (s *TypeWithNamedParams) End() Pos {
	return s.RightParenPos
}

func (s *TypeWithNamedParams) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	for _, param := range s.Params {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTypeWithNamedParams(s)
}

func (s *TypeWithNamedParams) Type() string {
	return s.Name.Name
}

type ComplexType struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Name          *Ident
	Params        []ColumnType
}

func (c *ComplexType) Pos() Pos {
	return c.Name.NamePos
}

func (c *ComplexType) End() Pos {
	return c.RightParenPos
}

func (c *ComplexType) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	for _, param := range c.Params {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitComplexType(c)
}

func (c *ComplexType) Type() string {
	return c.Name.Name
}

type NestedType struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Name          *Ident
	Columns       []Expr
}

func (n *NestedType) Pos() Pos {
	return n.Name.NamePos
}

func (n *NestedType) End() Pos {
	return n.RightParenPos
}

func (n *NestedType) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Name.Accept(visitor); err != nil {
		return err
	}
	for _, column := range n.Columns {
		if err := column.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitNestedType(n)
}

func (n *NestedType) Type() string {
	return n.Name.Name
}

type EnumValue struct {
	Name  *StringLiteral
	Value *NumberLiteral
}

func (e *EnumValue) Pos() Pos {
	return e.Name.Pos()
}

func (e *EnumValue) End() Pos {
	return e.Value.End()
}

func (e *EnumValue) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	if err := e.Name.Accept(visitor); err != nil {
		return err
	}
	if err := e.Value.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitEnumValue(e)
}

type EnumType struct {
	Name    *Ident
	ListPos Pos
	ListEnd Pos
	Values  []EnumValue
}

func (e *EnumType) Pos() Pos {
	return e.ListPos
}

func (e *EnumType) End() Pos {
	return e.ListEnd
}

func (e *EnumType) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	if err := e.Name.Accept(visitor); err != nil {
		return err
	}
	for i := range e.Values {
		if err := e.Values[i].Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitEnumType(e)
}

func (e *EnumType) Type() string {
	return e.Name.Name
}

type ColumnTypeExpr struct {
	Name *Ident
}

func (c *ColumnTypeExpr) Pos() Pos {
	return c.Name.NamePos
}

func (c *ColumnTypeExpr) End() Pos {
	return c.Name.NameEnd
}

func (c *ColumnTypeExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Name.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitColumnTypeExpr(c)
}
