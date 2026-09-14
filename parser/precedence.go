package parser

import "strings"

// Keep the existing exported constants and their numeric values. Operators
// that ClickHouse binds equally share a level in operatorPrecedence below.
const (
	PrecedenceUnknown = iota
	PrecedenceArrow
	PrecedenceQuery
	PrecedenceOr
	PrecedenceAnd
	PrecedenceNot
	PrecedenceGlobal
	PrecedenceIs
	PrecedenceCompare
	PrecedenceBetweenLike
	precedenceIn
	PrecedenceConcat
	PrecedenceAddSub
	PrecedenceMulDivMod
	PrecedenceBracket
	PrecedenceDot
	PrecedenceDoubleColon
)

// operatorPrecedence is shared by the parser and formatter. Postfix operators
// associate in source order; comparisons, LIKE and IN share one binding level.
func operatorPrecedence(op TokenKind) int {
	switch strings.ToUpper(string(op)) {
	case "->":
		return PrecedenceArrow
	case "?":
		return PrecedenceQuery
	case "OR":
		return PrecedenceOr
	case "AND":
		return PrecedenceAnd
	case KeywordNot:
		return PrecedenceNot
	case "IS", "BETWEEN", "NOT BETWEEN":
		// BETWEEN's bounds stop before IS NULL and another BETWEEN.
		return PrecedenceIs
	case "=", "==", "!=", "<>", "<", "<=", ">", ">=",
		"IN", "NOT IN", "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE", "REGEXP":
		return PrecedenceCompare
	case "||":
		return PrecedenceConcat
	case "+", "-":
		return PrecedenceAddSub
	case "*", "/", "%":
		return PrecedenceMulDivMod
	case ".", ".:", "::", "(", "[":
		return PrecedenceBracket
	default:
		return PrecedenceUnknown
	}
}

func expressionPrecedence(expr Expr) int {
	switch e := expr.(type) {
	case *ColumnExpr:
		return expressionPrecedence(e.Expr)
	case *BinaryOperation:
		if precedence := operatorPrecedence(e.Operation); precedence != PrecedenceUnknown {
			return precedence
		}
	case *UnaryExpr:
		switch strings.ToUpper(string(e.Kind)) {
		case "+", "-":
			// This is the operand's lower bound; the formatter excludes an
			// equal-precedence binary operand when writing a unary sign.
			return PrecedenceMulDivMod
		case KeywordNot:
			return PrecedenceNot
		}
	case *TernaryOperation:
		return PrecedenceQuery
	case *BetweenClause:
		return operatorPrecedence(TokenKind(KeywordBetween))
	case *IsNullExpr, *IsNotNullExpr:
		return operatorPrecedence(TokenKind(KeywordIs))
	case *IndexOperation, *ObjectParams:
		return PrecedenceBracket
	}
	// Primary expressions and explicit parentheses already protect grouping.
	return PrecedenceDoubleColon + 1
}

func (f *Formatter) writeOperand(expr Expr, precedence int, parenEqual bool) {
	childPrecedence := expressionPrecedence(expr)
	paren := childPrecedence < precedence || (parenEqual && childPrecedence == precedence)
	if paren {
		f.WriteByte('(')
	}
	f.WriteExpr(expr)
	if paren {
		f.WriteByte(')')
	}
}
