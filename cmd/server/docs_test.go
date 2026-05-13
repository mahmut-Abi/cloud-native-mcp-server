package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/tooldoc"
)

func TestToolsReferenceServiceCountsMatchInventory(t *testing.T) {
	docPath := filepath.Join("..", "..", "docs", "TOOLS.md")
	docContent, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", docPath, err)
	}

	counts, err := tooldoc.CollectInventory(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to collect inventory: %v", err)
	}
	expected := map[string]int{
		"Kubernetes":    len(counts["kubernetes"]),
		"Helm":          len(counts["helm"]),
		"Grafana":       len(counts["grafana"]),
		"Prometheus":    len(counts["prometheus"]),
		"Loki":          len(counts["loki"]),
		"Kibana":        len(counts["kibana"]),
		"Elasticsearch": len(counts["elasticsearch"]),
		"Alertmanager":  len(counts["alertmanager"]),
		"Jaeger":        len(counts["jaeger"]),
		"Langfuse":      len(counts["langfuse"]),
		"OpenTelemetry": len(counts["opentelemetry"]),
		"Sentry":        len(counts["sentry"]),
		"Utilities":     len(counts["utilities"]),
	}

	docCounts := parseDocServiceCounts(t, string(docContent))
	for section, want := range expected {
		got, ok := docCounts[section]
		if !ok {
			t.Fatalf("section %q not found in docs/TOOLS.md", section)
		}
		if got != want {
			t.Fatalf("section %q count mismatch: docs=%d code=%d", section, got, want)
		}
	}
}

func parseDocServiceCounts(t *testing.T, content string) map[string]int {
	t.Helper()

	pattern := regexp.MustCompile(`## ([A-Za-z]+) \((\d+) tools\)`)
	matches := pattern.FindAllStringSubmatch(content, -1)
	counts := make(map[string]int, len(matches))
	for _, match := range matches {
		value, err := strconv.Atoi(match[2])
		if err != nil {
			t.Fatalf("failed to parse service count %q: %v", match[2], err)
		}
		counts[match[1]] = value
	}
	return counts
}

func TestGeneratedInventoryBlockMatchesGenerator(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	docPath := filepath.Join(repoRoot, "docs", "TOOLS.md")

	docContent, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", docPath, err)
	}

	inventory, err := tooldoc.CollectInventory(repoRoot)
	if err != nil {
		t.Fatalf("failed to collect inventory: %v", err)
	}

	generated := tooldoc.RenderGeneratedInventory(inventory)
	updated := tooldoc.ReplaceGeneratedInventory(string(docContent), generated)
	if updated != string(docContent) {
		t.Fatalf("generated inventory block is out of date; run `go run ./cmd/toolsdocgen`")
	}
}

func TestLLMVisibleDocsReferenceKnownTools(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	inventory, err := tooldoc.CollectInventory(repoRoot)
	if err != nil {
		t.Fatalf("failed to collect inventory: %v", err)
	}

	known := make(map[string]struct{})
	for _, names := range inventory {
		for _, name := range names {
			known[name] = struct{}{}
		}
	}

	files := []string{
		filepath.Join(repoRoot, "docs", "TOOLS.md"),
		filepath.Join(repoRoot, "website", "content", "en", "docs", "tools.md"),
		filepath.Join(repoRoot, "website", "content", "zh", "docs", "tools.md"),
		filepath.Join(repoRoot, "website", "content", "en", "tools.md"),
		filepath.Join(repoRoot, "website", "content", "zh", "tools.md"),
		filepath.Join(repoRoot, "website", "content", "en", "services", "opentelemetry.md"),
		filepath.Join(repoRoot, "website", "content", "zh", "services", "opentelemetry.md"),
	}

	toolGuideFiles, err := filepath.Glob(filepath.Join(repoRoot, "website", "content", "zh", "tools", "*.md"))
	if err != nil {
		t.Fatalf("failed to glob zh tools guides: %v", err)
	}
	files = append(files, toolGuideFiles...)

	toolPattern := regexp.MustCompile(`\b(?:kubernetes|helm|grafana|prometheus|loki|kibana|elasticsearch|alertmanager|jaeger|langfuse|sentry|opentelemetry|utilities)_[a-z0-9_]+\b`)

	for _, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", path, err)
		}

		seen := make(map[string]struct{})
		var unknown []string
		for _, match := range toolPattern.FindAllString(string(content), -1) {
			if _, ok := seen[match]; ok {
				continue
			}
			seen[match] = struct{}{}
			if _, ok := known[match]; !ok {
				unknown = append(unknown, match)
			}
		}

		if len(unknown) > 0 {
			t.Fatalf("%s references unknown tool names: %s", path, strings.Join(unknown, ", "))
		}
	}
}
