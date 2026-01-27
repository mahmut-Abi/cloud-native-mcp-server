package serverConfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware/hook"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/manager"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prompts"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/openapi"
)

const (
	serverName      = "kubernetes-mcp-server"
	serverVersion   = "0.0.1"
	shutdownTimeout = 3 * time.Second
)

type ServerConfig struct {
	disabledTools  map[string]bool
	serviceManager *manager.Manager
	auditStorage   middleware.AuditStorage
	allowedOrigins []string
	corsMaxAge     int
}

func (s *ServerConfig) InitHooks() *server.Hooks {
	logrus.WithFields(logrus.Fields{
		"component": "server",
		"operation": "init_hooks",
	}).Debug("Initializing server hooks")
	hooks := &server.Hooks{}
	hooks.AddOnRegisterSession(hook.SessionRegisterHookFunc())
	hooks.AddBeforeCallTool(hook.LogRequestHookFunc())
	hooks.AddAfterCallTool(hook.LogResponseHookFunc())
	logrus.WithFields(logrus.Fields{
		"component": "server",
		"operation": "init_hooks",
		"status":    "success",
	}).Debug("Server hooks initialized successfully")
	return hooks
}

func (s *ServerConfig) InitMCPServer(hooks *server.Hooks) *server.MCPServer {
	logrus.WithFields(logrus.Fields{
		"component": "server",
		"operation": "init_mcp_server",
	}).Debug("Initializing MCP server with tool, prompt, and resource capabilities")
	mcpServer := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithHooks(hooks),
		server.WithLogging(),
		server.WithRecovery(),
	)
	logrus.WithFields(logrus.Fields{
		"component": "server",
		"operation": "init_mcp_server",
		"status":    "success",
	}).Debug("MCP server initialized successfully")
	return mcpServer
}

func (s *ServerConfig) InitializeServices(appConfig *config.AppConfig) error {
	// Load CORS configuration
	if appConfig != nil {
		s.allowedOrigins = appConfig.Server.CORS.AllowedOrigins
		if len(s.allowedOrigins) == 0 {
			// Default to empty list (deny all origins) for security
			s.allowedOrigins = []string{}
			logrus.Warn("No CORS origins configured, CORS will deny all requests. Configure allowedOrigins explicitly if needed.")
		}
		s.corsMaxAge = appConfig.Server.CORS.MaxAge
		if s.corsMaxAge == 0 {
			s.corsMaxAge = 86400 // Default 24 hours
		}
	} else {
		// Default to empty list (deny all origins) for security
		s.allowedOrigins = []string{}
		s.corsMaxAge = 86400
		logrus.Warn("No app config provided, CORS will deny all requests. Configure allowedOrigins explicitly if needed.")
	}

	s.serviceManager = manager.NewManager()
	return s.serviceManager.Initialize(appConfig)
}

func (s *ServerConfig) AddToolsToServer(mcpServer *server.MCPServer) {
	if s.serviceManager == nil {
		logrus.Error("Service manager not initialized")
		return
	}

	logrus.Debug("Adding tools to MCP server via service manager")
	s.serviceManager.RegisterToolsAndHandlers(mcpServer)

	logrus.Debug("Adding resources to MCP server via service manager")
	s.serviceManager.RegisterResourcesAndHandlers(mcpServer)

	enabledServices := s.serviceManager.GetEnabledServices()
	logrus.Infof("Successfully registered %d enabled services", len(enabledServices))

	// Verify tool registration
	isValid, issues := s.serviceManager.VerifyToolRegistration()
	if !isValid {
		for _, issue := range issues {
			logrus.Errorf("Tool Registration Issue: %s", issue)
		}
	} else {
		logrus.Info("Tool registration verification passed - all tools have handlers")
	}

	// Log detailed registration report
	report := s.serviceManager.GetRegistrationReport()
	logrus.Debugf("Tool Registration Report: %+v", report)
}

func (s *ServerConfig) AddPromptsToServer(mcpServer *server.MCPServer) {
	logrus.WithFields(logrus.Fields{
		"component": "server",
		"operation": "add_prompts",
	}).Debug("Adding prompts to MCP server")
	mcpServer.AddPrompt(prompts.TestPodPrompt(), prompts.HandleTestPrompt)
	mcpServer.AddPrompt(prompts.K8sOpsPrompt(), prompts.HandleK8sOpsPrompt)
	logrus.Debug("Prompts added to MCP server successfully")
}

