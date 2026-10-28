package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_RowPolicyStatements(t *testing.T) {
	t.Parallel()
	queries := []string{
		"CREATE ROW POLICY p ON db.t FOR SELECT USING id = 42 TO r",
		"CREATE POLICY IF NOT EXISTS p ON db.t, db.t_local USING id = toUInt64(getSetting('filter_id')) TO r",
		"ALTER ROW POLICY IF EXISTS p ON db.t USING NONE TO NONE",
		"DROP POLICY IF EXISTS p ON db.t FROM local_directory ON CLUSTER 'c'",
		"SHOW CREATE ROW POLICY p ON db.t",
		"SHOW ROW POLICIES ON db.*",
		"SHOW CREATE ROW POLICIES ON *.*",
		"SHOW ROW POLICIES ON *.*",
	}
	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			statements, err := NewParser(query).ParseStmts()
			require.NoError(t, err)
			require.Len(t, statements, 1)
			formatted := Format(statements[0])
			again, err := NewParser(formatted).ParseStmts()
			require.NoError(t, err)
			require.Len(t, again, 1)
			require.Equal(t, formatted, Format(again[0]))
		})
	}
}

func TestParser_RowPolicyMixedStatements(t *testing.T) {
	query := "CREATE TABLE db.t (id UInt64) ENGINE = MergeTree ORDER BY id; CREATE ROW POLICY p ON db.t USING id = 42 TO r; SELECT id FROM db.t"
	statements, err := NewParser(query).ParseStmts()
	require.NoError(t, err)
	require.Len(t, statements, 3)
	require.True(t, strings.HasPrefix(Format(statements[1]), "CREATE ROW POLICY"))

	statements, err = NewParser("SELECT 1; CREATE ROW POLICY p ON *.*").ParseStmts()
	require.Error(t, err)
	require.Empty(t, statements, "a failed policy must not return a successfully parsed prefix")
}
