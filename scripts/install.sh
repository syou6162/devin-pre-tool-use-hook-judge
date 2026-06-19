#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INSTALL_DIR="${INSTALL_DIR:-${HOME}/.local/bin}"
TARGET_PROJECT="${1:-.}"
BINARY_NAME="devin-pre-tool-use-hook-judge"
BUILTIN="${BUILTIN:-validate_git_push}"
HOOKS_FILE="${TARGET_PROJECT%/}/.devin/hooks.v1.json"

echo "==> Building ${BINARY_NAME}..."
mkdir -p "${INSTALL_DIR}"
(
  cd "${REPO_ROOT}"
  go build -o "${INSTALL_DIR}/${BINARY_NAME}" .
)

if ! command -v "${INSTALL_DIR}/${BINARY_NAME}" >/dev/null 2>&1; then
  echo "error: failed to install ${BINARY_NAME} to ${INSTALL_DIR}" >&2
  exit 1
fi

echo "==> Installed ${INSTALL_DIR}/${BINARY_NAME}"

mkdir -p "$(dirname "${HOOKS_FILE}")"
cat > "${HOOKS_FILE}" <<EOF
{
  "PreToolUse": [
    {
      "matcher": "^exec$",
      "hooks": [
        {
          "type": "command",
          "command": "${INSTALL_DIR}/${BINARY_NAME} --builtin ${BUILTIN}",
          "timeout": 120
        }
      ]
    }
  ]
}
EOF

echo "==> Wrote ${HOOKS_FILE}"
echo
echo "Installation complete."
echo
echo "Next steps:"
echo "  1. Ensure ${INSTALL_DIR} is on your PATH"
echo "  2. Install Devin CLI if you have not already: https://docs.devin.ai/cli"
echo "  3. Run Devin CLI in ${TARGET_PROJECT} and verify PreToolUse hooks fire"
echo
echo "Optional environment variables:"
echo "  INSTALL_DIR   Install destination for the binary (default: ~/.local/bin)"
echo "  BUILTIN       Builtin validator name (default: validate_git_push)"
echo "  TARGET        Project directory passed as the first argument (default: .)"
