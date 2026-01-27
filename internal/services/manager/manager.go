package manager

import (
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/config"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/alertmanager"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/elasticsearch"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/grafana"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/helm"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/jaeger"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/kibana"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/kubernetes"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/opentelemetry"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/prometheus"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/utilities"
)

var logger = logrus.WithField("component", "manager")

// Manager manages all services and provides a unified interface
type Manager struct {
	registry             *services.Registry
	kubernetesService    *kubernetes.Service
	grafanaService       *grafana.Service
	prometheusService    *prometheus.Service
	kibanaService        *kibana.Service
	helmService          *helm.Service
	alertmanagerService  *alertmanager.Service
	elasticsearchService *elasticsearch.Service
	jaegerService        *jaeger.Service
	opentelemetryService *opentelemetry.Service
	utilitiesService     *utilities.Service
	disabledTools        map[string]bool
	disabledToolsMutex   sync.RWMutex     // Protect disabledTools from concurrent access
	serviceStatus        map[string]bool  // Track service initialization status
	serviceErrors        map[string]error // Track service initialization errors
	statusMutex          sync.RWMutex     // Protect serviceStatus and serviceErrors from concurrent access
}

// NewManager creates a new service manager
func NewManager() *Manager {
	return &Manager{
		registry:      services.NewRegistry(),
		disabledTools: make(map[string]bool),
		serviceStatus: make(map[string]bool),
		serviceErrors: make(map[string]error),
	}
}

// Initialize initializes all services with the configuration using parallel initialization
func (m *Manager) Initialize(appConfig *config.AppConfig) error {
	// Load disabled tools from configuration
	if appConfig != nil {
		m.disabledToolsMutex.Lock()
		for _, toolName := range appConfig.EnableDisable.DisabledTools {
			m.disabledTools[toolName] = true
		}
		m.disabledToolsMutex.Unlock()
		if len(m.disabledTools) > 0 {
			logger.Infof("Loaded %d disabled tools from configuration", len(m.disabledTools))
		}
	}

	// Apply service filters from configuration
	if appConfig != nil {
		m.ApplyServiceFilters(appConfig.EnableDisable.DisabledServices, appConfig.EnableDisable.EnabledServices)
	}

	// Create services
	m.kubernetesService = kubernetes.NewService()
	m.grafanaService = grafana.NewService()
	m.prometheusService = prometheus.NewService()
	m.kibanaService = kibana.NewService()
	m.helmService = helm.NewService()
	m.alertmanagerService = alertmanager.NewService()
	m.elasticsearchService = elasticsearch.NewService()
	m.jaegerService = jaeger.NewService()
	m.opentelemetryService = opentelemetry.NewService()
	m.utilitiesService = utilities.NewService()

	// Register services
	m.registry.Register(m.kubernetesService)
	m.registry.Register(m.grafanaService)
	m.registry.Register(m.prometheusService)
	m.registry.Register(m.kibanaService)
	m.registry.Register(m.helmService)
	m.registry.Register(m.alertmanagerService)
	m.registry.Register(m.elasticsearchService)
	m.registry.Register(m.jaegerService)
	m.registry.Register(m.opentelemetryService)
	m.registry.Register(m.utilitiesService)

	// Define empty config as default
	var cfg interface{} = appConfig
	if appConfig == nil {
		cfg = &config.AppConfig{}
	}

	// Initialize critical Kubernetes service (must be initialized first)
	if err := m.kubernetesService.Initialize(cfg); err != nil {
		logger.WithError(err).Error("Critical: Kubernetes service initialization failed")
		m.statusMutex.Lock()
		m.serviceStatus["kubernetes"] = false
		m.serviceErrors["kubernetes"] = err
		m.statusMutex.Unlock()
		return fmt.Errorf("failed to initialize kubernetes service: %w", err)
	}
	logger.Debug("Kubernetes service initialized successfully")
	m.statusMutex.Lock()
	m.serviceStatus["kubernetes"] = true
	m.statusMutex.Unlock()

	// Initialize optional services in parallel for faster startup
	var wg sync.WaitGroup
	optionalServices := []struct {
		name     string
		initFunc func() error
	}{
		{"grafana", func() error { return m.grafanaService.Initialize(cfg) }},
		{"prometheus", func() error { return m.prometheusService.Initialize(cfg) }},
		{"kibana", func() error { return m.kibanaService.Initialize(cfg) }},
		{"elasticsearch", func() error { return m.elasticsearchService.Initialize(cfg) }},
		{"helm", func() error { return m.helmService.Initialize(cfg) }},
		{"alertmanager", func() error { return m.alertmanagerService.Initialize(cfg) }},
		{"jaeger", func() error { return m.jaegerService.Initialize(cfg) }},
		{"opentelemetry", func() error { return m.opentelemetryService.Initialize(cfg) }},
		{"utilities", func() error { return m.utilitiesService.Initialize(cfg) }},
	}

	for _, svc := range optionalServices {
		wg.Add(1)
		go func(s struct {
			name     string
			initFunc func() error
		}) {
			defer wg.Done()
			m.initializeOptionalService(s.name, s.initFunc)
		}(svc)
	}
	wg.Wait()

	m.logRegisteredTools()
	m.LogServiceStatus()
	logger.Info("Service manager initialization completed")
	return nil
}

