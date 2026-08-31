[agentbox-sdk-monorepo](../README.md) / SandboxNetworkRule

# Type Alias: SandboxNetworkRule

> **SandboxNetworkRule** = `object`

Per-domain rule applied to egress requests.

## Properties

### transform?

> `optional` **transform?**: [`SandboxNetworkTransform`](SandboxNetworkTransform.md) \| [`SandboxNetworkTransformResolver`](SandboxNetworkTransformResolver.md)

Transform applied to requests matching this rule.

Accepts either a static object or a callback that receives a
[SandboxNetworkTransformContext](SandboxNetworkTransformContext.md) of placeholder strings — use the
callback to inject a workload identity token the proxy mints per request.

#### Example

```ts
{
  transform: ({ iam }) => ({
    headers: { Authorization: `Bearer ${iam.tokens.aws}` },
  }),
}
```
