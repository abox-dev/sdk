[agentbox-sdk-monorepo](../README.md) / SandboxConnectOpts

# Type Alias: SandboxConnectOpts

> **SandboxConnectOpts** = `ConnectionOpts` & `object`

Options for connecting to a Sandbox.

## Type Declaration

### timeoutMs?

> `optional` **timeoutMs?**: `number`

Timeout for the sandbox in **milliseconds**.
For running sandboxes, the timeout will update only if the new timeout is longer than the existing one.
Maximum time a sandbox can be kept alive is 24 hours (86_400_000 milliseconds) for Pro users and 1 hour (3_600_000 milliseconds) for Hobby users.

#### Default

```ts
300_000 // 5 minutes
```
