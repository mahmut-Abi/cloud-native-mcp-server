package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// HandleGetMirrorConfiguration returns a handler function for getting mirror configuration.
func HandleGetMirrorConfiguration(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_mirror_configuration").Debug("Handler invoked")

		// Get mirror configuration
		mirrorConfig := c.GetMirrorConfiguration()

		logrus.Debug("helm_get_mirror_configuration succeeded")

		// Serialize to JSON for better readability
		jsonData, err := marshalIndentJSON(mirrorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleClearCache clears the Helm cache
func HandleClearCache(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_clear_cache").Debug("Handler invoked")

		if err := c.ClearCache(); err != nil {
			return nil, fmt.Errorf("failed to clear cache: %w", err)
		}

		logrus.Debug("helm_clear_cache succeeded")
		return mcp.NewToolResultText("Helm cache cleared successfully"), nil
	}
}

// HandleGetCacheStats returns cache statistics
func HandleGetCacheStats(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_cache_stats").Debug("Handler invoked")

		stats, err := c.GetCacheStats()
		if err != nil {
			return nil, fmt.Errorf("failed to get cache stats: %w", err)
		}

		logrus.Debug("helm_cache_stats succeeded")
		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(stats)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleHelmHealthCheck handles Helm service health diagnostics
func HandleHelmHealthCheck(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_health_check").Debug("Handler invoked")

		checkClient := getOptionalBoolParam(request, "checkClient")
		checkKubernetes := getOptionalBoolParam(request, "checkKubernetes")
		checkRepositories := getOptionalBoolParam(request, "checkRepositories")
		checkCache := getOptionalBoolParam(request, "checkCache")

		// Default to all checks if none specified
		if !checkClient && !checkKubernetes && !checkRepositories && !checkCache {
			checkClient = true
			checkKubernetes = true
			checkRepositories = true
			checkCache = true
		}

		health := map[string]interface{}{
			"service":     "helm",
			"checks":      map[string]interface{}{},
			"initialized": c != nil,
		}

		if c != nil {
			// Try to list releases as a basic connectivity test
			if checkClient {
				_, listErr := c.ListReleases(false, "")
				clientStatus := "initialized"
				if listErr != nil {
					clientStatus = "error: " + listErr.Error()
				}
				health["checks"].(map[string]interface{})["client"] = map[string]interface{}{
					"status": clientStatus,
				}
			}

			if checkKubernetes {
				// Try to list releases across all namespaces as a K8s connectivity test
				_, k8sErr := c.ListReleases(true, "")
				k8sStatus := "connected"
				if k8sErr != nil {
					k8sStatus = "error: " + k8sErr.Error()
				}
				health["checks"].(map[string]interface{})["kubernetes"] = map[string]interface{}{
					"status": k8sStatus,
				}
			}

			if checkRepositories {
				repos, err := c.ListRepositories()
				repoStatus := "unknown"
				if err != nil {
					repoStatus = "error: " + err.Error()
				} else if repos != nil {
					repoStatus = fmt.Sprintf("%d repositories configured", len(repos))
				}
				health["checks"].(map[string]interface{})["repositories"] = map[string]interface{}{
					"status": repoStatus,
					"count":  len(repos),
				}
			}

			if checkCache {
				stats, err := c.GetCacheStats()
				cacheStatus := map[string]interface{}{}
				if err != nil {
					cacheStatus["status"] = "error: " + err.Error()
				} else if stats != nil {
					cacheStatus = map[string]interface{}{
						"status":     "ok",
						"cachePath":  stats.CachePath,
						"indexFiles": stats.IndexFiles,
					}
				}
				health["checks"].(map[string]interface{})["cache"] = cacheStatus
			}
		} else {
			health["checks"].(map[string]interface{})["client"] = map[string]interface{}{
				"status":  "not_initialized",
				"message": "Helm client is nil - check service initialization",
			}
		}

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(health)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}
