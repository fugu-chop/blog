package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fugu-chop/blog/internal/config"
	"github.com/fugu-chop/blog/internal/templates"
)

type Executer interface {
	Execute(http.ResponseWriter, *http.Request, any)
}

/*
TemplateCloner extracts the Clone() method on the html/template.Template type
to an interface.
*/
type TemplateCloner interface {
	Clone() (*template.Template, error)
}

/*
Template is a type that encapsulates a *template.Template type and
a method to write that template to a io.ResponseWriter. The htmlTpl
field should implement the TemplateCloner interface (usually by passing in
a html/template.Template type).
*/
type Template struct {
	htmlTpl TemplateCloner
}

/*
GenerateTemplate takes in an arbitrary number of strings which represent gohtml
template files and calls various helper methods.

It parses the validity of the template before returning a Template type (an invalid
template will cause a panic).
*/
func GenerateTemplate(patterns ...string) Template {
	inbuiltPatterns := []string{config.LayoutTemplate}
	patterns = append(inbuiltPatterns, patterns...)
	return must(parseFS(templates.FS, patterns...))
}

/*
Execute writes to http.ResponseWriter and sets the `Content-Type` header to `text/html`.

Internally it clones the underlying *template.Template type to which it is attached. This
helps avoid concurrency issues where a template is being used across different web requests
(different goroutines).
*/
func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	// Cloning ensures we don't end up sharing a template across goroutines
	// This can avoid issues with CSRF tokens overwritten between different
	// requests as each request is a different goroutine
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

/*
ParseFS attempts to open FileSystem and apply templates sequentially.

Templates are passed to the `patterns` parameter and applied in the
order they are passed. This enables use of templating within .gohtml templates.
*/
func parseFS(fs fs.FS, patterns ...string) (Template, error) {
	/*
		The html/template package does not include the path in a template name
		we use the base filepath to fix this for nesting templates within folders
	*/
	tpl := template.New(filepath.Base(patterns[0]))
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

/*
Must ensures that templates can be parsed before they are used.

A function that parses a template should be passed to the `err` parameter.
*/
func must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
