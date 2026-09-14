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

func TestAlterTableEndIncludesSettings(t *testing.T) {
	sql := "ALTER TABLE t DROP PARTITION p SETTINGS mutations_sync=1"
	stmt := parseOneStmt(t, sql).(*AlterTable)
	require.NotNil(t, stmt.Settings)
	// End() used to discard the Settings end and stop at the last alter clause
	require.Equal(t, stmt.Settings.End(), stmt.End())
	require.Greater(t, stmt.End(), stmt.AlterExprs[0].End())
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

func TestKeywordArgFunctionPositions(t *testing.T) {
	// The keyword forms reuse UnaryExpr/BinaryOperation, so their spans must match
	// what the same nodes report anywhere else.
	sql := "SELECT trim(BOTH ' ' FROM s)"
	stmt := parseOneStmt(t, sql).(*SelectQuery)
	fn := stmt.SelectItems[0].Expr.(*FunctionExpr)
	require.Equal(t, Pos(7), fn.Pos())
	// Spans are complete half-open source ranges, so the closing paren is included.
	require.Equal(t, Pos(len(sql)), fn.End())

	from := fn.Params.Items.Items[0].(*ColumnExpr).Expr.(*BinaryOperation)
	require.Equal(t, "BOTH", sql[12:16])
	require.Equal(t, "s", sql[26:27])
	// the operation spans its operands: BOTH ... s
	require.Equal(t, Pos(12), from.Pos())
	require.Equal(t, Pos(27), from.End())

	// the characters literal runs from inside the opening quote to the closing one
	modifier := from.LeftExpr.(*UnaryExpr)
	require.Equal(t, Pos(12), modifier.Pos())
	require.Equal(t, modifier.Expr.End(), modifier.End())
	require.Equal(t, Pos(17), modifier.Expr.Pos())
	require.Equal(t, Pos(20), modifier.Expr.End())
}

func TestJoinLocalityIsInsideTheJoinSpan(t *testing.T) {
	sql := "SELECT * FROM t1 GLOBAL LEFT JOIN t2 ON t1.a = t2.a"
	stmt := parseOneStmt(t, sql).(*SelectQuery)
	join := stmt.From.Expr.(*JoinExpr).Right.(*JoinExpr)
	// the span starts at GLOBAL, so a source-rewriting consumer keeps the locality
	require.Equal(t, "GLOBAL", sql[17:23])
	require.Equal(t, Pos(17), join.Pos())
	require.Equal(t, Pos(len(sql)), join.End())
}

func TestSourceSpansCoverReportedNodes(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		wantSlice string
		get       func(Expr) Expr
	}{
		{"SELECT *", "*", func(e Expr) Expr { return e.(*SelectQuery).SelectItems[0].Expr }},
		{"SELECT ?", "?", func(e Expr) Expr { return e.(*SelectQuery).SelectItems[0].Expr }},
		{"SELECT 'abc'", "'abc'", func(e Expr) Expr { return e.(*SelectQuery).SelectItems[0].Expr }},
		{"SELECT f(1)", "f(1)", func(e Expr) Expr { return e.(*SelectQuery).SelectItems[0].Expr }},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneStmt(t, tc.sql)
			node := tc.get(stmt)
			require.GreaterOrEqual(t, node.Pos(), Pos(0))
			require.LessOrEqual(t, node.Pos(), node.End())
			require.LessOrEqual(t, node.End(), Pos(len(tc.sql)))
			require.Equal(t, tc.wantSlice, tc.sql[node.Pos():node.End()])
		})
	}
}

func TestSelectQuerySpanContainsSetOperations(t *testing.T) {
	sql := "SELECT 1 INTERSECT SELECT 2 UNION ALL SELECT 3"
	root := parseOneStmt(t, sql).(*SelectQuery)
	require.Equal(t, Pos(len(sql)), root.End())
	require.NotNil(t, root.Intersect)
	require.NotNil(t, root.Intersect.UnionAll)
	require.LessOrEqual(t, root.Pos(), root.Intersect.Pos())
	require.LessOrEqual(t, root.Intersect.End(), root.End())
	require.LessOrEqual(t, root.Intersect.UnionAll.End(), root.End())
}
