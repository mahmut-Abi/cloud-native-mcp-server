package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
		"OpenTelemetry": len(counts["opentelemetry"]),
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
