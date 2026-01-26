package handlers

import (
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
)

func TestRequireStringParam(t *testing.T) {
	tests := []struct {
		name    string
		req     mcp.CallToolRequest
		param   string
		wantVal string
		wantErr bool
	}{
		{
			name: "valid string param",
			req: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{"test": "value"},
				},
			},
			param:   "test",
			wantVal: "value",
			wantErr: false,
		},
		{
			name:    "missing required param",
			req:     mcp.CallToolRequest{},
			param:   "test",
			wantVal: "",
			wantErr: true,
		},
		{
			name: "non-string param",
			req: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{"test": 123},
				},
			},
			param:   "test",
			wantVal: "",
			wantErr: true,
		},
		{
			name: "empty string param",
			req: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{"test": ""},
				},
			},
			param:   "test",
			wantVal: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := requireStringParam(tt.req, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("requireStringParam() error = %v, wantErr %v", err, tt.wantErr)
			}
			if val != tt.wantVal {
				t.Errorf("requireStringParam() = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

func TestGetOptionalStringParam(t *testing.T) {
	tests := []struct {
		name    string
		req     mcp.CallToolRequest
		param   string
		wantVal string
	}{
		{
			name: "param exists",
			req: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{"test": "value"},
				},
			},
			param:   "test",
			wantVal: "value",
		},
		{
			name:    "param missing",
			req:     mcp.CallToolRequest{},
			param:   "test",
			wantVal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := getOptionalStringParam(tt.req, tt.param)
			if val != tt.wantVal {
				t.Errorf("getOptionalStringParam() = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

func TestGetOptionalBoolParam(t *testing.T) {
	tests := []struct {
		name    string
		req     mcp.CallToolRequest
		param   string
		wantVal bool
	}{
		{
			name: "bool param true",
			req: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{"flag": true},
				},
			},
			param:   "flag",
			wantVal: true,
		},
		{
			name:    "param missing",
			req:     mcp.CallToolRequest{},
			param:   "flag",
			wantVal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := getOptionalBoolParam(tt.req, tt.param)
			if val != tt.wantVal {
				t.Errorf("getOptionalBoolParam() = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

func TestGetOptionalIntParam(t *testing.T) {
	tests := []struct {
		name    string
		req     mcp.CallToolRequest
		param   string
		wantVal int
	}{
		{
			name: "int param",
			req: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{"count": 42.0},
				},
			},
			param:   "count",
			wantVal: 42,
		},
		{
			name:    "param missing",
			req:     mcp.CallToolRequest{},
			param:   "count",
			wantVal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := getOptionalIntParam(tt.req, tt.param)
			if val != tt.wantVal {
				t.Errorf("getOptionalIntParam() = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

// Temporarily commented out - tests require proper client mocking
// func TestHandleListReleases(t *testing.T) {
// 	handler := HandleListReleases(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing namespace parameter")
// 	}
// }
//
// func TestHandleGetRelease(t *testing.T) {
// 	handler := HandleGetRelease(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleListRepositories(t *testing.T) {
// 	handler := HandleListRepositories(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleInstallRelease(t *testing.T) {
// 	handler := HandleInstallRelease(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing chart parameter")
// 	}
// }
//
// func TestHandleUninstallRelease(t *testing.T) {
// 	handler := HandleUninstallRelease(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleUpgradeRelease(t *testing.T) {
// 	handler := HandleUpgradeRelease(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default", "name": "test"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing chart parameter")
// 	}
// }
//
// func TestHandleRollbackRelease(t *testing.T) {
// 	handler := HandleRollbackRelease(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default", "name": "test"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing revision parameter")
// 	}
// }
//
// func TestHandleGetReleaseHistory(t *testing.T) {
// 	handler := HandleGetReleaseHistory(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]any{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleGetReleaseValues(t *testing.T) {
// 	handler := HandleGetReleaseValues(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleGetReleaseManifest(t *testing.T) {
// 	handler := HandleGetReleaseManifest(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleSearchCharts(t *testing.T) {
// 	handler := HandleSearchCharts(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing query parameter")
// 	}
// }
//
// func TestHandleGetChartInfo(t *testing.T) {
// 	handler := HandleGetChartInfo(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing chart parameter")
// 	}
// }
//
// func TestHandleTemplateChart(t *testing.T) {
// 	handler := HandleTemplateChart(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing chart parameter")
// 	}
// }
//
// func TestHandleAddRepository(t *testing.T) {
// 	handler := HandleAddRepository(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleRemoveRepository(t *testing.T) {
// 	handler := HandleRemoveRepository(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleUpdateRepositories(t *testing.T) {
// 	handler := HandleUpdateRepositories(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleGetMirrorConfiguration(t *testing.T) {
// 	handler := HandleGetMirrorConfiguration(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleListReleasesPaginated(t *testing.T) {
// 	handler := HandleListReleasesPaginated(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing namespace parameter")
// 	}
// }
//
// func TestHandleGetReleaseStatus(t *testing.T) {
// 	handler := HandleGetReleaseStatus(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleGetRecentFailures(t *testing.T) {
// 	handler := HandleGetRecentFailures(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleClusterOverview(t *testing.T) {
// 	handler := HandleGetClusterOverview(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleGetReleaseSummary(t *testing.T) {
// 	handler := HandleGetReleaseSummary(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// // Temporarily commented out - missing handler function
// // func TestHandleListReleasesSummary(t *testing.T) {
// // 	handler := HandleListReleasesSummary(nil)
// //
// // 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// // 	req := mcp.CallToolRequest{}
// // 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
// //
// // 	_, err := handler(nil, req)
// // 	_ = err // Expected to fail due to nil client
// // }
//
// func TestHandleFindReleasesByLabels(t *testing.T) {
// 	handler := HandleFindReleasesByLabels(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing labels parameter")
// 	}
// }
//
// func TestHandleGetResourcesOfRelease(t *testing.T) {
// 	handler := HandleGetResourcesOfRelease(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleClearCache(t *testing.T) {
// 	handler := HandleClearCache(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleGetCacheStats(t *testing.T) {
// 	handler := HandleGetCacheStats(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleQuickInfo(t *testing.T) {
// 	handler := HandleGetQuickInfo(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleFindReleasesByChart(t *testing.T) {
// 	handler := HandleFindReleasesByChart(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing chart parameter")
// 	}
// }
//
// // Temporarily commented out - missing handler function
// // func TestHandleListReleasesInNamespace(t *testing.T) {
// // 	handler := HandleListReleasesInNamespace(nil)
// //
// // 	// Test with missing required parameters
// // 	req := mcp.CallToolRequest{}
// // 	req.Params.Arguments = map[string]interface{}{}
// //
// // 	_, err := handler(nil, req)
// // 	if err == nil {
// // 		t.Fatalf("expected error for missing namespace parameter")
// // 	}
// // }
//
// func TestHandleFindBrokenReleases(t *testing.T) {
// 	handler := HandleFindBrokenReleases(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
//
// func TestHandleValidateRelease(t *testing.T) {
// 	handler := HandleValidateRelease(nil)
//
// 	// Test with missing required parameters
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	if err == nil {
// 		t.Fatalf("expected error for missing name parameter")
// 	}
// }
//
// func TestHandleHelmHealthCheck(t *testing.T) {
// 	handler := HandleHelmHealthCheck(nil)
//
// 	// Test with valid parameters (will fail due to nil client but tests parameter validation)
// 	req := mcp.CallToolRequest{}
// 	req.Params.Arguments = map[string]interface{}{"namespace": "default"}
//
// 	_, err := handler(nil, req)
// 	_ = err // Expected to fail due to nil client
// }
