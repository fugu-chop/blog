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

func New(port string) (*Server, error) {
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
		addr:   port,
		Mux:    mux,
		server: server,
	}

	s.mount()

	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	log.Printf("starting server on localhost at port %s...", s.addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("Shuttingdown server at address %s", s.addr)

	return s.server.Shutdown(ctx)
}

func (s *Server) mount() {
	log.Print("registering routes on server")

	s.Mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s route invoked", r.Method, r.URL)
		fmt.Fprintf(w, "Hello there!")
	})
}
