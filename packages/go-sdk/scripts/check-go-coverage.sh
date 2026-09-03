#!/usr/bin/env bash
set -euo pipefail

root_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
profile=$(mktemp)
trap 'rm -f "$profile"' EXIT

cd "$root_dir"
test_packages=$(go list ./... | grep -v -E '/(internal/gen|tests)/')
go test -covermode=atomic -coverprofile="$profile" $test_packages
coverage=$(go tool cover -func="$profile" | awk '/^total:/ {gsub(/%/, "", $3); print $3}')

awk -v coverage="$coverage" 'BEGIN { if (coverage + 0 < 90) exit 1 }' || {
  printf 'Go statement coverage %s%% is below 90%%\n' "$coverage" >&2
  exit 1
}
printf 'Go statement coverage: %s%%\n' "$coverage"