// initializeOptionalService initializes an optional service with error handling
func (m *Manager) initializeOptionalService(serviceName string, initFunc func() error) {
	svcLogger := logger.WithField("service", serviceName)
	if err := initFunc(); err != nil {
		svcLogger.WithError(err).Warn("Optional service initialization failed, continuing with degraded functionality")
		m.statusMutex.Lock()
		m.serviceStatus[serviceName] = false
		m.serviceErrors[serviceName] = err
		m.statusMutex.Unlock()
	} else {
		svcLogger.Debug("Service initialized successfully")
		m.statusMutex.Lock()
		m.serviceStatus[serviceName] = true
		m.statusMutex.Unlock()
	}
}

// InitializeParallel initializes services with parallel execution for optional services
func (m *Manager) InitializeParallel(appConfig *config.AppConfig) error {
	// Create all services
	m.kubernetesService = kubernetes.NewService()
	m.grafanaService = grafana.NewService()
	m.prometheusService = prometheus.NewService()
	m.kibanaService = kibana.NewService()
	m.helmService = helm.NewService()
	m.alertmanagerService = alertmanager.NewService()
	m.elasticsearchService = elasticsearch.NewService()
	m.jaegerService = jaeger.NewService()
	m.opentelemetryService = opentelemetry.NewService()
	m.utilitiesService = utilities.NewService()

	// Register services
	m.registry.Register(m.kubernetesService)
	m.registry.Register(m.grafanaService)
	m.registry.Register(m.prometheusService)
	m.registry.Register(m.kibanaService)
	m.registry.Register(m.helmService)
	m.registry.Register(m.alertmanagerService)
	m.registry.Register(m.elasticsearchService)
	m.registry.Register(m.jaegerService)
	m.registry.Register(m.opentelemetryService)
	m.registry.Register(m.utilitiesService)

	// Block on critical Kubernetes service (must be initialized first)
	if err := m.kubernetesService.Initialize(appConfig); err != nil {
		logger.WithError(err).Error("Critical: Kubernetes service initialization failed")
		return fmt.Errorf("failed to initialize kubernetes service: %w", err)
	}
	logger.Debug("Kubernetes service initialized successfully")

	// Parallelize optional services for faster startup
	var wg sync.WaitGroup
	optionalServices := []struct {
		name     string
		initFunc func() error
	}{
		{"grafana", func() error { return m.grafanaService.Initialize(appConfig) }},
		{"prometheus", func() error { return m.prometheusService.Initialize(appConfig) }},
		{"kibana", func() error { return m.kibanaService.Initialize(appConfig) }},
		{"elasticsearch", func() error { return m.elasticsearchService.Initialize(appConfig) }},
		{"helm", func() error { return m.helmService.Initialize(appConfig) }},
		{"alertmanager", func() error { return m.alertmanagerService.Initialize(appConfig) }},
		{"jaeger", func() error { return m.jaegerService.Initialize(appConfig) }},
		{"opentelemetry", func() error { return m.opentelemetryService.Initialize(appConfig) }},
		{"utilities", func() error { return m.utilitiesService.Initialize(appConfig) }},
	}

	for _, svc := range optionalServices {
		wg.Add(1)
		go func(s struct {
			name     string
			initFunc func() error
		}) {
			defer wg.Done()
			m.initializeOptionalService(s.name, s.initFunc)
		}(svc)
	}
	wg.Wait()

	logger.Info("Service manager initialization completed successfully")
	return nil
}

