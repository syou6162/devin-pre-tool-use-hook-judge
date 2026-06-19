#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="devin-pre-tool-use-hook-judge"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Building ${BINARY_NAME}..."
(cd "${REPO_ROOT}" && go build -o "/tmp/${BINARY_NAME}" .)

mkdir -p "${INSTALL_DIR}"
install -m 0755 "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
rm -f "/tmp/${BINARY_NAME}"

echo "Installed ${INSTALL_DIR}/${BINARY_NAME}"
echo
echo "Add hook configuration from ${REPO_ROOT}/.devin/hooks.v1.json to your project."
echo "Example:"
echo "  cp ${REPO_ROOT}/.devin/hooks.v1.json .devin/hooks.v1.json"
