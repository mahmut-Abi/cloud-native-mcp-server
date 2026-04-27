package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/jsonpath"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/constants"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kubernetes/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/sanitize"
)

// Type alias for PaginationInfo from client package
type PaginationInfo = client.PaginationInfo

var (
	ErrMissingRequiredParam = errors.New("missing required parameter")
	ErrInvalidJSONPath      = errors.New("invalid jsonpath expression")
	ErrJSONPathExecution    = errors.New("jsonpath execution error")
	ErrInvalidManifest      = errors.New("invalid manifest format")
	ErrCommandExecutionFail = errors.New("command execution failed")
)

func applyJSONPath(input any, expr string) (any, error) {
	jp := jsonpath.New("mcp-jsonpath")
	jp.AllowMissingKeys(true)

	expr = normalizeJSONPathExpression(expr)
	if err := jp.Parse(expr); err != nil {
		return nil, fmt.Errorf("%w '%s': %v", ErrInvalidJSONPath, expr, err)
	}

	var buf bytes.Buffer
	if err := jp.Execute(&buf, input); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJSONPathExecution, err)
	}

	resultStr := buf.String()
	lines := strings.Split(strings.ReplaceAll(resultStr, "\r\n", "\n"), "\n")
	var final []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			final = append(final, line)
		}
	}
	if len(final) == 0 {
		final = []string{}
	}
	return final, nil
}

func getRequestArguments(request mcp.CallToolRequest) map[string]any {
	args := request.GetArguments()
	if args == nil {
		return nil
	}

	nested, ok := args["params"].(map[string]any)
	if !ok || len(nested) == 0 {
		return args
	}

	merged := make(map[string]any, len(nested)+len(args))
	for k, v := range nested {
		merged[k] = v
	}
	for k, v := range args {
		if k == "params" {
			continue
		}
		merged[k] = v
	}
	return merged
}

// Helper function to validate required string parameter
func requireStringParam(request mcp.CallToolRequest, param string) (string, error) {
	value, ok := getRequestArguments(request)[param].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
	}
	return sanitize.SanitizeFilterValue(value), nil
}

func requireArgument(request mcp.CallToolRequest, param string) (any, error) {
	value, ok := getRequestArguments(request)[param]
	if !ok || value == nil {
		return nil, fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
	}
	return value, nil
}

// requireRawStringParam returns the original string value without sanitization.
// Use this for JSON payloads where quotes, braces, and whitespace are significant.
func requireRawStringParam(request mcp.CallToolRequest, param string) (string, error) {
	value, ok := getRequestArguments(request)[param].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
	}
	return value, nil
}

// Helper function to get optional string parameter
func getOptionalStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := getRequestArguments(request)[param].(string)
	return sanitize.SanitizeFilterValue(value)
}

// getOptionalRawStringParam returns the original optional string value without sanitization.
func getOptionalRawStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := getRequestArguments(request)[param].(string)
	return value
}

func requireJSONObjectParam(request mcp.CallToolRequest, param string) (map[string]any, error) {
	value, err := requireArgument(request, param)
	if err != nil {
		return nil, err
	}

	switch typed := value.(type) {
	case map[string]interface{}:
		return typed, nil
	case string:
		var result map[string]any
		if strings.TrimSpace(typed) == "" {
			return nil, fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
		}
		if err := json.Unmarshal([]byte(typed), &result); err != nil {
			return nil, fmt.Errorf("failed to parse %s JSON object: %w", param, err)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("%s must be a JSON object", param)
	}
}

func getOptionalJSONObjectParam(request mcp.CallToolRequest, param string) (map[string]any, bool, error) {
	value, ok := getRequestArguments(request)[param]
	if !ok || value == nil {
		return nil, false, nil
	}

	switch typed := value.(type) {
	case map[string]interface{}:
		return typed, true, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil, false, nil
		}
		var result map[string]any
		if err := json.Unmarshal([]byte(typed), &result); err != nil {
			return nil, true, fmt.Errorf("failed to parse %s JSON object: %w", param, err)
		}
		return result, true, nil
	default:
		return nil, true, fmt.Errorf("%s must be a JSON object", param)
	}
}

func requireRawJSONParam(request mcp.CallToolRequest, param string) ([]byte, error) {
	value, err := requireArgument(request, param)
	if err != nil {
		return nil, err
	}

	switch typed := value.(type) {
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil, fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
		}
		return []byte(typed), nil
	case map[string]interface{}, []interface{}:
		data, err := optimize.GlobalJSONPool.MarshalToBytes(typed)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize %s JSON payload: %w", param, err)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("%s must be a JSON object, JSON array, or JSON string", param)
	}
}

func getOptionalStringArrayParam(request mcp.CallToolRequest, param string) ([]string, error) {
	value, ok := getRequestArguments(request)[param]
	if !ok || value == nil {
		return nil, nil
	}

	switch typed := value.(type) {
	case []string:
		return typed, nil
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("%s must contain only strings", param)
			}
			result = append(result, str)
		}
		return result, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil, nil
		}

		var result []string
		if strings.HasPrefix(strings.TrimSpace(typed), "[") {
			if err := json.Unmarshal([]byte(typed), &result); err != nil {
				return nil, fmt.Errorf("failed to parse %s JSON array: %w", param, err)
			}
			return result, nil
		}

		parts := strings.Split(typed, ",")
		result = make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("%s must be an array of strings", param)
	}
}

func normalizeJSONPathExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return expr
	}
	if strings.HasPrefix(expr, "{") {
		return expr
	}
	if strings.HasPrefix(expr, ".") {
		return "{" + expr + "}"
	}
	return "{." + expr + "}"
}

func getOptionalSearchKinds(request mcp.CallToolRequest) ([]string, error) {
	kind := getOptionalStringParam(request, "kind")
	if kind != "" {
		return []string{kind}, nil
	}

	resourceTypes, err := getOptionalStringArrayParam(request, "resourceTypes")
	if err != nil {
		return nil, err
	}
	if len(resourceTypes) > 0 {
		return resourceTypes, nil
	}

	return []string{
		"Pod",
		"Deployment",
		"StatefulSet",
		"DaemonSet",
		"Service",
		"Job",
		"CronJob",
		"ConfigMap",
		"Secret",
		"Ingress",
	}, nil
}

// Helper function to marshal JSON response using pooled encoder
func marshalJSONResponse(data any) (*mcp.CallToolResult, error) {
	jsonResponse, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}
	return mcp.NewToolResultText(string(jsonResponse)), nil
}

// Helper function to create optimized JSON response for LLM
func marshalOptimizedResponse(data any, toolName string) (*mcp.CallToolResult, error) {
	result := FormatResponseForLLM(data, toolName)
	return result, nil
}

// Helper function to create error response
func createErrorResponse(message string) *mcp.CallToolResult {
	return mcp.NewToolResultError(message)
}

func getBoolParam(request mcp.CallToolRequest, param string, defaultValue bool) bool {
	if value, ok := getRequestArguments(request)[param]; ok {
		switch typed := value.(type) {
		case bool:
			return typed
		case string:
			switch strings.ToLower(strings.TrimSpace(typed)) {
			case "true", "1", "yes", "y", "on":
				return true
			case "false", "0", "no", "n", "off":
				return false
			}
		}
	}
	return defaultValue
}

