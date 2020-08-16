package main

//  rootVIII
//    recursively search and query a Linux filesystem
//    by matching substring and/or permissions

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// QueryFS stores details related to the current search.
type QueryFS struct {
	Substring   string
	Permissions string
	IsTerm      bool
	IsPerm      bool
}

func (q QueryFS) processed(fileName string, processedDirectories []string) bool {
	for i := 0; i < len(processedDirectories); i++ {
		if processedDirectories[i] != fileName {
			continue
		}
		return true
	}
	return false
}

// Query recursively searches entire file-system starting from provided path.
func (q QueryFS) Query(path string, dirs []string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		newPath := fmt.Sprintf("%s/%s", path, f.Name())
		if f.IsDir() {
			if !q.processed(newPath, dirs) {
				q.evaluate(newPath)
				dirs = append(dirs, newPath)
				q.Query(newPath, dirs)
			}
		} else {
			q.evaluate(newPath)
		}
	}
}

func (q QueryFS) evaluate(path string) {
	if q.IsTerm {
		if !strings.Contains(path, q.Substring) {
			goto end
		}
	}
	if q.IsPerm {
		fstat, err := os.Stat(path)
		if err != nil || q.Permissions != fmt.Sprintf("%v", fstat.Mode().Perm()) {
			goto end
		}
	}
	fmt.Printf("%s\n", path)
end:
}

func displayError(e error) {
	fmt.Printf("%v\n", e)
	os.Exit(2)
}

func main() {
	path := flag.String("d", "", "directory")
	term := flag.String("s", "", "substring")
	permissions := flag.String("p", "", "permissions")
	flag.Parse()
	if len(os.Args) < 5 {
		displayError(fmt.Errorf("invalid arguments provided"))
	}

	fmode, err := os.Lstat(*path)
	if err != nil {
		displayError(err)
	}
	if !fmode.Mode().IsDir() {
		displayError(fmt.Errorf("invalid directory provided for start point %s", *path))
	}

	start := *path
	if start != "/" && start[len(*path)-1:] == "/" {
		start = start[:len(*path)-1]
	}

	var qfs = &QueryFS{
		Substring:   *term,
		Permissions: *permissions,
		IsTerm:      len(*term) > 0,
		IsPerm:      len(*permissions) > 0,
	}
	qfs.Query(start, []string{})
}
