# Upstream attribution

AgentBox SDK is derived from the E2B SDK repository maintained by FOUNDRYLABS, INC.:

- Source: https://github.com/e2b-dev/E2B
- Source commit: `5a56c87e9db0e221b138662805af7743e75f1082`
- Imported without Git history and adapted for the AgentBox public API and runtime protocol.

The root repository retains the upstream Apache License 2.0. The JavaScript, Python, CLI, and Code Interpreter packages retain their upstream MIT license and FOUNDRYLABS, INC. copyright notice, with the RetailDriver LLC notice added for the AgentBox changes.

The build-only Python OpenAPI generator remains pinned under its upstream package name because replacing the generator would not change the generated runtime API. It is used only inside the codegen container and is not included in wheel, sdist, or runtime dependency metadata.
