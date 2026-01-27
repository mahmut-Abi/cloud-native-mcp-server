// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains index pattern-related handlers.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

// HandleGetIndexPatterns handles Kibana index patterns retrieval requests.
func HandleGetIndexPatterns(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get index patterns handler")

		// Get index patterns
		indexPatterns, err := c.GetIndexPatterns(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get index patterns: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(indexPatterns)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format index patterns: %v", err)),
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

// HandleGetIndexPattern handles specific Kibana index pattern retrieval requests.
func HandleGetIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get index pattern handler")

		// Extract parameters
		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		// Get index pattern ID parameter
		indexPatternID, ok := req.GetArguments()["index_pattern_id"].(string)
		if !ok || indexPatternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Index pattern ID parameter is required"),
				},
			}, nil
		}

		// Get index pattern
		indexPattern, err := c.GetIndexPattern(ctx, indexPatternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get index pattern: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(indexPattern)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format index pattern: %v", err)),
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

// HandleCreateIndexPattern handles creating a new index pattern
func HandleCreateIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create index pattern handler")

		title, _ := req.GetArguments()["title"].(string)
		timeField, _ := req.GetArguments()["timeField"].(string)

		if title == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("title is required"),
				},
			}, nil
		}

		pattern, err := c.CreateIndexPattern(ctx, title, timeField)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create index pattern: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(pattern)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
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

// HandleUpdateIndexPattern handles updating an index pattern
func HandleUpdateIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update index pattern handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)
		title, _ := req.GetArguments()["title"].(string)
		timeField, _ := req.GetArguments()["timeField"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		pattern, err := c.UpdateIndexPattern(ctx, patternID, title, timeField)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update index pattern: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(pattern)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
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

// HandleDeleteIndexPattern handles deleting an index pattern
func HandleDeleteIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete index pattern handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		err := c.DeleteIndexPattern(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete index pattern: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted index pattern: %s", patternID)),
			},
		}, nil
	}
}

// HandleSetDefaultIndexPattern handles setting the default index pattern
func HandleSetDefaultIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana set default index pattern handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		err := c.SetDefaultIndexPattern(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to set default index pattern: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully set default index pattern: %s", patternID)),
			},
		}, nil
	}
}

// HandleRefreshIndexPatternFields handles refreshing index pattern fields
func HandleRefreshIndexPatternFields(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana refresh index pattern fields handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		err := c.RefreshIndexPatternFields(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to refresh index pattern fields: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully refreshed fields for index pattern: %s", patternID)),
			},
		}, nil
	}
}

// HandleGetIndexPatternFields handles index pattern fields retrieval requests.
func HandleGetIndexPatternFields(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patternID := ""
		if v, ok := req.GetArguments()["patternID"]; ok {
			if s, ok := v.(string); ok {
				patternID = s
			}
		}

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("patternID is required"),
				},
			}, nil
		}

		logrus.WithField("patternID", patternID).Debug("Executing Kibana get index pattern fields handler")

		fields, err := c.GetIndexPatternFields(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get index pattern fields: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(fields)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format fields: %v", err)),
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
