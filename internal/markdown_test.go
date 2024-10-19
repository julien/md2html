package internal_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/julien/md2html/internal"
)

func Test_FileString(t *testing.T) {
	var (
		now     = time.Now()
		y, m, d = now.Date()
		f       = internal.File{
			Name:    "test",
			ModTime: now,
			Path:    "test",
		}
		s        = f.String()
		expected = fmt.Sprintf("%s-%d-%d-%d", f.Name, y, m, d)
	)

	if s != expected {
		t.Fatalf("got %s want %s", s, expected)
	}
}

// func Test_FileContents(t *testing.T) {
// }
