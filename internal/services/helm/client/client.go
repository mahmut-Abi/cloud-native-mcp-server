// Package client provides Helm client operations for the MCP server.
// It implements the Helm client for managing Helm releases, charts, and repositories.
package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/rest"

	"helm.sh/helm/v3/cmd/helm/search"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
)

// ClientOptions represents configuration options for the Helm client.
type ClientOptions struct {
	// KubeconfigPath is the path to the kubeconfig file for Helm operations.
	// If empty, it will use the default kubeconfig (~/.kube/config or $KUBECONFIG).
	KubeconfigPath string

	// Namespace is the default namespace for Helm operations.
	// If empty, it will use the namespace specified in the kubeconfig.
	Namespace string

	// Debug enables debug mode for Helm operations.
	Debug bool

	// RESTConfig is the Kubernetes REST configuration.
	// If provided, it will be used instead of building a new configuration from KubeconfigPath.
	RESTConfig *rest.Config

	// Optimizer for handling Helm repository optimization
	Optimizer *RepositoryOptimizer
}

// PaginationInfo represents pagination metadata for Helm responses
type PaginationInfo struct {
	ContinueToken   string `json:"continueToken"`
	RemainingCount  int64  `json:"remainingCount"`
	CurrentPageSize int64  `json:"currentPageSize"`
	HasMore         bool   `json:"hasMore"`
}

// Client represents a Helm client for managing Helm releases and charts.
type Client struct {
	settings     *cli.EnvSettings
	options      ClientOptions
	actionConfig *action.Configuration
	restConfig   *rest.Config
	mu           sync.Mutex
	optimizer    *RepositoryOptimizer
}

// NewClient creates a new Helm client with the given options.
func NewClient(options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	client := &Client{
		options: *options,
	}

	// Create Helm settings
	settings := cli.New()
	if options.KubeconfigPath != "" {
		// Check if kubeconfig file exists before using it
		if _, err := os.Stat(options.KubeconfigPath); err == nil {
			settings.KubeConfig = options.KubeconfigPath
			logrus.Debugf("Helm using kubeconfig: %s", options.KubeconfigPath)
		} else {
			logrus.Debugf("Helm kubeconfig not found at %s, will attempt in-cluster config", options.KubeconfigPath)
		}
	} else {
		logrus.Debug("Helm kubeconfig path is empty, attempting in-cluster configuration")
	}
	if options.Namespace != "" {
		settings.SetNamespace(options.Namespace)
	}
	client.settings = settings

	// Create action configuration
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		settings.Namespace(),
		os.Getenv("HELM_DRIVER"),
		logrus.Debugf,
	); err != nil {
		return nil, fmt.Errorf("failed to initialize Helm action configuration: %w", err)
	}
	client.actionConfig = actionConfig

	// Store REST configuration if provided
	if options.RESTConfig != nil {
		client.restConfig = options.RESTConfig
	}

	// Initialize optimizer from options
	if options.Optimizer != nil {
		client.optimizer = options.Optimizer
	} else {
		client.optimizer = NewRepositoryOptimizer(nil, 300, 3, false)
	}

	return client, nil
}

