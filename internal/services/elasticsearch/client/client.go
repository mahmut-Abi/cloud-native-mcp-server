// Package client provides Elasticsearch HTTP client operations for the MCP server.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

// ClientOptions represents Elasticsearch client configuration
type ClientOptions struct {
	Addresses     []string
	Username      string
	Password      string
	BearerToken   string
	APIKey        string
	Timeout       time.Duration
	TLSSkipVerify bool
}

// Client represents an Elasticsearch HTTP client
type Client struct {
	httpClient *http.Client
	addresses  []string
	authType   string
	username   string
	password   string
	token      string
}

// NewClient creates a new Elasticsearch client
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		opts = &ClientOptions{
			Addresses: []string{"http://localhost:9200"},
			Timeout:   30 * time.Second,
		}
	}

	if len(opts.Addresses) == 0 {
		opts.Addresses = []string{"http://localhost:9200"}
	}

	for i, addr := range opts.Addresses {
		if !strings.HasPrefix(addr, "http") {
			opts.Addresses[i] = "http://" + addr
		}
	}

	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(opts.Timeout)

	c := &Client{
		httpClient: httpClient,
		addresses:  opts.Addresses,
		username:   opts.Username,
		password:   opts.Password,
	}

	if opts.APIKey != "" {
		c.authType = "apikey"
		c.token = opts.APIKey
	} else if opts.BearerToken != "" {
		c.authType = "bearer"
		c.token = opts.BearerToken
	} else if opts.Username != "" && opts.Password != "" {
		c.authType = "basic"
	}

	return c, nil
}

// Health checks cluster health
func (c *Client) Health(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.get(ctx, "/_cluster/health")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Indices returns list of indices
func (c *Client) Indices(ctx context.Context) ([]string, error) {
	resp, err := c.get(ctx, "/_cat/indices?format=json")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var indices []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&indices); err != nil {
		return nil, err
	}

	var names []string
	for _, idx := range indices {
		if name, ok := idx["index"]; ok {
			names = append(names, fmt.Sprintf("%v", name))
		}
	}
	return names, nil
}

// IndexStats returns statistics
func (c *Client) IndexStats(ctx context.Context, indexName string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/%s/_stats", indexName)
	resp, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Nodes returns node information
func (c *Client) Nodes(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.get(ctx, "/_nodes")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Info returns cluster info
func (c *Client) Info(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.get(ctx, "/")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) addAuthHeaders(req *http.Request) {
	switch c.authType {
	case "basic":
		req.SetBasicAuth(c.username, c.password)
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+c.token)
	case "apikey":
		req.Header.Set("Authorization", "ApiKey "+c.token)
	}
}

func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	return c.doRequest(ctx, "GET", path, nil)
}

// PaginationInfo represents pagination metadata for Elasticsearch responses
type PaginationInfo struct {
	ContinueToken   string `json:"continueToken"`
	RemainingCount  int64  `json:"remainingCount"`
	CurrentPageSize int64  `json:"currentPageSize"`
	HasMore         bool   `json:"hasMore"`
}

// IndicesFull returns detailed indices information with all metadata
func (c *Client) IndicesFull(ctx context.Context, indexPattern string) ([]map[string]interface{}, error) {
	path := "/_cat/indices?format=json"
	if indexPattern != "" && indexPattern != "*" {
		path += "&index=" + url.QueryEscape(indexPattern)
	}

	resp, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var indices []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&indices); err != nil {
		return nil, err
	}
	return indices, nil
}