func (s *ServerConfig) InitSSEServers(mcpServer *server.MCPServer, addr string, appConfig *config.AppConfig) map[string]*server.SSEServer {
	logrus.Debug("Initializing multiple SSE servers for different services")

	sseServers := make(map[string]*server.SSEServer)

	// Default paths if not configured
	kubernetesPath := "/api/kubernetes/sse"
	grafanaPath := "/api/grafana/sse"
	prometheusPath := "/api/prometheus/sse"
	kibanaPath := "/api/kibana/sse"
	helmPath := "/api/helm/sse"
	elasticsearchPath := "/api/elasticsearch/sse"
	alertmanagerPath := "/api/alertmanager/sse"
	jaegerPath := "/api/jaeger/sse"
	aggregatePath := "/api/aggregate/sse"
	utilitiesPath := "/api/utilities/sse"

	// Override with config if provided
	if appConfig != nil {
		if appConfig.Server.SSEPaths.Kubernetes != "" {
			kubernetesPath = appConfig.Server.SSEPaths.Kubernetes
		}
		if appConfig.Server.SSEPaths.Grafana != "" {
			grafanaPath = appConfig.Server.SSEPaths.Grafana
		}
		if appConfig.Server.SSEPaths.Prometheus != "" {
			prometheusPath = appConfig.Server.SSEPaths.Prometheus
		}
		if appConfig.Server.SSEPaths.Kibana != "" {
			kibanaPath = appConfig.Server.SSEPaths.Kibana
		}
		if appConfig.Server.SSEPaths.Helm != "" {
			helmPath = appConfig.Server.SSEPaths.Helm
		}
		if appConfig.Server.SSEPaths.Elasticsearch != "" {
			elasticsearchPath = appConfig.Server.SSEPaths.Elasticsearch
		}
		if appConfig.Server.SSEPaths.Alertmanager != "" {
			alertmanagerPath = appConfig.Server.SSEPaths.Alertmanager
		}
		if appConfig.Server.SSEPaths.Jaeger != "" {
			jaegerPath = appConfig.Server.SSEPaths.Jaeger
		}
		// Check for aggregate SSE path in config
		if appConfig.Server.SSEPaths.Aggregate != "" {
			aggregatePath = appConfig.Server.SSEPaths.Aggregate
		}
		// Check for utilities SSE path in config
		if appConfig.Server.SSEPaths.Utilities != "" {
			utilitiesPath = appConfig.Server.SSEPaths.Utilities
		}
	}

	// Create service-specific MCP servers
	kubernetesServer := s.createServiceMCPServer("kubernetes", mcpServer)
	grafanaServer := s.createServiceMCPServer("grafana", mcpServer)
	prometheusServer := s.createServiceMCPServer("prometheus", mcpServer)
	kibanaServer := s.createServiceMCPServer("kibana", mcpServer)
	helmServer := s.createServiceMCPServer("helm", mcpServer)
	elasticsearchServer := s.createServiceMCPServer("elasticsearch", mcpServer)
	alertmanagerServer := s.createServiceMCPServer("alertmanager", mcpServer)
	jaegerServer := s.createServiceMCPServer("jaeger", mcpServer)
	opentelemetryServer := s.createServiceMCPServer("opentelemetry", mcpServer)

	// Create aggregated MCP server with all services
	aggregateServer := s.createAggregateMCPServer(mcpServer, kubernetesServer, grafanaServer, prometheusServer, kibanaServer, helmServer, elasticsearchServer, alertmanagerServer, jaegerServer, opentelemetryServer)

	// Create SSE servers for each service
	sseServers["kubernetes"] = server.NewSSEServer(kubernetesServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(kubernetesPath),
		server.WithMessageEndpoint(kubernetesPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	sseServers["grafana"] = server.NewSSEServer(grafanaServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(grafanaPath),
		server.WithMessageEndpoint(grafanaPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	sseServers["prometheus"] = server.NewSSEServer(prometheusServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(prometheusPath),
		server.WithMessageEndpoint(prometheusPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	sseServers["kibana"] = server.NewSSEServer(kibanaServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(kibanaPath),
		server.WithMessageEndpoint(kibanaPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)
	sseServers["helm"] = server.NewSSEServer(helmServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(helmPath),
		server.WithMessageEndpoint(helmPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	// Create aggregated SSE server
	sseServers["aggregate"] = server.NewSSEServer(aggregateServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(aggregatePath),
		server.WithMessageEndpoint(aggregatePath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	// Create elasticsearch SSE server
	sseServers["elasticsearch"] = server.NewSSEServer(elasticsearchServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(elasticsearchPath),
		server.WithMessageEndpoint(elasticsearchPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	// Create alertmanager SSE server
	sseServers["alertmanager"] = server.NewSSEServer(alertmanagerServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(alertmanagerPath),
		server.WithMessageEndpoint(alertmanagerPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	// Create jaeger SSE server
	sseServers["jaeger"] = server.NewSSEServer(jaegerServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(jaegerPath),
		server.WithMessageEndpoint(jaegerPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	// Create opentelemetry SSE server
	opentelemetryPath := appConfig.Server.SSEPaths.OpenTelemetry
	if opentelemetryPath == "" {
		opentelemetryPath = "/api/opentelemetry/sse"
	}
	sseServers["opentelemetry"] = server.NewSSEServer(opentelemetryServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(opentelemetryPath),
		server.WithMessageEndpoint(opentelemetryPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	// Create utilities server and SSE server
	utilitiesServer := s.createServiceMCPServer("utilities", mcpServer)
	sseServers["utilities"] = server.NewSSEServer(utilitiesServer,
		server.WithStaticBasePath(""),
		server.WithSSEEndpoint(utilitiesPath),
		server.WithMessageEndpoint(utilitiesPath+"/message"),
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	logrus.WithFields(logrus.Fields{
		"kubernetes_path":    kubernetesPath,
		"grafana_path":       grafanaPath,
		"prometheus_path":    prometheusPath,
		"kibana_path":        kibanaPath,
		"helm_path":          helmPath,
		"elasticsearch_path": elasticsearchPath,
		"alertmanager_path":  alertmanagerPath,
		"jaeger_path":        jaegerPath,
		"aggregate_path":     aggregatePath,
		"utilities_path":     utilitiesPath,
	}).Info("Multiple SSE servers initialized successfully")

	return sseServers
}

func (s *ServerConfig) InitStreamableHTTPServers(mcpServer *server.MCPServer, addr string, appConfig *config.AppConfig) map[string]*server.StreamableHTTPServer {
	logrus.Debug("Initializing multiple StreamableHTTP servers for different services")

	streamableHTTPServers := make(map[string]*server.StreamableHTTPServer)

	// Default paths if not configured
	kubernetesPath := "/api/kubernetes/streamable-http"
	grafanaPath := "/api/grafana/streamable-http"
	prometheusPath := "/api/prometheus/streamable-http"
	kibanaPath := "/api/kibana/streamable-http"
	helmPath := "/api/helm/streamable-http"
	elasticsearchPath := "/api/elasticsearch/streamable-http"
	alertmanagerPath := "/api/alertmanager/streamable-http"
	jaegerPath := "/api/jaeger/streamable-http"
	aggregatePath := "/api/aggregate/streamable-http"
	utilitiesPath := "/api/utilities/streamable-http"

	// Override with config if provided
	if appConfig != nil {
		if appConfig.Server.StreamableHTTPPaths.Kubernetes != "" {
			kubernetesPath = appConfig.Server.StreamableHTTPPaths.Kubernetes
		}
		if appConfig.Server.StreamableHTTPPaths.Grafana != "" {
			grafanaPath = appConfig.Server.StreamableHTTPPaths.Grafana
		}
		if appConfig.Server.StreamableHTTPPaths.Prometheus != "" {
			prometheusPath = appConfig.Server.StreamableHTTPPaths.Prometheus
		}
		if appConfig.Server.StreamableHTTPPaths.Kibana != "" {
			kibanaPath = appConfig.Server.StreamableHTTPPaths.Kibana
		}
		if appConfig.Server.StreamableHTTPPaths.Helm != "" {
			helmPath = appConfig.Server.StreamableHTTPPaths.Helm
		}
		if appConfig.Server.StreamableHTTPPaths.Elasticsearch != "" {
			elasticsearchPath = appConfig.Server.StreamableHTTPPaths.Elasticsearch
		}
		if appConfig.Server.StreamableHTTPPaths.Alertmanager != "" {
			alertmanagerPath = appConfig.Server.StreamableHTTPPaths.Alertmanager
		}
		if appConfig.Server.StreamableHTTPPaths.Jaeger != "" {
			jaegerPath = appConfig.Server.StreamableHTTPPaths.Jaeger
		}
		if appConfig.Server.StreamableHTTPPaths.Aggregate != "" {
			aggregatePath = appConfig.Server.StreamableHTTPPaths.Aggregate
		}
		if appConfig.Server.StreamableHTTPPaths.Utilities != "" {
			utilitiesPath = appConfig.Server.StreamableHTTPPaths.Utilities
		}
	}

	// Create service-specific MCP servers
	kubernetesServer := s.createServiceMCPServer("kubernetes", mcpServer)
	grafanaServer := s.createServiceMCPServer("grafana", mcpServer)
	prometheusServer := s.createServiceMCPServer("prometheus", mcpServer)
	kibanaServer := s.createServiceMCPServer("kibana", mcpServer)
	helmServer := s.createServiceMCPServer("helm", mcpServer)
	elasticsearchServer := s.createServiceMCPServer("elasticsearch", mcpServer)
	alertmanagerServer := s.createServiceMCPServer("alertmanager", mcpServer)
	jaegerServer := s.createServiceMCPServer("jaeger", mcpServer)
	opentelemetryServer := s.createServiceMCPServer("opentelemetry", mcpServer)

	// Create aggregated MCP server with all services
	aggregateServer := s.createAggregateMCPServer(mcpServer, kubernetesServer, grafanaServer, prometheusServer, kibanaServer, helmServer, elasticsearchServer, alertmanagerServer, jaegerServer, opentelemetryServer)

	// Create utilities server
	utilitiesServer := s.createServiceMCPServer("utilities", mcpServer)

	// Create StreamableHTTPServer for each service
	streamableHTTPServers["kubernetes"] = server.NewStreamableHTTPServer(kubernetesServer,
		server.WithEndpointPath(kubernetesPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["grafana"] = server.NewStreamableHTTPServer(grafanaServer,
		server.WithEndpointPath(grafanaPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["prometheus"] = server.NewStreamableHTTPServer(prometheusServer,
		server.WithEndpointPath(prometheusPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["kibana"] = server.NewStreamableHTTPServer(kibanaServer,
		server.WithEndpointPath(kibanaPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["helm"] = server.NewStreamableHTTPServer(helmServer,
		server.WithEndpointPath(helmPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["aggregate"] = server.NewStreamableHTTPServer(aggregateServer,
		server.WithEndpointPath(aggregatePath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["elasticsearch"] = server.NewStreamableHTTPServer(elasticsearchServer,
		server.WithEndpointPath(elasticsearchPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["alertmanager"] = server.NewStreamableHTTPServer(alertmanagerServer,
		server.WithEndpointPath(alertmanagerPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["jaeger"] = server.NewStreamableHTTPServer(jaegerServer,
		server.WithEndpointPath(jaegerPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	// Create opentelemetry StreamableHTTP server
	opentelemetryPath := appConfig.Server.StreamableHTTPPaths.OpenTelemetry
	if opentelemetryPath == "" {
		opentelemetryPath = "/api/opentelemetry/streamable-http"
	}
	streamableHTTPServers["opentelemetry"] = server.NewStreamableHTTPServer(opentelemetryServer,
		server.WithEndpointPath(opentelemetryPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	streamableHTTPServers["utilities"] = server.NewStreamableHTTPServer(utilitiesServer,
		server.WithEndpointPath(utilitiesPath),
		server.WithHeartbeatInterval(60*time.Second),
		server.WithStateLess(true),
	)

	logrus.WithFields(logrus.Fields{
		"kubernetes_path":    kubernetesPath,
		"grafana_path":       grafanaPath,
		"prometheus_path":    prometheusPath,
		"kibana_path":        kibanaPath,
		"helm_path":          helmPath,
		"elasticsearch_path": elasticsearchPath,
		"alertmanager_path":  alertmanagerPath,
		"jaeger_path":        jaegerPath,
		"opentelemetry_path": opentelemetryPath,
		"aggregate_path":     aggregatePath,
		"utilities_path":     utilitiesPath,
	}).Info("Multiple StreamableHTTP servers initialized successfully")

	return streamableHTTPServers
}

func (s *ServerConfig) createServiceMCPServer(serviceName string, baseMcpServer *server.MCPServer) *server.MCPServer {
	logrus.Debugf("Creating service-specific MCP server for: %s", serviceName)

	// Create a new MCP server for the specific service
	serviceServer := server.NewMCPServer(
		fmt.Sprintf("kubernetes-mcp-server-%s", serviceName),
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add service-specific tools and handlers
	if s.serviceManager != nil {
		switch serviceName {
		case "kubernetes":
			if kubernetesService := s.serviceManager.GetKubernetesService(); kubernetesService != nil && kubernetesService.IsEnabled() {
				tools := kubernetesService.GetTools()
				handlers := kubernetesService.GetHandlers()
				for _, tool := range tools {
					if handler, exists := handlers[tool.Name]; exists {
						serviceServer.AddTool(tool, handler)
					}
				}
				// Add prompts for Kubernetes service
				serviceServer.AddPrompt(prompts.TestPodPrompt(), prompts.HandleTestPrompt)
				serviceServer.AddPrompt(prompts.K8sOpsPrompt(), prompts.HandleK8sOpsPrompt)
			}
		case "grafana":
			if grafanaService := s.serviceManager.GetGrafanaService(); grafanaService != nil && grafanaService.IsEnabled() {
				tools := grafanaService.GetTools()
				handlers := grafanaService.GetHandlers()
				for _, tool := range tools {
					if handler, exists := handlers[tool.Name]; exists {
						serviceServer.AddTool(tool, handler)
					}
				}
			}
		case "prometheus":
			if prometheusService := s.serviceManager.GetPrometheusService(); prometheusService != nil && prometheusService.IsEnabled() {
				tools := prometheusService.GetTools()
				handlers := prometheusService.GetHandlers()
				for _, tool := range tools {
					if handler, exists := handlers[tool.Name]; exists {
						serviceServer.AddTool(tool, handler)
					}
				}
				// Add resources for Prometheus service
				resources := prometheusService.GetResources()
				resourceHandlers := prometheusService.GetResourceHandlers()
				for _, resource := range resources {
					if handler, exists := resourceHandlers[resource.URI]; exists {
						serviceServer.AddResource(resource, handler)
					}
				}
			}
		case "kibana":
			if kibanaService := s.serviceManager.GetKibanaService(); kibanaService != nil && kibanaService.IsEnabled() {
				tools := kibanaService.GetTools()
				handlers := kibanaService.GetHandlers()
				for _, tool := range tools {
					if handler, exists := handlers[tool.Name]; exists {
						serviceServer.AddTool(tool, handler)
					}
				}
			}
		case "elasticsearch":
			if elasticsearchService := s.serviceManager.GetElasticsearchService(); elasticsearchService != nil && elasticsearchService.IsEnabled() {
				tools := elasticsearchService.GetTools()
				handlers := elasticsearchService.GetHandlers()
				for _, tool := range tools {
					if handler, exists := handlers[tool.Name]; exists {
						serviceServer.AddTool(tool, handler)
					}
				}
			}
		case "alertmanager":
			if alertmanagerService := s.serviceManager.GetAlertmanagerService(); alertmanagerService != nil && alertmanagerService.IsEnabled() {
				tools := alertmanagerService.GetTools()
				handlers := alertmanagerService.GetHandlers()
				for _, tool := range tools {
					if handler, exists := handlers[tool.Name]; exists {
						serviceServer.AddTool(tool, handler)
					}
				}
			}
		case "jaeger":
			if jaegerService := s.serviceManager.GetJaegerService(); jaegerService != nil && jaegerService.IsEnabled() {
				tools := jaegerService.GetTools()
				handlers := jaegerService.GetHandlers()
				for _, tool := range tools {
					if handler, exists := handlers[tool.Name]; exists {
						serviceServer.AddTool(tool, handler)
					}
				}
			}
		}
	}

	return serviceServer
}

// createAggregateMCPServer creates an MCP server that aggregates all service capabilities
func (s *ServerConfig) createAggregateMCPServer(baseMcpServer *server.MCPServer, kubernetesServer, grafanaServer, prometheusServer, kibanaServer, helmServer, elasticsearchServer, alertmanagerServer, jaegerServer, opentelemetryServer *server.MCPServer) *server.MCPServer {
	logrus.Debug("Creating aggregated MCP server with all services")

	// Create a new MCP server for aggregated services
	aggregateServer := server.NewMCPServer(
		"kubernetes-mcp-server-aggregate",
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add all tools, prompts, and resources from all services
	if s.serviceManager != nil {
		// Add Kubernetes service capabilities
		if kubernetesService := s.serviceManager.GetKubernetesService(); kubernetesService != nil && kubernetesService.IsEnabled() {
			tools := kubernetesService.GetTools()
			handlers := kubernetesService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
			// Add prompts for Kubernetes service
			aggregateServer.AddPrompt(prompts.TestPodPrompt(), prompts.HandleTestPrompt)
			aggregateServer.AddPrompt(prompts.K8sOpsPrompt(), prompts.HandleK8sOpsPrompt)
		}

		// Add Grafana service capabilities
		if grafanaService := s.serviceManager.GetGrafanaService(); grafanaService != nil && grafanaService.IsEnabled() {
			tools := grafanaService.GetTools()
			handlers := grafanaService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
		}

		// Add Prometheus service capabilities
		if prometheusService := s.serviceManager.GetPrometheusService(); prometheusService != nil && prometheusService.IsEnabled() {
			tools := prometheusService.GetTools()
			handlers := prometheusService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
			// Add resources for Prometheus service
			resources := prometheusService.GetResources()
			resourceHandlers := prometheusService.GetResourceHandlers()
			for _, resource := range resources {
				if handler, exists := resourceHandlers[resource.URI]; exists {
					aggregateServer.AddResource(resource, handler)
				}
			}
		}

		// Add Kibana service capabilities
		if kibanaService := s.serviceManager.GetKibanaService(); kibanaService != nil && kibanaService.IsEnabled() {
			tools := kibanaService.GetTools()
			handlers := kibanaService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
		}

		// Add Helm service capabilities
		if helmService := s.serviceManager.GetHelmService(); helmService != nil && helmService.IsEnabled() {
			tools := helmService.GetTools()
			handlers := helmService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
			// Add additional tools for Helm service
			additionalTools := helmService.GetAdditionalTools()
			additionalHandlers := helmService.GetAdditionalHandlers()
			for _, tool := range additionalTools {
				if handler, exists := additionalHandlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
		}

		// Add Utilities service capabilities
		if utilitiesService := s.serviceManager.GetUtilitiesService(); utilitiesService != nil && utilitiesService.IsEnabled() {
			tools := utilitiesService.GetTools()
			handlers := utilitiesService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
		}

		// Add Elasticsearch service capabilities
		if elasticsearchService := s.serviceManager.GetElasticsearchService(); elasticsearchService != nil && elasticsearchService.IsEnabled() {
			tools := elasticsearchService.GetTools()
			handlers := elasticsearchService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
		}

		// Add Alertmanager service capabilities
		if alertmanagerService := s.serviceManager.GetAlertmanagerService(); alertmanagerService != nil && alertmanagerService.IsEnabled() {
			tools := alertmanagerService.GetTools()
			handlers := alertmanagerService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
		}

		// Add Jaeger service capabilities
		if jaegerService := s.serviceManager.GetJaegerService(); jaegerService != nil && jaegerService.IsEnabled() {
			tools := jaegerService.GetTools()
			handlers := jaegerService.GetHandlers()
			for _, tool := range tools {
				if handler, exists := handlers[tool.Name]; exists {
					aggregateServer.AddTool(tool, handler)
				}
			}
		}
	}

	return aggregateServer
}

func (s *ServerConfig) swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	// Serve Swagger UI static content
	swaggerHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *,
        *:before,
        *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"> </script>
<script>
    window.onload = function() {
        const ui = SwaggerUIBundle({
            url: '/api/openapi.json',
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.presets.standalone
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        });
        window.ui = ui;
    };
</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(swaggerHTML))
	if err != nil {
		logrus.Errorf("Failed to serve Swagger UI: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *ServerConfig) openAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Generate OpenAPI specification
	generator := openapi.NewGenerator(s.serviceManager.GetServiceRegistry())
	spec, err := generator.Generate()
	if err != nil {
		logrus.Errorf("Failed to generate OpenAPI spec: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert to JSON
	data, err := json.Marshal(spec)
	if err != nil {
		logrus.Errorf("Failed to marshal OpenAPI spec: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		logrus.Errorf("Failed to write OpenAPI spec: %v", err)
	}
}

// Efficient client IP extraction with caching in context
func getClientIP(r *http.Request) string {
	// Check if we already extracted the IP and stored it in the request context
	if ip, ok := r.Context().Value("client_ip").(string); ok && ip != "" {
		return ip
	}

	// Extract the IP
	var ip string

	// Check X-Forwarded-For header - most common for proxied requests
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For may contain multiple IPs, take the first one
		parts := strings.SplitN(forwarded, ",", 2)
		ip = strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}

	// Check X-Real-IP header - used by some proxies
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	addr := r.RemoteAddr
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		return host
	}

	// If we couldn't split the host:port, return the RemoteAddr as-is
	return addr
}

// Precomputed JSON response for health check handler
var healthCheckResponse = []byte(`{"status":"healthy"}`)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)

	// Only log in debug mode
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.WithField("client_ip", clientIP).Debug("Health check endpoint called")
	}

	// Set headers and write response in one go
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	// Use precomputed response instead of JSON encoder
	_, err := w.Write(healthCheckResponse)

	if err != nil {
		logrus.WithField("client_ip", clientIP).Errorf("Failed to write health check response: %v", err)
	} else if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.WithField("client_ip", clientIP).Debug("Health check response sent successfully")
	}
}

func (s *ServerConfig) StartServer(server *http.Server, addr string) {
	logrus.Infof("🚀 MCP Server starting on %s", addr)
	logrus.Debugf("Server configuration: %+v", server)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("Server error: %v", err)
	}
}

func (s *ServerConfig) WaitForShutdown(server *http.Server) {
	logrus.Debug("Setting up signal handlers for graceful shutdown")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logrus.Infof("Received signal %v, shutting down server...", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	logrus.Debugf("Starting graceful shutdown with timeout: %v", shutdownTimeout)
	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("Server shutdown error: %v", err)
	} else {
		logrus.Info("Server shutdown completed")
	}
}

func (s *ServerConfig) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		// Only log CORS processing in debug mode to reduce log volume
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			logrus.WithField("client_ip", clientIP).Debugf("CORS middleware processing request: %s %s", r.Method, r.URL.Path)
		}

		// Set CORS headers - using Header() only once for performance
		headers := w.Header()

		// Use configured CORS origins or default to all origins for development
		allowedOrigins := []string{"*"}
		if s != nil && len(s.allowedOrigins) > 0 {
			allowedOrigins = s.allowedOrigins
		}

		// Set Access-Control-Allow-Origin header
		if len(allowedOrigins) > 0 {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// No origin header, allow all
				headers.Set("Access-Control-Allow-Origin", allowedOrigins[0])
			} else {
				// Check if origin is in allowed list
				allowed := false
				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == "*" || allowedOrigin == origin {
						allowed = true
						break
					}
				}
				if allowed {
					headers.Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		// Set allowed methods
		allowedMethods := []string{"HEAD", "GET", "POST", "OPTIONS"}
		headers.Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

		// Set allowed headers
		allowedHeaders := []string{"Content-Type", "Authorization", "X-Requested-With"}
		headers.Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

		// Set max age
		maxAge := 86400 // 24 hours
		if s != nil && s.corsMaxAge > 0 {
			maxAge = s.corsMaxAge
		}
		headers.Set("Access-Control-Max-Age", strconv.Itoa(maxAge))

		// Fast path for OPTIONS requests
		if r.Method == "OPTIONS" {
			if logrus.IsLevelEnabled(logrus.DebugLevel) {
				logrus.WithField("client_ip", clientIP).Debug("Handling OPTIONS preflight request")
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
	buffer bytes.Buffer
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size

	// Capture response body for error responses (up to 4KB)
	if rw.status >= 400 && rw.buffer.Len() < 4096 {
		remaining := 4096 - rw.buffer.Len()
		if len(b) <= remaining {
			rw.buffer.Write(b)
		} else {
			rw.buffer.Write(b[:remaining])
		}
	}

	return size, err
}

func (rw *responseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		clientIP := getClientIP(r)

		// Avoid expensive debug logging unless enabled
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			logrus.WithField("client_ip", clientIP).Debugf("Request started: %s %s", r.Method, r.URL.Path)
		}

		// Create a response writer wrapper to capture status code and body
		rw := &responseWriter{ResponseWriter: w, status: 0}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Log at appropriate level based on status code
		fields := logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      rw.status,
			"size":        rw.size,
			"duration_ms": duration.Milliseconds(),
			"client_ip":   clientIP,
		}

		// Log based on status code
		if rw.status >= 500 {
			logrus.WithFields(fields).Error("Request failed")
		} else if rw.status >= 400 {
			// For 4xx errors, try to log the response body for debugging
			if rw.buffer.Len() > 0 && rw.buffer.Len() <= 4096 {
				fields["response_body"] = rw.buffer.String()
			}
			logrus.WithFields(fields).Warn("Request error")
		} else {
			logrus.WithFields(fields).Info("Request processed")
		}

		// Detailed timing info only at debug level
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			logrus.WithField("client_ip", clientIP).Debugf("Request completed in %v", duration)
		}
	})
}

// ApplyServiceFilters applies service enable/disable filters
func (s *ServerConfig) ApplyServiceFilters(disabledServices, enabledServices, disabledTools string) error {
	if s.serviceManager == nil {
		return fmt.Errorf("service manager not initialized")
	}

	// Parse service lists
	disabledSvcs := parseList(disabledServices)
	enabledSvcs := parseList(enabledServices)
	disabledToolList := parseList(disabledTools)

	// If specific services are enabled, disable all others
	allServices := []string{"kubernetes", "grafana", "prometheus", "kibana", "helm", "elasticsearch", "utilities"}
	if len(enabledSvcs) > 0 {
		for _, svc := range allServices {
			if !enabledSvcs[svc] {
				disabledSvcs[svc] = true
			}
		}
	}

	// Apply service disablement
	// Convert maps to slices for manager.ApplyServiceFilters
	disabledSvcList := make([]string, 0, len(disabledSvcs))
	for svc := range disabledSvcs {
		disabledSvcList = append(disabledSvcList, svc)
	}
	enabledSvcList := make([]string, 0, len(enabledSvcs))
	for svc := range enabledSvcs {
		enabledSvcList = append(enabledSvcList, svc)
	}
	s.serviceManager.ApplyServiceFilters(disabledSvcList, enabledSvcList)

	// Apply tool disablement
	s.disabledTools = disabledToolList

	return nil
}

// parseList parses a comma-separated list into a map
func parseList(list string) map[string]bool {
	result := make(map[string]bool)
	if list == "" {
		return result
	}
	items := strings.Split(list, ",")
	for _, item := range items {
		result[strings.TrimSpace(item)] = true
	}
	return result
}

// GetEnabledServices returns all enabled services
func (s *ServerConfig) GetEnabledServices() map[string]services.Service {
	if s.serviceManager == nil {
		return make(map[string]services.Service)
	}
	return s.serviceManager.GetEnabledServices()
}

// GetAllTools returns all tools from enabled services
func (s *ServerConfig) GetAllTools() []mcp.Tool {
	if s.serviceManager == nil {
		return []mcp.Tool{}
	}
	return s.serviceManager.GetAllTools()
}

// GetAllToolsIncludingDisabled returns all tools including from disabled services
func (s *ServerConfig) GetAllToolsIncludingDisabled() []mcp.Tool {
	if s.serviceManager == nil {
		return nil
	}
	allServices := s.serviceManager.GetServiceRegistry().GetAllServices()
	var tools []mcp.Tool
	for _, service := range allServices {
		if svcTools := service.GetTools(); svcTools != nil {
			tools = append(tools, svcTools...)
		}
	}
	return tools
}

// SetupMultipleRoutes sets up HTTP routes for the server
func (s *ServerConfig) SetupMultipleRoutes(mux *http.ServeMux, sseServers map[string]*server.SSEServer, streamableHTTPServers map[string]*server.StreamableHTTPServer, mode string, appConfig *config.AppConfig, baseMcpServer *server.MCPServer) {
	logrus.Debug("Setting up HTTP routes for SSE and StreamableHTTP servers")

	// Initialize audit storage if not already done
	if s.auditStorage == nil && appConfig != nil && appConfig.Audit.Enabled {
		var err error
		s.auditStorage, err = middleware.CreateAuditStorage(appConfig)
		if err != nil {
			s.auditStorage = nil
			logrus.WithFields(logrus.Fields{
				"component": "server",
				"operation": "init_audit_storage",
				"error":     err,
			}).Error("Failed to create audit storage")
		} else {
			logrus.WithFields(logrus.Fields{
				"component": "server",
				"operation": "init_audit_storage",
				"storage":   appConfig.Audit.Storage,
				"format":    appConfig.Audit.Format,
				"max_logs":  appConfig.Audit.MaxLogs,
			}).Debug("Audit storage initialized")
		}
	}

	mux.HandleFunc("/health", healthCheckHandler)

	// Add OpenAPI documentation endpoints
	mux.HandleFunc("/api/openapi.json", s.openAPIHandler)
	mux.HandleFunc("/api/docs", s.swaggerUIHandler)

	// Add audit log endpoints if audit is enabled
	if appConfig != nil && appConfig.Audit.Enabled && s.auditStorage != nil {
		mux.HandleFunc("/api/audit/logs", s.auditLogsHandler())
		mux.HandleFunc("/api/audit/stats", s.auditStatsHandler())
		logrus.Debug("Audit endpoints registered")
	}

	// Default paths if not configured
	kubernetesPath := "/api/kubernetes/sse"
	grafanaPath := "/api/grafana/sse"
	prometheusPath := "/api/prometheus/sse"
	kibanaPath := "/api/kibana/sse"
	helmPath := "/api/helm/sse"
	aggregatePath := "/api/aggregate/sse"

	// Override with config if provided
	if appConfig != nil {
		if appConfig.Server.SSEPaths.Kubernetes != "" {
			kubernetesPath = appConfig.Server.SSEPaths.Kubernetes
		}
		if appConfig.Server.SSEPaths.Grafana != "" {
			grafanaPath = appConfig.Server.SSEPaths.Grafana
		}
		if appConfig.Server.SSEPaths.Prometheus != "" {
			prometheusPath = appConfig.Server.SSEPaths.Prometheus
		}
		if appConfig.Server.SSEPaths.Kibana != "" {
			kibanaPath = appConfig.Server.SSEPaths.Kibana
		}
		if appConfig.Server.SSEPaths.Helm != "" {
			helmPath = appConfig.Server.SSEPaths.Helm
		}
		if appConfig.Server.SSEPaths.Aggregate != "" {
			aggregatePath = appConfig.Server.SSEPaths.Aggregate
		}
	}

	if mode == "sse" || mode == "http" {
		// Setup routes for each service
		for serviceName, sseServer := range sseServers {
			logrus.Debugf("Setting up routes for service: %s", serviceName)

			// Set up endpoints based on service name and configuration
			var sseEndpoint, messageEndpoint string
			switch serviceName {
			case "kubernetes":
				sseEndpoint = kubernetesPath
			case "grafana":
				sseEndpoint = grafanaPath
			case "prometheus":
				sseEndpoint = prometheusPath
			case "kibana":
				sseEndpoint = kibanaPath
			case "helm":
				sseEndpoint = helmPath
			case "aggregate":
				sseEndpoint = aggregatePath
			default:
				sseEndpoint = "/api/" + serviceName + "/sse"
			}
			messageEndpoint = sseEndpoint + "/message"

			// Set up message handler
			messageHandler := sseServer.MessageHandler()

			// Apply audit middleware if enabled
			if appConfig != nil && appConfig.Audit.Enabled && s.auditStorage != nil {
				auditConfig := middleware.AuditMiddlewareConfig{
					Enabled: true,
					Storage: s.auditStorage,
				}
				messageHandler = middleware.AuditMiddleware(auditConfig)(messageHandler)
			}

			// Apply authentication middleware if enabled
			if appConfig != nil && appConfig.Auth.Enabled {
				authConfig := middleware.AuthConfig{
					Enabled:     appConfig.Auth.Enabled,
					Mode:        appConfig.Auth.Mode,
					APIKey:      appConfig.Auth.APIKey,
					BearerToken: appConfig.Auth.BearerToken,
					Username:    appConfig.Auth.Username,
					Password:    appConfig.Auth.Password,
				}
				messageHandler = middleware.AuthMiddleware(authConfig)(messageHandler)
			}

			// Apply CORS and logging middleware
			messageHandler = s.corsMiddleware(loggingMiddleware(messageHandler))

			// Apply security middleware
			securityConfig := middleware.DefaultSecurityConfig()
			messageHandler = middleware.SecurityMiddleware(securityConfig)(messageHandler)

			mux.Handle(messageEndpoint, messageHandler)

			if mode == "sse" {
				// Set up SSE handler
				sseHandler := sseServer.SSEHandler()

				// Apply audit middleware if enabled
				if appConfig != nil && appConfig.Audit.Enabled && s.auditStorage != nil {
					auditConfig := middleware.AuditMiddlewareConfig{
						Enabled: true,
						Storage: s.auditStorage,
					}
					sseHandler = middleware.AuditMiddleware(auditConfig)(sseHandler)
				}

				// Apply authentication middleware if enabled
				if appConfig != nil && appConfig.Auth.Enabled {
					authConfig := middleware.AuthConfig{
						Enabled:     appConfig.Auth.Enabled,
						Mode:        appConfig.Auth.Mode,
						APIKey:      appConfig.Auth.APIKey,
						BearerToken: appConfig.Auth.BearerToken,
						Username:    appConfig.Auth.Username,
						Password:    appConfig.Auth.Password,
					}
					sseHandler = middleware.AuthMiddleware(authConfig)(sseHandler)
				}

				// Apply CORS and logging middleware
				sseHandler = s.corsMiddleware(loggingMiddleware(sseHandler))

				// Apply security middleware
				securityConfig := middleware.DefaultSecurityConfig()
				sseHandler = middleware.SecurityMiddleware(securityConfig)(sseHandler)

				mux.Handle(sseEndpoint, sseHandler)
			}

			if mode == "sse" {
				logrus.WithFields(logrus.Fields{
					"service":          serviceName,
					"sse_endpoint":     sseEndpoint,
					"message_endpoint": messageEndpoint,
				}).Debug("SSE routes configured for service")
			}
		}
	}

	// Setup routes for StreamableHTTP mode
	if mode == "streamable-http" && streamableHTTPServers != nil {
		logrus.Debug("Setting up StreamableHTTP routes")

		// Default paths if not configured
		kubernetesPath := "/api/kubernetes/streamable-http"
		grafanaPath := "/api/grafana/streamable-http"
		prometheusPath := "/api/prometheus/streamable-http"
		kibanaPath := "/api/kibana/streamable-http"
		helmPath := "/api/helm/streamable-http"
		aggregatePath := "/api/aggregate/streamable-http"
		utilitiesPath := "/api/utilities/streamable-http"

		// Override with config if provided
		if appConfig != nil {
			if appConfig.Server.StreamableHTTPPaths.Kubernetes != "" {
				kubernetesPath = appConfig.Server.StreamableHTTPPaths.Kubernetes
			}
			if appConfig.Server.StreamableHTTPPaths.Grafana != "" {
				grafanaPath = appConfig.Server.StreamableHTTPPaths.Grafana
			}
			if appConfig.Server.StreamableHTTPPaths.Prometheus != "" {
				prometheusPath = appConfig.Server.StreamableHTTPPaths.Prometheus
			}
			if appConfig.Server.StreamableHTTPPaths.Kibana != "" {
				kibanaPath = appConfig.Server.StreamableHTTPPaths.Kibana
			}
			if appConfig.Server.StreamableHTTPPaths.Helm != "" {
				helmPath = appConfig.Server.StreamableHTTPPaths.Helm
			}
			if appConfig.Server.StreamableHTTPPaths.Aggregate != "" {
				aggregatePath = appConfig.Server.StreamableHTTPPaths.Aggregate
			}
			if appConfig.Server.StreamableHTTPPaths.Utilities != "" {
				utilitiesPath = appConfig.Server.StreamableHTTPPaths.Utilities
			}
		}

		for serviceName, httpServer := range streamableHTTPServers {
			logrus.Debugf("Setting up StreamableHTTP routes for service: %s", serviceName)

			var httpPath string
			switch serviceName {
			case "kubernetes":
				httpPath = kubernetesPath
			case "grafana":
				httpPath = grafanaPath
			case "prometheus":
				httpPath = prometheusPath
			case "kibana":
				httpPath = kibanaPath
			case "helm":
				httpPath = helmPath
			case "aggregate":
				httpPath = aggregatePath
			case "utilities":
				httpPath = utilitiesPath
			default:
				httpPath = "/api/" + serviceName + "/streamable-http"
			}

			// Create handler with middleware
			var httpHandler http.Handler = httpServer

			// Apply audit middleware if enabled
			if appConfig != nil && appConfig.Audit.Enabled && s.auditStorage != nil {
				auditConfig := middleware.AuditMiddlewareConfig{
					Enabled: true,
					Storage: s.auditStorage,
				}
				httpHandler = middleware.AuditMiddleware(auditConfig)(httpHandler)
			}

			// Apply authentication middleware if enabled
			if appConfig != nil && appConfig.Auth.Enabled {
				authConfig := middleware.AuthConfig{
					Enabled:     appConfig.Auth.Enabled,
					Mode:        appConfig.Auth.Mode,
					APIKey:      appConfig.Auth.APIKey,
					BearerToken: appConfig.Auth.BearerToken,
					Username:    appConfig.Auth.Username,
					Password:    appConfig.Auth.Password,
				}
				httpHandler = middleware.AuthMiddleware(authConfig)(httpHandler)
			}

			// Apply CORS and logging middleware
			httpHandler = s.corsMiddleware(loggingMiddleware(httpHandler))

			// Apply security middleware
			securityConfig := middleware.DefaultSecurityConfig()
			httpHandler = middleware.SecurityMiddleware(securityConfig)(httpHandler)

			mux.Handle(httpPath, httpHandler)

			logrus.WithFields(logrus.Fields{
				"service":       serviceName,
				"http_endpoint": httpPath,
			}).Debug("StreamableHTTP route configured for service")
		}
	}

	logrus.WithField("mode", mode).Debug("HTTP routes setup completed for all services")
}

// auditLogsHandler returns handler for retrieving audit logs
func (s *ServerConfig) auditLogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.auditStorage == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "{\"error\":\"audit storage not initialized\"}")
			return
		}

		// Build query criteria from URL parameters
		criteria := make(map[string]interface{})
		if userID := r.URL.Query().Get("user_id"); userID != "" {
			criteria["user_id"] = userID
		}
		if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
			criteria["service_name"] = serviceName
		}
		if toolName := r.URL.Query().Get("tool_name"); toolName != "" {
			criteria["tool_name"] = toolName
		}
		if status := r.URL.Query().Get("status"); status != "" {
			criteria["status"] = status
		}

		logs, err := s.auditStorage.Query(criteria)
		if err != nil {
			logrus.WithError(err).Error("Failed to query audit logs")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "{\"error\":\"failed to query logs\"}")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if logs == nil {
			logs = []middleware.AuditLogEntry{}
		}
		_ = json.NewEncoder(w).Encode(logs)
	}
}

// auditStatsHandler returns handler for retrieving audit statistics
func (s *ServerConfig) auditStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.auditStorage == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "{\"error\":\"audit storage not initialized\"}")
			return
		}

		// Parse time range from query parameters
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		startTime := time.Now().Add(-24 * time.Hour)
		endTime := time.Now()

		if startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				startTime = t
			}
		}

		if endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				endTime = t
			}
		}

		stats, err := s.auditStorage.GetStats(startTime, endTime)
		if err != nil {
			logrus.WithError(err).Error("Failed to get audit stats")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "{\"error\":\"failed to get stats\"}")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(stats)
	}
}

// Shutdown gracefully shuts down the server and all resources
func (s *ServerConfig) Shutdown() error {
	logrus.Info("Shutting down server configuration...")
	var errs []error

	// Close audit storage
	if s.auditStorage != nil {
		if err := s.auditStorage.Close(); err != nil {
			errs = append(errs, fmt.Errorf("audit storage close error: %w", err))
		}
	}

	// Shutdown service manager
	if s.serviceManager != nil {
		if err := s.serviceManager.Shutdown(); err != nil {
			errs = append(errs, fmt.Errorf("service manager shutdown error: %w", err))
		}
	}

	if len(errs) > 0 {
		logrus.Warnf("Shutdown completed with %d errors", len(errs))
		for _, err := range errs {
			logrus.Warnf("  - %v", err)
		}
		return fmt.Errorf("shutdown completed with %d errors", len(errs))
	}

	logrus.Info("Server configuration shutdown completed")
	return nil
}
