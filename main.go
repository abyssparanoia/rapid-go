package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

var addr = ":3001"

func main() {
	// Dependency
	d := Dependency{}
	d.Inject()

	// Routing
	r := chi.NewRouter()
	Routing(r, d)

	//server
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Run
	fmt.Printf("[START] server. port: %s\n", addr)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("[CLOSED] server closed with error: %s\n", err)
		}
	}()

	// graceful shuttdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	fmt.Printf("SIGNAL %d received, so server shutting down now...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d.GracefulShutdown(ctx, &server)
}