// IndicesPaginated returns indices with pagination support
func (c *Client) IndicesPaginated(ctx context.Context, continueToken string, limit int, indexPattern string, includeHealth bool) ([]map[string]interface{}, *PaginationInfo, error) {
	// Get all indices first (Elasticsearch doesn't have server-side pagination for _cat API)
	allIndices, err := c.IndicesFull(ctx, indexPattern)
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination
	totalCount := int64(len(allIndices))
	offset := 0
	if continueToken != "" {
		if _, err := fmt.Sscanf(continueToken, "offset:%d", &offset); err != nil {
			offset = 0
		}
	}

	start := offset
	end := offset + limit
	if start >= len(allIndices) {
		allIndices = []map[string]interface{}{}
		end = start
	} else if end > len(allIndices) {
		allIndices = allIndices[start:]
		end = len(allIndices)
	} else {
		allIndices = allIndices[start:end]
	}

	// Create summary format to reduce size
	summaries := make([]map[string]interface{}, len(allIndices))
	for i, idx := range allIndices {
		summary := map[string]interface{}{
			"name":       idx["index"],
			"status":     idx["status"],
			"uuid":       idx["uuid"],
			"pri":        idx["pri"],
			"rep":        idx["rep"],
			"docs.count": idx["docs.count"],
			"store.size": idx["store.size"],
		}

		if includeHealth {
			summary["health"] = idx["health"]
		}

		summaries[i] = summary
	}

	// Build pagination info
	pagination := &PaginationInfo{
		CurrentPageSize: int64(len(summaries)),
		RemainingCount:  totalCount - int64(end),
		HasMore:         end < len(allIndices)+offset,
	}

	if pagination.HasMore {
		pagination.ContinueToken = fmt.Sprintf("offset:%d", end)
	}

	return summaries, pagination, nil
}

// NodesSummary returns optimized node information
func (c *Client) NodesSummary(ctx context.Context, role string, includeMetrics bool, limit int) ([]map[string]interface{}, error) {
	nodes, err := c.Nodes(ctx)
	if err != nil {
		return nil, err
	}

	var nodeSummaries []map[string]interface{}

	if nodesMap, ok := nodes["nodes"].(map[string]interface{}); ok {
		count := 0
		for nodeID, nodeInfo := range nodesMap {
			if count >= limit {
				break
			}

			if nodeData, ok := nodeInfo.(map[string]interface{}); ok {
				// Apply role filter if specified
				if role != "" {
					if roles, ok := nodeData["roles"].([]interface{}); ok {
						hasRole := false
						for _, r := range roles {
							if rStr, ok := r.(string); ok && rStr == role {
								hasRole = true
								break
							}
						}
						if !hasRole {
							continue
						}
					}
				}

				summary := map[string]interface{}{
					"node_id": nodeID,
					"name":    nodeData["name"],
					"host":    nodeData["host"],
					"ip":      nodeData["ip"],
					"version": nodeData["version"],
				}

				if roles, ok := nodeData["roles"].([]interface{}); ok {
					summary["roles"] = roles
				}

				if includeMetrics {
					if os, ok := nodeData["os"].(map[string]interface{}); ok {
						if cpu, ok := os["available_processors"].(float64); ok {
							summary["cpu_cores"] = int64(cpu)
						}
					}
					if jvm, ok := nodeData["jvm"].(map[string]interface{}); ok {
						if mem, ok := jvm["mem"].(map[string]interface{}); ok {
							if heapUsed, ok := mem["heap_used_percent"].(float64); ok {
								summary["heap_used_percent"] = heapUsed
							}
						}
					}
				}

				nodeSummaries = append(nodeSummaries, summary)
				count++
			}
		}
	}

	return nodeSummaries, nil
}

// ClusterHealthSummary returns lightweight health information
func (c *Client) ClusterHealthSummary(ctx context.Context, level string, includeIndices bool) (map[string]interface{}, error) {
	health, err := c.Health(ctx)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"status":          health["status"],
		"number_of_nodes": health["number_of_nodes"],
		"active_shards":   health["active_shards"],
		"cluster_name":    health["cluster_name"],
		"timestamp":       time.Now().Format(time.RFC3339),
	}

	if level == "detailed" || level == "indices" {
		summary["timed_out"] = health["timed_out"]
		summary["delayed_unassigned_shards"] = health["delayed_unassigned_shards"]
		summary["number_of_data_nodes"] = health["number_of_data_nodes"]
		summary["active_primary_shards"] = health["active_primary_shards"]
		summary["initializing_shards"] = health["initializing_shards"]
		summary["unassigned_shards"] = health["unassigned_shards"]
	}

	if includeIndices && level == "indices" {
		// Get per-index health (this adds more data)
		indices, err := c.IndicesFull(ctx, "")
		if err == nil {
			var indexHealth []map[string]interface{}
			for _, idx := range indices {
				if idx["health"] != nil {
					indexHealth = append(indexHealth, map[string]interface{}{
						"name":   idx["index"],
						"health": idx["health"],
						"status": idx["status"],
					})
				}
			}
			if len(indexHealth) > 0 {
				summary["indices_health"] = indexHealth
			}
		}
	}

	return summary, nil
}

