package handlers

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// ToolRecommendation intelligent tool recommendation system
type ToolRecommendation struct {
	PrimaryTools     []string               // Recommended primary tools
	AlternativeTools []string               // Alternative tools
	Reason           string                 // Recommendation reason
	Context          map[string]interface{} // Context information
	Tips             []string               // Usage tips
	EstimatedSize    string                 // Estimated response size
}

// ToolAdvisor intelligent tool advisor
type ToolAdvisor struct {
	toolPatterns map[string]ToolProfile
}

// ToolProfile tool profile
type ToolProfile struct {
	Name          string
	Category      string   // summary, standard, detail
	DataSize      string   // small, medium, large
	UseCase       string   // monitoring, troubleshooting, inventory
	CommonParams  []string // Common parameters
	Alternatives  []string // Alternative tools
	Prerequisites []string // Prerequisites
}

// NewToolAdvisor creates new tool advisor
func NewToolAdvisor() *ToolAdvisor {
	return &ToolAdvisor{
		toolPatterns: map[string]ToolProfile{
			"kubernetes_list_resources": {
				Name:          "kubernetes_list_resources",
				Category:      "standard",
				DataSize:      "medium",
				UseCase:       "inventory",
				CommonParams:  []string{"kind", "namespace", "labelSelector", "limit"},
				Alternatives:  []string{"kubernetes_list_resources_summary", "kubernetes_list_resources_full"},
				Prerequisites: []string{},
			},
			"kubernetes_list_resources_summary": {
				Name:          "kubernetes_list_resources_summary",
				Category:      "summary",
				DataSize:      "small",
				UseCase:       "inventory",
				CommonParams:  []string{"kind", "namespace", "labelSelector", "limit"},
				Alternatives:  []string{"kubernetes_list_resources"},
				Prerequisites: []string{},
			},
			"kubernetes_list_resources_full": {
				Name:          "kubernetes_list_resources_full",
				Category:      "detail",
				DataSize:      "large",
				UseCase:       "troubleshooting",
				CommonParams:  []string{"kind", "namespace", "labelSelector", "limit"},
				Alternatives:  []string{"kubernetes_list_resources"},
				Prerequisites: []string{},
			},
			"kubernetes_get_events": {
				Name:          "kubernetes_get_events",
				Category:      "standard",
				DataSize:      "medium",
				UseCase:       "troubleshooting",
				CommonParams:  []string{"namespace", "fieldSelector", "limit"},
				Alternatives:  []string{"kubernetes_get_recent_events", "kubernetes_get_events_detail"},
				Prerequisites: []string{},
			},
			"kubernetes_get_recent_events": {
				Name:          "kubernetes_get_recent_events",
				Category:      "summary",
				DataSize:      "small",
				UseCase:       "monitoring",
				CommonParams:  []string{"namespace", "limit"},
				Alternatives:  []string{"kubernetes_get_events"},
				Prerequisites: []string{},
			},
			"kubernetes_get_events_detail": {
				Name:          "kubernetes_get_events_detail",
				Category:      "detail",
				DataSize:      "large",
				UseCase:       "troubleshooting",
				CommonParams:  []string{"namespace", "fieldSelector", "limit"},
				Alternatives:  []string{"kubernetes_get_events"},
				Prerequisites: []string{},
			},
			"kubernetes_get_pod_logs": {
				Name:          "kubernetes_get_pod_logs",
				Category:      "detail",
				DataSize:      "medium",
				UseCase:       "troubleshooting",
				CommonParams:  []string{"name", "namespace", "container", "tailLines"},
				Alternatives:  []string{"kubernetes_describe_resource"},
				Prerequisites: []string{"pod name"},
			},
			"kubernetes_describe_resource": {
				Name:          "kubernetes_describe_resource",
				Category:      "detail",
				DataSize:      "large",
				UseCase:       "troubleshooting",
				CommonParams:  []string{"kind", "name", "namespace"},
				Alternatives:  []string{"kubernetes_get_resource_summary"},
				Prerequisites: []string{"resource name"},
			},
			"kubernetes_get_resource_summary": {
				Name:          "kubernetes_get_resource_summary",
				Category:      "summary",
				DataSize:      "small",
				UseCase:       "monitoring",
				CommonParams:  []string{"kind", "name", "namespace"},
				Alternatives:  []string{"kubernetes_get_resource"},
				Prerequisites: []string{"resource name"},
			},
			"kubernetes_get_resources_detail": {
				Name:          "kubernetes_get_resources_detail",
				Category:      "detail",
				DataSize:      "large",
				UseCase:       "troubleshooting",
				CommonParams:  []string{"kind", "names", "namespace"},
				Alternatives:  []string{"kubernetes_list_resources"},
				Prerequisites: []string{"resource names"},
			},
		},
	}
}

