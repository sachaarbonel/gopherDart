package main

import (
	"bytes"
	"go/ast"
	"go/token"
	"log"
	"testing"

	"github.com/lologarithm/gopherDart/dart"
)

// TestIf tests the if statement.
func TestIf(t *testing.T) {
	lib := dart.NewLibrary()
	lib.Name = "Test"
	ctx := &LibraryContext{
		Name:        lib.Name,
		Indentation: "",
		Class:       nil,
	}
	stmt := &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{&ast.BasicLit{Value: "1"}},
		},
		Cond: &ast.BinaryExpr{
			X: &ast.Ident{
				Name: "x",
			},
			Op: token.EQL,
			Y: &ast.BasicLit{
				Value: "1",
			},
		},
		Body: &ast.BlockStmt{},
		Else: &ast.BlockStmt{}, // else branch; or nil
	}
	buf := &bytes.Buffer{}
	printStmt(stmt, buf, "", ctx)
	log.Printf("Buffer:\n%s", buf.String())
}
