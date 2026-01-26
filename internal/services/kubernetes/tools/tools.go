package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// GetResourceSummaryTool retrieves essential summary information about a specific Kubernetes resource
func GetResourceSummaryTool() mcp.Tool {
	logrus.Debug("Creating GetResourceSummaryTool")
	return mcp.NewTool("kubernetes_get_resource_summary",
		mcp.WithDescription("⚠️ PRIORITY: Optimized for LLM efficiency: Returns only essential fields (name, namespace, kind, status, age, labels). 90-95% smaller than full resource. ⚡ Best for: single resource status check, quick confirmation, resource identification. 📋 Use cases: pod_troubleshooting, health_check, log_analysis. 🔄 Workflow: list_resources_summary to discover problematic resources → use this tool to get detailed snapshot → combine with get_events and get_pod_logs for deep analysis when needed."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Resource kind/type - must be the same as in standard kubectl (Pod, Service, Deployment, ConfigMap, Secret, Namespace, Node, etc.). Use exact case-sensitive names as they appear in Kubernetes API.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the specific resource to get summary for. Must match metadata.name exactly. Use list_resources_summary first to find correct resource names if uncertain.")),
		mcp.WithString("namespace",
			mcp.Description("Required for namespaced resources (Pod, Service, Deployment, etc.). Omit for cluster-scoped resources (Node, PersistentVolume, ClusterRole, etc.")),
		mcp.WithString("includeLabels",
			mcp.Description("Optional comma-separated label keys to include (e.g., 'app,version,env'). If omitted, includes up to 10 labels automatically. Use to focus on specific labels important for your use case.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the tool itself (true/false).")),
	)
}

// GetResourceTool retrieves detailed information about a specific Kubernetes resource
func GetResourceTool() mcp.Tool {
	logrus.Debug("Creating GetResourceTool")
	return mcp.NewTool("kubernetes_get_resource",
		mcp.WithDescription("Retrieve detailed metadata and status of a Kubernetes resource. This tool returns the complete YAML/JSON representation of the resource, similar to 'kubectl get <resource> <name> -o yaml'. Use this when you need to inspect the current state, configuration, or status of a specific Kubernetes object. For troubleshooting, pay attention to the status fields which contain runtime information about the resource."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Resource kind - the type of Kubernetes object to retrieve. Common examples: Pod (for containers), Service (for networking), Deployment (for application workloads), ConfigMap (for configuration), Secret (for sensitive data), Ingress (for HTTP routing), PersistentVolume/PersistentVolumeClaim (for storage), Namespace (for resource grouping). Use exact case-sensitive names as they appear in Kubernetes API.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the specific resource instance to retrieve. This must match the metadata.name field of the resource exactly. Names are case-sensitive and must follow Kubernetes naming conventions (lowercase alphanumeric with hyphens and dots allowed).")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace where the resource is located. Required for namespaced resources (Pod, Service, Deployment, ConfigMap, Secret, etc.). Not applicable for cluster-scoped resources (Node, ClusterRole, PersistentVolume, etc.). If unsure about the namespace, use list_resources tool first to find the resource location. Default namespace is 'default' if not specified for namespaced resources.")),
		mcp.WithString("jsonpath",
			mcp.Description("JSONPath expression to extract specific fields from the resource instead of returning the full object. Useful for getting specific values like status, labels, or nested properties. Examples: '{.status.phase}' (for Pod phase), '{.spec.containers[*].name}' (for container names), '{.metadata.labels}' (for labels), '{.status.conditions[?(@.type==\"Ready\")].status}' (for readiness status). Enclose the expression in single quotes and use curly braces.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the API call itself. Set to 'true' to see detailed request/response information, 'false' or omit for normal output. Only use when debugging tool execution issues.")),
	)
}

