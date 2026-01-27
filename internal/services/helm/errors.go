package helm

import (
	"fmt"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/errors"
)

// Helm-specific error codes
const (
	// Connection errors
	ErrCodeConnectionFailed = "HELM_CONNECTION_FAILED"
	ErrCodeKubeconfigFailed = "HELM_KUBECONFIG_FAILED"

	// Release errors
	ErrCodeReleaseNotFound       = "HELM_RELEASE_NOT_FOUND"
	ErrCodeReleaseCreateFailed   = "HELM_RELEASE_CREATE_FAILED"
	ErrCodeReleaseUpdateFailed   = "HELM_RELEASE_UPDATE_FAILED"
	ErrCodeReleaseDeleteFailed   = "HELM_RELEASE_DELETE_FAILED"
	ErrCodeReleaseRollbackFailed = "HELM_RELEASE_ROLLBACK_FAILED"
	ErrCodeReleaseStatusFailed   = "HELM_RELEASE_STATUS_FAILED"
	ErrCodeReleaseHistoryFailed  = "HELM_RELEASE_HISTORY_FAILED"
	ErrCodeReleaseValuesFailed   = "HELM_RELEASE_VALUES_FAILED"
	ErrCodeReleaseManifestFailed = "HELM_RELEASE_MANIFEST_FAILED"
	ErrCodeReleaseTestFailed     = "HELM_RELEASE_TEST_FAILED"

	// Chart errors
	ErrCodeChartNotFound       = "HELM_CHART_NOT_FOUND"
	ErrCodeChartPullFailed     = "HELM_CHART_PULL_FAILED"
	ErrCodeChartLoadFailed     = "HELM_CHART_LOAD_FAILED"
	ErrCodeChartTemplateFailed = "HELM_CHART_TEMPLATE_FAILED"
	ErrCodeChartSearchFailed   = "HELM_CHART_SEARCH_FAILED"
	ErrCodeChartInfoFailed     = "HELM_CHART_INFO_FAILED"

	// Repository errors
	ErrCodeRepositoryNotFound     = "HELM_REPOSITORY_NOT_FOUND"
	ErrCodeRepositoryAddFailed    = "HELM_REPOSITORY_ADD_FAILED"
	ErrCodeRepositoryRemoveFailed = "HELM_REPOSITORY_REMOVE_FAILED"
	ErrCodeRepositoryUpdateFailed = "HELM_REPOSITORY_UPDATE_FAILED"
	ErrCodeRepositoryListFailed   = "HELM_REPOSITORY_LIST_FAILED"
	ErrCodeRepositoryIndexFailed  = "HELM_REPOSITORY_INDEX_FAILED"

	// Dependency errors
	ErrCodeDependencyFailed = "HELM_DEPENDENCY_FAILED"

	// Validation errors
	ErrCodeValidationFailed = "HELM_VALIDATION_FAILED"
	ErrCodeLintFailed       = "HELM_LINT_FAILED"

	// Timeout errors
	ErrCodeInstallTimeout  = "HELM_INSTALL_TIMEOUT"
	ErrCodeUpgradeTimeout  = "HELM_UPGRADE_TIMEOUT"
	ErrCodeDeleteTimeout   = "HELM_DELETE_TIMEOUT"
	ErrCodeRollbackTimeout = "HELM_ROLLBACK_TIMEOUT"
	ErrCodeTestTimeout     = "HELM_TEST_TIMEOUT"
	ErrCodeSearchTimeout   = "HELM_SEARCH_TIMEOUT"
	ErrCodeUpdateTimeout   = "HELM_UPDATE_TIMEOUT"

	// Configuration errors
	ErrCodeValuesInvalid    = "HELM_VALUES_INVALID"
	ErrCodeNamespaceInvalid = "HELM_NAMESPACE_INVALID"
	ErrCodeNameInvalid      = "HELM_NAME_INVALID"

	// API errors
	ErrCodeAPIError      = "HELM_API_ERROR"
	ErrCodeInternalError = "HELM_INTERNAL_ERROR"
)

// ConnectionFailedError creates a connection failed error
func ConnectionFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeConnectionFailed, "failed to connect to Kubernetes cluster").
		WithHTTPStatus(503)
}

// KubeconfigFailedError creates a kubeconfig failed error
func KubeconfigFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeKubeconfigFailed, "failed to load kubeconfig").
		WithHTTPStatus(500)
}

