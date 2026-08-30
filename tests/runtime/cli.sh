#!/usr/bin/env bash
set -euo pipefail

: "${AGENTBOX_CLI_BIN:?AGENTBOX_CLI_BIN is required}"
: "${AGENTBOX_RUNTIME_TMP:?AGENTBOX_RUNTIME_TMP is required}"

export HOME="$AGENTBOX_RUNTIME_TMP/home"
export NO_COLOR=1
mkdir -p "$HOME"

sandbox_id=
template_name="sdk-cli-$(date +%s)"
template_created=false

cleanup() {
  if [[ -n "$sandbox_id" ]]; then
    "$AGENTBOX_CLI_BIN" sandbox kill "$sandbox_id" >/dev/null 2>&1 || true
  fi
  if [[ "$template_created" == true ]]; then
    "$AGENTBOX_CLI_BIN" template delete "$template_name" --yes \
      >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

"$AGENTBOX_CLI_BIN" configure ab_config_value >/dev/null

python3 - "$HOME/.agentbox/config.json" <<'PY'
import json
import stat
import sys
from pathlib import Path

config = Path(sys.argv[1])
assert stat.S_IMODE(config.parent.stat().st_mode) == 0o700
assert stat.S_IMODE(config.stat().st_mode) == 0o600
assert json.loads(config.read_text()) == {"apiKey": "ab_config_value"}
PY

# The valid environment key must override the deliberately unusable configured key.
sandbox_id=$(
  "$AGENTBOX_CLI_BIN" sandbox create base --detach --timeout 180 \
    | tail -n 1 \
    | tr -d '\r'
)
[[ "$sandbox_id" =~ ^[a-z0-9]+$ ]]

"$AGENTBOX_CLI_BIN" sandbox list --format json \
  | node -e '
    let input = "";
    process.stdin.on("data", (chunk) => (input += chunk));
    process.stdin.on("end", () => {
      const id = process.argv[1];
      const rows = JSON.parse(input);
      if (!rows.some((row) => row.sandboxId === id)) process.exit(1);
    });
  ' "$sandbox_id"

"$AGENTBOX_CLI_BIN" sandbox info "$sandbox_id" --format json \
  | node -e '
    let input = "";
    process.stdin.on("data", (chunk) => (input += chunk));
    process.stdin.on("end", () => {
      if (JSON.parse(input).sandboxId !== process.argv[1]) process.exit(1);
    });
  ' "$sandbox_id"

exec_output=$(
  "$AGENTBOX_CLI_BIN" sandbox exec "$sandbox_id" printf cli-agentbox
)
[[ "$exec_output" == "cli-agentbox" ]]
"$AGENTBOX_CLI_BIN" sandbox logs "$sandbox_id" --format json >/dev/null
"$AGENTBOX_CLI_BIN" sandbox pause "$sandbox_id" >/dev/null
"$AGENTBOX_CLI_BIN" sandbox resume "$sandbox_id" >/dev/null
"$AGENTBOX_CLI_BIN" sandbox kill "$sandbox_id" >/dev/null
sandbox_id=

template_dir="$AGENTBOX_RUNTIME_TMP/template"
mkdir -p "$template_dir"
cp "$AGENTBOX_RUNTIME_DOCKERFILE" "$template_dir/Agentbox.Dockerfile"

"$AGENTBOX_CLI_BIN" template build "$template_name" \
  --path "$template_dir" \
  --dockerfile Agentbox.Dockerfile \
  --cpu-count 1 \
  --memory-mb 512
template_created=true

"$AGENTBOX_CLI_BIN" template list --format json \
  | node -e '
    let input = "";
    process.stdin.on("data", (chunk) => (input += chunk));
    process.stdin.on("end", () => {
      const name = process.argv[1];
      const rows = JSON.parse(input);
      if (!rows.some((row) => row.names.some((item) => item.endsWith(`/${name}`)))) {
        process.exit(1);
      }
    });
  ' "$template_name"

"$AGENTBOX_CLI_BIN" template delete "$template_name" --yes >/dev/null
template_created=false

echo "CLI runtime smoke passed"
