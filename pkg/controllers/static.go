package controllers

import (
	"net/http"

	"github.com/fugu-chop/blog/pkg/views"
)

/*
StaticHandler writes the data in a views.Template type to
an io.ResponseWriter interface.
*/
func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}
