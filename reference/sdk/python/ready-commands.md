## `agentbox.template.readycmd`

**Classes:**

- [**ReadyCmd**](#agentbox.template.readycmd.ReadyCmd) – Wrapper class for ready check commands.

**Functions:**

- [**wait_for_file**](#agentbox.template.readycmd.wait_for_file) – Wait for a file to exist.
- [**wait_for_port**](#agentbox.template.readycmd.wait_for_port) – Wait for a port to be listening.
- [**wait_for_process**](#agentbox.template.readycmd.wait_for_process) – Wait for a process with a specific name to be running.
- [**wait_for_timeout**](#agentbox.template.readycmd.wait_for_timeout) – Wait for a specified timeout before considering the sandbox ready.
- [**wait_for_url**](#agentbox.template.readycmd.wait_for_url) – Wait for a URL to return a specific HTTP status code.

### `agentbox.template.readycmd.ReadyCmd`

```python
ReadyCmd(cmd)
```

Wrapper class for ready check commands.

**Functions:**

- [**get_cmd**](#agentbox.template.readycmd.ReadyCmd.get_cmd) – 

#### `agentbox.template.readycmd.ReadyCmd.get_cmd`

```python
get_cmd()
```

### `agentbox.template.readycmd.wait_for_file`

```python
wait_for_file(filename)
```

Wait for a file to exist.

Uses shell test command to check file existence.

:param filename: Path to the file to wait for

:return: ReadyCmd that checks for the file

Example
```python
from agentbox import Template, wait_for_file

template = (
    Template()
    .from_base_image()
    .set_start_cmd('./init.sh', wait_for_file('/tmp/ready'))
)
```

### `agentbox.template.readycmd.wait_for_port`

```python
wait_for_port(port)
```

Wait for a port to be listening.

Uses `ss` command to check if a port is open and listening.

:param port: Port number to wait for

:return: ReadyCmd that checks for the port

Example
```python
from agentbox import Template, wait_for_port

template = (
    Template()
    .from_python_image()
    .set_start_cmd('python -m http.server 8000', wait_for_port(8000))
)
```

### `agentbox.template.readycmd.wait_for_process`

```python
wait_for_process(process_name)
```

Wait for a process with a specific name to be running.

Uses `pgrep` to check if a process exists.

:param process_name: Name of the process to wait for

:return: ReadyCmd that checks for the process

Example
```python
from agentbox import Template, wait_for_process

template = (
    Template()
    .from_base_image()
    .set_start_cmd('./my-daemon', wait_for_process('my-daemon'))
)
```

### `agentbox.template.readycmd.wait_for_timeout`

```python
wait_for_timeout(timeout)
```

Wait for a specified timeout before considering the sandbox ready.

Uses `sleep` command to wait for a fixed duration.

:param timeout: Time to wait in **milliseconds** (minimum: 1000ms / 1 second)

:return: ReadyCmd that waits for the specified duration

Example
```python
from agentbox import Template, wait_for_timeout

template = (
    Template()
    .from_node_image()
    .set_start_cmd('npm start', wait_for_timeout(5000))  # Wait 5 seconds
)
```

### `agentbox.template.readycmd.wait_for_url`

```python
wait_for_url(url, status_code=200)
```

Wait for a URL to return a specific HTTP status code.

Uses `curl` to make HTTP requests and check the response status.

:param url: URL to check (e.g., 'http://localhost:3000/health')
:param status_code: Expected HTTP status code (default: 200)

:return: ReadyCmd that checks the URL

Example
```python
from agentbox import Template, wait_for_url

template = (
    Template()
    .from_node_image()
    .set_start_cmd('npm start', wait_for_url('http://localhost:3000/health'))
)
```