// GetIndexDetailAdvanced returns comprehensive index information
func (c *Client) GetIndexDetailAdvanced(ctx context.Context, indexName string, includeMappings, includeSettings, includeStats, includeSegments bool, outputFormat string) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"index":   indexName,
		"queryAt": time.Now().Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"includeMappings": includeMappings,
			"includeSettings": includeSettings,
			"includeStats":    includeStats,
			"includeSegments": includeSegments,
			"outputFormat":    outputFormat,
		},
	}

	// Get mappings
	if includeMappings {
		if mappings, err := c.getIndexMappings(ctx, indexName); err == nil {
			result["mappings"] = mappings
		}
	}

	// Get settings
	if includeSettings {
		if settings, err := c.getIndexSettings(ctx, indexName); err == nil {
			result["settings"] = settings
		}
	}

	// Get stats
	if includeStats {
		if stats, err := c.IndexStats(ctx, indexName); err == nil {
			result["stats"] = stats
		}
	}

	// Get segments
	if includeSegments {
		if segments, err := c.getIndexSegments(ctx, indexName); err == nil {
			result["segments"] = segments
		}
	}

	// Apply output format optimization
	switch outputFormat {
	case "compact":
		// Remove verbose fields for compact output
		if result["mappings"] != nil {
			delete(result, "mappings")
		}
		if result["settings"] != nil {
			delete(result, "settings")
		}
		if result["segments"] != nil {
			delete(result, "segments")
		}
	case "verbose":
		// Add raw data for complete analysis (would need additional API calls)
		result["note"] = "Verbose output would include additional raw API responses"
	}

	return result, nil
}

// GetClusterDetailAdvanced returns comprehensive cluster information
func (c *Client) GetClusterDetailAdvanced(ctx context.Context, includeNodes, includeIndices, includeSettings, includeStats, includeShards bool, outputFormat string) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"cluster": "elasticsearch",
		"queryAt": time.Now().Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"includeNodes":    includeNodes,
			"includeIndices":  includeIndices,
			"includeSettings": includeSettings,
			"includeStats":    includeStats,
			"includeShards":   includeShards,
			"outputFormat":    outputFormat,
		},
	}

	// Get basic cluster info
	if info, err := c.Info(ctx); err == nil {
		result["info"] = info
	}

	// Get health
	if health, err := c.Health(ctx); err == nil {
		result["health"] = health
	}

	// Get detailed node information
	if includeNodes {
		if nodes, err := c.Nodes(ctx); err == nil {
			result["nodes"] = nodes
		}
	}

	// Get indices overview
	if includeIndices {
		if indices, err := c.IndicesFull(ctx, ""); err == nil {
			// Create summary to reduce size
			var indicesSummary []map[string]interface{}
			for _, idx := range indices {
				indicesSummary = append(indicesSummary, map[string]interface{}{
					"name":       idx["index"],
					"status":     idx["status"],
					"health":     idx["health"],
					"docs.count": idx["docs.count"],
					"store.size": idx["store.size"],
					"pri":        idx["pri"],
					"rep":        idx["rep"],
				})
			}
			result["indices_summary"] = indicesSummary
		}
	}

	// Get cluster statistics
	if includeStats {
		if stats, err := c.getClusterStats(ctx); err == nil {
			result["cluster_stats"] = stats
		}
	}

	// Get shard allocation
	if includeShards {
		if shards, err := c.getShardAllocation(ctx); err == nil {
			result["shard_allocation"] = shards
		}
	}

	// Get cluster settings
	if includeSettings {
		if settings, err := c.getClusterSettings(ctx); err == nil {
			result["cluster_settings"] = settings
		}
	}

	// Apply output format optimization
	if outputFormat == "compact" {
		// Keep only essential information
		if result["health"] != nil {
			if health, ok := result["health"].(map[string]interface{}); ok {
				result["health"] = map[string]interface{}{
					"status":          health["status"],
					"number_of_nodes": health["number_of_nodes"],
					"active_shards":   health["active_shards"],
				}
			}
		}
		if result["nodes"] != nil {
			delete(result, "nodes") // Remove detailed nodes in compact mode
		}
		if result["indices_summary"] != nil {
			delete(result, "indices_summary") // Remove indices in compact mode
		}
	}

	return result, nil
}

