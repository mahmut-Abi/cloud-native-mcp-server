package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestUpdateDashboardUsesGrafanaSaveEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/db" {
			t.Errorf("Expected path /api/dashboards/db, got %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if _, ok := reqBody["dashboard"].(map[string]interface{}); !ok {
			t.Fatalf("Expected dashboard object in request body, got %#v", reqBody)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"id":      1,
			"uid":     "test-dashboard",
			"slug":    "test-dashboard",
			"version": 1,
			"dashboard": map[string]interface{}{
				"id":    1,
				"uid":   "test-dashboard",
				"title": "Test Dashboard",
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dashboard, err := client.UpdateDashboard(context.Background(), DashboardUpdateRequest{
		Dashboard: map[string]interface{}{
			"title": "Test Dashboard",
			"uid":   "test-dashboard",
		},
		Overwrite: true,
		Message:   "test update",
	})
	if err != nil {
		t.Fatalf("UpdateDashboard failed: %v", err)
	}

	if dashboard.Title != "Test Dashboard" {
		t.Errorf("Expected dashboard title Test Dashboard, got %s", dashboard.Title)
	}
}

func TestUpdateDashboardSynthesizesDashboardFromSaveResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/dashboards/db" {
			t.Errorf("Expected path /api/dashboards/db, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      27,
			"uid":     "synth-dashboard",
			"url":     "/d/synth-dashboard/synth-dashboard",
			"version": 4,
			"status":  "success",
			"slug":    "synth-dashboard",
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dashboard, err := client.UpdateDashboard(context.Background(), DashboardUpdateRequest{
		FolderUID: "ops-folder",
		Overwrite: true,
		Dashboard: map[string]interface{}{
			"title": "Synth Dashboard",
			"uid":   "synth-dashboard",
			"tags":  []interface{}{"mcp", "smoke"},
		},
	})
	if err != nil {
		t.Fatalf("UpdateDashboard failed: %v", err)
	}

	if dashboard.UID != "synth-dashboard" {
		t.Errorf("Expected UID synth-dashboard, got %s", dashboard.UID)
	}
	if dashboard.Title != "Synth Dashboard" {
		t.Errorf("Expected title Synth Dashboard, got %s", dashboard.Title)
	}
	if dashboard.FolderUID != "ops-folder" {
		t.Errorf("Expected folder UID ops-folder, got %s", dashboard.FolderUID)
	}
	if dashboard.Version != 4 {
		t.Errorf("Expected version 4, got %d", dashboard.Version)
	}
	if dashboard.URL != "/d/synth-dashboard/synth-dashboard" {
		t.Errorf("Expected dashboard URL to be populated, got %s", dashboard.URL)
	}
	if len(dashboard.Tags) != 2 || dashboard.Tags[0] != "mcp" || dashboard.Tags[1] != "smoke" {
		t.Errorf("Expected tags to be synthesized from request, got %#v", dashboard.Tags)
	}
	if dashboard.Dashboard == nil {
		t.Errorf("Expected original dashboard payload to be retained")
	}
}

func TestCreateFolderUsesFoldersEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders" {
			t.Errorf("Expected path /api/folders, got %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody["title"] != "Operations" {
			t.Fatalf("Expected title Operations, got %#v", reqBody["title"])
		}
		if reqBody["uid"] != "ops-folder" {
			t.Fatalf("Expected uid ops-folder, got %#v", reqBody["uid"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    12,
			"uid":   "ops-folder",
			"title": "Operations",
			"url":   "/dashboards/f/ops-folder/operations",
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	folder, err := client.CreateFolder(context.Background(), FolderCreateRequest{
		UID:   "ops-folder",
		Title: "Operations",
	})
	if err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	if folder.UID != "ops-folder" {
		t.Errorf("Expected UID ops-folder, got %s", folder.UID)
	}
}

func TestUpdateFolderUsesFolderUIDEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders/ops-folder" {
			t.Errorf("Expected path /api/folders/ops-folder, got %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody["title"] != "Operations Updated" {
			t.Fatalf("Expected updated title, got %#v", reqBody["title"])
		}
		if reqBody["version"] != float64(7) {
			t.Fatalf("Expected version 7, got %#v", reqBody["version"])
		}
		if reqBody["overwrite"] != true {
			t.Fatalf("Expected overwrite=true, got %#v", reqBody["overwrite"])
		}
		if _, ok := reqBody["uid"]; ok {
			t.Fatalf("UID should not be serialized into request body: %#v", reqBody)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      12,
			"uid":     "ops-folder",
			"title":   "Operations Updated",
			"version": 8,
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	folder, err := client.UpdateFolder(context.Background(), FolderUpdateRequest{
		UID:       "ops-folder",
		Title:     "Operations Updated",
		Version:   7,
		Overwrite: true,
	})
	if err != nil {
		t.Fatalf("UpdateFolder failed: %v", err)
	}

	if folder.Version != 8 {
		t.Errorf("Expected version 8, got %d", folder.Version)
	}
}

func TestDeleteFolderIncludesForceDeleteRulesQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders/ops-folder" {
			t.Errorf("Expected path /api/folders/ops-folder, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("forceDeleteRules"); got != "true" {
			t.Fatalf("Expected forceDeleteRules=true, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Folder deleted",
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.DeleteFolder(context.Background(), "ops-folder", true)
	if err != nil {
		t.Fatalf("DeleteFolder failed: %v", err)
	}

	if result["message"] != "Folder deleted" {
		t.Errorf("Expected delete message, got %#v", result["message"])
	}
}

func TestGetDashboardVersionUsesVersionDetailEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/test-dashboard/versions/3" {
			t.Errorf("Expected path /api/dashboards/uid/test-dashboard/versions/3, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":            33,
			"version":       3,
			"createdBy":     "agent",
			"message":       "before migration",
			"parentVersion": 2,
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	version, err := client.GetDashboardVersion(context.Background(), "test-dashboard", 3)
	if err != nil {
		t.Fatalf("GetDashboardVersion failed: %v", err)
	}

	if version.Version != 3 {
		t.Errorf("Expected version 3, got %d", version.Version)
	}
}

func TestGetDashboardVersionsHandlesWrappedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/test-dashboard/versions" {
			t.Errorf("Expected path /api/dashboards/uid/test-dashboard/versions, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("Expected limit=5, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"dashboardUID": "test-dashboard",
			"versions": []map[string]interface{}{
				{
					"id":        2,
					"version":   2,
					"createdBy": "admin",
					"message":   "smoke_update",
				},
				{
					"id":        1,
					"version":   1,
					"createdBy": "admin",
					"message":   "smoke_create",
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	versions, err := client.GetDashboardVersions(context.Background(), "test-dashboard", 5, 0)
	if err != nil {
		t.Fatalf("GetDashboardVersions failed: %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("Expected 2 versions, got %d", len(versions))
	}
	if versions[0].Version != 2 || versions[1].Version != 1 {
		t.Errorf("Unexpected version ordering: %#v", versions)
	}
}

func TestRestoreDashboardVersionUsesRestoreEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/test-dashboard/restore" {
			t.Errorf("Expected path /api/dashboards/uid/test-dashboard/restore, got %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody["version"] != float64(3) {
			t.Fatalf("Expected restore version 3, got %#v", reqBody["version"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Dashboard restored",
			"uid":     "test-dashboard",
			"version": 3,
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.RestoreDashboardVersion(context.Background(), "test-dashboard", 3)
	if err != nil {
		t.Fatalf("RestoreDashboardVersion failed: %v", err)
	}

	if result["message"] != "Dashboard restored" {
		t.Errorf("Expected restore message, got %#v", result["message"])
	}
}

func TestGetPluginsUsesPluginsEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/plugins" {
			t.Errorf("Expected path /api/plugins, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"id":            "grafana-lokiexplore-app",
				"name":          "Grafana Logs Drilldown",
				"type":          "app",
				"enabled":       true,
				"pinned":        true,
				"defaultNavUrl": "/plugins/grafana-lokiexplore-app/page/grafana-logs-drilldown",
			},
			{
				"id":      "prometheus",
				"name":    "Prometheus",
				"type":    "datasource",
				"enabled": true,
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	plugins, err := client.GetPlugins(context.Background())
	if err != nil {
		t.Fatalf("GetPlugins failed: %v", err)
	}
	if len(plugins) != 2 {
		t.Fatalf("Expected 2 plugins, got %d", len(plugins))
	}
	if plugins[0].ID != "grafana-lokiexplore-app" {
		t.Errorf("Expected first plugin id grafana-lokiexplore-app, got %s", plugins[0].ID)
	}
}

func TestGetPluginUsesSettingsEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/plugins/grafana-lokiexplore-app/settings" {
			t.Errorf("Expected path /api/plugins/grafana-lokiexplore-app/settings, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":            "grafana-lokiexplore-app",
			"name":          "Grafana Logs Drilldown",
			"type":          "app",
			"enabled":       true,
			"defaultNavUrl": "/plugins/grafana-lokiexplore-app/page/grafana-logs-drilldown",
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	plugin, err := client.GetPlugin(context.Background(), "grafana-lokiexplore-app")
	if err != nil {
		t.Fatalf("GetPlugin failed: %v", err)
	}
	if plugin.Name != "Grafana Logs Drilldown" {
		t.Errorf("Expected plugin name Grafana Logs Drilldown, got %s", plugin.Name)
	}
}

func TestGenerateLogsDrilldownLinkUsesPluginNavURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/plugins/grafana-lokiexplore-app/settings" {
			t.Errorf("Expected path /api/plugins/grafana-lokiexplore-app/settings, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":            "grafana-lokiexplore-app",
			"name":          "Grafana Logs Drilldown",
			"type":          "app",
			"enabled":       true,
			"defaultNavUrl": "/plugins/grafana-lokiexplore-app/page/grafana-logs-drilldown",
		})
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{URL: server.URL, APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	deeplink, err := client.GenerateLogsDrilldownLink(context.Background(), map[string]string{
		"from":          "now-1h",
		"to":            "now",
		"datasourceUid": "loki",
	})
	if err != nil {
		t.Fatalf("GenerateLogsDrilldownLink failed: %v", err)
	}

	expected := server.URL + "/plugins/grafana-lokiexplore-app/page/grafana-logs-drilldown?from=now-1h&to=now&var-datasource=loki"
	if deeplink.URL != expected {
		t.Errorf("Expected %s, got %s", expected, deeplink.URL)
	}
}

func TestRenderDashboardPanelUsesRenderEndpointOutsideAPIBase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/render/d-solo/test-dashboard" {
			t.Errorf("Expected path /render/d-solo/test-dashboard, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("panelId"); got != "7" {
			t.Errorf("Expected panelId=7, got %q", got)
		}
		if got := r.URL.Query().Get("width"); got != "1000" {
			t.Errorf("Expected width=1000, got %q", got)
		}
		if got := r.URL.Query().Get("height"); got != "500" {
			t.Errorf("Expected height=500, got %q", got)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header, got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Accept") != "image/png" {
			t.Errorf("Expected Accept image/png, got %q", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("png-data"))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:    server.URL,
		APIKey: "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	image, err := client.RenderDashboardPanel(context.Background(), "test-dashboard", 7, map[string]string{
		"width":  "1000",
		"height": "500",
	})
	if err != nil {
		t.Fatalf("RenderDashboardPanel failed: %v", err)
	}

	if image.ContentType != "image/png" {
		t.Errorf("Expected content type image/png, got %q", image.ContentType)
	}
	if string(image.ImageData) != "png-data" {
		t.Errorf("Expected image payload png-data, got %q", string(image.ImageData))
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

func TestRetryOnTransientStatusThenSuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"message":"temporary unavailable"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:            server.URL,
		APIKey:         "test-api-key",
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dashboards, err := client.GetDashboards(context.Background())
	if err != nil {
		t.Fatalf("GetDashboards failed: %v", err)
	}
	if len(dashboards) != 0 {
		t.Fatalf("expected 0 dashboards, got %d", len(dashboards))
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestNoRetryForNonIdempotentMethod(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"message":"temporary unavailable"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:            server.URL,
		APIKey:         "test-api-key",
		MaxRetries:     3,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.CreateDatasource(context.Background(), CreateDatasourceRequest{
		Name: "Test Datasource",
		Type: "prometheus",
		URL:  "http://localhost:9090",
	})
	if err == nil {
		t.Fatalf("expected error for non-idempotent request")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-idempotent request, got %d", attempts)
	}
}
