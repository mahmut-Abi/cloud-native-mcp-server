// Package middleware provides HTTP middleware components for the cloud-native MCP server.
package middleware

import (
	"log"
	"net/http"
	"sync"
)

// BackendAuthHandler is a function that parses backend service credentials from
// HTTP request headers, creates a service-specific client, and returns a new
// request with the client stored in the request context.
//
// If the required headers are not present, it returns the original request
// unchanged — the tool handler will return an appropriate error when it tries
// to access the missing client from context.
type BackendAuthHandler func(r *http.Request) (*http.Request, error)

var (
	backendHandlers   = make(map[string]BackendAuthHandler)
	backendHandlersMu sync.RWMutex
)

// RegisterBackendAuthHandler registers a per-service handler that parses
// backend credentials from HTTP headers and injects a client into the context.
// Each service calls this during initialization to register its own handler.
func RegisterBackendAuthHandler(serviceName string, handler BackendAuthHandler) {
	backendHandlersMu.Lock()
	defer backendHandlersMu.Unlock()
	backendHandlers[serviceName] = handler
	log.Printf("Registered backend auth handler for service: %s", serviceName)
}

// BackendAuthMiddleware creates an HTTP middleware that extracts backend
// service credentials from request headers for a specific service.
//
// The middleware looks up the registered BackendAuthHandler for the given
// serviceName. If found, the handler parses headers and may inject a client
// into the request context. If no handler is registered or headers are
// missing, the request passes through unchanged.
func BackendAuthMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			backendHandlersMu.RLock()
			handler, ok := backendHandlers[serviceName]
			backendHandlersMu.RUnlock()

			if !ok {
				// No backend auth handler registered for this service
				next.ServeHTTP(w, r)
				return
			}

			newReq, err := handler(r)
			if err != nil {
				log.Printf("[backend-auth] %s handler failed: %v", serviceName, err)
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, newReq)
		})
	}
}
