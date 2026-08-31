# Template commands

## agentbox template

manage sandbox templates

### Usage

```bash
agentbox template [options] [command]
```

## agentbox template create

build Dockerfile as a Sandbox template. This command reads a Dockerfile and builds it directly.

### Usage

```bash
agentbox template create [options] <template-name>
```

### Options

- `-p, --path &lt;path&gt;`: change root directory where command is executed to &lt;path&gt; directory
- `-d, --dockerfile &lt;file&gt;`: specify path to Dockerfile. By default AgentBox tries to find agentbox.Dockerfile or Dockerfile in root directory.
- `-c, --cmd &lt;start-command&gt;`: specify command that will be executed when the sandbox is started.
- `--ready-cmd &lt;ready-command&gt;`: specify command that will need to exit 0 for the template to be ready.
- `--cpu-count &lt;cpu-count&gt;`: specify the number of CPUs that will be used to run the sandbox. The default value is 2.
- `--memory-mb &lt;memory-mb&gt;`: specify the amount of memory in megabytes that will be used to run the sandbox. Must be an even number. The default value is 1024.
- `--no-cache`: skip cache when building the template.

## agentbox template delete

delete sandbox template

### Usage

```bash
agentbox template delete [options] [template]
```

### Options

- `-s, --select`: select sandbox template from interactive list
- `-t, --project &lt;project-id&gt;`: specify the project ID for the operation. See https://docs.agentbox.ru/en/cli/configuration/.
- `-y, --yes`: skip manual delete confirmation

## agentbox template list

list sandbox templates

### Usage

```bash
agentbox template list [options]
```

### Options

- `-t, --project &lt;project-id&gt;`: specify the project ID for the operation. See https://docs.agentbox.ru/en/cli/configuration/.
- `-f, --format &lt;format&gt;`: output format, eg. json, pretty

