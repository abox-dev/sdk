Use pnpm for JavaScript packages, uv for Python packages, and Go modules for Go packages.
Use English exclusively in source code, comments, documentation, commit messages, and GitHub pull request titles and descriptions.
Keep the JavaScript, Python sync/async, and Go SDKs behaviorally aligned, including their Code Interpreter APIs.
Use only Go syntax and runtime dependencies compatible with the `go` directive in `packages/go-sdk/go.mod`.
Run format checks, lint, type checks, unit tests, deterministic generation, builds, package-install checks, and the Go race and coverage checks before committing. Handwritten Go code must keep at least 90% statement coverage; generated packages are excluded from the threshold.
The API and envd snapshots under spec/ are generated from mono/infra. Do not edit them manually. Update them with `make sync-specs MONO_DIR=/path/to/mono`, then run `make generate`.
Generated clients must depend only on checked-in snapshots and never fetch network content during generation. Do not edit generated Go files under `packages/go-sdk/internal/gen` manually.
Public APIs, package artifacts, examples, errors, environment variables, and headers must use AgentBox naming. Upstream names are allowed only in licenses, attribution, pinned build-only codegen tooling, and wire/protobuf namespaces that are required by the runtime protocol.
Default development credentials may be stored in `.env.local` or `~/.agentbox/config.json`; never print or commit them.

When a new Go version is released:

- Treat the `packages/go-sdk/go.mod` directive as the only source of the minimum supported Go version. Raise it only when intentionally ending support for an older release.
- Update the exact local patch version in `.tool-versions` and the Go builder in `codegen.Dockerfile`; generated code must still compile with the minimum version from `packages/go-sdk/go.mod`.
- Keep `.github/workflows/ci.yml` on a continuous matrix of every Go minor release from the `packages/go-sdk/go.mod` minimum through the current stable release. Add new minors when released; remove old minors only when raising the minimum.
- Update the Go toolchain used by `.github/workflows/release.yml` and the supported-version instructions in `RELEASING.md`.
- For a patch release, update exact tooling pins. For a minor release, add the CI matrix entry. When raising the minimum, update `packages/go-sdk/go.mod`, all pins, CI, and documentation in the same change.
