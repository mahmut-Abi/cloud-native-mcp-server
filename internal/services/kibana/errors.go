package kibana

import (
	"fmt"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/errors"
)

// Kibana-specific error codes
const (
	// Connection errors
	ErrCodeConnectionFailed = "KIBANA_CONNECTION_FAILED"
	ErrCodeAuthentication   = "KIBANA_AUTHENTICATION"
	ErrCodeUnauthorized     = "KIBANA_UNAUTHORIZED"
	ErrCodeForbidden        = "KIBANA_FORBIDDEN"

	// Space errors
	ErrCodeSpaceNotFound     = "KIBANA_SPACE_NOT_FOUND"
	ErrCodeSpaceCreateFailed = "KIBANA_SPACE_CREATE_FAILED"
	ErrCodeSpaceUpdateFailed = "KIBANA_SPACE_UPDATE_FAILED"
	ErrCodeSpaceDeleteFailed = "KIBANA_SPACE_DELETE_FAILED"

	// Dashboard errors
	ErrCodeDashboardNotFound     = "KIBANA_DASHBOARD_NOT_FOUND"
	ErrCodeDashboardCreateFailed = "KIBANA_DASHBOARD_CREATE_FAILED"
	ErrCodeDashboardUpdateFailed = "KIBANA_DASHBOARD_UPDATE_FAILED"
	ErrCodeDashboardDeleteFailed = "KIBANA_DASHBOARD_DELETE_FAILED"
	ErrCodeDashboardCloneFailed  = "KIBANA_DASHBOARD_CLONE_FAILED"

	// Visualization errors
	ErrCodeVisualizationNotFound     = "KIBANA_VISUALIZATION_NOT_FOUND"
	ErrCodeVisualizationCreateFailed = "KIBANA_VISUALIZATION_CREATE_FAILED"
	ErrCodeVisualizationUpdateFailed = "KIBANA_VISUALIZATION_UPDATE_FAILED"
	ErrCodeVisualizationDeleteFailed = "KIBANA_VISUALIZATION_DELETE_FAILED"
	ErrCodeVisualizationCloneFailed  = "KIBANA_VISUALIZATION_CLONE_FAILED"

	// Index pattern errors
	ErrCodeIndexPatternNotFound     = "KIBANA_INDEX_PATTERN_NOT_FOUND"
	ErrCodeIndexPatternCreateFailed = "KIBANA_INDEX_PATTERN_CREATE_FAILED"
	ErrCodeIndexPatternUpdateFailed = "KIBANA_INDEX_PATTERN_UPDATE_FAILED"
	ErrCodeIndexPatternDeleteFailed = "KIBANA_INDEX_PATTERN_DELETE_FAILED"

	// Alert rule errors
	ErrCodeAlertRuleNotFound     = "KIBANA_ALERT_RULE_NOT_FOUND"
	ErrCodeAlertRuleCreateFailed = "KIBANA_ALERT_RULE_CREATE_FAILED"
	ErrCodeAlertRuleUpdateFailed = "KIBANA_ALERT_RULE_UPDATE_FAILED"
	ErrCodeAlertRuleDeleteFailed = "KIBANA_ALERT_RULE_DELETE_FAILED"

	// Connector errors
	ErrCodeConnectorNotFound     = "KIBANA_CONNECTOR_NOT_FOUND"
	ErrCodeConnectorCreateFailed = "KIBANA_CONNECTOR_CREATE_FAILED"
	ErrCodeConnectorUpdateFailed = "KIBANA_CONNECTOR_UPDATE_FAILED"
	ErrCodeConnectorDeleteFailed = "KIBANA_CONNECTOR_DELETE_FAILED"

	// Data view errors
	ErrCodeDataViewNotFound     = "KIBANA_DATA_VIEW_NOT_FOUND"
	ErrCodeDataViewCreateFailed = "KIBANA_DATA_VIEW_CREATE_FAILED"
	ErrCodeDataViewUpdateFailed = "KIBANA_DATA_VIEW_UPDATE_FAILED"
	ErrCodeDataViewDeleteFailed = "KIBANA_DATA_VIEW_DELETE_FAILED"

	// Query errors
	ErrCodeQueryFailed = "KIBANA_QUERY_FAILED"
	ErrCodeQueryError  = "KIBANA_QUERY_ERROR"

	// API errors
	ErrCodeAPIError        = "KIBANA_API_ERROR"
	ErrCodeInvalidResponse = "KIBANA_INVALID_RESPONSE"
	ErrCodeRateLimited     = "KIBANA_RATE_LIMITED"
	ErrCodeServerError     = "KIBANA_SERVER_ERROR"
)

