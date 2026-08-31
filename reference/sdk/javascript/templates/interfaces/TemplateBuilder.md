[agentbox-sdk-monorepo](../README.md) / TemplateBuilder

# Interface: TemplateBuilder

Main builder state for constructing templates.
Provides methods for customizing the template environment.

## Methods

### addMcpServer()

> **addMcpServer**(`servers`): `TemplateBuilder`

Install MCP servers using mcp-gateway.
Note: Requires a base image with mcp-gateway pre-installed (e.g., mcp-gateway).

#### Parameters

##### servers

keyof McpServer \| keyof McpServer[]

MCP server name(s)

#### Returns

`TemplateBuilder`

#### Throws

If the base template is not mcp-gateway

#### Example

```ts
template.addMcpServer('exa')
template.addMcpServer(['brave', 'firecrawl', 'duckduckgo'])
```

***

### aptInstall()

> **aptInstall**(`packages`, `options?`): `TemplateBuilder`

Install Debian/Ubuntu packages using apt-get.

#### Parameters

##### packages

`string` \| `string`[]

Package name(s)

##### options?

###### fixMissing?

`boolean`

###### noInstallRecommends?

`boolean`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.aptInstall('vim')
template.aptInstall(['git', 'curl', 'wget'])
template.aptInstall(['vim'], { noInstallRecommends: true })
template.aptInstall(['vim'], { fixMissing: true })
```

***

### betaDevContainerPrebuild()

> **betaDevContainerPrebuild**(`devcontainerDirectory`): `TemplateBuilder`

Prebuild a devcontainer from the specified directory.

#### Parameters

##### devcontainerDirectory

`string`

Path to the devcontainer directory

#### Returns

`TemplateBuilder`

#### Example

```ts
template
 .gitClone('https://myrepo.com/project.git', '/my-devcontainer')
 .betaDevContainerPrebuild('/my-devcontainer')
```

***

### betaSetDevContainerStart()

> **betaSetDevContainerStart**(`devcontainerDirectory`): `TemplateFinal`

Start a devcontainer from the specified directory.

#### Parameters

##### devcontainerDirectory

`string`

Path to the devcontainer directory

#### Returns

`TemplateFinal`

#### Example

```ts
template
 .gitClone('https://myrepo.com/project.git', '/my-devcontainer')
 .startDevcontainer('/my-devcontainer')

// Prebuild and start
template
 .gitClone('https://myrepo.com/project.git', '/my-devcontainer')
 .betaDevContainerPrebuild('/my-devcontainer')
 // Other instructions...
 .betaSetDevContainerStart('/my-devcontainer')
```

***

### bunInstall()

> **bunInstall**(`packages?`, `options?`): `TemplateBuilder`

Install Bun packages using bun.

#### Parameters

##### packages?

`string` \| `string`[]

Package name(s) or undefined for package.json

##### options?

Install options

###### dev?

`boolean`

###### g?

`boolean`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.bunInstall('express')
template.bunInstall(['lodash', 'axios'])
template.bunInstall('tsx', { g: true })
template.bunInstall('typescript', { dev: true })
template.bunInstall()  // Installs from package.json
```

***

### copy()

> **copy**(`src`, `dest`, `options?`): `TemplateBuilder`

Copy files or directories into the template.

#### Parameters

##### src

`PathLike` \| `PathLike`[]

Source path(s)

##### dest

`PathLike`

Destination path

##### options?

Copy options

###### forceUpload?

`true`

###### gzip?

`boolean`

###### mode?

`number`

###### resolveSymlinks?

`boolean`

###### user?

`string`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.copy('requirements.txt', '/home/user/')
template.copy(['app.ts', 'config.ts'], '/app/', { mode: 0o755 })
```

***

### copyItems()

> **copyItems**(`items`): `TemplateBuilder`

Copy multiple items with individual options.

#### Parameters

##### items

[`CopyItem`](../type-aliases/CopyItem.md)[]

Array of copy items

#### Returns

`TemplateBuilder`

#### Example

```ts
template.copyItems([
  { src: 'app.ts', dest: '/app/' },
  { src: 'config.ts', dest: '/app/', mode: 0o644 }
])
```

***

### gitClone()

> **gitClone**(`url`, `path?`, `options?`): `TemplateBuilder`

Clone a Git repository.

#### Parameters

##### url

`string`

Repository URL

##### path?

`PathLike`

Optional destination path

##### options?

Clone options

###### branch?

`string`

###### depth?

`number`

###### user?

`string`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.gitClone('https://github.com/user/repo.git', '/app/repo')
template.gitClone('https://github.com/user/repo.git', undefined, {
  branch: 'main',
  depth: 1
})
template.gitClone('https://github.com/user/repo.git', '/app/repo', {
  user: 'root'
})
```

***

### makeDir()

> **makeDir**(`path`, `options?`): `TemplateBuilder`

Create directories.

#### Parameters

##### path

`PathLike` \| `PathLike`[]

Directory path(s)

##### options?

Directory options

###### mode?

`number`

###### user?

`string`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.makeDir('/app/data', { mode: 0o755 })
template.makeDir(['/app/logs', '/app/cache'])
template.makeDir('/app/data', { mode: 0o755, user: 'root' })
```

***

### makeSymlink()

> **makeSymlink**(`src`, `dest`, `options?`): `TemplateBuilder`

Create a symbolic link.

#### Parameters

##### src

`PathLike`

Source path (target)

##### dest

`PathLike`

Destination path (symlink location)

##### options?

Symlink options

###### force?

`boolean`

###### user?

`string`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.makeSymlink('/usr/bin/python3', '/usr/bin/python')
template.makeSymlink('/usr/bin/python3', '/usr/bin/python', { user: 'root' })
template.makeSymlink('/usr/bin/python3', '/usr/bin/python', { force: true })
```

