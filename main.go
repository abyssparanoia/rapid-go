package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abyssparanoia/rapid-go/src/lib/log"
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

	// set default logger
	handler := log.NewLogger(r)

	//server
	server := http.Server{
		Addr:    addr,
		Handler: handler,
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

	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("failed to gracefully shutdown: %s\n", err)
	}

	fmt.Printf("server shutdown completed\n")

}
