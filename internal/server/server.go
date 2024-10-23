package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
