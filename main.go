package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	clickhouse "github.com/AfterShip/clickhouse-sql-parser/parser"
)

// version is overridable at build time via
// -ldflags "-X main.version=..." (see the Makefile); when empty the
// binary reports the module version the Go toolchain recorded, so a
// `go install ...@vX.Y.Z` build needs no hand-maintained constant.
var version string

const help = `
Usage: clickhouse-sql-parser [YOUR SQL STRING] -f [YOUR SQL FILE] -format -beautify
`

func getVersion() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	// A build from a source checkout has no module version; report the
	// VCS revision the toolchain stamped instead.
	var revision, dirty string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			if setting.Value == "true" {
				dirty = "-dirty"
			}
		}
	}
	if revision == "" {
		return "devel"
	}
	if len(revision) > 12 {
		revision = revision[:12]
	}
	return "devel-" + revision + dirty
}

var options struct {
	help     bool
	file     string
	format   bool
	beautify bool
	version  bool
}

func init() {
	flag.BoolVar(&options.format, "format", false, "Print formatted ClickHouse SQL")
	flag.BoolVar(&options.beautify, "beautify", false, "Beautify print the ClickHouse SQL")
	flag.StringVar(&options.file, "f", "", "Parse SQL from file")
	flag.BoolVar(&options.help, "h", false, "Print help message")
	flag.BoolVar(&options.version, "v", false, "Print version")
}

func main() {
	flag.Parse()
	if options.version {
		fmt.Println(getVersion())
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
			fmt.Fprintf(os.Stderr, "read file error: %s\n", err.Error())
			os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "parse statements error: %s\n", err.Error())
		os.Exit(1)
	}
	if !options.format && !options.beautify { // print AST
		bytes, _ := json.MarshalIndent(stmts, "", "  ") // nolint
		fmt.Println(string(bytes))
	} else { // format SQL
		for _, stmt := range stmts {
			if options.beautify {
				formatter := clickhouse.NewFormatter()
				formatter.WithBeautify()
				formatter.WriteExpr(stmt)
				fmt.Println(formatter.String())
			} else {
				fmt.Println(clickhouse.Format(stmt))
			}
		}
	}
}
