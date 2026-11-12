package parser

type OrderDirection string

const (
	OrderDirectionNone OrderDirection = ""
	OrderDirectionAsc  OrderDirection = "ASC"
	OrderDirectionDesc OrderDirection = "DESC"
)

type SelectItem struct {
	Expr Expr
	// Please refer: https://clickhouse.com/docs/en/sql-reference/statements/select#select-modifiers
	Modifiers []*FunctionExpr
	Alias     *Ident
}

func (s *SelectItem) Pos() Pos {
	return s.Expr.Pos()
}

func (s *SelectItem) End() Pos {
	if s.Alias != nil {
		return s.Alias.End()
	}
	if len(s.Modifiers) > 0 {
		return s.Modifiers[len(s.Modifiers)-1].End()
	}
	return s.Expr.End()
}

func (s *SelectItem) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Expr.Accept(visitor); err != nil {
		return err
	}
	for _, modifier := range s.Modifiers {
		if err := modifier.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Alias != nil {
		if err := s.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSelectItem(s)
}

type JoinTableExpr struct {
	Table        *TableExpr
	StatementEnd Pos
	SampleRatio  *SampleClause
	HasFinal     bool
}

func (j *JoinTableExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if err := j.Table.Accept(visitor); err != nil {
		return err
	}
	if j.SampleRatio != nil {
		if err := j.SampleRatio.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitJoinTableExpr(j)
}

func (j *JoinTableExpr) Pos() Pos {
	return j.Table.Pos()
}

func (j *JoinTableExpr) End() Pos {
	return j.StatementEnd
}

type TableArgListExpr struct {
	LeftParenPos  Pos
	RightParenPos Pos
	Args          []Expr
}

func (t *TableArgListExpr) Pos() Pos {
	return t.LeftParenPos
}

func (t *TableArgListExpr) End() Pos {
	return t.RightParenPos
}

func (t *TableArgListExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	for _, arg := range t.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTableArgListExpr(t)
}

type TableFunctionExpr struct {
	Name Expr
	Args *TableArgListExpr
}

func (t *TableFunctionExpr) Pos() Pos {
	return t.Name.Pos()
}

func (t *TableFunctionExpr) End() Pos {
	return t.Args.End()
}

func (t *TableFunctionExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Name.Accept(visitor); err != nil {
		return err
	}
	if err := t.Args.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTableFunctionExpr(t)
}

type Fill struct {
	FillPos   Pos
	From      Expr // optional
	To        Expr // optional
	Step      Expr // optional
	Staleness Expr // optional
}

func (f *Fill) Pos() Pos {
	return f.FillPos
}

func (f *Fill) End() Pos {
	if f.Staleness != nil {
		return f.Staleness.End()
	}
	if f.Step != nil {
		return f.Step.End()
	}
	if f.To != nil {
		return f.To.End()
	}
	if f.From != nil {
		return f.From.End()
	}
	return f.FillPos + Pos(len("FILL"))
}

func (f *Fill) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if f.From != nil {
		if err := f.From.Accept(visitor); err != nil {
			return err
		}
	}
	if f.To != nil {
		if err := f.To.Accept(visitor); err != nil {
			return err
		}
	}
	if f.Step != nil {
		if err := f.Step.Accept(visitor); err != nil {
			return err
		}
	}
	if f.Staleness != nil {
		if err := f.Staleness.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitFill(f)
}

type OrderExpr struct {
	OrderPos  Pos
	Expr      Expr
	Alias     *Ident
	Direction OrderDirection
	Fill      *Fill // optional WITH FILL clause
}

func (o *OrderExpr) Pos() Pos {
	return o.OrderPos
}

func (o *OrderExpr) End() Pos {
	if o.Fill != nil {
		return o.Fill.End()
	}
	if o.Alias != nil {
		return o.Alias.End()
	}
	return o.Expr.End()
}

func (o *OrderExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.Expr.Accept(visitor); err != nil {
		return err
	}
	if o.Alias != nil {
		if err := o.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Fill != nil {
		if err := o.Fill.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitOrderByExpr(o)
}

type InterpolateItem struct {
	Column *Ident
	Expr   Expr // optional AS expression
}

func (i *InterpolateItem) Pos() Pos {
	return i.Column.Pos()
}

func (i *InterpolateItem) End() Pos {
	if i.Expr != nil {
		return i.Expr.End()
	}
	return i.Column.End()
}

func (i *InterpolateItem) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	if err := i.Column.Accept(visitor); err != nil {
		return err
	}
	if i.Expr != nil {
		if err := i.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitInterpolateItem(i)
}

type InterpolateClause struct {
	InterpolatePos Pos
	ListEnd        Pos
	Items          []*InterpolateItem // can be nil for INTERPOLATE without columns
}

func (i *InterpolateClause) Pos() Pos {
	return i.InterpolatePos
}

func (i *InterpolateClause) End() Pos {
	return i.ListEnd
}

func (i *InterpolateClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(i)
	defer visitor.Leave(i)
	for _, item := range i.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitInterpolateClause(i)
}

type OrderByClause struct {
	OrderPos    Pos
	ListEnd     Pos
	Items       []Expr
	Interpolate *InterpolateClause // optional INTERPOLATE clause
}

func (o *OrderByClause) Pos() Pos {
	return o.OrderPos
}

func (o *OrderByClause) End() Pos {
	if o.Interpolate != nil {
		return o.Interpolate.End()
	}
	return o.ListEnd
}

func (o *OrderByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	for _, item := range o.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	if o.Interpolate != nil {
		if err := o.Interpolate.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitOrderByListExpr(o)
}

type SettingExpr struct {
	SettingsPos Pos
	Name        *Ident
	Expr        Expr
}

func (s *SettingExpr) Pos() Pos {
	return s.SettingsPos
}

func (s *SettingExpr) End() Pos {
	return s.Expr.End()
}

func (s *SettingExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Name.Accept(visitor); err != nil {
		return err
	}
	if err := s.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSettingsExpr(s)
}

type SettingsClause struct {
	SettingsPos Pos
	ListEnd     Pos
	Items       []*SettingExpr
}

func (s *SettingsClause) Pos() Pos {
	return s.SettingsPos
}

func (s *SettingsClause) End() Pos {
	return s.ListEnd
}

func (s *SettingsClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	for _, item := range s.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSettingsExprList(s)
}

type WindowFunctionExpr struct {
	Function *FunctionExpr
	OverPos  Pos
	OverExpr Expr
}

func (w *WindowFunctionExpr) Pos() Pos {
	return w.Function.Pos()
}

func (w *WindowFunctionExpr) End() Pos {
	return w.OverExpr.End()
}

func (w *WindowFunctionExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.Function.Accept(visitor); err != nil {
		return err
	}
	if err := w.OverExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFunctionExpr(w)
}

type WithClause struct {
	WithPos Pos
	EndPos  Pos
	CTEs    []*CTEStmt
}

func (w *WithClause) Pos() Pos {
	return w.WithPos
}

func (w *WithClause) End() Pos {
	return w.EndPos
}

func (w *WithClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	for _, cte := range w.CTEs {
		if err := cte.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitWithExpr(w)
}

type TopClause struct {
	TopPos   Pos
	TopEnd   Pos
	Number   *NumberLiteral
	WithTies bool
}

func (t *TopClause) Pos() Pos {
	return t.TopPos
}

func (t *TopClause) End() Pos {
	return t.TopEnd
}

func (t *TopClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Number.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTopExpr(t)
}

type TableExpr struct {
	TablePos Pos
	TableEnd Pos
	Alias    *AliasExpr
	Expr     Expr
	HasFinal bool
}

func (t *TableExpr) Pos() Pos {
	return t.TablePos
}

func (t *TableExpr) End() Pos {
	return t.TableEnd
}

func (t *TableExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(t)
	defer visitor.Leave(t)
	if err := t.Expr.Accept(visitor); err != nil {
		return err
	}
	if t.Alias != nil {
		if err := t.Alias.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitTableExpr(t)
}

type OnClause struct {
	OnPos Pos
	On    *ColumnExprList
}

func (o *OnClause) Pos() Pos {
	return o.OnPos
}

func (o *OnClause) End() Pos {
	return o.On.End()
}

func (o *OnClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(o)
	defer visitor.Leave(o)
	if err := o.On.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitOnExpr(o)
}

type UsingClause struct {
	UsingPos Pos
	Using    *ColumnExprList
}

func (u *UsingClause) Pos() Pos {
	return u.UsingPos
}

func (u *UsingClause) End() Pos {
	return u.Using.End()
}

func (u *UsingClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(u)
	defer visitor.Leave(u)
	if err := u.Using.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUsingExpr(u)
}

type JoinExpr struct {
	JoinPos     Pos
	Left        Expr
	Right       Expr
	Modifiers   []string
	Constraints Expr
}

func (j *JoinExpr) Pos() Pos {
	return j.JoinPos
}

func (j *JoinExpr) End() Pos {
	// Return the rightmost position
	if j.Right != nil {
		return j.Right.End()
	}
	if j.Constraints != nil {
		return j.Constraints.End()
	}
	return j.Left.End()
}

func (j *JoinExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if err := j.Left.Accept(visitor); err != nil {
		return err
	}
	if j.Right != nil {
		if err := j.Right.Accept(visitor); err != nil {
			return err
		}
	}
	if j.Constraints != nil {
		if err := j.Constraints.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitJoinExpr(j)
}

type JoinConstraintClause struct {
	ConstraintPos Pos
	On            *ColumnExprList
	Using         *ColumnExprList
}

func (j *JoinConstraintClause) Pos() Pos {
	return j.ConstraintPos
}

func (j *JoinConstraintClause) End() Pos {
	if j.On != nil {
		return j.On.End()
	}
	return j.Using.End()
}

func (j *JoinConstraintClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(j)
	defer visitor.Leave(j)
	if j.On != nil {
		if err := j.On.Accept(visitor); err != nil {
			return err
		}
	}
	if j.Using != nil {
		if err := j.Using.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitJoinConstraintExpr(j)
}

type FromClause struct {
	FromPos Pos
	Expr    Expr
}

func (f *FromClause) Pos() Pos {
	return f.FromPos
}

func (f *FromClause) End() Pos {
	return f.Expr.End()
}

func (f *FromClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitFromExpr(f)
}

type WhereClause struct {
	WherePos Pos
	Expr     Expr
}

func (w *WhereClause) Pos() Pos {
	return w.WherePos
}

func (w *WhereClause) End() Pos {
	return w.Expr.End()
}

func (w *WhereClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWhereExpr(w)
}

type PrewhereClause struct {
	PrewherePos Pos
	Expr        Expr
}

func (w *PrewhereClause) Pos() Pos {
	return w.PrewherePos
}

func (w *PrewhereClause) End() Pos {
	return w.Expr.End()
}

func (w *PrewhereClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if err := w.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPrewhereExpr(w)
}

type GroupByClause struct {
	GroupByPos    Pos
	GroupByEnd    Pos
	AggregateType string
	Expr          Expr
	WithCube      bool
	WithRollup    bool
	WithTotals    bool
}

func (g *GroupByClause) Pos() Pos {
	return g.GroupByPos
}

func (g *GroupByClause) End() Pos {
	return g.GroupByEnd
}

func (g *GroupByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(g)
	defer visitor.Leave(g)
	if g.Expr != nil {
		if err := g.Expr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitGroupByExpr(g)
}

type HavingClause struct {
	HavingPos Pos
	Expr      Expr
}

func (h *HavingClause) Pos() Pos {
	return h.HavingPos
}

func (h *HavingClause) End() Pos {
	return h.Expr.End()
}

func (h *HavingClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(h)
	defer visitor.Leave(h)
	if err := h.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitHavingExpr(h)
}

type LimitClause struct {
	LimitPos Pos
	Limit    Expr
	Offset   Expr
}

func (l *LimitClause) Pos() Pos {
	return l.LimitPos
}

func (l *LimitClause) End() Pos {
	if l.Offset != nil {
		return l.Offset.End()
	}
	return l.Limit.End()
}

func (l *LimitClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(l)
	defer visitor.Leave(l)
	if l.Limit != nil {
		if err := l.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if l.Offset != nil {
		if err := l.Offset.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitLimitExpr(l)
}

type LimitByClause struct {
	Limit  *LimitClause
	ByExpr *ColumnExprList
}

func (l *LimitByClause) Pos() Pos {
	return l.Limit.Pos()
}

func (l *LimitByClause) End() Pos {
	if l.ByExpr != nil {
		return l.ByExpr.End()
	}
	if l.Limit != nil {
		return l.Limit.End()
	}
	return l.Limit.End()
}

func (l *LimitByClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(l)
	defer visitor.Leave(l)
	if l.Limit != nil {
		if err := l.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if l.ByExpr != nil {
		if err := l.ByExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitLimitByExpr(l)
}

type WindowExpr struct {
	LeftParenPos  Pos
	RightParenPos Pos
	WindowName    *Ident
	PartitionBy   *PartitionByClause
	OrderBy       *OrderByClause
	Frame         *WindowFrameClause
}

func (w *WindowExpr) Pos() Pos {
	return w.LeftParenPos
}

func (w *WindowExpr) End() Pos {
	return w.RightParenPos
}

func (w *WindowExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	if w.WindowName != nil {
		if err := w.WindowName.Accept(visitor); err != nil {
			return err
		}
	}
	if w.PartitionBy != nil {
		if err := w.PartitionBy.Accept(visitor); err != nil {
			return err
		}
	}
	if w.OrderBy != nil {
		if err := w.OrderBy.Accept(visitor); err != nil {
			return err
		}
	}
	if w.Frame != nil {
		if err := w.Frame.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitWindowConditionExpr(w)
}

type WindowDefinition struct {
	Name  *Ident
	AsPos Pos
	Expr  *WindowExpr
}

func (w *WindowDefinition) Pos() Pos {
	if w == nil || w.Name == nil {
		return 0
	}
	return w.Name.Pos()
}

func (w *WindowDefinition) End() Pos {
	if w == nil || w.Expr == nil {
		return 0
	}
	return w.Expr.End()
}

type WindowClause struct {
	WindowPos Pos
	EndPos    Pos
	Windows   []*WindowDefinition
}

func (w *WindowClause) Pos() Pos {
	return w.WindowPos
}

func (w *WindowClause) End() Pos {
	if w.EndPos != 0 {
		return w.EndPos
	}
	if len(w.Windows) == 0 {
		return w.WindowPos
	}
	return w.Windows[len(w.Windows)-1].End()
}

func (w *WindowClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(w)
	defer visitor.Leave(w)
	for _, window := range w.Windows {
		if window == nil {
			continue
		}
		if window.Name != nil {
			if err := window.Name.Accept(visitor); err != nil {
				return err
			}
		}
		if window.Expr != nil {
			if err := window.Expr.Accept(visitor); err != nil {
				return err
			}
		}
	}
	return visitor.VisitWindowExpr(w)
}

type WindowFrameClause struct {
	FramePos Pos
	Type     string
	Extend   Expr
}

func (f *WindowFrameClause) Pos() Pos {
	return f.FramePos
}

func (f *WindowFrameClause) End() Pos {
	return f.Extend.End()
}

func (f *WindowFrameClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Extend.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameExpr(f)
}

type WindowFrameExtendExpr struct {
	Expr      Expr
	Direction string
	EndPos    Pos
}

func (f *WindowFrameExtendExpr) Pos() Pos {
	return f.Expr.Pos()
}

func (f *WindowFrameExtendExpr) End() Pos {
	if f.EndPos != 0 {
		return f.EndPos
	}
	return f.Expr.End()
}

func (f *WindowFrameExtendExpr) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameExtendExpr(f)
}

type WindowFrameCurrentRow struct {
	CurrentPos Pos
	RowEnd     Pos
}

func (f *WindowFrameCurrentRow) Pos() Pos {
	return f.CurrentPos
}

func (f *WindowFrameCurrentRow) End() Pos {
	return f.RowEnd
}

func (f *WindowFrameCurrentRow) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	return visitor.VisitWindowFrameCurrentRow(f)
}

type WindowFrameUnbounded struct {
	UnboundedPos Pos
	UnboundedEnd Pos
	Direction    string
}

func (f *WindowFrameUnbounded) Pos() Pos {
	return f.UnboundedPos
}

func (f *WindowFrameUnbounded) End() Pos {
	return f.UnboundedEnd
}

func (f *WindowFrameUnbounded) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	return visitor.VisitWindowFrameUnbounded(f)
}

type WindowFrameNumber struct {
	Number    *NumberLiteral
	EndPos    Pos
	Direction string
}

func (f *WindowFrameNumber) Pos() Pos {
	return f.Number.Pos()
}

func (f *WindowFrameNumber) End() Pos {
	return f.EndPos
}

func (f *WindowFrameNumber) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Number.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameNumber(f)
}

type WindowFrameParam struct {
	Param     *QueryParam
	EndPos    Pos
	Direction string
}

func (f *WindowFrameParam) Pos() Pos {
	return f.Param.Pos()
}

func (f *WindowFrameParam) End() Pos {
	return f.EndPos
}

func (f *WindowFrameParam) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Param.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitWindowFrameParam(f)
}

type SelectQuery struct {
	SelectPos    Pos
	StatementEnd Pos
	// InnerQuery makes this node a parenthesized group wrapping that query:
	// (SELECT 1 UNION ALL SELECT 2). When set, every other clause field is
	// empty except Settings, Format and the set-operation fields below, which
	// bind clauses that follow the closing ')'. SelectPos and StatementEnd
	// span the parentheses and any trailing clause.
	InnerQuery    *SelectQuery `json:",omitempty"`
	With          *WithClause
	Top           *TopClause
	HasDistinct   bool
	DistinctOn    *DistinctOn
	SelectItems   []*SelectItem
	From          *FromClause
	Window        *WindowClause
	Prewhere      *PrewhereClause
	Where         *WhereClause
	GroupBy       *GroupByClause
	WithTotal     bool
	Having        *HavingClause
	OrderBy       *OrderByClause
	LimitBy       *LimitByClause
	Limit         *LimitClause
	Settings      *SettingsClause
	Format        *FormatClause
	UnionAll      *SelectQuery
	UnionDistinct *SelectQuery
	Except        *SelectQuery
	Intersect     *SelectQuery
}

func (s *SelectQuery) Pos() Pos {
	return s.SelectPos
}

func (s *SelectQuery) End() Pos {
	return s.StatementEnd
}

func (s *SelectQuery) Accept(visitor ASTVisitor) error { // nolint:funlen
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.InnerQuery != nil {
		if err := s.InnerQuery.Accept(visitor); err != nil {
			return err
		}
	}
	if s.With != nil {
		if err := s.With.Accept(visitor); err != nil {
			return err
		}
	}
	if s.DistinctOn != nil {
		if err := s.DistinctOn.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Top != nil {
		if err := s.Top.Accept(visitor); err != nil {
			return err
		}
	}
	if s.SelectItems != nil {
		for _, item := range s.SelectItems {
			if err := item.Accept(visitor); err != nil {
				return err
			}
		}
	}
	if s.From != nil {
		if err := s.From.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Window != nil {
		if err := s.Window.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Prewhere != nil {
		if err := s.Prewhere.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Where != nil {
		if err := s.Where.Accept(visitor); err != nil {
			return err
		}
	}
	if s.GroupBy != nil {
		if err := s.GroupBy.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Having != nil {
		if err := s.Having.Accept(visitor); err != nil {
			return err
		}
	}
	if s.OrderBy != nil {
		if err := s.OrderBy.Accept(visitor); err != nil {
			return err
		}
	}
	if s.LimitBy != nil {
		if err := s.LimitBy.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Limit != nil {
		if err := s.Limit.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Settings != nil {
		if err := s.Settings.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Format != nil {
		if err := s.Format.Accept(visitor); err != nil {
			return err
		}
	}
	if s.UnionAll != nil {
		if err := s.UnionAll.Accept(visitor); err != nil {
			return err
		}
	}
	if s.UnionDistinct != nil {
		if err := s.UnionDistinct.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Except != nil {
		if err := s.Except.Accept(visitor); err != nil {
			return err
		}
	}
	if s.Intersect != nil {
		if err := s.Intersect.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSelectQuery(s)
}

type DistinctOn struct {
	Idents        []*NestedIdentifier
	DistinctOnPos Pos
	DistinctOnEnd Pos
}

func (s *DistinctOn) Pos() Pos {
	return s.DistinctOnPos
}

func (s *DistinctOn) End() Pos {
	return s.DistinctOnEnd
}

func (s *DistinctOn) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	for _, ident := range s.Idents {
		if err := ident.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitDistinctOn(s)
}

type SubQuery struct {
	HasParen bool
	Select   *SelectQuery
}

func (s *SubQuery) Pos() Pos {
	return s.Select.Pos()
}

func (s *SubQuery) End() Pos {
	return s.Select.End()
}

func (s *SubQuery) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if s.Select != nil {
		if err := s.Select.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSubQueryExpr(s)
}

type CTEStmt struct {
	CTEPos Pos
	Expr   Expr
	Alias  Expr
}

func (c *CTEStmt) Pos() Pos {
	return c.CTEPos
}

func (c *CTEStmt) End() Pos {
	return c.Expr.End()
}

func (c *CTEStmt) Accept(visitor ASTVisitor) error {
	visitor.Enter(c)
	defer visitor.Leave(c)
	if err := c.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := c.Alias.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCTEExpr(c)
}

type FormatClause struct {
	FormatPos Pos
	Format    *Ident
}

func (f *FormatClause) Pos() Pos {
	return f.FormatPos
}

func (f *FormatClause) End() Pos {
	return f.Format.End()
}

func (f *FormatClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(f)
	defer visitor.Leave(f)
	if err := f.Format.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitFormatExpr(f)
}

type SampleClause struct {
	SamplePos Pos
	Ratio     *RatioExpr
	Offset    *RatioExpr
}

func (s *SampleClause) Pos() Pos {
	return s.SamplePos
}

func (s *SampleClause) End() Pos {
	if s.Offset != nil {
		return s.Offset.End()
	}
	return s.Ratio.End()
}

func (s *SampleClause) Accept(visitor ASTVisitor) error {
	visitor.Enter(s)
	defer visitor.Leave(s)
	if err := s.Ratio.Accept(visitor); err != nil {
		return err
	}
	if s.Offset != nil {
		if err := s.Offset.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSampleRatioExpr(s)
}
