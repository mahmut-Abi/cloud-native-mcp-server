package grafana

import (
	"fmt"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/errors"
)

// Grafana-specific error codes
const (
	// Connection errors
	ErrCodeConnectionFailed = "GRAFANA_CONNECTION_FAILED"
	ErrCodeAuthentication   = "GRAFANA_AUTHENTICATION"
	ErrCodeUnauthorized     = "GRAFANA_UNAUTHORIZED"
	ErrCodeForbidden        = "GRAFANA_FORBIDDEN"

	// Dashboard errors
	ErrCodeDashboardNotFound     = "GRAFANA_DASHBOARD_NOT_FOUND"
	ErrCodeDashboardCreateFailed = "GRAFANA_DASHBOARD_CREATE_FAILED"
	ErrCodeDashboardUpdateFailed = "GRAFANA_DASHBOARD_UPDATE_FAILED"
	ErrCodeDashboardDeleteFailed = "GRAFANA_DASHBOARD_DELETE_FAILED"
	ErrCodeDashboardImportFailed = "GRAFANA_DASHBOARD_IMPORT_FAILED"
	ErrCodeDashboardExportFailed = "GRAFANA_DASHBOARD_EXPORT_FAILED"
	ErrCodeDashboardRenderFailed = "GRAFANA_DASHBOARD_RENDER_FAILED"
	ErrCodePanelNotFound         = "GRAFANA_PANEL_NOT_FOUND"

	// Data source errors
	ErrCodeDatasourceNotFound     = "GRAFANA_DATASOURCE_NOT_FOUND"
	ErrCodeDatasourceCreateFailed = "GRAFANA_DATASOURCE_CREATE_FAILED"
	ErrCodeDatasourceUpdateFailed = "GRAFANA_DATASOURCE_UPDATE_FAILED"
	ErrCodeDatasourceDeleteFailed = "GRAFANA_DATASOURCE_DELETE_FAILED"
	ErrCodeDatasourceTestFailed   = "GRAFANA_DATASOURCE_TEST_FAILED"
	ErrCodeDatasourceHealthFailed = "GRAFANA_DATASOURCE_HEALTH_FAILED"

	// Folder errors
	ErrCodeFolderNotFound     = "GRAFANA_FOLDER_NOT_FOUND"
	ErrCodeFolderCreateFailed = "GRAFANA_FOLDER_CREATE_FAILED"
	ErrCodeFolderUpdateFailed = "GRAFANA_FOLDER_UPDATE_FAILED"
	ErrCodeFolderDeleteFailed = "GRAFANA_FOLDER_DELETE_FAILED"

	// Alert rule errors
	ErrCodeAlertRuleNotFound     = "GRAFANA_ALERT_RULE_NOT_FOUND"
	ErrCodeAlertRuleCreateFailed = "GRAFANA_ALERT_RULE_CREATE_FAILED"
	ErrCodeAlertRuleUpdateFailed = "GRAFANA_ALERT_RULE_UPDATE_FAILED"
	ErrCodeAlertRuleDeleteFailed = "GRAFANA_ALERT_RULE_DELETE_FAILED"

	// Annotation errors
	ErrCodeAnnotationNotFound     = "GRAFANA_ANNOTATION_NOT_FOUND"
	ErrCodeAnnotationCreateFailed = "GRAFANA_ANNOTATION_CREATE_FAILED"
	ErrCodeAnnotationUpdateFailed = "GRAFANA_ANNOTATION_UPDATE_FAILED"
	ErrCodeAnnotationDeleteFailed = "GRAFANA_ANNOTATION_DELETE_FAILED"

	// User/Team errors
	ErrCodeUserNotFound     = "GRAFANA_USER_NOT_FOUND"
	ErrCodeTeamNotFound     = "GRAFANA_TEAM_NOT_FOUND"
	ErrCodeRoleNotFound     = "GRAFANA_ROLE_NOT_FOUND"
	ErrCodePermissionDenied = "GRAFANA_PERMISSION_DENIED"

	// API errors
	ErrCodeAPIError        = "GRAFANA_API_ERROR"
	ErrCodeInvalidResponse = "GRAFANA_INVALID_RESPONSE"
	ErrCodeRateLimited     = "GRAFANA_RATE_LIMITED"
	ErrCodeServerError     = "GRAFANA_SERVER_ERROR"
)

// ConnectionFailedError creates a connection failed error
func ConnectionFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeConnectionFailed, "failed to connect to Grafana").
		WithHTTPStatus(503)
}

// AuthenticationError creates an authentication error
func AuthenticationError() *errors.ServiceError {
	return errors.New(ErrCodeAuthentication, "Grafana authentication failed").
		WithHTTPStatus(401)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError() *errors.ServiceError {
	return errors.New(ErrCodeUnauthorized, "unauthorized access to Grafana").
		WithHTTPStatus(401)
}

// ForbiddenError creates a forbidden error
func ForbiddenError() *errors.ServiceError {
	return errors.New(ErrCodeForbidden, "forbidden access to Grafana resource").
		WithHTTPStatus(403)
}

// DashboardNotFoundError creates a dashboard not found error
func DashboardNotFoundError(uid string) *errors.ServiceError {
	return errors.New(ErrCodeDashboardNotFound, fmt.Sprintf("dashboard not found: %s", uid)).
		WithHTTPStatus(404).
		WithContext("dashboard_uid", uid)
}

// DashboardCreateFailedError creates a dashboard creation failed error
func DashboardCreateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDashboardCreateFailed, "failed to create dashboard").
		WithHTTPStatus(500)
}

