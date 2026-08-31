[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxNetworkInfo

# Type Alias: SandboxNetworkInfo

> **SandboxNetworkInfo** = `object`

Network configuration as returned by the sandbox info endpoint. Mirrors
[SandboxNetworkOpts](SandboxNetworkOpts.md) but with `allowOut`/`denyOut` always materialized
to plain string arrays.

## Properties

### allowOut?

> `optional` **allowOut?**: `string`[]

***

### allowPublicTraffic?

> `optional` **allowPublicTraffic?**: `boolean`

***

### denyOut?

> `optional` **denyOut?**: `string`[]

***

### maskRequestHost?

> `optional` **maskRequestHost?**: `string`

***

### rules?

> `optional` **rules?**: `Record`\<`string`, [`SandboxNetworkRuleInfo`](SandboxNetworkRuleInfo.md)[]\>
