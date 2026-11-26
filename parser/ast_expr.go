package parser

type OperationExpr struct {
	OperationPos Pos
	Kind         TokenKind
}

func (o *OperationExpr) Pos() Pos {
	return o.OperationPos
}

func (o *OperationExpr) End() Pos {
	return o.OperationPos + Pos(len(o.Kind))
}

func (o *OperationExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	return visitor.VisitOperationExpr(o)
}

type TernaryOperation struct {
	Condition Expr
	TrueExpr  Expr
	FalseExpr Expr
}

func (t *TernaryOperation) Pos() Pos {
	return t.Condition.Pos()
}

func (t *TernaryOperation) End() Pos {
	return t.FalseExpr.End()
}

func (t *TernaryOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.TrueExpr.Accept(visitor); err != nil {
		return err
	}
	if err := t.FalseExpr.Accept(visitor); err != nil {
		return err
	}
	if err := t.Condition.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTernaryExpr(t)
}

type BinaryOperation struct {
	LeftExpr  Expr
	Operation TokenKind
	RightExpr Expr
	HasGlobal bool
	HasNot    bool
}

func (p *BinaryOperation) Pos() Pos {
	return p.LeftExpr.Pos()
}

func (p *BinaryOperation) End() Pos {
	return p.RightExpr.End()
}

func (p *BinaryOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(p)
	defer visitor.Leave(p)
	if err := p.LeftExpr.Accept(visitor); err != nil {
		return err
	}
	if err := p.RightExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitBinaryExpr(p)
}

type IndexOperation struct {
	Object    Expr
	Operation TokenKind
	Index     Expr
}

func (i *IndexOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.Object.Accept(visitor); err != nil {
		return err
	}
	if err := i.Index.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIndexOperation(i)
}

func (i *IndexOperation) Pos() Pos {
	return i.Object.Pos()
}

func (i *IndexOperation) End() Pos {
	return i.Index.End()
}

type ParamExprList struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Items         *ColumnExprList
	ColumnArgList *ColumnArgList
}

func (f *ParamExprList) Pos() Pos {
	return f.LeftParenPos
}

func (f *ParamExprList) End() Pos {
	return f.RightParenPos
}

func (f *ParamExprList) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Items.Accept(visitor); err != nil {
		return err
	}
	if f.ColumnArgList != nil {
		if err := f.ColumnArgList.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitParamExprList(f)
}

type NamedParameterExpr struct {
	NamePos Pos
	Name    *Ident
	Value   Expr
}

func (n *NamedParameterExpr) Pos() Pos {
	return n.NamePos
}

func (n *NamedParameterExpr) End() Pos {
	return n.Value.End()
}

func (n *NamedParameterExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Name.Accept(visitor); err != nil {
		return err
	}
	if err := n.Value.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNamedParameterExpr(n)
}

type ArrayParamList struct {
	LeftBracketPos  Pos
	RightBracketPos Pos
	Items           *ColumnExprList
}

func (a *ArrayParamList) Pos() Pos {
	return a.LeftBracketPos
}

func (a *ArrayParamList) End() Pos {
	return a.RightBracketPos
}

func (a *ArrayParamList) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Items.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitArrayParamList(a)
}

type ObjectParams struct {
	Object Expr
	Params *ArrayParamList
}

func (o *ObjectParams) Pos() Pos {
	return o.Object.Pos()
}

func (o *ObjectParams) End() Pos {
	return o.Params.End()
}

func (o *ObjectParams) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.Object.Accept(visitor); err != nil {
		return err
	}
	if err := o.Params.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitObjectParams(o)
}

type FunctionExpr struct {
	Name   *Ident
	Params *ParamExprList
}

func (f *FunctionExpr) Pos() Pos {
	return f.Name.NamePos
}

func (f *FunctionExpr) End() Pos {
	params := f.Params
	// Bare EXCEPT col uses the column's bounds as its synthetic paren
	// positions, so the right position is already exclusive.
	if items := params.Items; items != nil && params.LeftParenPos == items.Pos() && params.RightParenPos == items.End() {
		return params.RightParenPos
	}
	return params.RightParenPos + 1
}

