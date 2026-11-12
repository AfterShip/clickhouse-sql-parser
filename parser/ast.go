package parser

type Expr interface {
	Pos() Pos
	End() Pos
	FormatSQL(*Formatter)
	Accept(visitor ASTVisitor) error
}

type DDL interface {
	Expr
	Type() string
}
