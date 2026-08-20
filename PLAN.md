# Plan: Close ClickHouse parser grammar gaps

Status: **implemented and verified.** Both gaps below are fixed on this
branch; see each gap's “Implemented” section for the change and the
“Verification” section at the end for the results.

## Why

The goal is a pure-Go, dialect-faithful parser that fully covers the
statement shapes below. The gaps were found by running the parser against a
corpus of real migration statements; each item is a statement shape this
parser rejected, and fixing them removes the need for any keyword-prefix
fallback when classifying statements.

## Validation method

Each candidate statement was checked with the full ClickHouse server via
Docker (`clickhouse/clickhouse-server:latest`, version 26.7.4.58), not
`clickhouse local`: `clickhouse local` uses a restricted parser and wrongly
rejects statements the real engine accepts (it rejected
`CREATE OR REPLACE MATERIALIZED VIEW`, which the full server parses).
Candidates were run as real DDL against the server (with prerequisite
tables created), so a pass means the engine actually accepts the syntax.

## Gap 1 — `CREATE OR REPLACE MATERIALIZED VIEW`

**Statement:** `CREATE OR REPLACE MATERIALIZED VIEW mv TO dest AS SELECT * FROM src`

**Validated:** accepted by ClickHouse server 26.7.4.58 (DDL executed
successfully).

**Implemented:**
- `parser/parser_table.go` — `parseDDL` guard now accepts `MATERIALIZED`
  after `CREATE OR REPLACE` (error message updated to
  `TEMPORARY|TABLE|VIEW|FUNCTION|DICTIONARY|MATERIALIZED`).
- `parser/ast.go` — `CreateMaterializedView` gained an `OrReplace` field
  (mirroring `CreateView.OrReplace`).
- `parser/parser_view.go` — `parseCreateMaterializedView` accepts the
  `orReplace` flag and sets it on the node.
- `parser/format.go` — `FormatSQL` renders `CREATE OR REPLACE MATERIALIZED
  VIEW` when the flag is set.
- Existing materialized-view fixtures' goldens regenerated with the new
  field.

## Gap 2 — TTL with `GROUP BY` action (TTL delete-by-group)

**Statement:**
`CREATE TABLE t (id UInt64, created DateTime, x UInt64) ENGINE = MergeTree() ORDER BY id TTL created + INTERVAL 1 DAY GROUP BY id SET x = sum(x)`

**Validated:** accepted by ClickHouse server 26.7.4.58. Engine constraints
confirmed: the `SET` assignment must contain an aggregate function
(`SET x = 1` fails with `BAD_TTL_EXPRESSION`), `SET` is optional after
`GROUP BY` (`GROUP BY id` alone creates the table), and `GROUP BY` keys must
be a prefix of the primary key.

**Implemented:**
- `parser/ast.go` — new nodes `TTLPolicyGroupBy` (expr list + `Set`
  assignments) and `TTLPolicySetExpr` on `TTLPolicyRule`, with
  `Pos`/`End`/`Accept`.
- `parser/parser_table.go` — `tryParseTTLPolicy` handles a leading `GROUP
  BY` without an action keyword; new `parseTTLPolicyGroupBy` parses the key
  list (comma-separated, kept in a `ColumnExprList`) and the optional
  `SET <col> = <agg expr>[, ...]` assignments.
- `parser/parser_common.go` — `parseIdentOrKeyword` helper for assignment
  target names (which may be keywords).
- `parser/format.go` — `FormatSQL` for `TTLPolicyGroupBy`/`TTLPolicySetExpr`
  (one `SET`, comma-joined assignments).
- `parser/ast_visitor.go` / `parser/walk.go` — visitor methods and walk
  cases for both new node types.

## Test coverage

Fixtures added under `parser/testdata/ddl/` with regenerated goldens in
`output/`, `format/`, and `format/beautify/`:

- `create_or_replace_materialized_view.sql`
- `create_table_ttl_group_by.sql` — three engine-validated variants: single
  `SET`, `SET`-less `GROUP BY id`, and multi-key/multi-`SET`
  (`GROUP BY id, created SET x = sum(x), total = count()`).

## Verification

- `make test` — full suite passes
  (`ok github.com/AfterShip/clickhouse-sql-parser/parser`, 58.5% coverage).
- `make update_test` — goldens regenerated and committed to the working
  tree; existing TTL/MV goldens updated only by the new fields.
- `gofmt -l` — clean. (`golangci-lint` is not installed in this
  environment, so `make lint` cannot run; gofmt and a `Walk` smoke check
  were used instead.)
- `Walk` smoke test on the new TTL AST reaches all new nodes
  (`TTLPolicyRule: 1, TTLPolicyGroupBy: 1, TTLPolicySetExpr: 2`).
- CLI round-trip: both statements parse and re-format byte-identically.
- Engine parity: every accepted form above (single/multi key, with/without
  `SET`, single/multi assignment) was executed against
  `clickhouse/clickhouse-server:latest` and matches the parser's accepted
  input.