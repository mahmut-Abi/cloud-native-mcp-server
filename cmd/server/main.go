package main

import (
	"os"

	appconfig "github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/sirupsen/logrus"

	server "github.com/mark3labs/mcp-go/server"
)

func main() {
	// Set reasonable defaults for GOMAXPROCS
	setupGOMAXPROCS()

	config := parseFlags()

	// Load YAML/env config if provided (load once and save for later use)
	var appConfig *appconfig.AppConfig
	if config.ConfigPath != "" {
		logrus.Debugf("Loading configuration from: %s", config.ConfigPath)
		if ac, err := appconfig.Load(config.ConfigPath); err != nil {
			logrus.WithField("path", config.ConfigPath).WithError(err).Warn("Failed to load config file; proceeding with flags/env only")
			logrus.Info("Tip: Run with --help to see all available configuration options")
		} else {
			applyAppConfig(config, ac)
			appConfig = ac
			logrus.Info("Configuration loaded successfully from file")
		}
	} else {
		logrus.Debug("No config file specified, using flags and environment variables only")
	}

	// Set default mode if not specified
	if config.Mode == "" {
		config.Mode = "sse"
		logrus.Debug("No mode specified, defaulting to 'sse'")
	}

	setupLogging(config.LogLevel)

	// Initialize metrics system
	initMetrics(config.Mode, config.Addr)

	// Initialize MCP server
	mcpServer, sc := initMCPServer(appConfig)

	// Handle list mode if specified
	if config.ListMode != "" {
		opts := ListDisplayOptions{
			Format:        config.Format,
			ServiceFilter: config.ServiceName,
			Verbose:       config.Verbose,
		}

		switch config.ListMode {
		case "services":
			enabled := sc.GetEnabledServices()
			if err := DisplayServices(enabled, opts); err != nil {
				logrus.Fatalf("Failed to display services: %v", err)
			}
			os.Exit(0)
		case "tools":
			allTools := sc.GetAllToolsIncludingDisabled()
			if err := DisplayTools(allTools, opts); err != nil {
				logrus.Fatalf("Failed to display tools: %v", err)
			}
			os.Exit(0)
		default:
			logrus.Fatalf("invalid --list mode: %s (expected services|tools)", config.ListMode)
		}
	}

	if config.Mode == "stdio" {
		logrus.Info("Starting in stdio mode")
		if err := server.ServeStdio(mcpServer); err != nil {
			logrus.Fatalf("stdio server error: %v", err)
		}
		return
	}

	// Start HTTP server
	srv, err := startHTTPServer(config, appConfig, mcpServer)
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}

	// Wait for graceful shutdown
	gracefulShutdown(srv, sc)
}
