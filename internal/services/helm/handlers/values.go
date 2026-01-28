package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/client"
)

// HandleGetReleaseValues returns a handler function for getting Helm release values.
func HandleGetReleaseValues(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_release_values").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		all := getOptionalBoolParam(request, "all")
		values, err := c.GetReleaseValuesAsMap(name, namespace, all)
		if err != nil {
			return nil, fmt.Errorf("failed to get release values for %s in namespace %s: %w", name, namespace, err)
		}
		logrus.WithField("release", name).Debug("helm_get_release_values succeeded")
		jsonData, err := marshalIndentJSON(values)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}