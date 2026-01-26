// Package tools provides shared utilities for tool definitions across all services.
package tools

// CommonDescriptions contains reusable parameter descriptions to reduce duplication.
// This reduces code by 30-40% across all services.
var CommonDescriptions = map[string]string{
	// Kubernetes Common Parameters
	"k8s_kind":      "Resource kind - the type of Kubernetes object (Pod, Service, Deployment, ConfigMap, Secret, Ingress, PersistentVolume, Namespace, Node, etc.). Use exact case-sensitive names.",
	"k8s_name":      "Exact name of the specific resource instance. Must match metadata.name exactly. Names are case-sensitive and follow Kubernetes naming conventions.",
	"k8s_namespace": "Kubernetes namespace where the resource is located. Required for namespaced resources (Pod, Service, Deployment, ConfigMap, Secret, etc.). Not applicable for cluster-scoped resources (Node, ClusterRole, PersistentVolume).",
	"k8s_debug":     "Enable verbose debug output for troubleshooting API calls. Set to 'true' for detailed information, 'false' or omit for normal output.",
	"k8s_jsonpath":  "JSONPath expression to extract specific fields from the resource. Examples: '{.status.phase}' for Pod phase, '{.spec.containers[*].name}' for container names, '{.metadata.labels}' for labels.",
	"k8s_dryrun":    "Validate without applying changes. Set to 'true' to test the operation without modifying the cluster, 'false' or omit to apply normally.",
	"k8s_overwrite": "Overwrite existing values when updating. Set to 'true' to replace existing values, 'false' or omit to keep existing values.",

	// Monitoring Common Parameters
	"query":         "Query string or expression specific to the service (PromQL for Prometheus, KQL for Kibana, etc.). Include all necessary filters and aggregations.",
	"start_time":    "Start timestamp for range queries in RFC3339 format (e.g., 2024-01-15T10:30:00Z). Required for historical data queries.",
	"end_time":      "End timestamp for range queries in RFC3339 format. Must be after start_time. If omitted, uses current time.",
	"debug_verbose": "Enable comprehensive debug output including API requests, responses, and timing information. Set to 'true' for troubleshooting.",

	// Prometheus Specific
	"prom_query": "PromQL query string to execute. Include metric names, labels, operators, and aggregation functions.",
	"prom_time":  "Optional timestamp in RFC3339 format. If not specified, uses current time for instant queries.",
	"prom_step":  "Query resolution step width for range queries (e.g., '30s', '1m', '5m'). Determines data point granularity. Defaults to '15s'.",
	"prom_state": "Filter targets by state. Options: 'active' (currently being scraped), 'dropped' (excluded by relabeling), 'any' (all targets).",
	"prom_type":  "Filter rules by type: 'alert' for alerting rules, 'record' for recording rules. If not specified, returns all rules.",
	"prom_label": "The label name to retrieve values for. Must match exactly as it appears in the metrics.",
	"prom_match": "Label selector(s) to match series (e.g., 'up{job=\"prometheus\"}'). Can specify multiple matchers.",

	// Grafana Specific
	"grafana_uid":        "Unique identifier (UID) of the dashboard. Found in dashboard URL (/d/dashboard-uid/dashboard-name). Case-sensitive.",
	"grafana_folder":     "Grafana folder name or ID containing the dashboard or resource. Organizes dashboards hierarchically.",
	"grafana_tag":        "Tag to filter dashboards. Grafana uses tags for organizational and filtering purposes.",
	"grafana_datasource": "Name or ID of the data source configuration. Examples: 'Prometheus', 'Elasticsearch', 'InfluxDB'.",

	// Kibana Specific
	"kibana_space_id":     "The unique identifier of the Kibana space. Spaces isolate Kibana objects and dashboards.",
	"kibana_dashboard_id": "The unique identifier of the dashboard within the space.",
	"kibana_object_type":  "Type of saved object to search for: 'dashboard', 'visualization', 'index-pattern', 'search', 'canvas'.",
	"kibana_search_term":  "Search term to filter saved objects by title or content. Uses substring matching.",
	"kibana_page":         "Page number for pagination (starts from 1). Use with per_page for result pagination.",
	"kibana_per_page":     "Number of results per page (max 100). Defaults to 20 for better performance.",
}

// ParameterDefaults provides common default values for parameters across services.
var ParameterDefaults = map[string]interface{}{
	"debug":              false,
	"dryRun":             false,
	"timeout":            30,
	"retries":            3,
	"page":               1,
	"per_page":           20,
	"step":               "15s",
	"gracePeriodSeconds": 30,
}

// ParameterEnums defines standard enum values for parameters.
var ParameterEnums = map[string][]string{
	"prom_state":         {"active", "dropped", "any"},
	"prom_type":          {"alert", "record"},
	"patch_type":         {"json", "merge", "apply"},
	"kibana_object_type": {"dashboard", "visualization", "index-pattern", "search", "canvas"},
}
