package jaeger

import (
	"fmt"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/errors"
)

// Jaeger-specific error codes
const (
	// Connection errors
	ErrCodeConnectionFailed = "JAEGER_CONNECTION_FAILED"
	ErrCodeAuthentication   = "JAEGER_AUTHENTICATION"
	ErrCodeUnauthorized     = "JAEGER_UNAUTHORIZED"
	ErrCodeForbidden        = "JAEGER_FORBIDDEN"

	// Trace errors
	ErrCodeTraceNotFound     = "JAEGER_TRACE_NOT_FOUND"
	ErrCodeTraceQueryFailed  = "JAEGER_TRACE_QUERY_FAILED"
	ErrCodeTraceSearchFailed = "JAEGER_TRACE_SEARCH_FAILED"
	ErrCodeTraceInvalidID    = "JAEGER_TRACE_INVALID_ID"

	// Service errors
	ErrCodeServiceNotFound   = "JAEGER_SERVICE_NOT_FOUND"
	ErrCodeServiceListFailed = "JAEGER_SERVICE_LIST_FAILED"

	// Operation errors
	ErrCodeOperationNotFound   = "JAEGER_OPERATION_NOT_FOUND"
	ErrCodeOperationListFailed = "JAEGER_OPERATION_LIST_FAILED"

	// Span errors
	ErrCodeSpanNotFound = "JAEGER_SPAN_NOT_FOUND"

	// Dependency errors
	ErrCodeDependencyQueryFailed = "JAEGER_DEPENDENCY_QUERY_FAILED"

	// API errors
	ErrCodeAPIError        = "JAEGER_API_ERROR"
	ErrCodeInvalidResponse = "JAEGER_INVALID_RESPONSE"
	ErrCodeRateLimited     = "JAEGER_RATE_LIMITED"
	ErrCodeServerError     = "JAEGER_SERVER_ERROR"
)

// ConnectionFailedError creates a connection failed error
func ConnectionFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeConnectionFailed, "failed to connect to Jaeger").
		WithHTTPStatus(503)
}

// AuthenticationError creates an authentication error
func AuthenticationError() *errors.ServiceError {
	return errors.New(ErrCodeAuthentication, "Jaeger authentication failed").
		WithHTTPStatus(401)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError() *errors.ServiceError {
	return errors.New(ErrCodeUnauthorized, "unauthorized access to Jaeger").
		WithHTTPStatus(401)
}

// ForbiddenError creates a forbidden error
func ForbiddenError() *errors.ServiceError {
	return errors.New(ErrCodeForbidden, "forbidden access to Jaeger resource").
		WithHTTPStatus(403)
}

// TraceNotFoundError creates a trace not found error
func TraceNotFoundError(traceID string) *errors.ServiceError {
	return errors.New(ErrCodeTraceNotFound, fmt.Sprintf("trace not found: %s", traceID)).
		WithHTTPStatus(404).
		WithContext("trace_id", traceID)
}

// TraceQueryFailedError creates a trace query failed error
func TraceQueryFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeTraceQueryFailed, "failed to query trace").
		WithHTTPStatus(500)
}

// TraceSearchFailedError creates a trace search failed error
func TraceSearchFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeTraceSearchFailed, "failed to search traces").
		WithHTTPStatus(500)
}

// TraceInvalidIDError creates a trace invalid ID error
func TraceInvalidIDError(traceID string) *errors.ServiceError {
	return errors.New(ErrCodeTraceInvalidID, fmt.Sprintf("invalid trace ID: %s", traceID)).
		WithHTTPStatus(400).
		WithContext("trace_id", traceID)
}

// ServiceNotFoundError creates a service not found error
func ServiceNotFoundError(service string) *errors.ServiceError {
	return errors.New(ErrCodeServiceNotFound, fmt.Sprintf("service not found: %s", service)).
		WithHTTPStatus(404).
		WithContext("service", service)
}

// ServiceListFailedError creates a service list failed error
func ServiceListFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeServiceListFailed, "failed to list services").
		WithHTTPStatus(500)
}

// OperationNotFoundError creates an operation not found error
func OperationNotFoundError(operation string) *errors.ServiceError {
	return errors.New(ErrCodeOperationNotFound, fmt.Sprintf("operation not found: %s", operation)).
		WithHTTPStatus(404).
		WithContext("operation", operation)
}

// OperationListFailedError creates an operation list failed error
func OperationListFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeOperationListFailed, "failed to list operations").
		WithHTTPStatus(500)
}

// SpanNotFoundError creates a span not found error
func SpanNotFoundError(spanID string) *errors.ServiceError {
	return errors.New(ErrCodeSpanNotFound, fmt.Sprintf("span not found: %s", spanID)).
		WithHTTPStatus(404).
		WithContext("span_id", spanID)
}

// DependencyQueryFailedError creates a dependency query failed error
func DependencyQueryFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDependencyQueryFailed, "failed to query dependencies").
		WithHTTPStatus(500)
}

// APIError creates a generic API error
func APIError(statusCode int, message string) *errors.ServiceError {
	return errors.New(ErrCodeAPIError, fmt.Sprintf("Jaeger API error (status %d): %s", statusCode, message)).
		WithHTTPStatus(statusCode).
		WithContext("status_code", statusCode)
}

// InvalidResponseError creates an invalid response error
func InvalidResponseError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeInvalidResponse, "invalid Jaeger API response").
		WithHTTPStatus(502)
}

// RateLimitedError creates a rate limited error
func RateLimitedError() *errors.ServiceError {
	return errors.New(ErrCodeRateLimited, "Jaeger API rate limit exceeded").
		WithHTTPStatus(429)
}

// ServerError creates a server error
func ServerError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeServerError, "Jaeger server error").
		WithHTTPStatus(500)
}
