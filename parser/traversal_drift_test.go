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