// getNestedString extracts nested string from map safely
func getNestedString(obj map[string]any, path string) string {
	if obj == nil || path == "" {
		return ""
	}
	parts := strings.Split(path, ".")
	current := obj
	for i, part := range parts {
		if i == len(parts)-1 {
			if val, ok := current[part].(string); ok {
				return val
			}
			return ""
		}
		if next, ok := current[part].(map[string]any); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

// HandleDescribeResource handles resource description requests (similar to kubectl describe).
func HandleDescribeResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{"tool": "describe_resource", "kind": kind, "name": name, "ns": namespace, "debug": debug}).Debug("Handler invoked")

		result, err := client.GetResource(ctx, kind, name, namespace)
		if err != nil {
			return nil, err
		}
		logrus.Debug("describe_resource succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleGetResourceUsage handles resource usage information requests (CPU/Memory).
func HandleGetResourceUsage(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resourceType, err := requireStringParam(request, "resourceType")
		if err != nil {
			return nil, err
		}
		name := getOptionalStringParam(request, "name")
		namespace := getOptionalStringParam(request, "namespace")
		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{"tool": "get_resource_usage", "resourceType": resourceType, "name": name, "ns": namespace, "debug": debug}).Debug("Handler invoked")

		// Use the new GetResourceUsage method to get actual metrics
		result, err := client.GetResourceUsage(ctx, resourceType, name, namespace)
		if err != nil {
			return nil, err
		}
		logrus.Debug("get_resource_usage succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleGetRecentEvents handles recent critical events retrieval with optimized output.
func HandleGetRecentEvents(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := getOptionalStringParam(request, "namespace")
		fieldSelector := getOptionalStringParam(request, "fieldSelector")
		debug := getOptionalStringParam(request, "debug")

		// More conservative default limit for recent events
		limit := int64(20)
		if v, ok := request.GetArguments()["limit"]; ok {
			if f, ok := v.(float64); ok {
				limit = int64(f)
				if limit <= 0 || limit > 100 {
					limit = 20
				}
			}
		}

		logrus.WithFields(logrus.Fields{"tool": "get_recent_events", "ns": namespace, "fieldSelector": fieldSelector, "limit": limit, "debug": debug}).Debug("Handler invoked")

		// Create field selector that focuses on important events only
		selector := fieldSelector
		if selector == "" {
			selector = "type!=Normal" // By default, exclude normal events
		} else {
			selector = fmt.Sprintf("%s,type!=Normal", selector)
		}

		// Use paginated listing to prevent context overflow
		resources, err := client.ListResourcesWithPagination(ctx, "Event", namespace, "", selector, "", limit)
		if err != nil {
			return nil, err
		}

		// Extract only essential fields from events
		var recentEvents []map[string]interface{}
		for _, event := range resources {
			recentEvents = append(recentEvents, map[string]interface{}{
				"type":      getNestedString(event, "type"),
				"reason":    getNestedString(event, "reason"),
				"message":   getNestedString(event, "message"),
				"timestamp": getNestedString(event, "lastTimestamp"),
				"object":    fmt.Sprintf("%s/%s", getNestedString(event, "involvedObject.kind"), getNestedString(event, "involvedObject.name")),
				"namespace": getNestedString(event, "involvedObject.namespace"),
			})
		}

		// Get pagination info
		paginationInfo, err := client.GetPaginationInfo(ctx, "Event", namespace, "", selector, "", limit)
		if err != nil {
			logrus.WithError(err).Warn("Failed to get pagination info for recent events")
			paginationInfo = &PaginationInfo{ContinueToken: "", RemainingCount: 0, CurrentPageSize: 0, HasMore: false}
		}

		response := map[string]interface{}{
			"events": recentEvents,
			"count":  len(recentEvents),
			"pagination": map[string]interface{}{
				"continueToken":   paginationInfo.ContinueToken,
				"remainingCount":  paginationInfo.RemainingCount,
				"currentPageSize": paginationInfo.CurrentPageSize,
				"hasMore":         paginationInfo.HasMore,
			},
		}

		logrus.WithFields(logrus.Fields{"count": len(recentEvents), "hasMore": paginationInfo.HasMore}).Debug("get_recent_events succeeded")
		return marshalJSONResponse(response)
	}
}

// HandleGetEvents handles events retrieval requests for troubleshooting.
func HandleGetEvents(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := getOptionalStringParam(request, "namespace")
		fieldSelector := getOptionalStringParam(request, "fieldSelector")
		debug := getOptionalStringParam(request, "debug")

		limit := int64(constants.DefaultLimit)
		if v, ok := request.GetArguments()["limit"]; ok {
			if f, ok := v.(float64); ok {
				limit = int64(f)
				if limit <= 0 || limit > constants.MaxLimit {
					if limit > constants.MaxLimit {
						logrus.WithField("requested", limit).WithField("max", constants.MaxLimit).Warn("Event limit too high, resetting to safe maximum")
						limit = constants.MaxLimit
					} else {
						limit = constants.DefaultLimit
					}
				}
				if limit > constants.WarningLimit {
					logrus.WithField("limit", limit).Warn("Large event limit may cause context overflow, consider using get_recent_events for critical events only")
				}
			}
		}

		logrus.WithFields(logrus.Fields{"tool": "get_events", "ns": namespace, "fieldSelector": fieldSelector, "limit": limit, "debug": debug}).Debug("Handler invoked")

		// Use paginated listing to prevent context overflow
		resources, err := client.ListResourcesWithPagination(ctx, "Event", namespace, "", fieldSelector, "", limit)
		if err != nil {
			return nil, err
		}

		// Get pagination info
		paginationInfo, err := client.GetPaginationInfo(ctx, "Event", namespace, "", fieldSelector, "", limit)
		if err != nil {
			logrus.WithError(err).Warn("Failed to get pagination info for events")
			paginationInfo = &PaginationInfo{ContinueToken: "", RemainingCount: 0, CurrentPageSize: 0, HasMore: false}
		}

		// Create response with pagination metadata
		response := map[string]interface{}{
			"events": resources,
			"count":  len(resources),
			"pagination": map[string]interface{}{
				"continueToken":   paginationInfo.ContinueToken,
				"remainingCount":  paginationInfo.RemainingCount,
				"currentPageSize": paginationInfo.CurrentPageSize,
				"hasMore":         paginationInfo.HasMore,
			},
		}

		logrus.WithFields(logrus.Fields{"count": len(resources), "hasMore": paginationInfo.HasMore}).Debug("get_events succeeded")
		return marshalOptimizedResponse(response, "get_events")
	}
}

// HandleGetResourceDetails handles detailed resource information requests.
func HandleGetResourceDetails(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{"tool": "get_resource_details", "kind": kind, "name": name, "ns": namespace, "debug": debug}).Debug("Handler invoked")

		result, err := client.GetResource(ctx, kind, name, namespace)
		if err != nil {
			return nil, err
		}
		logrus.Debug("get_resource_details succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleContainerLogs handles log requests for a container and returns the log content.
func HandleContainerLogs(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := getOptionalStringParam(request, "namespace")
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		container := getOptionalStringParam(request, "container")
		logrus.WithFields(logrus.Fields{"tool": "get_pod_logs", "pod": name, "ns": namespace, "container": container}).Debug("Handler invoked")

		tailLines := int64(constants.DefaultTailLines)
		if v, ok := request.GetArguments()["tailLines"]; ok {
			if f, ok := v.(float64); ok {
				tailLines = int64(f)
				if tailLines < 0 || tailLines > 200 { // Maximum limit to prevent excessive output
					if tailLines > 200 {
						logrus.WithField("requested", tailLines).Warn("Log tail lines too high, resetting to safe maximum")
						tailLines = 200
					} else {
						tailLines = constants.DefaultTailLines
					}
				}
			}
		}

		result, err := client.GetContainerLog(ctx, name, namespace, container, tailLines)
		if err != nil {
			return nil, err
		}
		// Smart log processing with size monitoring
		logSize := len(result)

		// Pre-process large logs to prevent context overflow
		processedLogs := result
		truncationInfo := map[string]interface{}{}

		if logSize > 10000 { // 10KB threshold for logs
			logrus.WithFields(logrus.Fields{
				"pod":       name,
				"size":      logSize,
				"tailLines": tailLines,
			}).Warn("Large log response detected, applying smart truncation")

			// Split log into lines and keep recent portion
			lines := strings.Split(result, "\n")
			limitedLines := lines
			if len(lines) > 200 { // Maximum 200 lines
				limitedLines = lines[len(lines)-200:] // Keep last 200 lines
				truncationInfo["truncated"] = true
				truncationInfo["originalLines"] = len(lines)
				truncationInfo["retainedLines"] = len(limitedLines)
			}
			processedLogs = strings.Join(limitedLines, "\n")

			// Check character count limit
			if len(processedLogs) > 50000 { // 50KB character limit
				processedLogs = processedLogs[len(processedLogs)-50000:] // Keep last 50KB
				truncationInfo["charTruncated"] = true
				truncationInfo["originalChars"] = logSize
				truncationInfo["retainedChars"] = len(processedLogs)
				if _, exists := truncationInfo["retainedLines"]; !exists {
					truncationInfo["retainedLines"] = "estimate"
				}
			}

			if len(truncationInfo) > 0 {
				truncationInfo["reason"] = "Context overflow prevention"
			}
		}

		logData := map[string]interface{}{
			"logs": processedLogs,
			"metadata": map[string]interface{}{
				"pod":           name,
				"namespace":     namespace,
				"container":     container,
				"tailLines":     tailLines,
				"originalSize":  logSize,
				"processedSize": len(processedLogs),
			},
		}

		// Add truncation information if applied
		if len(truncationInfo) > 0 {
			logData["metadata"].(map[string]interface{})["truncation"] = truncationInfo
		}

		logrus.WithFields(logrus.Fields{
			"pod":           name,
			"originalSize":  logSize,
			"processedSize": len(processedLogs),
			"tailLines":     tailLines,
		}).Debug("get_pod_logs with smart processing succeeded")

		return marshalOptimizedResponse(logData, "get_pod_logs")
	}
}

// HandlePortForward handles port forwarding requests to a pod.
func HandlePortForward(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		podName, err := requireStringParam(request, "podName")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		localPort := int32(0)
		if v, ok := request.GetArguments()["localPort"]; ok {
			if f, ok := v.(float64); ok {
				localPort = int32(f)
			} else {
				return nil, fmt.Errorf("localPort must be a number")
			}
		} else {
			return nil, fmt.Errorf("missing required parameter: localPort")
		}

		podPort := int32(0)
		if v, ok := request.GetArguments()["podPort"]; ok {
			if f, ok := v.(float64); ok {
				podPort = int32(f)
			} else {
				return nil, fmt.Errorf("podPort must be a number")
			}
		} else {
			return nil, fmt.Errorf("missing required parameter: podPort")
		}

		address := getOptionalStringParam(request, "address")
		if address == "" {
			address = "localhost"
		}
		debug := getOptionalStringParam(request, "debug")

		logrus.WithFields(logrus.Fields{"tool": "port_forward", "pod": podName, "ns": namespace, "localPort": localPort, "podPort": podPort, "address": address, "debug": debug}).Debug("Handler invoked")

		err = client.PortForward(ctx, podName, namespace, localPort, podPort, address)
		if err != nil {
			return nil, err
		}
		return marshalJSONResponse(map[string]any{
			"status":    "ok",
			"message":   "port forwarding established",
			"address":   address,
			"localPort": localPort,
			"namespace": namespace,
			"podName":   podName,
			"podPort":   podPort,
		})
	}
}

