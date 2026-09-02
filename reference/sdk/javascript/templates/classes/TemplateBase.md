[agentbox-sdk-monorepo](../README.md) / TemplateBase

# Class: TemplateBase

Builder for AgentBox sandbox templates, and the entrypoint for the template API.

Exposed as [Template](../variables/Template.md), which can be called as a factory.

## Extends

- `ClientFactory`

## Implements

- `TemplateFromImage`
- [`TemplateBuilder`](../interfaces/TemplateBuilder.md)
- `TemplateFinal`

## Constructors

### Constructor

> **new TemplateBase**(`options?`): `TemplateBase`

#### Parameters

##### options?

`TemplateOptions`

#### Returns

`TemplateBase`

#### Overrides

`ClientFactory.constructor`

## Methods

### aptInstall()

> **aptInstall**(`packages`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.aptInstall('vim')
template.aptInstall(['git', 'curl', 'wget'])
template.aptInstall(['vim'], { noInstallRecommends: true })
template.aptInstall(['vim'], { fixMissing: true })
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`aptInstall`](../interfaces/TemplateBuilder.md#aptinstall)

***

### betaDevContainerPrebuild()

> **betaDevContainerPrebuild**(`devcontainerDirectory`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Prebuild a devcontainer from the specified directory.

#### Parameters

##### devcontainerDirectory

`string`

Path to the devcontainer directory

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template
 .gitClone('https://myrepo.com/project.git', '/my-devcontainer')
 .betaDevContainerPrebuild('/my-devcontainer')
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`betaDevContainerPrebuild`](../interfaces/TemplateBuilder.md#betadevcontainerprebuild)

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

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`betaSetDevContainerStart`](../interfaces/TemplateBuilder.md#betasetdevcontainerstart)

***

### bunInstall()

> **bunInstall**(`packages?`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.bunInstall('express')
template.bunInstall(['lodash', 'axios'])
template.bunInstall('tsx', { g: true })
template.bunInstall('typescript', { dev: true })
template.bunInstall()  // Installs from package.json
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`bunInstall`](../interfaces/TemplateBuilder.md#buninstall)

***

### copy()

> **copy**(`src`, `dest`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.copy('requirements.txt', '/home/user/')
template.copy(['app.ts', 'config.ts'], '/app/', { mode: 0o755 })
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`copy`](../interfaces/TemplateBuilder.md#copy)

***

### copyItems()

> **copyItems**(`items`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Copy multiple items with individual options.

#### Parameters

##### items

[`CopyItem`](../type-aliases/CopyItem.md)[]

Array of copy items

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.copyItems([
  { src: 'app.ts', dest: '/app/' },
  { src: 'config.ts', dest: '/app/', mode: 0o644 }
])
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`copyItems`](../interfaces/TemplateBuilder.md#copyitems)

***

### fromAlpineImage()

> **fromAlpineImage**(`variant?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from an Alpine-based Docker image.

#### Parameters

##### variant?

`string` = `'3.24'`

Alpine variant (default: '3.24')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromAlpineImage('3.24')
```

#### Implementation of

`TemplateFromImage.fromAlpineImage`

***

### fromArchImage()

> **fromArchImage**(`variant?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from an Arch Linux-based Docker image.

Defaults to `latest`: Arch is a rolling release and template provisioning
runs `pacman -Syu`, so pinning a tag would not change the built result.

#### Parameters

##### variant?

`string` = `'latest'`

Arch Linux variant (default: 'latest')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromArchImage('base-devel')
```

#### Implementation of

`TemplateFromImage.fromArchImage`

***

### fromAWSRegistry()

> **fromAWSRegistry**(`image`, `credentials`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a Docker image in AWS ECR.

#### Parameters

##### image

`string`

Full ECR image path

##### credentials

AWS credentials

###### accessKeyId

`string`

###### region

`string`

###### secretAccessKey

`string`

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromAWSRegistry(
  '123456789.dkr.ecr.us-west-2.amazonaws.com/myimage:latest',
  {
    accessKeyId: 'AKIA...',
    secretAccessKey: '...',
    region: 'us-west-2'
  }
)
```

#### Implementation of

`TemplateFromImage.fromAWSRegistry`

***

### fromBaseImage()

> **fromBaseImage**(): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from AgentBox's pinned default base image.

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromBaseImage()
```

#### Implementation of

`TemplateFromImage.fromBaseImage`

***

### fromBunImage()

> **fromBunImage**(`variant?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a Bun-based Docker image.

#### Parameters

##### variant?

`string` = `'latest'`

Bun variant (default: 'latest')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromBunImage('1.3')
```

#### Implementation of

`TemplateFromImage.fromBunImage`

***

### fromDebianImage()

> **fromDebianImage**(`variant?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a Debian-based Docker image.

#### Parameters

##### variant?

`string` = `'stable'`

Debian variant (default: 'stable')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromDebianImage('bookworm')
```

#### Implementation of

`TemplateFromImage.fromDebianImage`

***

### fromDockerfile()

> **fromDockerfile**(`dockerfileContentOrPath`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Parse a Dockerfile and convert it to Template SDK format.

#### Parameters

##### dockerfileContentOrPath

`string`

Dockerfile content or path

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromDockerfile('Dockerfile')
Template().fromDockerfile('FROM python:3\nRUN pip install numpy')
```

#### Implementation of

`TemplateFromImage.fromDockerfile`

***

### fromFedoraImage()

> **fromFedoraImage**(`variant?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a Fedora-based Docker image.

#### Parameters

##### variant?

`string` = `'44'`

Fedora variant (default: '44')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromFedoraImage('44')
```

#### Implementation of

`TemplateFromImage.fromFedoraImage`

***

### fromGCPRegistry()

> **fromGCPRegistry**(`image`, `credentials`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a Docker image in Google Container Registry.

#### Parameters

##### image

`string`

Full GCR/GAR image path

##### credentials

GCP service account credentials

###### serviceAccountJSON

`string` \| `object`

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromGCPRegistry(
  'gcr.io/myproject/myimage:latest',
  { serviceAccountJSON: 'path/to/service-account.json' }
)
```

#### Implementation of

`TemplateFromImage.fromGCPRegistry`

***

### fromImage()

> **fromImage**(`baseImage`, `credentials?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a custom Docker image.

#### Parameters

##### baseImage

`string`

Docker image name

##### credentials?

Optional credentials for private registries

###### password

`string`

###### username

`string`

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromImage('python:3')

// With credentials (optional)
Template().fromImage('myregistry.com/myimage:latest', {
  username: 'user',
  password: 'pass'
})
```

#### Implementation of

`TemplateFromImage.fromImage`

***

### fromNodeImage()

> **fromNodeImage**(`variant?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a Node.js-based Docker image.

#### Parameters

##### variant?

`string` = `'lts'`

Node.js variant (default: 'lts')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromNodeImage('24')
```

#### Implementation of

`TemplateFromImage.fromNodeImage`

***

### fromPythonImage()

> **fromPythonImage**(`version?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from a Python-based Docker image.

#### Parameters

##### version?

`string` = `'3'`

Python version (default: '3')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromPythonImage('3')
```

#### Implementation of

`TemplateFromImage.fromPythonImage`

***

### fromTemplate()

> **fromTemplate**(`template`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from an existing AgentBox template.

#### Parameters

##### template

`string`

AgentBox template ID or alias

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromTemplate('my-base-template')
```

#### Implementation of

`TemplateFromImage.fromTemplate`

***

### fromUbuntuImage()

> **fromUbuntuImage**(`variant?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Start from an Ubuntu-based Docker image.

#### Parameters

##### variant?

`string` = `'latest'`

Ubuntu variant (default: 'latest')

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
Template().fromUbuntuImage('24.04')
```

#### Implementation of

`TemplateFromImage.fromUbuntuImage`

***

### gitClone()

> **gitClone**(`url`, `path?`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`gitClone`](../interfaces/TemplateBuilder.md#gitclone)

***

### makeDir()

> **makeDir**(`path`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.makeDir('/app/data', { mode: 0o755 })
template.makeDir(['/app/logs', '/app/cache'])
template.makeDir('/app/data', { mode: 0o755, user: 'root' })
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`makeDir`](../interfaces/TemplateBuilder.md#makedir)

***

### makeSymlink()

> **makeSymlink**(`src`, `dest`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.makeSymlink('/usr/bin/python3', '/usr/bin/python')
template.makeSymlink('/usr/bin/python3', '/usr/bin/python', { user: 'root' })
template.makeSymlink('/usr/bin/python3', '/usr/bin/python', { force: true })
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`makeSymlink`](../interfaces/TemplateBuilder.md#makesymlink)

***

### npmInstall()

> **npmInstall**(`packages?`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.npmInstall('express')
template.npmInstall(['lodash', 'axios'])
template.npmInstall('tsx', { g: true })
template.npmInstall('typescript', { dev: true })
template.npmInstall()  // Installs from package.json
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`npmInstall`](../interfaces/TemplateBuilder.md#npminstall)

***

### pipInstall()

> **pipInstall**(`packages?`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.pipInstall('numpy')  // Installs globally (default)
template.pipInstall(['pandas', 'scikit-learn'])
template.pipInstall('numpy', { g: false })  // Install for user only
template.pipInstall()  // Installs from current directory
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`pipInstall`](../interfaces/TemplateBuilder.md#pipinstall)

***

### remove()

> **remove**(`path`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.remove('/tmp/cache', { recursive: true, force: true })
template.remove('/tmp/cache', { recursive: true, force: true, user: 'root' })
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`remove`](../interfaces/TemplateBuilder.md#remove)

***

### rename()

> **rename**(`src`, `dest`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.rename('/tmp/old.txt', '/tmp/new.txt')
template.rename('/tmp/old.txt', '/tmp/new.txt', { user: 'root' })
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`rename`](../interfaces/TemplateBuilder.md#rename)

***

### runCmd()

#### Call Signature

> **runCmd**(`command`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

##### Example

```ts
template.runCmd('apt-get update')
template.runCmd(['pip install numpy', 'pip install pandas'])
template.runCmd('apt-get install vim', { user: 'root' })
```

##### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`runCmd`](../interfaces/TemplateBuilder.md#runcmd)

#### Call Signature

> **runCmd**(`commands`, `options?`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

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

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

##### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`runCmd`](../interfaces/TemplateBuilder.md#runcmd)

***

### setEnvs()

> **setEnvs**(`envs`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Set environment variables.
Note: Environment variables defined here are available only during template build.

#### Parameters

##### envs

`Record`\<`string`, `string`\>

Environment variables

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.setEnvs({ NODE_ENV: 'production', PORT: '8080' })
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`setEnvs`](../interfaces/TemplateBuilder.md#setenvs)

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

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`setReadyCmd`](../interfaces/TemplateBuilder.md#setreadycmd)

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

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`setStartCmd`](../interfaces/TemplateBuilder.md#setstartcmd)

***

### setUser()

> **setUser**(`user`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Set the user for subsequent commands.

#### Parameters

##### user

`string`

Username

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.setUser('root')
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`setUser`](../interfaces/TemplateBuilder.md#setuser)

***

### setWorkdir()

> **setWorkdir**(`workdir`): [`TemplateBuilder`](../interfaces/TemplateBuilder.md)

Set the working directory.

#### Parameters

##### workdir

`PathLike`

Working directory path

#### Returns

[`TemplateBuilder`](../interfaces/TemplateBuilder.md)

#### Example

```ts
template.setWorkdir('/app')
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`setWorkdir`](../interfaces/TemplateBuilder.md#setworkdir)

***

### skipCache()

> **skipCache**(): `this`

Skip cache for all subsequent build instructions from this point.

#### Returns

`this`

#### Example

```ts
Template().skipCache().fromPythonImage('3')
```

#### Implementation of

[`TemplateBuilder`](../interfaces/TemplateBuilder.md).[`skipCache`](../interfaces/TemplateBuilder.md#skipcache)

***

### assignTags()

> `static` **assignTags**(`targetName`, `tags`, `options?`): `Promise`\<[`TemplateTagInfo`](../type-aliases/TemplateTagInfo.md)\>

Assign tag(s) to an existing template build.

#### Parameters

##### targetName

`string`

Template name in 'name:tag' format (the source build to tag from)

##### tags

`string` \| `string`[]

Tag or tags to assign

##### options?

`ConnectionOpts`

Authentication options

#### Returns

`Promise`\<[`TemplateTagInfo`](../type-aliases/TemplateTagInfo.md)\>

Tag info with buildId and assigned tags

#### Example

```ts
// Assign a single tag
await Template.assignTags('my-template:v1.0', 'production')

// Assign multiple tags
await Template.assignTags('my-template:v1.0', ['production', 'stable'])
```

***

### build()

> `static` **build**(`template`, `name`, `options?`): `Promise`\<[`BuildInfo`](../type-aliases/BuildInfo.md)\>

Build and deploy a template to AgentBox infrastructure.

#### Parameters

##### template

[`TemplateClass`](../type-aliases/TemplateClass.md)

The template to build

##### name

`string`

Template name in 'name' or 'name:tag' format

##### options?

[`BuildOptions`](../type-aliases/BuildOptions.md)

Optional build configuration options

#### Returns

`Promise`\<[`BuildInfo`](../type-aliases/BuildInfo.md)\>

#### Example

```ts
const template = Template().fromPythonImage('3')

// Build with single tag in name
await Template.build(template, 'my-python-env:v1.0')

// Build with multiple tags
await Template.build(template, 'my-python-env', { tags: ['v1.0', 'stable'] })
```

***

### buildInBackground()

> `static` **buildInBackground**(`template`, `name`, `options?`): `Promise`\<[`BuildInfo`](../type-aliases/BuildInfo.md)\>

Build and deploy a template to AgentBox infrastructure without waiting for completion.

#### Parameters

##### template

[`TemplateClass`](../type-aliases/TemplateClass.md)

The template to build

##### name

`string`

Template name in 'name' or 'name:tag' format

##### options?

[`BuildOptions`](../type-aliases/BuildOptions.md)

Optional build configuration options

#### Returns

`Promise`\<[`BuildInfo`](../type-aliases/BuildInfo.md)\>

#### Example

```ts
const template = Template().fromPythonImage('3')

// Build with single tag in name
const data = await Template.buildInBackground(template, 'my-python-env:v1.0')

// Build with multiple tags
const data = await Template.buildInBackground(template, 'my-python-env', { tags: ['v1.0', 'stable'] })
```

***

### exists()

> `static` **exists**(`name`, `options?`): `Promise`\<`boolean`\>

Check if a template with the given name exists.

#### Parameters

##### name

`string`

Template name to check

##### options?

`ConnectionOpts`

Authentication options

#### Returns

`Promise`\<`boolean`\>

True if the name exists, false otherwise

#### Example

```ts
const exists = await Template.exists('my-python-env')
if (exists) {
  console.log('Template exists!')
}
```

***

### getBuildStatus()

> `static` **getBuildStatus**(`data`, `options?`): `Promise`\<[`TemplateBuildStatusResponse`](../type-aliases/TemplateBuildStatusResponse.md)\>

Get the status of a build.

#### Parameters

##### data

`Pick`\<[`BuildInfo`](../type-aliases/BuildInfo.md), `"templateId"` \| `"buildId"`\>

Build identifiers

##### options?

[`GetBuildStatusOptions`](../type-aliases/GetBuildStatusOptions.md)

Authentication options

#### Returns

`Promise`\<[`TemplateBuildStatusResponse`](../type-aliases/TemplateBuildStatusResponse.md)\>

#### Example

```ts
const status = await Template.getBuildStatus(data, { logsOffset: 0 })
```

***

### getTags()

> `static` **getTags**(`templateId`, `options?`): `Promise`\<[`TemplateTag`](../type-aliases/TemplateTag.md)[]\>

Get all tags for a template.

#### Parameters

##### templateId

`string`

Template ID or name

##### options?

`ConnectionOpts`

Authentication options

#### Returns

`Promise`\<[`TemplateTag`](../type-aliases/TemplateTag.md)[]\>

Array of tag details including tag name, buildId, and creation date

#### Example

```ts
const tags = await Template.getTags('my-template')
for (const tag of tags) {
  console.log(`Tag: ${tag.tag}, Build: ${tag.buildId}, Created: ${tag.createdAt}`)
}
```

***

### removeTags()

> `static` **removeTags**(`name`, `tags`, `options?`): `Promise`\<`void`\>

Remove tag(s) from a template.

#### Parameters

##### name

`string`

Template name

##### tags

`string` \| `string`[]

Tag or tags to remove

##### options?

`ConnectionOpts`

Authentication options

#### Returns

`Promise`\<`void`\>

#### Example

```ts
// Remove a single tag
await Template.removeTags('my-template', 'production')

// Remove multiple tags from a template
await Template.removeTags('my-template', ['production', 'staging'])
```

***

### toDockerfile()

> `static` **toDockerfile**(`template`): `string`

Convert a template to Dockerfile format.
Note: Templates based on other AgentBox templates cannot be converted to Dockerfile.

#### Parameters

##### template

[`TemplateClass`](../type-aliases/TemplateClass.md)

The template to convert

#### Returns

`string`

Dockerfile string representation

#### Throws

Error if the template is based on another AgentBox template

***

### toJSON()

> `static` **toJSON**(`template`, `computeHashes?`): `Promise`\<`string`\>

Convert a template to JSON representation.

#### Parameters

##### template

[`TemplateClass`](../type-aliases/TemplateClass.md)

The template to convert

##### computeHashes?

`boolean` = `true`

Whether to compute file hashes for cache invalidation

#### Returns

`Promise`\<`string`\>

JSON string representation of the template
