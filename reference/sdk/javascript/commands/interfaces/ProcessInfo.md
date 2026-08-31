[agentbox-sdk-monorepo](../README.md) / ProcessInfo

# Interface: ProcessInfo

Information about a command, PTY session or start command running in the sandbox as process.

## Properties

### args

> **args**: `string`[]

Command arguments.

***

### cmd

> **cmd**: `string`

Command that was executed.

***

### cwd?

> `optional` **cwd?**: `string`

Executed command working directory.

***

### envs

> **envs**: `Record`\<`string`, `string`\>

Environment variables used for the command.

***

### pid

> **pid**: `number`

Process ID.

***

### tag?

> `optional` **tag?**: `string`

Custom tag used for identifying special commands like start command in the custom template.
