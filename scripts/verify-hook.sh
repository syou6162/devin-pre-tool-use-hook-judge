#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY="${REPO_ROOT}/devin-pre-tool-use-hook-judge"

echo "==> Building local binary for hook verification..."
(
  cd "${REPO_ROOT}"
  go build -o "${BINARY}" .
)

pass() {
  echo "PASS: $1"
}

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

assert_exit() {
  local expected="$1"
  local actual="$2"
  local name="$3"
  if [[ "${actual}" -ne "${expected}" ]]; then
    fail "${name}: expected exit ${expected}, got ${actual}"
  fi
}

assert_json_field() {
  local json="$1"
  local field="$2"
  local expected="$3"
  local name="$4"
  local actual
  actual="$(printf '%s' "${json}" | jq -r ".${field}")"
  if [[ "${actual}" != "${expected}" ]]; then
    fail "${name}: expected ${field}=${expected}, got ${actual}"
  fi
}

echo "==> Case 1: missing config should block"
set +e
output="$(
  printf '%s' '{"hook_event_name":"PreToolUse","tool_name":"exec","tool_input":{"command":"git push"}}' \
    | "${BINARY}" 2>/dev/null
)"
exit_code=$?
set -e
assert_exit 2 "${exit_code}" "missing config"
assert_json_field "${output}" "decision" "block" "missing config"
pass "missing config blocks safely"

echo "==> Case 2: invalid input should block"
set +e
output="$(
  printf '%s' '{"hook_event_name":"PreToolUse","tool_name":"","tool_input":{}}' \
    | "${BINARY}" --builtin validate_git_push 2>/dev/null
)"
exit_code=$?
set -e
assert_exit 2 "${exit_code}" "invalid input"
assert_json_field "${output}" "decision" "block" "invalid input"
pass "invalid input blocks safely"

echo "==> Case 3: hooks.v1.json is valid JSON"
jq empty "${REPO_ROOT}/.devin/hooks.v1.json"
pass "hooks.v1.json parses as JSON"

echo
echo "Hook protocol verification completed successfully."
echo "Note: full Devin CLI integration requires the devin binary and credentials."
