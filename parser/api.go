package parser

import (
	"errors"
	"fmt"
)

// ParseStmt parses exactly one statement from the input. It returns an error
// if the input contains no statement or more than one statement. To parse a
// script with multiple statements, use NewParser(sql).ParseStmts().
func ParseStmt(sql string) (Expr, error) {
	stmts, err := NewParser(sql).ParseStmts()
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 {
		return nil, errors.New("no statement found in input")
	}
	if len(stmts) > 1 {
		return nil, fmt.Errorf("expected exactly one statement, but found %d", len(stmts))
	}
	return stmts[0], nil
}

// ParseExpr parses a single expression fragment, such as a column reference,
// a function call or an arithmetic expression — e.g. `toDate(created_at) + 1`.
// The whole input must be consumed by the expression.
func ParseExpr(sql string) (Expr, error) {
	p := NewParser(sql)
	if err := p.lexer.consumeToken(); err != nil {
		return nil, p.wrapError(err)
	}
	expr, err := p.parseExpr(p.Pos())
	if err != nil {
		return nil, p.wrapError(err)
	}
	if p.last() != nil {
		return nil, p.wrapError(fmt.Errorf("unexpected token after expression: %q", p.lastTokenString()))
	}
	return expr, nil
}

// FormatBeautify renders an expression into multi-line indented SQL. It is a
// convenience for NewFormatter().WithBeautify(); use Format for compact
// single-line SQL.
func FormatBeautify(expr Expr) string {
	formatter := NewFormatter().WithBeautify()
	formatter.WriteExpr(expr)
	return formatter.String()
}
