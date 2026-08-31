[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxNetworkRuleInfo

# Type Alias: SandboxNetworkRuleInfo

> **SandboxNetworkRuleInfo** = `object`

Per-domain rule as returned by the sandbox info endpoint. Mirrors
[SandboxNetworkRule](SandboxNetworkRule.md) but with `transform` always materialized to the
static [SandboxNetworkTransform](SandboxNetworkTransform.md) shape — no callback variant.

## Properties

### transform?

> `optional` **transform?**: [`SandboxNetworkTransform`](SandboxNetworkTransform.md)