// HandleCreateResource handles resource creation requests.
func HandleCreateResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		apiVersion, err := requireStringParam(request, "apiVersion")
		if err != nil {
			return nil, err
		}
		metadata, err := requireJSONObjectParam(request, "metadata")
		if err != nil {
			return nil, err
		}
		metadataJSON, err := optimize.GlobalJSONPool.MarshalToBytes(metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize metadata JSON: %w", err)
		}

		specJSON := ""
		spec, hasSpec, err := getOptionalJSONObjectParam(request, "spec")
		if err != nil {
			return nil, err
		}
		if hasSpec {
			specBytes, err := optimize.GlobalJSONPool.MarshalToBytes(spec)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize spec JSON: %w", err)
			}
			specJSON = string(specBytes)
		}
		logrus.WithFields(logrus.Fields{"tool": "create_resource", "kind": kind, "apiVersion": apiVersion}).Debug("Handler invoked")

		result, err := client.CreateResource(ctx, kind, apiVersion, string(metadataJSON), specJSON)
		if err != nil {
			return nil, err
		}
		logrus.Debug("create_resource succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleUpdateResource handles update requests for Kubernetes resources.
func HandleUpdateResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		name := getOptionalStringParam(request, "name")
		manifest, err := requireRawStringParam(request, "manifest")
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidManifest, err)
		}
		logrus.WithFields(logrus.Fields{"tool": "update_resource", "kind": kind, "name": name, "ns": namespace}).Debug("Handler invoked")

		result, err := client.UpdateResource(ctx, kind, name, namespace, manifest)
		if err != nil {
			return nil, err
		}
		logrus.Debug("update_resource succeeded")
		return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
	}
}

// HandlePatchResource handles patch requests for Kubernetes resources.
func HandlePatchResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		patchBytes, err := requireRawJSONParam(request, "patch")
		if err != nil {
			return nil, err
		}
		patchType := getOptionalStringParam(request, "patchType")
		if patchType == "" {
			patchType = "merge"
		}
		logrus.WithFields(logrus.Fields{"tool": "patch_resource", "kind": kind, "name": name, "ns": namespace, "patchType": patchType}).Debug("Handler invoked")

		result, err := client.PatchResource(ctx, kind, name, namespace, patchBytes, patchType)
		if err != nil {
			return nil, err
		}
		logrus.Debug("patch_resource succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleContainerExec handles command execution requests in containers.
func HandleContainerExec(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "podName")
		if err != nil {
			return nil, err
		}
		container, err := requireStringParam(request, "containerName")
		if err != nil {
			return nil, err
		}
		commandEncoded, err := requireRawStringParam(request, "command")
		if err != nil {
			return nil, err
		}
		commandArgs := strings.Fields(commandEncoded)
		if len(commandArgs) == 0 {
			return nil, fmt.Errorf("parsed command is empty")
		}
		logrus.WithFields(logrus.Fields{"tool": "pod_exec", "pod": name, "ns": namespace, "container": container}).Debug("Handler invoked")

		result, err := client.ExecCommand(ctx, name, namespace, container, commandArgs)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCommandExecutionFail, err)
		}
		logrus.Debug("pod_exec succeeded")
		return mcp.NewToolResultText(result), nil
	}
}

// HandleGetResourceSummary handles resource summary requests with minimal output.
func HandleGetResourceSummary(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		includeLabels := getOptionalStringParam(request, "includeLabels")
		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{"tool": "get_resource_summary", "kind": kind, "name": name, "ns": namespace, "debug": debug}).Debug("Handler invoked")

		// Get the full resource first
		resource, err := client.GetResource(ctx, kind, name, namespace)
		if err != nil {
			return nil, err
		}

		// Parse includeLabels
		var labelKeys []string
		if includeLabels != "" {
			labelKeys = strings.Split(includeLabels, ",")
			for i := range labelKeys {
				labelKeys[i] = strings.TrimSpace(labelKeys[i])
			}
		}

		// Extract summary using the existing summary functionality
		summaries := client.ExtractResourceSummaries([]map[string]interface{}{resource}, labelKeys)
		if len(summaries) == 0 {
			return createErrorResponse("failed to extract resource summary"), nil
		}

		response := map[string]interface{}{
			"summary": summaries[0],
			"kind":    kind,
			"name":    name,
		}
		if namespace != "" {
			response["namespace"] = namespace
		}

		logrus.Debug("get_resource_summary succeeded")
		return marshalJSONResponse(response)
	}
}

