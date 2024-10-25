package internal

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"html/template"

	"github.com/gomarkdown/markdown"
	markdownhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const md = ".md"

//go:embed layout.html
var layout string
var (
	tpl          *template.Template
	errNoContent = errors.New("no content")
	repl         = strings.NewReplacer(md, ".html")
)

type config struct {
	src string
	dir bool
}

func createConfig(src string) (config, error) {
	abs, err := filepath.Abs(src)
	if err != nil {
		return config{}, err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return config{}, err
	}

	return config{
		src: abs,
		dir: info.IsDir(),
	}, nil
}

func init() {
	var err error
	tpl, err = template.New("layout").Parse(layout)
	if err != nil {
		panic(err)
	}
}

func Run(src, dst string) error {
	cfg, err := createConfig(src)
	if err != nil {
		return err
	}

	if !cfg.dir {
		name, err := writeFile(dst, cfg.src)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "wrote %s\n", name)
		return nil
	}

	var (
		fch, ech = find(src)
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

			name, err := writeFile(dst, f)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			fmt.Fprintf(os.Stdout, "wrote %s\n", name)
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

func find(name string) (chan string, chan error) {
	var (
		sys = os.DirFS(name)
		fch = make(chan string, 1)
		ech = make(chan error, 1)
	)

	base, err := filepath.Abs(name)
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

			if d.IsDir() || !strings.HasSuffix(path, md) {
				return nil
			}

			abs, err := filepath.Abs(base + string(os.PathSeparator) + path)
			if err != nil {
				ech <- err
				return err
			}

			fch <- abs
			return nil
		})
	}()
	return fch, ech
}

func writeFile(dst, name string) (string, error) {
	b, err := mdToHTML(name)
	if err != nil {
		return "", err
	}

	if len(b) == 0 {
		return "", errNoContent
	}

	abs, err := filepath.Abs(dst)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(abs, os.ModePerm); err != nil {
		return "", err
	}

	out, err := os.Create(abs + string(os.PathSeparator) + repl.Replace(filepath.Base(name)))
	if err != nil {
		return "", err
	}
	defer out.Close()

	return out.Name(), tpl.Execute(out, template.HTML(b))
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func mdToHTML(name string) ([]byte, error) {
	var (
		p      = parser.NewWithExtensions(parser.CommonExtensions | parser.NoEmptyLineBeforeBlock)
		r      = markdownhtml.NewRenderer(markdownhtml.RendererOptions{Flags: markdownhtml.CommonFlags | markdownhtml.HrefTargetBlank})
		b, err = readFile(name)
	)

	if err != nil {
		return nil, err
	}
	return markdown.Render(p.Parse(b), r), nil
}
