[@abox-dev/sdk](../README.md) / Template

# Variable: Template

> `const` **Template**: `CallableTemplate`\<*typeof* [`TemplateBase`](../classes/TemplateBase.md)\>

Builder and API entrypoint for AgentBox sandbox templates.

`Template` is the [TemplateBase](../classes/TemplateBase.md) class, wrapped so it can also be
called as a factory returning a builder. The statics (`Template.build`,
`Template.exists`, …) resolve their connection options off the class they are
called on — so a subclass can bind its own defaults.

## Param

**options**

Optional builder options, e.g. the file context path used to
  resolve relative paths passed to `copy`

## Returns

A template builder

## Example

```ts
import { Template } from '@abox-dev/sdk'

const template = Template()
  .fromPythonImage('3')
  .copy('requirements.txt', '/app/')
  .pipInstall()

await Template.build(template, 'my-python-app:v1.0')
```
