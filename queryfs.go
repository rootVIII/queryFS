package main

//    recursively search and query a Linux filesystem
//    by matching substring and/or permissions

/*
MIT License

Copyright (c) 2020 rootVIII

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// QueryFS stores details related to the current search.
type QueryFS struct {
	Substring    string
	Permissions  string
	OwnerGroup   string
	IsTerm       bool
	IsPerm       bool
	IsOwnerGroup bool
	Group        map[int]string
	Passwd       map[int]string
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
		var newPath string
		if path != "/" {
			newPath = fmt.Sprintf("%s/%s", path, f.Name())
		} else {
			newPath = fmt.Sprintf("%s%s", path, f.Name())
		}

		if f.IsDir() {
			if !q.processed(newPath, dirs) {
				q.evaluate(newPath, true)
				dirs = append(dirs, newPath)
				q.Query(newPath, dirs)
			}
		} else {
			q.evaluate(newPath, false)
		}
	}
}

func (q QueryFS) evaluate(path string, isDirectory bool) {
	if q.IsTerm {
		if !strings.Contains(path, q.Substring) {
			goto end
		}
	}
	if q.IsPerm {
		fstat, err := os.Stat(path)
		if err != nil {
			goto end
		}
		var permissions string
		if isDirectory {
			tmp := []byte(fmt.Sprintf("%v", fstat.Mode().Perm()))
			tmp[0] = 0x64
			permissions = string(tmp)
		} else {
			permissions = fmt.Sprintf("%v", fstat.Mode().Perm())
		}
		if q.Permissions != permissions {
			goto end
		}
	}
	if q.IsOwnerGroup {
		fstat, err := os.Stat(path)
		if err != nil {
			goto end
		}
		if stat, ok := fstat.Sys().(*syscall.Stat_t); ok {
			og := fmt.Sprintf("%s:%s", q.Passwd[int(stat.Uid)], q.Group[int(stat.Gid)])
			if q.OwnerGroup != og {
				goto end
			}
		}
	}
	fmt.Printf("%s\n", path)
end:
}

func (q *QueryFS) parseIDS(idFile string, wgroup *sync.WaitGroup) {
	defer wgroup.Done()
	for _, line := range strings.Split(readIn(idFile), "\n") {
		if len(line) > 0 {
			lineItems := strings.Split(line, ":")
			lineNo, _ := strconv.Atoi(lineItems[2])
			if idFile[len(idFile)-6:] != "passwd" {
				q.Group[int(lineNo)] = lineItems[0]
			} else {
				q.Passwd[int(lineNo)] = lineItems[0]
			}
		}
	}
}

// SetIDS sets contents of /etc/passwd and /etc/group if querying with owner:group.
func (q *QueryFS) SetIDS() {
	if q.IsOwnerGroup {
		q.Group = make(map[int]string)
		q.Passwd = make(map[int]string)

		var wg sync.WaitGroup
		wg.Add(2)
		go q.parseIDS("/etc/group", &wg)
		go q.parseIDS("/etc/passwd", &wg)
		wg.Wait()
	}
}

func readIn(absPath string) string {
	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		displayError(err)
	}
	return string(file)
}

func displayError(e error) {
	fmt.Printf("%v\n", e)
	os.Exit(1)
}

func main() {
	path := flag.String("d", "", "directory")
	term := flag.String("s", "", "substring")
	og := flag.String("o", "", "owner:group")
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
		Substring:    *term,
		Permissions:  *permissions,
		OwnerGroup:   *og,
		IsTerm:       len(*term) > 0,
		IsPerm:       len(*permissions) > 0,
		IsOwnerGroup: len(*og) > 0,
	}
	qfs.SetIDS()
	qfs.Query(start, []string{})
}
