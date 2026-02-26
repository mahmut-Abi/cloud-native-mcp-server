#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-http://127.0.0.1:8080}"
SSE_PATH="${SSE_PATH:-/api/aggregate/sse}"
API_KEY="${API_KEY:-}"
CONNECT_TIMEOUT="${CONNECT_TIMEOUT:-5}"
SSE_MAX_TIME="${SSE_MAX_TIME:-8}"
POST_MAX_TIME="${POST_MAX_TIME:-8}"

if [[ "${BASE_URL}" != http://* && "${BASE_URL}" != https://* ]]; then
  echo "ERROR: BASE_URL must start with http:// or https://, got: ${BASE_URL}" >&2
  exit 1
fi

sse_url="${BASE_URL%/}${SSE_PATH}"
if [[ -n "${API_KEY}" ]]; then
  separator="?"
  if [[ "${sse_url}" == *\?* ]]; then
    separator="&"
  fi
  sse_url="${sse_url}${separator}api_key=${API_KEY}"
fi

tmp_stream="$(mktemp)"
tmp_body="$(mktemp)"
cleanup() {
  rm -f "${tmp_stream}" "${tmp_body}"
}
trap cleanup EXIT

echo "[1/3] Opening SSE stream: ${sse_url}"
set +e
curl -sS -N \
  --connect-timeout "${CONNECT_TIMEOUT}" \
  --max-time "${SSE_MAX_TIME}" \
  -H "Accept: text/event-stream" \
  "${sse_url}" > "${tmp_stream}"
curl_rc=$?
set -e

if [[ ! -s "${tmp_stream}" ]]; then
  echo "ERROR: SSE endpoint returned no response bytes (curl exit: ${curl_rc})" >&2
  exit 1
fi

if ! grep -q '^event: endpoint' "${tmp_stream}"; then
  echo "ERROR: did not receive 'event: endpoint' from SSE stream" >&2
  echo "---- first 30 lines ----" >&2
  sed -n '1,30p' "${tmp_stream}" >&2
  echo "------------------------" >&2
  exit 1
fi

message_endpoint="$(awk '/^data:/{sub(/^data:[[:space:]]*/,""); print; exit}' "${tmp_stream}")"
if [[ -z "${message_endpoint}" ]]; then
  echo "ERROR: received endpoint event but missing data line" >&2
  exit 1
fi

if [[ "${message_endpoint}" == http://* || "${message_endpoint}" == https://* ]]; then
  message_url="${message_endpoint}"
else
  message_url="${BASE_URL%/}${message_endpoint}"
fi

if [[ -n "${API_KEY}" && "${message_url}" != *"api_key="* ]]; then
  separator="?"
  if [[ "${message_url}" == *\?* ]]; then
    separator="&"
  fi
  message_url="${message_url}${separator}api_key=${API_KEY}"
fi

echo "[2/3] Sending initialize request to: ${message_url}"
initialize_payload='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"sse-smoke-test","version":"1.0.0"}}}'

status_code="$(curl -sS \
  --connect-timeout "${CONNECT_TIMEOUT}" \
  --max-time "${POST_MAX_TIME}" \
  -o "${tmp_body}" \
  -w '%{http_code}' \
  -X POST \
  -H "Content-Type: application/json" \
  --data "${initialize_payload}" \
  "${message_url}")"

if [[ "${status_code}" != "202" && "${status_code}" != "200" ]]; then
  echo "ERROR: initialize request failed with HTTP ${status_code}" >&2
  echo "---- response body ----" >&2
  cat "${tmp_body}" >&2
  echo >&2
  echo "-----------------------" >&2
  exit 1
fi

echo "[3/3] Initialize accepted with HTTP ${status_code}"
echo "PASS: SSE handshake and message endpoint flow are healthy."