// RegisterToolsAndHandlers registers all tools and handlers with the MCP server
func (m *Manager) RegisterToolsAndHandlers(mcpServer *server.MCPServer) {
	// Get all enabled services
	services := m.registry.GetEnabledServices()

	// Track statistics
	totalTools := 0
	missingHandlers := 0
	missingHandlerNames := []string{}

	// Register tools and handlers for each service
	for _, service := range services {
		tools := service.GetTools()
		handlers := service.GetHandlers()

		for _, tool := range tools {
			totalTools++

			// Skip disabled tools
			m.disabledToolsMutex.RLock()
			disabled := m.disabledTools[tool.Name]
			m.disabledToolsMutex.RUnlock()

			if disabled {
				logger.Debugf("Skipping disabled tool: %s", tool.Name)
				continue
			}
			if handler, exists := handlers[tool.Name]; exists {
				mcpServer.AddTool(tool, handler)
			} else {
				logger.Errorf("Tool '%s' has no handler defined, skipping registration", tool.Name)
				missingHandlers++
				missingHandlerNames = append(missingHandlerNames, tool.Name)
			}
		}
	}

	// Report statistics
	logger.Infof("Registered %d tools from %d services", totalTools, len(services))
	if missingHandlers > 0 {
		logger.Warnf("Warning: %d tools are missing handlers and will not be available:", missingHandlers)
		for _, name := range missingHandlerNames {
			logger.Warnf("  - %s", name)
		}
	}
}

// RegisterResourcesAndHandlers registers all resources and handlers with the MCP server
func (m *Manager) RegisterResourcesAndHandlers(mcpServer *server.MCPServer) {
	// Get all enabled services
	services := m.registry.GetEnabledServices()

	// Register resources and handlers for each service
	for _, service := range services {
		// Check if service supports resources (using type assertion)
		if resourceService, ok := service.(interface {
			GetResources() []mcp.Resource
			GetResourceHandlers() map[string]server.ResourceHandlerFunc
		}); ok {
			resources := resourceService.GetResources()
			resourceHandlers := resourceService.GetResourceHandlers()

			for _, resource := range resources {
				if handler, exists := resourceHandlers[resource.URI]; exists {
					mcpServer.AddResource(resource, handler)
				}
			}
		}
	}
}

// GetEnabledServices returns all enabled services
func (m *Manager) GetEnabledServices() map[string]services.Service {
	return m.registry.GetEnabledServices()
}

// GetAllTools returns all tools from enabled services (including disabled tools)
func (m *Manager) GetAllTools() []mcp.Tool {
	return m.registry.GetAllTools()
}

// GetAllEnabledTools returns all enabled tools (excluding disabled tools)
func (m *Manager) GetAllEnabledTools() []mcp.Tool {
	allTools := m.registry.GetAllTools()
	var enabledTools []mcp.Tool
	for _, tool := range allTools {
		if !m.disabledTools[tool.Name] {
			enabledTools = append(enabledTools, tool)
		}
	}
	return enabledTools
}

// GetAllHandlers returns all handlers from enabled services
func (m *Manager) GetAllHandlers() map[string]server.ToolHandlerFunc {
	return m.registry.GetAllHandlers()
}

// GetKubernetesService returns the Kubernetes service
func (m *Manager) GetKubernetesService() *kubernetes.Service {
	return m.kubernetesService
}

// GetGrafanaService returns the Grafana service
func (m *Manager) GetGrafanaService() *grafana.Service {
	return m.grafanaService
}

// GetPrometheusService returns the Prometheus service
func (m *Manager) GetPrometheusService() *prometheus.Service {
	return m.prometheusService
}

// GetKibanaService returns the Kibana service
func (m *Manager) GetKibanaService() *kibana.Service {
	return m.kibanaService
}

// GetHelmService returns the Helm service
func (m *Manager) GetHelmService() *helm.Service {
	return m.helmService
}

// GetElasticsearchService returns the Elasticsearch service
func (m *Manager) GetElasticsearchService() *elasticsearch.Service {
	return m.elasticsearchService
}

// GetAlertmanagerService returns the Alertmanager service
func (m *Manager) GetAlertmanagerService() *alertmanager.Service {
	return m.alertmanagerService
}

// GetJaegerService returns the Jaeger service
func (m *Manager) GetJaegerService() *jaeger.Service {
	return m.jaegerService
}

// GetOpenTelemetryService returns the OpenTelemetry service
func (m *Manager) GetOpenTelemetryService() *opentelemetry.Service {
	return m.opentelemetryService
}

// GetUtilitiesService returns the Utilities service
func (m *Manager) GetUtilitiesService() *utilities.Service {
	return m.utilitiesService
}

// GetServiceStatus returns initialization status for a service
func (m *Manager) GetServiceStatus(serviceName string) bool {
	m.statusMutex.RLock()
	defer m.statusMutex.RUnlock()
	status, exists := m.serviceStatus[serviceName]
	if !exists {
		return false
	}
	return status
}

