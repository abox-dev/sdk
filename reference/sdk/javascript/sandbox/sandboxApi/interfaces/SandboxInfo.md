[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxInfo

# Interface: SandboxInfo

Information about a sandbox.

## Properties

### allowInternetAccess?

> `optional` **allowInternetAccess?**: `boolean`

Whether internet access was explicitly enabled or disabled for the sandbox.

***

### cpuCount

> **cpuCount**: `number`

Sandbox CPU count.

***

### endAt

> **endAt**: `Date`

Sandbox expiration date.

***

### envdVersion

> **envdVersion**: `string`

Envd version.

***

### lifecycle?

> `optional` **lifecycle?**: [`SandboxInfoLifecycle`](../type-aliases/SandboxInfoLifecycle.md)

Sandbox lifecycle configuration.

***

### memoryMB

> **memoryMB**: `number`

Sandbox Memory size in MiB.

***

### metadata

> **metadata**: `Record`\<`string`, `string`\>

Saved sandbox metadata.

***

### name?

> `optional` **name?**: `string`

Template name.

***

### network?

> `optional` **network?**: [`SandboxNetworkInfo`](../type-aliases/SandboxNetworkInfo.md)

Sandbox network configuration.

***

### sandboxDomain?

> `optional` **sandboxDomain?**: `string`

Sandbox domain.

***

### sandboxId

> **sandboxId**: `string`

Sandbox ID.

***

### startedAt

> **startedAt**: `Date`

Sandbox start time.

***

### state

> **state**: [`SandboxState`](../type-aliases/SandboxState.md)

Sandbox state.

#### String

can be `running` or `paused`

***

### templateId

> **templateId**: `string`

Template ID.
