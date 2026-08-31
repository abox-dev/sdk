# GET /envs

Environment variables

## Parameters

- **`Agentbox-Sandbox-Id`** · `string` · header · required

  Identifier of the sandbox that receives the request.

- **`Agentbox-Sandbox-Port`** · `integer` · header · required

  Internal envd HTTP port exposed through the sandbox proxy.

## Responses

### 200

Environment variables

Content-Type: `application/json`

Schema: `EnvVars`

- **`*`** · `string` · additional property
