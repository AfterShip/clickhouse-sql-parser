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

// operatorPrecedence defines expression binding levels. Postfix operators
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
