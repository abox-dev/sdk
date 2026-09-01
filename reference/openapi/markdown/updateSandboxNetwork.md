# PUT /sandboxes/{sandboxID}/network

Update sandbox network

Update the network configuration for a running sandbox. Replaces the current egress rules with the provided configuration. Omitting field clears it.

## Parameters

- **`sandboxID`** · `string` · path · required

## Request body

Required: yes

### application/json

Schema: `SandboxNetworkUpdateConfig`

Network configuration update for a running sandbox. Replaces the current egress rules with the provided configuration. Omitting a field clears it.

- **`allowOut`** · `array<string>` · optional

  List of allowed destinations for egress traffic. Each entry can be a CIDR block (e.g. "8.8.8.8/32"), a bare IP address (e.g. "8.8.8.8"), or a domain name (e.g. "example.com", "*.example.com"). Allowed entries always take precedence over denied entries.

- **`denyOut`** · `array<string>` · optional

  List of denied CIDR blocks or IP addresses for egress traffic. Domain names are not supported for deny rules.

- **`rules`** · `object` · optional

  Per-domain transform rules. Replaces all existing rules when provided.

- **`rules.*`** · `array<SandboxNetworkRule>` · additional property

- **`rules.*.transform`** · `SandboxNetworkTransform` · optional

  Transformations applied to matching egress requests before forwarding.

- **`rules.*.transform.headers`** · `object` · optional

  HTTP headers to inject or override in matching requests. An existing header with the same name is replaced. Values are plain strings; secret resolution happens client-side before sending to the API.

- **`rules.*.transform.headers.*`** · `string` · additional property

- **`allow_internet_access`** · `boolean` · optional

  Allow sandbox to access the internet. When set to false, it behaves the same as specifying denyOut to 0.0.0.0/0 in the network config.

## Responses

### 204

Successfully updated the sandbox network configuration

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 404

Not found

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 409

Conflict

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error
