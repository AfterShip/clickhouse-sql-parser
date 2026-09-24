package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectFormatFunctionAndOutputClause(t *testing.T) {
	for _, sql := range []string{
		"SELECT format('{}', 1)",
		"SELECT 1, format('{}', 1)",
		"SELECT 1, format('{}', 1) FORMAT JSON",
		"SELECT 1 FORMAT JSON",
	} {
		t.Run(sql, func(t *testing.T) {
			stmts, err := NewParser(sql).ParseStmts()
			require.NoError(t, err)
			require.Len(t, stmts, 1)
			formatted := Format(stmts[0])
			again, err := NewParser(formatted).ParseStmts()
			require.NoError(t, err)
			require.Len(t, again, 1)
			require.Equal(t, formatted, Format(again[0]))
		})
	}
}
