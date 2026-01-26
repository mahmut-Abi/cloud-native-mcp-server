package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateDatasource(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources" {
			t.Errorf("Expected path /api/datasources, got %s", r.URL.Path)
		}

		// Verify authorization header
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Error("Expected Authorization header")
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
	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test CreateDatasource
	req := CreateDatasourceRequest{
		Name:      "Test Datasource",
		Type:      "prometheus",
		URL:       "http://localhost:9090",
		Access:    "proxy",
		IsDefault: false,
	}

	datasource, err := client.CreateDatasource(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateDatasource failed: %v", err)
	}

	// Verify response
	if datasource.ID != 1 {
		t.Errorf("Expected ID 1, got %d", datasource.ID)
	}
	if datasource.UID != "test-uid" {
		t.Errorf("Expected UID test-uid, got %s", datasource.UID)
	}
	if datasource.Name != "Test Datasource" {
		t.Errorf("Expected name Test Datasource, got %s", datasource.Name)
	}
	if datasource.Type != "prometheus" {
		t.Errorf("Expected type prometheus, got %s", datasource.Type)
	}
}

func TestCreateDatasourceWithJSONData(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request body contains jsonData
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if jsonData, ok := reqBody["jsonData"].(map[string]interface{}); ok {
			if jsonData["httpMethod"] != "POST" {
				t.Errorf("Expected httpMethod POST in jsonData, got %v", jsonData["httpMethod"])
			}
		} else {
			t.Error("Expected jsonData in request")
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"id":   2,
			"uid":  "test-uid-2",
			"name": "Test Datasource 2",
			"type": "prometheus",
			"url":  "http://localhost:9090",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test CreateDatasource with jsonData
	req := CreateDatasourceRequest{
		Name: "Test Datasource 2",
		Type: "prometheus",
		URL:  "http://localhost:9090",
		JSONData: map[string]interface{}{
			"httpMethod": "POST",
		},
	}

	datasource, err := client.CreateDatasource(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateDatasource failed: %v", err)
	}

	if datasource.ID != 2 {
		t.Errorf("Expected ID 2, got %d", datasource.ID)
	}
}

func TestCreateDatasourceError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"message": "Invalid datasource configuration",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test CreateDatasource with error
	req := CreateDatasourceRequest{
		Name: "Test Datasource",
		Type: "prometheus",
		URL:  "http://localhost:9090",
	}

	_, err = client.CreateDatasource(context.Background(), req)
	if err == nil {
		t.Error("Expected error from CreateDatasource, got nil")
	}
}

func TestUpdateDatasource(t *testing.T) {
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
	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test UpdateDatasource
	req := UpdateDatasourceRequest{
		UID:       "test-uid",
		Name:      "Updated Datasource",
		Type:      "prometheus",
		URL:       "http://localhost:9091",
		Access:    "proxy",
		IsDefault: true,
	}

	datasource, err := client.UpdateDatasource(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateDatasource failed: %v", err)
	}

	// Verify response
	if datasource.Name != "Updated Datasource" {
		t.Errorf("Expected name Updated Datasource, got %s", datasource.Name)
	}
	if datasource.URL != "http://localhost:9091" {
		t.Errorf("Expected URL http://localhost:9091, got %s", datasource.URL)
	}
	if datasource.IsDefault != true {
		t.Errorf("Expected IsDefault true, got %v", datasource.IsDefault)
	}
}

func TestUpdateDatasourceError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"message": "Datasource not found",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test UpdateDatasource with error
	req := UpdateDatasourceRequest{
		UID:  "non-existent-uid",
		Name: "Test Datasource",
		Type: "prometheus",
		URL:  "http://localhost:9090",
	}

	_, err = client.UpdateDatasource(context.Background(), req)
	if err == nil {
		t.Error("Expected error from UpdateDatasource, got nil")
	}
}

func TestDeleteDatasource(t *testing.T) {
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
	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test DeleteDatasource
	err = client.DeleteDatasource(context.Background(), "test-uid")
	if err != nil {
		t.Fatalf("DeleteDatasource failed: %v", err)
	}
}

func TestDeleteDatasourceError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"message": "Datasource not found",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test DeleteDatasource with error
	err = client.DeleteDatasource(context.Background(), "non-existent-uid")
	if err == nil {
		t.Error("Expected error from DeleteDatasource, got nil")
	}
}