// DashboardUpdateFailedError creates a dashboard update failed error
func DashboardUpdateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDashboardUpdateFailed, "failed to update dashboard").
		WithHTTPStatus(500)
}

// DashboardDeleteFailedError creates a dashboard deletion failed error
func DashboardDeleteFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDashboardDeleteFailed, "failed to delete dashboard").
		WithHTTPStatus(500)
}

// DashboardRenderFailedError creates a dashboard render failed error
func DashboardRenderFailedError(err error, panelID int) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDashboardRenderFailed, "failed to render dashboard panel").
		WithHTTPStatus(500).
		WithContext("panel_id", panelID)
}

// DatasourceNotFoundError creates a datasource not found error
func DatasourceNotFoundError(uid string) *errors.ServiceError {
	return errors.New(ErrCodeDatasourceNotFound, fmt.Sprintf("datasource not found: %s", uid)).
		WithHTTPStatus(404).
		WithContext("datasource_uid", uid)
}

// DatasourceCreateFailedError creates a datasource creation failed error
func DatasourceCreateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDatasourceCreateFailed, "failed to create datasource").
		WithHTTPStatus(500)
}

// DatasourceUpdateFailedError creates a datasource update failed error
func DatasourceUpdateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDatasourceUpdateFailed, "failed to update datasource").
		WithHTTPStatus(500)
}

// DatasourceDeleteFailedError creates a datasource deletion failed error
func DatasourceDeleteFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDatasourceDeleteFailed, "failed to delete datasource").
		WithHTTPStatus(500)
}

// DatasourceTestFailedError creates a datasource test failed error
func DatasourceTestFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDatasourceTestFailed, "datasource health check failed").
		WithHTTPStatus(500)
}

// FolderNotFoundError creates a folder not found error
func FolderNotFoundError(uid string) *errors.ServiceError {
	return errors.New(ErrCodeFolderNotFound, fmt.Sprintf("folder not found: %s", uid)).
		WithHTTPStatus(404).
		WithContext("folder_uid", uid)
}

// AlertRuleNotFoundError creates an alert rule not found error
func AlertRuleNotFoundError(uid string) *errors.ServiceError {
	return errors.New(ErrCodeAlertRuleNotFound, fmt.Sprintf("alert rule not found: %s", uid)).
		WithHTTPStatus(404).
		WithContext("alert_rule_uid", uid)
}

// AlertRuleCreateFailedError creates an alert rule creation failed error
func AlertRuleCreateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeAlertRuleCreateFailed, "failed to create alert rule").
		WithHTTPStatus(500)
}

// AlertRuleUpdateFailedError creates an alert rule update failed error
func AlertRuleUpdateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeAlertRuleUpdateFailed, "failed to update alert rule").
		WithHTTPStatus(500)
}

// AlertRuleDeleteFailedError creates an alert rule deletion failed error
func AlertRuleDeleteFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeAlertRuleDeleteFailed, "failed to delete alert rule").
		WithHTTPStatus(500)
}

// AnnotationNotFoundError creates an annotation not found error
func AnnotationNotFoundError(id int) *errors.ServiceError {
	return errors.New(ErrCodeAnnotationNotFound, fmt.Sprintf("annotation not found: %d", id)).
		WithHTTPStatus(404).
		WithContext("annotation_id", id)
}

// AnnotationCreateFailedError creates an annotation creation failed error
func AnnotationCreateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeAnnotationCreateFailed, "failed to create annotation").
		WithHTTPStatus(500)
}

// AnnotationUpdateFailedError creates an annotation update failed error
func AnnotationUpdateFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeAnnotationUpdateFailed, "failed to update annotation").
		WithHTTPStatus(500)
}

// AnnotationDeleteFailedError creates an annotation deletion failed error
func AnnotationDeleteFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeAnnotationDeleteFailed, "failed to delete annotation").
		WithHTTPStatus(500)
}

// APIError creates a generic API error
func APIError(statusCode int, message string) *errors.ServiceError {
	return errors.New(ErrCodeAPIError, fmt.Sprintf("Grafana API error (status %d): %s", statusCode, message)).
		WithHTTPStatus(statusCode).
		WithContext("status_code", statusCode)
}

// InvalidResponseError creates an invalid response error
func InvalidResponseError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeInvalidResponse, "invalid Grafana API response").
		WithHTTPStatus(502)
}

// RateLimitedError creates a rate limited error
func RateLimitedError() *errors.ServiceError {
	return errors.New(ErrCodeRateLimited, "Grafana API rate limit exceeded").
		WithHTTPStatus(429)
}

// ServerError creates a server error
func ServerError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeServerError, "Grafana server error").
		WithHTTPStatus(500)
}
