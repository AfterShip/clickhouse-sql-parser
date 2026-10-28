package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatter_NumericTupleFieldBoundary(t *testing.T) {
	for _, query := range []string{
		"SELECT 0 .A",
		"SELECT -1 .A",
		"ALTER POLICY p ON db.t USING 0 .A",
	} {
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