// HandleGetResource handles resource retrieval requests.
func HandleGetResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		jsonpath := getOptionalRawStringParam(request, "jsonpath")
		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{"tool": "get_resource", "kind": kind, "name": name, "ns": namespace, "jsonpath": jsonpath, "debug": debug}).Debug("Handler invoked")

		resource, err := client.GetResource(ctx, kind, name, namespace)
		if err != nil {
			return nil, err
		}

		var result interface{} = resource

		// Apply JSONPath filter if provided
		if jsonpath != "" {
			filtered, err := applyJSONPath(resource, jsonpath)
			if err != nil {
				logrus.WithError(err).Warn("JSONPath filtering failed")
				return createErrorResponse(err.Error()), nil
			}
			result = filtered
		}

		logrus.Debug("get_resource succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleListResources handles listing resources requests with pagination support.
func HandleListResources(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		labelSelector := getOptionalStringParam(request, "labelSelector")
		fieldSelector := getOptionalStringParam(request, "fieldSelector")
		jsonpath := getOptionalRawStringParam(request, "jsonpath")
		continueToken := getOptionalStringParam(request, "continueToken")
		debug := getOptionalStringParam(request, "debug")
		jsonpaths, err := getOptionalStringArrayParam(request, "jsonpaths")
		if err != nil {
			return createErrorResponse(err.Error()), nil
		}

		// Parse limit parameter with conservative default to prevent context overflow
		limit := int64(constants.DefaultLimit)
		if v, ok := request.GetArguments()["limit"]; ok {
			if f, ok := v.(float64); ok {
				limit = int64(f)
				if limit <= 0 || limit > constants.MaxLimit {
					if limit > constants.MaxLimit {
						logrus.WithField("requested", limit).WithField("max", constants.MaxLimit).Warn("Limit too high, resetting to safe maximum")
						limit = constants.MaxLimit
					} else {
						limit = constants.DefaultLimit // Reset to default if out of bounds
					}
				}
				if limit > constants.WarningLimit {
					logrus.WithField("limit", limit).Warn("Large limit may cause context overflow, consider using summary tools or pagination")
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"tool":      "list_resources",
			"kind":      kind,
			"ns":        namespace,
			"labels":    labelSelector,
			"fields":    fieldSelector,
			"jsonpath":  jsonpath,
			"jsonpaths": jsonpaths,
			"continue":  continueToken,
			"limit":     limit,
			"debug":     debug,
		}).Debug("Handler invoked")

		resources, err := client.ListResourcesWithPagination(ctx, kind, namespace, labelSelector, fieldSelector, continueToken, limit)
		if err != nil {
			return createErrorResponse(err.Error()), nil
		}

		// Get pagination info
		paginationInfo, err := client.GetPaginationInfo(ctx, kind, namespace, labelSelector, fieldSelector, continueToken, limit)
		if err != nil {
			logrus.WithError(err).Warn("Failed to get pagination info")
			paginationInfo = &PaginationInfo{ContinueToken: "", RemainingCount: 0, CurrentPageSize: 0, HasMore: false}
		}

		// Wrap resources into {"items": [...] } to support JSONPath like {.items[*].metadata.name}
		wrapped := map[string]any{"items": resources}
		var result any = wrapped

		// Apply JSONPath filter if provided (single jsonpath for backward compatibility)
		if jsonpath != "" {
			filtered, err := applyJSONPath(wrapped, jsonpath)
			if err != nil {
				logrus.WithError(err).Warn("JSONPath filtering failed, returning error response")
				return createErrorResponse(err.Error()), nil
			} else {
				result = filtered
			}
		} else if len(jsonpaths) > 0 {
			// Handle multiple JSONPath expressions
			expressions := jsonpaths
			var tableResult []map[string]any

			// Process each resource - convert resources to []interface{} for processing
			resourceList := make([]any, len(resources))
			for i, resource := range resources {
				resourceList[i] = resource
			}

			for _, resource := range resourceList {
				row := make(map[string]any)
				for i, expr := range expressions {
					expr = strings.TrimSpace(expr)
					if expr == "" {
						continue
					}

					fieldValue, err := applyJSONPath(resource, expr)
					if err != nil {
						logrus.WithError(err).WithField("expression", expr).Debug("JSONPath expression failed")
						row[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("error: %s", err.Error())
					} else {
						row[fmt.Sprintf("field_%d", i)] = fieldValue
					}
				}
				tableResult = append(tableResult, row)
			}
			result = map[string]any{
				"expressions": expressions,
				"data":        tableResult,
			}
		}

		// Add pagination metadata to the response
		response := map[string]any{
			"data": result,
			"pagination": map[string]any{
				"continueToken":   paginationInfo.ContinueToken,
				"remainingCount":  paginationInfo.RemainingCount,
				"currentPageSize": paginationInfo.CurrentPageSize,
				"hasMore":         paginationInfo.HasMore,
			},
			"count": len(resources),
		}

		logrus.WithFields(logrus.Fields{
			"count":     len(resources),
			"hasMore":   paginationInfo.HasMore,
			"remaining": paginationInfo.RemainingCount,
		}).Debug("list_resources succeeded")
		return marshalOptimizedResponse(response, "list_resources")
	}
}

// HandleListResourcesSummary handles listing resources with summary output for LLM efficiency
func HandleListResourcesSummary(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		labelSelector := getOptionalStringParam(request, "labelSelector")
		includeLabels := getOptionalStringParam(request, "includeLabels")
		limitStr := getOptionalStringParam(request, "limit")
		continueToken := getOptionalStringParam(request, "continueToken")

		limit := int64(constants.DefaultLimit)
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				if l > constants.MaxLimit {
					logrus.WithField("requested", l).WithField("max", constants.MaxLimit).Warn("Summary limit too high, resetting to safe maximum")
					limit = constants.MaxLimit
				} else {
					limit = int64(l)
				}
				if l > constants.WarningLimit {
					logrus.WithField("limit", l).Warn("Large summary limit may cause context overflow")
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"tool":     "list_resources_summary",
			"kind":     kind,
			"ns":       namespace,
			"labels":   labelSelector,
			"limit":    limit,
			"continue": continueToken,
		}).Debug("Handler invoked")

		// Use paginated listing to avoid loading too much data
		resources, err := client.ListResourcesWithPagination(ctx, kind, namespace, labelSelector, "", continueToken, limit)
		if err != nil {
			return createErrorResponse(err.Error()), nil
		}

		// Get pagination info
		paginationInfo, err := client.GetPaginationInfo(ctx, kind, namespace, labelSelector, "", continueToken, limit)
		if err != nil {
			logrus.WithError(err).Warn("Failed to get pagination info for summary")
			paginationInfo = &PaginationInfo{ContinueToken: "", RemainingCount: 0, CurrentPageSize: 0, HasMore: false}
		}

		// Parse includeLabels
		var labelKeys []string
		if includeLabels != "" {
			labelKeys = strings.Split(includeLabels, ",")
			for i := range labelKeys {
				labelKeys[i] = strings.TrimSpace(labelKeys[i])
			}
		}

		// Extract summaries (already limited by pagination)
		summaries := client.ExtractResourceSummaries(resources, labelKeys)

		response := map[string]interface{}{
			"items": summaries,
			"count": len(summaries),
			"pagination": map[string]interface{}{
				"continueToken":   paginationInfo.ContinueToken,
				"remainingCount":  paginationInfo.RemainingCount,
				"currentPageSize": paginationInfo.CurrentPageSize,
				"hasMore":         paginationInfo.HasMore,
			},
		}

		// Apply caching optimization for summary responses
		finalResponse := response
		if logrus.GetLevel() >= logrus.InfoLevel {
			logrus.WithFields(logrus.Fields{
				"count":     len(summaries),
				"hasMore":   paginationInfo.HasMore,
				"itemsSize": len(summaries),
				"responseSize": func() int {
					if data, err := optimize.GlobalJSONPool.MarshalToBytes(response); err == nil {
						return len(data)
					}
					return 0
				}(),
			}).Info("list_resources_summary response prepared")
		} else {
			logrus.WithFields(logrus.Fields{
				"count":   len(summaries),
				"hasMore": paginationInfo.HasMore,
			}).Debug("list_resources_summary succeeded")
		}

		return marshalOptimizedResponse(finalResponse, "list_resources_summary")
	}
}

