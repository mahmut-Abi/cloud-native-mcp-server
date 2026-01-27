package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ClientOptions
		wantErr bool
	}{
		{
			name: "valid client",
			opts: &ClientOptions{
				URL:     "http://localhost:5601",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			opts: &ClientOptions{
				URL:     "",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid URL",
			opts: &ClientOptions{
				URL:     "://invalid-url",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "with API key",
			opts: &ClientOptions{
				URL:     "http://localhost:5601",
				APIKey:  "test-api-key",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "with basic auth",
			opts: &ClientOptions{
				URL:      "http://localhost:5601",
				Username: "user",
				Password: "pass",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "with space",
			opts: &ClientOptions{
				URL:     "http://localhost:5601",
				Space:   "my-space",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "default timeout",
			opts: &ClientOptions{
				URL:     "http://localhost:5601",
				Timeout: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() should return non-nil client")
			}
		})
	}
}

func TestGetSpaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		URL:     server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	spaces, err := client.GetSpaces(ctx)
	if err != nil {
		t.Errorf("GetSpaces() error = %v", err)
		return
	}

	if spaces == nil {
		t.Error("GetSpaces() should return non-nil spaces")
	}
}

func TestGetIndexPatterns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page": 1,
			"per_page": 20,
			"total": 1,
			"saved_objects": [{
				"id": "test-pattern",
				"type": "index-pattern",
				"attributes": {
					"title": "test-*",
					"timeFieldName": "@timestamp"
				}
			}]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		URL:     server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	patterns, err := client.GetIndexPatterns(ctx)
	if err != nil {
		t.Errorf("GetIndexPatterns() error = %v", err)
		return
	}

	if patterns == nil {
		t.Error("GetIndexPatterns() should return non-nil patterns")
	}
}

func TestGetDashboards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page": 1,
			"per_page": 20,
			"total": 1,
			"saved_objects": [{
				"id": "test-dashboard",
				"type": "dashboard",
				"attributes": {
					"title": "Test Dashboard"
				}
			}]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		URL:     server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	dashboards, err := client.GetDashboards(ctx)
	if err != nil {
		t.Errorf("GetDashboards() error = %v", err)
		return
	}

	if dashboards == nil {
		t.Error("GetDashboards() should return non-nil dashboards")
	}
}

func TestGetVisualizations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page": 1,
			"per_page": 20,
			"total": 1,
			"saved_objects": [{
				"id": "test-visualization",
				"type": "visualization",
				"attributes": {
					"title": "Test Visualization"
				}
			}]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		URL:     server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	visualizations, err := client.GetVisualizations(ctx)
	if err != nil {
		t.Errorf("GetVisualizations() error = %v", err)
		return
	}

	if visualizations == nil {
		t.Error("GetVisualizations() should return non-nil visualizations")
	}
}

func TestSearchSavedObjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page": 1,
			"per_page": 20,
			"total": 0,
			"saved_objects": []
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		URL:     server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	result, err := client.SearchSavedObjects(ctx, "dashboard", "", 1, 20)
	if err != nil {
		t.Errorf("SearchSavedObjects() error = %v", err)
		return
	}

	if result == nil {
		t.Error("SearchSavedObjects() should return non-nil result")
	}
}

func TestCreateSavedObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "new-id",
			"type": "dashboard",
			"attributes": {
				"title": "New Dashboard"
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		URL:     server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	attrs := map[string]interface{}{
		"title": "New Dashboard",
	}
	obj, err := client.CreateSavedObject(ctx, "dashboard", attrs, nil, "")
	if err != nil {
		t.Errorf("CreateSavedObject() error = %v", err)
		return
	}

	if obj == nil {
		t.Error("CreateSavedObject() should return non-nil object")
	}
}

func TestDeleteSavedObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "deleted-id",
			"type": "dashboard"
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		URL:     server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	err := client.DeleteSavedObject(ctx, "dashboard", "test-id", false)
	if err != nil {
		t.Errorf("DeleteSavedObject() error = %v", err)
	}
}