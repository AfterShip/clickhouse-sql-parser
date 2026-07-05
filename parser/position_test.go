package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func parseOneStmt(t *testing.T, sql string) Expr {
	t.Helper()
	stmts, err := NewParser(sql).ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	return stmts[0]
}

func TestRatioExprPositions(t *testing.T) {
	sql := "SELECT * FROM t SAMPLE 1/2"
	stmt := parseOneStmt(t, sql).(*SelectQuery)
	sample := stmt.From.Expr.(*JoinTableExpr).SampleRatio
	require.NotNil(t, sample)
	ratio := sample.Ratio
	require.NotNil(t, ratio.Denominator)
	// numerator "1" at offset 23, denominator "2" at offset 25 — the spans
	// must not overlap (the denominator used to inherit the numerator's pos)
	require.Equal(t, Pos(23), ratio.Numerator.Pos())
	require.Equal(t, Pos(24), ratio.Numerator.End())
	require.Equal(t, Pos(25), ratio.Denominator.Pos())
	require.Equal(t, Pos(26), ratio.Denominator.End())
}

func TestGroupByClauseEnd(t *testing.T) {
	sql := "SELECT a FROM t GROUP BY a WITH TOTALS"
	stmt := parseOneStmt(t, sql).(*SelectQuery)
	require.NotNil(t, stmt.GroupBy)
	// the clause ends at the TOTALS keyword, not at the next token's start
	require.Equal(t, Pos(len(sql)), stmt.GroupBy.End())

	sql = "SELECT a FROM t GROUP BY a"
	stmt = parseOneStmt(t, sql).(*SelectQuery)
	require.Equal(t, Pos(len(sql)), stmt.GroupBy.End())

	sql = "SELECT a FROM t GROUP BY ALL"
	stmt = parseOneStmt(t, sql).(*SelectQuery)
	require.Equal(t, Pos(len(sql)), stmt.GroupBy.End())
}

func TestIsNullExprPositions(t *testing.T) {
	sql := "SELECT a IS NULL"
	stmt := parseOneStmt(t, sql).(*SelectQuery)
	isNull := stmt.SelectItems[0].Expr.(*IsNullExpr)
	// the node spans `a IS NULL`: from the operand to the end of NULL
	require.Equal(t, Pos(7), isNull.Pos())
	require.Equal(t, Pos(len(sql)), isNull.End())
	require.Equal(t, Pos(9), isNull.IsPos) // position of IS

	sql = "SELECT a IS NOT NULL"
	stmt = parseOneStmt(t, sql).(*SelectQuery)
	isNotNull := stmt.SelectItems[0].Expr.(*IsNotNullExpr)
	require.Equal(t, Pos(7), isNotNull.Pos())
	require.Equal(t, Pos(len(sql)), isNotNull.End())
	require.Equal(t, Pos(9), isNotNull.IsPos)
}

func TestAlterDetachPartitionPos(t *testing.T) {
	sql := "ALTER TABLE t DETACH PARTITION p"
	stmt := parseOneStmt(t, sql).(*AlterTable)
	detach := stmt.AlterExprs[0].(*AlterTableDetachPartition)
	// the clause starts at the DETACH keyword, like every sibling clause
	require.Equal(t, Pos(14), detach.Pos())
}

func TestAlterDropPartitionEndIncludesSettings(t *testing.T) {
	sql := "ALTER TABLE t DROP PARTITION p SETTINGS mutations_sync=1"
	stmt := parseOneStmt(t, sql).(*AlterTable)
	drop := stmt.AlterExprs[0].(*AlterTableDropPartition)
	require.NotNil(t, drop.Settings)
	// End() used to discard the Settings end and stop at the partition
	require.Equal(t, drop.Settings.End(), drop.End())
	require.Greater(t, drop.End(), drop.Partition.End())
}

func TestDictionaryAttributeEnd(t *testing.T) {
	sql := "CREATE DICTIONARY d (user_id UInt64 IS_OBJECT_ID) PRIMARY KEY user_id SOURCE(CLICKHOUSE()) LAYOUT(FLAT()) LIFETIME(300)"
	stmt := parseOneStmt(t, sql).(*CreateDictionary)
	attrs := stmt.Schema.Attributes
	require.Len(t, attrs, 1)
	// the attribute ends at IS_OBJECT_ID; the old implementation returned
	// NamePos + len("IS_OBJECT_ID"), landing in the middle of the type
	require.Equal(t, "IS_OBJECT_ID", sql[36:48])
	require.Equal(t, Pos(48), attrs[0].End())
}
