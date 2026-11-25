#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

export PATH="/home/ninofelino12/go/bin:$PATH:$HOME/google-cloud-sdk/bin:$PATH"

# Project ID Firestore (default: arduino-ee66f sesuai Firebase config yang diberikan)
export GOOGLE_APPLICATION_CREDENTIALS="${GOOGLE_APPLICATION_CREDENTIALS:-/home/ninofelino12/service-account.json}"
export FIREBASE_PROJECT_ID="${FIREBASE_PROJECT_ID:-bizflow-dashboard-vlo72}"
export FIRESTORE_COLLECTION="${FIRESTORE_COLLECTION:-items}"
export PORT="${PORT:-8080}"

export GOCACHE="${GOCACHE:-/home/ninofelino12/.cache/go-build}"
mkdir -p "$GOCACHE"
echo "=== Backend dengan Firestore ==="
echo "Project ID    : $FIREBASE_PROJECT_ID"
echo "Collection    : $FIRESTORE_COLLECTION"
echo "Port          : $PORT"
echo "Go version    : $(go version)"
echo "gcloud account: $(gcloud config get-value account 2>/dev/null || echo 'not logged in')"
echo "==============================="

echo "[1/3] Download dependency..."
go mod download

echo "[2/3] Build binary..."
go build -o backend .

echo "[3/3] Run server..."
echo "Ctrl+C untuk berhenti"
./backend
