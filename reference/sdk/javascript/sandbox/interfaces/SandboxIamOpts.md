[agentbox-sdk-monorepo](../README.md) / SandboxIamOpts

# Interface: SandboxIamOpts

Sandbox workload identity configuration. A non-empty `tokens` map enables
workload identity for the sandbox.

## Properties

### tokens?

> `optional` **tokens?**: `Record`\<`string`, [`SandboxIamToken`](SandboxIamToken.md)\>

Named workload-token definitions, keyed by a caller-chosen token name.
Each value contains the token `audience` and `tokenType`.

A name is interpolated into the `'${agentbox.identity.tokens.<name>}'`
placeholder a network transform resolves, so it cannot be empty or contain
`{`, `}` or control characters.
