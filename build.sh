#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if ! command -v go >/dev/null 2>&1; then
  echo "Go tidak ditemukan. Install Go terlebih dahulu."
  exit 1
fi

echo "=== Build backend (portable) ==="
echo "Go version: $(go version)"
echo "[1/3] Tidy module..."
go mod tidy

echo "[2/3] Download module..."
go mod download

echo "[3/3] Build binary..."
go build -o backend .

echo "BUILD OK: $(pwd)/backend"
