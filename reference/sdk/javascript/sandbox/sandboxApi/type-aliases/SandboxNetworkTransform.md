[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxNetworkTransform

# Type Alias: SandboxNetworkTransform

> **SandboxNetworkTransform** = `object`

Transform applied to egress requests matching a [SandboxNetworkRule](SandboxNetworkRule.md).

## Properties

### headers?

> `optional` **headers?**: `Record`\<`string`, `string`\>

Headers to inject into the outbound request. Values override any headers
already present on the request.