// HandleDeleteResource handles resource deletion requests.
func HandleDeleteResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		logrus.WithFields(logrus.Fields{"tool": "delete_resource", "kind": kind, "name": name, "ns": namespace}).Debug("Handler invoked")

		err = client.DeleteResource(ctx, kind, name, namespace)
		if err != nil {
			return nil, err
		}
		logrus.Debug("delete_resource succeeded")
		return marshalJSONResponse(map[string]any{
			"status":    "ok",
			"message":   "resource deleted successfully",
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		})
	}
}

// HandleCheckPermissions handles permission checking requests.
func HandleCheckPermissions(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		verb, err := requireStringParam(request, "verb")
		if err != nil {
			return nil, err
		}
		resourceName := getOptionalStringParam(request, "resourceName")
		namespace := getOptionalStringParam(request, "namespace")
		resourceGroup := getOptionalStringParam(request, "resourceGroup")
		resourceResource := getOptionalStringParam(request, "resourceResource")
		subresource := getOptionalStringParam(request, "subresource")
		logrus.WithFields(logrus.Fields{"tool": "check_permissions", "verb": verb, "group": resourceGroup, "resource": resourceResource, "subresource": subresource, "name": resourceName, "ns": namespace}).Debug("Handler invoked")

		result, err := client.CheckPermissions(ctx, verb, resourceName, resourceGroup, resourceResource, subresource, namespace)
		if err != nil {
			return nil, err
		}
		message := "denied"
		if result {
			message = "allowed"
		}
		logrus.WithField("allowed", result).Debug("check_permissions succeeded")
		return marshalJSONResponse(map[string]any{
			"allowed":       result,
			"message":       message,
			"verb":          verb,
			"resourceName":  resourceName,
			"resourceGroup": resourceGroup,
			"resource":      resourceResource,
			"subresource":   subresource,
			"namespace":     namespace,
		})
	}
}

// HandleTest handles test requests with confirmation.
func HandleTest(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "test_tool").Debug("Handler invoked")
		confirmed := getBoolParam(request, "confirmed", false)
		if !confirmed {
			logrus.Info("test handler: confirmed is false")
			return createErrorResponse("this operation requires confirmation; set confirmed=true to continue"), nil
		}

		logrus.Info("test handler: confirmed is true")
		return marshalJSONResponse(map[string]any{
			"status":    "ok",
			"confirmed": true,
			"message":   "operation confirmed",
		})
	}
}

// HandleScaleResource scales a namespaced resource to target replicas.
func HandleScaleResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}
		replicasVal, ok := request.GetArguments()["replicas"]
		if !ok {
			return nil, fmt.Errorf("missing required parameter: replicas")
		}
		var replicas int32
		switch v := replicasVal.(type) {
		case float64:
			replicas = int32(v)
		case int:
			replicas = int32(v)
		default:
			return nil, fmt.Errorf("replicas must be a number")
		}
		logrus.WithFields(logrus.Fields{"tool": "scale_resource", "kind": kind, "name": name, "ns": namespace, "replicas": replicas}).Debug("Handler invoked")

		if err := client.ScaleResourceByKind(ctx, kind, name, namespace, replicas); err != nil {
			return nil, err
		}
		logrus.Debug("scale_resource succeeded")
		return marshalJSONResponse(map[string]any{
			"status":    "ok",
			"message":   "resource scaled successfully",
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
			"replicas":  replicas,
		})
	}
}

// HandleGetRolloutStatus returns rollout status for a workload resource.
func HandleGetRolloutStatus(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		timeoutSeconds := 120
		if v, ok := request.GetArguments()["timeoutSeconds"]; ok {
			switch typed := v.(type) {
			case float64:
				timeoutSeconds = int(typed)
			case int:
				timeoutSeconds = typed
			}
		}

		result, err := client.GetRolloutStatus(ctx, kind, name, namespace, timeoutSeconds)
		if err != nil {
			return nil, err
		}
		result["timeoutSeconds"] = timeoutSeconds
		return marshalJSONResponse(result)
	}
}

// HandleCordonNode marks a node unschedulable.
func HandleCordonNode(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeName, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		if err := client.CordonNode(ctx, nodeName); err != nil {
			return nil, err
		}

		return marshalJSONResponse(map[string]any{
			"status":  "ok",
			"message": "node cordoned successfully",
			"name":    nodeName,
		})
	}
}

// HandleUncordonNode marks a node schedulable.
func HandleUncordonNode(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeName, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		if err := client.UncordonNode(ctx, nodeName); err != nil {
			return nil, err
		}

		return marshalJSONResponse(map[string]any{
			"status":  "ok",
			"message": "node uncordoned successfully",
			"name":    nodeName,
		})
	}
}

// HandleDrainNode cordons and drains a node.
func HandleDrainNode(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeName, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		deleteEmptyDir := getBoolParam(request, "deleteEmptyDir", false)
		ignoreDaemonsets := getBoolParam(request, "ignoreDaemonsets", true)
		gracePeriodSeconds := int32(30)
		timeoutSeconds := int32(120)

		if v, ok := request.GetArguments()["gracePeriodSeconds"]; ok {
			switch typed := v.(type) {
			case float64:
				gracePeriodSeconds = int32(typed)
			case int:
				gracePeriodSeconds = int32(typed)
			}
		}
		if v, ok := request.GetArguments()["timeoutSeconds"]; ok {
			switch typed := v.(type) {
			case float64:
				timeoutSeconds = int32(typed)
			case int:
				timeoutSeconds = int32(typed)
			}
		}

		if err := client.DrainNode(ctx, nodeName, deleteEmptyDir, ignoreDaemonsets, gracePeriodSeconds, timeoutSeconds); err != nil {
			return nil, err
		}

		return marshalJSONResponse(map[string]any{
			"status":             "ok",
			"message":            "node drained successfully",
			"name":               nodeName,
			"deleteEmptyDir":     deleteEmptyDir,
			"ignoreDaemonsets":   ignoreDaemonsets,
			"gracePeriodSeconds": gracePeriodSeconds,
			"timeoutSeconds":     timeoutSeconds,
		})
	}
}

