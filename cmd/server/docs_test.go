package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestToolsReferenceServiceCountsMatchInventory(t *testing.T) {
	docPath := filepath.Join("..", "..", "docs", "TOOLS.md")
	docContent, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", docPath, err)
	}

	counts := collectToolCounts(t)
	expected := map[string]int{
		"Kubernetes":    counts["kubernetes"],
		"Helm":          counts["helm"],
		"Grafana":       counts["grafana"],
		"Prometheus":    counts["prometheus"],
		"Kibana":        counts["kibana"],
		"Elasticsearch": counts["elasticsearch"],
		"Alertmanager":  counts["alertmanager"],
		"Jaeger":        counts["jaeger"],
		"OpenTelemetry": counts["opentelemetry"],
		"Utilities":     counts["utilities"],
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

func collectToolCounts(t *testing.T) map[string]int {
	t.Helper()

	files, err := filepath.Glob(filepath.Join("..", "..", "internal", "services", "*", "tools", "*.go"))
	if err != nil {
		t.Fatalf("failed to glob tool files: %v", err)
	}

	namePattern := regexp.MustCompile(`NewTool\("([a-z0-9_]+)"|Name:\s*"([a-z0-9_]+)"`)
	seen := make(map[string]map[string]bool)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("failed to read %s: %v", file, err)
		}

		matches := namePattern.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			name := match[1]
			if name == "" {
				name = match[2]
			}
			if name == "" || name == "test" {
				continue
			}

			service := strings.SplitN(name, "_", 2)[0]
			if seen[service] == nil {
				seen[service] = make(map[string]bool)
			}
			seen[service][name] = true
		}
	}

	counts := make(map[string]int, len(seen))
	for service, names := range seen {
		counts[service] = len(names)
	}
	return counts
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
