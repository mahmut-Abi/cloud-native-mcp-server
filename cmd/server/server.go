package main

import (
	"context"
	"net"
	"net/http"
	"runtime"

	appconfig "github.com/mahmut-Abi/k8s-mcp-server/internal/config"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/config/serverConfig"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/middleware"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/observability/metrics"
	"github.com/sirupsen/logrus"

	server "github.com/mark3labs/mcp-go/server"
)

// createServer creates and configures the HTTP server
func createServer(config *CLIConfig, handler http.Handler) *http.Server {
	logrus.WithFields(logrus.Fields{"addr": config.Addr, "read": config.ReadTimeout, "write": config.WriteTimeout, "idle": config.IdleTimeout}).Debug("Creating HTTP server")

	return &http.Server{
		Addr:           config.Addr,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
		BaseContext: func(_ net.Listener) context.Context {
			return context.Background()
		},
	}
}

// startHTTPServer starts the HTTP server
func startHTTPServer(config *CLIConfig, appConfig *appconfig.AppConfig, mcpServer *server.MCPServer) (*http.Server, error) {
	logrus.Infof("Starting k8s-mcp-server on %s (mode=%s)", config.Addr, config.Mode)

	var sseServers map[string]*server.SSEServer
	var streamableHTTPServers map[string]*server.StreamableHTTPServer

	// Only initialize servers based on mode
	switch config.Mode {
	case "sse":
		sseServers = initSSEServers(mcpServer, config.Addr, appConfig)
	case "streamable-http":
		streamableHTTPServers = initStreamableHTTPServers(mcpServer, config.Addr, appConfig)
	}

	mux := http.NewServeMux()
	setupMultipleRoutes(mux, sseServers, streamableHTTPServers, config.Mode, appConfig, mcpServer)

	// Register metrics endpoint with optional authentication
	metricsHandler := metrics.Handler()
	if appConfig != nil && appConfig.Auth.Enabled {
		authConfig := middleware.AuthConfig{
			Enabled:     true,
			Mode:        appConfig.Auth.Mode,
			APIKey:      appConfig.Auth.APIKey,
			BearerToken: appConfig.Auth.BearerToken,
			Username:    appConfig.Auth.Username,
			Password:    appConfig.Auth.Password,
		}
		metricsHandler = middleware.AuthMiddleware(authConfig)(metricsHandler)
		logrus.Info("Metrics endpoint protected with authentication")
	} else {
		logrus.Warn("Metrics endpoint is not protected - consider enabling authentication for production")
	}
	mux.Handle("/metrics", metricsHandler)

	// Wrap mux with metrics middleware
	handler := middleware.MetricsMiddleware("k8s-mcp-server")(mux)

	srv := createServer(config, handler)

	// Start server in goroutine
	go func() {
		logrus.Infof("Server listening on %s", config.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	return srv, nil
}

// initSSEServers initializes SSE servers
func initSSEServers(mcpServer *server.MCPServer, addr string, appConfig *appconfig.AppConfig) map[string]*server.SSEServer {
	var sc serverConfig.ServerConfig
	return sc.InitSSEServers(mcpServer, addr, appConfig)
}

// initStreamableHTTPServers initializes streamable HTTP servers
func initStreamableHTTPServers(mcpServer *server.MCPServer, addr string, appConfig *appconfig.AppConfig) map[string]*server.StreamableHTTPServer {
	var sc serverConfig.ServerConfig
	return sc.InitStreamableHTTPServers(mcpServer, addr, appConfig)
}

// setupMultipleRoutes sets up multiple routes for the server
func setupMultipleRoutes(mux *http.ServeMux, sseServers map[string]*server.SSEServer, streamableHTTPServers map[string]*server.StreamableHTTPServer, mode string, appConfig *appconfig.AppConfig, mcpServer *server.MCPServer) {
	var sc serverConfig.ServerConfig
	sc.SetupMultipleRoutes(mux, sseServers, streamableHTTPServers, mode, appConfig, mcpServer)
}

// initMCPServer initializes the MCP server
func initMCPServer(appConfig *appconfig.AppConfig) (*server.MCPServer, *serverConfig.ServerConfig) {
	var sc serverConfig.ServerConfig
	mcpHooks := sc.InitHooks()
	mcpServer := sc.InitMCPServer(mcpHooks)

	// Initialize services
	if err := sc.InitializeServices(appConfig); err != nil {
		logrus.Fatalf("Failed to initialize services: %v", err)
	}

	sc.AddToolsToServer(mcpServer)
	sc.AddPromptsToServer(mcpServer)

	return mcpServer, &sc
}

// initMetrics initializes the metrics system
func initMetrics(mode, addr string) {
	gitCommit := "unknown"
	if commit, err := runGitCommand("rev-parse", "HEAD"); err == nil {
		gitCommit = commit
	}
	metrics.Init("dev", gitCommit, getRuntimeVersion(), mode, addr)
}

// getRuntimeVersion returns the Go runtime version
func getRuntimeVersion() string {
	return runtime.Version()
}