// RecommendTools recommends tools based on usage scenario
func (advisor *ToolAdvisor) RecommendTools(scenario string, params map[string]interface{}) ToolRecommendation {
	scenario = strings.ToLower(strings.TrimSpace(scenario))

	logrus.WithFields(logrus.Fields{
		"scenario": scenario,
		"params":   params,
	}).Debug("Generating tool recommendations")

	switch scenario {
	case "pod_troubleshooting", "debug_pod":
		return advisor.recommendForPodTroubleshooting(params)
	case "health_check", "quick_check":
		return advisor.recommendForHealthCheck(params)
	case "inventory", "list":
		return advisor.recommendForInventory(params)
	case "event_analysis", "debug_events":
		return advisor.recommendForEventAnalysis(params)
	case "log_analysis":
		return advisor.recommendForLogAnalysis(params)
	case "resource_monitoring":
		return advisor.recommendForResourceMonitoring(params)
	case "full_analysis", "comprehensive":
		return advisor.recommendForFullAnalysis(params)
	default:
		return advisor.recommendForGeneralUse(params)
	}
}

// recommendForPodTroubleshooting Pod troubleshooting recommendation
func (advisor *ToolAdvisor) recommendForPodTroubleshooting(params map[string]interface{}) ToolRecommendation {
	tools := []string{}
	tips := []string{
		"Start with summary tools to quickly understand resource status",
		"Then get detailed information for in-depth analysis",
		"Finally check related events and logs for troubleshooting",
	}

	// Check if there are specific resources
	if namespace, ok := params["namespace"].(string); ok && namespace != "" {
		tools = append(tools, "kubernetes_list_resources_summary -kind='Pod' -namespace='"+namespace+"' -limit=20")
		if podName, ok := params["podName"].(string); ok && podName != "" {
			tools = append(tools, "kubernetes_get_resource_summary -kind='Pod' -name='"+podName+"' -namespace='"+namespace+"'")
			tools = append(tools, "kubernetes_get_recent_events -namespace='"+namespace+"' -fieldSelector='involvedObject.name="+podName+"'")
			tools = append(tools, "kubernetes_get_pod_logs -name='"+podName+"' -namespace='"+namespace+"' -tailLines=100")
		}
	} else {
		tools = append(tools, "kubernetes_list_resources_summary -kind='Pod' -limit=20")
		tools = append(tools, "kubernetes_get_recent_events -limit=30")
	}

	return ToolRecommendation{
		PrimaryTools:     tools,
		AlternativeTools: []string{"kubernetes_list_resources -kind='Pod' -limit=10", "kubernetes_describe_resource -kind='Pod' -name='<pod-name>'"},
		Reason:           "Pod troubleshooting workflow: from overview to details, from status to events and logs",
		Context: map[string]interface{}{
			"scenario": "pod_troubleshooting",
			"workflow": "summary -> details -> events -> logs",
		},
		Tips:          tips,
		EstimatedSize: "medium",
	}
}

// recommendForHealthCheck health check recommendation
func (advisor *ToolAdvisor) recommendForHealthCheck(params map[string]interface{}) ToolRecommendation {
	tools := []string{
		"kubernetes_list_resources_summary -kind='Pod' -limit=20",
		"kubernetes_list_resources_summary -kind='Deployment' -limit=10",
		"kubernetes_get_recent_events -limit=20",
	}

	tips := []string{
		"Use summary tools to check overall cluster health status",
		"Focus on abnormal Pods and critical events",
		"If problems found, get detailed information further",
	}

	return ToolRecommendation{
		PrimaryTools:     tools,
		AlternativeTools: []string{},
		Reason:           "Quick health check: key resource summary and important events",
		Context: map[string]interface{}{
			"scenario": "health_check",
			"focus":    "critical_resources",
		},
		Tips:          tips,
		EstimatedSize: "small",
	}
}

