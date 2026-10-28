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
