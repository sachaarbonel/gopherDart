// This is a test.
package main

import (
	"os"
	"path/filepath"
)

func main() {
	dir := os.Args[1]
	// err := RemoveContents("lib")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	transpile(dir)

	// files, err := ioutil.ReadDir("standlib")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, f := range files {
	// 	exec.Command("cp", "standlib/"+f.Name(), path.Join(dir, "lib")).Run()
	// }

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