// recommendForInventory inventory recommendation
func (advisor *ToolAdvisor) recommendForInventory(params map[string]interface{}) ToolRecommendation {
	tools := []string{}

	// Check specified resource type
	if kind, ok := params["kind"].(string); ok && kind != "" {
		tools = append(tools, fmt.Sprintf("kubernetes_list_resources_summary -kind='%s' -limit=50", kind))
		if namespace, ok := params["namespace"].(string); ok && namespace != "" {
			tools = append(tools, fmt.Sprintf("kubernetes_list_resources_summary -kind='%s' -namespace='%s' -limit=50", kind, namespace))
		}
	} else {
		tools = append(tools, "kubernetes_list_resources_summary -kind='Pod' -limit=30")
		tools = append(tools, "kubernetes_list_resources_summary -kind='Deployment' -limit=20")
		tools = append(tools, "kubernetes_list_resources_summary -kind='Service' -limit=20")
	}

	tips := []string{
		"Use summary tools for efficient resource inventory",
		"Filter results by namespace and labelSelector",
		"If complete configuration needed, use detailed tool versions",
	}

	return ToolRecommendation{
		PrimaryTools:     tools,
		AlternativeTools: []string{"kubernetes_list_resources -kind='Pod' -limit=20", "kubernetes_list_resources_full -kind='Pod' -limit=5"},
		Reason:           "Resource inventory: summary first, view details on demand",
		Context: map[string]interface{}{
			"scenario": "inventory",
			"approach": "summary_first",
		},
		Tips:          tips,
		EstimatedSize: "small",
	}
}

// recommendForEventAnalysis event analysis recommendation
func (advisor *ToolAdvisor) recommendForEventAnalysis(params map[string]interface{}) ToolRecommendation {
	tools := []string{"kubernetes_get_recent_events -limit=30"}

	if namespace, ok := params["namespace"].(string); ok && namespace != "" {
		tools = append(tools, fmt.Sprintf("kubernetes_get_events -namespace='%s' -limit=50", namespace))
	} else {
		tools = append(tools, "kubernetes_get_events -limit=50")
	}

	if resourceType, ok := params["resourceType"].(string); ok && resourceType != "" {
		tools = append(tools, fmt.Sprintf("kubernetes_list_resources_summary -kind='%s' -limit=20", resourceType))
	}

	tips := []string{
		"Start analysis from recent critical events",
		"Use filters to focus on specific resources or event types",
		"View detailed event information when necessary",
	}

	return ToolRecommendation{
		PrimaryTools:     tools,
		AlternativeTools: []string{"kubernetes_get_events_detail -limit=100"},
		Reason:           "Event analysis: from summary to details, focus on critical anomalies",
		Context: map[string]interface{}{
			"scenario": "event_analysis",
			"focus":    "problematic_events",
		},
		Tips:          tips,
		EstimatedSize: "medium",
	}
}

// recommendForLogAnalysis log analysis recommendation
func (advisor *ToolAdvisor) recommendForLogAnalysis(params map[string]interface{}) ToolRecommendation {
	tips := []string{
		"Limit log lines to avoid context overflow",
		"Use tailLines parameter to control log size",
		"Combine with event analysis for comprehensive diagnosis",
	}

	if podName, ok := params["podName"].(string); ok && podName != "" {
		namespace, _ := params["namespace"].(string)
		if namespace == "" {
			namespace = "default"
		}

		tools := []string{
			fmt.Sprintf("kubernetes_get_resource_summary -kind='Pod' -name='%s' -namespace='%s'", podName, namespace),
			fmt.Sprintf("kubernetes_get_pod_logs -name='%s' -namespace='%s' -tailLines=100", podName, namespace),
			fmt.Sprintf("kubernetes_get_recent_events -namespace='%s' -fieldSelector='involvedObject.name=%s'", namespace, podName),
		}

		return ToolRecommendation{
			PrimaryTools:     tools,
			AlternativeTools: []string{"kubernetes_describe_resource -kind='Pod' -name='" + podName + "' -namespace='" + namespace + "'"},
			Reason:           "Log analysis: analyze logs and related events after checking Pod status",
			Context: map[string]interface{}{
				"scenario": "log_analysis",
				"focus":    "specific_pod",
			},
			Tips:          tips,
			EstimatedSize: "medium",
		}
	}

	return ToolRecommendation{
		PrimaryTools:     []string{},
		AlternativeTools: []string{"kubernetes_get_pod_logs -name='<pod-name>' -namespace='<namespace>'"},
		Reason:           "Need to specify Pod name for log analysis",
		Context: map[string]interface{}{
			"scenario": "log_analysis",
			"requires": "pod_name",
		},
		Tips:          append(tips, "Please provide Pod name to get specific logs"),
		EstimatedSize: "medium",
	}
}

// recommendForResourceMonitoring resource monitoring recommendation
func (advisor *ToolAdvisor) recommendForResourceMonitoring(params map[string]interface{}) ToolRecommendation {
	tools := []string{
		"kubernetes_list_resources_summary -kind='Pod' -limit=30",
		"kubernetes_list_resources_summary -kind='Deployment' -limit=20",
		"kubernetes_get_recent_events -limit=15",
	}
	alternatives := []string{
		"kubernetes_get_resource_usage -resourceType='node'",
		"kubernetes_get_resource_usage -resourceType='pod' -namespace='<namespace>'",
	}

	tips := []string{
		"Monitor resource status changes and error events",
		"Use summary tools to periodically check resource health status",
		"Get detailed status through specific resource names",
	}

	return ToolRecommendation{
		PrimaryTools:     tools,
		AlternativeTools: alternatives,
		Reason:           "Resource monitoring: periodically check key resource status and abnormal events",
		Context: map[string]interface{}{
			"scenario": "resource_monitoring",
			"approach": "proactive_monitoring",
		},
		Tips:          tips,
		EstimatedSize: "small",
	}
}

