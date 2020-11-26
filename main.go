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
	"os"

	"github.com/rootVIII/queryfs/query"
)

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

	var qfs = &query.QFS{
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
