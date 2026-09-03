# Go SDK parity matrix

This matrix tracks the union of the JavaScript and Python SDK contracts. Python
sync/async duplicates map to one context-aware Go contract plus cancellation,
concurrency, and race coverage.

| Contract area | JavaScript/Python source suites | Hermetic Go coverage | KVM coverage |
| --- | --- | --- | --- |
| Configuration precedence, custom headers, proxy, pooling, User-Agent, logging, request timeouts | core connection/config tests | `config_test.go` | release smoke |
| Authentication, rate limits, HTTP/Connect error mapping, cancellation | core API/envd error tests | `config_test.go`, `sandbox_test.go`, `envd_test.go` | release smoke |
| Create, connect/resume, list/filter/page, info, kill, timeout, keepalive, pause | sandbox API tests | `sandbox_test.go` | `integration/sdk_test.go` |
| Network rules, IAM payloads, metrics, structured logs | sandbox network/IAM/metrics/log tests | `sandbox_test.go` | `integration/sdk_test.go` |
| Forks, snapshots, signed upload/download URLs | sandbox fork/snapshot/signature tests | `sandbox_test.go` | core KVM lifecycle |
| Foreground/background commands, attach/list, stdin/EOF, signals, output streams, exit errors | command and command-handle tests | `envd_test.go` | `integration/sdk_test.go` |
| PTY create/attach/input/resize/kill | PTY tests | `envd_test.go` | KVM command transport |
| Text/binary/stream reads and writes, batch writes, list/stat/metadata/exists/mkdir/move/remove/watch | filesystem and watch-handle tests | `envd_test.go` | `integration/sdk_test.go` |
| Base images/templates, private registries, Dockerfile parsing, copy, packages, env/user/workdir/start/ready/cache | template builder/parser tests | `template_test.go` | `integration/sdk_test.go` |
| Build request/upload/start/poll/log/status, visibility, tags, list/info/delete | template API/build tests | `template_test.go` | `integration/sdk_test.go` |
| Code execution languages, contexts, env vars, callbacks, request/execution timeouts | Code Interpreter sandbox tests | `codeinterpreter/client_test.go` | `integration/sdk_test.go` |
| Execution/result/log/error/raw MIME models, charts, unknown fields | Code Interpreter messaging/chart tests | `codeinterpreter/client_test.go` | Code Interpreter KVM smoke |
| Consumer compatibility and release installation | package artifact tests | `scripts/test-go-consumer.sh` | `../../tests/runtime/go/main.go` |

The CI coverage gate is 90% statement coverage across handwritten Go packages;
`internal/gen` is excluded. Every supported Go minor runs the complete hermetic
suite, and the newest minor additionally runs the race detector.
