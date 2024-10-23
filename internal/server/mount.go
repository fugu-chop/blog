package server

import (
	"log"

	"github.com/fugu-chop/blog/internal/controllers"
	"github.com/fugu-chop/blog/internal/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*
mount registers various routes on the *chi.Mux router. It parses various templates
that are embedded in the binary via fs.FS and uses HandlerFuncs provided in the
`controllers` package to write these to io.ResponseWriter.
*/
func (s *Server) mount() {
	log.Print("registering routes on server")

	s.Mux.Use(middleware.RequestID)
	s.Mux.Use(middleware.RealIP)
	s.Mux.Use(middleware.Logger)
	s.Mux.Use(middleware.Recoverer)

	// Ensure template can be parsed before attempting to use
	homeTpl := views.GenerateTemplate("root/home.gohtml")
	s.Mux.Get("/", controllers.StaticHandler(homeTpl))

	aboutTpl := views.GenerateTemplate("root/about.gohtml")
	s.Mux.Get("/about", controllers.StaticHandler(aboutTpl))

	blogTpl := views.GenerateTemplate("root/blog.gohtml")
	s.Mux.Get("/blog", controllers.StaticHandler(blogTpl))

	// Project related templates
	projectsTpl := views.GenerateTemplate("root/projects.gohtml")
	websiteTpl := views.GenerateTemplate("projects/website.gohtml")

	s.Mux.Route("/projects", func(r chi.Router) {
		r.Get("/", controllers.StaticHandler(projectsTpl))
		r.Get("/website", controllers.StaticHandler(websiteTpl))
	})

	// make sure to register cookies only for admin page for posting blog
	// use gorilla/csrf to generate csrf token middleware
	// Add sessions to headers
	// Store user in context to minimise db queries
}
