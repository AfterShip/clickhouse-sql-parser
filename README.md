# ClickHouse SQL Parser 

The goal of this project is to build a ClickHouse SQL parser in Go with the following key features:

- Parse ClickHouse SQL into AST
- Beautify ClickHouse SQL format

This project is inspired by [memefish](https://github.com/cloudspannerecosystem/memefish) which is a SQL parser for Spanner in Go.
## How to use

You can use it as your Go library or CLI tool, see the following examples:

- Use clickhouse-sql-parser as a Go library

```Go
package main

import (
    clickhouse "github.com/AfterShip/clickhouse-sql-parser/parser"
)

query := "SELECT * FROM clickhouse"
parser := clickhouse.NewParser(query)
// Parse query into AST
statements, err := newParser.ParseStatements()
if err != nil {
    return nil, err
}

// To beautify query, it's as simple as:
for _, statement := range statements {
    fmt.Println(statement.String(0))
}
```

- Use clickhouse-sql-parser as a CLI tool

```bash
$ go install github.com/AfterShip/clickhouse-sql-parser@latest
## Parse query into AST
$ clickhouse-sql-parser "SELECT * FROM clickhouse WHERE a=100"

## Beautify query
$ clickhouse-sql-parser -format "SELECT * FROM clickhouse WHERE a=100"

## Parse query from file
$ clickhouse-sql-parser -file ./test.sql
```

## Contact us

Feel free to open an issue or discussion if you have any issues or questions.