// ListReleases returns a list of releases based on the provided options.
func (c *Client) ListReleases(allNamespaces bool, namespace string) ([]*release.Release, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation if specified and not using all namespaces
	if namespace != "" && !allNamespaces {
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	client := action.NewList(c.actionConfig)
	client.AllNamespaces = allNamespaces

	// If listing all namespaces, ensure we have cluster-wide access
	if allNamespaces {
		// Re-initialize with empty namespace to access all namespaces
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			"", // Empty namespace for cluster-wide access
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace
			if originalNamespace != "" {
				c.settings.SetNamespace(originalNamespace)
				if err := c.actionConfig.Init(
					c.settings.RESTClientGetter(),
					originalNamespace,
					os.Getenv("HELM_DRIVER"),
					logrus.Debugf,
				); err != nil {
					logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
				}
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration for cluster-wide access: %w", err)
		}
	}

	releases, err := client.Run()

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	} else if allNamespaces && originalNamespace != "" {
		// Restore original namespace when listing all namespaces
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	return releases, err
}

// GetRelease returns a specific release by name in the specified namespace.
func (c *Client) GetRelease(name string, namespace string) (*release.Release, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation
	if namespace != "" {
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	client := action.NewGet(c.actionConfig)
	release, err := client.Run(name)

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	return release, err
}

// InstallRelease installs a Helm chart with the given values in the specified namespace.
func (c *Client) InstallRelease(name string, chartPath string, namespace string, values map[string]interface{}) (*release.Release, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logrus.Debugf("Installing Helm release '%s' from chart '%s' in namespace '%s'", name, chartPath, namespace)

	client := action.NewInstall(c.actionConfig)
	client.ReleaseName = name
	client.Namespace = namespace
	client.CreateNamespace = true

	// Load chart
	logrus.Debugf("Locating chart '%s'", chartPath)
	chart, err := client.LocateChart(chartPath, c.settings)
	if err != nil {
		logrus.Errorf("Failed to locate chart '%s': %v", chartPath, err)
		return nil, fmt.Errorf("could not locate chart '%s': %w", chartPath, err)
	}
	logrus.Debugf("Located chart at path '%s'", chart)

	// Load chart
	logrus.Debugf("Loading chart from path '%s'", chart)
	chartRequested, err := loader.Load(chart)
	if err != nil {
		logrus.Errorf("Failed to load chart from path '%s': %v", chart, err)
		return nil, fmt.Errorf("failed to load chart from path '%s': %w", chart, err)
	}
	logrus.Debugf("Loaded chart '%s' version '%s'", chartRequested.Name(), chartRequested.Metadata.Version)

	// Install release
	logrus.Debugf("Installing release '%s' in namespace '%s'", name, namespace)
	release, err := client.Run(chartRequested, values)
	if err != nil {
		logrus.Errorf("Failed to install release '%s' from chart '%s' in namespace '%s': %v", name, chartPath, namespace, err)
		return nil, fmt.Errorf("failed to install release '%s' from chart '%s' in namespace '%s': %w", name, chartPath, namespace, err)
	}

	logrus.Infof("Successfully installed release '%s' version %d in namespace '%s'", name, release.Version, namespace)
	return release, nil
}

// UpgradeRelease upgrades a Helm release with the given values in the specified namespace.
func (c *Client) UpgradeRelease(name string, chartPath string, namespace string, values map[string]interface{}) (*release.Release, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logrus.Debugf("Upgrading Helm release '%s' with chart '%s' in namespace '%s'", name, chartPath, namespace)

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation
	if namespace != "" {
		logrus.Debugf("Setting namespace to '%s' for upgrade operation", namespace)
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			logrus.Errorf("Failed to reinitialize Helm action configuration for namespace '%s': %v", namespace, err)
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration for namespace '%s': %w", namespace, err)
		}
	}

	client := action.NewUpgrade(c.actionConfig)
	client.Namespace = namespace

	// Load chart
	logrus.Debugf("Locating chart '%s'", chartPath)
	chart, err := client.LocateChart(chartPath, c.settings)
	if err != nil {
		logrus.Errorf("Failed to locate chart '%s': %v", chartPath, err)
		// Restore original namespace before returning error
		if namespace != "" && originalNamespace != namespace {
			logrus.Debugf("Restoring original namespace '%s'", originalNamespace)
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
		}
		return nil, fmt.Errorf("could not locate chart '%s': %w", chartPath, err)
	}
	logrus.Debugf("Located chart at path '%s'", chart)

	// Load chart
	logrus.Debugf("Loading chart from path '%s'", chart)
	chartRequested, err := loader.Load(chart)
	if err != nil {
		logrus.Errorf("Failed to load chart from path '%s': %v", chart, err)
		// Restore original namespace before returning error
		if namespace != "" && originalNamespace != namespace {
			logrus.Debugf("Restoring original namespace '%s'", originalNamespace)
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
		}
		return nil, fmt.Errorf("failed to load chart from path '%s': %w", chart, err)
	}
	logrus.Debugf("Loaded chart '%s' version '%s'", chartRequested.Name(), chartRequested.Metadata.Version)

	// Upgrade release
	logrus.Debugf("Upgrading release '%s' in namespace '%s'", name, namespace)
	release, err := client.Run(name, chartRequested, values)
	if err != nil {
		logrus.Errorf("Failed to upgrade release '%s' with chart '%s' in namespace '%s': %v", name, chartPath, namespace, err)
		// Restore original namespace before returning error
		if namespace != "" && originalNamespace != namespace {
			logrus.Debugf("Restoring original namespace '%s'", originalNamespace)
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
		}
		return nil, fmt.Errorf("failed to upgrade release '%s' with chart '%s' in namespace '%s': %w", name, chartPath, namespace, err)
	}

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		logrus.Debugf("Restoring original namespace '%s'", originalNamespace)
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	logrus.Infof("Successfully upgraded release '%s' to version %d in namespace '%s'", name, release.Version, namespace)
	return release, nil
}

// UninstallRelease uninstalls a Helm release in the specified namespace.
func (c *Client) UninstallRelease(name string, namespace string) (*release.UninstallReleaseResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logrus.Debugf("Uninstalling Helm release '%s' in namespace '%s'", name, namespace)

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation
	if namespace != "" {
		logrus.Debugf("Setting namespace to '%s' for uninstall operation", namespace)
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			logrus.Errorf("Failed to reinitialize Helm action configuration for namespace '%s': %v", namespace, err)
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration for namespace '%s': %w", namespace, err)
		}
	}

	client := action.NewUninstall(c.actionConfig)
	logrus.Debugf("Running uninstall for release '%s' in namespace '%s'", name, namespace)
	result, err := client.Run(name)
	if err != nil {
		logrus.Errorf("Failed to uninstall release '%s' in namespace '%s': %v", name, namespace, err)
		// Restore original namespace before returning error
		if namespace != "" && originalNamespace != namespace {
			logrus.Debugf("Restoring original namespace '%s'", originalNamespace)
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
		}
		return nil, fmt.Errorf("failed to uninstall release '%s' in namespace '%s': %w", name, namespace, err)
	}

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		logrus.Debugf("Restoring original namespace '%s'", originalNamespace)
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if result != nil && result.Release != nil {
		logrus.Infof("Successfully uninstalled release '%s' version %d in namespace '%s'", name, result.Release.Version, namespace)
	} else {
		logrus.Infof("Successfully uninstalled release '%s' in namespace '%s'", name, namespace)
	}
	return result, nil
}

// RollbackRelease rolls back a Helm release to a specific revision in the specified namespace.
func (c *Client) RollbackRelease(name string, revision int, namespace string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation
	if namespace != "" {
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	client := action.NewRollback(c.actionConfig)
	client.Version = revision
	err := client.Run(name)

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	return err
}

// GetReleaseHistory returns the revision history of a release.
func (c *Client) GetReleaseHistory(name string) ([]*release.Release, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client := action.NewHistory(c.actionConfig)
	client.Max = 20
	return client.Run(name)
}

// ListRepositories lists the configured Helm repositories.
func (c *Client) ListRepositories() ([]*Repository, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	repoFile := c.settings.RepositoryConfig

	// Ensure the repositories file exists
	if _, err := os.Stat(repoFile); os.IsNotExist(err) {
		// Create parent directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(repoFile), 0755); err != nil {
			return nil, fmt.Errorf("failed to create repository directory: %w", err)
		}
		// Create empty repositories file
		if err := repo.NewFile().WriteFile(repoFile, 0644); err != nil {
			return nil, fmt.Errorf("failed to create repository file: %w", err)
		}
	}

	// Read repositories file
	repositories, err := repo.LoadFile(repoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load repositories file: %w", err)
	}

	result := make([]*Repository, 0, len(repositories.Repositories))
	for _, r := range repositories.Repositories {
		if r != nil {
			result = append(result, &Repository{
				Name: r.Name,
				URL:  r.URL,
			})
		}
	}
	return result, nil
}

