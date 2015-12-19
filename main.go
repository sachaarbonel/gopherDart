// This is a test.
package main

import (
	"fmt"
	//"github.com/lologarithm/gopherDart/dart"
	"go/ast"
	"go/parser"
	"go/token"
	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var fset *token.FileSet

func main() {
	dir := os.Args[1]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read dir: %s", err)
		return
	}

	lib := NewLibrary()
	parsed := make([]*ast.File, len(files))
	count := 0
	fset = token.NewFileSet()
	for _, fi := range files {
		if strings.Contains(fi.Name(), ".go") && !strings.Contains(fi.Name(), "_test") {
			f, err := parser.ParseFile(fset, filepath.Join(dir, fi.Name()), nil, 0)
			if err == nil {
				parsed[count] = f
				count++
			}

		}
	}

	parsed = parsed[:count]
	info := &types.Info{Defs: make(map[*ast.Ident]types.Object), Uses: make(map[*ast.Ident]types.Object), Types: make(map[ast.Expr]types.TypeAndValue)}
	cfg := &types.Config{}
	_, err = cfg.Check("gopherDart", fset, parsed, info)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range parsed {
		LoadToLibrary(f, lib, info)
	}

	ioutil.WriteFile("lib.dart", convert(lib), 0644)
}

func convert(lib *Library) []byte {
	return Print(lib)
}

func report(n ast.Node) {
	pos := n.Pos()
	fmt.Println("Problem at " + fset.Position(pos).String() + " with ")
}
