package parser

var intervalType = NewSet("SECOND", "MINUTE", "HOUR", "DAY", "WEEK", "MONTH", "QUARTER", "YEAR")

type OpType string

const (
	// Comparison operators
	opTypeEQ    TokenKind = "="
	opTypeNE    TokenKind = "!="
	opTypeLT    TokenKind = "<"
	opTypeLE    TokenKind = "<="
	opTypeGT    TokenKind = ">"
	opTypeGE    TokenKind = ">="
	opTypeQuery           = "?"

	// Arithmetic operators
	opTypePlus  TokenKind = "+"
	opTypeMinus TokenKind = "-"
	opTypeMul   TokenKind = "*"
	opTypeDiv   TokenKind = "/"
	opTypeMod   TokenKind = "%"

	// Logical operators
	opTypeAnd TokenKind = "AND"
	opTypeOr  TokenKind = "OR"
)
