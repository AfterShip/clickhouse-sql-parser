package parser

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// The package has two independent traversal engines: Accept/ASTVisitor and
// Walk/WalkFunc. Each encodes every node's children separately, so a new AST
// node type added to one can silently be forgotten in the other. This test
// statically asserts that every type with an Accept method also has a Visit
// method on the ASTVisitor interface and a case in Walk's type switch (and
// vice versa), so the engines cannot drift.
func TestTraversalEnginesCoverSameNodeTypes(t *testing.T) {
	entries, err := os.ReadDir(".")
	require.NoError(t, err)

	fset := token.NewFileSet()
	acceptTypes := map[string]bool{}
	visitorTypes := map[string]bool{}
	walkTypes := map[string]bool{}

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := goparser.ParseFile(fset, name, nil, 0)
		require.NoError(t, err)
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				switch {
				case d.Name.Name == "Accept" && d.Recv != nil && len(d.Recv.List) == 1:
					if star, ok := d.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := star.X.(*ast.Ident); ok {
							acceptTypes[ident.Name] = true
						}
					}
				case d.Name.Name == "Walk" && d.Recv == nil:
					ast.Inspect(d.Body, func(n ast.Node) bool {
						cc, ok := n.(*ast.CaseClause)
						if !ok {
							return true
						}
						for _, expr := range cc.List {
							if star, ok := expr.(*ast.StarExpr); ok {
								if ident, ok := star.X.(*ast.Ident); ok {
									walkTypes[ident.Name] = true
								}
							}
						}
						return true
					})
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok || ts.Name.Name != "ASTVisitor" {
						continue
					}
					iface, ok := ts.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}
					for _, method := range iface.Methods.List {
						ft, ok := method.Type.(*ast.FuncType)
						if !ok || ft.Params == nil || len(ft.Params.List) != 1 {
							continue
						}
						if star, ok := ft.Params.List[0].Type.(*ast.StarExpr); ok {
							if ident, ok := star.X.(*ast.Ident); ok {
								visitorTypes[ident.Name] = true
							}
						}
					}
				}
			}
		}
	}

	require.NotEmpty(t, acceptTypes)
	require.NotEmpty(t, visitorTypes)
	require.NotEmpty(t, walkTypes)

	require.Empty(t, diffSet(acceptTypes, visitorTypes),
		"types with an Accept method but no ASTVisitor Visit method")
	require.Empty(t, diffSet(visitorTypes, acceptTypes),
		"types with an ASTVisitor Visit method but no Accept method")
	require.Empty(t, diffSet(acceptTypes, walkTypes),
		"types with an Accept method but no case in Walk's type switch")
	require.Empty(t, diffSet(walkTypes, acceptTypes),
		"types with a case in Walk's type switch but no Accept method")
}

// diffSet returns the members of a that are not in b, sorted.
func diffSet(a, b map[string]bool) []string {
	var diff []string
	for name := range a {
		if !b[name] {
			diff = append(diff, name)
		}
	}
	sort.Strings(diff)
	return diff
}

// VisitJoinTableExpr must be called even when the table expression carries a
// SAMPLE clause; Accept used to return early after visiting the sample ratio.
func TestVisitJoinTableExprWithSampleRatio(t *testing.T) {
	stmts, err := NewParser("SELECT * FROM t SAMPLE 1/10").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	var visitedJoinTable, visitedSample bool
	visitor := &DefaultASTVisitor{
		Visit: func(expr Expr) error {
			switch expr.(type) {
			case *JoinTableExpr:
				visitedJoinTable = true
			case *SampleClause:
				visitedSample = true
			}
			return nil
		},
	}
	require.NoError(t, stmts[0].Accept(visitor))
	require.True(t, visitedSample, "SampleClause was not visited")
	require.True(t, visitedJoinTable, "JoinTableExpr was not visited")
}

// TestTraversalEnginesVisitSameFields statically asserts that, for every node
// type, the set of child fields referenced by its Accept method matches the
// set referenced by its case in Walk's type switch. The type-level test above
// cannot catch a child field that one engine traverses and the other forgot
// (e.g. Walk missing InsertStmt.Values while Accept visits it).
func TestTraversalEnginesVisitSameFields(t *testing.T) {
	entries, err := os.ReadDir(".")
	require.NoError(t, err)

	fset := token.NewFileSet()
	acceptFields := map[string]map[string]bool{}
	walkFields := map[string]map[string]bool{}

	// collectSelectors records every selector `<base>.<Field>` in node whose
	// base is the identifier baseName, into out.
	collectSelectors := func(node ast.Node, baseName string, out map[string]bool) {
		ast.Inspect(node, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == baseName {
				out[sel.Sel.Name] = true
			}
			return true
		})
	}

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := goparser.ParseFile(fset, name, nil, 0)
		require.NoError(t, err)
		for _, decl := range file.Decls {
			d, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			switch {
			case d.Name.Name == "Accept" && d.Recv != nil && len(d.Recv.List) == 1:
				star, ok := d.Recv.List[0].Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				typeIdent, ok := star.X.(*ast.Ident)
				if !ok {
					continue
				}
				recvName := d.Recv.List[0].Names[0].Name
				fields := map[string]bool{}
				collectSelectors(d.Body, recvName, fields)
				acceptFields[typeIdent.Name] = fields
			case d.Name.Name == "Walk" && d.Recv == nil:
				ast.Inspect(d.Body, func(n ast.Node) bool {
					cc, ok := n.(*ast.CaseClause)
					if !ok {
						return true
					}
					fields := map[string]bool{}
					for _, stmt := range cc.Body {
						collectSelectors(stmt, "n", fields)
					}
					for _, expr := range cc.List {
						star, ok := expr.(*ast.StarExpr)
						if !ok {
							continue
						}
						if ident, ok := star.X.(*ast.Ident); ok {
							walkFields[ident.Name] = fields
						}
					}
					return true
				})
			}
		}
	}

	require.NotEmpty(t, acceptFields)
	require.NotEmpty(t, walkFields)

	var problems []string
	for typeName, aFields := range acceptFields {
		wFields, ok := walkFields[typeName]
		if !ok {
			continue // type-level coverage is asserted by the test above
		}
		for _, field := range diffSet(aFields, wFields) {
			problems = append(problems,
				typeName+"."+field+" is traversed by Accept but not by Walk")
		}
		for _, field := range diffSet(wFields, aFields) {
			problems = append(problems,
				typeName+"."+field+" is traversed by Walk but not by Accept")
		}
	}
	sort.Strings(problems)
	require.Empty(t, problems, "child-field traversal drift between Accept and Walk")
}

