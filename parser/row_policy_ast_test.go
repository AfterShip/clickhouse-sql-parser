package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRowPolicyTargetBindings(t *testing.T) {
	cases := []struct {
		query string
		want  [][3]string
	}{
		{"CREATE POLICY p1, p2 ON db.t", [][3]string{{"p1", "db", "t"}, {"p2", "db", "t"}}},
		{"CREATE POLICY p ON db.t, other.t2", [][3]string{{"p", "db", "t"}, {"p", "other", "t2"}}},
		{"CREATE POLICY p1 ON db.t, p2 ON other.t2", [][3]string{{"p1", "db", "t"}, {"p2", "other", "t2"}}},
		{"CREATE POLICY p ON t, t2", [][3]string{{"p", "", "t"}, {"p", "", "t2"}}},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			statements, err := NewParser(tc.query).ParseStmts()
			require.NoError(t, err)
			create, ok := statements[0].(*CreateRowPolicy)
			require.True(t, ok, "got %T", statements[0])
			var got [][3]string
			for _, target := range create.Targets {
				database := ""
				if target.Table.Database != nil {
					database = target.Table.Database.Name
				}
				got = append(got, [3]string{target.PolicyName.Name, database, target.Table.Table.Name})
			}
			require.Equal(t, tc.want, got)
			create.Targets[0].PolicyName.Name = "changed"
			require.Equal(t, tc.want[1][0], create.Targets[1].PolicyName.Name)
			create.Targets[0].Table.Table.Name = "changed"
			require.Equal(t, tc.want[1][2], create.Targets[1].Table.Table.Name)
		})
	}
}

func TestAlterRowPolicyPresence(t *testing.T) {
	parse := func(query string) *AlterRowPolicy {
		t.Helper()
		statements, err := NewParser(query).ParseStmts()
		require.NoError(t, err)
		alter, ok := statements[0].(*AlterRowPolicy)
		require.True(t, ok, "got %T", statements[0])
		return alter
	}
	// Omitted ALTER clauses must not reset the existing mode, filter or recipients.
	omitted := parse("ALTER ROW POLICY p ON db.t")
	require.Empty(t, omitted.Mode)
	require.Nil(t, omitted.To)
	require.Empty(t, omitted.Filters)

	// Explicit NONE must remain distinct from omission and a false expression.
	cleared := parse("ALTER ROW POLICY p ON db.t AS PERMISSIVE USING NONE TO NONE")
	require.Equal(t, RowPolicyPermissive, cleared.Mode)
	require.NotNil(t, cleared.To)
	require.False(t, cleared.To.All)
	require.Empty(t, cleared.To.Principals)
	require.Len(t, cleared.Filters, 1)
	require.True(t, cleared.Filters[0].Using.None)
	require.Nil(t, cleared.Filters[0].Using.Expr)
	require.Nil(t, cleared.Filters[0].WithCheck)

	zero := parse("ALTER ROW POLICY p ON db.t USING 0")
	require.False(t, zero.Filters[0].Using.None)
	number, ok := zero.Filters[0].Using.Expr.(*NumberLiteral)
	require.True(t, ok, "got %T", zero.Filters[0].Using.Expr)
	require.Equal(t, "0", number.Literal)

	checkOnly := parse("ALTER ROW POLICY p ON db.t WITH CHECK id = 42")
	require.Nil(t, checkOnly.Filters[0].Using)
	require.NotNil(t, checkOnly.Filters[0].WithCheck.Expr)
	require.Contains(t, Format(checkOnly), "WITH CHECK id = 42")

	repeated := parse("ALTER ROW POLICY p ON db.t USING 1, USING 0")
	require.Len(t, repeated.Filters, 2)
	require.Equal(t, "1", Format(repeated.Filters[0].Using))
	require.Equal(t, "0", Format(repeated.Filters[1].Using))
}

