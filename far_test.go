package far_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/soyart/far"
)

func TestFindAndReplace(t *testing.T) {
	type testFile struct {
		path     string
		original string
		expected string
	}

	type testCase struct {
		old   string
		new   string
		root  string
		files []testFile
	}

	tests := []testCase{
		{
			old:  "ref-1",
			new:  "something-else",
			root: "testassets",
			files: []testFile{
				{
					path: "testassets/foo.txt",

					original: `# foo.txt

This is ref-1

This is also ref-1`,

					expected: `# foo.txt

This is something-else

This is also something-else`,
				},
				{
					path: "testassets/baz/baz",
					original: `# baz

Some noise

ref-1

ref-2`,
					expected: `# baz

Some noise

something-else

ref-2`,
				},
			},
		},
	}

	for i := range tests {
		tc := &tests[i]
		for j := range tc.files {
			f := &tc.files[j]

			// stat to remember old perm
			stat, err := os.Stat(f.path)
			if err != nil && !os.IsNotExist(err) {
				panic(err)
			}
			if stat != nil && stat.IsDir() {
				panic(fmt.Errorf("unexpected dir '%s' in test cases", f.path))
			}
			perm := os.ModePerm
			if stat != nil {
				perm = stat.Mode().Perm()
			}
			err = os.RemoveAll(f.path)
			if err != nil && !os.IsNotExist(err) {
				panic(err)
			}

			mkdir := filepath.Dir(f.path)
			err = os.MkdirAll(mkdir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile(f.path, []byte(f.original), perm)
			if err != nil {
				panic(err)
			}
		}

		err := far.FindAndReplace(tc.root, tc.old, tc.new)
		if err != nil {
			t.Fatal(err)
		}

		for j := range tc.files {
			f := &tc.files[j]
			actual, err := os.ReadFile(f.path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(actual, []byte(f.expected)) {
				t.Logf("unexpected value for path '%s'", f.path)
				t.Logf("expected: %s", f.expected)
				t.Logf("actual: %s", actual)
				t.Fatal("unexpected value")
			}
		}
	}
}
