[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxNetworkRules

# Type Alias: SandboxNetworkRules

> **SandboxNetworkRules** = `Record`\<`string`, [`SandboxNetworkRule`](SandboxNetworkRule.md)[]\> \| `Map`\<`string`, [`SandboxNetworkRule`](SandboxNetworkRule.md)[]\>

Map of host (or CIDR / IP) to ordered list of rules applied to outbound
requests for that host. Accepts either a plain object or a `Map`.
Registering a host here does not allow egress on its own — the host must
also appear in [SandboxNetworkOpts.allowOut](SandboxNetworkOpts.md#allowout).
