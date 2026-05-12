---
title: "Sentry Service"
weight: 11
---

# Sentry Service

The Sentry service provides 9 read-only tools for issue triage, project discovery, and issue event inspection.

## Overview

The Sentry service in Cloud Native MCP Server lets AI assistants query Sentry organizations, projects, issues, and issue events through the official Sentry REST API. It is intended for troubleshooting and error-triage workflows, not write operations.

### Key Capabilities

{{< columns >}}
### 🚩 Issue Triage
List issues in a compact form first, then inspect a single issue in detail.
<--->

### 🧭 Project Discovery
Browse organizations and projects to find the correct Sentry scope.
{{< /columns >}}

{{< columns >}}
### 🧾 Event Inspection
Drill into events attached to an issue to understand the failing request or crash instance.
<--->

### 🔐 Token Validation
Verify that the configured Sentry token and base URL are valid.
{{< /columns >}}

---

## Available Tools (9)

- **sentry_test_connection**: Verify Sentry connectivity
- **sentry_list_organizations**: List organizations visible to the token
- **sentry_list_projects**: List projects in an organization
- **sentry_get_project**: Get a project by organization and project slug
- **sentry_list_issues_summary**: Compact issue discovery view
- **sentry_list_issues**: Full issue listing with filters
- **sentry_get_issue**: Get a specific issue
- **sentry_list_issue_events**: List events for an issue
- **sentry_get_issue_event**: Get a specific issue event

## Configuration Example

```yaml
sentry:
  enabled: true
  url: "https://sentry.io"
  authToken: "sntrys_..."
  organization: "acme"
  project: "frontend"
  timeoutSec: 30
```

## Next Steps

- [Configuration Guide](/docs/configuration/) for environment variables and routes
- [Jaeger Service](/services/jaeger/) for distributed tracing
- [Langfuse Service](/services/langfuse/) for LLM-specific observability
