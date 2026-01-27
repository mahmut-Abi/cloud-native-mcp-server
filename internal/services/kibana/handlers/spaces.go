// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains space-related handlers.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

// HandleGetSpaces handles Kibana spaces retrieval requests.
func HandleGetSpaces(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get spaces handler")

		// Get spaces
		spaces, err := c.GetSpaces(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get spaces: %v", err)),
				},
			}, nil
		}

		// Format result using optimized JSON pool
		resultJSON, err := marshalIndentJSON(spaces)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format spaces: %v", err)),
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

// HandleGetSpace handles specific Kibana space retrieval requests.
func HandleGetSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get space handler")

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

		// Get space ID parameter
		spaceID, ok := req.GetArguments()["space_id"].(string)
		if !ok || spaceID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Space ID parameter is required"),
				},
			}, nil
		}

		// Get space
		space, err := c.GetSpace(ctx, spaceID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get space: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(space)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format space: %v", err)),
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

// HandleCreateSpace handles creating a new Kibana space
func HandleCreateSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create space handler")

		spaceID, _ := req.GetArguments()["id"].(string)
		name, _ := req.GetArguments()["name"].(string)
		description, _ := req.GetArguments()["description"].(string)
		color, _ := req.GetArguments()["color"].(string)
		initials, _ := req.GetArguments()["initials"].(string)

		var disabledFeatures []string
		if feats, ok := req.GetArguments()["disabledFeatures"].([]interface{}); ok {
			for _, f := range feats {
				if s, ok := f.(string); ok {
					disabledFeatures = append(disabledFeatures, s)
				}
			}
		}

		space := client.Space{
			ID:               spaceID,
			Name:             name,
			Description:      description,
			Color:            color,
			Initials:         initials,
			DisabledFeatures: disabledFeatures,
		}

		created, err := c.CreateSpace(ctx, space)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create space: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(created)
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

// HandleUpdateSpace handles updating an existing Kibana space
func HandleUpdateSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update space handler")

		spaceID, _ := req.GetArguments()["space_id"].(string)
		name, _ := req.GetArguments()["name"].(string)
		description, _ := req.GetArguments()["description"].(string)
		color, _ := req.GetArguments()["color"].(string)
		initials, _ := req.GetArguments()["initials"].(string)

		var disabledFeatures []string
		if feats, ok := req.GetArguments()["disabledFeatures"].([]interface{}); ok {
			for _, f := range feats {
				if s, ok := f.(string); ok {
					disabledFeatures = append(disabledFeatures, s)
				}
			}
		}

		space := client.Space{
			Name:             name,
			Description:      description,
			Color:            color,
			Initials:         initials,
			DisabledFeatures: disabledFeatures,
		}

		updated, err := c.UpdateSpace(ctx, spaceID, space)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update space: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(updated)
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

// HandleDeleteSpace handles deleting a Kibana space
func HandleDeleteSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete space handler")

		spaceID, _ := req.GetArguments()["space_id"].(string)
		force := false
		if f, ok := req.GetArguments()["force"].(bool); ok {
			force = f
		}

		if spaceID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("space_id is required"),
				},
			}, nil
		}

		err := c.DeleteSpace(ctx, spaceID, force)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete space: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted space: %s", spaceID)),
			},
		}, nil
	}
}

// HandleSpacesSummary handles getting spaces summary with LLM optimization
func HandleSpacesSummary(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := parseLimitWithWarnings(request, "kibana_spaces_summary")

		logrus.WithFields(logrus.Fields{
			"tool":  "kibana_spaces_summary",
			"limit": limit,
		}).Debug("Handler invoked")

		spaces, err := c.SpacesSummary(ctx, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get spaces summary: %w", err)
		}

		response := map[string]interface{}{
			"spaces": spaces,
			"count":  len(spaces),
			"metadata": map[string]interface{}{
				"tool":         "kibana_spaces_summary",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "kibana_spaces_summary")
	}
}