// DescribeResourceTool retrieves detailed description of a specific Kubernetes resource (similar to kubectl describe)
func DescribeResourceTool() mcp.Tool {
	logrus.Debug("Creating DescribeResourceTool")
	return mcp.NewTool("kubernetes_describe_resource",
		mcp.WithDescription("Get comprehensive human-readable description and status information of a Kubernetes resource (similar to 'kubectl describe'). This tool provides formatted output with events, conditions, and detailed status information that's especially useful for troubleshooting. Unlike get_resource which returns raw YAML/JSON, this provides interpreted information including related events, resource relationships, and computed status. Use this tool when you need to understand what's happening with a resource, diagnose issues, or get a complete operational overview. The output includes resource configuration, current status, recent events, and any error conditions."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Resource kind - the type of Kubernetes object to describe. Common examples: Pod (for containers and their status), Service (for networking configuration), Deployment (for application workload status), ConfigMap (for configuration data), Secret (for sensitive data), Ingress (for HTTP routing rules), PersistentVolume/PersistentVolumeClaim (for storage), Node (for cluster infrastructure), Namespace (for resource grouping). Use exact case-sensitive names as they appear in Kubernetes API. The describe command works best with workload resources like Pods, Deployments, and Services.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the resource instance to describe. This must match the resource name exactly as it appears in Kubernetes. Use 'list_resources' tool first if you need to find the correct resource name. For Pods created by Deployments, the name will include generated suffixes.")),
		mcp.WithString("namespace",
			mcp.Description("Namespace where the resource exists. This is required for namespaced resources (Pod, Service, Deployment, ConfigMap, Secret, etc.) but should be omitted for cluster-scoped resources (Node, PersistentVolume, ClusterRole, etc.). If unsure whether a resource is namespaced, try without namespace first - the error will indicate if namespace is required. Use 'default' namespace if not specified during resource creation.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the tool itself. Set to 'true' to see detailed execution information, 'false' or omit for normal output. Only use when the tool itself is not working as expected.")),
	)
}

// GetRecentEventsTool retrieves recent cluster events with optimized output
func GetRecentEventsTool() mcp.Tool {
	logrus.Debug("Creating GetRecentEventsTool")
	return mcp.NewTool("kubernetes_get_recent_events",
		mcp.WithDescription("⚠️ PRIORITY: Optimized for LLM efficiency: Returns only recent critical events (warnings, errors, failed pods) with essential information. 80-90% smaller than full events. 🎯 Best for: problem diagnosis, cluster health monitoring, quick troubleshooting. 📋 Use cases: pod_troubleshooting, health_check, event_analysis, resource_monitoring. 🔄 Workflow: use this tool first to identify critical events → combine with list_resources_summary to view related resources → use get_events_detail for deep analysis."),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to filter events from. If not specified, shows critical events from all namespaces (requires cluster-wide permissions). For focused troubleshooting, specify the namespace where the problematic resources are located.")),
		mcp.WithString("fieldSelector",
			mcp.Description("Field selector to filter events. Common examples: 'involvedObject.name=my-pod' (events for specific pod), 'type=Warning' (only warnings), 'reason=Failed' (failure events). Combine with commas for AND logic.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of recent events to return (default: 20, max: 100). This tool is optimized for recent events, so higher limits are not recommended.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the tool itself (true/false).")),
	)
}

// GetEventsTool retrieves cluster events for troubleshooting
func GetEventsTool() mcp.Tool {
	logrus.Debug("Creating GetEventsTool")
	return mcp.NewTool("kubernetes_get_events",
		mcp.WithDescription("Retrieve Kubernetes cluster events for troubleshooting and monitoring purposes. Events provide valuable insights into cluster activities, resource state changes, and error conditions. This tool is essential for diagnosing issues with pods, deployments, services, and other resources. Events show chronological activities like pod scheduling, image pulling, container creation, failures, and warnings. Use this tool when investigating why resources are not working as expected, such as pods stuck in pending state, failed deployments, or service connectivity issues. Events are automatically cleaned up after a retention period (typically 1 hour), so recent events are most relevant for troubleshooting."),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to filter events from. If not specified, events from all namespaces will be returned (requires cluster-wide permissions). For focused troubleshooting, specify the namespace where the problematic resources are located. Common namespaces include 'default', 'kube-system' (for cluster components), 'kube-public', or custom application namespaces. Use this to narrow down events when you know which namespace contains the resources you're investigating.")),
		mcp.WithString("fieldSelector",
			mcp.Description("Field selector to filter events based on specific criteria. This allows precise filtering of events related to specific resources or conditions. Common examples: 'involvedObject.name=my-pod' (events for a specific pod), 'involvedObject.kind=Pod' (all pod-related events), 'type=Warning' (only warning events), 'type=Normal' (only normal events), 'reason=Failed' (events with Failed reason), 'involvedObject.namespace=my-namespace' (events for resources in specific namespace). You can combine multiple selectors with commas. This is particularly useful when troubleshooting specific resources or looking for particular types of issues.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of events to return in the response. Default is 100 events if not specified. Set a higher limit (e.g., 500, 1000) when you need to see more historical events for comprehensive troubleshooting. Set a lower limit (e.g., 20, 50) for quick checks or when you only need recent events. Be mindful that very high limits may return large amounts of data and take longer to process. Events are typically returned in reverse chronological order (newest first), so limiting helps focus on the most recent activities.")),
		mcp.WithString("debug",
			mcp.Description("Enable detailed debug output for troubleshooting the tool itself (true/false). When set to 'true', provides additional logging information about the API calls, authentication, and processing steps. Use this when the get_events tool itself is not working as expected or when you need to understand the underlying Kubernetes API interactions. This is separate from the Kubernetes events themselves and is used for debugging the tool's operation.")),
	)
}

// GetResourceUsageTool retrieves resource usage information (CPU/Memory)
func GetResourceUsageTool() mcp.Tool {
	logrus.Debug("Creating GetResourceUsageTool")
	return mcp.NewTool("kubernetes_get_resource_usage",
		mcp.WithDescription("Retrieve real-time resource usage metrics (CPU and Memory) for Kubernetes nodes or pods. This tool queries the metrics server to show current resource consumption, which is essential for performance monitoring, capacity planning, and troubleshooting resource-related issues. The output includes current usage values and percentages relative to requests/limits. Note: This requires the metrics-server to be installed and running in your cluster. Use this tool when you need to: monitor resource consumption, identify resource-hungry pods, check node capacity, troubleshoot performance issues, or validate resource requests/limits. For nodes, shows total usage across all pods. For pods, shows per-container breakdown when available."),
		mcp.WithString("resourceType", mcp.Required(),
			mcp.Description("Type of Kubernetes resource to check usage metrics for. Valid values: 'node' (for cluster node resource usage including CPU, memory, and capacity information across all pods running on each node), 'pod' (for individual pod resource usage showing CPU and memory consumption per container). Use 'node' to get cluster-wide resource overview and identify resource pressure. Use 'pod' to analyze specific application resource consumption and identify resource-hungry containers.")),
		mcp.WithString("name",
			mcp.Description("Specific name of the resource instance to check usage for. When resourceType is 'node': provide the exact node name (e.g., 'worker-node-1', 'ip-10-0-1-123.ec2.internal'). When resourceType is 'pod': provide the exact pod name. If omitted, shows usage metrics for ALL resources of the specified type, which is useful for getting an overview. Use list_resources tool first if you're unsure about exact names. Node names can be found with 'kubectl get nodes', pod names with 'kubectl get pods'.")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace where the pod is located. This parameter is REQUIRED when resourceType is 'pod' since pods are namespaced resources. Ignored when resourceType is 'node' since nodes are cluster-scoped. If you're unsure about the namespace, use the list_resources tool first to find pods and their namespaces. Common namespaces include: 'default' (default namespace), 'kube-system' (system pods), 'kube-public' (publicly readable), or custom application namespaces. Example: if checking pod 'nginx-deployment-abc123' in namespace 'production', set this to 'production'.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the metrics API call and tool execution. Set to 'true' to see detailed request/response information, API endpoints being called, authentication details, and any errors from the metrics server. Set to 'false' or omit for normal output showing only the resource usage metrics. Use debug mode when: metrics server seems unavailable, getting authentication errors, tool returns unexpected results, or when troubleshooting cluster metrics collection issues.")),
	)
}

// PortForwardTool creates port forwarding to a pod
func PortForwardTool() mcp.Tool {
	logrus.Debug("Creating PortForwardTool")
	return mcp.NewTool("kubernetes_port_forward",
		mcp.WithDescription("Create port forwarding from a local port to a pod port, similar to 'kubectl port-forward'. This tool establishes a network tunnel that allows you to access services running inside a pod from your local machine. This is particularly useful for debugging applications, accessing databases, web interfaces, or APIs that are not exposed through Kubernetes services. The port forward session remains active until explicitly stopped. Use this when you need direct access to a pod's network interface for development, testing, or troubleshooting purposes. Make sure the target pod is running and the specified pod port is actually listening for connections."),
		mcp.WithString("podName", mcp.Required(),
			mcp.Description("Exact name of the target pod to forward traffic to. The pod must be in 'Running' state for port forwarding to work. Use 'list_resources' tool with kind='Pod' to find available pod names if needed. Pod names are case-sensitive and must match exactly as they appear in the cluster. If the pod restarts or gets recreated, you'll need to establish a new port forward session.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Kubernetes namespace where the target pod is located. This is required as pods are namespaced resources. Common namespaces include 'default', 'kube-system', 'kube-public', or custom application namespaces. Use 'list_resources' or 'get_resource_details' tools to verify the pod's namespace if uncertain. Namespace names are case-sensitive.")),
		mcp.WithNumber("localPort", mcp.Required(),
			mcp.Description("Local port number on your machine to bind the port forward to. This is the port you'll connect to locally (e.g., http://localhost:8080). Choose a port that's not already in use on your local system. Common ranges: 8000-8999 for web services, 5432 for PostgreSQL, 3306 for MySQL, 6379 for Redis. The port must be between 1-65535 and available for binding.")),
		mcp.WithNumber("podPort", mcp.Required(),
			mcp.Description("Target port number inside the pod that you want to access. This must be a port that the application inside the pod is actually listening on. Check the pod's container specifications, service definitions, or use 'describe_resource' to find the correct port. Common examples: 80/8080 for web servers, 443 for HTTPS, 3000 for Node.js apps, 5000 for Python Flask, 8000 for Django. The port must be between 1-65535.")),
		mcp.WithString("address",
			mcp.Description("Local IP address to bind the port forward to. Defaults to 'localhost' (127.0.0.1) which only allows local connections. Use '0.0.0.0' to allow connections from other machines on your network (security risk - use carefully). For most debugging and development purposes, the default 'localhost' is recommended for security. IPv6 addresses are also supported (e.g., '::1' for IPv6 localhost).")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting port forward setup and connection issues. Set to 'true' to see detailed information about the port forward session establishment, traffic flow, and any connection errors. Set to 'false' or omit for normal output. Debug mode is helpful when diagnosing connectivity issues or when the port forward fails to establish properly.")),
	)
}

// CreateResourceTool creates any Kubernetes resource
func CreateResourceTool() mcp.Tool {
	logrus.Debug("Creating CreateResourceTool")
	return mcp.NewTool("kubernetes_create_resource",
		mcp.WithDescription("Create a new Kubernetes resource by providing the complete resource manifest. This tool accepts a structured resource definition and creates it in the cluster, similar to 'kubectl apply' or 'kubectl create'. Use this tool when you need to deploy new applications, create configuration objects, set up RBAC resources, or provision any Kubernetes resource type. The tool requires you to provide the complete resource specification including apiVersion, kind, metadata, and spec fields. Before creating resources, ensure you have the necessary permissions and that the target namespace exists. For complex deployments, consider creating dependencies first (e.g., ConfigMaps and Secrets before Deployments that reference them). This operation is idempotent - if a resource with the same name already exists, the creation will fail unless you use the update_resource tool instead."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("The Kubernetes resource kind (type) to create. This must match exactly with valid Kubernetes API resource kinds. Common examples include: 'Pod' (for single containers), 'Deployment' (for scalable application workloads), 'Service' (for networking and load balancing), 'ConfigMap' (for configuration data), 'Secret' (for sensitive information like passwords and certificates), 'Namespace' (for resource isolation), 'Ingress' (for HTTP/HTTPS routing), 'PersistentVolume' and 'PersistentVolumeClaim' (for storage), 'ServiceAccount' (for pod identity), 'Role' and 'ClusterRole' (for RBAC permissions), 'StatefulSet' (for stateful applications), 'DaemonSet' (for node-level services), 'Job' and 'CronJob' (for batch workloads). The kind is case-sensitive and must use the exact capitalization as defined in the Kubernetes API (e.g., 'ConfigMap', not 'configmap').")),
		mcp.WithString("apiVersion", mcp.Required(),
			mcp.Description("The Kubernetes API version for the resource type being created. This determines which version of the resource schema to use and must match the kind specified. Common API versions include: 'v1' (for core resources like Pod, Service, ConfigMap, Secret, Namespace, PersistentVolume, PersistentVolumeClaim), 'apps/v1' (for Deployment, StatefulSet, DaemonSet, ReplicaSet), 'batch/v1' (for Job), 'batch/v1beta1' (for CronJob), 'networking.k8s.io/v1' (for Ingress, NetworkPolicy), 'rbac.authorization.k8s.io/v1' (for Role, ClusterRole, RoleBinding, ClusterRoleBinding), 'storage.k8s.io/v1' (for StorageClass), 'autoscaling/v2' (for HorizontalPodAutoscaler). Use 'kubectl api-resources' to verify the correct apiVersion for specific resource types in your cluster. The format is typically 'group/version' for non-core resources and just 'version' for core resources.")),
		mcp.WithString("metadata", mcp.Required(),
			mcp.Description("Resource metadata as a properly formatted JSON string containing essential resource identification and organizational information. This must include at minimum the 'name' field, and 'namespace' field for namespaced resources. Required structure: '{\"name\": \"resource-name\", \"namespace\": \"target-namespace\"}'. Optional but commonly used fields include: 'labels' (for resource organization and selection, e.g., '{\"app\": \"web\", \"version\": \"v1.0\"}'), 'annotations' (for non-identifying metadata, e.g., '{\"description\": \"Web server deployment\"}'), 'generateName' (for auto-generated names when name is not specified). Example for a namespaced resource: '{\"name\": \"my-app\", \"namespace\": \"production\", \"labels\": {\"app\": \"web\", \"tier\": \"frontend\"}, \"annotations\": {\"description\": \"Production web server\"}}'. For cluster-scoped resources (like ClusterRole, PersistentVolume), omit the namespace field. Resource names must follow DNS subdomain format: lowercase alphanumeric characters, hyphens, and dots, starting and ending with alphanumeric characters, and be no more than 253 characters long.")),
		mcp.WithString("spec",
			mcp.Description("Resource specification as a properly formatted JSON string defining the desired state and configuration of the resource. This field contains the main configuration that defines how the resource should behave and is specific to each resource kind. For Deployment: include 'replicas', 'selector', and 'template' with pod specification. For Service: include 'selector', 'ports', and 'type'. For ConfigMap: include 'data' or 'binaryData' fields. For Secret: include 'data' (base64 encoded) or 'stringData' (plain text) fields. For PersistentVolumeClaim: include 'accessModes', 'resources', and optionally 'storageClassName'. Some resources like ConfigMap and Secret may have spec content in metadata instead. Example for Deployment spec: '{\"replicas\": 3, \"selector\": {\"matchLabels\": {\"app\": \"web\"}}, \"template\": {\"metadata\": {\"labels\": {\"app\": \"web\"}}, \"spec\": {\"containers\": [{\"name\": \"web\", \"image\": \"nginx:1.20\", \"ports\": [{\"containerPort\": 80}]}]}}}'. Refer to Kubernetes API documentation for exact schema requirements for each resource kind. This field may be optional for some resource types that store their configuration in metadata instead of spec.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting resource creation issues. Set to 'true' to see detailed information about the API request, response, and any validation errors. Set to 'false' or omit for normal output showing only success/failure status. Debug output includes the complete resource manifest being sent to the API server, HTTP request details, and detailed error messages if creation fails. Use this when resources fail to create and you need to understand why, such as validation errors, permission issues, or API server problems. Debug information can help identify issues like malformed JSON, missing required fields, invalid field values, or insufficient permissions.")),
	)
}

// ListResourcesTool lists Kubernetes resources of a given kind
func ListResourcesTool() mcp.Tool {
	logrus.Debug("Creating ListResourcesTool")
	return mcp.NewTool("kubernetes_list_resources",
		mcp.WithDescription("📊 STANDARD TOOL: List with filtering and pagination - use when need more detail than summary. 🔍 Standard scenarios: need more information than summary but not full configuration. ⚠️ Recommendation: prioritize list_resources_summary, use this tool when more fields are needed, use list_resources_full for complete configuration. 📋 Use cases: inventory, log_analysis. 🔧 Features: supports labelSelector/fieldSelector filtering, JSONPath field extraction, pagination continuation."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind/type to list - the category of resources you want to discover. Common resource types include: 'Pod' (running containers and applications), 'Service' (network services and load balancing), 'Deployment' (application deployments and replica management), 'ConfigMap' (configuration data), 'Secret' (sensitive information like passwords and certificates), 'Ingress' (HTTP/HTTPS routing rules), 'PersistentVolume' and 'PersistentVolumeClaim' (storage resources), 'Namespace' (resource organization and isolation), 'Node' (cluster infrastructure), 'DaemonSet' (node-wide services), 'StatefulSet' (stateful applications), 'Job' and 'CronJob' (batch workloads), 'ServiceAccount' (identity and permissions), 'Role' and 'ClusterRole' (RBAC permissions), 'CustomResource' (custom resource definitions). Use exact case-sensitive names as they appear in Kubernetes API (e.g., 'Pod' not 'pod'). If unsure about available resource types, try common ones first.")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scope the resource listing. For namespaced resources (Pod, Service, Deployment, ConfigMap, Secret, etc.), this filters results to only show resources within the specified namespace. If omitted for namespaced resources, shows resources from ALL namespaces (requires cluster-wide list permissions). For cluster-scoped resources (Node, ClusterRole, PersistentVolume, etc.), this parameter is ignored. Common namespaces include: 'default' (default namespace for user resources), 'kube-system' (system components and cluster services), 'kube-public' (publicly accessible resources), 'kube-node-lease' (node heartbeat data), or custom application namespaces. Use this to narrow down results when working with specific applications or environments. Leave empty to get a cluster-wide view.")),
		mcp.WithString("labelSelector",
			mcp.Description("Label selector to filter resources based on their metadata labels. Labels are key-value pairs attached to resources for organization and selection. Use this powerful filtering mechanism to find resources matching specific criteria. Syntax examples: 'app=nginx' (resources with label app=nginx), 'env=production' (production environment resources), 'app=nginx,env=prod' (multiple labels with AND logic), 'tier in (frontend,backend)' (resources with tier label having specific values), 'app!=legacy' (exclude resources with app=legacy), 'version' (resources that have a 'version' label regardless of value), '!debug' (resources that do NOT have a 'debug' label). Common label patterns: 'app' (application name), 'version' (application version), 'env' (environment like dev/staging/prod), 'tier' (application tier like frontend/backend), 'release' (deployment release). Combine multiple selectors with commas for AND logic.")),
		mcp.WithString("fieldSelector",
			mcp.Description("Field selector to filter resources based on specific field values. This provides more precise filtering compared to label selectors. Common examples: 'involvedObject.name=my-pod' (events for a specific pod), 'involvedObject.kind=Pod' (all pod-related events), 'type=Warning' (only warning events), 'type=Normal' (only normal events), 'reason=Failed' (events with Failed reason), 'involvedObject.namespace=my-namespace' (events for resources in specific namespace), 'status.phase=Running' (running pods), 'status.phase!=Failed' (exclude failed resources). You can combine multiple selectors with commas. Field selectors work on actual resource fields rather than labels and are useful for filtering by runtime state or resource properties.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of resources to return in a single page (default: 100, max: 500). This parameter enables pagination and prevents context overflow by limiting response size. Use smaller limits (10-50) for quick overviews or when you only need a few resources. Use larger limits (100-500) for comprehensive analysis. When 'hasMore' is true in the response, use the 'continueToken' to fetch the next page. The pagination is handled by the Kubernetes API server, so this is more efficient than client-side limiting.")),
		mcp.WithString("continueToken",
			mcp.Description("Pagination token from a previous response to fetch the next page of results. When the response indicates 'hasMore': true, use the provided 'continueToken' to get the next batch of resources. This enables efficient pagination through large result sets without loading all data into memory. Leave empty for the first page request. The token is opaque and should be used exactly as provided in the previous response's 'pagination.continueToken' field.")),
		mcp.WithString("jsonpath",
			mcp.Description("Single JSONPath expression to extract and format specific fields from each resource. This allows you to customize the output format and display only the information you need. IMPORTANT: Due to Kubernetes JSONPath limitations, you CANNOT create JSON objects directly - only text output is supported. The expression must start with '{' and end with '}'. VALID examples for text output: '{.metadata.name}' (show only resource names), '{.status.phase}' (show status for Pods), '{.metadata.namespace}' (show namespaces), '{range .items[*]}{.metadata.name}{\"\\t\"}{.status.phase}{\"\\n\"}{end}' (tabular format with name and status), '{range .items[*]}namespace={.metadata.namespace},name={.metadata.name},kind={.kind}{\"\\n\"}{end}' (comma-separated key=value pairs). For multiple fields, use text concatenation with separators like commas, tabs {\"\\t\"}, or newlines {\"\\n\"}. INVALID examples that will fail: '{\"name\": \"{.metadata.name}\"}' (JSON object syntax), '{range .items[*]}{\"namespace\": \"{.metadata.namespace}\"}{end}' (JSON formatting). Use simple text formatting with delimiters that can be parsed later if needed. Escape sequences: {\"\\t\"} for tabs, {\"\\n\"} for newlines, {\"\\\"\"} for quotes. Always test simple expressions first before using complex range operations.")),
		mcp.WithString("jsonpaths",
			mcp.Description("Multiple JSONPath expressions provided as a JSON string containing an array of JSONPath expressions. This allows extracting multiple specific fields from each resource in a structured format. Each expression in the array will be evaluated separately and the results combined. The parameter should be a valid JSON array of strings, for example: '[\"metadata.name\", \"metadata.namespace\", \"status.phase\"]' or '[\"spec.replicas\", \"status.readyReplicas\", \"metadata.labels.app\"]'. Each JSONPath expression should NOT include the surrounding curly braces (unlike the single jsonpath parameter) - just provide the path like 'metadata.name' instead of '{.metadata.name}'. This is useful when you need multiple specific fields without having to construct complex range expressions. The output will include results for each specified path. If a path doesn't exist for a particular resource, it will be omitted or show as null in the output. Cannot be used together with the single 'jsonpath' parameter - choose one approach based on your needs.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the tool execution and API interactions. Set to 'true' to see detailed information about the Kubernetes API calls, authentication process, request/response details, pagination tokens, and any filtering operations being applied. Set to 'false' or omit for normal output showing only the resource information. Debug mode is helpful when: the tool is not returning expected results, you're getting authentication or permission errors, pagination is not working as expected, or you're troubleshooting connectivity issues. Normal users should leave this unset or set to 'false' for cleaner output.")),
	)
}

// ListResourcesSummaryTool lists Kubernetes resources with minimal summary output
func ListResourcesSummaryTool() mcp.Tool {
	logrus.Debug("Creating ListResourcesSummaryTool")
	return mcp.NewTool("kubernetes_list_resources_summary",
		mcp.WithDescription("⚠️ PRIORITY: Optimized for LLM efficiency: Returns only essential fields (name, namespace, status, age). 80-90% smaller than detailed version. 🚀 Best for: quick browsing, resource discovery, health checks, initial diagnosis. 📋 Use cases: pod_troubleshooting, health_check, inventory, resource_monitoring. 🔄 Workflow: use this tool first to discover resources → use get_resource_summary to view details → combine with get_recent_events to analyze problems when needed."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Resource kind to list (Pod, Deployment, Service, ConfigMap, etc.). Use exact case-sensitive names as they appear in Kubernetes API. Common types: Pod, Service, Deployment, ConfigMap, Secret, Namespace, Node, Ingress, StatefulSet, DaemonSet, Job, CronJob.")),
		mcp.WithString("namespace",
			mcp.Description("Optional namespace filter. Omit for cluster-wide listing across all namespaces (requires cluster-wide permissions). For namespaced resources, this limits results to the specified namespace. Ignored for cluster-scoped resources like Node, PersistentVolume, ClusterRole.")),
		mcp.WithString("labelSelector",
			mcp.Description("Optional label selector for filtering resources (e.g., 'app=nginx', 'env=production', 'tier in (frontend,backend)'). Use combination with commas for AND logic: 'app=nginx,env=prod'. This helps narrow down results to specific applications or environments.")),
		mcp.WithString("includeLabels",
			mcp.Description("Optional comma-separated label keys to include in the summary output (e.g., 'app,version,env'). When specified, only these labels will be included for each resource. If omitted, a limited set of labels (max 10) will be included automatically. Useful for reducing output size while maintaining essential context.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of resources to return (default: 100, max: 500). This enables server-side pagination to prevent context overflow. Use smaller values (10-50) for quick overviews, larger values (100-500) for comprehensive analysis. Pagination is handled by Kubernetes API for efficiency.")),
		mcp.WithString("continueToken",
			mcp.Description("Pagination token from previous response to fetch the next page. When response indicates 'hasMore': true, use the provided 'continueToken' to get the next batch. Leave empty for the first request. This enables efficient traversal of large result sets without loading all data.")),
	)
}

// GetResourceDetailsTool retrieves detailed information about a specific Kubernetes resource
func GetResourceDetailsTool() mcp.Tool {
	logrus.Debug("Creating GetResourceDetailsTool")
	return mcp.NewTool("kubernetes_get_resource_details",
		mcp.WithDescription("Retrieve comprehensive detailed information about a specific Kubernetes resource including its current configuration, status, metadata, and relationships. This tool provides an in-depth view of a single resource instance, showing all fields and nested properties in a structured format. Unlike the describe_resource tool which provides human-readable formatted output, this tool returns the complete raw resource data including all API fields, custom resource definitions, and technical details. Use this tool when you need to: examine the exact configuration of a resource, inspect status conditions and detailed state information, analyze resource relationships and ownership, extract specific field values for automation, or understand the complete resource structure for troubleshooting. This is particularly useful for debugging complex resource configurations, understanding why a resource is in a particular state, or when you need programmatic access to specific resource properties."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("The Kubernetes resource kind (type) to retrieve detailed information for. This must be an exact, case-sensitive match with valid Kubernetes API resource kinds. Common examples include: 'Pod' (for containerized applications and their runtime state), 'Service' (for network service definitions and endpoints), 'Deployment' (for declarative application deployments and rollout status), 'ConfigMap' (for configuration data and key-value pairs), 'Secret' (for sensitive information like passwords, tokens, and certificates), 'Namespace' (for resource grouping and isolation), 'Ingress' (for HTTP/HTTPS routing rules and TLS configuration), 'PersistentVolume' and 'PersistentVolumeClaim' (for storage resources and binding status), 'ServiceAccount' (for pod identity and token information), 'Node' (for cluster infrastructure and node conditions), 'DaemonSet' (for node-level workloads), 'StatefulSet' (for stateful applications with persistent identity), 'Job' and 'CronJob' (for batch workloads and scheduling). The kind must match exactly as it appears in the Kubernetes API specification.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The exact name of the specific resource instance to retrieve detailed information for. This must match the metadata.name field of the resource exactly as it exists in the cluster. Resource names are case-sensitive and must follow Kubernetes naming conventions (lowercase alphanumeric characters, hyphens, and dots are allowed, but no spaces or special characters). For resources created by controllers (like Pods created by Deployments), the name will include generated suffixes or prefixes. If you're unsure about the exact resource name, use the 'list_resources' tool first to discover available resources and their exact names. Examples: 'nginx-deployment-7fb96c846b-xyz12' (for a Pod), 'my-app-service' (for a Service), 'web-app-deployment' (for a Deployment). The name must exist in the specified namespace (for namespaced resources) or in the cluster (for cluster-scoped resources).")),
		mcp.WithString("namespace",
			mcp.Description("The Kubernetes namespace where the resource is located. This parameter is REQUIRED for namespaced resources (such as Pod, Service, Deployment, ConfigMap, Secret, Ingress, PersistentVolumeClaim, ServiceAccount, Role, RoleBinding, etc.) but should be OMITTED for cluster-scoped resources (such as Node, PersistentVolume, ClusterRole, ClusterRoleBinding, Namespace itself, etc.). If you're unsure whether a resource type is namespaced or cluster-scoped, try the operation without specifying a namespace first - the error message will indicate if a namespace is required. Common namespace examples: 'default' (the default namespace if none was specified during resource creation), 'kube-system' (for Kubernetes system components), 'kube-public' (for publicly accessible resources), or custom application namespaces like 'production', 'staging', 'development'. Use the 'list_resources' tool to discover which namespaces contain your target resources if uncertain.")),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive debug output for troubleshooting the tool execution and API interactions. Set to 'true' to see detailed information including: Kubernetes API endpoints being called, authentication and authorization details, request and response headers and bodies, error messages and stack traces, timing information for API calls, and internal tool processing steps. Set to 'false' or omit this parameter for normal operation with standard output showing only the resource details. Debug mode is particularly useful when: the tool is not returning expected results, you're getting authentication or permission errors, the resource seems to exist but cannot be retrieved, you need to understand the underlying API calls for automation purposes, or when reporting issues with the tool itself. Note that debug output may contain sensitive information and should be used carefully in production environments.")),
	)
}

// UpdateResourceTool updates a Kubernetes resource by replacement
func UpdateResourceTool() mcp.Tool {
	logrus.Debug("Creating UpdateResourceTool")
	return mcp.NewTool("kubernetes_update_resource",
		mcp.WithDescription("Update an existing Kubernetes resource by replacing it with a new complete manifest definition. This tool performs a full resource replacement operation, similar to 'kubectl replace' or 'kubectl apply'. Use this when you need to modify existing resources such as updating image versions in Deployments, changing service configurations, modifying resource limits, updating environment variables, or applying configuration changes. The tool requires the complete resource manifest including all required fields. Important: This operation replaces the entire resource, so ensure your manifest includes all desired configuration. For partial updates or strategic merges, consider using patch operations instead. The resource must already exist - use create_resource for new resources. Always verify the current resource state with get_resource before updating to avoid overwriting recent changes made by other users or controllers."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind (type) of the resource to update. This must match the existing resource's kind exactly. Common examples: 'Deployment' (for updating application deployments like image versions, replica counts, environment variables), 'Service' (for modifying service configurations, ports, selectors), 'ConfigMap' (for updating configuration data), 'Secret' (for updating sensitive information), 'Ingress' (for modifying routing rules), 'StatefulSet' (for stateful application updates), 'DaemonSet' (for node-level service updates), 'Pod' (though direct pod updates are limited - most fields are immutable after creation). The kind is case-sensitive and must use exact Kubernetes API capitalization (e.g., 'ConfigMap', not 'configmap'). Use get_resource or list_resources tools first to verify the existing resource's kind if uncertain.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the existing resource to update. This must match the resource's metadata.name field precisely as it appears in the cluster. Resource names are case-sensitive and must follow Kubernetes naming conventions (lowercase alphanumeric characters, hyphens, and dots allowed). Use list_resources or get_resource tools to find the correct resource name if you're unsure. For resources created by higher-level controllers (like Pods created by Deployments), be aware that direct updates may be overwritten by the controller - update the controlling resource instead (e.g., update the Deployment rather than individual Pods).")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace where the target resource exists. This is REQUIRED for namespaced resources (Pod, Service, Deployment, ConfigMap, Secret, Ingress, StatefulSet, DaemonSet, etc.) and must match the namespace where the resource currently resides. OMIT this parameter for cluster-scoped resources (Node, ClusterRole, ClusterRoleBinding, PersistentVolume, StorageClass, etc.). If you're unsure whether a resource type is namespaced, try the operation without namespace first - the error message will indicate if namespace is required. Use get_resource or list_resources tools to verify the correct namespace of existing resources. Common namespaces include 'default' (default namespace), 'kube-system' (system components), or custom application namespaces.")),
		mcp.WithString("manifest", mcp.Required(),
			mcp.Description("Complete Kubernetes resource manifest as a JSON string containing the full updated resource definition. This must include ALL required fields: 'apiVersion' (e.g., 'apps/v1' for Deployments, 'v1' for Services/ConfigMaps), 'kind' (must match the kind parameter), 'metadata' (including 'name' and 'namespace' if applicable), and 'spec' (the resource specification with your desired changes). IMPORTANT: Include the current 'metadata.resourceVersion' from the existing resource to ensure you're updating the latest version and avoid conflicts. The manifest should represent the complete desired state after the update, not just the changes. For complex resources, retrieve the current manifest using get_resource, modify the necessary fields, and provide the complete updated manifest here. Ensure proper JSON formatting with correct nesting, data types (strings in quotes, numbers without quotes, arrays with square brackets, objects with curly braces), and valid Kubernetes field names and values.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the update operation and API interactions. Set to 'true' to see detailed information about the API request being sent, response from Kubernetes API server, validation errors, field conflicts, and any issues during the update process. Set to 'false' or omit for normal output showing only the update result. Use debug mode when: the update operation fails with unclear errors, you need to verify what data is being sent to the API, there are validation or schema errors, you suspect API permission issues, or when troubleshooting resource version conflicts. Debug output helps identify issues with manifest formatting, missing required fields, or API communication problems.")),
	)
}

// DeleteResourceTool deletes a Kubernetes resource
func DeleteResourceTool() mcp.Tool {
	logrus.Debug("Creating DeleteResourceTool")
	destructive := true
	return mcp.NewTool("kubernetes_delete_resource",
		mcp.WithDescription("Permanently delete a Kubernetes resource from the cluster. This is an IRREVERSIBLE and DESTRUCTIVE operation that completely removes the specified resource and all its associated data. Use this tool with extreme caution as deleted resources cannot be recovered unless you have backups. Before deletion, consider: 1) Backing up the resource manifest using get_resource tool, 2) Checking if other resources depend on this one, 3) Verifying you have the correct resource name and namespace, 4) Understanding the impact on running applications. Deleting certain resources can cause service disruptions (Pods will terminate, Services will stop routing traffic, PersistentVolumes may lose data, Deployments will stop managing replicas). For graceful application shutdown, consider scaling down Deployments to 0 replicas first. Some resources like Namespaces will delete ALL contained resources. Always double-check the resource details before proceeding with deletion."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind (type) to delete. This must be the exact, case-sensitive resource type as it appears in Kubernetes API. Common examples: 'Pod' (will terminate the running container immediately), 'Service' (will stop network routing to pods), 'Deployment' (will delete the deployment and all its managed pods), 'ConfigMap' (will remove configuration data that pods may be using), 'Secret' (will remove sensitive data like certificates and passwords), 'PersistentVolumeClaim' (may cause data loss if pods are using the storage), 'Ingress' (will stop HTTP/HTTPS routing), 'Namespace' (DANGER: deletes ALL resources within the namespace), 'StatefulSet' (will delete stateful application and associated pods), 'DaemonSet' (will remove the daemon from all nodes). Be especially careful with storage-related resources (PVC, PV) as deletion may result in permanent data loss. Use list_resources or get_resource_details tools first to confirm the exact kind if uncertain.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the specific resource instance to delete. This must match the metadata.name field of the resource exactly - names are case-sensitive and must be spelled precisely. Use the list_resources tool or get_resource_details tool first to verify the correct resource name if you're unsure. For resources created by higher-level controllers (like Pods created by Deployments), deleting individual instances may cause them to be recreated automatically - consider deleting the parent resource instead. Resource names must follow Kubernetes naming conventions (lowercase alphanumeric characters, hyphens, and dots). Double-check this value as deletion is permanent and cannot be undone. Typos in the name will either fail the deletion or potentially delete the wrong resource.")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace where the target resource is located. This is REQUIRED for namespaced resources including Pod, Service, Deployment, ConfigMap, Secret, PersistentVolumeClaim, Ingress, and most application-level resources. This parameter should be OMITTED for cluster-scoped resources like Node, ClusterRole, PersistentVolume, StorageClass, and Namespace itself. If you're unsure whether a resource type requires a namespace, use the list_resources tool first to see how the resource is structured. Common namespaces include: 'default' (default namespace for resources), 'kube-system' (critical system components - be very careful), 'kube-public' (publicly accessible resources), or custom application namespaces. Namespace names are case-sensitive. Deleting resources from 'kube-system' can break cluster functionality - only proceed if you're absolutely certain of the impact.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the deletion operation. Set to 'true' to see detailed information about the deletion process, API calls, authentication, and any errors or warnings during the operation. Set to 'false' or omit for normal output showing only the deletion result. Debug mode is helpful when: deletion fails unexpectedly, you need to understand why a resource cannot be deleted (due to finalizers, dependencies, or permissions), you want to see the exact API calls being made, or when troubleshooting authentication/authorization issues. Debug output may include sensitive information like resource configurations, so use carefully in production environments.")),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				DestructiveHint: &destructive,
			},
		),
	)
}

// ContainerLogsTool retrieves logs from a Pod container
func ContainerLogsTool() mcp.Tool {
	logrus.Debug("Creating ContainerLogsTool")
	return mcp.NewTool("kubernetes_get_pod_logs",
		mcp.WithDescription("Retrieve and display container logs from a Kubernetes Pod, similar to 'kubectl logs'. This tool is essential for debugging applications, monitoring application behavior, troubleshooting startup issues, and investigating runtime errors. Container logs contain stdout and stderr output from the application processes running inside containers. Use this tool when you need to: diagnose why a pod is failing or crashing, monitor application output and error messages, investigate performance issues or unexpected behavior, check application startup sequences, or analyze error patterns. The tool can retrieve logs from current running containers or previous container instances (if the container has restarted). For multi-container pods, you must specify the container name. Logs are returned in chronological order and can be limited to recent entries for better performance."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the Pod from which to retrieve container logs. The pod name must match exactly as it appears in Kubernetes and is case-sensitive. Pod names typically follow patterns like 'deployment-name-random-suffix' for pods created by Deployments, or custom names for manually created pods. Use 'list_resources' tool with kind='Pod' first if you're unsure of the exact pod name. The pod can be in any state (Running, Pending, Failed, Succeeded) but must exist in the cluster. For pods created by controllers like Deployments, the name includes generated suffixes (e.g., 'nginx-deployment-abc123-xyz789').")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Kubernetes namespace where the target Pod is located. This is required since Pods are namespaced resources. Common namespaces include 'default' (default namespace for user workloads), 'kube-system' (system components and cluster services), 'kube-public' (publicly accessible resources), or custom application namespaces like 'production', 'staging', 'development'. If you're unsure about the namespace, use 'list_resources' tool with kind='Pod' to discover pods across namespaces. Namespace names are case-sensitive and must match exactly.")),
		mcp.WithString("container",
			mcp.Description("Name of the specific container within the Pod to retrieve logs from. This parameter is REQUIRED for multi-container pods since each container has separate logs. For single-container pods, this parameter is optional and will default to the only container. Container names are defined in the Pod specification under spec.containers[].name field. Common container names include 'app', 'main', 'web', 'api', or descriptive names like 'nginx', 'redis', 'database'. Use 'get_resource' or 'describe_resource' tools to inspect the pod and find container names if needed. If you specify a non-existent container name, the operation will fail with an error.")),
		mcp.WithNumber("tailLines",
			mcp.Description("Maximum number of recent log lines to retrieve from the end of the log stream. This helps limit output size and focus on recent activity. Default is 100 lines if not specified. Common values: 50 (quick check of recent activity), 100-200 (standard troubleshooting), 500-1000 (detailed investigation), 0 (retrieve all available logs - use carefully as this can be very large). Higher values may take longer to retrieve and consume more memory. For very active applications, even 100 lines might represent only a few seconds of activity. Set to a higher value when you need more historical context for debugging complex issues.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the log retrieval operation itself. Set to 'true' to see detailed information about the API calls, authentication, pod discovery, and any errors encountered while accessing logs. Set to 'false' or omit for normal output showing only the container logs. Debug mode is useful when: the tool fails to retrieve logs, you're getting authentication errors, the pod or container cannot be found, or when you need to understand the underlying Kubernetes API interactions. This debug output is separate from and in addition to the actual container logs.")),
	)
}

// ContainerExecTool executes commands in a Pod container
func ContainerExecTool() mcp.Tool {
	logrus.Debug("Creating ContainerExecTool")
	return mcp.NewTool("kubernetes_pod_exec",
		mcp.WithDescription("Execute commands inside a running container within a Kubernetes Pod, similar to 'kubectl exec'. This tool provides direct access to the container's runtime environment for debugging, troubleshooting, and administrative tasks. Use this tool when you need to: investigate application issues by examining files or processes inside containers, run diagnostic commands to check connectivity or resource usage, access application logs or configuration files directly, perform maintenance tasks like clearing caches or temporary files, test network connectivity from within the container, or install debugging tools temporarily. The pod must be in 'Running' state for command execution to work. Commands are executed with the same user privileges as the container's main process unless the container runs as root. Be cautious with destructive commands as they can affect the running application. For security reasons, avoid executing commands that modify critical system files or expose sensitive information."),
		mcp.WithString("podName", mcp.Required(),
			mcp.Description("Exact name of the target Pod where the command will be executed. The pod must be in 'Running' state for command execution to succeed. Pod names are case-sensitive and must match exactly as they appear in the cluster. For pods created by Deployments or other controllers, the name typically includes generated suffixes (e.g., 'nginx-deployment-7fb96c846b-xyz12'). Use 'list_resources' tool with kind='Pod' to discover available pod names if needed. The pod must exist and be accessible - if it's in Pending, Failed, or Succeeded state, command execution will fail. Multi-container pods require specifying the target container name in the containerName parameter.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Kubernetes namespace where the target Pod is located. This is required since Pods are namespaced resources. Common namespaces include: 'default' (default namespace for user workloads), 'kube-system' (system components - be careful with commands here), 'kube-public' (publicly accessible resources), or custom application namespaces like 'production', 'staging', 'development'. If you're unsure about the pod's namespace, use 'list_resources' tool with kind='Pod' to find pods across namespaces. Namespace names are case-sensitive and must match exactly as they exist in the cluster. Executing commands in system namespaces like 'kube-system' requires extra caution as it may affect cluster operations.")),
		mcp.WithString("containerName",
			mcp.Description("Name of the specific container within the Pod where the command should be executed. This parameter is REQUIRED for multi-container pods since you must specify which container to target. For single-container pods, this parameter is optional and will default to the only available container. Container names are defined in the Pod specification under spec.containers[].name field. Common container names include 'app', 'main', 'web', 'api', 'sidecar', or descriptive names like 'nginx', 'redis', 'database', 'proxy'. Use 'get_resource' or 'describe_resource' tools to inspect the pod specification and find the correct container names. If you specify a non-existent container name, the operation will fail with a 'container not found' error. Each container has its own filesystem and process space, so choose the container that contains the files or processes you need to access.")),
		mcp.WithString("command", mcp.Required(),
			mcp.Description("Command to execute inside the container, specified as a properly formatted JSON array of strings where each array element represents a command argument. This format ensures proper argument parsing and handles spaces, special characters, and quoting correctly. Examples: '[\"ls\", \"-la\", \"/app\"]' (list files with details), '[\"ps\", \"aux\"]' (show running processes), '[\"cat\", \"/etc/hostname\"]' (read file contents), '[\"env\"]' (show environment variables), '[\"curl\", \"-I\", \"http://localhost:8080/health\"]' (test HTTP endpoint), '[\"df\", \"-h\"]' (check disk usage), '[\"netstat\", \"-tuln\"]' (show network connections), '[\"tail\", \"-n\", \"100\", \"/var/log/app.log\"]' (show recent log entries), '[\"find\", \"/app\", \"-name\", \"*.conf\"]' (find configuration files), '[\"wget\", \"--spider\", \"https://external-api.com\"]' (test external connectivity). Each command part must be a separate quoted string in the array. Do not use shell operators like pipes (|), redirects (>), or background execution (&) as they require shell interpretation - use '[\"sh\", \"-c\", \"your shell command here\"]' for complex shell operations. Common debugging commands: 'ls' (list files), 'cat' (read files), 'ps' (processes), 'top' (resource usage), 'netstat' (network), 'curl'/'wget' (HTTP requests), 'env' (environment), 'df' (disk space).")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting command execution and container access issues. Set to 'true' to see detailed information about the execution process including: connection establishment to the pod, container selection process, command parsing and validation, execution environment details, and any errors during command execution. Set to 'false' or omit for normal output showing only the command results. Debug mode is helpful when: commands fail to execute with unclear errors, you're getting permission or access denied errors, the container or pod cannot be found, network connectivity issues prevent execution, or when you need to understand the execution environment. Debug output helps identify issues like incorrect container names, pod states that prevent execution, or API communication problems. Use debug mode when troubleshooting tool behavior, not for debugging the applications inside containers.")),
	)
}

// CheckPermissionsTool verifies user permissions for Kubernetes resources
func CheckPermissionsTool() mcp.Tool {
	logrus.Debug("Creating CheckPermissionsTool")
	return mcp.NewTool("kubernetes_check_permissions",
		mcp.WithDescription("Verify and check user permissions for Kubernetes API operations using the SubjectAccessReview API, similar to 'kubectl auth can-i'. This tool is essential for troubleshooting authorization issues, validating RBAC configurations, and understanding what operations your current user/service account can perform on Kubernetes resources. Use this tool when you encounter permission denied errors, need to verify access before attempting operations, want to audit user capabilities, or are debugging RBAC policies. The tool checks whether the current authentication context (user, service account, or group) has permission to perform specific actions on Kubernetes resources. This is particularly useful for security auditing, troubleshooting failed operations due to insufficient permissions, and validating that RBAC rules are working as expected. The check is performed against the live cluster state and considers all applicable Role, ClusterRole, RoleBinding, and ClusterRoleBinding resources."),
		mcp.WithString("verb", mcp.Required(),
			mcp.Description("Kubernetes API verb (action) to check permission for. This represents the operation you want to verify access to. Common verbs include: 'get' (retrieve individual resource details), 'list' (retrieve multiple resources or resource listings), 'create' (create new resources), 'update' (modify existing resources), 'patch' (partially update resources), 'delete' (remove resources), 'deletecollection' (remove multiple resources), 'watch' (monitor resource changes), 'proxy' (proxy requests to resources). Special verbs include: 'use' (for resources like PodSecurityPolicies), 'bind' (for role binding), 'escalate' (for privilege escalation), 'impersonate' (for user impersonation). The verb must match exactly as defined in Kubernetes RBAC rules. Use lowercase as verbs are case-sensitive. Examples: check 'get' permission before using get_resource tool, check 'list' before using list_resources tool, check 'create' before using create_resource tool, check 'delete' before using delete_resource tool.")),
		mcp.WithString("resourceGroup", mcp.Required(),
			mcp.Description("Kubernetes API group that contains the resource type you want to check permissions for. This specifies which API group the resource belongs to in the Kubernetes API hierarchy. Common API groups include: '' or 'v1' (core API group for basic resources like Pod, Service, ConfigMap, Secret, Namespace, PersistentVolume, PersistentVolumeClaim), 'apps' or 'apps/v1' (for application workloads like Deployment, StatefulSet, DaemonSet, ReplicaSet), 'batch' or 'batch/v1' (for Job resources), 'batch/v1beta1' (for CronJob), 'networking.k8s.io' or 'networking.k8s.io/v1' (for Ingress, NetworkPolicy), 'rbac.authorization.k8s.io' or 'rbac.authorization.k8s.io/v1' (for Role, ClusterRole, RoleBinding, ClusterRoleBinding), 'storage.k8s.io/v1' (for StorageClass), 'autoscaling/v2' (for HorizontalPodAutoscaler), 'metrics.k8s.io' (for metrics server resources). For core resources use empty string '' or 'v1'. For other resources, use the full API group name. Use 'kubectl api-resources' to find the correct API group for specific resource types.")),
		mcp.WithString("resourceResource", mcp.Required(),
			mcp.Description("Kubernetes resource type (kind in lowercase plural form) to check permissions for. This must be the plural form of the resource name as it appears in the Kubernetes API. Common resource types include: 'pods' (for Pod resources), 'services' (for Service resources), 'deployments' (for Deployment resources), 'configmaps' (for ConfigMap resources), 'secrets' (for Secret resources), 'namespaces' (for Namespace resources), 'nodes' (for Node resources), 'persistentvolumes' (for PersistentVolume resources), 'persistentvolumeclaims' (for PersistentVolumeClaim resources), 'ingresses' (for Ingress resources), 'roles' (for Role resources), 'clusterroles' (for ClusterRole resources), 'rolebindings' (for RoleBinding resources), 'clusterrolebindings' (for ClusterRoleBinding resources), 'statefulsets' (for StatefulSet resources), 'daemonsets' (for DaemonSet resources), 'jobs' (for Job resources), 'cronjobs' (for CronJob resources). Use lowercase plural forms as they appear in RBAC rules and API paths. You can find the correct resource names using 'kubectl api-resources' command.")),
		mcp.WithString("subresource",
			mcp.Description("Kubernetes subresource to check permissions for, if you need to verify access to specific subresources rather than the main resource. Subresources are specific aspects or operations on a main resource that have separate permission controls. Common subresources include: 'status' (for updating resource status fields, often restricted to controllers), 'scale' (for scaling operations on Deployments, StatefulSets, ReplicaSets), 'log' or 'logs' (for accessing container logs), 'exec' (for executing commands in containers), 'portforward' (for port forwarding to pods), 'proxy' (for proxying requests to pods or services), 'attach' (for attaching to running containers), 'binding' (for pod binding operations), 'eviction' (for pod eviction), 'finalizers' (for managing resource finalizers). Leave empty to check permissions on the main resource itself. Subresource permissions are often more restrictive and may be granted separately from main resource permissions. Example: you might have 'get' permission on 'pods' but not on 'pods/log' subresource.")),
		mcp.WithString("resourceName",
			mcp.Description("Specific name of an individual resource instance to check permissions for. This parameter allows checking permissions on a named resource rather than the resource type in general. When specified, the permission check will verify whether you can perform the specified verb on this exact resource instance. This is useful for resources that have name-specific RBAC rules or when checking permissions on particular sensitive resources. For example, checking if you can delete a specific pod named 'critical-app-pod' versus checking if you can delete pods in general. Resource names are case-sensitive and must match exactly. Leave empty to check permissions on the resource type as a whole. Note: some RBAC policies grant permissions to specific named resources, so checking with and without a resource name may yield different results. This is particularly relevant for sensitive resources where access might be restricted to specific instances.")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scope the permission check to for namespaced resources. This parameter is REQUIRED when checking permissions on namespaced resources (like pods, services, deployments, configmaps, secrets, etc.) and should be OMITTED for cluster-scoped resources (like nodes, persistentvolumes, clusterroles, etc.). Namespace-specific permissions are common in RBAC configurations where users might have different access levels in different namespaces. For example, you might have full access in 'development' namespace but only read access in 'production' namespace. Common namespaces include: 'default' (default namespace for user resources), 'kube-system' (system components - often restricted), 'kube-public' (publicly accessible resources), or custom application namespaces like 'production', 'staging', 'development'. If you're unsure whether a resource is namespaced, try the check without namespace first - the tool will indicate if namespace is required. Leave empty for cluster-scoped resources or to check cluster-level permissions on namespaced resource types.")),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive debug output for troubleshooting the permission check operation and understanding RBAC evaluation. Set to 'true' to see detailed information including: the exact SubjectAccessReview request being sent to the Kubernetes API, response details from the authorization system, information about the current authentication context (user, groups, service account), any error messages or warnings during the permission check, and timing information for the API call. Set to 'false' or omit for normal output showing only the permission result (allowed/denied). Debug mode is particularly useful when: permission checks return unexpected results, you're troubleshooting RBAC configuration issues, you need to understand which user context is being evaluated, you're getting authentication errors, or when you need detailed information for security auditing. Debug output may contain sensitive information about your authentication context and should be used carefully in shared environments.")),
	)
}

// TestTool is a test tool requiring user confirmation for testing MCP tool functionality
func TestTool() mcp.Tool {
	logrus.Debug("Creating TestTool")
	confirmed := true
	return mcp.NewTool("kubernetes_test_tool",
		mcp.WithDescription("Test tool for validating MCP (Model Context Protocol) tool functionality and user confirmation workflows. This tool is designed for testing purposes to verify that the MCP integration is working correctly and that user confirmation mechanisms are functioning as expected. Use this tool when you need to: validate the MCP connection and tool execution pipeline, test user confirmation prompts and workflows, verify that tool annotations are being processed correctly, debug MCP tool parameter handling and validation, or ensure that the Kubernetes MCP server is responding properly to tool requests. This tool requires explicit user confirmation before execution to demonstrate confirmation workflows. The tool performs basic validation and returns success/failure status to help diagnose MCP integration issues. It is safe to use and does not modify any Kubernetes resources or cluster state."),
		mcp.WithBoolean("confirmed",
			mcp.Description("User confirmation flag that must be explicitly set to 'true' to proceed with tool execution. This is a required safety parameter to demonstrate confirmation workflows and prevent accidental execution. The default value is 'false', which will cause the tool to abort execution with a clear error message. Set this to 'true' only after you have read and understood the tool's purpose and confirmed that you want to proceed with the test operation. This parameter validates the MCP confirmation mechanism and ensures that users must explicitly acknowledge potentially impactful operations. The confirmation requirement helps establish proper workflows for tools that may require user awareness or approval before execution."),
			mcp.Required(),
			mcp.DefaultBool(false),
		),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive verbose debug output for troubleshooting MCP tool execution, parameter processing, and communication flows. Set to 'true' to see detailed information about: MCP request parsing and validation, parameter processing and type conversion, tool execution lifecycle and timing, response formatting and serialization, error handling and stack traces, authentication and authorization steps. Set to 'false' or omit for normal output showing only the test result and basic execution status. Debug mode is particularly useful when: MCP tools are not executing as expected, parameters are not being passed correctly, there are authentication or communication issues, you need to understand the tool execution flow for development purposes, or when reporting issues with the MCP integration. Debug output helps identify problems with tool definitions, parameter validation, or the underlying MCP protocol implementation.")),
		mcp.WithToolAnnotation(
			mcp.ToolAnnotation{
				Title:           "Confirm Test Tool Execution",
				DestructiveHint: &confirmed,
				OpenWorldHint:   &confirmed,
			},
		),
	)
}

// ScaleResourceTool scales a namespaced resource (e.g., Deployment, StatefulSet)
func ScaleResourceTool() mcp.Tool {
	logrus.Debug("Creating ScaleResourceTool")
	return mcp.NewTool("kubernetes_scale_resource",
		mcp.WithDescription("Scale a Kubernetes workload resource by adjusting the number of replica instances, similar to 'kubectl scale'. This tool modifies the replica count for scalable resources like Deployments, StatefulSets, and ReplicaSets to handle varying load demands, perform maintenance, or optimize resource usage. Scaling is a fundamental operation for managing application availability and resource consumption. Use this tool when you need to: increase replicas to handle higher traffic loads, decrease replicas to save resources during low-demand periods, scale down to zero for maintenance or cost optimization, quickly respond to performance issues, or test application behavior under different replica counts. The scaling operation is performed by updating the resource's spec.replicas field. For Deployments, this triggers a rolling update process. For StatefulSets, scaling respects ordered startup/shutdown procedures. Always consider the impact on application availability and resource constraints when scaling. Monitor resource usage and application performance after scaling operations to ensure desired outcomes."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind (type) that supports scaling operations. This must be a scalable resource type with replica management capabilities. Supported resource kinds include: 'Deployment' (most common - for stateless applications that can be scaled horizontally with rolling updates), 'StatefulSet' (for stateful applications requiring ordered scaling with persistent identity and storage), 'ReplicaSet' (lower-level replica controller, though Deployments are preferred for most use cases), 'ReplicationController' (legacy resource, Deployments are recommended instead). The kind is case-sensitive and must match exactly as defined in Kubernetes API (e.g., 'Deployment' not 'deployment'). Note: Not all Kubernetes resources support scaling - resources like Pod, Service, ConfigMap, Secret cannot be scaled directly. Use 'list_resources' or 'get_resource_details' tools first to verify the resource type and current replica count if uncertain. For most application workloads, 'Deployment' is the appropriate choice.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the scalable resource instance to modify. The name must match the metadata.name field of the resource precisely as it exists in the cluster. Resource names are case-sensitive and must follow Kubernetes naming conventions (lowercase alphanumeric characters, hyphens, and dots allowed). Use 'list_resources' tool with the appropriate kind to discover available scalable resources and their exact names if you're unsure. For example, if you want to scale a Deployment named 'web-app-deployment', use that exact string. The resource must already exist in the specified namespace - use 'create_resource' tool first if the resource doesn't exist. Verify the current state of the resource with 'get_resource_details' or 'describe_resource' tools before scaling to understand the current replica count and resource status.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Kubernetes namespace where the target scalable resource is located. This is required since scalable resources like Deployments, StatefulSets, and ReplicaSets are namespaced resources. The namespace must exist and contain the specified resource. Common namespaces include: 'default' (default namespace for user workloads), 'kube-system' (system components - be very careful when scaling these), 'kube-public' (publicly accessible resources), or custom application namespaces like 'production', 'staging', 'development', 'web-app', 'api-services'. Namespace names are case-sensitive and must match exactly. Use 'list_resources' tool to verify which namespace contains your target resource if uncertain. Be especially cautious when scaling resources in 'kube-system' namespace as these are often critical cluster components that could affect cluster stability if scaled incorrectly.")),
		mcp.WithNumber("replicas", mcp.Required(),
			mcp.Description("Target number of replica instances (pods) that the resource should maintain after scaling. This must be a non-negative integer (0 or greater). Common scaling scenarios: Set to 0 to completely stop the application (useful for maintenance, cost savings, or troubleshooting), Set to 1 for minimal resource usage while keeping the application available, Set to 2-3 for basic high availability and load distribution, Set to higher values (5, 10, 20+) for high-traffic applications requiring horizontal scaling. Consider these factors when choosing replica count: Available cluster resources (CPU, memory, storage), Application resource requirements per replica, Load balancing and traffic distribution needs, Budget and cost constraints, Disaster recovery and availability requirements. For StatefulSets, scaling up creates new instances with persistent identity; scaling down removes the highest-numbered instances first. For Deployments, scaling triggers rolling updates if the pod template has changed. Monitor cluster resource usage after scaling to ensure sufficient capacity. The operation may fail if there are insufficient cluster resources to support the requested replica count.")),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive debug output for troubleshooting the scaling operation and understanding the process details. Set to 'true' to see detailed information including: Kubernetes API calls and responses, current resource state before scaling, scaling operation progress and status, any errors or warnings during the scaling process, resource validation and permission checks, timing information for the scaling operation. Set to 'false' or omit for normal output showing only the scaling result and final status. Debug mode is particularly helpful when: scaling operations fail or behave unexpectedly, you need to understand why scaling is taking longer than expected, there are resource constraints or quota limitations, you're troubleshooting RBAC permission issues, or when you want to monitor the scaling process in detail for learning or automation purposes. Debug output may include sensitive cluster information, so use carefully in production environments.")),
	)
}

// GetAPIVersionsTool retrieves available Kubernetes API versions
func GetAPIVersionsTool() mcp.Tool {
	logrus.Debug("Creating GetAPIVersionsTool")
	return mcp.NewTool("kubernetes_get_api_versions",
		mcp.WithDescription("Retrieve all available Kubernetes API versions supported by the cluster, similar to 'kubectl api-versions'. This tool discovers which API versions are available for creating and managing resources in your specific Kubernetes cluster. Different clusters may support different API versions depending on the Kubernetes version and installed extensions. Use this tool when you need to: verify which API versions are supported before creating resources, troubleshoot API compatibility issues, check for deprecated API versions that may affect your resources, discover available API groups and their versions, or understand the API capabilities of your cluster. The output shows both core API versions (like v1) and extension API versions (like apps/v1, networking.k8s.io/v1). This information is essential for writing correct resource manifests and understanding which features are available in your cluster environment."),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting API discovery issues. Set to 'true' to see detailed information about the API server communication, discovery process, and any errors encountered. Set to 'false' or omit for normal output showing only the available API versions. Debug mode helps identify connectivity issues, authentication problems, or API server configuration issues that might prevent proper API version discovery.")),
	)
}

// GetAPIResourcesTool retrieves available Kubernetes API resources
func GetAPIResourcesTool() mcp.Tool {
	logrus.Debug("Creating GetAPIResourcesTool")
	return mcp.NewTool("kubernetes_get_api_resources",
		mcp.WithDescription("Retrieve comprehensive information about all available Kubernetes API resources (kinds) supported by the cluster, similar to 'kubectl api-resources'. This tool provides detailed information about each resource type including their names, API groups, whether they are namespaced, and their short names. Use this tool when you need to: discover what resource types are available in your cluster, find the correct API group and version for specific resources, determine whether a resource is namespaced or cluster-scoped, find short names for resources to use in kubectl commands, explore custom resources installed by operators or CRDs, troubleshoot resource creation issues by verifying resource availability, or understand the complete API surface of your Kubernetes cluster. The output includes core resources (Pod, Service, etc.) and any custom resources defined by Custom Resource Definitions (CRDs) or API extensions."),
		mcp.WithString("apiGroup",
			mcp.Description("Filter results to show only resources from a specific API group. This allows you to focus on resources from particular API groups rather than seeing all available resources. Common API groups include: '' or 'core' (core resources like Pod, Service, ConfigMap, Secret, Namespace), 'apps' (application workloads like Deployment, StatefulSet, DaemonSet), 'batch' (batch workloads like Job, CronJob), 'networking.k8s.io' (networking resources like Ingress, NetworkPolicy), 'rbac.authorization.k8s.io' (RBAC resources like Role, ClusterRole), 'storage.k8s.io' (storage resources like StorageClass), 'autoscaling' (scaling resources like HorizontalPodAutoscaler), 'metrics.k8s.io' (metrics resources), or custom API groups from installed operators. Leave empty to show resources from all API groups. Use this filter when you're working with specific APIs or need to understand what resources are available in a particular domain.")),
		mcp.WithBoolean("namespaced",
			mcp.Description("Filter results to show only namespaced resources (true) or only cluster-scoped resources (false). Leave unset to show both types. Namespaced resources (like Pod, Service, Deployment, ConfigMap, Secret) exist within specific namespaces and are isolated from resources in other namespaces. Cluster-scoped resources (like Node, ClusterRole, PersistentVolume, Namespace itself) exist at the cluster level and are not confined to any namespace. Set to true when you need to find resources that can be created within namespaces, or false when looking for cluster-wide resources. This filter is particularly useful for understanding resource scope and planning RBAC policies, as namespaced and cluster-scoped resources often require different permission configurations.")),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive debug output for troubleshooting API resource discovery and understanding the cluster's API capabilities. Set to 'true' to see detailed information about: API server communication and discovery requests, processing of API groups and versions, resource metadata extraction and formatting, any errors or warnings during discovery, timing information for API calls. Set to 'false' or omit for normal output showing only the API resource information. Debug mode is helpful when: resource discovery fails or returns unexpected results, you're troubleshooting API server connectivity, you need to understand which API extensions are installed, or when debugging issues with custom resources or CRDs. Debug output may include sensitive cluster configuration details.")),
	)
}

// GetResourcesDetailTool retrieves detailed information for multiple resources efficiently
func GetResourcesDetailTool() mcp.Tool {
	logrus.Debug("Creating GetResourcesDetailTool")
	return mcp.NewTool("kubernetes_get_resources_detail",
		mcp.WithDescription("🔍 Efficiently retrieve detailed information for multiple specific Kubernetes resources in a single request. This tool is optimized for getting comprehensive details about multiple resources without overwhelming the context. Use this when you need detailed information about several specific resources identified from a previous list operation. The tool supports batch retrieval with pagination awareness and provides full resource details including configuration, status, events, and relationships. Ideal for: investigating specific resources identified in a list, getting complete details for a subset of resources, performing detailed analysis on selected items, or preparing for resource modification operations."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Resource kind/type - must be the same for all resources in this request (e.g., 'Pod', 'Deployment', 'Service'). Use exact case-sensitive names as they appear in Kubernetes API. This constraint ensures efficient API batching while preventing resource type mixing that could cause context overflow.")),
		mcp.WithArray("names", mcp.Required(),
			mcp.Description("Array of exact resource names to retrieve detailed information for. All resources must be of the same kind specified in the 'kind' parameter. Names are case-sensitive and must match the metadata.name field exactly. Use this to get detailed info for multiple resources efficiently in one API call. Recommended to limit to 10-20 resources per request to avoid context overflow while still being more efficient than individual resource calls.")),
		mcp.WithString("namespace",
			mcp.Description("Required namespace for namespaced resources (Pod, Service, Deployment, ConfigMap, Secret, etc.). All resources must be in the same namespace. Omit only for cluster-scoped resources (Node, PersistentVolume, ClusterRole, etc.). This constraint enables efficient namespace-scoped API operations.")),
		mcp.WithBoolean("includeEvents",
			mcp.Description("Include related events for each resource (default: false). When set to true, retrieves recent events related to each resource, providing valuable context for troubleshooting. Events are automatically filtered and limited to prevent excessive output. Useful for understanding resource status, troubleshooting issues, or identifying recent changes affecting the resources.")),
		mcp.WithBoolean("includeStatus",
			mcp.Description("Include detailed status information (default: true). When false, focuses primarily on configuration and metadata, reducing output size. Status information includes conditions, readiness states, and runtime details that are essential for understanding resource state but can increase response size significantly.")),
		mcp.WithString("debug",
			mcp.Description("Enable detailed debug output for troubleshooting the batch retrieval operation. Set to 'true' to see information about individual API calls, caching behavior, resource processing, and any issues encountered. Set to 'false' or omit for normal output showing only the resource details. Debug mode helps understand the efficiency improvements and identify any bottlenecks in the batch operation.")),
	)
}

// GetEventsDetailTool retrieves detailed events with full information
func GetEventsDetailTool() mcp.Tool {
	logrus.Debug("Creating GetEventsDetailTool")
	return mcp.NewTool("kubernetes_get_events_detail",
		mcp.WithDescription("📋 Retrieve comprehensive Kubernetes events with complete information for thorough analysis and troubleshooting. This tool provides full event details including all fields, timestamps, and messages without optimization for context size. Use this when you need complete event information for detailed investigation, audit purposes, or comprehensive troubleshooting. For quick checks, use the optimized get_recent_events tool instead. This tool includes pagination support to handle large event volumes while maintaining full detail preservation."),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to filter events from. If not specified, events from all namespaces will be returned (requires cluster-wide permissions). For focused troubleshooting, specify the namespace where the problematic resources are located. Common namespaces include 'default', 'kube-system' (for cluster components), 'kube-public', or custom application namespaces.")),
		mcp.WithString("fieldSelector",
			mcp.Description("Field selector to filter events based on specific criteria. This allows precise filtering of events related to specific resources or conditions. Common examples: 'involvedObject.name=my-pod' (events for a specific pod), 'involvedObject.kind=Pod' (all pod-related events), 'type=Warning' (only warning events), 'type=Normal' (only normal events), 'reason=Failed' (events with Failed reason). You can combine multiple selectors with commas.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of events to return (default: 50, max: 200). This tool returns full event details, so be conservative with limits to avoid context overflow. Use higher limits only when you need comprehensive event analysis. Default is optimized for detailed investigation without overwhelming context.")),
		mcp.WithString("continueToken",
			mcp.Description("Pagination token from a previous response to fetch the next page of events. When the response indicates 'hasMore': true, use the provided 'continueToken' to get the next batch of events. This enables efficient traversal through large event sets while maintaining full detail.")),
		mcp.WithBoolean("includeNormalEvents",
			mcp.Description("Include normal operational events (default: false). When set to false, only returns warning and error events for focused troubleshooting. When set to true, includes all events including normal operational events. Normal events are typically high-volume and less critical for troubleshooting.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the events retrieval operation. Set to 'true' to see detailed information about API calls, filtering, pagination, and any issues encountered. Set to 'false' or omit for normal output showing only the events.")),
	)
}

// ListResourcesFullTool lists resources with complete details (non-optimized)
func ListResourcesFullTool() mcp.Tool {
	logrus.Debug("Creating ListResourcesFullTool")
	return mcp.NewTool("kubernetes_list_resources_full",
		mcp.WithDescription("📄 List Kubernetes resources with complete details and full object information. This tool returns the entire resource configuration without any optimization for context size. Use this sparingly when you need complete resource details for analysis, backup, or comprehensive configuration review. For most use cases, prefer the summary tools or use pagination with reasonable limits. This tool is ideal for: exporting resource configurations, performing detailed compliance checks, creating backups, or when you need full YAML/JSON representation of multiple resources."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind/type to list with full details (e.g., 'Pod', 'Deployment', 'Service'). Use exact case-sensitive names as they appear in Kubernetes API. This tool will return complete objects for all matching resources, so use with caution.")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scope the resource listing. For namespaced resources, this filters results to only show resources within the specified namespace. If omitted for namespaced resources, shows resources from ALL namespaces (requires cluster-wide permissions). For cluster-scoped resources, this parameter is ignored.")),
		mcp.WithString("labelSelector",
			mcp.Description("Label selector to filter resources based on their metadata labels. Use this powerful filtering mechanism to find resources matching specific criteria. Syntax: 'app=nginx', 'env=production', or 'app=nginx,env=prod' for multiple labels.")),
		mcp.WithString("fieldSelector",
			mcp.Description("Field selector to filter resources based on specific field values. Provides more precise filtering compared to label selectors. Common examples: 'status.phase=Running' (running pods), 'status.phase!=Failed' (exclude failed resources).")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of resources to return (default: 10, max: 50). This tool returns full resource details, so limits are conservative to prevent context overflow. Default limit is intentionally low. Increase only when you specifically need more full resources.")),
		mcp.WithString("continueToken",
			mcp.Description("Pagination token from previous response to fetch the next page of results. When response indicates 'hasMore': true, use the provided 'continueToken' to get the next batch of full resources.")),
		mcp.WithBoolean("includeStatus",
			mcp.Description("Include detailed status information (default: true). When false, reduces output size by excluding runtime status fields while keeping configuration. Useful for configuration-focused analysis.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the full resource listing operation. Set to 'true' to see detailed API information, processing steps, and any issues. Set to 'false' or omit for normal output.")),
	)
}

// GetResourceDetailAdvancedTool retrieves comprehensive detailed information with enhanced formatting
func GetResourceDetailAdvancedTool() mcp.Tool {
	logrus.Debug("Creating GetResourceDetailAdvancedTool")
	return mcp.NewTool("kubernetes_get_resource_detail_advanced",
		mcp.WithDescription("🔍 Advanced resource detail retrieval with enhanced formatting and optional components. Get comprehensive information about a specific Kubernetes resource with configuration, status, events, relationships, and diagnostics. This is the ultimate detail tool for thorough analysis and troubleshooting. ⚠️ NOTE: This tool returns detailed information and may produce large output. Recommend using summary tools first to confirm resource scope, then use this tool for deep analysis."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Resource kind/type - exact case-sensitive name as in Kubernetes API (e.g., 'Pod', 'Deployment', 'Service', 'ConfigMap', etc.).")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact resource name from metadata.name field.")),
		mcp.WithString("namespace",
			mcp.Description("Required for namespaced resources (Pod, Service, Deployment, etc.). Omit for cluster-scoped resources (Node, PersistentVolume, etc.).")),
		mcp.WithBoolean("includeEvents",
			mcp.Description("Include related events for context and troubleshooting (default: false). Events help understand what happened to the resource.")),
		mcp.WithBoolean("includeRelationships",
			mcp.Description("Include owner/dependent relationships (default: false). Shows what this resource depends on or what depends on it.")),
		mcp.WithBoolean("includeDiagnostics",
			mcp.Description("Include diagnostic information and health checks (default: false). Provides additional insights for troubleshooting.")),
		mcp.WithBoolean("includeConfiguration",
			mcp.Description("Include full configuration details (default: true). When false, focuses on status and metadata only.")),
		mcp.WithString("outputFormat",
			mcp.Description("Output format preference: 'compact', 'structured', or 'verbose' (default: 'structured'). Controls detail level and organization.")),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive debug output for troubleshooting the detail retrieval process (true/false).")),
	)
}

// ============ Troubleshooting Tools ============

// GetUnhealthyResourcesTool finds pods and resources in unhealthy states
func GetUnhealthyResourcesTool() mcp.Tool {
	logrus.Debug("Creating GetUnhealthyResourcesTool")
	return mcp.NewTool("kubernetes_get_unhealthy_resources",
		mcp.WithDescription("Find Kubernetes resources in unhealthy states (crash, pending, failed, etc.)"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to scan. Empty = all namespaces")),
		mcp.WithArray("resourceTypes",
			mcp.Description("Resource types to check (Pod, Job, Deployment, StatefulSet, DaemonSet). Default: all")),
	)
}

// GetNodeConditionsTool retrieves node conditions and health status
func GetNodeConditionsTool() mcp.Tool {
	logrus.Debug("Creating GetNodeConditionsTool")
	return mcp.NewTool("kubernetes_get_node_conditions",
		mcp.WithDescription("Get detailed node conditions (Ready, MemoryPressure, DiskPressure, PIDPressure, etc.)"),
		mcp.WithString("nodeName", mcp.Required(),
			mcp.Description("Exact node name to get conditions for")),
	)
}

// AnalyzeIssueTool performs AI-powered issue analysis
func AnalyzeIssueTool() mcp.Tool {
	logrus.Debug("Creating AnalyzeIssueTool")
	return mcp.NewTool("kubernetes_analyze_issue",
		mcp.WithDescription("AI-powered Kubernetes resource issue analysis with recommendations"),
		mcp.WithString("issueType", mcp.Required(),
			mcp.Description("Issue type: pod_crash, pod_pending, deployment_unavailable, job_failed")),
		mcp.WithString("resourceKind", mcp.Required(),
			mcp.Description("Resource kind (Pod, Deployment, Job, etc.)")),
		mcp.WithString("resourceName", mcp.Required(),
			mcp.Description("Resource name to analyze")),
		mcp.WithString("namespace",
			mcp.Description("Resource namespace (required for namespaced resources)")),
	)
}

// ============ Search Tools ============

// SearchResourcesTool searches for Kubernetes resources by name with fuzzy matching
func SearchResourcesTool() mcp.Tool {
	logrus.Debug("Creating SearchResourcesTool")
	return mcp.NewTool("kubernetes_search_resources",
		mcp.WithDescription("🔍 Search for Kubernetes resources by name with fuzzy matching. This tool helps you find resources when you only remember part of the name or want to discover resources matching a pattern. Supports multiple search strategies including contains, startsWith, endsWith, and regex. Perfect for resource discovery, troubleshooting, and cluster exploration. 🎯 Best for: resource discovery, fuzzy name matching, quick resource location. 📋 Use cases: forgot complete resource name, finding resources with specific patterns, batch resource management."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind/type to search (e.g., 'Pod', 'Deployment', 'Service', 'ConfigMap', 'Secret', 'Ingress', 'StatefulSet', 'DaemonSet', 'Job', 'CronJob'). Use exact case-sensitive names as they appear in Kubernetes API.")),
		mcp.WithString("query", mcp.Required(),
			mcp.Description("Search query string to match against resource names. The search will find resources whose names contain this string. For example, 'nginx' will find 'nginx-deployment', 'nginx-pod-123', 'my-nginx-service', etc. The search is case-insensitive by default.")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scope the search. For namespaced resources, this filters results to only show resources within the specified namespace. If omitted, searches across ALL namespaces (requires cluster-wide permissions). For cluster-scoped resources, this parameter is ignored.")),
		mcp.WithString("searchMode",
			mcp.Description("Search strategy to use (default: 'contains'). Options: 'contains' (name contains the query), 'startsWith' (name starts with the query), 'endsWith' (name ends with the query), 'exact' (exact match), 'regex' (regular expression pattern). Use 'contains' for broad matching, 'startsWith' for prefix-based filtering, 'endsWith' for suffix-based filtering, 'exact' for precise matching, and 'regex' for complex patterns.")),
		mcp.WithBoolean("caseSensitive",
			mcp.Description("Whether the search should be case-sensitive (default: false). When false, the search ignores case differences, making it easier to find resources. When true, matches must respect the exact case of the query string.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of matching resources to return (default: 50, max: 200). This controls the size of the result set. Use smaller limits (10-20) for quick searches, larger limits (50-200) for comprehensive discovery.")),
		mcp.WithString("labelSelector",
			mcp.Description("Optional label selector to further filter search results. Use this to combine name-based search with label-based filtering. Syntax: 'app=nginx', 'env=production', or 'app=nginx,env=prod' for multiple labels.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the search operation (true/false).")),
	)
}