func TestRowPolicyPrincipalSets(t *testing.T) {
	cases := []struct {
		clause   string
		all      bool
		names    []string
		excluded []string
	}{
		{"TO NONE", false, nil, nil},
		{"TO ALL", true, nil, nil},
		{"TO ALL, ignored EXCEPT r1, CURRENT_USER()", true, nil, []string{"r1", "CURRENT_USER"}},
		{"TO r1, r2 EXCEPT r2", false, []string{"r1", "r2"}, []string{"r2"}},
		{"TO `ALL`, `CURRENT_USER`, currentUser()", false, []string{"`ALL`", "`CURRENT_USER`", "CURRENT_USER"}, nil},
		{"TO 'user`name@host', john@localhost", false, []string{"'user`name@host'", "john@localhost"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.clause, func(t *testing.T) {
			statements, err := NewParser("CREATE ROW POLICY p ON db.t USING 1 " + tc.clause).ParseStmts()
			require.NoError(t, err)
			create, ok := statements[0].(*CreateRowPolicy)
			require.True(t, ok, "got %T", statements[0])
			set := create.To
			require.Equal(t, tc.all, set.All)
			var names, excluded []string
			for _, p := range set.Principals {
				names = append(names, Format(p))
			}
			for _, p := range set.Except {
				excluded = append(excluded, Format(p))
			}
			require.Equal(t, tc.names, names)
			require.Equal(t, tc.excluded, excluded)
			formatted := Format(statements[0])
			again, err := NewParser(formatted).ParseStmts()
			require.NoError(t, err)
			require.Equal(t, formatted, Format(again[0]))
		})
	}
}

func TestRowPolicyWildcardIdentity(t *testing.T) {
	statements, err := NewParser("CREATE POLICY p ON db.*, db.`*`, *").ParseStmts()
	require.NoError(t, err)
	create, ok := statements[0].(*CreateRowPolicy)
	require.True(t, ok, "got %T", statements[0])
	targets := create.Targets
	require.Len(t, targets, 3)
	require.Equal(t, "*", targets[0].Table.Table.Name)
	require.NotEqual(t, BackTicks, targets[0].Table.Table.QuoteType)
	require.Equal(t, BackTicks, targets[1].Table.Table.QuoteType)
	require.Nil(t, targets[2].Table.Database, "unqualified * must retain the implicit current database")
	require.Contains(t, Format(statements[0]), "db.`*`")
}

func TestRowPolicyGlobalInspectionScope(t *testing.T) {
	for _, prefix := range []string{"SHOW ROW POLICIES", "SHOW CREATE ROW POLICIES"} {
		statements, err := NewParser(prefix + " ON *.*").ParseStmts()
		require.NoError(t, err)
		show, ok := statements[0].(*ShowRowPolicies)
		require.True(t, ok, "got %T", statements[0])
		require.Equal(t, "*", show.Table.Database.Name)
		require.Equal(t, "*", show.Table.Table.Name)
	}
	for _, command := range []string{"CREATE", "ALTER", "DROP"} {
		_, err := NewParser(command + " ROW POLICY p ON *.*").ParseStmts()
		require.Error(t, err, command)
	}
	_, err := NewParser("SHOW ROW POLICIES ON *.events").ParseStmts()
	require.Error(t, err)
}

func TestRowPolicyTraversal(t *testing.T) {
	queries := []string{
		"CREATE POLICY p ON db.t IN local_directory ON CLUSTER c AS RESTRICTIVE USING id = 42 WITH CHECK id > 0 TO reader, CURRENT_USER EXCEPT user1@'host'",
		"ALTER POLICY p ON db.t RENAME TO renamed USING NONE TO NONE",
		"DROP POLICY p ON db.t FROM local_directory ON CLUSTER c",
		"SHOW CREATE POLICY p1 ON db.t, p2 ON other.t",
		"SHOW CREATE POLICY p",
		"SHOW POLICIES ON db.*",
	}
	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			statements, err := NewParser(query).ParseStmts()
			require.NoError(t, err)
			var walked []string
			Walk(statements[0], func(expr Expr) bool {
				walked = append(walked, fmt.Sprintf("%T", expr))
				return true
			})
			visitor := &rowPolicyTestVisitor{}
			require.NoError(t, statements[0].Accept(visitor))
			require.Equal(t, walked, visitor.entered)
			require.Zero(t, visitor.depth)
		})
	}
	statements, err := NewParser(queries[0]).ParseStmts()
	require.NoError(t, err)
	stop := errors.New("stop on condition")
	visitor := &rowPolicyTestVisitor{DefaultASTVisitor: DefaultASTVisitor{Visit: func(expr Expr) error {
		if _, ok := expr.(*RowPolicyCondition); ok {
			return stop
		}
		return nil
	}}}
	require.ErrorIs(t, statements[0].Accept(visitor), stop)
	require.Zero(t, visitor.depth)
	conditions := 0
	completed := Walk(statements[0], func(expr Expr) bool {
		if _, ok := expr.(*RowPolicyFilter); ok {
			return false
		}
		if _, ok := expr.(*RowPolicyCondition); ok {
			conditions++
		}
		return true
	})
	require.False(t, completed)
	require.Zero(t, conditions, "stopping at a filter must prevent visits to its conditions")
}

