package main

import (
	"fmt"
	"log"

	"github.com/AfterShip/clickhouse-sql-parser/parser"
)

func main() {
	// Test MODIFY SETTING
	sql1 := "ALTER TABLE example_table MODIFY SETTING max_part_loading_threads=8, max_parts_in_total=50000;"
	lexer1 := parser.NewLexer(sql1)
	p1 := parser.NewParser(lexer1)
	stmt1, err1 := p1.ParseDDL()
	if err1 != nil {
		log.Fatalf("Error parsing MODIFY SETTING: %v", err1)
	}
	fmt.Printf("MODIFY SETTING parsed successfully: %s\n", stmt1.String())

	// Test RESET SETTING
	sql2 := "ALTER TABLE example_table RESET SETTING max_part_loading_threads;"
	lexer2 := parser.NewLexer(sql2)
	p2 := parser.NewParser(lexer2)
	stmt2, err2 := p2.ParseDDL()
	if err2 != nil {
		log.Fatalf("Error parsing RESET SETTING: %v", err2)
	}
	fmt.Printf("RESET SETTING parsed successfully: %s\n", stmt2.String())

	// Test RESET multiple SETTINGs
	sql3 := "ALTER TABLE example_table RESET SETTING max_part_loading_threads, max_parts_in_total, another_setting;"
	lexer3 := parser.NewLexer(sql3)
	p3 := parser.NewParser(lexer3)
	stmt3, err3 := p3.ParseDDL()
	if err3 != nil {
		log.Fatalf("Error parsing RESET multiple SETTINGs: %v", err3)
	}
	fmt.Printf("RESET multiple SETTINGs parsed successfully: %s\n", stmt3.String())
}