package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/client"
)

// HandleListRepositories returns a handler function for listing Helm repositories.
func HandleListRepositories(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "list_helm_repos").Debug("Handler invoked")

		// Create a context with 2 minute timeout for repository listing
		listCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan struct {
			repos []*client.Repository
			err   error
		}, 1)

		// Run the list operation in a goroutine
		go func() {
			repos, err := c.ListRepositories()
			resultChan <- struct {
				repos []*client.Repository
				err   error
			}{repos, err}
		}()

		// Wait for either the operation to complete or the context to timeout
		var repos []*client.Repository
		var err error
		select {
		case result := <-resultChan:
			repos, err = result.repos, result.err
			if err != nil {
				return nil, fmt.Errorf("failed to list repositories: %w", err)
			}
		case <-listCtx.Done():
			return nil, fmt.Errorf("repository listing timed out after 2 minutes")
		}

		repoMaps := make([]map[string]interface{}, len(repos))
		for i, r := range repos {
			if r != nil {
				repoMaps[i] = map[string]interface{}{
					"name": r.Name,
					"url":  r.URL,
				}
			}
		}

		logrus.WithField("count", len(repos)).Debug("list_helm_repos succeeded")

		// Serialize to JSON for better readability
		jsonData, err := marshalIndentJSON(repoMaps)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleAddRepository returns a handler function for adding a Helm repository.
func HandleAddRepository(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_add_repository").Debug("Handler invoked")
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		url, err := requireStringParam(request, "url")
		if err != nil {
			return nil, err
		}

		// Create a context with 2 minute timeout for repository addition
		addCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan error, 1)

		// Run the add operation in a goroutine
		go func() {
			resultChan <- c.AddRepository(name, url)
		}()

		// Wait for either the operation to complete or the context to timeout
		select {
		case err := <-resultChan:
			if err != nil {
				return nil, fmt.Errorf("failed to add repository %s: %w", name, err)
			}
		case <-addCtx.Done():
			return nil, fmt.Errorf("repository addition timed out after 2 minutes")
		}

		logrus.WithField("repository", name).Debug("helm_add_repository succeeded")
		message := fmt.Sprintf("Successfully added Helm repository '%s' with URL '%s'", name, url)
		return mcp.NewToolResultText(message), nil
	}
}

// HandleRemoveRepository returns a handler function for removing a Helm repository.
func HandleRemoveRepository(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_remove_repository").Debug("Handler invoked")
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Create a context with 2 minute timeout for repository removal
		removeCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan error, 1)

		// Run the remove operation in a goroutine
		go func() {
			resultChan <- c.RemoveRepository(name)
		}()

		// Wait for either the operation to complete or the context to timeout
		select {
		case err := <-resultChan:
			if err != nil {
				return nil, fmt.Errorf("failed to remove repository %s: %w", name, err)
			}
		case <-removeCtx.Done():
			return nil, fmt.Errorf("repository removal timed out after 2 minutes")
		}

		logrus.WithField("repository", name).Debug("helm_remove_repository succeeded")
		message := fmt.Sprintf("Successfully removed Helm repository '%s'", name)
		return mcp.NewToolResultText(message), nil
	}
}

// HandleUpdateRepositories returns a handler function for updating Helm repositories.
func HandleUpdateRepositories(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_update_repositories").Debug("Handler invoked")

		// Create a context with 5 minute timeout for repository updates
		updateCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan error, 1)

		// Run the update in a goroutine
		go func() {
			resultChan <- c.UpdateRepositories()
		}()

		// Wait for either the update to complete or the context to timeout
		select {
		case err := <-resultChan:
			if err != nil {
				return nil, fmt.Errorf("failed to update repositories: %w", err)
			}
		case <-updateCtx.Done():
			return nil, fmt.Errorf("repository update timed out after 5 minutes")
		}

		logrus.Debug("helm_update_repositories succeeded")
		message := "Successfully updated all Helm repositories"
		return mcp.NewToolResultText(message), nil
	}
}