func TestRowPolicyStatementPositions(t *testing.T) {
	queries := []string{
		"CREATE POLICY p ON db.`table`",
		"CREATE POLICY p ON db.t AS PERMISSIVE",
		"CREATE POLICY p ON db.t USING t.`column`",
		"CREATE POLICY p ON db.t TO user1@'host'",
		"ALTER POLICY p ON db.t RENAME TO `new-name`",
		"DROP POLICY p ON db.t ON CLUSTER `my-cluster`",
		"SHOW CREATE POLICY p ON db.`table`",
		"CREATE POLICY `读策略` ON `数据库`.`表` TO `读者`",
	}
	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			statements, err := NewParser(query + "; -- trailing comment").ParseStmts()
			require.NoError(t, err)
			require.Equal(t, Pos(0), statements[0].Pos())
			require.Equal(t, Pos(len(query)), statements[0].End())
		})
	}
}

func TestParser_RowPolicyTruncatedInput(t *testing.T) {
	query := "CREATE ROW POLICY p ON db.t IN local_directory AS RESTRICTIVE USING has([1,2], id) TO ALL EXCEPT r"
	for i := range query {
		require.NotPanics(t, func() { _, _ = NewParser(query[:i]).ParseStmts() })
	}
	for _, suffix := range []string{" /*", " TO 'unterminated", " AS", " USING", " TO", " IN", " ON CLUSTER"} {
		_, err := NewParser("CREATE ROW POLICY p ON db.t" + suffix).ParseStmts()
		require.Error(t, err, suffix)
	}
	for _, query := range []string{
		"CREATE POLICY p ON db.t AS PERMISSIVE /*",
		"CREATE POLICY p ON db.t TO NONE /*",
		"CREATE POLICY p ON db.t USING NONE /*",
		"CREATE POLICY p ON db.t TO CURRENT_USER /*",
		"CREATE POLICY p ON db.t ON CLUSTER 'c' /*",
		"CREATE POLICY p ON db.t AS PERMISSIVE AS RESTRICTIVE",
		"CREATE POLICY p ON db.t IN a IN b",
		"CREATE POLICY p ON db.t TO 'user'@'host'",
		"CREATE POLICY p ON db.t ON CLUSTER a ON CLUSTER b",
	} {
		_, err := NewParser(query).ParseStmts()
		require.Error(t, err, query)
	}
}

func TestParser_RowPolicyContextualKeywords(t *testing.T) {
	for _, word := range []string{"permissive", "restrictive", "policies", "current_user", "currentuser"} {
		for _, query := range []string{
			"SELECT " + word + " FROM db.t",
			"SELECT x AS " + word + " FROM db.t",
			"SELECT x " + word + " FROM db.t",
			"SELECT " + word + "() FROM db.t",
			"CREATE TABLE db." + word + " (" + word + " UInt64) ENGINE = MergeTree ORDER BY " + word,
			"CREATE ROLE " + word,
		} {
			statements, err := NewParser(query).ParseStmts()
			require.NoError(t, err, query)
			require.Len(t, statements, 1)
			again, err := NewParser(Format(statements[0])).ParseStmts()
			require.NoError(t, err, query)
			require.Len(t, again, 1)
			require.Equal(t, Format(statements[0]), Format(again[0]))
		}
	}
}

func FuzzParser_RowPolicyRoundTrip(f *testing.F) {
	for _, seed := range []string{
		"CREATE POLICY p ON db.t USING x = 1 TO ALL EXCEPT r",
		"ALTER POLICY p ON db.t USING NONE TO NONE",
		"DROP POLICY p ON db.t",
		"SHOW CREATE POLICY p ON db.t",
		"SHOW ROW POLICIES ON *.*",
		"CREATE POLICY",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, query string) {
		if len(query) > 4096 || !strings.Contains(strings.ToUpper(query), "POLIC") {
			t.Skip()
		}
		statements, err := NewParser(query).ParseStmts()
		if err != nil {
			return
		}
		for _, statement := range statements {
			formatted := Format(statement)
			again, err := NewParser(formatted).ParseStmts()
			require.NoError(t, err)
			require.Len(t, again, 1)
			require.Equal(t, formatted, Format(again[0]))
		}
	})
}

type rowPolicyTestVisitor struct {
	DefaultASTVisitor
	entered []string
	depth   int
}

func (v *rowPolicyTestVisitor) Enter(expr Expr) {
	v.entered = append(v.entered, fmt.Sprintf("%T", expr))
	v.depth++
}

func (v *rowPolicyTestVisitor) Leave(Expr) { v.depth-- }
