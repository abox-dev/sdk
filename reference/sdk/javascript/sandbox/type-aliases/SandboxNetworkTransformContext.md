[agentbox-sdk-monorepo](../README.md) / SandboxNetworkTransformContext

# Type Alias: SandboxNetworkTransformContext

> **SandboxNetworkTransformContext** = `object`

Context passed to a [SandboxNetworkRule](SandboxNetworkRule.md) `transform` callback. Its
values are literal placeholder strings that the egress proxy resolves per
request, so the secret itself never leaves the platform.

## Properties

### iam

> **iam**: `object`

Workload identity placeholders.

#### tokens

> **tokens**: `Record`\<`string`, `string`\>

Placeholder for each workload token registered in
[SandboxOpts.iam](../interfaces/SandboxOpts.md#iam), keyed by token name. `tokens.aws` is the string
`'${agentbox.identity.tokens.aws}'`, which the egress proxy replaces with a
freshly minted token when it forwards the request.

Reading a name that is not registered throws
InvalidArgumentError — the proxy never turns an unregistered name
into a token, so a typo would surface as a confusing auth failure at the
destination. The four names the runtime reads off any object it
serializes, awaits or coerces (`toJSON`, `then`, `toString`, `valueOf`)
throw on use rather than on the read.