// AddRepository adds a new Helm repository.
func (c *Client) AddRepository(name, url string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	repoFile := c.settings.RepositoryConfig

	// Ensure the repositories file exists
	if _, err := os.Stat(repoFile); os.IsNotExist(err) {
		// Create parent directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(repoFile), 0755); err != nil {
			return fmt.Errorf("failed to create repository directory: %w", err)
		}
		// Create empty repositories file
		if err := repo.NewFile().WriteFile(repoFile, 0644); err != nil {
			return fmt.Errorf("failed to create repository file: %w", err)
		}
	}

	// Read repositories file
	repositories, err := repo.LoadFile(repoFile)
	if err != nil {
		return fmt.Errorf("failed to load repositories file: %w", err)
	}

	// Add repository
	newRepo := &repo.Entry{
		Name: name,
		URL:  url,
	}

	// Check if repository already exists
	for _, r := range repositories.Repositories {
		if r.Name == name {
			return fmt.Errorf("repository %s already exists", name)
		}
	}

	repositories.Add(newRepo)

	// Save repositories file
	if err := repositories.WriteFile(repoFile, 0644); err != nil {
		return fmt.Errorf("failed to save repositories file: %w", err)
	}

	return nil
}

// RemoveRepository removes a Helm repository.
func (c *Client) RemoveRepository(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	repoFile := c.settings.RepositoryConfig

	// Read repositories file
	repositories, err := repo.LoadFile(repoFile)
	if err != nil {
		return fmt.Errorf("failed to load repositories file: %w", err)
	}

	// Remove repository
	repositories.Remove(name)

	// Save repositories file
	if err := repositories.WriteFile(repoFile, 0644); err != nil {
		return fmt.Errorf("failed to save repositories file: %w", err)
	}

	return nil
}

// ReleaseToMap converts a Helm release to a map for JSON serialization.
func ReleaseToMap(r *release.Release) map[string]interface{} {
	if r == nil {
		return nil
	}

	result := map[string]interface{}{
		"name":         r.Name,
		"namespace":    r.Namespace,
		"revision":     r.Version,
		"status":       r.Info.Status.String(),
		"last_updated": r.Info.LastDeployed,
	}

	if r.Chart != nil {
		result["chart"] = ChartToMap(r.Chart)
	}

	if r.Info != nil {
		result["info"] = map[string]interface{}{
			"first_deployed": r.Info.FirstDeployed,
			"last_deployed":  r.Info.LastDeployed,
			"deleted":        r.Info.Deleted,
			"description":    r.Info.Description,
			"status":         r.Info.Status.String(),
			"notes":          r.Info.Notes,
		}
	}

	return result
}

// ChartToMap converts a Helm chart to a map for JSON serialization.
func ChartToMap(c *chart.Chart) map[string]interface{} {
	if c == nil {
		return nil
	}

	result := map[string]interface{}{
		"name":        c.Name(),
		"version":     c.Metadata.Version,
		"app_version": c.Metadata.AppVersion,
		"description": c.Metadata.Description,
	}

	if len(c.Metadata.Maintainers) > 0 {
		maintainers := make([]string, len(c.Metadata.Maintainers))
		for i, m := range c.Metadata.Maintainers {
			maintainers[i] = m.Name
		}
		result["maintainers"] = maintainers
	}

	return result
}

// Repository represents a simplified Helm repository.
type Repository struct {
	Name string
	URL  string
}

// ListReleasesAsMap returns releases as maps for easy JSON serialization.
func (c *Client) ListReleasesAsMap(allNamespaces bool, namespace string) ([]map[string]interface{}, error) {
	releases, err := c.ListReleases(allNamespaces, namespace)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(releases))
	for _, rel := range releases {
		if rel != nil {
			result = append(result, ReleaseToMap(rel))
		}
	}
	return result, nil
}

// GetReleaseAsMap returns a release as a map for easy JSON serialization.
func (c *Client) GetReleaseAsMap(name string, namespace string) (map[string]interface{}, error) {
	rel, err := c.GetRelease(name, namespace)
	if err != nil {
		return nil, err
	}
	return ReleaseToMap(rel), nil
}

