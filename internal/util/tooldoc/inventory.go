package tooldoc

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	BeginMarker = "<!-- BEGIN GENERATED TOOL INVENTORY -->"
	EndMarker   = "<!-- END GENERATED TOOL INVENTORY -->"
)

var (
	newToolPattern = regexp.MustCompile(`NewTool\("([a-z0-9_]+)"`)
	namePattern    = regexp.MustCompile(`Name:\s*"([a-z0-9_]+)"`)
)

var serviceDisplayNames = map[string]string{
	"alertmanager":  "Alertmanager",
	"elasticsearch": "Elasticsearch",
	"grafana":       "Grafana",
	"helm":          "Helm",
	"jaeger":        "Jaeger",
	"kibana":        "Kibana",
	"kubernetes":    "Kubernetes",
	"opentelemetry": "OpenTelemetry",
	"prometheus":    "Prometheus",
	"utilities":     "Utilities",
}

var serviceOrder = []string{
	"kubernetes",
	"helm",
	"grafana",
	"prometheus",
	"kibana",
	"elasticsearch",
	"alertmanager",
	"jaeger",
	"opentelemetry",
	"utilities",
}

// CollectInventory scans tool definition files and returns service -> sorted tool names.
func CollectInventory(repoRoot string) (map[string][]string, error) {
	files, err := filepath.Glob(filepath.Join(repoRoot, "internal", "services", "*", "tools", "*.go"))
	if err != nil {
		return nil, fmt.Errorf("glob tool files: %w", err)
	}

	seen := make(map[string]map[string]bool)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file, err)
		}
		text := string(content)

		for _, match := range newToolPattern.FindAllStringSubmatch(text, -1) {
			recordTool(seen, match[1])
		}
		for _, match := range namePattern.FindAllStringSubmatch(text, -1) {
			recordTool(seen, match[1])
		}
	}

	inventory := make(map[string][]string, len(seen))
	for service, names := range seen {
		tools := make([]string, 0, len(names))
		for name := range names {
			tools = append(tools, name)
		}
		sort.Strings(tools)
		inventory[service] = tools
	}

	return inventory, nil
}

// RenderGeneratedInventory renders the generated inventory markdown block.
func RenderGeneratedInventory(inventory map[string][]string) string {
	var b strings.Builder
	b.WriteString(BeginMarker)
	b.WriteString("\n")
	b.WriteString("## Generated Inventory\n\n")
	b.WriteString("This section is generated from `internal/services/**/tools/*.go`.\n")
	b.WriteString("Do not edit this block by hand.\n\n")

	for _, service := range serviceOrder {
		tools := inventory[service]
		if len(tools) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("### %s (%d tools)\n\n", serviceDisplayNames[service], len(tools)))
		for _, tool := range tools {
			b.WriteString(fmt.Sprintf("- `%s`\n", tool))
		}
		b.WriteString("\n")
	}

	b.WriteString(EndMarker)
	b.WriteString("\n")
	return b.String()
}

// ReplaceGeneratedInventory replaces the generated inventory block in a document.
func ReplaceGeneratedInventory(doc string, generated string) string {
	start := strings.Index(doc, BeginMarker)
	end := strings.Index(doc, EndMarker)
	if start == -1 || end == -1 || end < start {
		if strings.HasSuffix(doc, "\n") {
			return doc + "\n" + generated
		}
		return doc + "\n\n" + generated
	}

	end += len(EndMarker)
	for end < len(doc) && (doc[end] == '\n' || doc[end] == '\r') {
		end++
	}
	return doc[:start] + generated + doc[end:]
}

func recordTool(seen map[string]map[string]bool, name string) {
	if name == "" || name == "test" {
		return
	}

	service, _, ok := strings.Cut(name, "_")
	if !ok {
		return
	}

	if seen[service] == nil {
		seen[service] = make(map[string]bool)
	}
	seen[service][name] = true
}
