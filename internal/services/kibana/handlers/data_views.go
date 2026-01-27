// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains data view-related handlers.
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

// HandleGetDataViews handles listing data views.
func HandleGetDataViews(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)

		logrus.WithFields(logrus.Fields{
			"page":    page,
			"perPage": perPage,
		}).Debug("Executing Kibana get data views handler")

		dataViews, err := c.GetDataViews(ctx, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get data views: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataViews)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format data views: %v", err)),
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

// HandleGetDataView handles getting a specific data view.
func HandleGetDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataViewID, err := requireStringParam(req, "data_view_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("data_view_id", dataViewID).Debug("Executing Kibana get data view handler")

		dataView, err := c.GetDataView(ctx, dataViewID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get data view: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataView)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format data view: %v", err)),
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

// HandleCreateDataView handles creating a new data view.
func HandleCreateDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title := getOptionalStringParam(req, "title")
		name := getOptionalStringParam(req, "name")
		timeField := getOptionalStringParam(req, "timeField")
		allowNoIndex := false
		if ani, ok := req.GetArguments()["allowNoIndex"].(bool); ok {
			allowNoIndex = ani
		}

		if title == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("title is required"),
				},
			}, nil
		}

		logrus.WithField("title", title).Debug("Executing Kibana create data view handler")

		dataView, err := c.CreateDataView(ctx, title, name, timeField, nil, nil, allowNoIndex)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create data view: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataView)
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

// HandleUpdateDataView handles updating an existing data view.
func HandleUpdateDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataViewID, err := requireStringParam(req, "data_view_id")
		if err != nil {
			return nil, err
		}

		title := getOptionalStringParam(req, "title")
		name := getOptionalStringParam(req, "name")
		timeField := getOptionalStringParam(req, "timeField")

		logrus.WithField("data_view_id", dataViewID).Debug("Executing Kibana update data view handler")

		dataView, err := c.UpdateDataView(ctx, dataViewID, title, name, timeField)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update data view: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataView)
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

// HandleDeleteDataView handles deleting a data view.
func HandleDeleteDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataViewID, err := requireStringParam(req, "data_view_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("data_view_id", dataViewID).Debug("Executing Kibana delete data view handler")

		err = c.DeleteDataView(ctx, dataViewID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete data view: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted data view: %s", dataViewID)),
			},
		}, nil
	}
}