// GetReleaseValuesAsMap returns the values of a release as a map.
func (c *Client) GetReleaseValuesAsMap(name string, namespace string, all bool) (map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation
	if namespace != "" {
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	getCmd := action.NewGetValues(c.actionConfig)
	getCmd.AllValues = all
	values, err := getCmd.Run(name)

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	return values, err
}

// GetReleaseManifest returns the rendered manifest of a release in the specified namespace.
func (c *Client) GetReleaseManifest(name string, namespace string) (string, error) {
	rel, err := c.GetRelease(name, namespace)
	if err != nil {
		return "", err
	}
	return rel.Manifest, nil
}

// GetReleaseHistoryAsMap returns the release history as a map slice.
func (c *Client) GetReleaseHistoryAsMap(name string, namespace string, max int) ([]map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation
	if namespace != "" {
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	client := action.NewHistory(c.actionConfig)
	if max > 0 {
		client.Max = max
	}
	releases, err := client.Run(name)

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(releases))
	for _, rel := range releases {
		if rel != nil {
			result = append(result, ReleaseToMap(rel))
		}
	}
	return result, nil
}

// SearchChartsAsMap searches for charts in configured repositories.
func (c *Client) SearchChartsAsMap(keyword string, devel bool) ([]map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Load the repositories.yaml file
	repoFile := c.settings.RepositoryConfig
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("failed to load repository file: %w", err)
	}

	if len(f.Repositories) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Create a search index
	index := search.NewIndex()

	// Process each repository
	for _, re := range f.Repositories {
		n := re.Name
		// Get the index file path for this repository
		indexFilePath := filepath.Join(c.settings.RepositoryCache, helmpath.CacheIndexFile(n))

		// Load the index file
		ind, err := repo.LoadIndexFile(indexFilePath)
		if err != nil {
			logrus.Warnf("Repo %q is corrupt or missing. Try 'helm repo update'. Error: %v", n, err)
			continue
		}

		// Add repository to index
		index.AddRepo(n, ind, devel)
	}

	// Perform search
	var results []*search.Result
	if keyword == "" {
		results = index.All()
	} else {
		results, err = index.Search(keyword, 100, false) // Limit to 100 results
		if err != nil {
			return nil, fmt.Errorf("failed to search charts: %w", err)
		}
	}

	// Convert results to map format
	result := make([]map[string]interface{}, 0, len(results))
	for _, r := range results {
		if r.Chart != nil {
			chartInfo := map[string]interface{}{
				"name":        r.Chart.Name,
				"version":     r.Chart.Version,
				"app_version": r.Chart.AppVersion,
				"description": r.Chart.Description,
				"repository":  strings.Split(r.Name, "/")[0], // Extract repository name from full chart name
			}

			// Apply version filtering based on devel flag
			if !devel {
				// For stable releases, skip prereleases unless explicitly requested
				if v, err := semver.NewVersion(r.Chart.Version); err == nil {
					if v.Prerelease() != "" {
						continue // Skip prerelease versions
					}
				}
			}

			result = append(result, chartInfo)
		}
	}

	return result, nil
}

// GetChartInfoAsMap returns chart information as a map.
func (c *Client) GetChartInfoAsMap(chartRef string) (map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(chartRef); err != nil {
		return nil, fmt.Errorf("chart not found: %s", chartRef)
	}
	chartPath := chartRef
	chrt, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}
	info := map[string]interface{}{
		"name":        chrt.Metadata.Name,
		"version":     chrt.Metadata.Version,
		"app_version": chrt.Metadata.AppVersion,
		"description": chrt.Metadata.Description,
		"home":        chrt.Metadata.Home,
	}
	if len(chrt.Metadata.Maintainers) > 0 {
		maintainers := make([]map[string]string, len(chrt.Metadata.Maintainers))
		for i, m := range chrt.Metadata.Maintainers {
			maintainers[i] = map[string]string{"name": m.Name, "email": m.Email}
		}
		info["maintainers"] = maintainers
	}
	return info, nil
}

// GetReleaseStatusAsMap returns the status of a release as a map.
func (c *Client) GetReleaseStatusAsMap(name string, namespace string) (map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save original namespace
	originalNamespace := c.settings.Namespace()

	// Temporarily set the namespace for this operation
	if namespace != "" {
		c.settings.SetNamespace(namespace)
		// Re-initialize action config with new namespace
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			// Restore original namespace before returning error
			c.settings.SetNamespace(originalNamespace)
			if err := c.actionConfig.Init(
				c.settings.RESTClientGetter(),
				originalNamespace,
				os.Getenv("HELM_DRIVER"),
				logrus.Debugf,
			); err != nil {
				logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
			}
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	client := action.NewStatus(c.actionConfig)
	rel, err := client.Run(name)

	// Restore original namespace
	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":      rel.Name,
		"namespace": rel.Namespace,
		"status":    rel.Info.Status.String(),
		"version":   rel.Version,
		"notes":     rel.Info.Notes,
	}, nil
}

// TemplateChart renders a chart template without installing it.
func (c *Client) TemplateChart(name, chartRef, namespace, valuesFile string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create install client for templating
	client := action.NewInstall(c.actionConfig)

	// Configure for dry-run (template only)
	client.DryRun = true
	client.ReleaseName = name
	client.Replace = true    // Skip name check
	client.ClientOnly = true // Client-side templating only

	// Set namespace if provided
	if namespace != "" {
		client.Namespace = namespace
	} else {
		client.Namespace = c.settings.Namespace()
	}

	// Locate chart
	chartPath, err := client.LocateChart(chartRef, c.settings)
	if err != nil {
		return "", fmt.Errorf("failed to locate chart: %w", err)
	}

	// Load chart
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return "", fmt.Errorf("failed to load chart: %w", err)
	}

	// Load values if values file is provided
	var values map[string]interface{}
	if valuesFile != "" {
		// Load values from file
		vals, err := chartutil.ReadValuesFile(valuesFile)
		if err != nil {
			return "", fmt.Errorf("failed to load values file: %w", err)
		}
		values = vals.AsMap()
	} else {
		// Use empty values
		values = make(map[string]interface{})
	}

	// Run install with dry-run to generate templates
	rel, err := client.Run(chartRequested, values)
	if err != nil {
		return "", fmt.Errorf("failed to template chart: %w", err)
	}

	return rel.Manifest, nil
}

