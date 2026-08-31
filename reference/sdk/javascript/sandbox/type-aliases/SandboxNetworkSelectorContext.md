[agentbox-sdk-monorepo](../README.md) / SandboxNetworkSelectorContext

# Type Alias: SandboxNetworkSelectorContext

> **SandboxNetworkSelectorContext** = `object`

Context passed to [SandboxNetworkOpts.allowOut](SandboxNetworkOpts.md#allowout) and
[SandboxNetworkOpts.denyOut](SandboxNetworkOpts.md#denyout) when they are defined as functions.

## Properties

### allTraffic

> **allTraffic**: `string`

All traffic sentinel — equivalent to `'0.0.0.0/0'`.

***

### rules

> **rules**: `Map`\<`string`, [`SandboxNetworkRule`](SandboxNetworkRule.md)[]\>

Rules registered in [SandboxNetworkOpts.rules](SandboxNetworkOpts.md#rules).
