package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Mux    *http.ServeMux
	addr   string
	server *http.Server
}

func New(ctx context.Context, port string) (*Server, error) {
	port = ":" + port
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Minute,
		WriteTimeout:      time.Minute,
		MaxHeaderBytes:    8 * 1024,
	}

	s := &Server{
		addr: port,
		// Expose the Mux on the server as the Handler type on the
		// server doesn't implement the HandleFunc or Handle interfaces
		// i.e. we cannot register routes on the Handler type
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

	s.Mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s route invoked", r.Method, r.URL)
		fmt.Fprintf(w, "Hello there!")
	})

	// make sure to register cookies only for admin page for posting blog
	// use gorilla/csrf to generate csrf token middleware
	// Store user in context to minimise db queries
}