// recommendForFullAnalysis comprehensive analysis recommendation
func (advisor *ToolAdvisor) recommendForFullAnalysis(params map[string]interface{}) ToolRecommendation {
	tools := []string{
		"kubernetes_list_resources_summary -kind='Pod' -limit=20",
		"kubernetes_list_resources_summary -kind='Deployment' -limit=10",
		"kubernetes_list_resources_summary -kind='Service' -limit=10",
		"kubernetes_get_recent_events -limit=30",
	}

	alternatives := []string{
		"kubernetes_get_events -limit=50",
		"kubernetes_list_resources -kind='Pod' -limit=10",
		"kubernetes_get_resource_usage -resourceType='node'",
	}

	tips := []string{
		"Comprehensive analysis requires multiple steps: first summary, then details, check logs and events when necessary",
		"Use pagination mechanism to gradually obtain large amounts of data",
		"Combine different tool perspectives for comprehensive diagnosis",
	}

	return ToolRecommendation{
		PrimaryTools:     tools,
		AlternativeTools: alternatives,
		Reason:           "Comprehensive analysis: multi-tool combination, hierarchical cluster status check",
		Context: map[string]interface{}{
			"scenario": "full_analysis",
			"approach": "multi_tool",
		},
		Tips:          tips,
		EstimatedSize: "large",
	}
}

// recommendForGeneralUse general scenario recommendation
func (advisor *ToolAdvisor) recommendForGeneralUse(params map[string]interface{}) ToolRecommendation {
	tools := []string{
		"kubernetes_list_resources_summary -kind='Pod' -limit=20",
		"kubernetes_get_recent_events -limit=20",
	}

	tips := []string{
		"Prioritize summary tools to avoid context overflow",
		"Get detailed resource information as needed",
		"Maintain appropriate request limits to ensure performance",
	}

	return ToolRecommendation{
		PrimaryTools:     tools,
		AlternativeTools: []string{"kubernetes_list_resources -kind='Pod' -limit=10"},
		Reason:           "General scenario: balance information volume and context safety",
		Context: map[string]interface{}{
			"scenario": "general",
			"approach": "balanced",
		},
		Tips:          tips,
		EstimatedSize: "small",
	}
}

// GetOptimalLimit recommends optimal parameters based on tool type and scenario
func (advisor *ToolAdvisor) GetOptimalLimit(toolName string, scenario string) int {
	profile, exists := advisor.toolPatterns[toolName]
	if !exists {
		return 20 // Default conservative limit
	}

	// Adjust based on tool type and scenario
	switch profile.Category {
	case "summary":
		if scenario == "inventory" || scenario == "health_check" {
			return 50 // Summary tools can be appropriately higher
		}
		return 30
	case "standard":
		return 20 // Standard tools keep moderate
	case "detail":
		return 5 // Detail tools keep very small limits
	default:
		return 20
	}
}

// AnalyzeContext analyzes user request and provides intelligent recommendations
func (advisor *ToolAdvisor) AnalyzeContext(requestText string, detectedParams map[string]interface{}) ToolRecommendation {
	requestText = strings.ToLower(requestText)

	logrus.WithFields(logrus.Fields{
		"requestText":    requestText,
		"detectedParams": detectedParams,
	}).Debug("Analyzing user context")

	// Detect intent
	scenarios := map[string][]string{
		"pod_troubleshooting": {"pod", "problem", "error", "crash", "pending", "failed"},
		"health_check":        {"status", "health", "check", "overview", "summary"},
		"inventory":           {"list", "all", "show", "inventory", "what", "discover"},
		"event_analysis":      {"event", "logs", "history", "timeline"},
		"log_analysis":        {"log", "output", "message", "debug"},
		"resource_monitoring": {"monitor", "watch", "resource", "usage"},
		"full_analysis":       {"full", "complete", "comprehensive", "detailed"},
	}

	bestMatch := "general"
	bestScore := 0

	for scenario, keywords := range scenarios {
		score := 0
		for _, keyword := range keywords {
			if strings.Contains(requestText, keyword) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			bestMatch = scenario
		}
	}

	return advisor.RecommendTools(bestMatch, detectedParams)
}

// Global tool advisor instance
var DefaultToolAdvisor = NewToolAdvisor()
