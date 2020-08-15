package main

//   rootVIII
//   recursively search a *NIX filesystem

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

// QueryFS stores details related to the current search.
type QueryFS struct {
	Term        string
	Permissions string
	Owner       string
	Group       string
	IsTerm      bool
	IsPerm      bool
	IsOwner     bool
	IsGroup     bool
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

// ListDirContents recursively searches the entire
// file-system starting from the provided path.
func (q QueryFS) ListDirContents(path string, dirs []string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		newPath := fmt.Sprintf("%s/%s", path, f.Name())
		if f.IsDir() {
			if !q.processed(newPath, dirs) {
				q.evaluate(newPath)
				dirs = append(dirs, newPath)
				q.ListDirContents(newPath, dirs)
			}
		} else {
			q.evaluate(newPath)
		}
	}
}

func (q QueryFS) evaluate(path string) {
	if q.IsTerm {
		if !strings.Contains(path, q.Term) {
			goto end
		}
	}
	if q.IsPerm {
		fstat, err := os.Stat(path)
		if err != nil || q.Permissions != fmt.Sprintf("%v", fstat.Mode().Perm()) {
			goto end
		}
	}
	if q.IsOwner || q.IsGroup {
		fstat, err := os.Stat(path)
		if err != nil {
			goto end
		}
		stat, _ := fstat.Sys().(*syscall.Stat_t)
		UID, GID := stat.Uid, stat.Gid
		fmt.Printf("UID: %d, GID: %d\n", UID, GID)
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
	term := flag.String("t", "", "term")
	permissions := flag.String("p", "", "permissions")
	owner := flag.String("o", "", "owner")
	group := flag.String("g", "", "group")
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
		Term:        *term,
		Permissions: *permissions,
		Owner:       *owner,
		Group:       *group,
		IsTerm:      len(*term) > 0,
		IsPerm:      len(*permissions) > 0,
		IsOwner:     len(*owner) > 0,
		IsGroup:     len(*group) > 0,
	}
	qfs.ListDirContents(start, []string{})
}
