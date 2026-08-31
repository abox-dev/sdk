[agentbox-sdk-monorepo](../README.md) / AgentBox

# Class: AgentBox

AgentBox client with an explicitly bound connection configuration.

The resources exposed by the client ([AgentBox.Sandbox](#sandbox) and
[AgentBox.Template](#template)) behave exactly like the top-level `Sandbox` and `Template` exports,
except the options passed to the client are used as the defaults instead of
the environment variables.
Per-call options still take precedence over the client's options.

Multiple clients are fully isolated from each other and from the top-level
env-configured exports.

## Example

```ts
import { AgentBox } from '@abox-dev/sdk'

const client = new AgentBox({ apiKey: 'ab_...' })

const sandbox = await client.Sandbox.create()
await client.Template.build(client.Template().fromPythonImage('3'), 'my-env')
```

## Constructors

### Constructor

> **new AgentBox**(`opts?`): `AgentBox`

Create a new client with the connection options bound to it.

#### Parameters

##### opts?

[`AgentBoxClientOpts`](../type-aliases/AgentBoxClientOpts.md)

connection options used as the defaults for every call made
  through this client's resource classes.

#### Returns

`AgentBox`

## Properties

### Sandbox

> `readonly` **Sandbox**: *typeof* [`Sandbox`](Sandbox.md)

`Sandbox` class bound to this client's connection configuration.

***

### Template

> `readonly` **Template**: `CallableTemplate`\<*typeof* `TemplateBase`\>

`Template` bound to this client's connection configuration. Both the
builder (`client.Template()`) and the statics
(`client.Template.build(...)`, `client.Template.exists(...)`, …) work like
the top-level `Template`.
