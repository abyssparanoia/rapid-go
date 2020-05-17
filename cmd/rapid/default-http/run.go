package defaulthttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"go.uber.org/zap"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"
)

func run() {

	e := &environment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger, err := log.New(e.Envrionment)
	if err != nil {
		panic(err)
	}

	// Dependency
	d := Dependency{}
	d.Inject(e, logger)

	// Routing
	r := chi.NewRouter()
	Routing(r, d)

	addr := fmt.Sprintf(":%s", e.Port)

	//server
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Run
	logger.Info(fmt.Sprintf("[START] server. port: %s\n", addr))
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("[CLOSED] server closed with error", zap.Error(err))
		}
	}()

	// graceful shuttdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	logger.Info(fmt.Sprintf("SIGNAL %d received, so server shutting down now...\n", <-quit))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error("failed to gracefully shutdown", zap.Error(err))
	}

	logger.Info("server shutdown completed")
}
