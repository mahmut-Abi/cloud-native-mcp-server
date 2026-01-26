// Package client provides Kubernetes API client functionality with caching and optimization.
// It offers high-level operations for interacting with Kubernetes clusters through
// dynamic and typed clients with built-in GroupVersionResource (GVR) caching.
package client

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned"
)

// ClientOptions holds configuration parameters for creating a Kubernetes client.
// It allows customization of timeouts, rate limiting, and caching behavior.
type ClientOptions struct {
	KubeconfigPath string        // Path to kubeconfig file (empty for default)
	Timeout        time.Duration // API request timeout
	QPS            float32       // Queries per second rate limit
	Burst          int           // Burst limit for rate limiting
	GVRCacheTTL    time.Duration // GroupVersionResource cache time-to-live
}

// Client provides high-level operations for interacting with Kubernetes clusters.
// It includes caching mechanisms for improved performance and supports both
// dynamic and typed client operations.
type Client struct {
	// Core Kubernetes clients
	clientset       kubernetes.Interface                           // Typed client for standard resources
	dynamicClient   dynamic.Interface                              // Dynamic client for any resource type
	discoveryClient discovery.DiscoveryInterface                   // API discovery client
	cacheDiscovery  *CacheableDiscovery                            // Enhanced cacheable discovery client
	authClient      authorizationv1client.AuthorizationV1Interface // Authorization client
	metricsClient   metricsv1beta1.Interface                       // Metrics client for resource usage
	restConfig      *rest.Config                                   // REST configuration
	kubeconfigPath  string                                         // Path to kubeconfig file

	// GVR cache for performance optimization
	gvrCache    map[string]schema.GroupVersionResource // Cache mapping kind to GVR
	gvrCacheMux sync.RWMutex                           // Mutex for thread-safe cache access
	cacheExpiry time.Time                              // Cache expiration timestamp
	cacheTTL    time.Duration                          // Cache time-to-live duration
}

// DefaultClientOptions returns default client options
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Timeout:     30 * time.Second,
		QPS:         100,
		Burst:       200,
		GVRCacheTTL: 15 * time.Minute,
	}
}

// NewClientWithOptions creates a new Kubernetes client with the specified options
func NewClientWithOptions(opts *ClientOptions) (*Client, error) {
	kubeconfigPath := resolveKubeconfigPath(opts.KubeconfigPath)

	var config *rest.Config
	var err error

	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			logrus.Debug("Not running in Kubernetes cluster, falling back to kubeconfig")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		} else {
			logrus.Info("Using in-cluster configuration (service account credentials)")
		}
	}
	if err != nil {
		return nil, err
	}

	// Apply rate limiting and timeout configuration
	if opts.QPS > 0 {
		config.QPS = opts.QPS
	}
	if opts.Burst > 0 {
		config.Burst = opts.Burst
	}
	if opts.Timeout > 0 {
		config.Timeout = opts.Timeout
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	metricsClient, err := metricsv1beta1.NewForConfig(config)
	if err != nil {
		// Metrics client is optional - don't fail if metrics server is not available
		metricsClient = nil
	}

	// Create enhanced cacheable discovery client
	cacheDiscovery := NewCacheableDiscovery(discoveryClient, 10*time.Minute)

	return &Client{
		clientset:       clientset,
		dynamicClient:   dynamicClient,
		discoveryClient: discoveryClient,
		cacheDiscovery:  cacheDiscovery,
		authClient:      clientset.AuthorizationV1(),
		metricsClient:   metricsClient,
		restConfig:      config,
		kubeconfigPath:  kubeconfigPath,
		gvrCache:        make(map[string]schema.GroupVersionResource, 100), // Pre-allocate size
		cacheTTL:        opts.GVRCacheTTL,
	}, nil
}

// GetKubeconfigPath returns the kubeconfig path used by this client
func (c *Client) GetKubeconfigPath() string {
	return c.kubeconfigPath
}

// GetRestConfig returns the REST config used by this client
func (c *Client) GetRestConfig() *rest.Config {
	return c.restConfig
}

// resolveKubeconfigPath resolves the kubeconfig path
// Priority:
// 1. Explicit path provided in config
// 2. KUBECONFIG environment variable or default kubeconfig location (if file exists)
// 3. Empty string to trigger InClusterConfig auto-detection (when in Pod)
func resolveKubeconfigPath(path string) string {
	if path != "" {
		return path
	}

	// Check KUBECONFIG environment variable or default kubeconfig location
	// Only return the path if the file actually exists
	if kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename(); kubeconfig != "" {
		if _, err := os.Stat(kubeconfig); err == nil {
			logrus.Debugf("Found kubeconfig at: %s", kubeconfig)
			return kubeconfig
		}
		logrus.Debugf("Kubeconfig file not found at: %s, will attempt in-cluster config", kubeconfig)
	}

	// Return empty string to trigger InClusterConfig detection
	logrus.Debug("No explicit kubeconfig found, attempting in-cluster configuration")
	return ""
}

