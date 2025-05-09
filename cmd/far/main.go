package main

import (
	"os"

	"github.com/soyart/far"
)

func main() {
	l := len(os.Args)
	if l < 3 {
		panic("expecting at least 2 arguments")
	}
	// progname <old> <new> [path]
	old, new, root := os.Args[1], os.Args[2], "."
	if l > 3 {
		root = os.Args[3]
	}
	err := far.FindAndReplace(root, old, new)
	if err != nil {
		panic(err)
	}
}
