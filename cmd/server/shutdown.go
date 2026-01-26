package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// gracefulShutdown handles graceful server shutdown with resource cleanup
func gracefulShutdown(server *http.Server, sc interface{}) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Increase shutdown timeout to 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clean up service resources if ServerConfig is provided
	if sc != nil {
		logrus.Info("Cleaning up service resources...")
		if serverConfig, ok := sc.(interface{ Shutdown() error }); ok {
			if err := serverConfig.Shutdown(); err != nil {
				logrus.WithError(err).Error("Failed to shutdown server configuration")
			}
		}
	}

	// Then shut down the server
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
		return
	}

	logrus.Info("Server exited cleanly")
}
