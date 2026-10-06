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

func TestLowercaseArrayJoinAST(t *testing.T) {
	type join struct {
		modifiers []string
		items     int
	}
	cases := []struct {
		sql   string
		joins []join
	}{
		{
			sql:   "SELECT * FROM t array join spec_containers",
			joins: []join{{modifiers: []string{"array", "JOIN"}, items: 1}},
		},
		{
			sql:   "SELECT * FROM t array join a AS x, b AS y",
			joins: []join{{modifiers: []string{"array", "JOIN"}, items: 2}},
		},
		{
			sql: "SELECT * FROM t left array join a AS x join u ON t.id = u.id",
			joins: []join{
				{modifiers: []string{"left", "array", "JOIN"}, items: 1},
				{modifiers: []string{"JOIN"}},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmts, err := NewParser(tc.sql).ParseStmts()
			require.NoError(t, err)
			require.Len(t, stmts, 1)

			query, ok := stmts[0].(*SelectQuery)
			require.True(t, ok)
			root, ok := query.From.Expr.(*JoinExpr)
			require.True(t, ok)
			require.IsType(t, &JoinTableExpr{}, root.Left)

			next := root.Right
			for _, want := range tc.joins {
				j, ok := next.(*JoinExpr)
				require.True(t, ok, "expected JoinExpr, got %T", next)
				require.Equal(t, want.modifiers, j.Modifiers)
				if want.items > 0 {
					list, ok := j.Left.(*ColumnExprList)
					require.True(t, ok, "expected ColumnExprList, got %T", j.Left)
					require.Len(t, list.Items, want.items)
				} else {
					require.IsType(t, &JoinTableExpr{}, j.Left)
				}
				next = j.Right
			}
			require.Nil(t, next)
		})
	}
}
