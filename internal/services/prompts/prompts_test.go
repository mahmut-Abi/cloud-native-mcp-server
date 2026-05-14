package prompts

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestRegistrationsHaveUniqueNames(t *testing.T) {
	seen := map[string]struct{}{}
	for _, registration := range Registrations() {
		name := registration.Prompt.Name
		if name == "" {
			t.Fatal("prompt name must not be empty")
		}
		if _, exists := seen[name]; exists {
			t.Fatalf("duplicate prompt name: %s", name)
		}
		seen[name] = struct{}{}
		if registration.Handler == nil {
			t.Fatalf("prompt %s has nil handler", name)
		}
	}
}

func TestRequiredServices(t *testing.T) {
	if got := RequiredServices(PromptName); len(got) != 0 {
		t.Fatalf("test prompt should not require services, got %v", got)
	}
	if got := RequiredServices(K8sOpsPromptName); len(got) != 1 || got[0] != "kubernetes" {
		t.Fatalf("unexpected required services for k8s prompt: %v", got)
	}
	if got := RequiredServices(PrometheusDiagnosisPromptName); len(got) != 1 || got[0] != "prometheus" {
		t.Fatalf("unexpected required services for prometheus prompt: %v", got)
	}
	if got := RequiredServices(MultiServiceRCAPromptName); len(got) != 0 {
		t.Fatalf("unexpected required services for multi-service RCA prompt: %v", got)
	}
}

func TestIsAvailable(t *testing.T) {
	isEnabled := func(name string) bool {
		switch name {
		case "kubernetes", "prometheus":
			return true
		default:
			return false
		}
	}

	if !IsAvailable(K8sOpsPromptName, isEnabled) {
		t.Fatalf("expected kubernetes prompt to be available")
	}
	if !IsAvailable(PrometheusDiagnosisPromptName, isEnabled) {
		t.Fatalf("expected prometheus prompt to be available")
	}
	if IsAvailable(ArgoCDDiagnosisPromptName, isEnabled) {
		t.Fatalf("expected argocd prompt to be unavailable")
	}
}

func TestHandleIncidentTriagePrompt(t *testing.T) {
	result, err := HandleIncidentTriagePrompt(context.Background(), mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Name: IncidentTriagePromptName,
			Arguments: map[string]string{
				"symptom":   "503 spikes",
				"namespace": "prod",
				"workload":  "api",
			},
		},
	})
	if err != nil {
		t.Fatalf("HandleIncidentTriagePrompt returned error: %v", err)
	}
	if result == nil || len(result.Messages) == 0 {
		t.Fatalf("expected prompt messages, got %#v", result)
	}
	text, ok := result.Messages[0].Content.(mcp.TextContent)
	if !ok {
		t.Fatalf("expected text content, got %#v", result.Messages[0].Content)
	}
	if text.Text == "" {
		t.Fatal("expected non-empty prompt text")
	}
	if !contains(text.Text, "kubernetes_get_unhealthy_resources") {
		t.Fatalf("expected incident triage prompt to mention summary-first tool, got %q", text.Text)
	}
}

func TestHandleQuestionResolutionPrompt(t *testing.T) {
	result, err := HandleQuestionResolutionPrompt(context.Background(), mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Name: QuestionResolutionPromptName,
			Arguments: map[string]string{
				"user_question": "Why is checkout API failing after deploy?",
			},
		},
	})
	if err != nil {
		t.Fatalf("HandleQuestionResolutionPrompt returned error: %v", err)
	}
	text, ok := result.Messages[0].Content.(mcp.TextContent)
	if !ok {
		t.Fatalf("expected text content, got %#v", result.Messages[0].Content)
	}
	if !contains(text.Text, "kubernetes") || !contains(text.Text, "prometheus") || !contains(text.Text, "langfuse") {
		t.Fatalf("expected routing guidance across services, got %q", text.Text)
	}
}

func TestRegistrationAppliesToServices(t *testing.T) {
	promRegistration := Registration{
		Prompt:            PrometheusDiagnosisPrompt(),
		Handler:           HandlePrometheusDiagnosisPrompt,
		RequiredServices: []string{"prometheus"},
	}

	if !registrationAppliesToServices(promRegistration, []string{"prometheus"}) {
		t.Fatal("expected prometheus prompt to match prometheus service")
	}
	if registrationAppliesToServices(promRegistration, []string{"kubernetes"}) {
		t.Fatal("did not expect prometheus prompt to match kubernetes-only service set")
	}

	k8sOnly := Registration{
		Prompt:            K8sOpsPrompt(),
		Handler:           HandleK8sOpsPrompt,
		RequiredServices: []string{"kubernetes"},
	}
	if !registrationAppliesToServices(k8sOnly, []string{"kubernetes", "prometheus"}) {
		t.Fatal("expected kubernetes prompt to match when kubernetes is present")
	}
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
