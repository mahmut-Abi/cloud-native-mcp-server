package openapi

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services"
)

func TestNewGenerator(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	if gen == nil {
		t.Error("NewGenerator should return non-nil generator")
	}
}

func TestGenerateOpenAPISpec(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	spec, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	if spec == nil {
		t.Error("Generate() should return non-nil spec")
	}

	// Check required fields
	if spec.OpenAPI == "" {
		t.Error("OpenAPI version should not be empty")
	}

	if spec.Info.Title == "" {
		t.Error("Info.Title should not be empty")
	}

	if spec.Paths == nil {
		t.Error("Paths should not be nil")
	}

	if spec.Components.Schemas == nil {
		t.Error("Components.Schemas should not be nil")
	}

	if spec.Tags == nil {
		t.Error("Tags should not be nil")
	}
}

func TestGenerateOpenAPISpecInfo(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	spec, _ := gen.Generate()
	if spec == nil {
		t.Fatal("Generate() returned nil spec")
	}

	if spec.Info.Title == "" {
		t.Error("Info.Title should not be empty")
	}

	if spec.Info.Version == "" {
		t.Error("Info.Version should not be empty")
	}

	if spec.Info.Description == "" {
		t.Error("Info.Description should not be empty")
	}
}

func TestGenerateOpenAPISpecServers(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	spec, _ := gen.Generate()
	if spec == nil {
		return
	}

	// Servers can be empty, that's ok
	if spec.Servers != nil && len(spec.Servers) > 0 {
		for _, server := range spec.Servers {
			if server.URL == "" {
				t.Error("Server URL should not be empty")
			}
		}
	}
}

func TestGenerateOpenAPISpecComponents(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	spec, _ := gen.Generate()
	if spec == nil {
		t.Fatal("Generate() returned nil spec")
	}

	// Check Components fields
	if spec.Components.Schemas == nil {
		t.Error("Components.Schemas should not be nil")
	}

	if spec.Components.SecuritySchemes == nil {
		t.Error("Components.SecuritySchemes should not be nil")
	}
}

func TestGenerateOpenAPISpecTags(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	spec, _ := gen.Generate()
	if spec == nil {
		return
	}

	// Tags should not be empty
	if spec.Tags == nil {
		t.Error("Tags should not be nil")
	}

	if len(spec.Tags) == 0 {
		t.Error("Tags should have at least one tag")
	}

	// Check tag structure
	for _, tag := range spec.Tags {
		if tag.Name == "" {
			t.Error("Tag name should not be empty")
		}
	}
}

func TestGenerateOpenAPISpecPaths(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	spec, _ := gen.Generate()
	if spec == nil {
		return
	}

	// Paths should not be empty
	if spec.Paths == nil {
		t.Error("Paths should not be nil")
	}

	if len(spec.Paths) == 0 {
		t.Error("Paths should have at least one path")
	}
}

func TestGenerateOpenAPISpecSecurity(t *testing.T) {
	registry := services.NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	gen := NewGenerator(registry)
	spec, _ := gen.Generate()
	if spec == nil {
		return
	}

	// Security can be nil, that's ok
	if spec.Security != nil && len(spec.Security) > 0 {
		for _, scheme := range spec.Security {
			if len(scheme) == 0 {
				t.Error("Security scheme should have at least one scheme")
			}
		}
	}
}