## `agentbox.template.logger`

**Classes:**

- [**DefaultBuildLogger**](#agentbox.template.logger.DefaultBuildLogger) – 
- [**DefaultBuildLoggerInitialState**](#agentbox.template.logger.DefaultBuildLoggerInitialState) – 
- [**LogEntry**](#agentbox.template.logger.LogEntry) – Represents a single log entry from the template build process.
- [**LogEntryEnd**](#agentbox.template.logger.LogEntryEnd) – Special log entry indicating the end of a build process.
- [**LogEntryStart**](#agentbox.template.logger.LogEntryStart) – Special log entry indicating the start of a build process.

**Functions:**

- [**default_build_logger**](#agentbox.template.logger.default_build_logger) – Create a default build logger with animated timer display.
- [**set_interval**](#agentbox.template.logger.set_interval) – Returns a stop function that can be called to cancel the interval.

**Attributes:**

- [**DEFAULT_LEVEL**](#agentbox.template.logger.DEFAULT_LEVEL) (<code>[LogEntryLevel](#agentbox.template.logger.LogEntryLevel)</code>) – Colored labels for each log level.
- [**LogEntryLevel**](#agentbox.template.logger.LogEntryLevel) – 
- [**TIMER_UPDATE_INTERVAL_MS**](#agentbox.template.logger.TIMER_UPDATE_INTERVAL_MS) – Default minimum log level to display.
- [**level_order**](#agentbox.template.logger.level_order) – 
- [**levels**](#agentbox.template.logger.levels) (<code>[Dict](#typing.Dict)[[LogEntryLevel](#agentbox.template.logger.LogEntryLevel), [tuple](#tuple)[[str](#str), [Style](#rich.style.Style)]]</code>) – Numeric ordering of log levels for comparison (lower = less severe).

### `agentbox.template.logger.DEFAULT_LEVEL`

```python
DEFAULT_LEVEL: LogEntryLevel = 'info'
```

Colored labels for each log level.

### `agentbox.template.logger.DefaultBuildLogger`

```python
DefaultBuildLogger(min_level=None)
```

**Functions:**

- [**logger**](#agentbox.template.logger.DefaultBuildLogger.logger) – 

#### `agentbox.template.logger.DefaultBuildLogger.logger`

```python
logger(log)
```

### `agentbox.template.logger.DefaultBuildLoggerInitialState`

Bases: <code>[TypedDict](#typing.TypedDict)</code>

**Attributes:**

- [**animation_frame**](#agentbox.template.logger.DefaultBuildLoggerInitialState.animation_frame) (<code>[int](#int)</code>) – 
- [**start_time**](#agentbox.template.logger.DefaultBuildLoggerInitialState.start_time) (<code>[float](#float)</code>) – 
- [**timer**](#agentbox.template.logger.DefaultBuildLoggerInitialState.timer) (<code>[Optional](#typing.Optional)[[Callable](#typing.Callable)[[], None]]</code>) – 

#### `agentbox.template.logger.DefaultBuildLoggerInitialState.animation_frame`

```python
animation_frame: int
```

#### `agentbox.template.logger.DefaultBuildLoggerInitialState.start_time`

```python
start_time: float
```

#### `agentbox.template.logger.DefaultBuildLoggerInitialState.timer`

```python
timer: Optional[Callable[[], None]]
```

### `agentbox.template.logger.LogEntry`

```python
LogEntry(timestamp, level, message)
```

Represents a single log entry from the template build process.

**Attributes:**

- [**level**](#agentbox.template.logger.LogEntry.level) (<code>[LogEntryLevel](#agentbox.template.logger.LogEntryLevel)</code>) – 
- [**message**](#agentbox.template.logger.LogEntry.message) (<code>[str](#str)</code>) – 
- [**timestamp**](#agentbox.template.logger.LogEntry.timestamp) (<code>[datetime](#datetime.datetime)</code>) – 

#### `agentbox.template.logger.LogEntry.level`

```python
level: LogEntryLevel
```

#### `agentbox.template.logger.LogEntry.message`

```python
message: str
```

#### `agentbox.template.logger.LogEntry.timestamp`

```python
timestamp: datetime
```

### `agentbox.template.logger.LogEntryEnd`

```python
LogEntryEnd(timestamp, level, message)
```

Bases: <code>[LogEntry](#agentbox.template.logger.LogEntry)</code>

Special log entry indicating the end of a build process.

**Attributes:**

- [**level**](#agentbox.template.logger.LogEntryEnd.level) (<code>[LogEntryLevel](#agentbox.template.logger.LogEntryLevel)</code>) – 
- [**message**](#agentbox.template.logger.LogEntryEnd.message) (<code>[str](#str)</code>) – 
- [**timestamp**](#agentbox.template.logger.LogEntryEnd.timestamp) (<code>[datetime](#datetime.datetime)</code>) – 

#### `agentbox.template.logger.LogEntryEnd.level`

```python
level: LogEntryLevel = field(default='debug', init=False)
```

#### `agentbox.template.logger.LogEntryEnd.message`

```python
message: str
```

#### `agentbox.template.logger.LogEntryEnd.timestamp`

```python
timestamp: datetime
```

### `agentbox.template.logger.LogEntryLevel`

```python
LogEntryLevel = Literal['debug', 'info', 'warn', 'error']
```

### `agentbox.template.logger.LogEntryStart`

```python
LogEntryStart(timestamp, level, message)
```

Bases: <code>[LogEntry](#agentbox.template.logger.LogEntry)</code>

Special log entry indicating the start of a build process.

**Attributes:**

- [**level**](#agentbox.template.logger.LogEntryStart.level) (<code>[LogEntryLevel](#agentbox.template.logger.LogEntryLevel)</code>) – 
- [**message**](#agentbox.template.logger.LogEntryStart.message) (<code>[str](#str)</code>) – 
- [**timestamp**](#agentbox.template.logger.LogEntryStart.timestamp) (<code>[datetime](#datetime.datetime)</code>) – 

#### `agentbox.template.logger.LogEntryStart.level`

```python
level: LogEntryLevel = field(default='debug', init=False)
```

#### `agentbox.template.logger.LogEntryStart.message`

```python
message: str
```

#### `agentbox.template.logger.LogEntryStart.timestamp`

```python
timestamp: datetime
```

### `agentbox.template.logger.TIMER_UPDATE_INTERVAL_MS`

```python
TIMER_UPDATE_INTERVAL_MS = 150
```

Default minimum log level to display.

### `agentbox.template.logger.default_build_logger`

```python
default_build_logger(min_level=None)
```

Create a default build logger with animated timer display.

:param min_level: Minimum log level to display (default: 'info')

:return: Logger function that accepts LogEntry instances

Example
```python
from agentbox import Template, default_build_logger

template = Template().from_python_image()

# Use with build - implementation would be in build_async module
# await Template.build(template,
#     alias='my-template',
#     on_build_logs=default_build_logger(min_level='debug')
# )
```

### `agentbox.template.logger.level_order`

```python
level_order = {'debug': 0, 'info': 1, 'warn': 2, 'error': 3}
```

### `agentbox.template.logger.levels`

```python
levels: Dict[LogEntryLevel, tuple[str, Style]] = {
    "error": ("ERROR", Style(color="red")),
    "warn": ("WARN ", Style(color="#FF4400")),
    "info": ("INFO ", Style(color="#FF8800")),
    "debug": ("DEBUG", Style(color="bright_black")),
}

```

Numeric ordering of log levels for comparison (lower = less severe).

### `agentbox.template.logger.set_interval`

```python
set_interval(func, interval)
```

Returns a stop function that can be called to cancel the interval.

Similar to JavaScript's setInterval.

:param func: Function to execute at each interval
:param interval: Interval duration in **seconds**

:return: Stop function that can be called to cancel the interval

