package framework

import (
	"fmt"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
)

// ServiceInitializer provides common service initialization patterns
type ServiceInitializer struct {
	serviceName string
}

// NewServiceInitializer creates a new service initializer for a specific service
func NewServiceInitializer(serviceName string) *ServiceInitializer {
	return &ServiceInitializer{serviceName: serviceName}
}

// InitConfig holds the common initialization configuration
type InitConfig struct {
	// Required indicates if the service must be enabled
	Required bool
	// URLValidator is a function to validate service URL
	URLValidator func(string) bool
	// ClientBuilder is a function to build the HTTP client using the config
	ClientBuilder func(*config.AppConfig) (interface{}, error)
}

// ServiceEnableChecker determines if a service should be enabled
type ServiceEnableChecker interface {
	IsEnabled(config interface{}) bool
	GetURL(config interface{}) string
}

// CommonServiceInit provides a reusable initialization pattern for services
type CommonServiceInit struct {
	initializer *ServiceInitializer
	config      *InitConfig
	checker     ServiceEnableChecker
}

// NewCommonServiceInit creates a new common service initializer
func NewCommonServiceInit(serviceName string, config *InitConfig, checker ServiceEnableChecker) *CommonServiceInit {
	return &CommonServiceInit{
		initializer: NewServiceInitializer(serviceName),
		config:      config,
		checker:     checker,
	}
}

// Initialize performs the common initialization steps
func (csi *CommonServiceInit) Initialize(cfg interface{}, setEnabled func(bool), setClient func(interface{})) error {
	// Step 1: Validate and cast config
	appConfig, ok := cfg.(*config.AppConfig)
	if !ok || appConfig == nil {
		// Service remains disabled if no config provided
		setEnabled(false)
		return nil
	}

	// Step 2: Check if service is enabled in configuration
	if !csi.checker.IsEnabled(appConfig) {
		setEnabled(false)
		return nil
	}

	// Step 3: Validate required configuration
	url := csi.checker.GetURL(appConfig)
	if url == "" || (csi.config.URLValidator != nil && !csi.config.URLValidator(url)) {
		errMsg := fmt.Sprintf("%s URL is required but not provided or invalid", csi.initializer.serviceName)
		if !csi.config.Required {
			setEnabled(false)
			return nil
		}
		return fmt.Errorf("%s", errMsg)
	}

	// Step 4: Build client if needed
	if csi.config.ClientBuilder != nil {
		client, err := csi.config.ClientBuilder(appConfig)
		if err != nil {
			return fmt.Errorf("failed to create %s client: %w", csi.initializer.serviceName, err)
		}
		setClient(client)
	}

	// Step 5: Enable the service
	setEnabled(true)
	return nil
}

// SimpleURLValidator provides basic URL validation
func SimpleURLValidator(url string) bool {
	return url != ""
}

// ServiceEnabled provides a simple enabled checker implementation
type ServiceEnabled struct {
	getEnabled func(*config.AppConfig) bool
	getURL     func(*config.AppConfig) string
}

// NewServiceEnabled creates a new service enabled checker
func NewServiceEnabled(getEnabled func(*config.AppConfig) bool, getURL func(*config.AppConfig) string) *ServiceEnabled {
	return &ServiceEnabled{
		getEnabled: getEnabled,
		getURL:     getURL,
	}
}

// IsEnabled checks if service is enabled
func (se *ServiceEnabled) IsEnabled(cfg interface{}) bool {
	appConfig, ok := cfg.(*config.AppConfig)
	if !ok || appConfig == nil {
		return false
	}
	return se.getEnabled(appConfig)
}

// GetURL gets the service URL
func (se *ServiceEnabled) GetURL(cfg interface{}) string {
	appConfig, ok := cfg.(*config.AppConfig)
	if !ok || appConfig == nil {
		return ""
	}
	return se.getURL(appConfig)
}
