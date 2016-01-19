// This is a test.
package main

import (
	"go/token"
	"log"
	"os"
	"path/filepath"
)

var fset *token.FileSet

func main() {
	dir := os.Args[1]
	err := RemoveContents("out")
	if err != nil {
		log.Fatal(err)
	}
	transPackage(dir)

}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
