package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLowercaseArrayJoinExpressionList(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1 FROM t ARRAY JOIN arrayConcat(xs, [x]) AS ys",
		"SELECT 1 FROM t array join arrayConcat(xs, [x]) AS ys",
	} {
		stmts, err := NewParser(sql).ParseStmts()
		require.NoError(t, err)
		require.Len(t, stmts, 1)
		require.Equal(t, Format(stmts[0]), Format(parseOneStmt(t, Format(stmts[0]))))
	}
}
