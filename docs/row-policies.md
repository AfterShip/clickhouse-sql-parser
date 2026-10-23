# Row policies

The parser supports `CREATE`, `ALTER`, and `DROP [ROW] POLICY`, plus `SHOW CREATE [ROW] POLICY` and `SHOW [ROW] POLICIES`. These produce typed AST nodes rather than opaque SQL strings. Parsing does not execute policies, resolve users or roles, or check whether a table exists.

```sql
CREATE ROW POLICY example_policy
ON db.example_table, db.example_table_local
AS PERMISSIVE
FOR SELECT USING id > 0
TO example_role;

ALTER ROW POLICY example_policy ON db.example_table
USING NONE TO NONE;

DROP ROW POLICY IF EXISTS example_policy
ON db.example_table, db.example_table_local;

SHOW CREATE ROW POLICY example_policy ON db.example_table;
SHOW ROW POLICIES ON db.*;
```

## Target bindings

`Targets` contains explicit `RowPolicyTarget` nodes. Each binds one `PolicyName` to a `TableIdentifier`.

The following forms are supported:

```sql
CREATE POLICY p1, p2 ON db.t USING 1 TO reader;
CREATE POLICY p ON db.t, other.t USING 1 TO reader;
CREATE POLICY p1 ON db.t, p2 ON other.t USING 1 TO reader;
```

A multi-name list cannot also have a multi-table target list. `p1, p2 ON t1, t2` is not a Cartesian product and is rejected.

Database policies use `db.*`, or `*` for the current database. A quoted table named `*` remains distinct from a wildcard. Policy declarations reject prefix wildcards and wildcard database names rather than silently narrowing them. `SHOW [CREATE] ROW POLICIES ON *.*` can inspect all policies.

`ON CLUSTER` is a statement-level clause. CREATE/ALTER can select an access storage with `IN`; DROP uses `FROM`.

## ALTER presence semantics

`CreateRowPolicy` and `AlterRowPolicy` embed `RowPolicyOptions`. The options retain omitted fields because omission on ALTER means keeping the existing value.

- `Mode == ""` means omitted. `RowPolicyPermissive` and `RowPolicyRestrictive` are explicit modes.
- No `Filters` means no filter update. Each `RowPolicyFilter` has optional `Using` and `WithCheck` conditions. A condition with `None == true` clears that condition; an expression containing zero is still an expression.
- `To == nil` means omitted. A non-nil `RowPolicySubjects` with no principals and `All == false` represents `TO NONE`.
- `RenameTo` is only available on ALTER and requires exactly one target.

CREATE's `OR REPLACE` is different from ALTER. Consumers must not interpret omitted CREATE fields as instructions to preserve the previous policy.

Repeated filter clauses remain ordered in `Filters`. They are not collapsed into a single expression.

## Principal sets

`RowPolicySubjects` represents `TO`, including named users/roles, `ALL`, `NONE`, `EXCEPT`, and `CURRENT_USER`. Dynamic current-user aliases normalize to `CURRENT_USER`; quoted principals with the same spelling remain names.

A `RowPolicyPrincipal` contains either `Name` and an optional `Host`, or `CurrentUser == true`. A complete single-quoted principal name can contain `@`; a separate host suffix follows an identifier, as in `user1@'localhost'` or a backtick-quoted identifier followed by `@'localhost'`.

## Filtering scope and traversal

ClickHouse currently implements SELECT row filtering. `FOR ALL` normalizes to SELECT. `FOR INSERT`, `FOR UPDATE`, and `FOR DELETE` are rejected.

The parser retains `WITH CHECK` syntax as a separate condition when present. That does not imply that ClickHouse enforces write-row checks. Consumers should validate the enforcement capabilities of their target server.

Filter expressions are visible through both `Accept`/`ASTVisitor` and `Walk`. These APIs inspect syntax; consumers must not evaluate request-dependent functions such as `getSetting` or `currentUser` while rewriting policy DDL.

Custom `ASTVisitor` implementations need the new policy visit methods. Embedding `DefaultASTVisitor` supplies the default methods.

Internal `ATTACH ROW POLICY` access-storage serialization is not part of the supported SQL statement grammar.
