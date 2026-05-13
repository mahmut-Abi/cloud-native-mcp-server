package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/opentelemetry/client"
)

type collectorModel struct {
	Receivers  map[string]map[string]interface{}
	Processors map[string]map[string]interface{}
	Exporters  map[string]map[string]interface{}
	Connectors map[string]map[string]interface{}
	Extensions map[string]map[string]interface{}
	Pipelines  []pipelineModel
	Telemetry  map[string]interface{}
	Warnings   []string
}

type pipelineModel struct {
	Name       string
	Signal     string
	Receivers  []string
	Processors []string
	Exporters  []string
}

type finding struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// HandleGetConfigSummary summarizes collector config into pipeline and component views.
func HandleGetConfigSummary(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		config, err := c.GetConfig(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve collector configuration: %v", err)), nil
		}

		result, err := json.MarshalIndent(summarizeCollectorConfig(buildCollectorModel(config)), "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize config summary: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}
}

// HandleGetCollectorSummary combines health, status, and config into one compact overview.
func HandleGetCollectorSummary(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		summary := map[string]interface{}{
			"errors": []string{},
		}

		health, err := c.GetHealth(ctx)
		if err != nil {
			summary["errors"] = append(summary["errors"].([]string), fmt.Sprintf("health: %v", err))
		} else {
			summary["health"] = summarizeHealth(health)
		}

		status, err := c.GetStatus(ctx)
		if err != nil {
			summary["errors"] = append(summary["errors"].([]string), fmt.Sprintf("status: %v", err))
		} else {
			summary["status"] = summarizeStatus(status)
		}

		config, err := c.GetConfig(ctx)
		if err != nil {
			summary["errors"] = append(summary["errors"].([]string), fmt.Sprintf("config: %v", err))
		} else {
			model := buildCollectorModel(config)
			configSummary := summarizeCollectorConfig(model)
			summary["config"] = map[string]interface{}{
				"components": configSummary["components"],
				"pipelines":  configSummary["pipelines"],
				"warnings":   configSummary["warnings"],
			}
			summary["signals"] = summarizeSignals(model.Pipelines)
		}

		result, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize collector summary: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}
}

// HandleAnalyzePipelineStatus analyzes collector pipelines for broken refs and common misconfigurations.
func HandleAnalyzePipelineStatus(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		signalFilter, _ := args["signal"].(string)
		pipelineFilter, _ := args["pipeline"].(string)

		config, err := c.GetConfig(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve collector configuration: %v", err)), nil
		}

		var health map[string]interface{}
		health, _ = c.GetHealth(ctx)

		var status map[string]interface{}
		status, _ = c.GetStatus(ctx)

		model := buildCollectorModel(config)
		analysis := analyzePipelines(model, health, status, signalFilter, pipelineFilter)

		result, err := json.MarshalIndent(analysis, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize pipeline analysis: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	}
}

func buildCollectorModel(config map[string]interface{}) collectorModel {
	model := collectorModel{
		Receivers:  namedConfigSection(config["receivers"]),
		Processors: namedConfigSection(config["processors"]),
		Exporters:  namedConfigSection(config["exporters"]),
		Connectors: namedConfigSection(config["connectors"]),
		Extensions: namedConfigSection(config["extensions"]),
	}

	service := asMap(config["service"])
	if len(service) == 0 {
		model.Warnings = append(model.Warnings, "No service section found in collector config.")
		return model
	}

	model.Telemetry = asMap(service["telemetry"])
	pipelines := asMap(service["pipelines"])
	if len(pipelines) == 0 {
		model.Warnings = append(model.Warnings, "No service.pipelines configured.")
		return model
	}

	names := sortedKeysFromMap(pipelines)
	for _, name := range names {
		pipelineConfig := asMap(pipelines[name])
		model.Pipelines = append(model.Pipelines, pipelineModel{
			Name:       name,
			Signal:     pipelineSignal(name),
			Receivers:  stringSlice(pipelineConfig["receivers"]),
			Processors: stringSlice(pipelineConfig["processors"]),
			Exporters:  stringSlice(pipelineConfig["exporters"]),
		})
	}

	if len(model.Receivers) == 0 {
		model.Warnings = append(model.Warnings, "No receivers configured.")
	}
	if len(model.Exporters) == 0 && len(model.Connectors) == 0 {
		model.Warnings = append(model.Warnings, "No exporters or connectors configured.")
	}

	return model
}

