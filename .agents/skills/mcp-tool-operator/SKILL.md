---
name: mcp-tool-operator
description: Use when an agent needs to work through MCP tools instead of shell commands or direct API calls. This skill is agent-neutral and covers tool discovery, argument shaping, response inspection, safe mutation patterns, and common failure handling for any MCP-capable environment.
---

# MCP Tool Operator

## Overview

Use this skill when the task should be executed through MCP tools.
It is generic by design and does not depend on a specific agent runtime, SDK, or server implementation.

The only assumptions are that the host can list tools, call a tool by name with JSON arguments, and inspect the returned tool result.

## Operating Pattern

1. Discover the available tool inventory from the running server.
2. Choose the smallest tool that can answer the question.
3. Prefer read-only tools before state-changing tools.
4. Send flat structured JSON arguments unless the host client imposes another wrapper.
5. Inspect the raw tool result before parsing or transforming it.
6. Verify the outcome with a second read tool when the action changes state.

## Tool Choice Rules

- Prefer summary, list, search, health, or paginated tools before full-detail tools.
- Prefer service-specific tools over generic helper tools.
- If multiple tools could work, choose the one that returns the least data needed to answer the request.
- If the user intent is ambiguous, start with discovery tools rather than mutation tools.

## Argument Rules

- Use the exact runtime tool name returned by the server inventory.
- Use the schema field names returned by the server as the canonical argument names.
- Send objects and arrays as structured JSON when the client supports it.
- Avoid JSON-encoding nested payloads into strings unless the server explicitly requires it.
- Avoid guessing aliases such as camelCase if the runtime inventory shows snake_case.

## Result Rules

- Do not assume every tool returns a JSON string.
- Some clients return parsed objects or arrays directly.
- Some clients return an MCP envelope with the payload in a text field.
- Inspect the raw result first, then parse only if needed.
- If parsing fails with messages similar to `[object Object]`, the usual problem is double parsing.

## Mutation Rules

- Do not mutate state unless the user explicitly asked for it or the task clearly requires it.
- Before a destructive action, verify target identity and scope.
- After a mutation, run a read tool to confirm the change.
- Prefer targeted patch or update operations over broad replace operations when both exist.

## Troubleshooting

- `Tool not found`
  - Re-read the live tool inventory
  - Check for runtime naming differences
  - Confirm the server has been restarted after recent tool additions

- Missing required parameter
  - Re-read the tool schema
  - Stop guessing argument names
  - Check whether the server expects a resource kind, namespace, or other locator field

- Parse failure
  - Inspect the raw tool result first
  - Avoid calling a JSON parser on objects that are already parsed

- Output too large
  - Switch to summary or paginated tools
  - Add filters
  - Request only the needed fields when the server supports field selection
