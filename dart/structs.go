package dart

import "go/ast"

// Library represents a collection of classes, variables, and functions.
type Library struct {
	Name      string
	Classes   map[string]*Class
	FuncTypes []*ast.GenDecl
	Funcs     []*ast.FuncDecl
	Vars      []*ast.GenDecl
}

func NewLibrary() *Library {
	return &Library{
		Name:    "",
		Classes: map[string]*Class{},
		Funcs:   []*ast.FuncDecl{},
		Vars:    []*ast.GenDecl{},
	}
}

// Class represents a dart class
type Class struct {
	Name    string
	Fields  []*ast.Field
	Methods []*ast.FuncDecl
}
