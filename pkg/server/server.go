package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fugu-chop/blog/pkg/controllers"
	"github.com/fugu-chop/blog/pkg/views"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*
The Server type encapsulates a multiplexer (using *chi.Mux)
and a standard http.Server type available in the standard library.
*/
type Server struct {
	Mux    *chi.Mux
	addr   string
	server *http.Server
}

/*
New creates a pointer to a Server type. It encapsulates the port
on which the server will listenAndServe, as well as the *chi.Mux
router. It sets some sensible defaults for timeouts and mounts
handlers to the router.
*/
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

/*
Start calls ListenAndServe on the http.Server in the Server type,
allowing it to start traffic. A failure on ListenAndServe that is not
a http.ErrServerClosed (i.e. a server shutdown) will result in a panic.
*/
func (s *Server) Start(ctx context.Context) error {
	log.Printf("starting server on localhost at port %s...", s.addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	return nil
}

/*
Shutdown calls Shutdown on the http.Server embedded within the Server type (i.e.
a graceful shutdown of the server).
*/
func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down server at address %s", s.addr)

	return s.server.Shutdown(ctx)
}

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

	projectsTpl := views.GenerateTemplate("root/projects.gohtml")
	s.Mux.Get("/projects", controllers.StaticHandler(projectsTpl))

	// make sure to register cookies only for admin page for posting blog
	// use gorilla/csrf to generate csrf token middleware
	// Add sessions to headers
	// Store user in context to minimise db queries
}
