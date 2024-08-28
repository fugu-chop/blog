package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fugu-chop/blog/pkg/server"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	svr, err := server.New(port)
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go svr.Start(ctx)

	<-signals
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown: %v", err)
	}
}