// UpdateRepositories updates all Helm repositories by downloading the latest index files.
func (c *Client) UpdateRepositories() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Load the repositories.yaml file
	repoFile := c.settings.RepositoryConfig
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no repositories configured. Add repositories before updating. Use 'helm_add_repository' tool to add repositories")
		}
		return fmt.Errorf("failed to load repository file: %w", err)
	}

	if len(f.Repositories) == 0 {
		logrus.Info("No repositories configured. Nothing to update. Use 'helm_add_repository' tool to add repositories")
		return nil
	}

	// Ensure repository cache directory exists
	if err := os.MkdirAll(c.settings.RepositoryCache, 0755); err != nil {
		return fmt.Errorf("failed to create repository cache directory: %w", err)
	}

	// Create chart repositories and download index files
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(c.settings))
		if err != nil {
			logrus.Warnf("Failed to create chart repository for %s: %v", cfg.Name, err)
			continue
		}

		// Set cache path
		r.CachePath = c.settings.RepositoryCache
		repos = append(repos, r)
	}

	if len(repos) == 0 {
		return fmt.Errorf("failed to create any chart repositories. Check if repositories are properly configured")
	}

	// Update each repository concurrently with individual timeouts
	logrus.Infof("Updating %d Helm chart repositories...", len(repos))
	errChan := make(chan error, len(repos))
	successCount := 0

	// Use a wait group to ensure all goroutines complete
	var wg sync.WaitGroup

	for _, re := range repos {
		wg.Add(1)
		go func(r *repo.ChartRepository) {
			defer wg.Done()
			logrus.Debugf("Updating repository: %s (%s)", r.Config.Name, r.Config.URL)

			// Create a context with 2 minute timeout for each repository update
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			// Create channels for result and error
			resultChan := make(chan string, 1)
			errChanInternal := make(chan error, 1)

			// Run the download in a goroutine
			go func() {
				result, err := r.DownloadIndexFile()
				if err != nil {
					errChanInternal <- err
				} else {
					resultChan <- result
				}
			}()

			// Wait for either completion or timeout
			select {
			case <-ctx.Done():
				logrus.Errorf("failed to update repository %q: timeout after 2 minutes", r.Config.Name)
				errChan <- fmt.Errorf("failed to update repository %q: timeout after 2 minutes", r.Config.Name)
			case err := <-errChanInternal:
				logrus.Errorf("failed to update repository %q: %v", r.Config.Name, err)
				errChan <- fmt.Errorf("failed to update repository %q: %v", r.Config.Name, err)
			case <-resultChan:
				logrus.Debugf("Successfully updated repository: %s", r.Config.Name)
				errChan <- nil
			}
		}(re)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Wait for all updates to complete
	for err := range errChan {
		if err != nil {
			logrus.Errorf("Repository update error: %v", err)
		} else {
			successCount++
		}
	}

	logrus.Infof("Repository update complete. Successfully updated %d/%d repositories.", successCount, len(repos))

	if successCount == 0 && len(repos) > 0 {
		return fmt.Errorf("failed to update any repositories. Check network connectivity and repository URLs")
	}

	return nil
}

// GetMirrorConfiguration returns information about configured mirrors.
func (c *Client) GetMirrorConfiguration() map[string]interface{} {
	if c.optimizer == nil {
		return map[string]interface{}{
			"enabled": false,
			"mirrors": map[string]string{},
		}
	}

	return map[string]interface{}{
		"enabled":     c.optimizer.IsMirrorEnabled(),
		"mirrors":     c.optimizer.ListMirrors(),
		"timeout_sec": int(c.optimizer.GetTimeout().Seconds()),
		"max_retries": c.optimizer.GetMaxRetry(),
	}
}

// ListReleasesPaginated lists Helm releases with pagination support
func (c *Client) ListReleasesPaginated(limit int, continueToken, namespace, status string) ([]map[string]interface{}, bool, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create list action
	list := action.NewList(c.actionConfig)
	list.All = true // List all namespaces, then filter if namespace is specified
	list.Limit = limit

	// Apply status filter if specified
	if status != "" {
		list.Filter = status
	}

	// Convert continue token to offset (simple implementation)
	offset := 0
	if continueToken != "" {
		if parsed, err := fmt.Sscanf(continueToken, "offset:%d", &offset); err != nil || parsed != 1 {
			offset = 0
		}
	}

	// Get all releases
	releases, err := list.Run()
	if err != nil {
		return nil, false, "", fmt.Errorf("failed to list releases: %w", err)
	}

	var filteredReleases []*release.Release
	for _, rel := range releases {
		// Filter by namespace if specified
		if namespace != "" && rel.Namespace != namespace {
			continue
		}

		// Status filter is already applied via list.Filter if specified
		filteredReleases = append(filteredReleases, rel)
	}

	// Apply pagination offset and limit
	totalCount := len(filteredReleases)
	startIdx := offset
	endIdx := offset + limit

	if startIdx >= totalCount {
		filteredReleases = []*release.Release{}
	} else if endIdx > totalCount {
		filteredReleases = filteredReleases[startIdx:]
	} else {
		filteredReleases = filteredReleases[startIdx:endIdx]
	}

	// Convert to summary format to reduce size
	summaries := make([]map[string]interface{}, len(filteredReleases))
	for i, rel := range filteredReleases {
		summaries[i] = map[string]interface{}{
			"name":       rel.Name,
			"namespace":  rel.Namespace,
			"status":     rel.Info.Status.String(),
			"chart":      rel.Chart.Metadata.Name,
			"version":    rel.Chart.Metadata.Version,
			"appVersion": rel.Chart.Metadata.AppVersion,
			"updated":    rel.Info.LastDeployed.Format(time.RFC3339),
			"labels":     rel.Labels,
		}
	}

	// Calculate pagination info
	hasMore := endIdx < totalCount
	var newContinueToken string
	if hasMore {
		newContinueToken = fmt.Sprintf("offset:%d", endIdx)
	}

	return summaries, hasMore, newContinueToken, nil
}

// FilterReleasesByLabels filters releases by specified label keys
func (c *Client) FilterReleasesByLabels(releases []map[string]interface{}, includeLabels string) []map[string]interface{} {
	if includeLabels == "" {
		return releases
	}

	labelKeys := strings.Split(includeLabels, ",")
	for i, key := range labelKeys {
		labelKeys[i] = strings.TrimSpace(key)
	}

	// Rebuild releases with only specified labels
	filtered := make([]map[string]interface{}, len(releases))
	for i, release := range releases {
		newRelease := make(map[string]interface{})

		// Copy all non-label fields
		for k, v := range release {
			if k != "labels" {
				newRelease[k] = v
			}
		}

		// Filter labels
		if labels, ok := release["labels"].(map[string]interface{}); ok {
			filteredLabels := make(map[string]interface{})
			for _, key := range labelKeys {
				if val, exists := labels[key]; exists {
					filteredLabels[key] = val
				}
			}
			newRelease["labels"] = filteredLabels
		}

		filtered[i] = newRelease
	}

	return filtered
}