func (f *FunctionExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Name.Accept(visitor); err != nil {
		return err
	}
	if err := f.Params.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitFunctionExpr(f)
}

type ColumnExpr struct {
	Expr  Expr
	Alias *Ident
}

func (c *ColumnExpr) Pos() Pos {
	return c.Expr.Pos()
}

func (c *ColumnExpr) End() Pos {
	if c.Alias != nil {
		return c.Alias.NameEnd
	}
	return c.Expr.End()
}

func (c *ColumnExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	if c.Alias != nil {
		if err := c.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnExpr(c)
}

type RatioExpr struct {
	Numerator *NumberLiteral
	// numberLiteral (SLASH numberLiteral)?
	Denominator *NumberLiteral
}

func (r *RatioExpr) Pos() Pos {
	return r.Numerator.NumPos
}

func (r *RatioExpr) End() Pos {
	if r.Denominator != nil {
		return r.Denominator.NumEnd
	}
	return r.Numerator.NumEnd
}

func (r *RatioExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(r)
	defer visitor.Leave(r)
	if err := r.Numerator.Accept(visitor); err != nil {
		return err
	}
	if r.Denominator != nil {
		if err := r.Denominator.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRatioExpr(r)
}

type IntervalExpr struct {
	// INTERVAL keyword position which might be omitted(IntervalPos = 0)
	IntervalPos Pos

	Expr Expr
	Unit *Ident
}

func (i *IntervalExpr) Pos() Pos {
	if i.IntervalPos != 0 {
		return i.IntervalPos
	}
	return i.Expr.Pos()
}

func (i *IntervalExpr) End() Pos {
	return i.Unit.End()
}

func (i *IntervalExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := i.Unit.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIntervalExpr(i)
}

type ColumnArgList struct {
	Distinct      bool
	LeftParenPos  Pos
	RightParenPos Pos
	Items         []Expr
}

func (c *ColumnArgList) Pos() Pos {
	return c.LeftParenPos
}

func (c *ColumnArgList) End() Pos {
	return c.RightParenPos
}

func (c *ColumnArgList) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, item := range c.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnArgList(c)
}

type ColumnExprList struct {
	ListPos     Pos
	ListEnd     Pos
	HasDistinct bool
	Items       []Expr
}

func (c *ColumnExprList) Pos() Pos {
	return c.ListPos
}

func (c *ColumnExprList) End() Pos {
	return c.ListEnd
}

func (c *ColumnExprList) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	for _, item := range c.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitColumnExprList(c)
}

type WhenClause struct {
	WhenPos Pos
	ThenPos Pos
	When    Expr
	Then    Expr
	ElsePos Pos
	Else    Expr
}

func (w *WhenClause) Pos() Pos {
	return w.WhenPos
}

func (w *WhenClause) End() Pos {
	if w.Else != nil {
		return w.Else.End()
	}
	return w.Then.End()
}

