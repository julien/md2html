package internal

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"html/template"
)

//go:embed templates/layout.html
var layout string

var ErrNoContent = errors.New("no content")

type Generator struct {
	tpl *template.Template
}

func New() (*Generator, error) {
	tpl, err := template.New("layout").Parse(layout)
	if err != nil {
		return nil, err
	}

	return &Generator{
		tpl: tpl,
	}, nil
}

func (g *Generator) Generate(src, dst string) error {
	var (
		fch, ech = Find(src)
		errs     = make([]error, 0)
	)

loop:
	for {
		select {
		case f, ok := <-fch:
			if !ok {
				fch = nil
				if ech == nil {
					break loop
				}
				continue
			}

			if _, err := g.writeFile(dst, f); err != nil {
				errs = append(errs, err)
				fmt.Printf("error: %v\n", err)
				continue
			}
		case e, ok := <-ech:
			if !ok {
				ech = nil
				if fch == nil {
					break loop
				}
				continue
			}
			errs = append(errs, e)
		}
	}

	return errors.Join(errs...)
}

func (g *Generator) writeFile(dst string, f File) (string, error) {
	b, err := ToHTML(f)
	if err != nil {
		return "", err
	}

	if len(b) == 0 {
		return "", ErrNoContent
	}

	abs, err := filepath.Abs(dst)
	if err != nil {
		return "", err
	}
	n := strings.Replace(f.Name, Extension, ".html", -1)
	out, err := os.Create(abs + string(os.PathSeparator) + n)
	if err != nil {
		return "", err
	}
	defer out.Close()

	return out.Name(), g.tpl.Execute(out, template.HTML(b))
}
