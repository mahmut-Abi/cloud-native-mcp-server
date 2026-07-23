package victoriametrics

import (
	"fmt"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/errors"
)

// VictoriaMetrics-specific error codes
const (
	// Connection errors
	ErrCodeConnectionFailed = "VICTORIAMETRICS_CONNECTION_FAILED"
	ErrCodeAuthentication   = "VICTORIAMETRICS_AUTHENTICATION"
	ErrCodeUnauthorized     = "VICTORIAMETRICS_UNAUTHORIZED"
	ErrCodeForbidden        = "VICTORIAMETRICS_FORBIDDEN"

	// Query errors
	ErrCodeQueryFailed      = "VICTORIAMETRICS_QUERY_FAILED"
	ErrCodeQueryBadData     = "VICTORIAMETRICS_QUERY_BAD_DATA"
	ErrCodeQueryTimeout     = "VICTORIAMETRICS_QUERY_TIMEOUT"
	ErrCodeQueryExecTimeout = "VICTORIAMETRICS_QUERY_EXECUTION_TIMEOUT"

	// Export errors
	ErrCodeExportFailed = "VICTORIAMETRICS_EXPORT_FAILED"

	// Import errors
	ErrCodeImportFailed = "VICTORIAMETRICS_IMPORT_FAILED"

	// Admin errors
	ErrCodeDeleteSeriesFailed = "VICTORIAMETRICS_DELETE_SERIES_FAILED"

	// Label errors
	ErrCodeLabelNotFound = "VICTORIAMETRICS_LABEL_NOT_FOUND"

	// Metadata errors
	ErrCodeMetadataFailed = "VICTORIAMETRICS_METADATA_FAILED"

	// API errors
	ErrCodeAPIError        = "VICTORIAMETRICS_API_ERROR"
	ErrCodeInvalidResponse = "VICTORIAMETRICS_INVALID_RESPONSE"
	ErrCodeRateLimited     = "VICTORIAMETRICS_RATE_LIMITED"
	ErrCodeServerError     = "VICTORIAMETRICS_SERVER_ERROR"
)

// ConnectionFailedError creates a connection failed error
func ConnectionFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeConnectionFailed, "failed to connect to VictoriaMetrics").
		WithHTTPStatus(503)
}

// AuthenticationError creates an authentication error
func AuthenticationError() *errors.ServiceError {
	return errors.New(ErrCodeAuthentication, "VictoriaMetrics authentication failed").
		WithHTTPStatus(401)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError() *errors.ServiceError {
	return errors.New(ErrCodeUnauthorized, "unauthorized access to VictoriaMetrics").
		WithHTTPStatus(401)
}

// ForbiddenError creates a forbidden error
func ForbiddenError() *errors.ServiceError {
	return errors.New(ErrCodeForbidden, "forbidden access to VictoriaMetrics resource").
		WithHTTPStatus(403)
}

// QueryFailedError creates a query failed error
func QueryFailedError(err error, query string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeQueryFailed, "VictoriaMetrics query failed").
		WithHTTPStatus(400).
		WithContext("query", query)
}

// QueryBadDataError creates a query bad data error
func QueryBadDataError(err error, query string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeQueryBadData, "VictoriaMetrics query returned bad data").
		WithHTTPStatus(400).
		WithContext("query", query)
}

// QueryTimeoutError creates a query timeout error
func QueryTimeoutError(query string, timeout string) *errors.ServiceError {
	return fmt.Errorf("VictoriaMetrics query timed out after %s: %s", timeout, query).
		(*errors.ServiceError)
}

// ExportFailedError creates an export failed error
func ExportFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeExportFailed, "VictoriaMetrics export failed").
		WithHTTPStatus(500)
}

// ImportFailedError creates an import failed error
func ImportFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeImportFailed, "VictoriaMetrics import failed").
		WithHTTPStatus(500)
}

// DeleteSeriesFailedError creates a delete series failed error
func DeleteSeriesFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDeleteSeriesFailed, "failed to delete VictoriaMetrics series").
		WithHTTPStatus(500)
}

// LabelNotFoundError creates a label not found error
func LabelNotFoundError(label string) *errors.ServiceError {
	return errors.New(ErrCodeLabelNotFound, fmt.Sprintf("label '%s' not found in VictoriaMetrics", label)).
		WithHTTPStatus(404)
}

// MetadataFailedError creates a metadata failed error
func MetadataFailedError(err error, metric string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeMetadataFailed, "failed to get VictoriaMetrics metadata").
		WithHTTPStatus(500).
		WithContext("metric", metric)
}

// APIError creates an API error
func APIError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeAPIError, "VictoriaMetrics API error").
		WithHTTPStatus(500)
}

// InvalidResponseError creates an invalid response error
func InvalidResponseError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeInvalidResponse, "invalid VictoriaMetrics response").
		WithHTTPStatus(502)
}

// RateLimitedError creates a rate limited error
func RateLimitedError() *errors.ServiceError {
	return errors.New(ErrCodeRateLimited, "VictoriaMetrics API rate limit exceeded").
		WithHTTPStatus(429)
}

// ServerError creates a server error
func ServerError() *errors.ServiceError {
	return errors.New(ErrCodeServerError, "VictoriaMetrics server error").
		WithHTTPStatus(500)
}