// SearchIndices provides advanced index search capabilities
func (c *Client) SearchIndices(ctx context.Context, query, healthStatus, indexStatus string, minDocCount, maxDocCount int, sortBy, sortOrder string, limit int, continueToken string) ([]map[string]interface{}, *PaginationInfo, error) {
	// Get all indices first
	allIndices, err := c.IndicesFull(ctx, "")
	if err != nil {
		return nil, nil, err
	}

	var filteredIndices []map[string]interface{}

	// Apply filters
	for _, idx := range allIndices {
		// Query filter
		if query != "" && query != "*" {
			if name, ok := idx["index"].(string); ok {
				if !strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
					continue
				}
			}
		}

		// Health status filter
		if healthStatus != "" {
			if health, ok := idx["health"].(string); ok && health != healthStatus {
				continue
			}
		}

		// Index status filter
		if indexStatus != "" {
			if status, ok := idx["status"].(string); ok && status != indexStatus {
				continue
			}
		}

		// Document count filters
		if minDocCount > 0 || maxDocCount > 0 {
			if docsCount, ok := idx["docs.count"].(string); ok {
				if count, err := strconv.ParseInt(docsCount, 10, 64); err == nil {
					if minDocCount > 0 && count < int64(minDocCount) {
						continue
					}
					if maxDocCount > 0 && count > int64(maxDocCount) {
						continue
					}
				}
			}
		}

		filteredIndices = append(filteredIndices, idx)
	}

	// Apply sorting
	sort.Slice(filteredIndices, func(i, j int) bool {
		var iVal, jVal interface{}

		switch sortBy {
		case "name":
			iVal = filteredIndices[i]["index"]
			jVal = filteredIndices[j]["index"]
		case "docs":
			iVal = filteredIndices[i]["docs.count"]
			jVal = filteredIndices[j]["docs.count"]
		case "size":
			iVal = filteredIndices[i]["store.size"]
			jVal = filteredIndices[j]["store.size"]
		case "health":
			iVal = filteredIndices[i]["health"]
			jVal = filteredIndices[j]["health"]
		default:
			iVal = filteredIndices[i]["index"]
			jVal = filteredIndices[j]["index"]
		}

		// Simple string comparison
		iStr := fmt.Sprintf("%v", iVal)
		jStr := fmt.Sprintf("%v", jVal)

		if sortOrder == "desc" {
			return iStr > jStr
		}
		return iStr < jStr
	})

	// Apply pagination
	totalCount := int64(len(filteredIndices))
	offset := 0
	if continueToken != "" {
		if _, err := fmt.Sscanf(continueToken, "offset:%d", &offset); err != nil {
			offset = 0
		}
	}

	start := offset
	end := offset + limit
	if start >= len(filteredIndices) {
		filteredIndices = []map[string]interface{}{}
		end = start
	} else if end > len(filteredIndices) {
		filteredIndices = filteredIndices[start:]
		end = len(filteredIndices)
	} else {
		filteredIndices = filteredIndices[start:end]
	}

	// Build pagination info
	pagination := &PaginationInfo{
		CurrentPageSize: int64(len(filteredIndices)),
		RemainingCount:  totalCount - int64(end),
		HasMore:         end < len(filteredIndices)+offset,
	}

	if pagination.HasMore {
		pagination.ContinueToken = fmt.Sprintf("offset:%d", end)
	}

	return filteredIndices, pagination, nil
}

// Helper methods for advanced queries
func (c *Client) getIndexMappings(ctx context.Context, indexName string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/%s/_mapping", indexName)
	resp, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) getIndexSettings(ctx context.Context, indexName string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/%s/_settings", indexName)
	resp, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) getIndexSegments(ctx context.Context, indexName string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/%s/_segments", indexName)
	resp, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) getClusterStats(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.get(ctx, "/_cluster/stats")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) getShardAllocation(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.get(ctx, "/_cat/shards?format=json")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) getClusterSettings(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.get(ctx, "/_cluster/settings")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	if len(c.addresses) == 0 {
		return nil, fmt.Errorf("no elasticsearch addresses")
	}

	baseURL := c.addresses[0]
	urlStr, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.addAuthHeaders(req)

	logrus.WithFields(logrus.Fields{
		"method": method,
		"path":   path,
	}).Debug("Elasticsearch request")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}
