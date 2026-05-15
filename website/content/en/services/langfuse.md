---
title: "Langfuse Service"
weight: 10
---

# Langfuse Service

The Langfuse service provides LLM observability and evaluation workflows with 37 tools for traces, sessions, prompts, scores, datasets, models, annotation queues, metrics, project management, membership management, and API key management.

## Overview

The Langfuse service in Cloud Native MCP Server lets AI assistants inspect prompt executions, user sessions, observations, and evaluation results from Langfuse using the Public API. It is useful when you want to trace how an LLM interaction behaved end to end without leaving the MCP surface.

### Key Capabilities

{{< columns >}}
### 🧭 Trace Discovery
Browse traces and drill into a single trace when you need full detail.
<--->

### 🧵 Session Correlation
Inspect sessions to connect multiple traces in one user flow.
{{< /columns >}}

{{< columns >}}
### 📝 Prompt Intelligence
Look up prompt versions, labels, and resolved prompt content.
<--->

### 📏 Evaluation Insights
Inspect scores and query Langfuse metrics for trend analysis.
{{< /columns >}}

---

## Available Tools (37)

### Health
- **langfuse_check_health**: Check Langfuse API and database health

### Traces
- **langfuse_list_traces_summary**: Compact trace discovery view
- **langfuse_list_traces**: Full trace listing with filters
- **langfuse_get_trace**: Get a specific trace by ID

### Annotation Queues and Datasets
- **langfuse_list_annotation_queues**: List annotation queues
- **langfuse_get_annotation_queue**: Get one annotation queue
- **langfuse_list_annotation_queue_items**: List queue items
- **langfuse_list_datasets**: List datasets
- **langfuse_get_dataset**: Get one dataset
- **langfuse_list_dataset_runs**: List dataset runs
- **langfuse_get_dataset_run**: Get one dataset run

### Sessions and Observations
- **langfuse_list_sessions**: List sessions
- **langfuse_get_session**: Get a specific session
- **langfuse_list_observations**: List observations
- **langfuse_get_observation**: Get a specific observation

### Models and Score Configurations
- **langfuse_list_llm_connections**: List LLM connections
- **langfuse_list_models**: List models
- **langfuse_get_model**: Get one model
- **langfuse_list_score_configs**: List score configurations
- **langfuse_get_score_config**: Get one score configuration

### Prompts, Scores, and Metrics
- **langfuse_list_prompts**: List prompt versions
- **langfuse_get_prompt**: Get a prompt by name
- **langfuse_list_scores**: List scores and evaluations
- **langfuse_get_score**: Get a specific score
- **langfuse_get_metrics**: Run a Langfuse metrics query

### Projects
- **langfuse_get_project**: Get the project associated with the configured credentials
- **langfuse_list_organization_projects**: List organization projects
- **langfuse_create_project**: Create a project
- **langfuse_update_project**: Update a project
- **langfuse_delete_project**: Delete a project
- **langfuse_list_project_memberships**: List project memberships
- **langfuse_upsert_project_membership**: Create or update a project membership
- **langfuse_delete_project_membership**: Delete a project membership

### API Key Management
- **langfuse_list_organization_api_keys**: List organization API keys
- **langfuse_list_project_api_keys**: List project API keys
- **langfuse_create_project_api_key**: Create a project API key
- **langfuse_delete_project_api_key**: Delete a project API key

---

## Configuration Example

```yaml
langfuse:
  enabled: true
  url: "https://cloud.langfuse.com"
  username: "pk-lf-..."
  password: "sk-lf-..."
  timeoutSec: 30
```

Langfuse uses HTTP Basic Auth on the Public API:

- Username: Langfuse public key
- Password: Langfuse secret key

`publicKey` and `secretKey` are still accepted as deprecated aliases for `username` and `password`. The API key management tools require organization-scoped Langfuse credentials.

## Next Steps

- [OpenTelemetry Service](/services/opentelemetry/) for infrastructure telemetry
- [Jaeger Service](/services/jaeger/) for distributed trace troubleshooting
- [Configuration Guide](/docs/configuration/) for environment variables and endpoints
