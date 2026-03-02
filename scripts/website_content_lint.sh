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

# 3) Chinese pages should not point to English doc/getting-started absolute paths.
report_fail \
  "Cross-language path mismatch in Chinese content" \
  '\]\(/docs/|\]\(/getting-started/' \
  'website/content/zh'

# 4) Soft check: key docs/posts should include a description field for SEO quality.
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
