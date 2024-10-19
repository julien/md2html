package internal

import (
	"errors"
	"fmt"
	_ "html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const Extension = ".md"

var (
	ErrNoFiles = errors.New("no .md file found")
	replacer   = strings.NewReplacer(" ", "", Extension, "")
)

type File struct {
	Name    string
	ModTime time.Time
	Path    string
}

func (f File) String() string {
	return fmt.Sprintf("%s-%s",
		replacer.Replace(f.Name),
		replacer.Replace(f.ModTime.Format(time.DateOnly)))
}

func (f File) Contents() ([]byte, error) {
	fi, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(fi)
	if err != nil {
		return nil, err
	}
	if err := fi.Close(); err != nil {
		return nil, err
	}
	return b, nil
}

func Find(dir string) (chan File, chan error) {
	var (
		sys = os.DirFS(dir)
		fch = make(chan File, 1)
		ech = make(chan error, 1)
	)

	base, err := filepath.Abs(dir)
	if err != nil {
		ech <- err
		close(fch)
		close(ech)
		return fch, ech
	}

	go func() {
		defer func() {
			close(fch)
			close(ech)
		}()

		fs.WalkDir(sys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				ech <- err
				return err
			}

			if d.IsDir() || !strings.HasSuffix(path, Extension) {
				return nil
			}

			i, err := d.Info()
			if err != nil {
				ech <- err
				return err
			}

			abs, err := filepath.Abs(base + string(os.PathSeparator) + path)
			if err != nil {
				ech <- err
				return err
			}

			fch <- File{ModTime: i.ModTime(), Name: path, Path: abs}
			return nil
		})
	}()
	return fch, ech
}

func ToHTML(f File) ([]byte, error) {
	var (
		p      = parser.NewWithExtensions(parser.CommonExtensions | parser.NoEmptyLineBeforeBlock)
		r      = html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})
		b, err = f.Contents()
	)
	if err != nil {
		return nil, err
	}
	return markdown.Render(p.Parse(b), r), nil
}
