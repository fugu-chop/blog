package views

import (
	"fmt"
	"html/template"
	"io/fs"
)

func ParseFS(fs fs.FS, patterns ...string) (template.Template, error) {
	tpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return template.Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return *tpl, nil
}

func Must(t template.Template, err error) template.Template {
	if err != nil {
		panic(err)
	}
	return t
}
