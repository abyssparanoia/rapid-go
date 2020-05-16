package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"

	env "github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

func main() {

	// initilize environment variables
	e := &environment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger, err := log.New(e.Envrionment)
	if err != nil {
		panic(err)
	}

	server := newDefaultServer(logger, e)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for s := range c {
		logger.Info("Signal recieved", zap.String("signal", s.String()))

		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			shutdown := make(chan bool, 1)

			go func() {

				server.GracefulStop()

				shutdown <- true
			}()

			select {
			case <-shutdown:
				logger.Info("Gracefully stop")
			case <-ctx.Done():
				logger.Info("Force stop")
				server.Stop()
			}

			cancel()

			logger.Info("Exit")
			return

		case syscall.SIGHUP:
		default:
			return

		}
	}

}
