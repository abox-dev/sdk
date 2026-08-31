[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxOnTimeout

# Type Alias: SandboxOnTimeout

> **SandboxOnTimeout** = `"pause"` \| `"kill"` \| \{ `action`: `"pause"`; `keepMemory?`: `boolean`; \} \| \{ `action`: `"kill"`; \}

What happens when the sandbox timeout is reached. Either the bare action
(`'pause'` / `'kill'`), or an object form that also controls the pause
snapshot kind via `keepMemory`.

The object form is a discriminated union on `action`: `keepMemory` is only
accepted alongside `action: 'pause'`. Passing `keepMemory` with
`action: 'kill'` is a compile-time type error.

## Union Members

`"pause"`

***

`"kill"`

***

### Type Literal

\{ `action`: `"pause"`; `keepMemory?`: `boolean`; \}

#### action

> **action**: `"pause"`

Auto-pause the sandbox when the timeout is reached.

#### keepMemory?

> `optional` **keepMemory?**: `boolean`

Whether the timeout auto-pause keeps a full memory snapshot.

When `false`, the auto-pause drops the in-memory state and persists only
the filesystem (a filesystem-only snapshot); resuming such a sandbox
cold-boots (reboots) it from disk, losing running processes and open
connections.

Cannot be combined with `autoResume`: auto-resume wakes a paused sandbox
on inbound traffic by restoring its memory snapshot in place, so the
request that woke it hits an already-running process. A filesystem-only
snapshot has no memory to restore — resuming cold-boots it — so it can't
be woken transparently by traffic and must be resumed explicitly via
`connect()`.

Left unset, the flag is omitted from the create request and the API's own
default (currently enabled) applies.

***

### Type Literal

\{ `action`: `"kill"`; \}

#### action

> **action**: `"kill"`

Kill the sandbox when the timeout is reached.
