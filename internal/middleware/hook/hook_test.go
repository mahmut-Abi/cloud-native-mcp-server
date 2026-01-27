package hook

import (
	"context"
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/observability/metrics"
	mcp "github.com/mark3labs/mcp-go/mcp"
	dto "github.com/prometheus/client_model/go"
)

func init() {
	// Initialize metrics for all tests
	metrics.Init("test", "test", "go1.25", "sse", "0.0.0.0:8080")
}

func TestSessionRegisterHookFunc(t *testing.T) {
	hook := SessionRegisterHookFunc()
	if hook == nil {
		t.Error("Expected hook function, got nil")
	}

	mockSession := &mockClientSession{
		sessionID:   "test-session-123",
		initialized: true,
	}

	ctx := context.Background()
	hook(ctx, mockSession)
}

func TestLogRequestHookFunc(t *testing.T) {
	hook := LogRequestHookFunc()
	if hook == nil {
		t.Error("Expected hook function, got nil")
	}

	req := &mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "test-tool",
		},
	}

	ctx := context.Background()
	hook(ctx, "test-id", req)
}

func TestLogResponseHookFunc(t *testing.T) {
	hook := LogResponseHookFunc()
	if hook == nil {
		t.Error("Expected hook function, got nil")
	}

	req := &mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "test-tool",
		},
	}

	result := &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "test result",
			},
		},
	}

	ctx := context.Background()
	hook(ctx, "test-id", req, result)
}

func TestInitializationHookFunc(t *testing.T) {
	hook := InitializationHookFunc()
	if hook == nil {
		t.Error("Expected hook function, got nil")
	}

	ctx := context.Background()
	err := hook(ctx, "test-id", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

type mockClientSession struct {
	sessionID   string
	initialized bool
	notifyChan  chan mcp.JSONRPCNotification
}

func (m *mockClientSession) SessionID() string {
	return m.sessionID
}

func (m *mockClientSession) Initialized() bool {
	return m.initialized
}

func (m *mockClientSession) Initialize() {
	m.initialized = true
}

func (m *mockClientSession) NotificationChannel() chan<- mcp.JSONRPCNotification {
	if m.notifyChan == nil {
		m.notifyChan = make(chan mcp.JSONRPCNotification, 10)
	}
	return m.notifyChan
}

func TestLogRequestHookFunc_RecordsStartTime(t *testing.T) {
	hook := LogRequestHookFunc()
	req := &mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "kubernetes_list_pods",
		},
	}

	ctx := context.Background()
	requestID := "test-request-id"

	// Call request hook
	hook(ctx, requestID, req)

	// Verify that start time was recorded
	if _, exists := toolCallStartTimes[requestID]; !exists {
		t.Error("Expected start time to be recorded for request ID")
	}

	// Clean up
	delete(toolCallStartTimes, requestID)
}

func TestLogResponseHookFunc_RecordsMetrics(t *testing.T) {
	requestHook := LogRequestHookFunc()
	responseHook := LogResponseHookFunc()

	req := &mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "kubernetes_list_pods",
		},
	}

	result := &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "test result",
			},
		},
	}

	ctx := context.Background()
	requestID := "test-request-id"

	// Call request hook
	requestHook(ctx, requestID, req)

	// Call response hook
	responseHook(ctx, requestID, req, result)

	// Verify metrics were recorded
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check tool_calls_total
	found := false
	for _, m := range metricsList {
		if m.GetName() == "tool_calls_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "kubernetes") &&
					hasLabels(labels, "tool_name", "kubernetes_list_pods") &&
					hasLabels(labels, "status", "success") {
					found = true
					if metric.Counter.GetValue() != 1 {
						t.Errorf("tool_calls_total = %v, want 1", metric.Counter.GetValue())
					}
					break
				}
			}
			break
		}
	}
	if !found {
		t.Error("tool_calls_total metric not found")
	}

	// Check tool_call_duration_seconds
	foundDuration := false
	for _, m := range metricsList {
		if m.GetName() == "tool_call_duration_seconds" {
			foundDuration = true
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "kubernetes") &&
					hasLabels(labels, "tool_name", "kubernetes_list_pods") &&
					hasLabels(labels, "status", "success") {
					if metric.Histogram.GetSampleCount() != 1 {
						t.Errorf("tool_call_duration_seconds count = %v, want 1", metric.Histogram.GetSampleCount())
					}
					break
				}
			}
			break
		}
	}
	if !foundDuration {
		t.Error("tool_call_duration_seconds metric not found")
	}
}

