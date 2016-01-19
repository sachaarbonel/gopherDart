// This is a test.
package main

import (
	"go/token"
	"os"
)

var fset *token.FileSet

func main() {
	dir := os.Args[1]
	transPackage(dir)
}
