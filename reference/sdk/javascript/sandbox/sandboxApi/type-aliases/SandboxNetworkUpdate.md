[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxNetworkUpdate

# Type Alias: SandboxNetworkUpdate

> **SandboxNetworkUpdate** = `object`

Subset of [SandboxNetworkOpts](SandboxNetworkOpts.md) accepted by [SandboxApi.updateNetwork](../classes/SandboxApi.md#updatenetwork).
The update endpoint replaces all egress rules atomically — fields that are
omitted are cleared on the server.

## Properties

### allowInternetAccess?

> `optional` **allowInternetAccess?**: `boolean`

Allow sandbox to access the internet. When set to `false`, it behaves the
same as specifying `denyOut: ['0.0.0.0/0']` in the network config.

***

### allowOut?

> `optional` **allowOut?**: [`SandboxNetworkSelector`](SandboxNetworkSelector.md)

See [SandboxNetworkOpts.allowOut](SandboxNetworkOpts.md#allowout).

***

### denyOut?

> `optional` **denyOut?**: [`SandboxNetworkSelector`](SandboxNetworkSelector.md)

See [SandboxNetworkOpts.denyOut](SandboxNetworkOpts.md#denyout).

***

### rules?

> `optional` **rules?**: [`SandboxNetworkRules`](SandboxNetworkRules.md)

See [SandboxNetworkOpts.rules](SandboxNetworkOpts.md#rules). A `transform` callback works here
too, but the update payload carries no `iam` config, so token names cannot
be checked against the sandbox's registered tokens — every name resolves to
its placeholder and a typo only surfaces at the destination.
