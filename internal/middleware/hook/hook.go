package hook

import (
	"context"
	"strings"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/observability/metrics"
	mcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

var toolCallStartTimes = make(map[any]time.Time)

// SessionRegisterHookFunc creates a hook function for session registration events
func SessionRegisterHookFunc() server.OnRegisterSessionHookFunc {
	return func(ctx context.Context, session server.ClientSession) {
		logrus.WithFields(logrus.Fields{
			"sessionID":   session.SessionID(),
			"initialized": session.Initialized(),
		}).Info("Session registered")

		// Debug information for session registration
		logrus.WithFields(logrus.Fields{
			"sessionID":   session.SessionID(),
			"initialized": session.Initialized(),
			"context":     ctx,
		}).Debug("Session registration debug info")
	}
}

// LogRequestHookFunc creates a hook function to log incoming tool requests
func LogRequestHookFunc() server.OnBeforeCallToolFunc {
	return func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		// Record start time for metrics
		toolCallStartTimes[id] = time.Now()

		fields := logrus.Fields{
			"id":     id,
			"method": message.Method,
		}

		if message.Params.Name != "" {
			fields["tool_name"] = message.Params.Name
			if message.Params.Meta != nil {
				fields["meta"] = message.Params.Meta
			}
			if message.Params.Arguments != nil {
				fields["arguments"] = message.Params.Arguments
			}
		}

		logrus.WithFields(fields).Info("Tool request received")

		// Debug information for tool requests
		debugFields := logrus.Fields{
			"id":             id,
			"method":         message.Method,
			"full_message":   message,
			"context_values": ctx,
		}
		logrus.WithFields(debugFields).Debug("Tool request debug info")
	}
}

// LogResponseHookFunc creates a hook function to log tool responses
func LogResponseHookFunc() server.OnAfterCallToolFunc {
	return func(ctx context.Context, id any, message *mcp.CallToolRequest, result any) {
		// Calculate duration and record metrics
		duration := 0.0
		if startTime, ok := toolCallStartTimes[id]; ok {
			duration = time.Since(startTime).Seconds()
			delete(toolCallStartTimes, id)
		}

		// Determine status and service name
		status := "success"
		hasError := false
		contentItems := 0

		switch typedResult := result.(type) {
		case *mcp.CallToolResult:
			if typedResult != nil {
				hasError = typedResult.IsError
				if typedResult.Content != nil {
					contentItems = len(typedResult.Content)
				}
			}
		case mcp.CallToolResult:
			hasError = typedResult.IsError
			if typedResult.Content != nil {
				contentItems = len(typedResult.Content)
			}
		}

		if hasError {
			status = "error"
		}

		serviceName := "unknown"
		toolName := "unknown"
		if message != nil && message.Params.Name != "" {
			toolName = message.Params.Name
			// Extract service name from tool name (e.g., "kubernetes_list_pods" -> "kubernetes")
			parts := strings.Split(toolName, "_")
			if len(parts) > 0 {
				serviceName = parts[0]
			}
		}

		// Record metrics
		metrics.RecordToolCall(serviceName, toolName, status, duration)

		fields := logrus.Fields{
			"id":       id,
			"hasError": hasError,
		}

		if message != nil && message.Params.Name != "" {
			fields["tool_name"] = message.Params.Name
		}

		if contentItems > 0 {
			fields["content_items"] = contentItems
		}

		logrus.WithFields(fields).Info("Tool response sent")

		// Debug information for tool responses
		debugFields := logrus.Fields{
			"id":              id,
			"hasError":        hasError,
			"full_result":     result,
			"request_message": message,
			"context_values":  ctx,
		}

		if callToolResult, ok := result.(*mcp.CallToolResult); ok && callToolResult != nil && callToolResult.Content != nil {
			debugFields["response_content"] = callToolResult.Content
		}
		logrus.WithFields(debugFields).Debug("Tool response debug info")
	}
}

// InitializationHookFunc creates a hook function for initialization requests
func InitializationHookFunc() server.OnRequestInitializationFunc {
	return func(ctx context.Context, id any, message any) error {
		logrus.WithFields(logrus.Fields{
			"id":          id,
			"messageType": "initialization",
		}).Info("Initialization request received")

		// Debug information for initialization requests
		logrus.WithFields(logrus.Fields{
			"id":             id,
			"messageType":    "initialization",
			"full_message":   message,
			"context_values": ctx,
		}).Debug("Initialization request debug info")

		return nil
	}
}
