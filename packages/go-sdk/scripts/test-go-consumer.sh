#!/usr/bin/env bash
set -euo pipefail

root_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
consumer_dir=$(mktemp -d)
trap 'rm -rf "$consumer_dir"' EXIT

cd "$consumer_dir"
go mod init example.com/agentbox-consumer >/dev/null
go mod edit -replace github.com/abox-dev/sdk/packages/go-sdk="$root_dir"
go mod edit -require github.com/abox-dev/sdk/packages/go-sdk@v0.0.0
cp "$root_dir/tests/consumer/consumer_test.go" .
go mod tidy
go test ./...