// HandleWaitForResource waits until a resource reaches a condition.
func HandleWaitForResource(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		namespace := getOptionalStringParam(request, "namespace")
		condition := getOptionalStringParam(request, "condition")
		if condition == "" {
			condition = "ready"
		}

		timeoutSeconds := 120
		pollIntervalSeconds := 5
		if v, ok := request.GetArguments()["timeoutSeconds"]; ok {
			switch typed := v.(type) {
			case float64:
				timeoutSeconds = int(typed)
			case int:
				timeoutSeconds = typed
			}
		}
		if v, ok := request.GetArguments()["pollIntervalSeconds"]; ok {
			switch typed := v.(type) {
			case float64:
				pollIntervalSeconds = int(typed)
			case int:
				pollIntervalSeconds = typed
			}
		}

		result, err := client.WaitForResource(ctx, kind, name, namespace, condition, timeoutSeconds, pollIntervalSeconds)
		if err != nil {
			return nil, err
		}
		return marshalJSONResponse(result)
	}
}

// HandleRestartWorkload triggers a rollout restart for supported workload kinds.
func HandleRestartWorkload(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace, err := request.RequireString("namespace")
		if err != nil {
			return nil, err
		}

		waitForReady := getBoolParam(request, "waitForReady", false)
		timeoutSeconds := 120
		if v, ok := request.GetArguments()["timeoutSeconds"]; ok {
			switch typed := v.(type) {
			case float64:
				timeoutSeconds = int(typed)
			case int:
				timeoutSeconds = typed
			}
		}

		resource, restartedAt, err := client.RestartWorkload(ctx, kind, name, namespace)
		if err != nil {
			return nil, err
		}

		response := map[string]any{
			"status":      "ok",
			"message":     "workload restart triggered",
			"kind":        kind,
			"name":        name,
			"namespace":   namespace,
			"restartedAt": restartedAt,
			"resource":    resource,
		}

		if waitForReady {
			waitCondition := "available"
			if strings.EqualFold(kind, "StatefulSet") {
				waitCondition = "ready"
			}
			waitResult, err := client.WaitForResource(ctx, kind, name, namespace, waitCondition, timeoutSeconds, 5)
			if err != nil {
				return nil, err
			}
			response["wait"] = waitResult
		}

		return marshalJSONResponse(response)
	}
}

