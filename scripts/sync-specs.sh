#!/usr/bin/env bash
set -euo pipefail

mono_dir=${1:?Usage: sync-specs.sh /path/to/mono}
repo_root=$(cd "$(dirname "$0")/.." && pwd)

copy_spec() {
  source_path=$1
  target_path=$2
  install -m 0644 "$mono_dir/$source_path" "$repo_root/$target_path"
}

copy_spec infra/spec/openapi.yml spec/openapi.yml
copy_spec infra/packages/envd/spec/envd.yaml spec/envd/envd.yaml
copy_spec infra/packages/envd/spec/filesystem/filesystem.proto spec/envd/filesystem/filesystem.proto
copy_spec infra/packages/envd/spec/process/process.proto spec/envd/process/process.proto

revision=$(git -C "$mono_dir" rev-parse HEAD)
openapi_sha=$(shasum -a 256 "$repo_root/spec/openapi.yml" | cut -d ' ' -f 1)
envd_sha=$(shasum -a 256 "$repo_root/spec/envd/envd.yaml" | cut -d ' ' -f 1)
filesystem_sha=$(shasum -a 256 "$repo_root/spec/envd/filesystem/filesystem.proto" | cut -d ' ' -f 1)
process_sha=$(shasum -a 256 "$repo_root/spec/envd/process/process.proto" | cut -d ' ' -f 1)

tmp_file=$(mktemp)
trap 'rm -f "$tmp_file"' EXIT
jq -n \
  --arg revision "$revision" \
  --arg openapi "$openapi_sha" \
  --arg envd "$envd_sha" \
  --arg filesystem "$filesystem_sha" \
  --arg process "$process_sha" \
  '{source:"agentbox/mono", revision:$revision, files:{
    "spec/openapi.yml":$openapi,
    "spec/envd/envd.yaml":$envd,
    "spec/envd/filesystem/filesystem.proto":$filesystem,
    "spec/envd/process/process.proto":$process
  }}' > "$tmp_file"
install -m 0644 "$tmp_file" "$repo_root/spec/source.json"
