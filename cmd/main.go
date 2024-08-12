package main

import (
	"log"
	"os"

	"github.com/fugu-chop/blog/pkg/server"
)

func main() {
	// mux := http.NewServeMux()

	// mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello there!")
	// })

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	svr, err := server.New(port)
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	svr.Start()
}
