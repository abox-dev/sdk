[agentbox-sdk-monorepo](../README.md) / SandboxNetworkSelector

# Type Alias: SandboxNetworkSelector

> **SandboxNetworkSelector** = `string`[] \| ((`ctx`) => `string`[])

Egress rule list, either a static array of CIDR blocks / IP addresses /
hostnames, or a callback that receives `{ allTraffic, rules }` and returns
the same.
