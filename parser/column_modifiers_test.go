package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_ColumnModifierOrder(t *testing.T) {
	modifiers := map[string]string{
		"comment": "cOmMeNt 'example'",
		"codec":   "cOdEc(ZSTD(1))",
		"ttl":     "tTl event_time + INTERVAL 1 DAY",
	}
	orders := [][]string{
		{"comment", "codec", "ttl"},
		{"comment", "ttl", "codec"},
		{"codec", "comment", "ttl"},
		{"codec", "ttl", "comment"},
		{"ttl", "comment", "codec"},
		{"ttl", "codec", "comment"},
	}
	forms := []struct{ name, query, columnType string }{
		{"create", "CREATE TABLE db.example (event_time DateTime, %s) ENGINE = MergeTree ORDER BY event_time", "String"},
		{"add", "ALTER TABLE db.example ADD COLUMN %s", "String"},
		{"modify", "ALTER TABLE db.example MODIFY COLUMN %s", "String"},
		{"modify without type", "ALTER TABLE db.example MODIFY COLUMN %s", ""},
	}
	for _, form := range forms {
		for _, order := range orders {
			t.Run(form.name+"/"+strings.Join(order, "_"), func(t *testing.T) {
				var clauses []string
				for _, name := range order {
					clauses = append(clauses, modifiers[name])
				}
				nameAndType := strings.TrimSpace("value " + form.columnType)
				declaration := nameAndType + " " + strings.Join(clauses, " ")
				query := fmt.Sprintf(form.query, declaration)
				statements, err := NewParser(query).ParseStmts()
				require.NoError(t, err)
				require.Len(t, statements, 1)
				column := modifierTestColumn(t, statements[0])
				require.Equal(t, form.columnType, Format(column.Type))
				if form.columnType == "" {
					require.Nil(t, column.Type)
				}
				require.NotNil(t, column.Comment)
				require.NotNil(t, column.Codec)
				require.NotNil(t, column.TTL)
				require.Equal(t, Pos(strings.Index(query, declaration)+len(declaration)), column.End())
				canonical := nameAndType + " COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY"
				require.Equal(t, canonical, Format(column))
				for _, beautify := range []bool{false, true} {
					formatter := NewFormatter()
					if beautify {
						formatter.WithBeautify()
					}
					statements[0].FormatSQL(formatter)
					reparsed, err := NewParser(formatter.String()).ParseStmts()
					require.NoError(t, err)
					require.Len(t, reparsed, 1)
					reparsedColumn := modifierTestColumn(t, reparsed[0])
					require.Equal(t, column.Type == nil, reparsedColumn.Type == nil)
					require.Equal(t, form.columnType, Format(reparsedColumn.Type))
					require.NotNil(t, reparsedColumn.Comment)
					require.NotNil(t, reparsedColumn.Codec)
					require.NotNil(t, reparsedColumn.TTL)
					require.Equal(t, canonical, Format(reparsedColumn))
				}
				validFormatSQL(t, Format(statements[0])+";\n")
			})
		}
	}
}

func TestParser_InvalidColumnModifiers(t *testing.T) {
	for _, tc := range []struct{ name, modifiers, errorText string }{
		{"comment twice", "COMMENT '' COMMENT 'example'", "duplicate COMMENT"},
		{"separated comments", "COMMENT 'first' CODEC(ZSTD(1)) COMMENT 'second'", "duplicate COMMENT"},
		{"codec twice", "CODEC(ZSTD(1)) CODEC(LZ4)", "duplicate CODEC"},
		{"separated codecs", "CODEC(ZSTD(1)) COMMENT 'example' CODEC(LZ4)", "duplicate CODEC"},
		{"ttl twice", "TTL event_time + INTERVAL 1 DAY TTL event_time + INTERVAL 2 DAY", "duplicate TTL"},
		{"separated ttls", "TTL event_time + INTERVAL 1 DAY COMMENT 'example' TTL event_time + INTERVAL 2 DAY", "duplicate TTL"},
		{"missing comment", "COMMENT", ""},
		{"incomplete codec", "CODEC(", ""},
		{"missing ttl", "TTL", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, query := range []string{
				"CREATE TABLE db.example (event_time DateTime, value String " + tc.modifiers + ") ENGINE = MergeTree ORDER BY event_time",
				"ALTER TABLE db.example MODIFY COLUMN value " + tc.modifiers,
			} {
				_, err := NewParser(query).ParseStmts()
				require.Error(t, err, query)
				if tc.errorText != "" {
					require.ErrorContains(t, err, tc.errorText, query)
				}
			}
		})
	}
}

func modifierTestColumn(t *testing.T, expression Expr) *ColumnDef {
	t.Helper()
	var columns []*ColumnDef
	Walk(expression, func(node Expr) bool {
		if column, ok := node.(*ColumnDef); ok && Format(column.Name) == "value" {
			columns = append(columns, column)
		}
		return true
	})
	require.Len(t, columns, 1)
	return columns[0]
}

func TestParser_ColumnValueExprWithoutType(t *testing.T) {
	forms := []struct{ name, query string }{
		{"create", "CREATE TABLE db.example (event_time DateTime, %s) ENGINE = MergeTree ORDER BY event_time"},
		{"add", "ALTER TABLE db.example ADD COLUMN %s"},
		{"modify", "ALTER TABLE db.example MODIFY COLUMN %s"},
	}
	for _, tc := range []struct{ name, declaration, columnType, canonical string }{
		{"default", "value DEFAULT event_time + 1", "", "value DEFAULT event_time + 1"},
		{"materialized", "value mAtErIaLiZeD toDate(event_time)", "", "value MATERIALIZED toDate(event_time)"},
		{"alias", "value ALIAS event_time", "", "value ALIAS event_time"},
		{"parenthesized default", "value DEFAULT (1) COMMENT 'example'", "", "value DEFAULT (1) COMMENT 'example'"},
		{"keyword type", "value Date DEFAULT today()", "Date", "value Date DEFAULT today()"},
		{"parametric type", "value Nullable(String) ALIAS 'x'", "Nullable(String)", "value Nullable(String) ALIAS 'x'"},
	} {
		for _, form := range forms {
			query := fmt.Sprintf(form.query, tc.declaration)
			t.Run(form.name+"/"+tc.name, func(t *testing.T) {
				statements, err := NewParser(query).ParseStmts()
				require.NoError(t, err)
				require.Len(t, statements, 1)
				column := modifierTestColumn(t, statements[0])
				require.Equal(t, tc.columnType, Format(column.Type))
				if tc.columnType == "" {
					require.Nil(t, column.Type)
				}
				require.Equal(t, tc.canonical, Format(column))
				require.Equal(t, Pos(strings.Index(query, tc.declaration)+len(tc.declaration)), column.End())
				validFormatSQL(t, Format(statements[0])+";\n")
			})
		}
	}
	for _, declaration := range []string{"value DEFAULT", "value MATERIALIZED", "value ALIAS"} {
		for _, form := range forms {
			query := fmt.Sprintf(form.query, declaration)
			_, err := NewParser(query).ParseStmts()
			require.Error(t, err, query)
		}
	}
}