// ReleaseNotFoundError creates a release not found error
func ReleaseNotFoundError(name, namespace string) *errors.ServiceError {
	return errors.New(ErrCodeReleaseNotFound, fmt.Sprintf("release not found: %s in namespace %s", name, namespace)).
		WithHTTPStatus(404).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseCreateFailedError creates a release creation failed error
func ReleaseCreateFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseCreateFailed, fmt.Sprintf("failed to create release: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseUpdateFailedError creates a release update failed error
func ReleaseUpdateFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseUpdateFailed, fmt.Sprintf("failed to update release: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseDeleteFailedError creates a release deletion failed error
func ReleaseDeleteFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseDeleteFailed, fmt.Sprintf("failed to delete release: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseRollbackFailedError creates a release rollback failed error
func ReleaseRollbackFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseRollbackFailed, fmt.Sprintf("failed to rollback release: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseStatusFailedError creates a release status failed error
func ReleaseStatusFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseStatusFailed, fmt.Sprintf("failed to get release status: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseHistoryFailedError creates a release history failed error
func ReleaseHistoryFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseHistoryFailed, fmt.Sprintf("failed to get release history: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseValuesFailedError creates a release values failed error
func ReleaseValuesFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseValuesFailed, fmt.Sprintf("failed to get release values: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseManifestFailedError creates a release manifest failed error
func ReleaseManifestFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseManifestFailed, fmt.Sprintf("failed to get release manifest: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ReleaseTestFailedError creates a release test failed error
func ReleaseTestFailedError(err error, name, namespace string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeReleaseTestFailed, fmt.Sprintf("release test failed: %s in namespace %s", name, namespace)).
		WithHTTPStatus(500).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// ChartNotFoundError creates a chart not found error
func ChartNotFoundError(chart string) *errors.ServiceError {
	return errors.New(ErrCodeChartNotFound, fmt.Sprintf("chart not found: %s", chart)).
		WithHTTPStatus(404).
		WithContext("chart", chart)
}

// ChartPullFailedError creates a chart pull failed error
func ChartPullFailedError(err error, chart string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeChartPullFailed, fmt.Sprintf("failed to pull chart: %s", chart)).
		WithHTTPStatus(500).
		WithContext("chart", chart)
}

// ChartLoadFailedError creates a chart load failed error
func ChartLoadFailedError(err error, chart string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeChartLoadFailed, fmt.Sprintf("failed to load chart: %s", chart)).
		WithHTTPStatus(500).
		WithContext("chart", chart)
}

// ChartTemplateFailedError creates a chart template failed error
func ChartTemplateFailedError(err error, chart string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeChartTemplateFailed, fmt.Sprintf("failed to template chart: %s", chart)).
		WithHTTPStatus(500).
		WithContext("chart", chart)
}

// ChartSearchFailedError creates a chart search failed error
func ChartSearchFailedError(err error, keyword string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeChartSearchFailed, fmt.Sprintf("failed to search charts: %s", keyword)).
		WithHTTPStatus(500).
		WithContext("keyword", keyword)
}

// ChartInfoFailedError creates a chart info failed error
func ChartInfoFailedError(err error, chart string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeChartInfoFailed, fmt.Sprintf("failed to get chart info: %s", chart)).
		WithHTTPStatus(500).
		WithContext("chart", chart)
}

// RepositoryNotFoundError creates a repository not found error
func RepositoryNotFoundError(name string) *errors.ServiceError {
	return errors.New(ErrCodeRepositoryNotFound, fmt.Sprintf("repository not found: %s", name)).
		WithHTTPStatus(404).
		WithContext("repository", name)
}

// RepositoryAddFailedError creates a repository add failed error
func RepositoryAddFailedError(err error, name string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeRepositoryAddFailed, fmt.Sprintf("failed to add repository: %s", name)).
		WithHTTPStatus(500).
		WithContext("repository", name)
}

// RepositoryRemoveFailedError creates a repository remove failed error
func RepositoryRemoveFailedError(err error, name string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeRepositoryRemoveFailed, fmt.Sprintf("failed to remove repository: %s", name)).
		WithHTTPStatus(500).
		WithContext("repository", name)
}

// RepositoryUpdateFailedError creates a repository update failed error
func RepositoryUpdateFailedError(err error, name string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeRepositoryUpdateFailed, fmt.Sprintf("failed to update repository: %s", name)).
		WithHTTPStatus(500).
		WithContext("repository", name)
}

// RepositoryListFailedError creates a repository list failed error
func RepositoryListFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeRepositoryListFailed, "failed to list repositories").
		WithHTTPStatus(500)
}

