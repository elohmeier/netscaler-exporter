#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/config/grafana/provisioning/dashboards"

cd "$SCRIPT_DIR"

# Install dependencies if vendor directory doesn't exist
if [[ ! -d vendor ]]; then
    echo "Installing dependencies..."
    jb install
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Generate dashboards
echo "Generating environment.json..."
jsonnet -J vendor environment.jsonnet > "${OUTPUT_DIR}/environment.json"

echo "Generating ha-pair.json..."
jsonnet -J vendor ha-pair.jsonnet > "${OUTPUT_DIR}/ha-pair.json"

echo "Generating chain.json..."
jsonnet -J vendor chain.jsonnet > "${OUTPUT_DIR}/chain.json"

echo "Generating problems.json..."
jsonnet -J vendor problems.jsonnet > "${OUTPUT_DIR}/problems.json"

echo "Done! Dashboards written to: ${OUTPUT_DIR}"
