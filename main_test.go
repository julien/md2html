package main

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_Run(t *testing.T) {
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
			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			var (
				dir = "testdata"
				src = wd + string(os.PathSeparator) + dir + string(os.PathSeparator) + tc.src
				dst = wd + string(os.PathSeparator) + dir + string(os.PathSeparator) + tc.dst
			)

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

			if err := Run(abssrc, absdst); err != nil && !tc.fail {
				t.Fatal(err)
			}
		})
	}
}
