// This is a test.
package main

import (
	//"github.com/lologarithm/gopherDart/dart"
	"go/token"
	"os"
)

var fset *token.FileSet

func main() {
	dir := os.Args[1]
	transPackage(dir)
}
