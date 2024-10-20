package controllers

import (
	"net/http"

	"github.com/fugu-chop/blog/pkg/views"
)

/*
StaticHandler writes the data in a views.Executer interface type to
an io.ResponseWriter interface.
*/
func StaticHandler(tpl views.Executer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}
