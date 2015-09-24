// This is a test.
package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lologarithm/gopherDart/dart"
)

func main() {
	dir := os.Args[1]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read dir: %s", err)
		return
	}

	lib := dart.NewLibrary()
	for _, f := range files {
		if strings.Contains(f.Name(), ".go") && !strings.Contains(f.Name(), "_test") {
			parse(dir+f.Name(), lib)
		}
	}

	ioutil.WriteFile("lib.dart", convert(lib), 0644)
}

func convert(lib *dart.Library) []byte {
	return Print(lib)
}

func parse(fn string, lib *dart.Library) {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fn, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	// ast.Print(fset, f)
	LoadToLibrary(f, lib)
}
