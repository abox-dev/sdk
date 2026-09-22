# Go envd user selection

The Go SDK supports explicit per-operation users for command launch, PTY creation,
and filesystem operations. Behavior follows the TypeScript and Python SDKs;
Go uses option structs, contexts, and channels for the corresponding operations.
Explicit user selection and the unified filesystem options API are introduced in
v0.2.0.

## User semantics

- `CommandOptions.User` and `PTYOptions.User` select the OS user of a new process.
  Selection works with collecting output, callbacks, and streaming delivery.
- Filesystem `User` selects the home directory used for path resolution and the
  owner of created filesystem objects. It does not provide OS permission isolation
  for filesystem RPCs. In particular, watching `~` selects that user's home; envd
  does not switch its own UID to install a filesystem watcher.
- Empty `User`, nil options, and zero-valued options all use the template default
  on envd 0.4.0 and newer. Older envd versions receive the compatibility username
  `user`, matching omitted users in TypeScript and Python.
- Go intentionally treats an empty string as omission. TypeScript/Python can
  distinguish omission from an explicit empty string on old envd; Go has no
  separate option to suppress that compatibility default.
- Existing-process operations (`Connect`, input, stdin close, signals, process
  listing, and PTY resize) do not select or change the process user. This matches
  TypeScript and Python; `CommandConnectOptions` only configures output delivery.
- Every selection belongs to one request. Concurrent operations on the same
  sandbox can select different users without changing sandbox or client defaults.

`USER` and `HOME` environment variables do not replace process identity selection.
An unknown RPC username returns `AuthenticationError`; a username containing a
colon or control character returns `InvalidArgumentError` before any request.
There is no fallback or retry with another identity. Stream errors are available
through the command/watch handle according to the existing handle contract.

## Public API

All filesystem methods take an options pointer as their final argument. Pass nil
for defaults. The API intentionally replaces the v0.1.8 signatures, without
parallel `WithOptions` methods.

| Operations                                             | Options                                                     |
| ------------------------------------------------------ | ----------------------------------------------------------- |
| `Commands.Run`, `Commands.Start`                       | `*CommandOptions` with `User`                               |
| `PTY.Create`                                           | `*PTYOptions` with `User`                                   |
| `Files.Read`, `ReadBytes`, `ReadText`, `ReadTo`        | `*FileOptions`                                              |
| `Files.Stat`, `Exists`, `MakeDir`, `Rename`, `Remove`  | `*FileOptions`                                              |
| `Files.List`                                           | `*ListFilesOptions`: `User`, `Depth` (zero means one level) |
| `Files.Write`, `WriteText`, `WriteBytes`, `WriteBatch` | `*WriteFileOptions`                                         |
| `Files.Watch`                                          | `*WatchOptions` with `User`                                 |
| `Files.SignedReadURL`, `SignedWriteURL`                | `*FileURLOptions`: `User`, `Expiration`                     |

For example:

```go
result, err := sandbox.Commands.Run(ctx, "id", &agentbox.CommandOptions{
    User: "root",
    Args: []string{"-un"},
})
info, err := sandbox.Files.Stat(ctx, "~", &agentbox.FileOptions{User: "root"})
entries, err := sandbox.Files.List(ctx, "~", &agentbox.ListFilesOptions{
    User: "root",
    Depth: 2,
})
watcher, err := sandbox.Files.Watch(ctx, "~", &agentbox.WatchOptions{User: "root"})
```

`ReadTo` takes `(ctx, path, writer, options)`. Signed URL methods take
`(path, options)`; a zero `Expiration` preserves non-expiring URL behavior.
`WriteBatch` applies the shared user and upload timeout to each file separately;
per-file metadata overrides common metadata keys without modifying caller maps.

Consumers migrating from v0.1.8 replace positional usernames with options, pass
nil to filesystem RPCs that previously had no options, and move listing depth
and signed URL expiration into their option structs.

## Transport and scope

Envd RPCs use `Authorization: Basic base64(username:)`, as implemented by
[TypeScript](../packages/js-sdk/src/envd/rpc.ts) and
[Python](../packages/python-sdk/agentbox/envd/utils.py).
HTTP file operations continue to use the `username` query parameter; signed URLs
include it in the signature. Access tokens, traffic tokens, and sandbox-routing
headers remain independent of the selected user.

The implementation sets identity on individual envd requests, not on the shared
HTTP client or control-plane client. It uses the existing envd wire contract and
does not change generated clients, protobuf schemas, or backend behavior.

## Validation

- Hermetic request tests cover command launch, PTY, every filesystem method,
  signed URLs, nil/empty options, and old/new envd defaults.
- Invalid names fail before transport; unknown RPC users produce typed errors
  without retry. Concurrent collecting/streaming commands return their own user.
- Scope tests cover process reattachment/control, platform requests, and unrelated
  HTTP requests. Existing watch cancellation/body-close tests use an explicit user.
- `TestEnvdUsersKVM` checks actual command/PTY identity for `root` and `user`,
  selected-user home paths, ownership, watch events, and file CRUD in a real
  sandbox. The test cleans up its sandbox.

The independent [command streaming](go-command-streaming.md) contract continues to
apply. Consumers can combine user selection and streaming in the same options.
