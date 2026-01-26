package prometheus

import (
	"fmt"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/errors"
)

// Prometheus-specific error codes
const (
	// Connection errors
	ErrCodeConnectionFailed = "PROMETHEUS_CONNECTION_FAILED"
	ErrCodeAuthentication   = "PROMETHEUS_AUTHENTICATION"
	ErrCodeUnauthorized     = "PROMETHEUS_UNAUTHORIZED"
	ErrCodeForbidden        = "PROMETHEUS_FORBIDDEN"

	// Query errors
	ErrCodeQueryFailed      = "PROMETHEUS_QUERY_FAILED"
	ErrCodeQueryBadData     = "PROMETHEUS_QUERY_BAD_DATA"
	ErrCodeQueryTimeout     = "PROMETHEUS_QUERY_TIMEOUT"
	ErrCodeQueryExecTimeout = "PROMETHEUS_QUERY_EXECUTION_TIMEOUT"

	// Target errors
	ErrCodeTargetNotFound = "PROMETHEUS_TARGET_NOT_FOUND"

	// Rule errors
	ErrCodeRuleNotFound = "PROMETHEUS_RULE_NOT_FOUND"

	// Label errors
	ErrCodeLabelNotFound = "PROMETHEUS_LABEL_NOT_FOUND"

	// Series errors
	ErrCodeSeriesNotFound = "PROMETHEUS_SERIES_NOT_FOUND"

	// TSDB errors
	ErrCodeTSDBSnapshotFailed = "PROMETHEUS_TSDB_SNAPSHOT_FAILED"
	ErrCodeTSDBCleanupFailed  = "PROMETHEUS_TSDB_CLEANUP_FAILED"

	// API errors
	ErrCodeAPIError        = "PROMETHEUS_API_ERROR"
	ErrCodeInvalidResponse = "PROMETHEUS_INVALID_RESPONSE"
	ErrCodeRateLimited     = "PROMETHEUS_RATE_LIMITED"
	ErrCodeServerError     = "PROMETHEUS_SERVER_ERROR"
)

// ConnectionFailedError creates a connection failed error
func ConnectionFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeConnectionFailed, "failed to connect to Prometheus").
		WithHTTPStatus(503)
}

// AuthenticationError creates an authentication error
func AuthenticationError() *errors.ServiceError {
	return errors.New(ErrCodeAuthentication, "Prometheus authentication failed").
		WithHTTPStatus(401)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError() *errors.ServiceError {
	return errors.New(ErrCodeUnauthorized, "unauthorized access to Prometheus").
		WithHTTPStatus(401)
}

// ForbiddenError creates a forbidden error
func ForbiddenError() *errors.ServiceError {
	return errors.New(ErrCodeForbidden, "forbidden access to Prometheus resource").
		WithHTTPStatus(403)
}

// QueryFailedError creates a query failed error
func QueryFailedError(err error, query string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeQueryFailed, "Prometheus query failed").
		WithHTTPStatus(400).
		WithContext("query", query)
}

// QueryBadDataError creates a query bad data error
func QueryBadDataError(err error, query string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeQueryBadData, "Prometheus query returned bad data").
		WithHTTPStatus(400).
		WithContext("query", query)
}

// QueryTimeoutError creates a query timeout error
func QueryTimeoutError(query string) *errors.ServiceError {
	return errors.New(ErrCodeQueryTimeout, fmt.Sprintf("Prometheus query timeout: %s", query)).
		WithHTTPStatus(504).
		WithContext("query", query)
}

// QueryExecutionTimeoutError creates a query execution timeout error
func QueryExecutionTimeoutError(query string) *errors.ServiceError {
	return errors.New(ErrCodeQueryExecTimeout, fmt.Sprintf("Prometheus query execution timeout: %s", query)).
		WithHTTPStatus(504).
		WithContext("query", query)
}

// TargetNotFoundError creates a target not found error
func TargetNotFoundError(target string) *errors.ServiceError {
	return errors.New(ErrCodeTargetNotFound, fmt.Sprintf("target not found: %s", target)).
		WithHTTPStatus(404).
		WithContext("target", target)
}

// RuleNotFoundError creates a rule not found error
func RuleNotFoundError(rule string) *errors.ServiceError {
	return errors.New(ErrCodeRuleNotFound, fmt.Sprintf("rule not found: %s", rule)).
		WithHTTPStatus(404).
		WithContext("rule", rule)
}

// LabelNotFoundError creates a label not found error
func LabelNotFoundError(label string) *errors.ServiceError {
	return errors.New(ErrCodeLabelNotFound, fmt.Sprintf("label not found: %s", label)).
		WithHTTPStatus(404).
		WithContext("label", label)
}

// SeriesNotFoundError creates a series not found error
func SeriesNotFoundError(match string) *errors.ServiceError {
	return errors.New(ErrCodeSeriesNotFound, fmt.Sprintf("series not found: %s", match)).
		WithHTTPStatus(404).
		WithContext("match", match)
}

// TSDBSnapshotFailedError creates a TSDB snapshot failed error
func TSDBSnapshotFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeTSDBSnapshotFailed, "failed to create TSDB snapshot").
		WithHTTPStatus(500)
}

// TSDBCleanupFailedError creates a TSDB cleanup failed error
func TSDBCleanupFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeTSDBCleanupFailed, "failed to cleanup TSDB").
		WithHTTPStatus(500)
}

// APIError creates a generic API error
func APIError(statusCode int, message string) *errors.ServiceError {
	return errors.New(ErrCodeAPIError, fmt.Sprintf("Prometheus API error (status %d): %s", statusCode, message)).
		WithHTTPStatus(statusCode).
		WithContext("status_code", statusCode)
}

// InvalidResponseError creates an invalid response error
func InvalidResponseError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeInvalidResponse, "invalid Prometheus API response").
		WithHTTPStatus(502)
}

// RateLimitedError creates a rate limited error
func RateLimitedError() *errors.ServiceError {
	return errors.New(ErrCodeRateLimited, "Prometheus API rate limit exceeded").
		WithHTTPStatus(429)
}

// ServerError creates a server error
func ServerError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeServerError, "Prometheus server error").
		WithHTTPStatus(500)
}
