package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/grafana/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleCreateDatasource(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources" {
			t.Errorf("Expected path /api/datasources, got %s", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"id":        1,
			"uid":       "test-uid",
			"name":      "Test Datasource",
			"type":      "prometheus",
			"url":       "http://localhost:9090",
			"access":    "proxy",
			"isDefault": false,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create handler
	handler := HandleCreateDatasource(grafanaClient)

	// Create request
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name": "Test Datasource",
				"type": "prometheus",
				"url":  "http://localhost:9090",
			},
		},
	}

	// Call handler
	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("HandleCreateDatasource failed: %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Content) == 0 {
		t.Fatal("Expected content in result")
	}

	// Parse response - use type assertion with proper type
	var response map[string]interface{}
	switch content := result.Content[0].(type) {
	case *mcp.TextContent:
		if err := json.Unmarshal([]byte(content.Text), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
	case mcp.TextContent:
		if err := json.Unmarshal([]byte(content.Text), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
	default:
		t.Fatalf("Expected TextContent in result, got %T", result.Content[0])
	}

	// Verify response contains datasource
	if response["datasource"] == nil {
		t.Error("Expected datasource in response")
	}

	if response["message"] != "Datasource created successfully" {
		t.Errorf("Expected success message, got %v", response["message"])
	}
}

func TestHandleCreateDatasourceWithOptionalParams(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request body contains optional parameters
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if reqBody["database"] != "testdb" {
			t.Errorf("Expected database testdb, got %v", reqBody["database"])
		}
		if reqBody["user"] != "testuser" {
			t.Errorf("Expected user testuser, got %v", reqBody["user"])
		}
		if reqBody["isDefault"] != true {
			t.Errorf("Expected isDefault true, got %v", reqBody["isDefault"])
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"id":   2,
			"uid":  "test-uid-2",
			"name": "Test Datasource 2",
			"type": "mysql",
			"url":  "http://localhost:3306",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create handler
	handler := HandleCreateDatasource(grafanaClient)

	// Create request with optional parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name":      "Test Datasource 2",
				"type":      "mysql",
				"url":       "http://localhost:3306",
				"database":  "testdb",
				"user":      "testuser",
				"isDefault": true,
			},
		},
	}

	// Call handler
	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("HandleCreateDatasource failed: %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
}

func TestHandleCreateDatasourceMissingRequiredParam(t *testing.T) {
	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    "http://localhost:3000",
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create handler
	handler := HandleCreateDatasource(grafanaClient)

	// Create request without required parameter
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name": "Test Datasource",
				// Missing 'type' and 'url'
			},
		},
	}

	// Call handler
	_, err = handler(context.Background(), request)
	if err == nil {
		t.Error("Expected error for missing required parameter")
	}
}

func TestHandleUpdateDatasource(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources/test-uid" {
			t.Errorf("Expected path /api/datasources/test-uid, got %s", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"id":        1,
			"uid":       "test-uid",
			"name":      "Updated Datasource",
			"type":      "prometheus",
			"url":       "http://localhost:9091",
			"access":    "proxy",
			"isDefault": true,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create handler
	handler := HandleUpdateDatasource(grafanaClient)

	// Create request
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"uid":  "test-uid",
				"name": "Updated Datasource",
				"type": "prometheus",
				"url":  "http://localhost:9091",
			},
		},
	}

	// Call handler
	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("HandleUpdateDatasource failed: %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Parse response
	var response map[string]interface{}
	switch content := result.Content[0].(type) {
	case *mcp.TextContent:
		if err := json.Unmarshal([]byte(content.Text), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
	case mcp.TextContent:
		if err := json.Unmarshal([]byte(content.Text), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
	default:
		t.Fatalf("Expected TextContent in result, got %T", result.Content[0])
	}

	if response["message"] != "Datasource updated successfully" {
		t.Errorf("Expected success message, got %v", response["message"])
	}
}

func TestHandleUpdateDatasourceMissingRequiredParam(t *testing.T) {
	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    "http://localhost:3000",
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create handler
	handler := HandleUpdateDatasource(grafanaClient)

	// Create request without required parameter
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"uid": "test-uid",
				// Missing 'name', 'type', and 'url'
			},
		},
	}

	// Call handler
	_, err = handler(context.Background(), request)
	if err == nil {
		t.Error("Expected error for missing required parameter")
	}
}

func TestHandleDeleteDatasource(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources/uid/test-uid" {
			t.Errorf("Expected path /api/datasources/uid/test-uid, got %s", r.URL.Path)
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		response := map[string]string{
			"message": "Datasource deleted",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create handler
	handler := HandleDeleteDatasource(grafanaClient)

	// Create request
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"uid": "test-uid",
			},
		},
	}

	// Call handler
	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("HandleDeleteDatasource failed: %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Parse response
	var response map[string]interface{}
	switch content := result.Content[0].(type) {
	case *mcp.TextContent:
		if err := json.Unmarshal([]byte(content.Text), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
	case mcp.TextContent:
		if err := json.Unmarshal([]byte(content.Text), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
	default:
		t.Fatalf("Expected TextContent in result, got %T", result.Content[0])
	}

	if response["message"] != "Datasource deleted successfully" {
		t.Errorf("Expected success message, got %v", response["message"])
	}

	if response["uid"] != "test-uid" {
		t.Errorf("Expected uid test-uid, got %v", response["uid"])
	}
}

func TestHandleDeleteDatasourceMissingRequiredParam(t *testing.T) {
	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    "http://localhost:3000",
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create handler
	handler := HandleDeleteDatasource(grafanaClient)

	// Create request without required parameter
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	// Call handler
	_, err = handler(context.Background(), request)
	if err == nil {
		t.Error("Expected error for missing required parameter")
	}
}

func TestHandleDatasourceErrorCases(t *testing.T) {
	// Create a test server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"message": "Bad request",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	grafanaClient, err := client.NewClient(&client.ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test CreateDatasource with error
	createHandler := HandleCreateDatasource(grafanaClient)
	createRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name": "Test Datasource",
				"type": "prometheus",
				"url":  "http://localhost:9090",
			},
		},
	}
	_, err = createHandler(context.Background(), createRequest)
	if err == nil {
		t.Error("Expected error from HandleCreateDatasource")
	}

	// Test UpdateDatasource with error
	updateHandler := HandleUpdateDatasource(grafanaClient)
	updateRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"uid":  "test-uid",
				"name": "Test Datasource",
				"type": "prometheus",
				"url":  "http://localhost:9090",
			},
		},
	}
	_, err = updateHandler(context.Background(), updateRequest)
	if err == nil {
		t.Error("Expected error from HandleUpdateDatasource")
	}

	// Test DeleteDatasource with error
	deleteHandler := HandleDeleteDatasource(grafanaClient)
	deleteRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"uid": "test-uid",
			},
		},
	}
	_, err = deleteHandler(context.Background(), deleteRequest)
	if err == nil {
		t.Error("Expected error from HandleDeleteDatasource")
	}
}
