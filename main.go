package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	clickhouse "github.com/AfterShip/clickhouse-sql-parser/parser"
)

// Version is set at build time using -ldflags:
// go build -ldflags "-X main.Version=v0.4.17"
var Version = "dev"
const help = `
Usage: clickhouse-sql-parser [YOUR SQL STRING] -f [YOUR SQL FILE] -format
`

var options struct {
	help    bool
	file    string
	format  bool
	version bool
}

func init() {
	flag.BoolVar(&options.format, "format", false, "Beautify print the ClickHouse SQL")
	flag.StringVar(&options.file, "f", "", "Parse SQL from file")
	flag.BoolVar(&options.help, "h", false, "Print help message")
	flag.BoolVar(&options.version, "v", false, "Print version")
}

func main() {
	flag.Parse()
	if options.version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if len(os.Args) < 2 || options.help {
		fmt.Print(help)
		os.Exit(0)
	}

	var err error
	var inputBytes []byte
	if options.file != "" {
		inputBytes, err = os.ReadFile(options.file)
		if err != nil {
			panic(fmt.Sprintf("read file error: %s", err.Error()))
		}
	} else {
		if strings.HasPrefix(os.Args[len(os.Args)-1], "-") {
			fmt.Print(help)
			os.Exit(0)
		}
		inputBytes = []byte(os.Args[len(os.Args)-1])
	}
	parser := clickhouse.NewParser(string(inputBytes))
	stmts, err := parser.ParseStmts()
	if err != nil {
		fmt.Printf("parse statements error: %s\n", err.Error())
		os.Exit(1)
	}
	if !options.format { // print AST
		bytes, _ := json.MarshalIndent(stmts, "", "  ") // nolint
		fmt.Println(string(bytes))
	} else { // format SQL
		for _, stmt := range stmts {
			fmt.Println(stmt.String())
		}
	}
}
