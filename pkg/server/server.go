package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fugu-chop/blog/pkg/config"
	"github.com/fugu-chop/blog/pkg/controllers"
	"github.com/fugu-chop/blog/pkg/views"
	"github.com/fugu-chop/blog/pkg/views/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Mux    *chi.Mux
	addr   string
	server *http.Server
}

func New(ctx context.Context, port string) (*Server, error) {
	port = ":" + port
	mux := chi.NewRouter()

	server := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Minute,
		WriteTimeout:      time.Minute,
		MaxHeaderBytes:    8 * 1024,
	}

	s := &Server{
		addr:   port,
		Mux:    mux,
		server: server,
	}

	s.mount()

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	log.Printf("starting server on localhost at port %s...", s.addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down server at address %s", s.addr)

	return s.server.Shutdown(ctx)
}

func (s *Server) mount() {
	log.Print("registering routes on server")

	s.Mux.Use(middleware.RequestID)
	s.Mux.Use(middleware.RealIP)
	s.Mux.Use(middleware.Logger)
	s.Mux.Use(middleware.Recoverer)

	// Ensure template can be before attempting to use
	homeTpl := views.Must(views.ParseFS(templates.FS, config.LayoutTemplate, "home.gohtml"))
	s.Mux.Get("/", controllers.StaticHandler(homeTpl))

	s.Mux.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am Dean")
	})
	s.Mux.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is my blog")
	})
	s.Mux.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to sign in page")
	})
	s.Mux.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Signed in")
	})
	s.Mux.Post("/signout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Signed out")
	})
	s.Mux.Post("/post", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Posted an article")
	})

	// make sure to register cookies only for admin page for posting blog
	// use gorilla/csrf to generate csrf token middleware
	// Add sessions to headers
	// Store user in context to minimise db queries
}