func summarizeCollectorConfig(model collectorModel) map[string]interface{} {
	componentSummary := map[string]interface{}{
		"receivers": map[string]interface{}{
			"count": len(model.Receivers),
			"names": sortedKeysFromNamedSection(model.Receivers),
		},
		"processors": map[string]interface{}{
			"count": len(model.Processors),
			"names": sortedKeysFromNamedSection(model.Processors),
		},
		"exporters": map[string]interface{}{
			"count": len(model.Exporters),
			"names": sortedKeysFromNamedSection(model.Exporters),
		},
		"connectors": map[string]interface{}{
			"count": len(model.Connectors),
			"names": sortedKeysFromNamedSection(model.Connectors),
		},
		"extensions": map[string]interface{}{
			"count": len(model.Extensions),
			"names": sortedKeysFromNamedSection(model.Extensions),
		},
	}

	pipelines := make([]map[string]interface{}, 0, len(model.Pipelines))
	for _, pipeline := range model.Pipelines {
		pipelines = append(pipelines, map[string]interface{}{
			"name":       pipeline.Name,
			"signal":     pipeline.Signal,
			"receivers":  pipeline.Receivers,
			"processors": pipeline.Processors,
			"exporters":  pipeline.Exporters,
		})
	}

	return map[string]interface{}{
		"components": componentSummary,
		"pipelines":  pipelines,
		"signals":    summarizeSignals(model.Pipelines),
		"telemetry":  summarizeTelemetry(model.Telemetry),
		"warnings":   model.Warnings,
	}
}

func analyzePipelines(model collectorModel, health, status map[string]interface{}, signalFilter, pipelineFilter string) map[string]interface{} {
	filtered := make([]pipelineModel, 0, len(model.Pipelines))
	for _, pipeline := range model.Pipelines {
		if signalFilter != "" && !strings.EqualFold(signalFilter, pipeline.Signal) {
			continue
		}
		if pipelineFilter != "" && pipelineFilter != pipeline.Name {
			continue
		}
		filtered = append(filtered, pipeline)
	}

	validReceiverRefs := referenceSet(model.Receivers, model.Connectors)
	validExporterRefs := referenceSet(model.Exporters, model.Connectors)
	validProcessorRefs := referenceSet(model.Processors, nil)
	statusPipelines := asMap(status["pipelines"])

	results := make([]map[string]interface{}, 0, len(filtered))
	for _, pipeline := range filtered {
		findings := make([]finding, 0, 8)

		missingReceivers := missingRefs(pipeline.Receivers, validReceiverRefs)
		missingProcessors := missingRefs(pipeline.Processors, validProcessorRefs)
		missingExporters := missingRefs(pipeline.Exporters, validExporterRefs)

		if len(pipeline.Receivers) == 0 {
			findings = append(findings, finding{"error", "Pipeline has no receivers configured."})
		}
		if len(pipeline.Exporters) == 0 {
			findings = append(findings, finding{"error", "Pipeline has no exporters configured."})
		}
		if len(pipeline.Processors) == 0 {
			findings = append(findings, finding{"warning", "Pipeline has no processors configured; verify batching, memory limiting, and enrichment are handled elsewhere."})
		}
		if len(missingReceivers) > 0 {
			findings = append(findings, finding{"error", fmt.Sprintf("Receivers referenced but not defined: %s", strings.Join(missingReceivers, ", "))})
		}
		if len(missingProcessors) > 0 {
			findings = append(findings, finding{"error", fmt.Sprintf("Processors referenced but not defined: %s", strings.Join(missingProcessors, ", "))})
		}
		if len(missingExporters) > 0 {
			findings = append(findings, finding{"error", fmt.Sprintf("Exporters/connectors referenced but not defined: %s", strings.Join(missingExporters, ", "))})
		}
		if !containsProcessor(pipeline.Processors, "batch") {
			findings = append(findings, finding{"warning", "No batch processor configured in pipeline; throughput and exporter efficiency may suffer."})
		}
		if !containsProcessor(pipeline.Processors, "memory_limiter") {
			findings = append(findings, finding{"warning", "No memory_limiter processor configured in pipeline; collector resource pressure may be harder to control."})
		}
		if pipeline.Signal == "traces" && !containsAnyProcessor(pipeline.Processors, []string{"tail_sampling", "probabilistic_sampler"}) {
			findings = append(findings, finding{"info", "No collector-level sampling processor configured; if you expect sampling, verify it is implemented upstream or in another collector."})
		}

		entry := map[string]interface{}{
			"name":       pipeline.Name,
			"signal":     pipeline.Signal,
			"receivers":  pipeline.Receivers,
			"processors": pipeline.Processors,
			"exporters":  pipeline.Exporters,
			"findings":   findings,
		}
		if pipelineStatus := asMap(statusPipelines[pipeline.Name]); len(pipelineStatus) > 0 {
			entry["status"] = pipelineStatus
		}
		results = append(results, entry)
	}

	globalFindings := make([]finding, 0, 4)
	if len(model.Pipelines) == 0 {
		globalFindings = append(globalFindings, finding{"error", "Collector config has no pipelines."})
	}
	if len(model.Warnings) > 0 {
		for _, warning := range model.Warnings {
			globalFindings = append(globalFindings, finding{"warning", warning})
		}
	}
	if healthStatus := firstStringValue(health, "status", "state"); healthStatus != "" && !strings.EqualFold(healthStatus, "ok") && !strings.EqualFold(healthStatus, "healthy") {
		globalFindings = append(globalFindings, finding{"warning", fmt.Sprintf("Collector health endpoint reported %q.", healthStatus)})
	}

	return map[string]interface{}{
		"health":          summarizeHealth(health),
		"status":          summarizeStatus(status),
		"pipelineCount":   len(results),
		"signals":         summarizeSignals(filtered),
		"pipelines":       results,
		"globalFindings":  globalFindings,
		"availableFilter": map[string]interface{}{"signal": signalFilter, "pipeline": pipelineFilter},
	}
}

