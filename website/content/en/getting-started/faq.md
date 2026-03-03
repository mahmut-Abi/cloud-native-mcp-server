---
title: "FAQ"
weight: 20
description: "Frequently asked questions for first-time setup and production rollout."
---

# Getting Started FAQ

## Which mode should I choose first?

Use this default decision path:

- `streamable-http`: best fit for modern MCP transport in production
- `sse`: broad compatibility with existing MCP clients and integrations

If you are unsure, start with `sse`, complete validation, then migrate to `streamable-http`.

## Why do I get `401 unauthorized`?

Check the following in order:

1. `MCP_AUTH_ENABLED=true` is actually set in runtime environment.
2. `MCP_AUTH_MODE` is one of `apikey`, `bearer`, `basic`.
3. For `apikey` mode, `MCP_AUTH_API_KEY` is non-empty and matches request value.
4. API key is sent via `X-Api-Key` header or `api_key` query parameter.
5. Container/process was restarted after configuration changes.

## Is API key header or query parameter preferred?

Both are supported, but use request headers in production whenever possible:

```bash
curl -sS -N \
  -H "X-Api-Key: ChangeMe-Strong-Key-123!" \
  http://127.0.0.1:8080/api/aggregate/sse
```

Use query parameter mainly for quick local checks.

## How can I reduce response size for AI agents?

- Use summary-focused tools before requesting full payloads.
- Use pagination-capable tools for large result sets.
- Disable unnecessary services via `MCP_DISABLED_SERVICES`.
- Apply request rate limits to avoid burst overload.

## Can I enable only a subset of services?

Yes. You can explicitly enable or disable services:

```bash
export MCP_ENABLED_SERVICES="kubernetes,helm,prometheus"
export MCP_DISABLED_SERVICES="kibana,jaeger"
```

Keep configuration consistent to avoid overlapping expectations in team environments.

## What is a safe production baseline?

- Enable authentication and rotate secrets regularly.
- Protect upstream credentials (Grafana, Prometheus, Kibana, etc.).
- Enable structured logging and audit if required.
- Monitor `/health` and key service integrations continuously.
- Run load tests against realistic tool-call patterns.

## Where should I go next?

- [Troubleshooting]({{< relref "troubleshooting.md" >}})
- [Security Guide]({{< relref "/docs/security.md" >}})
- [Configuration Guide]({{< relref "/docs/configuration.md" >}})
- [Performance Guide]({{< relref "/docs/performance.md" >}})
