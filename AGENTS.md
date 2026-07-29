# AI Agents Guidline

This document outlines guidance for AI coding agent including project structure, coding style, testing, and contribution practices for the ClickHouse SQL Parser project.


## Development Commands

```shell
# Build the CLI binary
make

# Run tests
make test

# Update golden fixtures after intentional output changes
make update_test
```

After editing code, use `goimports` and `gofmt` to maintain code style, and run `make lint` to check for any issues before committing or requesting a review.

## Project Structure & Module Organization
- `main.go` is the CLI entry point (`clickhouse-sql-parser`) for AST output and SQL formatting.
- `parser/` contains core parser code: lexer (`lexer.go`), AST definitions (`ast.go`), traversal helpers (`walk.go`), and grammar-specific parser files (`parser_query.go`, `parser_table.go`, `parser_alter.go`, etc.).
- Tests live next to source as `*_test.go` files, with fixtures under `parser/testdata/`.
- Fixture groups are organized by SQL type (`basic/`, `query/`, `dml/`, `ddl/`), with generated expectations in `output/` (AST JSON) and `format/` (formatted SQL).

## Coding Style & Naming Conventions
- Use Go 1.21 conventions (`go.mod`) and keep code `gofmt`/`goimports` clean (enforced by lint).
- Naming is the most important style aspect, try you best to choose a clear and descriptive name for variables, functions, types, and files. For example, use `parseSelect` for a function that parses a SELECT statement, and `SelectStatement` for the corresponding AST node type.
- Place parsing logic in the matching module by statement family (for example, query parsing in `parser/parser_query.go`).
- Follow existing parser naming patterns such as `parseXxx` helpers and explicit AST type names.
- Keep AST `FormatSQL()` output deterministic; formatting changes must be reflected in golden files.
- You must go through the repository before adding new code to ensure consistency with existing patterns and styles. If you are unsure about where to place new code or how to format it, please refer to the existing codebase or ask for guidance.
- Reusing existing code and patterns is encouraged to maintain consistency and reduce redundancy. If you find a similar function or pattern in the codebase, consider adapting it for your needs instead of creating something new from scratch.

## Validate SQL with `clickhouse-local`

Before adding support for new syntax or fixing a parsing issue, **verify that the SQL is actually valid ClickHouse** by checking it against the real engine with `clickhouse-local`. This prevents the parser from accepting syntax that ClickHouse itself rejects (or rejecting syntax it accepts), which is the single most common source of incorrect grammar changes.

- Install it once if it is not already available:
  - macOS: `brew install clickhouse`
  - Linux/other: `curl https://clickhouse.com/ | sh` (produces a `./clickhouse` binary; invoke as `./clickhouse local`)
  - Verify the install with `clickhouse local --version`.
- Validate a statement without executing it against real data by asking ClickHouse to parse it:

  ```shell
  # Syntax-check only (does not run the query). Exit code 0 means the SQL parses.
  echo "SELECT * FROM system.one WHERE a = 100" | clickhouse local --query "EXPLAIN SYNTAX $(cat)"

  # Or check DDL/DML that references no data:
  clickhouse local --multiquery --query "CREATE TABLE t (id UInt64) ENGINE = Memory; DESCRIBE t;"
  ```

- Workflow:
  1. Confirm the SQL you intend to support parses (or, for a bug fix, that ClickHouse’s behavior matches your expectation) with `clickhouse-local`.
  2. Only then add or change the grammar/AST in this parser.
  3. Add the same SQL as a fixture under `parser/testdata/<category>/` so the behavior is locked in.
- If `clickhouse-local` is unavailable in your environment, state that explicitly in the PR description and cite the relevant ClickHouse documentation for the syntax instead, rather than guessing.

## Testing Guidelines
- Use Go’s `testing` package with `testify/require` assertions and `goldie` snapshot comparisons.
- Add new SQL cases as `.sql` files under the appropriate `parser/testdata/<category>/` directory.
- If expected outputs change, run `make update_test` and commit updated files in both `output/` and/or `format/`.
- Prefer descriptive test names (`TestParser_*`, `TestWalk_*`) and subtests for per-fixture coverage.

## Commit & Pull Request Guidelines
- Match existing commit style: concise, imperative subjects like `Add support for ...` or `Fix parsing failure ...`, optionally with issue refs (for example `(#235)`).
- Keep PRs focused; describe grammar/AST impact, include representative SQL examples, and note regenerated fixtures.
- Before requesting review, run `make lint` and `make test` locally to mirror CI expectations.


## Important Notes

- Always validate the target SQL against `clickhouse-local` before adding new syntax or fixing a parsing issue (see "Validate SQL with `clickhouse-local`" above).
- You must confirm it's correctly added to `visitor.go`, `walk.go` and `format.go` when adding a new expression or statement type. This ensures that the new AST node is properly traversed and formatted.
- Newly added test cases must be concise and cover the core functionality being tested first.
