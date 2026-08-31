# GET /sandboxes/{sandboxID}

Sandbox

Get a sandbox by id

## Parameters

- **`sandboxID`** · `string` · path · required

## Responses

### 200

Successfully returned the sandbox

Content-Type: `application/json`

Schema: `SandboxDetail`

- **`templateID`** · `string` · required

  Identifier of the template from which is the sandbox created

- **`alias`** · `string` · optional

  Alias of the template

- **`sandboxID`** · `string` · required

  Identifier of the sandbox

- **`startedAt`** · `string` · required

  Time when the sandbox was started

- **`endAt`** · `string` · required

  Time when the sandbox will expire

- **`envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

- **`allowInternetAccess`** · `boolean` · optional

  Whether internet access was explicitly enabled or disabled for the sandbox. Null means it was not explicitly set.

- **`domain`** · `string` · optional

  Base domain where the sandbox traffic is accessible

- **`cpuCount`** · `CPUCount` · required

  CPU cores for the sandbox

- **`memoryMB`** · `MemoryMB` · required

  Memory for the sandbox in MiB

- **`diskSizeMB`** · `DiskSizeMB` · required

  Disk size for the sandbox in MiB

- **`metadata`** · `SandboxMetadata` · optional

- **`metadata.*`** · `string` · additional property

  Metadata of the sandbox

- **`state`** · `SandboxState` · required

  State of the sandbox

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

- **`lifecycle`** · `SandboxLifecycle` · optional

  Sandbox lifecycle policy returned by sandbox info.

- **`lifecycle.autoResume`** · `boolean` · required

  Whether the sandbox can auto-resume.

- **`lifecycle.onTimeout`** · `SandboxOnTimeout` · required

  Action taken when the sandbox times out.

### 404

Not found

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 401

Authentication error

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
