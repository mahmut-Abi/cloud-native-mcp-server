#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v rg >/dev/null 2>&1; then
  echo "ERROR: 'rg' (ripgrep) is required for website content lint checks." >&2
  exit 2
fi

FAIL_COUNT=0
WARN_COUNT=0

print_section() {
  printf '\n== %s ==\n' "$1"
}

report_fail() {
  local title="$1"
  local pattern="$2"
  local path="$3"

  print_section "$title"
  if rg -n --glob '*.md' "$pattern" "$path"; then
    FAIL_COUNT=$((FAIL_COUNT + 1))
  else
    echo "OK"
  fi
}

report_warn_missing_description() {
  local title="$1"
  shift
  local -a files=("$@")
  local missing=()

  print_section "$title"

  for file in "${files[@]}"; do
    if ! awk '
      BEGIN { in_fm=0; found=0; has_fm=0 }
      NR==1 && /^---$/ { in_fm=1; has_fm=1; next }
      in_fm && /^---$/ { exit(found ? 0 : 1) }
      in_fm && /^description:[[:space:]]*[^[:space:]].*/ { found=1 }
      END {
        if (!has_fm) exit 0
        if (in_fm && !found) exit 1
      }
    ' "$file"; then
      missing+=("$file")
    fi
  done

  if ((${#missing[@]} == 0)); then
    echo "OK"
  else
    WARN_COUNT=$((WARN_COUNT + ${#missing[@]}))
    printf 'Missing description in %d file(s):\n' "${#missing[@]}"
    printf '  - %s\n' "${missing[@]}"
  fi
}

report_fail_files() {
  local title="$1"
  local pattern="$2"
  shift 2
  local -a files=("$@")

  print_section "$title"
  if rg -n "$pattern" "${files[@]}"; then
    FAIL_COUNT=$((FAIL_COUNT + 1))
  else
    echo "OK"
  fi
}

# 1) Legacy variables and endpoints that should no longer appear in docs/content.
report_fail \
  "Legacy MCP variables or endpoints" \
  'MCP_SERVER_|/api/[a-z-]+/http|/api/audit/query|/v1/mcp/list-tools|/v1/mcp/call-tool|POST /v1/mcp|https?://[^ ]*/v1/mcp|MCP_PROMETHEUS_URL|MCP_ELASTICSEARCH_URL|MCP_ALERTMANAGER_URL' \
  'website/content'

# 2) Root/project docs should not use legacy auth variable name.
report_fail \
  "Legacy auth env var names in markdown" \
  '\bMCP_API_KEY\b' \
  '.'

# 3) Core configuration/deployment docs should not drift to legacy key names.
CORE_DOC_FILES=(
  "docs/CONFIGURATION.md"
  "docs/DEPLOYMENT.md"
  "docs/PERFORMANCE.md"
  "website/content/en/configuration.md"
  "website/content/en/docs/configuration.md"
  "website/content/zh/configuration.md"
  "website/content/zh/docs/configuration.md"
  "website/content/en/deployment.md"
  "website/content/en/docs/deployment.md"
  "website/content/zh/deployment.md"
  "website/content/zh/docs/deployment.md"
  "website/content/en/performance.md"
  "website/content/en/docs/performance.md"
  "website/content/zh/performance.md"
  "website/content/zh/docs/performance.md"
  "website/content/zh/guides/performance/_index.md"
  "website/content/zh/guides/performance/optimization.md"
  "website/content/zh/guides/performance/benchmarking.md"
  "website/content/en/architecture.md"
  "website/content/en/docs/architecture.md"
  "website/content/zh/architecture.md"
  "website/content/zh/docs/architecture.md"
  "website/content/zh/guides/configuration/_index.md"
  "website/content/zh/guides/deployment/_index.md"
)

report_fail_files \
  "Legacy config key names in core docs" \
  '\b(max_connections|cache_ttl|query_timeout|header_name|api_keys|max_response_size|compression_enabled|json_pool_size|file_path|read_timeout|write_timeout|timeout_sec|bearer_token|tls_skip_verify|tls_cert_file|tls_key_file|tls_ca_file)\s*:|\bapi_key\s*:' \
  "${CORE_DOC_FILES[@]}"

report_fail_files \
  "Non-existent config validation flag in docs" \
  'cloud-native-mcp-server[^\n`]*--validate-config' \
  "${CORE_DOC_FILES[@]}"

report_fail_files \
  "Stale validation error text in docs" \
  'missing required field "api_key"|invalid service URL "grafana:3000"|invalid server mode "invalid"' \
  "${CORE_DOC_FILES[@]}"

# 4) Chinese pages should not point to English doc/getting-started absolute paths.
report_fail \
  "Cross-language path mismatch in Chinese content" \
  '\]\(/docs/|\]\(/getting-started/' \
  'website/content/zh'

# 5) Soft check: key docs/posts should include a description field for SEO quality.
mapfile -t KEY_FILES < <(
  find website/content/en/docs website/content/zh/docs website/content/en/posts website/content/zh/posts \
    -type f -name '*.md' ! -name '_index.md' | sort
)
report_warn_missing_description "SEO metadata completeness (description field)" "${KEY_FILES[@]}"

print_section "Summary"
printf 'Fail checks: %d\n' "$FAIL_COUNT"
printf 'Warnings: %d\n' "$WARN_COUNT"

if ((FAIL_COUNT > 0)); then
  echo "website content lint failed"
  exit 1
fi

echo "website content lint passed"