func summarizeHealth(health map[string]interface{}) map[string]interface{} {
	if len(health) == 0 {
		return map[string]interface{}{}
	}
	summary := map[string]interface{}{
		"keys": sortedKeysFromMap(health),
	}
	if status := firstStringValue(health, "status", "state"); status != "" {
		summary["status"] = status
	}
	if message := firstStringValue(health, "message"); message != "" {
		summary["message"] = message
	}
	return summary
}

func summarizeStatus(status map[string]interface{}) map[string]interface{} {
	if len(status) == 0 {
		return map[string]interface{}{}
	}
	summary := map[string]interface{}{
		"keys": sortedKeysFromMap(status),
	}
	if state := firstStringValue(status, "status", "state", "health"); state != "" {
		summary["state"] = state
	}
	if message := firstStringValue(status, "message"); message != "" {
		summary["message"] = message
	}
	if pipelines := asMap(status["pipelines"]); len(pipelines) > 0 {
		summary["pipelineCount"] = len(pipelines)
		summary["pipelines"] = sortedKeysFromMap(pipelines)
	}
	if components := asMap(status["components"]); len(components) > 0 {
		summary["componentGroups"] = sortedKeysFromMap(components)
	}
	return summary
}

func summarizeSignals(pipelines []pipelineModel) map[string]int {
	signals := map[string]int{
		"metrics": 0,
		"traces":  0,
		"logs":    0,
	}
	for _, pipeline := range pipelines {
		signals[pipeline.Signal]++
	}
	return signals
}

func summarizeTelemetry(telemetry map[string]interface{}) map[string]interface{} {
	if len(telemetry) == 0 {
		return map[string]interface{}{}
	}
	summary := map[string]interface{}{}
	if logs := asMap(telemetry["logs"]); len(logs) > 0 {
		logSummary := map[string]interface{}{}
		if level := firstStringValue(logs, "level"); level != "" {
			logSummary["level"] = level
		}
		if len(logSummary) > 0 {
			summary["logs"] = logSummary
		}
	}
	if metrics := asMap(telemetry["metrics"]); len(metrics) > 0 {
		metricSummary := map[string]interface{}{}
		if address := firstStringValue(metrics, "address"); address != "" {
			metricSummary["address"] = address
		}
		if readers, ok := metrics["readers"].([]interface{}); ok {
			metricSummary["readerCount"] = len(readers)
		}
		if len(metricSummary) > 0 {
			summary["metrics"] = metricSummary
		}
	}
	if traces := asMap(telemetry["traces"]); len(traces) > 0 {
		traceSummary := map[string]interface{}{}
		if level := firstStringValue(traces, "level"); level != "" {
			traceSummary["level"] = level
		}
		if len(traceSummary) > 0 {
			summary["traces"] = traceSummary
		}
	}
	return summary
}

func namedConfigSection(value interface{}) map[string]map[string]interface{} {
	section := asMap(value)
	result := make(map[string]map[string]interface{}, len(section))
	for name, raw := range section {
		result[name] = asMap(raw)
	}
	return result
}

func asMap(value interface{}) map[string]interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return typed
	default:
		return map[string]interface{}{}
	}
}

func stringSlice(value interface{}) []string {
	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...)
	case []interface{}:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func pipelineSignal(name string) string {
	prefix := strings.SplitN(name, "/", 2)[0]
	switch prefix {
	case "metrics", "traces", "logs":
		return prefix
	default:
		return name
	}
}

func sortedKeysFromMap(input map[string]interface{}) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysFromNamedSection(input map[string]map[string]interface{}) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func referenceSet(primary, extra map[string]map[string]interface{}) map[string]struct{} {
	result := make(map[string]struct{}, len(primary)+len(extra))
	for key := range primary {
		result[key] = struct{}{}
	}
	for key := range extra {
		result[key] = struct{}{}
	}
	return result
}

func missingRefs(refs []string, defined map[string]struct{}) []string {
	var missing []string
	for _, ref := range refs {
		if _, ok := defined[ref]; !ok {
			missing = append(missing, ref)
		}
	}
	return missing
}

func containsProcessor(processors []string, name string) bool {
	for _, processor := range processors {
		if processor == name || strings.HasPrefix(processor, name+"/") {
			return true
		}
	}
	return false
}

func containsAnyProcessor(processors []string, names []string) bool {
	for _, name := range names {
		if containsProcessor(processors, name) {
			return true
		}
	}
	return false
}

func firstStringValue(values map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key].(string); ok && value != "" {
			return value
		}
	}
	return ""
}