func (w *WhenClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.When.Accept(visitor); err != nil {
		return err
	}
	if err := w.Then.Accept(visitor); err != nil {
		return err
	}
	if w.Else != nil {
		if err := w.Else.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitWhenExpr(w)
}

type CaseExpr struct {
	CasePos Pos
	EndPos  Pos
	Expr    Expr // optional
	Whens   []*WhenClause
	ElsePos Pos
	Else    Expr
}

func (c *CaseExpr) Pos() Pos {
	return c.CasePos
}

func (c *CaseExpr) End() Pos {
	return c.EndPos
}

func (c *CaseExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if c.Expr != nil {
		if err := c.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	for _, when := range c.Whens {
		if err := when.Accept(visitor); err != nil {
			return err
		}
	}
	if c.Else != nil {
		if err := c.Else.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitCaseExpr(c)
}

type CastExpr struct {
	CastPos   Pos
	Expr      Expr
	Separator string
	AsPos     Pos
	AsType    Expr
}

func (c *CastExpr) Pos() Pos {
	return c.CastPos
}

func (c *CastExpr) End() Pos {
	return c.AsType.End()
}

func (c *CastExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := c.AsType.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCastExpr(c)
}

type IsNullExpr struct {
	IsPos   Pos // position of the IS keyword
	NullEnd Pos // end of the NULL keyword
	Expr    Expr
}

func (n *IsNullExpr) Pos() Pos {
	return n.Expr.Pos()
}

func (n *IsNullExpr) End() Pos {
	return n.NullEnd
}

func (n *IsNullExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIsNullExpr(n)
}

type IsNotNullExpr struct {
	IsPos   Pos // position of the IS keyword
	NullEnd Pos // end of the NULL keyword
	Expr    Expr
}

func (n *IsNotNullExpr) Pos() Pos {
	return n.Expr.Pos()
}

func (n *IsNotNullExpr) End() Pos {
	return n.NullEnd
}

func (n *IsNotNullExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIsNotNullExpr(n)
}

type AliasExpr struct {
	Expr     Expr
	AliasPos Pos
	Alias    Expr
}

func (a *AliasExpr) Pos() Pos {
	return a.AliasPos
}

func (a *AliasExpr) End() Pos {
	return a.Alias.End()
}

func (a *AliasExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(a)
	defer visitor.Leave(a)
	if err := a.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := a.Alias.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitAliasExpr(a)
}

type BetweenClause struct {
	Expr    Expr
	Not     bool
	Between Expr
	AndPos  Pos
	And     Expr
}

func (f *BetweenClause) Pos() Pos {
	if f.Expr != nil {
		return f.Expr.Pos()
	}
	return f.Between.Pos()
}

func (f *BetweenClause) End() Pos {
	return f.And.End()
}

func (f *BetweenClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if f.Expr != nil {
		if err := f.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	if err := f.Between.Accept(visitor); err != nil {
		return err
	}
	if err := f.And.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitBetweenClause(f)
}

type NotExpr struct {
	NotPos Pos
	Expr   Expr
}

func (n *NotExpr) Pos() Pos {
	return n.NotPos
}

func (n *NotExpr) End() Pos {
	return n.Expr.End()
}

func (n *NotExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNotExpr(n)
}

type NegateExpr struct {
	NegatePos Pos
	Expr      Expr
}

func (n *NegateExpr) Pos() Pos {
	return n.NegatePos
}

func (n *NegateExpr) End() Pos {
	return n.Expr.End()
}

func (n *NegateExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitNegateExpr(n)
}

type GlobalInOperation struct {
	GlobalPos Pos
	Expr      Expr
}

func (g *GlobalInOperation) Pos() Pos {
	return g.GlobalPos
}

func (g *GlobalInOperation) End() Pos {
	return g.Expr.End()
}

func (g *GlobalInOperation) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	if err := g.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitGlobalInExpr(g)
}

type IntervalFrom struct {
	Interval *Ident
	FromPos  Pos
	FromExpr Expr
}

func (i *IntervalFrom) Pos() Pos {
	return i.Interval.NamePos
}

func (i *IntervalFrom) End() Pos {
	return i.FromExpr.End()
}

func (i *IntervalFrom) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if i.Interval != nil {
		if err := i.Interval.Accept(visitor); err != nil {
			return err
		}
	}
	if err := i.FromExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitIntervalFrom(i)
}

type ExtractExpr struct {
	ExtractPos Pos
	ExtractEnd Pos
	Parameters []Expr
}

func (e *ExtractExpr) Pos() Pos {
	return e.ExtractPos
}

func (e *ExtractExpr) End() Pos {
	return e.ExtractEnd
}

func (e *ExtractExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(e)
	defer visitor.Leave(e)
	for _, param := range e.Parameters {
		if err := param.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitExtractExpr(e)
}

type UnaryExpr struct {
	UnaryPos Pos
	Kind     TokenKind
	Expr     Expr
}

func (n *UnaryExpr) Pos() Pos {
	return n.UnaryPos
}

func (n *UnaryExpr) End() Pos {
	return n.Expr.End()
}

func (n *UnaryExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(n)
	defer visitor.Leave(n)
	if err := n.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUnaryExpr(n)
}
