[agentbox-sdk-monorepo](../README.md) / SandboxNetworkOpts

# Type Alias: SandboxNetworkOpts

> **SandboxNetworkOpts** = `object`

## Properties

### allowOut?

> `optional` **allowOut?**: [`SandboxNetworkSelector`](SandboxNetworkSelector.md)

Allow outbound traffic from the sandbox to the specified addresses.
If `allowOut` is not specified, all outbound traffic is allowed.

Accepts either a static array of CIDR blocks, IP addresses, or hostnames,
or a callback that receives `{ allTraffic, rules }` and returns the same.
`allTraffic` is `'0.0.0.0/0'`; `rules` is a `Map` view of
[SandboxNetworkOpts.rules](#rules).

Examples:
- Static list: `["1.1.1.1", "8.8.8.0/24"]`
- Allow only rule-registered hosts:
  `({ rules }) => [...rules.keys()]`

***

### allowPublicTraffic?

> `optional` **allowPublicTraffic?**: `boolean`

Specify if the sandbox URLs should be accessible only with authentication.

#### Default

```ts
true
```

***

### denyOut?

> `optional` **denyOut?**: [`SandboxNetworkSelector`](SandboxNetworkSelector.md)

Deny outbound traffic from the sandbox to the specified addresses.

Accepts the same shapes as [allowOut](#allowout).

Examples:
- Static list: `["1.1.1.1", "8.8.8.0/24"]`
- Block all egress: `({ allTraffic }) => [allTraffic]`

***

### maskRequestHost?

> `optional` **maskRequestHost?**: `string`

Specify host mask which will be used for all sandbox requests in the header.
You can use the ${PORT} variable that will be replaced with the actual port number of the service.

#### Default

```ts
${PORT}-sandboxid.agentbox-runtime.ru
```

***

### rules?

> `optional` **rules?**: [`SandboxNetworkRules`](SandboxNetworkRules.md)

Per-domain transform rules applied to matching egress HTTP/HTTPS
requests. Keys are domains (e.g. `"api.example.com"`); values are
ordered lists of rules.

Registering a host here does not allow egress on its own — the host must
also appear in [allowOut](#allowout). Hosts registered here are exposed to the
`allowOut`/`denyOut` callbacks via `rules`.

A rule's `transform` can also be a callback receiving a
[SandboxNetworkTransformContext](SandboxNetworkTransformContext.md), which is how a workload identity
token from [SandboxOpts.iam](../interfaces/SandboxOpts.md#iam) gets injected without the SDK ever
seeing its value.

#### Example

```ts
await Sandbox.create({
  network: {
    allowOut: ({ rules }) => [...rules.keys()],
    rules: {
      'api.openai.com': [
        { transform: { headers: { Authorization: `Bearer ${token}` } } },
      ],
      'api.internal.example.com': [
        {
          transform: ({ iam }) => ({
            headers: { Authorization: `Bearer ${iam.tokens.aws}` },
          }),
        },
      ],
    },
  },
})
```