// GetServiceError returns initialization error for a service
func (m *Manager) GetServiceError(serviceName string) error {
	m.statusMutex.RLock()
	defer m.statusMutex.RUnlock()
	return m.serviceErrors[serviceName]
}

// GetAllServiceStatus returns all service statuses
func (m *Manager) GetAllServiceStatus() map[string]bool {
	m.statusMutex.RLock()
	defer m.statusMutex.RUnlock()
	result := make(map[string]bool, len(m.serviceStatus))
	for k, v := range m.serviceStatus {
		result[k] = v
	}
	return result
}

// GetAllServiceErrors returns all service errors
func (m *Manager) GetAllServiceErrors() map[string]error {
	m.statusMutex.RLock()
	defer m.statusMutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]error, len(m.serviceErrors))
	for k, v := range m.serviceErrors {
		result[k] = v
	}
	return result
}

// LogServiceStatus logs all service statuses
func (m *Manager) LogServiceStatus() {
	logger.Info("────────────────────────────────")
	logger.Info("Service Initialization Status")
	logger.Info("────────────────────────────────")

	successCount := 0
	failureCount := 0

	m.statusMutex.RLock()
	defer m.statusMutex.RUnlock()

	for serviceName, status := range m.serviceStatus {
		if status {
			logger.Infof("  ✓ %s: initialized", serviceName)
			successCount++
		} else {
			err := m.serviceErrors[serviceName]
			if err != nil {
				logger.Errorf("  ✗ %s: failed - %v", serviceName, err)
			} else {
				logger.Errorf("  ✗ %s: failed", serviceName)
			}
			failureCount++
		}
	}

	logger.Infof("Total: %d services initialized successfully, %d failed", successCount, failureCount)
	logger.Info("────────────────────────────────")
}

// logRegisteredTools logs registered tools
func (m *Manager) logRegisteredTools() {
	tools := m.GetAllEnabledTools()
	handlers := m.GetAllHandlers()
	enabledServices := m.registry.GetEnabledServices()
	logger.Infof("Tool Registration: %d services, %d tools (%d disabled)", len(enabledServices), len(tools), len(m.disabledTools))
	if len(m.disabledTools) > 0 {
		logger.Infof("Disabled tools: %v", m.disabledTools)
	}
	for _, tool := range tools {
		if _, exists := handlers[tool.Name]; !exists {
			logger.Warnf("Tool %s missing handler", tool.Name)
		}
	}
}

// VerifyToolRegistration checks if all tools have handlers
func (m *Manager) VerifyToolRegistration() (bool, []string) {
	var issues []string
	tools := m.GetAllTools()
	handlers := m.GetAllHandlers()
	if len(tools) == 0 {
		issues = append(issues, "No tools registered")
	}
	for _, tool := range tools {
		if _, exists := handlers[tool.Name]; !exists {
			issues = append(issues, fmt.Sprintf("Tool %s has no handler", tool.Name))
		}
	}
	return len(issues) == 0, issues
}

// GetRegistrationReport returns tool registration statistics
func (m *Manager) GetRegistrationReport() map[string]interface{} {
	tools := m.GetAllTools()
	handlers := m.GetAllHandlers()
	enabledServices := m.registry.GetEnabledServices()
	report := make(map[string]interface{})
	report["enabled_services"] = len(enabledServices)
	report["registered_tools"] = len(tools)
	var missing []string
	for _, tool := range tools {
		if _, exists := handlers[tool.Name]; !exists {
			missing = append(missing, tool.Name)
		}
	}
	report["missing_handlers"] = missing
	return report
}

