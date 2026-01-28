// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
)

// HandleGetRules handles Prometheus rules retrieval requests.
func HandleGetRules(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get rules handler")

		// Extract optional rule type filter
		var ruleType string
		if req.GetArguments() != nil {
			if t, exists := req.GetArguments()["type"]; exists {
				if typeStr, ok := t.(string); ok {
					ruleType = typeStr
				}
			}
		}

		// Get rules
		rules, err := c.GetRules(ctx, ruleType)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get rules: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(rules)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format rules: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetRulesSummary handles Prometheus rules summary requests (optimized version).
func HandleGetRulesSummary(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleType := ""
		if t, ok := req.GetArguments()["type"].(string); ok {
			ruleType = t
		}

		logrus.WithField("type", ruleType).Debug("Executing Prometheus rules summary handler")

		rules, err := c.GetRules(ctx, ruleType)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get rules: %v", err)),
				},
			}, nil
		}

		// Create summary
		summary := make([]map[string]interface{}, 0)
		for _, rg := range rules {
			for _, r := range rg.Rules {
				summary = append(summary, map[string]interface{}{
					"name":   r.Name,
					"type":   r.Type,
					"health": r.Health,
				})
			}
		}

		result := map[string]interface{}{
			"count": len(summary),
			"rules": summary,
		}

		resultJSON, err := marshalIndentJSON(result)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format result: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}