// getCachedGVR retrieves a GroupVersionResource from cache
func (c *Client) getCachedGVR(kind string) (*schema.GroupVersionResource, error) {
	if kind == "" {
		return nil, fmt.Errorf("kind is empty")
	}

	c.gvrCacheMux.RLock()
	defer c.gvrCacheMux.RUnlock()

	// Check if cache is expired
	if time.Now().After(c.cacheExpiry) {
		return nil, fmt.Errorf("cache expired")
	}

	if gvr, exists := c.gvrCache[toLower(kind)]; exists {
		return &gvr, nil
	}

	return nil, fmt.Errorf("GVR not found in cache for kind: %s", kind)
}

// updateGVRCache adds or updates a GroupVersionResource in the cache
func (c *Client) updateGVRCache(kind string, gvr schema.GroupVersionResource) {
	c.gvrCacheMux.Lock()
	defer c.gvrCacheMux.Unlock()

	// Update cache expiry
	if time.Now().After(c.cacheExpiry) {
		c.cacheExpiry = time.Now().Add(c.cacheTTL)
	}

	// Store in cache
	c.gvrCache[toLower(kind)] = gvr
	logrus.WithFields(logrus.Fields{
		"kind":     kind,
		"group":    gvr.Group,
		"version":  gvr.Version,
		"resource": gvr.Resource,
	}).Debug("Added GVR to cache")
}

// Helper functions
func toLower(s string) string {
	return strings.ToLower(s)
}

func containsSlash(s string) bool {
	return strings.Contains(s, "/")
}

// PortForward creates a port forward to a pod
func (c *Client) PortForward(ctx context.Context, podName, namespace string, localPort, podPort int32, address string) error {
	// Build the URL for port forward
	req := c.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward")

	// Create SPDY transport
	transport, upgrader, err := spdy.RoundTripperFor(c.restConfig)
	if err != nil {
		return fmt.Errorf("failed to create SPDY round tripper: %w", err)
	}

	// Parse URL
	u, err := url.Parse(req.URL().String())
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Create port forwarder with optimized HTTP client
	optimizedClient := optimize.NewOptimizedHTTPClient()
	optimizedClient.Transport = transport
	dialer := spdy.NewDialer(upgrader, optimizedClient, "POST", u)

	// Setup port mapping
	ports := []string{fmt.Sprintf("%d:%d", localPort, podPort)}

	// Create ready and stop channels
	readyChannel := make(chan struct{})
	stopChannel := make(chan struct{}, 1)

	// Setup output streams (we'll use dummy writers since we're running in background)
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	// Create port forwarder
	pf, err := portforward.New(dialer, ports, stopChannel, readyChannel, out, errOut)
	if err != nil {
		return fmt.Errorf("failed to create port forwarder: %w", err)
	}

	// Start port forwarding in a goroutine
	go func() {
		defer close(stopChannel)
		if err := pf.ForwardPorts(); err != nil {
			logrus.WithError(err).Error("Port forwarding failed")
		}
	}()

	// Wait for ready signal or timeout
	select {
	case <-readyChannel:
		logrus.WithFields(logrus.Fields{
			"pod":       podName,
			"namespace": namespace,
			"localPort": localPort,
			"podPort":   podPort,
			"address":   address,
		}).Info("Port forwarding established")
		return nil
	case <-time.After(30 * time.Second):
		close(stopChannel)
		return fmt.Errorf("timeout waiting for port forward to be ready")
	case <-ctx.Done():
		close(stopChannel)
		return ctx.Err()
	}
}

