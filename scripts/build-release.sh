#!/usr/bin/env bash
set -euo pipefail

root_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
release_dir=${1:-"$root_dir/release"}

case "$release_dir" in
  ""|/|"$root_dir")
    echo "Refusing to clean unsafe release directory: $release_dir" >&2
    exit 1
    ;;
esac

mkdir -p "$release_dir/npm" "$release_dir/pypi" "$release_dir/pypi-code-interpreter"
find "$release_dir" -type f -delete

pnpm --dir "$root_dir/packages/js-sdk" build
pnpm --dir "$root_dir/packages/code-interpreter-js" build
pnpm --dir "$root_dir/packages/cli" build

pnpm --dir "$root_dir/packages/js-sdk" pack --pack-destination "$release_dir/npm"
pnpm --dir "$root_dir/packages/code-interpreter-js" pack --pack-destination "$release_dir/npm"
pnpm --dir "$root_dir/packages/cli" pack --pack-destination "$release_dir/npm"

uv build --project "$root_dir/packages/python-sdk" \
  --out-dir "$release_dir/pypi" --no-create-gitignore
uv build --project "$root_dir/packages/code-interpreter-python" \
  --out-dir "$release_dir/pypi-code-interpreter" --no-create-gitignore

(
  cd "$release_dir"
  while IFS= read -r -d '' artifact; do
    digest=$(shasum -a 256 "$artifact" | cut -d ' ' -f 1)
    printf '%s  %s\n' "$digest" "${artifact##*/}"
  done < <(
    find npm pypi pypi-code-interpreter -type f -print0 | sort -z
  ) > SHA256SUMS
)
