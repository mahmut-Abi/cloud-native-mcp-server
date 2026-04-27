package hook

import (
	"context"
	"fmt"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prompts"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// NormalizeToolErrorMiddleware converts handler errors into MCP tool error results.
// This keeps tool failures at the MCP layer instead of surfacing them as transport-level errors.
func NormalizeToolErrorMiddleware() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := next(ctx, request)
			if err != nil {
				logrus.WithError(err).WithField("tool", request.Params.Name).Warn("Tool handler returned error")
				return mcp.NewToolResultError(err.Error()), nil
			}
			if result == nil {
				logrus.WithField("tool", request.Params.Name).Warn("Tool handler returned nil result")
				return mcp.NewToolResultError("tool returned no result"), nil
			}
			return result, nil
		}
	}
}

// PromptLoggingMiddleware logs prompt requests and results.
func PromptLoggingMiddleware() server.PromptHandlerMiddleware {
	return func(next server.PromptHandlerFunc) server.PromptHandlerFunc {
		return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			logrus.WithFields(logrus.Fields{
				"prompt":    request.Params.Name,
				"arguments": request.Params.Arguments,
			}).Info("Prompt request received")

			start := time.Now()
			result, err := next(ctx, request)
			duration := time.Since(start)

			if err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"prompt":   request.Params.Name,
					"duration": duration.String(),
				}).Warn("Prompt handler failed")
				return nil, err
			}

			messageCount := 0
			if result != nil {
				messageCount = len(result.Messages)
			}
			logrus.WithFields(logrus.Fields{
				"prompt":        request.Params.Name,
				"duration":      duration.String(),
				"message_count": messageCount,
			}).Info("Prompt response sent")

			return result, nil
		}
	}
}

// NewTaskHooks returns task lifecycle hooks for observability and diagnostics.
func NewTaskHooks() *server.TaskHooks {
	hooks := &server.TaskHooks{}

	hooks.AddOnTaskCreated(func(ctx context.Context, metrics server.TaskMetrics) {
		logrus.WithFields(logrus.Fields{
			"task_id":    metrics.TaskID,
			"tool":       metrics.ToolName,
			"session_id": metrics.SessionID,
			"status":     metrics.Status,
		}).Info("Task created")
	})

	hooks.AddOnTaskCompleted(func(ctx context.Context, metrics server.TaskMetrics) {
		logrus.WithFields(logrus.Fields{
			"task_id":    metrics.TaskID,
			"tool":       metrics.ToolName,
			"session_id": metrics.SessionID,
			"status":     metrics.Status,
			"duration":   metrics.Duration.String(),
		}).Info("Task completed")
	})

	hooks.AddOnTaskFailed(func(ctx context.Context, metrics server.TaskMetrics) {
		logrus.WithError(metrics.Error).WithFields(logrus.Fields{
			"task_id":    metrics.TaskID,
			"tool":       metrics.ToolName,
			"session_id": metrics.SessionID,
			"status":     metrics.Status,
			"duration":   metrics.Duration.String(),
		}).Warn("Task failed")
	})

	hooks.AddOnTaskCancelled(func(ctx context.Context, metrics server.TaskMetrics) {
		logrus.WithFields(logrus.Fields{
			"task_id":    metrics.TaskID,
			"tool":       metrics.ToolName,
			"session_id": metrics.SessionID,
			"status":     metrics.Status,
			"duration":   metrics.Duration.String(),
		}).Info("Task cancelled")
	})

	return hooks
}

// PromptAvailabilityMiddleware blocks prompts whose backing services are unavailable.
func PromptAvailabilityMiddleware(isKubernetesEnabled func() bool) server.PromptHandlerMiddleware {
	return func(next server.PromptHandlerFunc) server.PromptHandlerFunc {
		return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			if request.Params.Name == prompts.K8sOpsPromptName && !isKubernetesEnabled() {
				return nil, fmt.Errorf("prompt %q is unavailable because the kubernetes service is disabled", prompts.K8sOpsPromptName)
			}
			return next(ctx, request)
		}
	}
}