// ApplyServiceFilters applies service enable/disable filters
func (m *Manager) ApplyServiceFilters(disabled, enabled []string) {
	// Convert string slices to maps for easier lookup
	disabledMap := make(map[string]bool)
	for _, svc := range disabled {
		disabledMap[svc] = true
	}

	enabledMap := make(map[string]bool)
	for _, svc := range enabled {
		enabledMap[svc] = true
	}

	allServices := []string{"kubernetes", "grafana", "prometheus", "kibana", "helm", "elasticsearch", "alertmanager", "jaeger", "opentelemetry", "utilities"}

	// If specific services are enabled, disable all others
	if len(enabled) > 0 {
		for _, svc := range allServices {
			if !enabledMap[svc] {
				disabledMap[svc] = true
			}
		}
	}

	// Apply service disablement
	if disabledMap["kubernetes"] && m.kubernetesService != nil {
		m.kubernetesService = nil
	}
	if disabledMap["grafana"] && m.grafanaService != nil {
		m.grafanaService = nil
	}
	if disabledMap["prometheus"] && m.prometheusService != nil {
		m.prometheusService = nil
	}
	if disabledMap["kibana"] && m.kibanaService != nil {
		m.kibanaService = nil
	}
	if disabledMap["helm"] && m.helmService != nil {
		m.helmService = nil
	}
	if disabledMap["elasticsearch"] && m.elasticsearchService != nil {
		m.elasticsearchService = nil
	}
	if disabledMap["alertmanager"] && m.alertmanagerService != nil {
		m.alertmanagerService = nil
	}
	if disabledMap["jaeger"] && m.jaegerService != nil {
		m.jaegerService = nil
	}
	if disabledMap["opentelemetry"] && m.opentelemetryService != nil {
		m.opentelemetryService = nil
	}
	if disabledMap["utilities"] && m.utilitiesService != nil {
		m.utilitiesService = nil
	}
}

// GetDisabledTools returns a list of disabled tools
func (m *Manager) GetDisabledTools() map[string]bool {
	m.disabledToolsMutex.RLock()
	defer m.disabledToolsMutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]bool)
	for k, v := range m.disabledTools {
		result[k] = v
	}
	return result
}

// DisableTool disables a specific tool
func (m *Manager) DisableTool(toolName string) {
	m.disabledToolsMutex.Lock()
	defer m.disabledToolsMutex.Unlock()
	m.disabledTools[toolName] = true
	logger.Infof("Tool disabled: %s", toolName)
}

// EnableTool enables a specific tool
func (m *Manager) EnableTool(toolName string) {
	m.disabledToolsMutex.Lock()
	defer m.disabledToolsMutex.Unlock()
	delete(m.disabledTools, toolName)
	logger.Infof("Tool enabled: %s", toolName)
}

// FilterTools filters tools based on disabled tools
func (m *Manager) FilterTools(tools []mcp.Tool, disabledTools map[string]bool) []mcp.Tool {
	var filtered []mcp.Tool
	for _, tool := range tools {
		if !disabledTools[tool.Name] {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// FilterHandlers filters handlers based on disabled tools
func (m *Manager) FilterHandlers(handlers map[string]server.ToolHandlerFunc, disabledTools map[string]bool) map[string]server.ToolHandlerFunc {
	filtered := make(map[string]server.ToolHandlerFunc)
	for name, handler := range handlers {
		if !disabledTools[name] {
			filtered[name] = handler
		}
	}
	return filtered
}

// LogEnabledServices logs all currently enabled services and tools
func (m *Manager) LogEnabledServices() {
	enabledServices := m.GetEnabledServices()
	all := m.GetAllTools()

	logger.Info("────────────────────────────────")
	logger.Info("Enabled Services and Tools")
	logger.Info("────────────────────────────────")
	logger.Infof("Total Enabled Services: %d", len(enabledServices))
	for svc := range enabledServices {
		logger.Infof("  ✓ %s", svc)
	}
	logger.Infof("Total Registered Tools: %d", len(all))
}

// GetServiceRegistry returns the service registry
func (m *Manager) GetServiceRegistry() *services.Registry {
	return m.registry
}

// Closer is an interface for services that can be closed
type Closer interface {
	Close() error
}

// Shutdown gracefully shuts down all services and releases resources
func (m *Manager) Shutdown() error {
	logger.Info("Shutting down service manager...")
	var errs []error

	// Collect all services that might need to be closed
	services := []struct {
		name    string
		service interface{}
	}{
		{"kubernetes", m.kubernetesService},
		{"grafana", m.grafanaService},
		{"prometheus", m.prometheusService},
		{"kibana", m.kibanaService},
		{"helm", m.helmService},
		{"alertmanager", m.alertmanagerService},
		{"elasticsearch", m.elasticsearchService},
		{"jaeger", m.jaegerService},
		{"opentelemetry", m.opentelemetryService},
		{"utilities", m.utilitiesService},
	}

	// Close services that implement the Closer interface
	for _, svc := range services {
		if svc.service == nil {
			continue
		}
		if closer, ok := svc.service.(Closer); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, fmt.Errorf("%s service close error: %w", svc.name, err))
			}
		}
	}

	if len(errs) > 0 {
		logger.Warnf("Shutdown completed with %d errors", len(errs))
		for _, err := range errs {
			logger.Warnf("  - %v", err)
		}
		return fmt.Errorf("shutdown completed with %d errors", len(errs))
	}

	logger.Info("Service manager shutdown completed successfully")
	return nil
}