***

### npmInstall()

> **npmInstall**(`packages?`, `options?`): `TemplateBuilder`

Install Node.js packages using npm.

#### Parameters

##### packages?

`string` \| `string`[]

Package name(s) or undefined for package.json

##### options?

Install options

###### dev?

`boolean`

###### g?

`boolean`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.npmInstall('express')
template.npmInstall(['lodash', 'axios'])
template.npmInstall('tsx', { g: true })
template.npmInstall('typescript', { dev: true })
template.npmInstall()  // Installs from package.json
```

***

### pipInstall()

> **pipInstall**(`packages?`, `options?`): `TemplateBuilder`

Install Python packages using pip.

#### Parameters

##### packages?

`string` \| `string`[]

Package name(s) or undefined for current directory

##### options?

Install options

###### g?

`boolean`

Install globally as root (default: true). Set to false for user-only installation with --user flag

#### Returns

`TemplateBuilder`

#### Example

```ts
template.pipInstall('numpy')  // Installs globally (default)
template.pipInstall(['pandas', 'scikit-learn'])
template.pipInstall('numpy', { g: false })  // Install for user only
template.pipInstall()  // Installs from current directory
```

***

### remove()

> **remove**(`path`, `options?`): `TemplateBuilder`

Remove files or directories.

#### Parameters

##### path

`PathLike` \| `PathLike`[]

Path(s) to remove

##### options?

Remove options

###### force?

`boolean`

###### recursive?

`boolean`

###### user?

`string`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.remove('/tmp/cache', { recursive: true, force: true })
template.remove('/tmp/cache', { recursive: true, force: true, user: 'root' })
```

***

### rename()

> **rename**(`src`, `dest`, `options?`): `TemplateBuilder`

Rename or move a file or directory.

#### Parameters

##### src

`PathLike`

Source path

##### dest

`PathLike`

Destination path

##### options?

Rename options

###### force?

`boolean`

###### user?

`string`

#### Returns

`TemplateBuilder`

#### Example

```ts
template.rename('/tmp/old.txt', '/tmp/new.txt')
template.rename('/tmp/old.txt', '/tmp/new.txt', { user: 'root' })
```

***

### runCmd()

#### Call Signature

> **runCmd**(`command`, `options?`): `TemplateBuilder`

Run a shell command.

##### Parameters

###### command

`string`

Command string

###### options?

Command options

###### user?

`string`

##### Returns

`TemplateBuilder`

##### Example

```ts
template.runCmd('apt-get update')
template.runCmd(['pip install numpy', 'pip install pandas'])
template.runCmd('apt-get install vim', { user: 'root' })
```

#### Call Signature

> **runCmd**(`commands`, `options?`): `TemplateBuilder`

Run multiple shell commands.

##### Parameters

###### commands

`string`[]

Array of command strings

###### options?

Command options

###### user?

`string`

##### Returns

`TemplateBuilder`

#### Call Signature

> **runCmd**(`commandOrCommands`, `options?`): `TemplateBuilder`

Run command(s).

##### Parameters

###### commandOrCommands

`string` \| `string`[]

Command or commands

###### options?

Command options

###### user?

`string`

##### Returns

`TemplateBuilder`

***

### setEnvs()

> **setEnvs**(`envs`): `TemplateBuilder`

Set environment variables.
Note: Environment variables defined here are available only during template build.

#### Parameters

##### envs

`Record`\<`string`, `string`\>

Environment variables

#### Returns

`TemplateBuilder`

#### Example

```ts
template.setEnvs({ NODE_ENV: 'production', PORT: '8080' })
```

***

### setReadyCmd()

> **setReadyCmd**(`readyCommand`): `TemplateFinal`

Set or update the ready check command.

#### Parameters

##### readyCommand

`string` \| `ReadyCmd`

Command to check readiness

#### Returns

`TemplateFinal`

#### Example

```ts
// Using a string command
template.setReadyCmd('curl http://localhost:8000/health')

// Using ReadyCmd helpers
import { waitForPort, waitForFile, waitForProcess } from '@abox-dev/sdk'

template.setReadyCmd(waitForPort(3000))

template.setReadyCmd(waitForFile('/tmp/ready'))

template.setReadyCmd(waitForProcess('nginx'))
```

***

### setStartCmd()

> **setStartCmd**(`startCommand`, `readyCommand`): `TemplateFinal`

Set the start command and ready check.

#### Parameters

##### startCommand

`string`

Command to run on startup

##### readyCommand

`string` \| `ReadyCmd`

Command to check readiness

#### Returns

`TemplateFinal`

#### Example

```ts
// Using a string command
template.setStartCmd(
  'node app.js',
  'curl http://localhost:8000/health'
)

// Using ReadyCmd helpers
import { waitForPort, waitForURL } from '@abox-dev/sdk'

template.setStartCmd(
  'python -m http.server 8000',
  waitForPort(8000)
)

template.setStartCmd(
  'npm start',
  waitForURL('http://localhost:3000/health', 200)
)
```

***

### setUser()

> **setUser**(`user`): `TemplateBuilder`

Set the user for subsequent commands.

#### Parameters

##### user

`string`

Username

#### Returns

`TemplateBuilder`

#### Example

```ts
template.setUser('root')
```

***

### setWorkdir()

> **setWorkdir**(`workdir`): `TemplateBuilder`

Set the working directory.

#### Parameters

##### workdir

`PathLike`

Working directory path

#### Returns

`TemplateBuilder`

#### Example

```ts
template.setWorkdir('/app')
```

***

### skipCache()

> **skipCache**(): `this`

Skip cache for all subsequent build instructions from this point.

#### Returns

`this`

#### Example

```ts
template.skipCache().runCmd('apt-get update')
```
