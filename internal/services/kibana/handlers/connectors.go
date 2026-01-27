// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains connector-related handlers.
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

// HandleGetConnectors handles listing connectors.
func HandleGetConnectors(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)

		logrus.WithFields(logrus.Fields{
			"page":    page,
			"perPage": perPage,
		}).Debug("Executing Kibana get connectors handler")

		connectors, err := c.GetConnectors(ctx, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get connectors: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connectors)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format connectors: %v", err)),
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

// HandleGetConnector handles getting a specific connector.
func HandleGetConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana get connector handler")

		connector, err := c.GetConnector(ctx, connectorID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get connector: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connector)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format connector: %v", err)),
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

// HandleCreateConnector handles creating a new connector.
func HandleCreateConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := getOptionalStringParam(req, "name")
		connectorTypeID := getOptionalStringParam(req, "connectorTypeId")

		var config, secrets map[string]interface{}
		if c, ok := req.GetArguments()["config"].(map[string]interface{}); ok {
			config = c
		}
		if s, ok := req.GetArguments()["secrets"].(map[string]interface{}); ok {
			secrets = s
		}

		if name == "" || connectorTypeID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("name and connectorTypeId are required"),
				},
			}, nil
		}

		logrus.WithField("name", name).Debug("Executing Kibana create connector handler")

		connector, err := c.CreateConnector(ctx, name, connectorTypeID, config, secrets)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create connector: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connector)
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

// HandleUpdateConnector handles updating an existing connector.
func HandleUpdateConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		name := getOptionalStringParam(req, "name")

		var config, secrets map[string]interface{}
		if c, ok := req.GetArguments()["config"].(map[string]interface{}); ok {
			config = c
		}
		if s, ok := req.GetArguments()["secrets"].(map[string]interface{}); ok {
			secrets = s
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana update connector handler")

		connector, err := c.UpdateConnector(ctx, connectorID, name, config, secrets)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update connector: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connector)
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

// HandleDeleteConnector handles deleting a connector.
func HandleDeleteConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana delete connector handler")

		err = c.DeleteConnector(ctx, connectorID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete connector: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted connector: %s", connectorID)),
			},
		}, nil
	}
}

// HandleTestConnector handles testing a connector.
func HandleTestConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		var body map[string]interface{}
		if b, ok := req.GetArguments()["body"].(map[string]interface{}); ok {
			body = b
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana test connector handler")

		err = c.TestConnector(ctx, connectorID, body)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to test connector: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully tested connector: %s", connectorID)),
			},
		}, nil
	}
}

// HandleGetConnectorTypes handles listing available connector types.
func HandleGetConnectorTypes(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get connector types handler")

		connectorTypes, err := c.GetConnectorTypes(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get connector types: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connectorTypes)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format connector types: %v", err)),
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