// TestInsertStmtEndWithoutValues guards that End() does not panic on
// `INSERT INTO t FORMAT CSV`, where the data arrives out of band and the
// statement carries neither VALUES nor a SELECT.
func TestInsertStmtEndWithoutValues(t *testing.T) {
	sql := "INSERT INTO t FORMAT CSV"
	stmts, err := NewParser(sql).ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	insert := stmts[0].(*InsertStmt)
	require.Empty(t, insert.Values)
	require.Nil(t, insert.SelectExpr)
	require.Equal(t, Pos(len(sql)), insert.End())
}

// TestWalkVisitsInsertValues guards that Walk descends into INSERT ... VALUES
// rows; the InsertStmt case used to skip the Values field entirely.
func TestWalkVisitsInsertValues(t *testing.T) {
	stmts, err := NewParser("INSERT INTO t VALUES (1, 2), (3, 4)").ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	values := 0
	Walk(stmts[0], func(node Expr) bool {
		if _, ok := node.(*AssignmentValues); ok {
			values++
		}
		return true
	})
	require.Equal(t, 2, values, "both VALUES rows should be walked")
}

// TestDestinationTableSchemaVisitedOnce guards that a materialized view's TO
// destination column list is visited exactly once, inside the destination
// clause, by both traversal engines.
func TestDestinationTableSchemaVisitedOnce(t *testing.T) {
	sql := "CREATE MATERIALIZED VIEW mv TO dest (id UInt64) AS SELECT id FROM src"
	stmts, err := NewParser(sql).ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	acceptVisits := 0
	visitor := &DefaultASTVisitor{
		Visit: func(expr Expr) error {
			if _, ok := expr.(*TableSchemaClause); ok {
				acceptVisits++
			}
			return nil
		},
	}
	require.NoError(t, stmts[0].Accept(visitor))

	walkVisits := 0
	Walk(stmts[0], func(node Expr) bool {
		if _, ok := node.(*TableSchemaClause); ok {
			walkVisits++
		}
		return true
	})
	require.Equal(t, 1, acceptVisits, "Accept should visit the destination schema exactly once")
	require.Equal(t, walkVisits, acceptVisits, "Accept and Walk disagree on schema visits")
}

func TestTraversalVisitsJSONTypeOptions(t *testing.T) {
	sql := "CREATE TABLE t (j JSON(max_dynamic_paths=8, max_dynamic_types=9, a.b String, SKIP c.d, SKIP REGEXP 're')) ENGINE = Memory"
	stmts, err := NewParser(sql).ParseStmts()
	require.NoError(t, err)
	require.Len(t, stmts, 1)

	type visitedValues struct {
		idents  []string
		numbers []string
		strings []string
	}
	collect := func(values *visitedValues, expr Expr) {
		switch node := expr.(type) {
		case *Ident:
			values.idents = append(values.idents, node.Name)
		case *NumberLiteral:
			values.numbers = append(values.numbers, node.Literal)
		case *StringLiteral:
			values.strings = append(values.strings, node.Literal)
		}
	}

	var acceptValues visitedValues
	visitor := &DefaultASTVisitor{
		Visit: func(expr Expr) error {
			collect(&acceptValues, expr)
			return nil
		},
	}
	require.NoError(t, stmts[0].Accept(visitor))

	var walkValues visitedValues
	Walk(stmts[0], func(node Expr) bool {
		collect(&walkValues, node)
		return true
	})

	for _, values := range []visitedValues{acceptValues, walkValues} {
		require.Contains(t, values.idents, "a")
		require.Contains(t, values.idents, "b")
		require.Contains(t, values.idents, "String")
		require.Contains(t, values.idents, "c")
		require.Contains(t, values.idents, "d")
		require.Contains(t, values.numbers, "8")
		require.Contains(t, values.numbers, "9")
		require.Contains(t, values.strings, "re")
	}
}
