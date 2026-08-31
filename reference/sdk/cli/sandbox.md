# Sandbox commands

## agentbox sandbox

work with sandboxes

### Usage

```bash
agentbox sandbox [options] [command]
```

## agentbox sandbox connect

connect terminal to already running sandbox

### Usage

```bash
agentbox sandbox connect [options] <sandboxID>
```

### Options

- `-u, --user &lt;user&gt;`: user to start the terminal session as
- `-c, --cwd &lt;dir&gt;`: working directory for the terminal session
- `-e, --env &lt;KEY=VALUE&gt;`: set environment variable for the terminal session (repeatable)

## agentbox sandbox create

create sandbox and connect terminal to it

### Usage

```bash
agentbox sandbox create [options] [template]
```

### Options

- `-d, --detach`: create sandbox without connecting terminal to it
- `--lifecycle.ontimeout &lt;action&gt;`: action when sandbox timeout is reached: pause or kill
- `--lifecycle.autoresume`: enable sandbox auto-resume, requires --lifecycle.ontimeout pause
- `--timeout &lt;seconds&gt;`: sandbox timeout in seconds
- `-u, --user &lt;user&gt;`: user to start the terminal session as
- `-c, --cwd &lt;dir&gt;`: working directory for the terminal session
- `-e, --env &lt;KEY=VALUE&gt;`: set environment variable for the terminal session (repeatable)

## agentbox sandbox exec

execute a command in a running sandbox

### Usage

```bash
agentbox sandbox exec [options] <sandboxID> <command...>
```

### Options

- `-b, --background`: run in background and return immediately
- `-c, --cwd &lt;dir&gt;`: working directory
- `-u, --user &lt;user&gt;`: run as specified user
- `-e, --env &lt;KEY=VALUE&gt;`: set environment variable (repeatable)

## agentbox sandbox info

show information for a sandbox

### Usage

```bash
agentbox sandbox info [options] <sandboxID>
```

### Options

- `-f, --format &lt;format&gt;`: output format, eg. json, pretty

## agentbox sandbox kill

kill sandbox

### Usage

```bash
agentbox sandbox kill [options] [sandboxIDs...]
```

### Options

- `-a, --all`: kill all sandboxes
- `-s, --state &lt;state&gt;`: when used with -a/--all flag, filter by state, eg. running, paused. Defaults to running
- `-m, --metadata &lt;metadata&gt;`: when used with -a/--all flag, filter by metadata, eg. key1=value1

## agentbox sandbox list

list all sandboxes, by default it list only running ones

### Usage

```bash
agentbox sandbox list [options]
```

### Options

- `-s, --state &lt;state&gt;`: filter by state, eg. running, paused. Defaults to running
- `-m, --metadata &lt;metadata&gt;`: filter by metadata, eg. key1=value1
- `-l, --limit &lt;limit&gt;`: limit the number of sandboxes returned (default: 1000, 0 for no limit)
- `-f, --format &lt;format&gt;`: output format, eg. json, pretty

## agentbox sandbox logs

show logs for sandbox

### Usage

```bash
agentbox sandbox logs [options] <sandboxID>
```

### Options

- `--level &lt;level&gt;`: filter logs by level (DEBUG, INFO, WARN, ERROR). The logs with the higher levels will be also shown. Default: `INFO`.
- `-f, --follow`: keep streaming logs until the sandbox is closed
- `--format &lt;format&gt;`: specify format for printing logs (json, pretty) Default: `pretty`.
- `--loggers [loggers]`: filter logs by loggers. Specify multiple loggers by separating them with a comma.

## agentbox sandbox metrics

show metrics for sandbox

### Usage

```bash
agentbox sandbox metrics [options] <sandboxID>
```

### Options

- `-f, --follow`: keep streaming metrics until the sandbox is closed
- `--format &lt;format&gt;`: specify format for printing metrics (json, pretty) Default: `pretty`.

## agentbox sandbox pause

pause sandbox

### Usage

```bash
agentbox sandbox pause [options] <sandboxID>
```

## agentbox sandbox resume

resume paused sandbox

### Usage

```bash
agentbox sandbox resume [options] <sandboxID>
```

## agentbox sandbox snapshot

work with sandbox snapshots

### Usage

```bash
agentbox sandbox snapshot [options] [command]
```

## agentbox sandbox snapshot create

create a snapshot from a sandbox

### Usage

```bash
agentbox sandbox snapshot create [options] <sandboxID>
```

### Options

- `-n, --name &lt;name&gt;`: name for the snapshot template, reuses the existing template if it already exists

## agentbox sandbox snapshot delete

delete snapshots

### Usage

```bash
agentbox sandbox snapshot delete [options] <snapshotIDs...>
```

## agentbox sandbox snapshot list

list snapshots

### Usage

```bash
agentbox sandbox snapshot list [options] [sandboxID]
```

### Options

- `-n, --name &lt;name&gt;`: filter by snapshot name or ID, optionally tag-qualified, eg. my-snapshot:v1
- `-l, --limit &lt;limit&gt;`: limit the number of snapshots returned (default: 1000, 0 for no limit)
- `-f, --format &lt;format&gt;`: output format, eg. json, pretty

