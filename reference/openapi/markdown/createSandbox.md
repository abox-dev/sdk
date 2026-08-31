# POST /sandboxes

Create sandbox

Create a sandbox from the template

## Request body

Required: yes

### application/json

Schema: `NewSandbox`

- **`templateID`** · `string` · required

  Identifier of the required template

- **`timeout`** · `integer` · optional

  Time to live for the sandbox in seconds.

- **`autoPause`** · `boolean` · optional

  Automatically pauses the sandbox after the timeout

- **`autoPauseMemory`** · `boolean` · optional

  Controls the snapshot kind taken when the sandbox auto-pauses on timeout (only relevant when autoPause is true). When false, the auto-pause drops the in-memory state and persists only the filesystem (a filesystem-only snapshot); resuming it cold-boots (reboots) the sandbox from disk. Such a snapshot cannot be auto-resumed by traffic and must be resumed explicitly, so it cannot be combined with autoResume. Defaults to true (full memory snapshot).

- **`autoResume`** · `SandboxAutoResumeConfig` · optional

  Auto-resume configuration for paused sandboxes.

- **`autoResume.enabled`** · `SandboxAutoResumeEnabled` · required

  Auto-resume enabled flag for paused sandboxes. Default false.

- **`secure`** · `boolean` · optional

  Secure all system communication with sandbox

- **`allow_internet_access`** · `boolean` · optional

  Allow sandbox to access the internet. When set to false, it behaves the same as specifying denyOut to 0.0.0.0/0 in the network config.

- **`network`** · `SandboxNetworkConfig` · optional

- **`network.allowPublicTraffic`** · `boolean` · optional

  Specify if the sandbox URLs should be accessible only with authentication.

- **`network.allowOut`** · `array<string>` · optional

  List of allowed destinations for egress traffic. Each entry can be a CIDR block (e.g. "8.8.8.8/32"), a bare IP address (e.g. "8.8.8.8"), or a domain name (e.g. "example.com", "*.example.com"). Allowed entries always take precedence over denied entries.

- **`network.denyOut`** · `array<string>` · optional

  List of denied CIDR blocks or IP addresses for egress traffic. Domain names are not supported for deny rules.

- **`network.maskRequestHost`** · `string` · optional

  Specify host mask which will be used for all sandbox requests

- **`network.rules`** · `object` · optional

  Per-domain transform rules applied to matching egress HTTP/HTTPS requests. Keys are domains (e.g. "api.example.com", "example.com"). A domain listed here is not automatically allowed - use allowOut to permit the traffic.

- **`network.rules.*`** · `array<SandboxNetworkRule>` · additional property

- **`network.rules.*.transform`** · `SandboxNetworkTransform` · optional

  Transformations applied to matching egress requests before forwarding.

- **`network.rules.*.transform.headers`** · `object` · optional

  HTTP headers to inject or override in matching requests. An existing header with the same name is replaced. Values are plain strings; secret resolution happens client-side before sending to the API.

- **`network.rules.*.transform.headers.*`** · `string` · additional property

- **`metadata`** · `SandboxMetadata` · optional

- **`metadata.*`** · `string` · additional property

  Metadata of the sandbox

- **`envVars`** · `EnvVars` · optional

- **`envVars.*`** · `string` · additional property

  Environment variables for the sandbox

- **`mcp`** · `Mcp` · optional

  MCP configuration for the sandbox

- **`mcp.*`** · `object` · additional property

- **`iam`** · `SandboxIam` · optional

  Sandbox workload identity configuration. A non-empty, valid tokens map enables workload identity for the sandbox.

- **`iam.tokens`** · `SandboxIamTokens` · optional

  Named workload-token definitions, keyed by a caller-chosen token name.

- **`iam.tokens.*`** · `SandboxIamToken` · additional property

- **`iam.tokens.*.audience`** · `string` · required

  Audience of the workload token, stored exactly as provided.

- **`iam.tokens.*.tokenType`** · `string` · required

  Workload token type.

## Responses

### 201

The sandbox was created successfully

Content-Type: `application/json`

Schema: `Sandbox`

- **`templateID`** · `string` · required

  Identifier of the template from which is the sandbox created

- **`sandboxID`** · `string` · required

  Identifier of the sandbox

- **`alias`** · `string` · optional

  Alias of the template

- **`envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

- **`domain`** · `string` · optional

  Base domain where the sandbox traffic is accessible

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error
