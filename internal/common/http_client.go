package common

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/constants"
	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
)

// AuthConfig represents authentication configuration for HTTP clients
type AuthConfig struct {
	Enabled     bool
	Mode        string // apikey, bearer, basic
	APIKey      string
	BearerToken string
	Username    string
	Password    string
}

// TLSConfig represents TLS configuration for HTTP clients
type TLSConfig struct {
	SkipVerify bool
	CertFile   string
	KeyFile    string
	CAFile     string
}

// ClientOptions represents common options for creating HTTP clients
type ClientOptions struct {
	BaseURL   string            // Base URL for the service
	APIPath   string            // API path to append to base URL (e.g., "api/" or "api/v1/")
	Timeout   time.Duration     // Request timeout (default: 30s)
	Auth      AuthConfig        // Authentication configuration
	TLS       TLSConfig         // TLS configuration
	UserAgent string            // Custom user agent
	Headers   map[string]string // Additional headers
}

// BuildHTTPClient creates a configured HTTP client with common settings
func BuildHTTPClient(opts *ClientOptions) (*http.Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("client options cannot be nil")
	}

	if opts.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	// Validate URL
	_, err := url.Parse(opts.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Set default timeout
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = constants.DefaultHTTPTimeout
	}

	// Create HTTP client with timeout
	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(timeout)

	// Configure TLS if needed
	if opts.TLS.SkipVerify || opts.TLS.CertFile != "" {
		if transport, ok := httpClient.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: opts.TLS.SkipVerify,
			}

			// Load TLS certificates if provided
			if opts.TLS.CertFile != "" && opts.TLS.KeyFile != "" {
				cert, err := tls.LoadX509KeyPair(opts.TLS.CertFile, opts.TLS.KeyFile)
				if err != nil {
					return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
				}
				transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
			}
		}
	}

	return httpClient, nil
}

// BuildBaseURL constructs the full base URL with API path
func BuildBaseURL(baseURL, apiPath string) (string, error) {
	if baseURL == "" {
		return "", fmt.Errorf("base URL cannot be empty")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Ensure URL has proper path
	if !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
	}

	// Append API path if provided
	if apiPath != "" {
		parsedURL.Path += apiPath
	}

	return parsedURL.String(), nil
}

// GetDefaultHeaders returns default headers for HTTP requests
func GetDefaultHeaders(contentType string) map[string]string {
	headers := make(map[string]string)
	if contentType != "" {
		headers["Content-Type"] = contentType
	} else {
		headers["Content-Type"] = "application/json"
	}
	headers["Accept"] = "application/json"
	return headers
}

// ValidateClientOptions validates client options
func ValidateClientOptions(opts *ClientOptions) error {
	if opts == nil {
		return fmt.Errorf("client options cannot be nil")
	}

	if opts.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}

	// Validate authentication configuration
	if opts.Auth.Enabled {
		switch opts.Auth.Mode {
		case "apikey":
			if opts.Auth.APIKey == "" {
				return fmt.Errorf("API key is required for apikey authentication")
			}
		case "bearer":
			if opts.Auth.BearerToken == "" {
				return fmt.Errorf("bearer token is required for bearer authentication")
			}
		case "basic":
			if opts.Auth.Username == "" || opts.Auth.Password == "" {
				return fmt.Errorf("username and password are required for basic authentication")
			}
		default:
			return fmt.Errorf("invalid authentication mode: %s", opts.Auth.Mode)
		}
	}

	// Validate TLS configuration
	if opts.TLS.CertFile != "" && opts.TLS.KeyFile == "" {
		return fmt.Errorf("TLS key file is required when certificate file is provided")
	}
	if opts.TLS.KeyFile != "" && opts.TLS.CertFile == "" {
		return fmt.Errorf("TLS certificate file is required when key file is provided")
	}

	return nil
}