// GetReleaseStatus gets minimal release status information
func (c *Client) GetReleaseStatus(name, namespace string) (map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create status action
	status := action.NewStatus(c.actionConfig)

	// Get release status
	rel, err := status.Run(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get release status: %w", err)
	}

	return map[string]interface{}{
		"name":        rel.Name,
		"namespace":   rel.Namespace,
		"status":      rel.Info.Status.String(),
		"description": rel.Info.Description,
		"version":     rel.Version,
		"updated":     rel.Info.LastDeployed.Format(time.RFC3339),
		"notes":       rel.Info.Notes,
	}, nil
}

// GetReleaseHistoryPaginated gets release history with pagination
func (c *Client) GetReleaseHistoryPaginated(name, namespace string, limit int, continueToken string, includeStatus bool) ([]map[string]interface{}, bool, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create history action
	hist := action.NewHistory(c.actionConfig)

	// Get all history
	releases, err := hist.Run(name)
	if err != nil {
		return nil, false, "", fmt.Errorf("failed to get release history: %w", err)
	}

	// Convert to summary format
	summaries := make([]map[string]interface{}, len(releases))
	for i, rel := range releases {
		summary := map[string]interface{}{
			"revision":   rel.Version,
			"updated":    rel.Info.LastDeployed.Format(time.RFC3339),
			"status":     rel.Info.Status.String(),
			"chart":      rel.Chart.Metadata.Name,
			"version":    rel.Chart.Metadata.Version,
			"appVersion": rel.Chart.Metadata.AppVersion,
		}

		if includeStatus {
			summary["description"] = rel.Info.Description
			summary["notes"] = rel.Info.Notes
		}

		summaries[i] = summary
	}

	// Simple pagination
	totalCount := len(summaries)
	offset := 0
	if continueToken != "" {
		if parsed, err := fmt.Sscanf(continueToken, "offset:%d", &offset); err != nil || parsed != 1 {
			offset = 0
		}
	}

	startIdx := offset
	endIdx := offset + limit

	if startIdx >= totalCount {
		summaries = []map[string]interface{}{}
	} else if endIdx > totalCount {
		summaries = summaries[startIdx:]
	} else {
		summaries = summaries[startIdx:endIdx]
	}

	hasMore := endIdx < totalCount
	var newContinueToken string
	if hasMore {
		newContinueToken = fmt.Sprintf("offset:%d", endIdx)
	}

	return summaries, hasMore, newContinueToken, nil
}

// GetFailedReleases gets releases with failed or pending status
func (c *Client) GetFailedReleases(limit int, namespace string, includePending bool) ([]map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create list action
	list := action.NewList(c.actionConfig)
	list.All = true

	// Get all releases
	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	var failedReleases []map[string]interface{}
	for _, rel := range releases {
		// Filter by namespace if specified
		if namespace != "" && rel.Namespace != namespace {
			continue
		}

		// Check status
		status := rel.Info.Status.String()
		if status == "failed" || (includePending && (status == "pending-install" || status == "pending-upgrade" || status == "pending-rollback")) {
			failedRelease := map[string]interface{}{
				"name":        rel.Name,
				"namespace":   rel.Namespace,
				"status":      status,
				"chart":       rel.Chart.Metadata.Name,
				"version":     rel.Chart.Metadata.Version,
				"appVersion":  rel.Chart.Metadata.AppVersion,
				"updated":     rel.Info.LastDeployed.Format(time.RFC3339),
				"description": rel.Info.Description,
				"error":       rel.Info.Description, // Description often contains error info for failed releases
			}

			if rel.Info.Notes != "" {
				failedRelease["notes"] = rel.Info.Notes
			}

			failedReleases = append(failedReleases, failedRelease)

			// Limit results
			if len(failedReleases) >= limit {
				break
			}
		}
	}

	return failedReleases, nil
}

// GetClusterOverview gets cluster-wide Helm overview
func (c *Client) GetClusterOverview(includeNodes, includeStorage bool) (map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create list action
	list := action.NewList(c.actionConfig)
	list.All = true

	// Get all releases
	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	// Count by status and namespace
	statusCounts := make(map[string]int)
	namespaceCounts := make(map[string]int)
	chartCounts := make(map[string]int)

	for _, rel := range releases {
		status := rel.Info.Status.String()
		statusCounts[status]++

		namespaceCounts[rel.Namespace]++

		chartKey := fmt.Sprintf("%s:%s", rel.Chart.Metadata.Name, rel.Chart.Metadata.Version)
		chartCounts[chartKey]++
	}

	overview := map[string]interface{}{
		"totalReleases":      len(releases),
		"statusBreakdown":    statusCounts,
		"namespaceBreakdown": namespaceCounts,
		"chartBreakdown":     chartCounts,
		"generatedAt":        time.Now().Format(time.RFC3339),
	}

	// Add optional info
	if includeNodes || includeStorage {
		overview["additionalInfo"] = map[string]interface{}{
			"nodesIncluded":   includeNodes,
			"storageIncluded": includeStorage,
			"note":            "Node and storage info would require additional Kubernetes API calls",
		}
	}

	return overview, nil
}

