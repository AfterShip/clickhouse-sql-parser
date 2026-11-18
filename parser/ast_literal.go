package parser

type Ident struct {
	Name      string
	QuoteType int
	NamePos   Pos
	NameEnd   Pos
}

func (i *Ident) Pos() Pos {
	return i.NamePos
}

func (i *Ident) End() Pos {
	return i.NameEnd
}

func (i *Ident) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	return visitor.VisitIdent(i)
}

type UUID struct {
	Value *StringLiteral
}

func (u *UUID) Pos() Pos {
	return u.Value.LiteralPos
}

func (u *UUID) End() Pos {
	return u.Value.LiteralEnd
}

func (u *UUID) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	if u.Value != nil {
		if err := u.Value.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitUUID(u)
}

type NullLiteral struct {
	NullPos Pos
}

func (n *NullLiteral) Pos() Pos {
	return n.NullPos
}

func (n *NullLiteral) End() Pos {
	return n.NullPos + 4
}

func (n *NullLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	return visitor.VisitNullLiteral(n)
}

type NotNullLiteral struct {
	NotPos      Pos
	NullLiteral *NullLiteral
}

func (n *NotNullLiteral) Pos() Pos {
	return n.NotPos
}

func (n *NotNullLiteral) End() Pos {
	return n.NullLiteral.End()
}

func (n *NotNullLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.NullLiteral.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNotNullLiteral(n)
}

type NestedIdentifier struct {
	Ident    *Ident
	DotIdent *Ident
}

func (n *NestedIdentifier) Pos() Pos {
	return n.Ident.Pos()
}

func (n *NestedIdentifier) End() Pos {
	if n.DotIdent != nil {
		return n.DotIdent.End()
	}
	return n.Ident.End()
}

func (n *NestedIdentifier) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Ident.Accept(visitor); err != nil {
		return err
	}
	if n.DotIdent != nil {
		if err := n.DotIdent.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitNestedIdentifier(n)
}

type Path struct {
	Fields []*Ident
}

func (p *Path) Pos() Pos {
	if len(p.Fields) > 0 {
		return p.Fields[0].Pos()
	}
	return 0
}

func (p *Path) End() Pos {
	if len(p.Fields) > 0 {
		return p.Fields[len(p.Fields)-1].End()
	}
	return 0
}

func (p *Path) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	for _, ident := range p.Fields {
		if err := ident.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPath(p)
}

type TableIdentifier struct {
	Database *Ident
	Table    *Ident
}

func (t *TableIdentifier) Pos() Pos {
	if t.Database != nil {
		return t.Database.Pos()
	}
	return t.Table.Pos()
}

func (t *TableIdentifier) End() Pos {
	return t.Table.End()
}

func (t *TableIdentifier) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if t.Database != nil {
		if err := t.Database.Accept(visitor); err != nil {
			return err
		}
	}
	if err := t.Table.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTableIdentifier(t)
}

type KeyValue struct {
	Key   StringLiteral
	Value Expr
}

type MapLiteral struct {
	LBracePos Pos
	RBracePos Pos
	KeyValues []KeyValue
}

func (m *MapLiteral) Pos() Pos {
	return m.LBracePos
}

func (m *MapLiteral) End() Pos {
	return m.RBracePos
}

func (m *MapLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(m)
	defer visitor.Leave(m)
	for _, kv := range m.KeyValues {
		if err := kv.Key.Accept(visitor); err != nil {
			return err
		}
		if err := kv.Value.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitMapLiteral(m)
}

type QueryParam struct {
	LBracePos Pos
	RBracePos Pos
	Name      *Ident
	Type      ColumnType
}

func (q *QueryParam) Pos() Pos {
	return q.LBracePos
}

func (q *QueryParam) End() Pos {
	return q.RBracePos + 1
}

func (q *QueryParam) Accept(visitor ASTVisitor) error {
	visitor.Enter(q)
	defer visitor.Leave(q)
	if err := q.Name.Accept(visitor); err != nil {
		return err
	}
	if err := q.Type.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitQueryParam(q)
}

type TypedPlaceholder struct {
	LeftBracePos  Pos
	RightBracePos Pos
	Name          *Ident
	Type          ColumnType
}

func (t *TypedPlaceholder) Pos() Pos {
	return t.LeftBracePos
}

func (t *TypedPlaceholder) End() Pos {
	return t.RightBracePos
}

func (t *TypedPlaceholder) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Name.Accept(visitor); err != nil {
		return err
	}
	if err := t.Type.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTypedPlaceholder(t)
}

type Literal interface {
	Expr
}

type NumberLiteral struct {
	NumPos  Pos
	NumEnd  Pos
	Literal string
	Base    int
}

func (n *NumberLiteral) Pos() Pos {
	return n.NumPos
}

func (n *NumberLiteral) End() Pos {
	return n.NumEnd
}

func (n *NumberLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	return visitor.VisitNumberLiteral(n)
}

type StringLiteral struct {
	LiteralPos Pos
	LiteralEnd Pos
	Literal    string
}

func (s *StringLiteral) Pos() Pos {
	return s.LiteralPos
}

func (s *StringLiteral) End() Pos {
	return s.LiteralEnd
}

func (s *StringLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	return visitor.VisitStringLiteral(s)
}

type BoolLiteral struct {
	LiteralPos Pos
	LiteralEnd Pos
	Literal    string
}

func (b *BoolLiteral) Pos() Pos {
	return b.LiteralPos
}

func (b *BoolLiteral) End() Pos {
	return b.LiteralEnd
}

func (b *BoolLiteral) Accept(visitor ASTVisitor) error {
	visitor.Enter(b)
	defer visitor.Leave(b)
	return visitor.VisitBoolLiteral(b)
}

type PlaceHolder struct {
	PlaceholderPos Pos
	PlaceHolderEnd Pos
	Type           string
}

func (p *PlaceHolder) Pos() Pos {
	return p.PlaceholderPos
}

func (p *PlaceHolder) End() Pos {
	return p.PlaceHolderEnd
}

func (p *PlaceHolder) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	return visitor.VisitPlaceHolderExpr(p)
}
