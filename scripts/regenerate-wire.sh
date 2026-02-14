#!/bin/bash

# Script to regenerate wire_gen.go using Wire
# Usage: ./scripts/regenerate-wire.sh

set -e

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Wire package directory
WIRE_DIR="${PROJECT_ROOT}/internal/app"

echo "Regenerating wire_gen.go..."
echo "Wire directory: ${WIRE_DIR}"
echo ""

# Check if wire directory exists
if [ ! -d "$WIRE_DIR" ]; then
    echo "Error: Wire directory not found at $WIRE_DIR"
    exit 1
fi

# Check if wire.go exists
if [ ! -f "${WIRE_DIR}/wire.go" ]; then
    echo "Error: wire.go not found at ${WIRE_DIR}/wire.go"
    exit 1
fi

# Change to wire directory and run wire
cd "$WIRE_DIR"
go run -mod=mod github.com/google/wire/cmd/wire

echo ""
echo "wire_gen.go regenerated successfully!"