// FindReleasesByLabels finds releases by label selector
func (c *Client) FindReleasesByLabels(labelSelector string, limit int, namespace string) ([]map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Parse simple label selector (key=value format)
	selectors := strings.Split(labelSelector, ",")
	for i, sel := range selectors {
		selectors[i] = strings.TrimSpace(sel)
	}

	// Create list action
	list := action.NewList(c.actionConfig)
	list.All = true

	// Get all releases
	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	var matchingReleases []map[string]interface{}
	for _, rel := range releases {
		// Filter by namespace if specified
		if namespace != "" && rel.Namespace != namespace {
			continue
		}

		// Check label matches
		matched := true
		for _, selector := range selectors {
			parts := strings.Split(selector, "=")
			if len(parts) != 2 {
				continue // Skip invalid selector
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if val, exists := rel.Labels[key]; !exists || val != value {
				matched = false
				break
			}
		}

		if matched {
			releaseSummary := map[string]interface{}{
				"name":       rel.Name,
				"namespace":  rel.Namespace,
				"status":     rel.Info.Status.String(),
				"chart":      rel.Chart.Metadata.Name,
				"version":    rel.Chart.Metadata.Version,
				"appVersion": rel.Chart.Metadata.AppVersion,
				"updated":    rel.Info.LastDeployed.Format(time.RFC3339),
				"labels":     rel.Labels,
			}
			matchingReleases = append(matchingReleases, releaseSummary)

			if len(matchingReleases) >= limit {
				break
			}
		}
	}

	return matchingReleases, nil
}

// GetResourcesOfRelease gets summarized list of Kubernetes resources managed by a Helm release
func (c *Client) GetResourcesOfRelease(name, namespace string, includeStatus bool, limit int) ([]map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create get action to retrieve release
	get := action.NewGet(c.actionConfig)

	// Get release (contains manifest)
	_, err := get.Run(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get release manifest: %w", err)
	}

	// Parse manifest into resources (simplified implementation)
	// In a real implementation, you'd parse the YAML manifest
	// For now, return basic resource information
	resources := []map[string]interface{}{
		{
			"kind":      "Placeholder",
			"name":      "Resource parsing would require YAML parsing implementation",
			"namespace": namespace,
			"status":    "Unknown",
		},
	}

	// Limit results
	if len(resources) > limit {
		resources = resources[:limit]
	}

	return resources, nil
}

// ClearCache clears the Helm cache to force fresh queries
func (c *Client) ClearCache() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cachePath := c.settings.RepositoryCache
	if cachePath == "" {
		cachePath = helmpath.CachePath("")
	}

	if err := os.RemoveAll(cachePath); err != nil {
		return fmt.Errorf("failed to clear Helm cache: %w", err)
	}

	logrus.Info("Helm cache cleared successfully")
	return nil
}

// CacheStats represents cache statistics
type CacheStats struct {
	CachePath  string `json:"cachePath"`
	IndexFiles int    `json:"indexFiles"`
}

// GetCacheStats returns cache statistics
func (c *Client) GetCacheStats() (*CacheStats, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cachePath := c.settings.RepositoryCache
	if cachePath == "" {
		cachePath = helmpath.CachePath("")
	}

	indexFiles := 0
	if err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), "-index.yaml") {
			indexFiles++
		}
		return nil
	}); err != nil {
		indexFiles = -1
	}

	return &CacheStats{
		CachePath:  cachePath,
		IndexFiles: indexFiles,
	}, nil
}

