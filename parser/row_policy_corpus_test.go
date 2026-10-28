package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_RowPolicyNativeSyntax(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		query string
		valid bool
	}{
		{name: "create_basic", query: "CREATE ROW POLICY p ON db.t FOR SELECT USING id = 42 TO r", valid: true},
		{name: "create_alias", query: "CREATE POLICY p ON db.t USING id = 42 TO r", valid: true},
		{name: "create_if_missing", query: "CREATE ROW POLICY IF NOT EXISTS p ON db.t USING 1 TO r", valid: true},
		{name: "create_replace", query: "CREATE ROW POLICY OR REPLACE p ON db.t USING 1 TO r", valid: true},
		{name: "no_filter_or_roles", query: "CREATE ROW POLICY p ON db.t", valid: true},
		{name: "filter_none", query: "CREATE ROW POLICY p ON db.t USING NONE", valid: true},
		{name: "filter_false", query: "CREATE ROW POLICY p ON db.t USING 0", valid: true},
		{name: "restrictive_before_filter", query: "CREATE ROW POLICY p ON db.t AS RESTRICTIVE FOR SELECT USING x = 1 TO r", valid: true},
		{name: "restrictive_after_filter", query: "CREATE ROW POLICY p ON db.t USING x = 1 AS RESTRICTIVE TO r", valid: true},
		{name: "permissive", query: "CREATE ROW POLICY p ON db.t AS PERMISSIVE USING 1 TO r", valid: true},
		{name: "for_all", query: "CREATE ROW POLICY p ON db.t FOR ALL USING x = 1 TO r", valid: true},
		{name: "two_names_one_table", query: "CREATE ROW POLICY p1, p2 ON db.t USING 1 TO r", valid: true},
		{name: "one_name_two_tables", query: "CREATE ROW POLICY p ON db.t, db.t_local USING 1 TO r", valid: true},
		{name: "paired_names_tables", query: "CREATE ROW POLICY p1 ON db.t, p2 ON other.t USING 1 TO r", valid: true},
		{name: "no_database", query: "CREATE ROW POLICY p ON t USING 1 TO r", valid: true},
		{name: "database_wildcard", query: "CREATE ROW POLICY p ON db.* USING 1 TO r", valid: true},
		{name: "current_database_wildcard", query: "CREATE ROW POLICY p ON * USING 1 TO r", valid: true},
		{name: "quoted_literal_star", query: "CREATE ROW POLICY p ON db.`*` USING 1 TO r", valid: true},
		{name: "quoted_names", query: "CREATE ROW POLICY 'policy name' ON `my-db`.`my-table` USING `context.id` = 42 TO `role-name`, 'user-name'", valid: true},
		{name: "storage", query: "CREATE ROW POLICY p ON db.t IN local_directory USING 1 TO r", valid: true},
		{name: "cluster_before_target", query: "CREATE ROW POLICY p ON CLUSTER 'cluster-name' ON db.t USING 1 TO r", valid: true},
		{name: "cluster_after_target", query: "CREATE ROW POLICY p ON db.t ON CLUSTER 'cluster-name' USING 1 TO r", valid: true},
		{name: "cluster_at_end", query: "CREATE ROW POLICY p ON db.t USING 1 TO r ON CLUSTER 'cluster-name'", valid: true},
		{name: "to_none", query: "CREATE ROW POLICY p ON db.t USING 1 TO NONE", valid: true},
		{name: "to_all", query: "CREATE ROW POLICY p ON db.t USING 1 TO ALL", valid: true},
		{name: "to_all_except", query: "CREATE ROW POLICY p ON db.t USING 1 TO ALL EXCEPT r1, r2", valid: true},
		{name: "to_named_except", query: "CREATE ROW POLICY p ON db.t USING 1 TO r1, r2 EXCEPT r2", valid: true},
		{name: "to_current_user", query: "CREATE ROW POLICY p ON db.t USING 1 TO CURRENT_USER", valid: true},
		{name: "to_current_user_parens", query: "CREATE ROW POLICY p ON db.t USING 1 TO CURRENT_USER()", valid: true},
		{name: "to_currentuser_function", query: "CREATE ROW POLICY p ON db.t USING 1 TO currentUser()", valid: true},
		{name: "to_user_host", query: "CREATE ROW POLICY p ON db.t USING 1 TO john@'localhost'", valid: true},
		{name: "predicate_function_commas", query: "CREATE ROW POLICY p ON db.t USING has([1, 2, 3], id) AND id = toUInt64(getSetting('filter_id')) TO r", valid: true},
		{name: "predicate_subquery", query: "CREATE ROW POLICY p ON db.t USING id IN (SELECT id FROM db.lookup_table) TO r", valid: true},
		{name: "alter_filter", query: "ALTER ROW POLICY p ON db.t USING id = 99", valid: true},
		{name: "alter_if_exists", query: "ALTER POLICY IF EXISTS p ON db.t USING 0", valid: true},
		{name: "alter_filter_none", query: "ALTER ROW POLICY p ON db.t USING NONE", valid: true},
		{name: "alter_roles_none", query: "ALTER ROW POLICY p ON db.t TO NONE", valid: true},
		{name: "alter_mode", query: "ALTER ROW POLICY p ON db.t AS PERMISSIVE", valid: true},
		{name: "alter_rename", query: "ALTER ROW POLICY p ON db.t RENAME TO 'new-policy'", valid: true},
		{name: "alter_multiple", query: "ALTER ROW POLICY p1, p2 ON db.t TO ALL", valid: true},
		{name: "drop_basic", query: "DROP ROW POLICY p ON db.t", valid: true},
		{name: "drop_if_exists_alias", query: "DROP POLICY IF EXISTS p ON db.t", valid: true},
		{name: "drop_pairs", query: "DROP ROW POLICY p1 ON db.t, p2 ON other.t", valid: true},
		{name: "drop_storage_cluster", query: "DROP ROW POLICY p ON db.t FROM local_directory ON CLUSTER 'c'", valid: true},
		{name: "attach_id_roles", query: "ATTACH ROW POLICY p ON db.t USING 1 TO ID('11111111-1111-1111-1111-111111111111')", valid: false},
		{name: "attach_no_roles", query: "ATTACH ROW POLICY p ON db.t USING 1", valid: false},
		{name: "show_exact", query: "SHOW CREATE ROW POLICY p ON db.t", valid: true},
		{name: "show_alias", query: "SHOW CREATE POLICY p ON db.t", valid: true},
		{name: "show_multiple", query: "SHOW CREATE ROW POLICY p1, p2 ON db.t", valid: true},
		{name: "show_short_name", query: "SHOW CREATE ROW POLICY p", valid: true},
		{name: "show_all", query: "SHOW CREATE ROW POLICIES", valid: true},
		{name: "show_table", query: "SHOW CREATE ROW POLICIES ON db.t", valid: true},
		{name: "show_policies", query: "SHOW ROW POLICIES ON db.*", valid: true},
		{name: "legacy_global_wildcard", query: "CREATE ROW POLICY p ON *.* USING 1 TO r", valid: false},
		{name: "legacy_prefix_wildcard", query: "CREATE ROW POLICY p ON db.prefix* USING 1 TO r", valid: false},
		{name: "legacy_prefix_database", query: "CREATE ROW POLICY p ON db*.* USING 1 TO r", valid: false},
		{name: "with_check_only", query: "CREATE ROW POLICY p ON db.t WITH CHECK x = 1 TO r", valid: true},
		{name: "with_using_and_check", query: "CREATE ROW POLICY p ON db.t USING x = 1 WITH CHECK y = 2 TO r", valid: true},
		{name: "repeated_filter", query: "ALTER ROW POLICY p ON db.t USING x = 1, USING x = 2", valid: true},
		{name: "invalid_cartesian", query: "CREATE ROW POLICY p1, p2 ON db.t, db.t2 USING 1 TO r", valid: false},
		{name: "invalid_multi_name_then_pair", query: "CREATE ROW POLICY p1, p2 ON db.t, p3 ON db.t3 USING 1 TO r", valid: false},
		{name: "invalid_clusters_per_pair", query: "CREATE ROW POLICY p1 ON CLUSTER 'c1' ON db.t, p2 ON CLUSTER 'c2' ON db.t2 USING 1 TO r", valid: false},
		{name: "invalid_insert", query: "CREATE ROW POLICY p ON db.t FOR INSERT USING 1 TO r", valid: false},
		{name: "invalid_update", query: "CREATE ROW POLICY p ON db.t FOR UPDATE USING 1 TO r", valid: false},
		{name: "invalid_delete", query: "CREATE ROW POLICY p ON db.t FOR DELETE USING 1 TO r", valid: false},
		{name: "invalid_multi_rename", query: "ALTER ROW POLICY p1, p2 ON db.t RENAME TO p3", valid: false},
		{name: "invalid_drop_sync", query: "DROP ROW POLICY p ON db.t SYNC", valid: false},
		{name: "invalid_comment", query: "CREATE ROW POLICY p ON db.t USING 1 TO r COMMENT 'owner'", valid: false},
		{name: "invalid_attach_named_role", query: "ATTACH ROW POLICY p ON db.t USING 1 TO r", valid: false},
		{name: "invalid_create_or_replace_order", query: "CREATE OR REPLACE ROW POLICY p ON db.t USING 1 TO r", valid: false},
		{name: "invalid_top_alias", query: "CREATE ROW POLICY p ON db.t USING x AS y TO r", valid: false},
		{name: "invalid_missing_expression", query: "CREATE ROW POLICY p ON db.t USING TO r", valid: false},
		{name: "invalid_trailing_roles_comma", query: "CREATE ROW POLICY p ON db.t USING 1 TO r,", valid: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			statements, err := NewParser(tc.query).ParseStmts()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, statements, 1)
			compact := Format(statements[0])
			again, err := NewParser(compact).ParseStmts()
			require.NoError(t, err)
			require.Len(t, again, 1)
			require.Equal(t, compact, Format(again[0]))
			beautifier := NewFormatter().WithBeautify()
			beautifier.WriteExpr(statements[0])
			pretty, err := NewParser(beautifier.String()).ParseStmts()
			require.NoError(t, err)
			require.Equal(t, compact, Format(pretty[0]))
		})
	}
}
