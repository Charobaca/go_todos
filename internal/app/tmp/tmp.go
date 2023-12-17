package tmp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var TMPL *template.Template

func ParseTemplates() error {
	t := template.New("").Funcs(sprig.FuncMap())

	err := filepath.Walk("internal/templates", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			tmplBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = t.New(path).Funcs(sprig.FuncMap()).Parse(string(tmplBytes))
			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		return err
	}

	TMPL = t

	return nil
}