// GetReleaseSummary returns a brief summary of a Helm release
func (c *Client) GetReleaseSummary(name string, namespace string) (*ReleaseSummary, error) {
	originalNamespace := c.settings.Namespace()

	if namespace != "" {
		c.settings.SetNamespace(namespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			c.settings.SetNamespace(originalNamespace)
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	getAction := action.NewGet(c.actionConfig)
	rel, err := getAction.Run(name)

	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("release %s not found: %w", name, err)
	}

	return ExtractReleaseSummary(rel), nil
}

// GetListReleasesSummary returns a list of release summaries
func (c *Client) GetListReleasesSummary(namespace string, limit, offset int) ([]*ReleaseSummary, error) {
	originalNamespace := c.settings.Namespace()

	if namespace != "" {
		c.settings.SetNamespace(namespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			c.settings.SetNamespace(originalNamespace)
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	listAction := action.NewList(c.actionConfig)
	listAction.AllNamespaces = namespace == ""
	if limit > 0 {
		listAction.Limit = limit + offset
	}

	releases, err := listAction.Run()

	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	if offset > 0 && offset < len(releases) {
		releases = releases[offset:]
	}

	return ExtractReleaseSummaries(releases), nil
}

// GetQuickInfo returns a quick overview of all Helm releases
func (c *Client) GetQuickInfo() (map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	list := action.NewList(c.actionConfig)
	list.All = true

	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	statusCounts := make(map[string]int)
	namespaces := make(map[string]int)
	var totalRevisions int

	for _, rel := range releases {
		statusCounts[rel.Info.Status.String()]++
		namespaces[rel.Namespace]++
		totalRevisions += rel.Version
	}

	return map[string]interface{}{
		"totalReleases":   len(releases),
		"totalNamespaces": len(namespaces),
		"statusCounts":    statusCounts,
		"namespaces":      namespaces,
		"totalRevisions":  totalRevisions,
	}, nil
}

// FindReleasesByChart finds releases using a specific chart
func (c *Client) FindReleasesByChart(chartName, chartVersion string, limit int) ([]map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	list := action.NewList(c.actionConfig)
	list.All = true

	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	var matchingReleases []map[string]interface{}
	for _, rel := range releases {
		currentChart := rel.Chart.Metadata.Name
		currentVersion := rel.Chart.Metadata.Version

		if currentChart == chartName {
			if chartVersion == "" || currentVersion == chartVersion {
				releaseSummary := map[string]interface{}{
					"name":      rel.Name,
					"namespace": rel.Namespace,
					"chart":     currentChart + ":" + currentVersion,
					"version":   rel.Version,
					"status":    rel.Info.Status.String(),
					"updated":   rel.Info.LastDeployed.Format(time.RFC3339),
				}
				matchingReleases = append(matchingReleases, releaseSummary)

				if len(matchingReleases) >= limit {
					break
				}
			}
		}
	}

	return matchingReleases, nil
}

// FindBrokenReleases finds releases with failed or pending status
func (c *Client) FindBrokenReleases(namespace string, limit int) ([]map[string]interface{}, error) {
	originalNamespace := c.settings.Namespace()

	if namespace != "" {
		c.settings.SetNamespace(namespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			c.settings.SetNamespace(originalNamespace)
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	listAction := action.NewList(c.actionConfig)
	listAction.AllNamespaces = namespace == ""
	if limit > 0 {
		listAction.Limit = limit
	}

	releases, err := listAction.Run()

	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	var brokenReleases []map[string]interface{}
	for _, rel := range releases {
		status := rel.Info.Status.String()
		if status == "failed" || status == "pending" || status == "uninstalling" {
			releaseInfo := map[string]interface{}{
				"name":      rel.Name,
				"namespace": rel.Namespace,
				"chart":     rel.Chart.Metadata.Name + ":" + rel.Chart.Metadata.Version,
				"version":   rel.Version,
				"status":    status,
				"updated":   rel.Info.LastDeployed.Format(time.RFC3339),
				"message":   rel.Info.Description,
			}
			brokenReleases = append(brokenReleases, releaseInfo)
		}
	}

	return brokenReleases, nil
}

// ValidateRelease validates a release configuration
func (c *Client) ValidateRelease(name string, namespace string) (map[string]interface{}, error) {
	originalNamespace := c.settings.Namespace()

	if namespace != "" {
		c.settings.SetNamespace(namespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			c.settings.SetNamespace(originalNamespace)
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	getAction := action.NewGet(c.actionConfig)
	rel, err := getAction.Run(name)

	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	result := map[string]interface{}{
		"name":      name,
		"namespace": namespace,
		"valid":     true,
		"errors":    []string{},
		"warnings":  []string{},
	}

	if err != nil {
		result["valid"] = false
		result["errors"] = append(result["errors"].([]string), fmt.Sprintf("release not found: %v", err))
		return result, nil
	}

	status := rel.Info.Status.String()
	if status == "failed" {
		result["valid"] = false
		result["warnings"] = append(result["warnings"].([]string), "Release is in failed status")
	}

	result["chart"] = rel.Chart.Metadata.Name + ":" + rel.Chart.Metadata.Version
	result["version"] = rel.Version
	result["status"] = status

	return result, nil
}

// ListReleasesByNamespace lists all releases in a specific namespace
func (c *Client) ListReleasesByNamespace(namespace string) ([]*ReleaseSummary, error) {
	originalNamespace := c.settings.Namespace()

	if namespace != "" {
		c.settings.SetNamespace(namespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			c.settings.SetNamespace(originalNamespace)
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	listAction := action.NewList(c.actionConfig)

	releases, err := listAction.Run()

	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list releases in namespace %s: %w", namespace, err)
	}

	return ExtractReleaseSummaries(releases), nil
}

// CompareReleaseVersions compares two release revisions
func (c *Client) CompareReleaseVersions(name, namespace string, revision1, revision2 int) (map[string]interface{}, error) {
	originalNamespace := c.settings.Namespace()

	if namespace != "" {
		c.settings.SetNamespace(namespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			namespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			c.settings.SetNamespace(originalNamespace)
			return nil, fmt.Errorf("failed to reinitialize Helm action configuration: %w", err)
		}
	}

	listAction := action.NewList(c.actionConfig)
	listAction.AllNamespaces = namespace == ""

	releases, err := listAction.Run()

	if namespace != "" && originalNamespace != namespace {
		c.settings.SetNamespace(originalNamespace)
		if err := c.actionConfig.Init(
			c.settings.RESTClientGetter(),
			originalNamespace,
			os.Getenv("HELM_DRIVER"),
			logrus.Debugf,
		); err != nil {
			logrus.Errorf("Failed to restore original Helm action configuration: %v", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get release history: %w", err)
	}

	var rev1, rev2 *release.Release
	for _, rel := range releases {
		if rel.Name == name && rel.Namespace == namespace {
			if rel.Version == revision1 {
				rev1 = rel
			}
			if rel.Version == revision2 {
				rev2 = rel
			}
		}
	}

	result := map[string]interface{}{
		"release":   name,
		"namespace": namespace,
	}

	if rev1 == nil {
		result["error"] = fmt.Sprintf("revision %d not found", revision1)
		return result, nil
	}
	if rev2 == nil {
		result["error"] = fmt.Sprintf("revision %d not found", revision2)
		return result, nil
	}

	rev1Map := map[string]interface{}{
		"version":    rev1.Version,
		"status":     rev1.Info.Status.String(),
		"chart":      rev1.Chart.Metadata.Name + ":" + rev1.Chart.Metadata.Version,
		"appVersion": rev1.Chart.Metadata.AppVersion,
		"updated":    rev1.Info.LastDeployed.Format(time.RFC3339),
	}

	rev2Map := map[string]interface{}{
		"version":    rev2.Version,
		"status":     rev2.Info.Status.String(),
		"chart":      rev2.Chart.Metadata.Name + ":" + rev2.Chart.Metadata.Version,
		"appVersion": rev2.Chart.Metadata.AppVersion,
		"updated":    rev2.Info.LastDeployed.Format(time.RFC3339),
	}

	result["revision_1"] = rev1Map
	result["revision_2"] = rev2Map

	if rev1.Chart.Metadata.Name != rev2.Chart.Metadata.Name || rev1.Chart.Metadata.Version != rev2.Chart.Metadata.Version {
		result["chart_changed"] = true
	} else {
		result["chart_changed"] = false
	}

	return result, nil
}
