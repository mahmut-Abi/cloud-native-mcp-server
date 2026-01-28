package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/client"
)

// HandleSearchCharts returns a handler function for searching Helm charts.
func HandleSearchCharts(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_search_charts").Debug("Handler invoked")
		keyword, err := requireStringParam(request, "keyword")
		if err != nil {
			return nil, err
		}
		devel := getOptionalBoolParam(request, "devel")

		// Create a context with 2 minute timeout for chart search
		searchCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan struct {
			charts []map[string]interface{}
			err    error
		}, 1)

		// Run the search operation in a goroutine
		go func() {
			charts, err := c.SearchChartsAsMap(keyword, devel)
			resultChan <- struct {
				charts []map[string]interface{}
				err    error
			}{charts, err}
		}()

		// Wait for either the operation to complete or the context to timeout
		var charts []map[string]interface{}
		select {
		case result := <-resultChan:
			charts, err = result.charts, result.err
			if err != nil {
				return nil, fmt.Errorf("failed to search charts with keyword %s: %w", keyword, err)
			}
		case <-searchCtx.Done():
			return nil, fmt.Errorf("chart search timed out after 2 minutes")
		}

		logrus.WithField("keyword", keyword).Debug("helm_search_charts succeeded")
		jsonData, err := marshalIndentJSON(charts)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetChartInfo returns a handler function for getting Helm chart info.
func HandleGetChartInfo(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_chart_info").Debug("Handler invoked")
		chart, err := requireStringParam(request, "chart")
		if err != nil {
			return nil, err
		}
		info, err := c.GetChartInfoAsMap(chart)
		if err != nil {
			return nil, fmt.Errorf("failed to get chart info for %s: %w", chart, err)
		}
		logrus.WithField("chart", chart).Debug("helm_get_chart_info succeeded")
		jsonData, err := marshalIndentJSON(info)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleTemplateChart returns a handler function for templating a Helm chart.
func HandleTemplateChart(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_template_chart").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		chart, err := requireStringParam(request, "chart")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		valuesFile := getOptionalStringParam(request, "values_file")

		// Create a context with 2 minute timeout for chart templating
		templateCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan struct {
			manifest string
			err      error
		}, 1)

		// Run the template operation in a goroutine
		go func() {
			manifest, err := c.TemplateChart(name, chart, namespace, valuesFile)
			resultChan <- struct {
				manifest string
				err      error
			}{manifest, err}
		}()

		// Wait for either the operation to complete or the context to timeout
		var manifest string
		select {
		case result := <-resultChan:
			manifest, err = result.manifest, result.err
			if err != nil {
				return nil, fmt.Errorf("failed to template chart %s in namespace %s: %w", chart, namespace, err)
			}
		case <-templateCtx.Done():
			return nil, fmt.Errorf("chart templating timed out after 2 minutes")
		}

		logrus.WithField("chart", chart).Debug("helm_template_chart succeeded")
		return mcp.NewToolResultText(manifest), nil
	}
}
