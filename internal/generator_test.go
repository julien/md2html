package internal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/julien/md2html/internal"
)

func Test_New(t *testing.T) {
	tcs := []struct {
		desc string
		err  error
	}{
		{
			desc: "returns a new instance",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := internal.New()
			if err != tc.err {
				t.Fatalf("got %v want: %v", err, tc.err)
			}
			if g == nil {
				t.Fatalf("got nil, want a new instance")
			}
		})
	}
}

func Test_Generate(t *testing.T) {
	tcs := []struct {
		desc string
		src  string
		dst  string
		fail bool
	}{
		{
			desc: "happy path",
			src:  "ok",
			dst:  "tmp",
		},
		{
			desc: "no content",
			src:  "nocontent",
			dst:  "tmp",
			fail: true,
		},
		{
			desc: "non exisiting",
			src:  "nonexisting",
			dst:  "tmp",
			fail: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := internal.New()
			if err != nil {
				t.Fatal(err)
			}

			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			testdir := "testdata"

			src := wd + string(os.PathSeparator) + testdir + string(os.PathSeparator) + tc.src
			dst := wd + string(os.PathSeparator) + testdir + string(os.PathSeparator) + tc.dst

			defer func() {
				if err := os.RemoveAll(dst); err != nil {
					t.Fatal(err)
				}
			}()

			abssrc, err := filepath.Abs(src)
			if err != nil {
				t.Fatal(err)
			}

			absdst, err := filepath.Abs(dst)
			if err != nil {
				t.Fatal(err)
			}

			if err := g.Generate(abssrc, absdst); err != nil && !tc.fail {
				t.Fatal(err)
			}
		})
	}
}
