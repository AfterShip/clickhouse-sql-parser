package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkParseSQLFiles benchmarks parsing all SQL files in the testdata/query directory
func BenchmarkParseSQLFiles(b *testing.B) {
	testFiles, err := filepath.Glob("testdata/query/*.sql")
	if err != nil {
		b.Fatalf("Failed to glob test files: %v", err)
	}

	for _, file := range testFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			b.Fatalf("Failed to read file %s: %v", file, err)
		}

		b.Run(filepath.Base(file), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				parser := NewParser(string(content))
				_, err := parser.ParseStmts()
				if err != nil {
					b.Fatalf("Failed to parse SQL from %s: %v", file, err)
				}
			}
		})
	}
}

// BenchmarkParseComplexQueries benchmarks parsing specifically complex SQL queries
func BenchmarkParseComplexQueries(b *testing.B) {
	complexQueries := []string{
		"testdata/query/select_with_multi_join.sql",
		"testdata/query/select_with_window_function.sql",
		"testdata/query/select_simple_with_with_clause.sql",
		"testdata/query/select_with_left_join.sql",
		"testdata/benchdata/posthog_huge_0.sql",
		"testdata/benchdata/posthog_huge_1.sql",
	}

	for _, queryFile := range complexQueries {
		content, err := os.ReadFile(queryFile)
		if err != nil {
			b.Fatalf("Failed to read file %s: %v", queryFile, err)
		}

		b.Run(queryFile, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				parser := NewParser(string(content))
				_, err := parser.ParseStmts()
				if err != nil {
					b.Fatalf("Failed to parse SQL from %s: %v", queryFile, err)
				}
			}
		})
	}
}