// HandleGetAPIVersions handles API versions retrieval requests.
func HandleGetAPIVersions(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{"tool": "get_api_versions", "debug": debug}).Debug("Handler invoked")

		result, err := client.GetAPIVersions(ctx)
		if err != nil {
			return nil, err
		}
		logrus.Debug("get_api_versions succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleGetAPIResources handles API resources retrieval requests.
func HandleGetAPIResources(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		apiGroup := getOptionalStringParam(request, "apiGroup")
		debug := getOptionalStringParam(request, "debug")

		var namespaced *bool
		if v, ok := request.GetArguments()["namespaced"]; ok {
			if b, ok := v.(bool); ok {
				namespaced = &b
			}
		}

		logrus.WithFields(logrus.Fields{"tool": "get_api_resources", "apiGroup": apiGroup, "namespaced": namespaced, "debug": debug}).Debug("Handler invoked")

		result, err := client.GetAPIResources(ctx, apiGroup, namespaced)
		if err != nil {
			return nil, err
		}
		logrus.Debug("get_api_resources succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleGetResourcesDetail handles detailed resource retrieval for multiple resources
func HandleGetResourcesDetail(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		includeEvents := getBoolParam(request, "includeEvents", false)
		includeStatus := getBoolParam(request, "includeStatus", true)
		debug := getOptionalStringParam(request, "debug")

		names, err := getOptionalStringArrayParam(request, "names")
		if err != nil {
			return createErrorResponse(err.Error()), nil
		}

		if len(names) == 0 {
			return createErrorResponse("names parameter is required and must be a non-empty array"), nil
		}

		logrus.WithFields(logrus.Fields{
			"tool":      "get_resources_detail",
			"kind":      kind,
			"names":     len(names),
			"namespace": namespace,
			"events":    includeEvents,
			"status":    includeStatus,
			"debug":     debug,
		}).Debug("Handler invoked")

		resources, err := client.GetResourcesDetail(ctx, kind, names, namespace, includeEvents, includeStatus)
		if err != nil {
			return createErrorResponse(err.Error()), nil
		}

		response := map[string]interface{}{
			"resources": resources,
			"count":     len(resources),
			"kind":      kind,
		}

		// Add metadata about the request for context
		response["metadata"] = map[string]interface{}{
			"requestedCount": len(names),
			"retrievedCount": len(resources),
			"includeEvents":  includeEvents,
			"includeStatus":  includeStatus,
			"namespace":      namespace,
		}

		logrus.WithFields(logrus.Fields{
			"requested": len(names),
			"retrieved": len(resources),
		}).Debug("get_resources_detail succeeded")

		return marshalJSONResponse(response)
	}
}

// HandleGetEventsDetail handles detailed events retrieval
func HandleGetEventsDetail(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := getOptionalStringParam(request, "namespace")
		fieldSelector := getOptionalStringParam(request, "fieldSelector")
		includeNormalEvents := getBoolParam(request, "includeNormalEvents", false)
		debug := getOptionalStringParam(request, "debug")

		// More conservative default for detailed events
		limit := int64(50)
		if v, ok := request.GetArguments()["limit"]; ok {
			if f, ok := v.(float64); ok {
				limit = int64(f)
				if limit <= 0 || limit > 200 {
					if limit > 200 {
						logrus.WithField("requested", limit).Warn("Event detail limit too high, resetting to maximum")
						limit = 200
					} else {
						limit = 50
					}
				}
				if limit > 100 {
					logrus.WithField("limit", limit).Warn("Large event detail limit may cause context overflow")
				}
			}
		}

		continueToken := getOptionalStringParam(request, "continueToken")

		logrus.WithFields(logrus.Fields{
			"tool":          "get_events_detail",
			"ns":            namespace,
			"fieldSelector": fieldSelector,
			"includeNormal": includeNormalEvents,
			"limit":         limit,
			"continue":      continueToken,
			"debug":         debug,
		}).Debug("Handler invoked")

		// Build field selector to exclude normal events unless requested
		selector := fieldSelector
		if !includeNormalEvents {
			if selector == "" {
				selector = "type!=Normal"
			} else {
				selector = fmt.Sprintf("%s,type!=Normal", selector)
			}
		}

		// Use paginated listing
		resources, err := client.ListResourcesWithPagination(ctx, "Event", namespace, "", selector, continueToken, limit)
		if err != nil {
			return nil, err
		}

		// Get pagination info
		paginationInfo, err := client.GetPaginationInfo(ctx, "Event", namespace, "", selector, continueToken, limit)
		if err != nil {
			logrus.WithError(err).Warn("Failed to get pagination info for detailed events")
			paginationInfo = &PaginationInfo{ContinueToken: "", RemainingCount: 0, CurrentPageSize: 0, HasMore: false}
		}

		response := map[string]interface{}{
			"events": resources,
			"count":  len(resources),
			"metadata": map[string]interface{}{
				"includeNormalEvents": includeNormalEvents,
				"fieldSelector":       selector,
			},
			"pagination": map[string]interface{}{
				"continueToken":   paginationInfo.ContinueToken,
				"remainingCount":  paginationInfo.RemainingCount,
				"currentPageSize": paginationInfo.CurrentPageSize,
				"hasMore":         paginationInfo.HasMore,
			},
		}

		logrus.WithFields(logrus.Fields{
			"count":   len(resources),
			"hasMore": paginationInfo.HasMore,
		}).Debug("get_events_detail succeeded")

		return marshalOptimizedResponse(response, "get_events_detail")
	}
}

// HandleListResourcesFull handles full resource listing without optimization
func HandleListResourcesFull(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		labelSelector := getOptionalStringParam(request, "labelSelector")
		fieldSelector := getOptionalStringParam(request, "fieldSelector")
		includeStatus := getBoolParam(request, "includeStatus", true)
		debug := getOptionalStringParam(request, "debug")
		continueToken := getOptionalStringParam(request, "continueToken")

		// Very conservative default for full resources
		limit := int64(10)
		if v, ok := request.GetArguments()["limit"]; ok {
			if f, ok := v.(float64); ok {
				limit = int64(f)
				if limit <= 0 || limit > 50 {
					if limit > 50 {
						logrus.WithField("requested", limit).Warn("Full resource limit too high, resetting to safe maximum")
						limit = 50
					} else {
						limit = 10
					}
				}
				if limit > 20 {
					logrus.WithField("limit", limit).Warn("Large full resource limit may cause significant context overflow")
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"tool":          "list_resources_full",
			"kind":          kind,
			"ns":            namespace,
			"labels":        labelSelector,
			"fields":        fieldSelector,
			"includeStatus": includeStatus,
			"limit":         limit,
			"continue":      continueToken,
			"debug":         debug,
		}).Debug("Handler invoked")

		resources, err := client.ListResourcesWithPagination(ctx, kind, namespace, labelSelector, fieldSelector, continueToken, limit)
		if err != nil {
			return createErrorResponse(err.Error()), nil
		}

		// Optionally remove status to reduce size if requested
		if !includeStatus {
			for _, resource := range resources {
				delete(resource, "status")
			}
		}

		// Get pagination info
		paginationInfo, err := client.GetPaginationInfo(ctx, kind, namespace, labelSelector, fieldSelector, continueToken, limit)
		if err != nil {
			logrus.WithError(err).Warn("Failed to get pagination info for full resources")
			paginationInfo = &PaginationInfo{ContinueToken: "", RemainingCount: 0, CurrentPageSize: 0, HasMore: false}
		}

		response := map[string]interface{}{
			"resources": resources,
			"count":     len(resources),
			"metadata": map[string]interface{}{
				"includeStatus": includeStatus,
				"fullDetails":   true,
			},
			"pagination": map[string]interface{}{
				"continueToken":   paginationInfo.ContinueToken,
				"remainingCount":  paginationInfo.RemainingCount,
				"currentPageSize": paginationInfo.CurrentPageSize,
				"hasMore":         paginationInfo.HasMore,
			},
		}

		logrus.WithFields(logrus.Fields{
			"count":   len(resources),
			"hasMore": paginationInfo.HasMore,
		}).Debug("list_resources_full succeeded")

		return marshalOptimizedResponse(response, "list_resources_full")
	}
}

// HandleGetResourceDetailAdvanced handles advanced resource detail retrieval with enhanced formatting
func HandleGetResourceDetailAdvanced(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kind, err := requireStringParam(request, "kind")
		if err != nil {
			return nil, err
		}
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace := getOptionalStringParam(request, "namespace")
		includeEvents := getBoolParam(request, "includeEvents", false)
		includeRelationships := getBoolParam(request, "includeRelationships", false)
		includeDiagnostics := getBoolParam(request, "includeDiagnostics", false)
		includeConfiguration := getBoolParam(request, "includeConfiguration", true)
		outputFormat := getOptionalStringParam(request, "outputFormat")
		debug := getOptionalStringParam(request, "debug")

		// Default to structured format if not specified
		if outputFormat == "" {
			outputFormat = "structured"
		}

		logrus.WithFields(logrus.Fields{
			"tool":                 "get_resource_detail_advanced",
			"kind":                 kind,
			"name":                 name,
			"namespace":            namespace,
			"includeEvents":        includeEvents,
			"includeRelationships": includeRelationships,
			"includeDiagnostics":   includeDiagnostics,
			"includeConfiguration": includeConfiguration,
			"outputFormat":         outputFormat,
			"debug":                debug,
		}).Debug("Handler invoked")

		// Get the base resource
		resource, err := client.GetResource(ctx, kind, name, namespace)
		if err != nil {
			return nil, err
		}

		// Build advanced detail response
		response := map[string]interface{}{
			"resource": map[string]interface{}{
				"kind":        kind,
				"name":        name,
				"namespace":   namespace,
				"retrievedAt": getCurrentTimestamp(),
			},
			"metadata": map[string]interface{}{
				"includeEvents":        includeEvents,
				"includeRelationships": includeRelationships,
				"includeDiagnostics":   includeDiagnostics,
				"includeConfiguration": includeConfiguration,
				"outputFormat":         outputFormat,
			},
		}

		// Add configuration if requested
		if includeConfiguration {
			response["configuration"] = resource
		} else {
			// Only include metadata and basic info
			if metadata, ok := resource["metadata"].(map[string]interface{}); ok {
				response["basicInfo"] = metadata
			}
			if spec, ok := resource["spec"].(map[string]interface{}); ok {
				response["spec"] = spec
			}
		}

		// Add events if requested
		if includeEvents {
			events, err := client.ListResourcesWithPagination(ctx, "Event", namespace,
				fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=%s", name, kind),
				"", "", 20)
			if err == nil {
				response["events"] = events
				response["eventCount"] = len(events)
			} else {
				logrus.WithError(err).Warn("Failed to retrieve events for advanced detail")
				response["events"] = []interface{}{}
				response["eventCount"] = 0
			}
		}

		// Add relationships if requested
		if includeRelationships {
			relationships := map[string]interface{}{}

			// Get owner references
			if metadata, ok := resource["metadata"].(map[string]interface{}); ok {
				if ownerRefs, exists := metadata["ownerReferences"]; exists {
					relationships["owners"] = ownerRefs
				}
			}

			// Try to find dependent resources (simple implementation)
			if labels, ok := resource["metadata"].(map[string]interface{})["labels"]; ok {
				if labelMap, ok := labels.(map[string]interface{}); ok {
					// Convert labels to selector format
					var selectors []string
					for k, v := range labelMap {
						selectors = append(selectors, fmt.Sprintf("%s=%s", k, v))
					}
					if len(selectors) > 0 {
						labelSelector := strings.Join(selectors, ",")
						// Look for pods with same labels (common case)
						if kind == "Deployment" || kind == "StatefulSet" || kind == "DaemonSet" {
							pods, err := client.ListResourcesWithPagination(ctx, "Pod", namespace, labelSelector, "", "", 10)
							if err == nil {
								relationships["dependents"] = map[string]interface{}{
									"pods":     pods,
									"podCount": len(pods),
								}
							}
						}
					}
				}
			}

			response["relationships"] = relationships
		}

		// Add diagnostics if requested
		if includeDiagnostics {
			diagnostics := map[string]interface{}{}

			// Check status conditions
			if status, ok := resource["status"].(map[string]interface{}); ok {
				if conditions, exists := status["conditions"]; exists {
					diagnostics["conditions"] = conditions
					diagnostics["statusHealth"] = analyzeHealthConditions(conditions)
				}
				if phase, exists := status["phase"]; exists {
					diagnostics["phase"] = phase
				}
			}

			// Add resource version for change detection
			if metadata, ok := resource["metadata"].(map[string]interface{}); ok {
				if resourceVersion, exists := metadata["resourceVersion"]; exists {
					diagnostics["resourceVersion"] = resourceVersion
				}
			}

			diagnostics["timestamp"] = getCurrentTimestamp()
			response["diagnostics"] = diagnostics
		}

		// Apply output formatting
		switch outputFormat {
		case "compact":
			// Remove verbose fields for compact output
			if response["configuration"] != nil {
				delete(response, "configuration")
			}
			delete(response, "metadata")
		case "verbose":
			// Add raw object for complete analysis
			response["rawObject"] = resource
		case "structured":
			// Default structured format - keep as is
		}

		logrus.WithField("kind", kind).WithField("name", name).Debug("get_resource_detail_advanced succeeded")

		// Use optimized response for large data
		data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, err
		}

		if len(data) > 50000 { // 50KB threshold
			logrus.Warn("Advanced detail response is large, using optimized formatting")
			return marshalOptimizedResponse(response, "get_resource_detail_advanced")
		}

		return marshalJSONResponse(response)
	}
}

