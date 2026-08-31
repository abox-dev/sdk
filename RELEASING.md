# Releasing AgentBox SDK

This document is the maintainer source of truth for changing, verifying, and
releasing the AgentBox SDK packages. The SDK repository owns package-level and
end-to-end SDK tests. The mono repository owns backend HTTP/Connect regression
tests and does not duplicate SDK test implementations.

All five public packages use one version:

- npm: `@abox-dev/sdk`, `@abox-dev/code-interpreter`, `@abox-dev/cli`;
- PyPI: `abox-sdk`, `abox-code-interpreter`.

## 1. Prepare the change

Start from an up-to-date `main` and develop the change on a branch. If the
backend API or envd contract changed, merge the corresponding mono change
first, then update the checked-in snapshots:

```bash
make sync-specs MONO_DIR=/path/to/mono
make generate
git diff --check
```

`make generate` uses only checked-in snapshots. Review the snapshot manifest,
mono revision, SHA-256 values, generated clients, and public API changes in the
same pull request. SDK-only changes do not require a spec sync.

## 2. Set the release version

Choose the next semantic version and update every package together:

```bash
pnpm release:version X.Y.Z
pnpm install --lockfile-only
uv lock --project packages/python-sdk
uv lock --project packages/code-interpreter-python
node scripts/check-release-versions.mjs vX.Y.Z
```

`release:version` updates the five workspace manifests and both Python
`pyproject.toml` files. Regenerate both Python lock files rather than editing
them by hand. Do not release the JavaScript and Python packages at different
versions.

## 3. Verify source and release artifacts

Install frozen dependencies and run the same checks as CI:

```bash
pnpm install --frozen-lockfile
uv sync --project packages/python-sdk --frozen
uv sync --project packages/code-interpreter-python --frozen
pnpm format:check
pnpm lint
pnpm typecheck
pnpm test
make generate
git diff --exit-code
./scripts/build-release.sh
node scripts/check-public-artifacts.mjs release
./scripts/test-release-artifacts.sh
```

`build-release.sh` creates the npm tarballs, Python wheel/sdist files, and
`SHA256SUMS`. The artifact tests install those files into clean temporary
environments; workspace links are not used.

## 4. Verify against the KVM runtime

Start the current backend from the sibling mono checkout with `make start` and
export its `AGENTBOX_API_KEY`, `AGENTBOX_API_URL`, and
`AGENTBOX_SANDBOX_URL`. Test the exact artifacts that are about to be
published:

```bash
./scripts/test-release-runtime.sh
```

The suite covers JavaScript, Python sync/async, Code Interpreter, CLI,
private traffic, and a temporary template build. It owns and removes its test
resources. A release must not be tagged if this suite fails.

## 5. Merge and tag

Merge the change only after the pull-request CI succeeds. Update local `main`,
verify that the working tree is clean, and create an annotated tag on that
commit:

```bash
git switch main
git pull --ff-only origin main
node scripts/check-release-versions.mjs vX.Y.Z
git tag -a vX.Y.Z -m "AgentBox SDK vX.Y.Z"
git push origin vX.Y.Z
```

Do not move or reuse a published tag. If a release needs a correction, publish
a new patch version.

The tag workflow builds the artifacts once, verifies them, publishes npm via
Trusted Publishing, publishes both PyPI projects via their GitHub environments,
and creates a GitHub Release with checksums. An existing registry file is
accepted only when its digest matches the newly built artifact.

## 6. Verify the published packages

After all registry jobs and the GitHub Release succeed, install the public
packages into clean environments and repeat the KVM suite:

```bash
./scripts/test-published-runtime.sh X.Y.Z
```

Verify the GitHub Release assets with its `SHA256SUMS`. Document public API
changes in the AgentBox documentation repository. No SDK-version update is
required in mono: its direct HTTP/Connect tests intentionally remain independent
from SDK releases.
