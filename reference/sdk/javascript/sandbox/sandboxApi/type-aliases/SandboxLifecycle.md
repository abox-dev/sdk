[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxLifecycle

# Type Alias: SandboxLifecycle

> **SandboxLifecycle** = `object`

## Properties

### autoResume?

> `optional` **autoResume?**: `boolean`

Auto-resume enabled flag.

Leave unset to let the API pick the behavior. Set `false` to opt out
explicitly and keep auto-resume off even if the API's default changes.
Can be `true` only when `onTimeout` is `pause`. Not supported when
`keepMemory` is `false` (a filesystem-only snapshot must be resumed
explicitly via `connect()`).

***

### onTimeout

> **onTimeout**: [`SandboxOnTimeout`](SandboxOnTimeout.md)

Action to take when sandbox timeout is reached. Accepts either `'pause'` /
`'kill'`, or `{ action, keepMemory }` to also control the pause snapshot kind.
Omitted from the create request when unset, leaving the API's default
(currently `kill`) in effect.
