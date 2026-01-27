package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	appconfig "github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/logging"
)

// CLIConfig holds all command line configuration
type CLIConfig struct {
	Addr         string
	Kubeconfig   string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ConfigPath   string
	Mode         string
	ListMode     string // services, tools, or empty
	Format       string // text, json, table, csv
	ServiceName  string // filter by service name
	Verbose      bool   // verbose output
	// Flags to track which parameters were explicitly set
	addrSet         bool
	kubeconfigSet   bool
	logLevelSet     bool
	readTimeoutSet  bool
	writeTimeoutSet bool
	idleTimeoutSet  bool
	modeSet         bool
}

// parseFlags parses command line flags and returns configuration
func parseFlags() *CLIConfig {
	var (
		addr         string
		kubeconfig   string
		logLevel     string
		readTimeout  int
		writeTimeout int
		idleTimeout  int
		configPath   string
		mode         string
		help         bool
		listMode     string
		format       string
		serviceName  string
		verbose      bool
	)

	// Set default kubeconfig path
	defaultKubeconfig := getDefaultKubeconfig()

	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "address to listen on")
	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig, "path to kubeconfig file")
	flag.StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error, fatal)")
	flag.IntVar(&readTimeout, "read-timeout", 0, "HTTP server read timeout in seconds (0 disables timeout)")
	flag.IntVar(&writeTimeout, "write-timeout", 0, "HTTP server write timeout in seconds (0 disables timeout; recommended for SSE)")
	flag.IntVar(&idleTimeout, "idle-timeout", 60, "HTTP server idle timeout in seconds (default: 60, suitable for streaming)")
	flag.StringVar(&configPath, "config", os.Getenv("MCP_CONFIG"), "path to YAML config file (env MCP_CONFIG)")
	flag.StringVar(&mode, "mode", os.Getenv("MCP_MODE"), "run mode: sse | http | streamable-http | stdio")
	flag.StringVar(&listMode, "list", "", "list mode: services or tools")
	flag.StringVar(&format, "output", "text", "output format for list command: text, json, table, csv")
	flag.StringVar(&serviceName, "service", "", "filter by service name")
	flag.BoolVar(&verbose, "verbose", false, "verbose output for tools descriptions")
	flag.BoolVar(&help, "help", false, "show help message")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	cfg := &CLIConfig{
		Addr:         addr,
		Kubeconfig:   kubeconfig,
		LogLevel:     logLevel,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		ConfigPath:   configPath,
		Mode:         mode,
		ListMode:     listMode,
		Format:       format,
		ServiceName:  serviceName,
		Verbose:      verbose,
	}

	// Track which flags were explicitly set
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "addr":
			cfg.addrSet = true
		case "kubeconfig":
			cfg.kubeconfigSet = true
		case "log-level":
			cfg.logLevelSet = true
		case "read-timeout":
			cfg.readTimeoutSet = true
		case "write-timeout":
			cfg.writeTimeoutSet = true
		case "idle-timeout":
			cfg.idleTimeoutSet = true
		case "mode":
			cfg.modeSet = true
		}
	})

	if cfg.WriteTimeout > 0 && cfg.WriteTimeout < time.Minute {
		logrus.WithField("writeTimeout", cfg.WriteTimeout).Warn("SSE may timeout with low write-timeout; set to 0 to disable or >=1m")
	}

	return cfg
}

// applyAppConfig merges AppConfig into CLI config. CLI flags take precedence.
func applyAppConfig(c *CLIConfig, ac *appconfig.AppConfig) {
	if ac == nil {
		return
	}

	logrus.Debug("Applying configuration from file")

	// Server - only apply if CLI flag was not explicitly set
	if !c.addrSet && ac.Server.Addr != "" {
		logrus.Debugf("Setting address from config: %s", ac.Server.Addr)
		c.Addr = ac.Server.Addr
	}
	if !c.readTimeoutSet && ac.Server.ReadTimeoutSec > 0 {
		logrus.Debugf("Setting read timeout from config: %ds", ac.Server.ReadTimeoutSec)
		c.ReadTimeout = time.Duration(ac.Server.ReadTimeoutSec) * time.Second
	}
	if !c.writeTimeoutSet && ac.Server.WriteTimeoutSec > 0 {
		logrus.Debugf("Setting write timeout from config: %ds", ac.Server.WriteTimeoutSec)
		c.WriteTimeout = time.Duration(ac.Server.WriteTimeoutSec) * time.Second
	}
	if !c.idleTimeoutSet && ac.Server.IdleTimeoutSec > 0 {
		logrus.Debugf("Setting idle timeout from config: %ds", ac.Server.IdleTimeoutSec)
		c.IdleTimeout = time.Duration(ac.Server.IdleTimeoutSec) * time.Second
	}
	if !c.logLevelSet && ac.Logging.Level != "" {
		logrus.Debugf("Setting log level from config: %s", ac.Logging.Level)
		c.LogLevel = ac.Logging.Level
	}
	if !c.modeSet && ac.Server.Mode != "" {
		logrus.Debugf("Setting mode from config: %s", ac.Server.Mode)
		c.Mode = ac.Server.Mode
	}
	// Kubeconfig - only apply if CLI flag was not explicitly set
	if !c.kubeconfigSet && ac.Kubernetes.Kubeconfig != "" {
		logrus.Debugf("Setting kubeconfig from config: %s", ac.Kubernetes.Kubeconfig)
		c.Kubeconfig = ac.Kubernetes.Kubeconfig
	}
	// Logging JSON format
	if ac.Logging.JSON {
		logrus.Debug("Enabling JSON logging format")
		logging.EnableJSONFormat()
	}
}

// getDefaultKubeconfig returns the default kubeconfig path
func getDefaultKubeconfig() string {
	if envConfig := os.Getenv("KUBECONFIG"); envConfig != "" {
		return envConfig
	}
	return filepath.Join(os.Getenv("HOME"), ".kube", "config")
}

// setupLogging initializes logging with the specified level
func setupLogging(logLevel string) {
	logging.InitStdoutLogger()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Warnf("Invalid log level %s, using info level", logLevel)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.WithField("level", level.String()).Debug("Logrus initialized")
}
