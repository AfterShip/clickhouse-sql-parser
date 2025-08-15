# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a ClickHouse SQL parser written in Go that parses ClickHouse SQL into AST (Abstract Syntax Tree) and provides SQL formatting capabilities. The project is inspired by memefish and is designed to work both as a Go library and a CLI tool.

## Build and Development Commands

### Build the CLI tool
```bash
make
# or
go build -o clickhouse-sql-parser main.go
```

### Run tests
```bash
make test
# Runs tests with coverage, race detection, and compatible flag
```

### Run compatible tests (for ClickHouse compatibility)
```bash
make test -compatible
# Tests against real ClickHouse SQL files from testdata/query/compatible/
```

### Update test golden files
```bash
make update_test
# Updates expected output files in testdata/*/output/ directories
```

### Run linting
```bash
make lint
# Uses golangci-lint with 20 minute timeout
```

### Run benchmarks
```bash
go test -bench=. -benchmem ./parser
```

## Architecture Overview

### Core Components

**Lexer (`parser/lexer.go`)**
- Tokenizes ClickHouse SQL input into tokens
- Handles keywords, identifiers, operators, literals, and comments
- Supports various token types including strings, numbers, and punctuation

**Parser (`parser/parser_*.go`)**
- Modular parser split across multiple files by functionality:
  - `parser_common.go` - Core parser logic and utilities
  - `parser_query.go` - SELECT statements and query parsing
  - `parser_table.go` - CREATE TABLE and table-related DDL
  - `parser_alter.go` - ALTER statements
  - `parser_drop.go` - DROP statements
  - `parser_view.go` - View-related statements
  - `parser_column.go` - Column definitions and operations

**AST (`parser/ast.go`)**
- Defines all AST node types implementing the `Expr` interface
- Each node provides `Pos()`, `End()`, `String()`, and `Accept()` methods
- Supports visitor pattern for AST traversal

**AST Traversal**
- **Walk Pattern** (`parser/walk.go`) - Recommended approach for AST traversal
  - `Walk(node, fn)` - Depth-first traversal
  - `Find(root, predicate)` - Find first matching node
  - `FindAll(root, predicate)` - Find all matching nodes
  - `Transform(root, transformer)` - Apply transformations
- **Visitor Pattern** (`parser/ast_visitor.go`) - More complex but powerful traversal

**Main Entry Point (`main.go`)**
- CLI tool supporting parsing to AST JSON or formatting SQL
- Accepts input from command line arguments or files

### Key Interfaces

- `Expr` - Base interface for all AST nodes
- `DDL` - Interface for Data Definition Language statements
- `ASTVisitor` - Visitor pattern interface for AST traversal
- `WalkFunc` - Function type for Walk pattern traversal

## Testing Strategy

The project uses a comprehensive testing approach:

**Golden File Testing**
- Test cases in `parser/testdata/` organized by category:
  - `basic/` - Simple test cases
  - `ddl/` - Data Definition Language tests
  - `dml/` - Data Manipulation Language tests  
  - `query/` - SELECT and query tests
- Expected outputs stored in `output/` subdirectories as `.golden.json` files
- Formatted SQL outputs in `format/` subdirectories

**Compatible Testing**
- Real ClickHouse SQL files in `testdata/query/compatible/1_stateful/`
- Run with `-compatible` flag to test against actual ClickHouse queries

**Benchmark Testing**
- Performance tests in `parser/benchmark_test.go`
- Tests parsing speed and memory allocation for various query types

## Development Guidelines

**Adding New SQL Features**
1. Add test cases to appropriate `testdata/` subdirectory
2. Implement lexer tokens if needed in `lexer.go`
3. Add AST node types to `ast.go` with all required methods
4. Implement parsing logic in appropriate `parser_*.go` file
5. Add visitor methods to `ast_visitor.go` if using visitor pattern
6. Update Walk functions in `walk.go` for new node types
7. Run `make update_test` to generate golden files

**Parser Module Organization**
- Keep parser functions organized by SQL statement type
- Use consistent naming: `parseXXX()` for parsing functions
- Implement proper error handling with descriptive messages
- Follow existing patterns for operator precedence and expression parsing

**AST Node Implementation**
- All nodes must implement `Pos()`, `End()`, `String()`, and `Accept()` methods
- String() method should regenerate valid ClickHouse SQL
- Accept() method must call visitor.Enter()/Leave() and visit all child nodes

** Walking the AST**

- For a new expression type, it should be also added to the `Walk` function in `walk.go`.