func TestLogResponseHookFunc_ErrorStatus(t *testing.T) {
	requestHook := LogRequestHookFunc()
	responseHook := LogResponseHookFunc()

	req := &mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "helm_list_releases",
		},
	}

	result := &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "error occurred",
			},
		},
	}

	ctx := context.Background()
	requestID := "test-error-id"

	// Call hooks
	requestHook(ctx, requestID, req)
	responseHook(ctx, requestID, req, result)

	// Verify metrics with error status
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metricsList {
		if m.GetName() == "tool_calls_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "helm") &&
					hasLabels(labels, "tool_name", "helm_list_releases") &&
					hasLabels(labels, "status", "error") {
					found = true
					if metric.Counter.GetValue() != 1 {
						t.Errorf("tool_calls_total = %v, want 1", metric.Counter.GetValue())
					}
					break
				}
			}
			break
		}
	}
	if !found {
		t.Error("tool_calls_total metric with error status not found")
	}
}

func TestLogResponseHookFunc_UnknownToolName(t *testing.T) {
	requestHook := LogRequestHookFunc()
	responseHook := LogResponseHookFunc()

	req := &mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "", // Empty tool name
		},
	}

	result := &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{},
	}

	ctx := context.Background()
	requestID := "test-unknown-id"

	// Call hooks
	requestHook(ctx, requestID, req)
	responseHook(ctx, requestID, req, result)

	// Verify metrics with unknown service/tool
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metricsList {
		if m.GetName() == "tool_calls_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "unknown") &&
					hasLabels(labels, "tool_name", "unknown") &&
					hasLabels(labels, "status", "success") {
					found = true
					break
				}
			}
			break
		}
	}
	if !found {
		t.Error("tool_calls_total metric with unknown service not found")
	}
}

func TestLogResponseHookFunc_NilRequest(t *testing.T) {
	responseHook := LogResponseHookFunc()

	result := &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{},
	}

	ctx := context.Background()
	requestID := "test-nil-request"

	// Call response hook with nil request (should not panic)
	responseHook(ctx, requestID, nil, result)

	// Verify metrics with unknown service/tool
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metricsList {
		if m.GetName() == "tool_calls_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "unknown") &&
					hasLabels(labels, "tool_name", "unknown") &&
					hasLabels(labels, "status", "success") {
					found = true
					break
				}
			}
			break
		}
	}
	if !found {
		t.Error("tool_calls_total metric with unknown service not found for nil request")
	}
}

func TestLogResponseHookFunc_MultipleCalls(t *testing.T) {

	requestHook := LogRequestHookFunc()

	responseHook := LogResponseHookFunc()

	// Make multiple calls

	for i := 0; i < 3; i++ {

		req := &mcp.CallToolRequest{

			Params: mcp.CallToolParams{

				Name: "kubernetes_list_pods",
			},
		}

		result := &mcp.CallToolResult{

			IsError: false,

			Content: []mcp.Content{},
		}

		ctx := context.Background()

		requestID := "test-request-id-multiple"

		requestHook(ctx, requestID, req)

		responseHook(ctx, requestID, req, result)

	}

	// Verify counter is at least 3 (may have more from previous tests)

	metricsList, err := metrics.Registry.Gather()

	if err != nil {

		t.Fatalf("Failed to gather metrics: %v", err)

	}

	for _, m := range metricsList {

		if m.GetName() == "tool_calls_total" {

			for _, metric := range m.GetMetric() {

				labels := metric.GetLabel()

				if hasLabels(labels, "service_name", "kubernetes") &&

					hasLabels(labels, "tool_name", "kubernetes_list_pods") &&

					hasLabels(labels, "status", "success") {

					if metric.Counter.GetValue() < 3 {

						t.Errorf("tool_calls_total = %v, want at least 3", metric.Counter.GetValue())

					}

					return

				}

			}

		}

	}

	t.Error("Expected metric not found")

}

func hasLabels(labels []*dto.LabelPair, name, value string) bool {
	for _, label := range labels {
		if label.GetName() == name && label.GetValue() == value {
			return true
		}
	}
	return false
}
