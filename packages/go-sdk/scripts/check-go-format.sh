#!/usr/bin/env bash
set -euo pipefail

unformatted=$(gofmt -l -- $(find . -name '*.go' -type f -not -path './.git/*'))
if [[ -n "$unformatted" ]]; then
  printf 'Go files need formatting:\n%s\n' "$unformatted" >&2
  exit 1
fi
