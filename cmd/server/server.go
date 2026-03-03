package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	appconfig "github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config/serverConfig"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/observability/metrics"
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
func startHTTPServer(config *CLIConfig, appConfig *appconfig.AppConfig, mcpServer *server.MCPServer, sc *serverConfig.ServerConfig) (*http.Server, error) {
	logrus.Infof("Starting cloud-native-mcp-server on %s (mode=%s)", config.Addr, config.Mode)
	if sc == nil {
		return nil, fmt.Errorf("server config is nil")
	}

	var sseServers map[string]*server.SSEServer
	var streamableHTTPServers map[string]*server.StreamableHTTPServer

	// Only initialize servers based on mode
	switch config.Mode {
	case "sse":
		sseServers = sc.InitSSEServers(mcpServer, config.Addr, appConfig)
	case "streamable-http":
		streamableHTTPServers = sc.InitStreamableHTTPServers(mcpServer, config.Addr, appConfig)
	default:
		return nil, fmt.Errorf("unsupported mode %q (supported: sse, streamable-http)", config.Mode)
	}

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, sseServers, streamableHTTPServers, config.Mode, appConfig, mcpServer)

	// Register metrics endpoint with optional authentication
	metricsHandler := metrics.Handler()
	if appConfig != nil && appConfig.Auth.Enabled {
		authConfig := middleware.AuthConfig{
			Enabled:             true,
			Mode:                appConfig.Auth.Mode,
			APIKey:              appConfig.Auth.APIKey,
			BearerToken:         appConfig.Auth.BearerToken,
			Username:            appConfig.Auth.Username,
			Password:            appConfig.Auth.Password,
			OIDCIssuerURL:       appConfig.Auth.OIDCIssuerURL,
			OIDCDiscoveryURL:    appConfig.Auth.OIDCDiscoveryURL,
			OIDCIssuer:          appConfig.Auth.OIDCIssuer,
			OIDCAudience:        appConfig.Auth.OIDCAudience,
			OIDCClientID:        appConfig.Auth.OIDCClientID,
			OIDCHTTPTimeoutSec:  appConfig.Auth.OIDCHTTPTimeoutSec,
			OIDCJWKSCacheTTLSec: appConfig.Auth.OIDCJWKSCacheTTLSec,
		}
		metricsHandler = middleware.AuthMiddleware(authConfig)(metricsHandler)
		logrus.Info("Metrics endpoint protected with authentication")
	} else {
		logrus.Warn("Metrics endpoint is not protected - consider enabling authentication for production")
	}
	mux.Handle("/metrics", metricsHandler)

	// Wrap mux with metrics middleware
	handler := middleware.MetricsMiddleware("cloud-native-mcp-server")(mux)

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
func initMetrics(buildInfo BuildInfo, mode, addr string) {
	metrics.Init(buildInfo.Version, buildInfo.Commit, buildInfo.GoVersion, mode, addr)
}