// GetResourceUsage retrieves resource usage metrics for nodes or pods
func (c *Client) GetResourceUsage(ctx context.Context, resourceType, name, namespace string) (map[string]any, error) {
	if c.metricsClient == nil {
		return nil, fmt.Errorf("metrics client not available - ensure metrics server is installed and accessible")
	}

	switch strings.ToLower(resourceType) {
	case "node":
		if name != "" {
			// Get specific node metrics
			nodeMetrics, err := c.metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to get node metrics for %s: %w", name, err)
			}
			return map[string]any{
				"kind": "NodeMetrics",
				"metadata": map[string]any{
					"name": nodeMetrics.Name,
				},
				"timestamp": nodeMetrics.Timestamp.Format(time.RFC3339),
				"window":    nodeMetrics.Window.Duration.String(),
				"usage": map[string]any{
					"cpu":    nodeMetrics.Usage.Cpu().String(),
					"memory": nodeMetrics.Usage.Memory().String(),
				},
			}, nil
		} else {
			// Get all node metrics
			nodeMetricsList, err := c.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to list node metrics: %w", err)
			}

			var nodes []map[string]any
			for _, nodeMetrics := range nodeMetricsList.Items {
				nodes = append(nodes, map[string]any{
					"name":      nodeMetrics.Name,
					"timestamp": nodeMetrics.Timestamp.Format(time.RFC3339),
					"window":    nodeMetrics.Window.Duration.String(),
					"usage": map[string]any{
						"cpu":    nodeMetrics.Usage.Cpu().String(),
						"memory": nodeMetrics.Usage.Memory().String(),
					},
				})
			}
			return map[string]any{
				"kind":  "NodeMetricsList",
				"items": nodes,
			}, nil
		}

	case "pod":
		if namespace == "" {
			return nil, fmt.Errorf("namespace is required for pod metrics")
		}

		if name != "" {
			// Get specific pod metrics
			podMetrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to get pod metrics for %s/%s: %w", namespace, name, err)
			}

			var containers []map[string]any
			for _, container := range podMetrics.Containers {
				containers = append(containers, map[string]any{
					"name": container.Name,
					"usage": map[string]any{
						"cpu":    container.Usage.Cpu().String(),
						"memory": container.Usage.Memory().String(),
					},
				})
			}

			return map[string]any{
				"kind": "PodMetrics",
				"metadata": map[string]any{
					"name":      podMetrics.Name,
					"namespace": podMetrics.Namespace,
				},
				"timestamp":  podMetrics.Timestamp.Format(time.RFC3339),
				"window":     podMetrics.Window.Duration.String(),
				"containers": containers,
			}, nil
		} else {
			// Get all pod metrics in namespace
			podMetricsList, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to list pod metrics in namespace %s: %w", namespace, err)
			}

			var pods []map[string]any
			for _, podMetrics := range podMetricsList.Items {
				var containers []map[string]any
				for _, container := range podMetrics.Containers {
					containers = append(containers, map[string]any{
						"name": container.Name,
						"usage": map[string]any{
							"cpu":    container.Usage.Cpu().String(),
							"memory": container.Usage.Memory().String(),
						},
					})
				}

				pods = append(pods, map[string]any{
					"name":       podMetrics.Name,
					"namespace":  podMetrics.Namespace,
					"timestamp":  podMetrics.Timestamp.Format(time.RFC3339),
					"window":     podMetrics.Window.Duration.String(),
					"containers": containers,
				})
			}
			return map[string]any{
				"kind":  "PodMetricsList",
				"items": pods,
			}, nil
		}

	default:
		return nil, fmt.Errorf("unsupported resource type for metrics: %s (supported: node, pod)", resourceType)
	}
}

// isInClusterConfig checks if running inside a Kubernetes pod
// Checks for:
// 1. KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT environment variables
// 2. Service account token file at /var/run/secrets/kubernetes.io/serviceaccount/token
// 3. Service account CA certificate at /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
func isInClusterConfig() bool {
	const (
		kubernetesServiceHostEnv = "KUBERNETES_SERVICE_HOST"
		kubernetesServicePortEnv = "KUBERNETES_SERVICE_PORT"
		serviceAccountTokenPath  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		serviceAccountCAPath     = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)

	_, hostOk := os.LookupEnv(kubernetesServiceHostEnv)
	_, portOk := os.LookupEnv(kubernetesServicePortEnv)

	if !hostOk || !portOk {
		return false
	}

	if _, err := os.Stat(serviceAccountTokenPath); err != nil {
		return false
	}

	if _, err := os.Stat(serviceAccountCAPath); err != nil {
		return false
	}

	return true
}

// ExtractResourceSummaries extracts lightweight summaries from resource objects
func (c *Client) ExtractResourceSummaries(objects []map[string]interface{}, selectedLabelKeys []string) []map[string]interface{} {
	var summaries []map[string]interface{}

	for _, obj := range objects {
		unstruct := &unstructured.Unstructured{Object: obj}
		summary := extractResourceSummary(unstruct, selectedLabelKeys)
		if summary != nil {
			summaries = append(summaries, summary)
		}
	}

	return summaries
}

// extractResourceSummary extracts essential fields from a resource
func extractResourceSummary(obj *unstructured.Unstructured, selectedLabelKeys []string) map[string]interface{} {
	if obj == nil {
		return nil
	}

	summary := map[string]interface{}{
		"name": obj.GetName(),
		"kind": obj.GetKind(),
		"age":  calculateResourceAge(obj.GetCreationTimestamp()),
	}

	if ns := obj.GetNamespace(); ns != "" {
		summary["namespace"] = ns
	}

	if status, _, _ := unstructured.NestedString(obj.Object, "status", "phase"); status != "" {
		summary["status"] = status
	}

	if labels := obj.GetLabels(); len(labels) > 0 {
		if len(selectedLabelKeys) > 0 {
			selected := make(map[string]string)
			for _, key := range selectedLabelKeys {
				if val, exists := labels[key]; exists {
					selected[key] = val
				}
			}
			if len(selected) > 0 {
				summary["labels"] = selected
			}
		} else if len(labels) <= 10 {
			summary["labels"] = labels
		}
	}

	return summary
}

// calculateResourceAge calculates time elapsed since creation
func calculateResourceAge(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "unknown"
	}

	duration := time.Since(timestamp.Time)

	if duration.Hours() < 1 {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1m"
		}
		return fmt.Sprintf("%dm", minutes)
	} else if duration.Hours() < 24 {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1h"
		}
		return fmt.Sprintf("%dh", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1d"
	}
	return fmt.Sprintf("%dd", days)
}
