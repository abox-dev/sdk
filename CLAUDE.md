Use pnpm for JavaScript packages and uv for Python packages.
Use English exclusively in source code, comments, documentation, commit messages, and GitHub pull request titles and descriptions.
Keep the JavaScript and Python SDKs, including sync and async Python APIs, behaviorally aligned.
Run format checks, lint, type checks, unit tests, deterministic generation, builds, and package-install checks before committing.
The API and envd snapshots under spec/ are generated from mono/infra. Do not edit them manually. Update them with `make sync-specs MONO_DIR=/path/to/mono`, then run `make generate`.
Generated clients must depend only on checked-in snapshots and never fetch network content during generation.
Public APIs, package artifacts, examples, errors, environment variables, and headers must use AgentBox naming. Upstream names are allowed only in licenses, attribution, pinned build-only codegen tooling, and wire/protobuf namespaces that are required by the runtime protocol.
Default development credentials may be stored in `.env.local` or `~/.agentbox/config.json`; never print or commit them.
