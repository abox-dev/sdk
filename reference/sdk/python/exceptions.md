## `agentbox.exceptions`

**Classes:**

- [**AuthenticationException**](#agentbox.exceptions.AuthenticationException) – Raised when authentication fails.
- [**BuildException**](#agentbox.exceptions.BuildException) – Raised when the build fails.
- [**FileNotFoundException**](#agentbox.exceptions.FileNotFoundException) – Raised when a file or directory is not found inside a sandbox.
- [**FileUploadException**](#agentbox.exceptions.FileUploadException) – Raised when the file upload fails.
- [**InvalidArgumentException**](#agentbox.exceptions.InvalidArgumentException) – Raised when an invalid argument is provided.
- [**NotEnoughSpaceException**](#agentbox.exceptions.NotEnoughSpaceException) – Raised when there is not enough disk space.
- [**RateLimitException**](#agentbox.exceptions.RateLimitException) – Raised when the API rate limit is exceeded.
- [**SandboxException**](#agentbox.exceptions.SandboxException) – Base class for all sandbox errors.
- [**SandboxNotFoundException**](#agentbox.exceptions.SandboxNotFoundException) – Raised when a sandbox is not found (e.g. it doesn't exist or is no longer running).
- [**TemplateException**](#agentbox.exceptions.TemplateException) – Exception raised when the template uses old envd version. It isn't compatible with the new SDK.
- [**TimeoutException**](#agentbox.exceptions.TimeoutException) – Raised when a timeout occurs.

**Functions:**

- [**format_request_timeout_error**](#agentbox.exceptions.format_request_timeout_error) – 
- [**format_sandbox_timeout_exception**](#agentbox.exceptions.format_sandbox_timeout_exception) – 

### `agentbox.exceptions.AuthenticationException`

Bases: <code>[Exception](#Exception)</code>

Raised when authentication fails.

### `agentbox.exceptions.BuildException`

Bases: <code>[Exception](#Exception)</code>

Raised when the build fails.

### `agentbox.exceptions.FileNotFoundException`

Bases: <code>[SandboxException](#agentbox.exceptions.SandboxException)</code>

Raised when a file or directory is not found inside a sandbox.

### `agentbox.exceptions.FileUploadException`

Bases: <code>[BuildException](#agentbox.exceptions.BuildException)</code>

Raised when the file upload fails.

### `agentbox.exceptions.InvalidArgumentException`

Bases: <code>[SandboxException](#agentbox.exceptions.SandboxException)</code>

Raised when an invalid argument is provided.

### `agentbox.exceptions.NotEnoughSpaceException`

Bases: <code>[SandboxException](#agentbox.exceptions.SandboxException)</code>

Raised when there is not enough disk space.

### `agentbox.exceptions.RateLimitException`

Bases: <code>[SandboxException](#agentbox.exceptions.SandboxException)</code>

Raised when the API rate limit is exceeded.

### `agentbox.exceptions.SandboxException`

Bases: <code>[Exception](#Exception)</code>

Base class for all sandbox errors.

Raised when a general sandbox exception occurs.

### `agentbox.exceptions.SandboxNotFoundException`

Bases: <code>[SandboxException](#agentbox.exceptions.SandboxException)</code>

Raised when a sandbox is not found (e.g. it doesn't exist or is no longer running).

### `agentbox.exceptions.TemplateException`

Bases: <code>[SandboxException](#agentbox.exceptions.SandboxException)</code>

Exception raised when the template uses old envd version. It isn't compatible with the new SDK.

### `agentbox.exceptions.TimeoutException`

Bases: <code>[SandboxException](#agentbox.exceptions.SandboxException)</code>

Raised when a timeout occurs.

The `unavailable` exception type is caused by sandbox timeout.

The `canceled` exception type is caused by exceeding request timeout.

The `deadline_exceeded` exception type is caused by exceeding the timeout for process, watch, etc.

The `unknown` exception type is sometimes caused by the sandbox timeout when the request is not processed correctly.

### `agentbox.exceptions.format_request_timeout_error`

```python
format_request_timeout_error()
```

### `agentbox.exceptions.format_sandbox_timeout_exception`

```python
format_sandbox_timeout_exception(message)
```