// getCurrentTimestamp returns current timestamp in ISO format
func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// analyzeHealthConditions analyzes resource conditions for health status
func analyzeHealthConditions(conditions interface{}) string {
	if conditionSlice, ok := conditions.([]interface{}); ok {
		for _, cond := range conditionSlice {
			if conditionMap, ok := cond.(map[string]interface{}); ok {
				if typ, exists := conditionMap["type"]; exists && typ == "Ready" {
					if status, exists := conditionMap["status"]; exists && status == "True" {
						return "Healthy"
					} else if status == "False" {
						return "Unhealthy"
					}
				}
			}
		}
	}
	return "Unknown"
}

// ============ Troubleshooting Handlers ============

// HandleGetUnhealthyResources handles finding unhealthy resources
func HandleGetUnhealthyResources(k8sClient *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := ""
		if v, ok := request.GetArguments()["namespace"].(string); ok {
			namespace = v
		}

		resourceTypes := []string{}
		if v, ok := request.GetArguments()["resourceTypes"].([]interface{}); ok {
			for _, rt := range v {
				if s, ok := rt.(string); ok {
					resourceTypes = append(resourceTypes, s)
				}
			}
		}

		logrus.WithField("namespace", namespace).Debug("Executing get_unhealthy_resources handler")

		unhealthy, err := k8sClient.GetUnhealthyResources(ctx, namespace, resourceTypes)
		if err != nil {
			return nil, fmt.Errorf("failed to get unhealthy resources: %w", err)
		}

		response := map[string]interface{}{
			"unhealthyResources": unhealthy,
			"count":              len(unhealthy),
		}

		data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// HandleGetNodeConditions handles retrieving node conditions
func HandleGetNodeConditions(k8sClient *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeName, err := requireStringParam(request, "nodeName")
		if err != nil {
			return nil, err
		}

		logrus.WithField("nodeName", nodeName).Debug("Executing get_node_conditions handler")

		conditions, err := k8sClient.GetNodeConditions(ctx, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed to get node conditions: %w", err)
		}

		data, err := optimize.GlobalJSONPool.MarshalToBytes(conditions)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// HandleAnalyzeIssue handles AI-powered issue analysis
func HandleAnalyzeIssue(k8sClient *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		issueType, err := requireStringParam(request, "issueType")
		if err != nil {
			return nil, err
		}

		resourceKind, err := requireStringParam(request, "resourceKind")
		if err != nil {
			return nil, err
		}

		resourceName, err := requireStringParam(request, "resourceName")
		if err != nil {
			return nil, err
		}

		namespace := ""
		if v, ok := request.GetArguments()["namespace"].(string); ok {
			namespace = v
		}

		logrus.WithFields(logrus.Fields{
			"issueType":    issueType,
			"resourceKind": resourceKind,
			"resourceName": resourceName,
			"namespace":    namespace,
		}).Debug("Executing analyze_issue handler")

		analysis, err := k8sClient.AnalyzeIssue(ctx, issueType, resourceKind, resourceName, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze issue: %w", err)
		}

		data, err := optimize.GlobalJSONPool.MarshalToBytes(analysis)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// ============ Search Handlers ============

// HandleSearchResources handles fuzzy search for resources by name
func HandleSearchResources(k8sClient *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kinds, err := getOptionalSearchKinds(request)
		if err != nil {
			return nil, err
		}

		query := getOptionalStringParam(request, "query")
		if query == "" {
			query = getOptionalStringParam(request, "name")
		}

		namespace := getOptionalStringParam(request, "namespace")
		searchMode := getOptionalStringParam(request, "searchMode")
		if searchMode == "" {
			searchMode = "contains"
		}

		caseSensitive := getBoolParam(request, "caseSensitive", false)

		limit := int64(50)
		if v, ok := request.GetArguments()["limit"].(float64); ok {
			limit = int64(v)
			if limit <= 0 || limit > 200 {
				limit = 50
			}
		}

		labelSelector := getOptionalStringParam(request, "labelSelector")
		debug := getOptionalStringParam(request, "debug")

		logrus.WithFields(logrus.Fields{
			"tool":          "search_resources",
			"kinds":         kinds,
			"query":         query,
			"namespace":     namespace,
			"searchMode":    searchMode,
			"caseSensitive": caseSensitive,
			"limit":         limit,
			"labelSelector": labelSelector,
			"debug":         debug,
		}).Debug("Handler invoked")

		var matchedResources []map[string]any
		queryStr := query
		if !caseSensitive {
			queryStr = strings.ToLower(query)
		}

		for _, kind := range kinds {
			resources, err := k8sClient.ListResources(ctx, kind, namespace, labelSelector, "")
			if err != nil {
				return nil, fmt.Errorf("failed to list resources for kind %s: %w", kind, err)
			}

			for _, resource := range resources {
				resourceMap := resource
				metadata, ok := resourceMap["metadata"].(map[string]any)
				if !ok {
					continue
				}

				name, ok := metadata["name"].(string)
				if !ok {
					continue
				}

				var matched bool
				searchName := name
				if !caseSensitive {
					searchName = strings.ToLower(name)
				}

				if queryStr == "" {
					matched = true
				} else {
					switch searchMode {
					case "contains":
						matched = strings.Contains(searchName, queryStr)
					case "startsWith":
						matched = strings.HasPrefix(searchName, queryStr)
					case "endsWith":
						matched = strings.HasSuffix(searchName, queryStr)
					case "exact":
						matched = searchName == queryStr
					case "regex":
						matched, err = regexMatch(searchName, queryStr)
						if err != nil {
							logrus.WithError(err).Warn("Regex match failed, skipping")
							continue
						}
					default:
						matched = strings.Contains(searchName, queryStr)
					}
				}

				if matched {
					matchedResources = append(matchedResources, resourceMap)
				}

				if len(matchedResources) >= int(limit) {
					break
				}
			}

			if len(matchedResources) >= int(limit) {
				break
			}
		}

		response := map[string]interface{}{
			"query":         query,
			"kinds":         kinds,
			"namespace":     namespace,
			"searchMode":    searchMode,
			"caseSensitive": caseSensitive,
			"matched":       len(matchedResources),
			"resources":     matchedResources,
		}

		if len(kinds) == 1 {
			response["kind"] = kinds[0]
		}

		data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		logrus.WithField("matchedCount", len(matchedResources)).Debug("search_resources succeeded")
		return mcp.NewToolResultText(string(data)), nil
	}
}

// regexMatch performs basic regex matching
func regexMatch(text, pattern string) (bool, error) {
	// Convert simple wildcard pattern to regex
	// * matches any sequence, ? matches any single character
	regexPattern := ""
	for _, char := range pattern {
		switch char {
		case '*':
			regexPattern += ".*"
		case '?':
			regexPattern += "."
		case '.', '^', '$', '+', '(', ')', '[', ']', '{', '}', '|', '\\':
			// Escape regex special characters
			regexPattern += "\\" + string(char)
		default:
			regexPattern += string(char)
		}
	}

	// Add anchors for exact matching
	regexPattern = "^" + regexPattern + "$"

	// Use Go's regex for matching
	matched, err := regexp.MatchString(regexPattern, text)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return matched, nil
}
