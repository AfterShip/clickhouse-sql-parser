package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShowAccess(t *testing.T) {
	stmts, err := NewParser("SHOW ACCESS").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	show, ok := stmts[0].(*ShowStmt)
	require.True(t, ok)
	require.Equal(t, "ACCESS", show.ShowType)
	require.Equal(t, "SHOW ACCESS", Format(show))

	for _, sql := range []string{"SHOW ACCESS JUNK", "SHOW ACCESS LIMIT 1", "SHOW ACCESS FROM system.users"} {
		_, err := NewParser(sql).ParseStmts()
		require.Error(t, err, sql)
	}
}
