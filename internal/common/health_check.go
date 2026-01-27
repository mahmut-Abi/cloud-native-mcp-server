package common

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HealthCheckConfig represents configuration for health checks
type HealthCheckConfig struct {
	ServiceName    string            // Name of the service (for logging)
	HealthEndpoint string            // Health check endpoint (e.g., "health", "status", "api/health")
	HTTPClient     *http.Client      // HTTP client to use for requests
	Method         string            // HTTP method (default: GET)
	Timeout        time.Duration     // Timeout for health check (default: 10s)
	ExpectedStatus int               // Expected HTTP status code (default: 200)
	Headers        map[string]string // Additional headers
}

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Healthy   bool          // Whether the service is healthy
	Latency   time.Duration // Time taken for the health check
	Status    string        // Status message
	Error     error         // Error if unhealthy
	Timestamp time.Time     // When the health check was performed
}

// HealthChecker provides common health check functionality
type HealthChecker struct {
	serviceName string
	logger      *logrus.Entry
}

// NewHealthChecker creates a new health checker for a service
func NewHealthChecker(serviceName string) *HealthChecker {
	return &HealthChecker{
		serviceName: serviceName,
		logger:      logrus.WithField("component", serviceName+"-health"),
	}
}

// CheckHealth performs a health check on the service
func (h *HealthChecker) CheckHealth(ctx context.Context, baseURL string, config *HealthCheckConfig) *HealthCheckResult {
	if config == nil {
		config = &HealthCheckConfig{
			ServiceName:    h.serviceName,
			HealthEndpoint: "health",
			Method:         "GET",
			Timeout:        10 * time.Second,
			ExpectedStatus: 200,
		}
	}

	// Override with provided config values
	if config.ServiceName == "" {
		config.ServiceName = h.serviceName
	}
	if config.Method == "" {
		config.Method = "GET"
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if config.ExpectedStatus == 0 {
		config.ExpectedStatus = 200
	}

	startTime := time.Now()
	result := &HealthCheckResult{
		Timestamp: startTime,
	}

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Build request URL
	reqURL := baseURL
	if config.HealthEndpoint != "" {
		reqURL += "/" + config.HealthEndpoint
	}

	// Create request
	req, err := http.NewRequestWithContext(timeoutCtx, config.Method, reqURL, nil)
	if err != nil {
		result.Healthy = false
		result.Status = "failed to create request"
		result.Error = fmt.Errorf("failed to create health check request: %w", err)
		result.Latency = time.Since(startTime)
		return result
	}

	// Set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Use provided HTTP client or create a default one
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	// Execute request
	h.logger.WithFields(logrus.Fields{
		"url":     reqURL,
		"method":  config.Method,
		"timeout": config.Timeout,
	}).Debug("Performing health check")

	resp, err := httpClient.Do(req)
	if err != nil {
		result.Healthy = false
		result.Status = "connection failed"
		result.Error = fmt.Errorf("health check request failed: %w", err)
		result.Latency = time.Since(startTime)
		h.logger.WithError(err).Debug("Health check failed")
		return result
	}
	defer func() { _ = resp.Body.Close() }()

	result.Latency = time.Since(startTime)

	// Check status code
	if resp.StatusCode != config.ExpectedStatus {
		result.Healthy = false
		result.Status = "unexpected status code"
		result.Error = fmt.Errorf("health check returned status %d (expected %d)", resp.StatusCode, config.ExpectedStatus)
		h.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"expected":    config.ExpectedStatus,
			"latency":     result.Latency,
		}).Debug("Health check failed")
		return result
	}

	// Health check successful
	result.Healthy = true
	result.Status = "healthy"
	h.logger.WithFields(logrus.Fields{
		"latency": result.Latency,
	}).Debug("Health check successful")

	return result
}

// SimpleHealthCheck performs a simple health check with minimal configuration
func (h *HealthChecker) SimpleHealthCheck(ctx context.Context, httpClient *http.Client, baseURL, healthEndpoint string) error {
	result := h.CheckHealth(ctx, baseURL, &HealthCheckConfig{
		ServiceName:    h.serviceName,
		HealthEndpoint: healthEndpoint,
		HTTPClient:     httpClient,
		Method:         "GET",
		ExpectedStatus: 200,
	})

	if !result.Healthy {
		return fmt.Errorf("%s health check failed: %w", h.serviceName, result.Error)
	}

	return nil
}