// ConnectionFailedError creates a connection failed error
func ConnectionFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeConnectionFailed, "failed to connect to Kibana").
		WithHTTPStatus(503)
}

// AuthenticationError creates an authentication error
func AuthenticationError() *errors.ServiceError {
	return errors.New(ErrCodeAuthentication, "Kibana authentication failed").
		WithHTTPStatus(401)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError() *errors.ServiceError {
	return errors.New(ErrCodeUnauthorized, "unauthorized access to Kibana").
		WithHTTPStatus(401)
}

// ForbiddenError creates a forbidden error
func ForbiddenError() *errors.ServiceError {
	return errors.New(ErrCodeForbidden, "forbidden access to Kibana resource").
		WithHTTPStatus(403)
}

// SpaceNotFoundError creates a space not found error
func SpaceNotFoundError(id string) *errors.ServiceError {
	return errors.New(ErrCodeSpaceNotFound, fmt.Sprintf("space not found: %s", id)).
		WithHTTPStatus(404).
		WithContext("space_id", id)
}

// DashboardNotFoundError creates a dashboard not found error
func DashboardNotFoundError(id string) *errors.ServiceError {
	return errors.New(ErrCodeDashboardNotFound, fmt.Sprintf("dashboard not found: %s", id)).
		WithHTTPStatus(404).
		WithContext("dashboard_id", id)
}

// VisualizationNotFoundError creates a visualization not found error
func VisualizationNotFoundError(id string) *errors.ServiceError {
	return errors.New(ErrCodeVisualizationNotFound, fmt.Sprintf("visualization not found: %s", id)).
		WithHTTPStatus(404).
		WithContext("visualization_id", id)
}

// IndexPatternNotFoundError creates an index pattern not found error
func IndexPatternNotFoundError(id string) *errors.ServiceError {
	return errors.New(ErrCodeIndexPatternNotFound, fmt.Sprintf("index pattern not found: %s", id)).
		WithHTTPStatus(404).
		WithContext("index_pattern_id", id)
}

// AlertRuleNotFoundError creates an alert rule not found error
func AlertRuleNotFoundError(id string) *errors.ServiceError {
	return errors.New(ErrCodeAlertRuleNotFound, fmt.Sprintf("alert rule not found: %s", id)).
		WithHTTPStatus(404).
		WithContext("alert_rule_id", id)
}

// ConnectorNotFoundError creates a connector not found error
func ConnectorNotFoundError(id string) *errors.ServiceError {
	return errors.New(ErrCodeConnectorNotFound, fmt.Sprintf("connector not found: %s", id)).
		WithHTTPStatus(404).
		WithContext("connector_id", id)
}

// DataViewNotFoundError creates a data view not found error
func DataViewNotFoundError(id string) *errors.ServiceError {
	return errors.New(ErrCodeDataViewNotFound, fmt.Sprintf("data view not found: %s", id)).
		WithHTTPStatus(404).
		WithContext("data_view_id", id)
}

// QueryFailedError creates a query failed error
func QueryFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeQueryFailed, "Kibana query failed").
		WithHTTPStatus(400)
}

// QueryError creates a query error
func QueryError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeQueryError, "Kibana query error").
		WithHTTPStatus(400)
}

// APIError creates a generic API error
func APIError(statusCode int, message string) *errors.ServiceError {
	return errors.New(ErrCodeAPIError, fmt.Sprintf("Kibana API error (status %d): %s", statusCode, message)).
		WithHTTPStatus(statusCode).
		WithContext("status_code", statusCode)
}

// InvalidResponseError creates an invalid response error
func InvalidResponseError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeInvalidResponse, "invalid Kibana API response").
		WithHTTPStatus(502)
}

// RateLimitedError creates a rate limited error
func RateLimitedError() *errors.ServiceError {
	return errors.New(ErrCodeRateLimited, "Kibana API rate limit exceeded").
		WithHTTPStatus(429)
}

// ServerError creates a server error
func ServerError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeServerError, "Kibana server error").
		WithHTTPStatus(500)
}
