#!/usr/bin/env bash
set -euo pipefail

root_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
release_dir=${1:-"$root_dir/release"}
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT

mkdir -p "$test_dir/checksums"
find "$release_dir"/npm \
  "$release_dir"/pypi \
  "$release_dir"/pypi-code-interpreter \
  -type f -exec cp {} "$test_dir/checksums" \;
cp "$release_dir/SHA256SUMS" "$test_dir/checksums/SHA256SUMS"
(
  cd "$test_dir/checksums"
  shasum -a 256 -c SHA256SUMS
)

mkdir -p "$test_dir/npm"
cd "$test_dir/npm"
npm init --yes >/dev/null
npm install "$release_dir"/npm/*.tgz >/dev/null
node --input-type=module <<'JS'
import { AgentBox, Sandbox, Template } from '@abox-dev/sdk'
import { Sandbox as CodeInterpreter } from '@abox-dev/code-interpreter'

if (!AgentBox || !Sandbox || !Template || !CodeInterpreter) {
  throw new Error('required AgentBox JavaScript exports are missing')
}
JS
./node_modules/.bin/agentbox --version

uv venv "$test_dir/python" --python 3.12 >/dev/null
uv pip install --python "$test_dir/python/bin/python" \
  "$release_dir"/pypi/*.whl \
  "$release_dir"/pypi-code-interpreter/*.whl >/dev/null
"$test_dir/python/bin/python" - <<'PY'
from agentbox import AgentBox, AsyncSandbox, Sandbox, Template
from agentbox_code_interpreter import AsyncSandbox as AsyncCodeInterpreter
from agentbox_code_interpreter import Sandbox as CodeInterpreter

assert all((AgentBox, AsyncSandbox, Sandbox, Template, AsyncCodeInterpreter, CodeInterpreter))
PY
