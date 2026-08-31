#!/usr/bin/env bash
set -euo pipefail

root_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
version=${1:-$(node -p 'require(process.argv[1]).version' \
  "$root_dir/packages/js-sdk/package.json")}
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT

: "${AGENTBOX_API_KEY:?AGENTBOX_API_KEY is required}"
: "${AGENTBOX_API_URL:?AGENTBOX_API_URL is required}"
: "${AGENTBOX_SANDBOX_URL:?AGENTBOX_SANDBOX_URL is required}"

mkdir -p "$test_dir/npm"
cd "$test_dir/npm"
npm init --yes >/dev/null
npm install \
  "@abox-dev/sdk@$version" \
  "@abox-dev/code-interpreter@$version" \
  "@abox-dev/cli@$version" >/dev/null
npm ls \
  "@abox-dev/sdk@$version" \
  "@abox-dev/code-interpreter@$version" \
  "@abox-dev/cli@$version" >/dev/null

cp "$root_dir/tests/runtime/core-js.mjs" "$test_dir/npm/core-js.mjs"
cp "$root_dir/tests/runtime/code-interpreter-js.mjs" \
  "$test_dir/npm/code-interpreter-js.mjs"
node "$test_dir/npm/core-js.mjs"
node "$test_dir/npm/code-interpreter-js.mjs"
AGENTBOX_CLI_BIN="$test_dir/npm/node_modules/.bin/agentbox" \
AGENTBOX_RUNTIME_TMP="$test_dir/cli" \
AGENTBOX_RUNTIME_DOCKERFILE="$root_dir/tests/runtime/Agentbox.Dockerfile" \
  "$root_dir/tests/runtime/cli.sh"

uv venv "$test_dir/python" --python 3.12 >/dev/null
uv pip install --python "$test_dir/python/bin/python" \
  "abox-sdk==$version" \
  "abox-code-interpreter==$version" >/dev/null
uv pip show --python "$test_dir/python/bin/python" \
  abox-sdk abox-code-interpreter >/dev/null
"$test_dir/python/bin/python" "$root_dir/tests/runtime/core-python.py"
"$test_dir/python/bin/python" \
  "$root_dir/tests/runtime/code_interpreter_python.py"