// RepositoryIndexFailedError creates a repository index failed error
func RepositoryIndexFailedError(err error, name string) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeRepositoryIndexFailed, fmt.Sprintf("failed to index repository: %s", name)).
		WithHTTPStatus(500).
		WithContext("repository", name)
}

// DependencyFailedError creates a dependency failed error
func DependencyFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeDependencyFailed, "failed to resolve dependencies").
		WithHTTPStatus(500)
}

// ValidationFailedError creates a validation failed error
func ValidationFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeValidationFailed, "validation failed").
		WithHTTPStatus(400)
}

// LintFailedError creates a lint failed error
func LintFailedError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeLintFailed, "lint failed").
		WithHTTPStatus(400)
}

// InstallTimeoutError creates an install timeout error
func InstallTimeoutError(name, namespace string) *errors.ServiceError {
	return errors.New(ErrCodeInstallTimeout, fmt.Sprintf("install timeout: %s in namespace %s", name, namespace)).
		WithHTTPStatus(504).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// UpgradeTimeoutError creates an upgrade timeout error
func UpgradeTimeoutError(name, namespace string) *errors.ServiceError {
	return errors.New(ErrCodeUpgradeTimeout, fmt.Sprintf("upgrade timeout: %s in namespace %s", name, namespace)).
		WithHTTPStatus(504).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// DeleteTimeoutError creates a delete timeout error
func DeleteTimeoutError(name, namespace string) *errors.ServiceError {
	return errors.New(ErrCodeDeleteTimeout, fmt.Sprintf("delete timeout: %s in namespace %s", name, namespace)).
		WithHTTPStatus(504).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// RollbackTimeoutError creates a rollback timeout error
func RollbackTimeoutError(name, namespace string) *errors.ServiceError {
	return errors.New(ErrCodeRollbackTimeout, fmt.Sprintf("rollback timeout: %s in namespace %s", name, namespace)).
		WithHTTPStatus(504).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// TestTimeoutError creates a test timeout error
func TestTimeoutError(name, namespace string) *errors.ServiceError {
	return errors.New(ErrCodeTestTimeout, fmt.Sprintf("test timeout: %s in namespace %s", name, namespace)).
		WithHTTPStatus(504).
		WithContext("release", name).
		WithContext("namespace", namespace)
}

// SearchTimeoutError creates a search timeout error
func SearchTimeoutError(keyword string) *errors.ServiceError {
	return errors.New(ErrCodeSearchTimeout, fmt.Sprintf("search timeout: %s", keyword)).
		WithHTTPStatus(504).
		WithContext("keyword", keyword)
}

// UpdateTimeoutError creates an update timeout error
func UpdateTimeoutError() *errors.ServiceError {
	return errors.New(ErrCodeUpdateTimeout, "repository update timeout").
		WithHTTPStatus(504)
}

// ValuesInvalidError creates a values invalid error
func ValuesInvalidError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeValuesInvalid, "invalid values").
		WithHTTPStatus(400)
}

// NamespaceInvalidError creates a namespace invalid error
func NamespaceInvalidError(namespace string) *errors.ServiceError {
	return errors.New(ErrCodeNamespaceInvalid, fmt.Sprintf("invalid namespace: %s", namespace)).
		WithHTTPStatus(400).
		WithContext("namespace", namespace)
}

// NameInvalidError creates a name invalid error
func NameInvalidError(name string) *errors.ServiceError {
	return errors.New(ErrCodeNameInvalid, fmt.Sprintf("invalid release name: %s", name)).
		WithHTTPStatus(400).
		WithContext("name", name)
}

// APIError creates a generic API error
func APIError(message string) *errors.ServiceError {
	return errors.New(ErrCodeAPIError, fmt.Sprintf("Helm API error: %s", message)).
		WithHTTPStatus(500)
}

// InternalError creates an internal error
func InternalError(err error) *errors.ServiceError {
	return errors.Wrap(err, ErrCodeInternalError, "Helm internal error").
		WithHTTPStatus(500)
}
