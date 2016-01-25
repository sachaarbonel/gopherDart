package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

func getCompileFiles(dir string) ([]string, error) {
	old_dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = os.Chdir(dir)
	defer os.Chdir(old_dir)                                   // woo using defers, change working directory back.
	cmd := exec.Command("go", "list", "-f", "'{{.GoFiles}}'") //heheheh
	var out []byte
	out, err = cmd.Output()
	if err != nil {
		return nil, err
	}
	filestr := string(out)
	re := regexp.MustCompile("[\\[\\]']")
	filestr = re.ReplaceAllLiteralString(filestr, "")
	re = regexp.MustCompile("\\s")
	return re.Split(filestr, -1), nil
}

func outputFile(fname string, toWrite []byte) error {
	//run_dir := os.Args[0]

	return ioutil.WriteFile(path.Join("lib", fname), toWrite, 0644)
}

func libName(dir string) string {
	st := path.Base(dir)
	return st

}

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
