import { Command, Option } from 'commander'
import { readFileSync, rmSync } from 'node:fs'
import { tmpdir } from 'node:os'
import path from 'node:path'
import { afterEach, describe, expect, test } from 'vitest'
import { commands2md } from '../../src/utils/commands2md'

const outputDir = path.join(
  tmpdir(),
  `agentbox-cli-reference-${process.pid}-${Date.now()}`
)

afterEach(() => rmSync(outputDir, { recursive: true, force: true }))

describe('commands2md defaults', () => {
  test('omits empty objects and serializes structured defaults', () => {
    const configure = new Command('configure')
    const sandbox = new Command('sandbox')
      .addOption(
        new Option('-e, --env <KEY=VALUE>', 'environment variables').default({})
      )
      .addOption(
        new Option('--labels <JSON>', 'labels').default({ zebra: 2, alpha: 1 })
      )
      .addOption(new Option('--ports <PORT>', 'ports').default([8080, 8081]))
      .addOption(new Option('--format <FORMAT>', 'format').default('pretty'))
    const template = new Command('template')

    commands2md([configure, sandbox, template], outputDir)

    const markdown = readFileSync(path.join(outputDir, 'sandbox.md'), 'utf8')
    expect(markdown).not.toContain('[object Object]')
    expect(markdown).toContain(
      '- `--labels &lt;JSON&gt;`: labels Default: `{"alpha":1,"zebra":2}`.'
    )
    expect(markdown).toContain(
      '- `--ports &lt;PORT&gt;`: ports Default: `[8080,8081]`.'
    )
    expect(markdown).toContain(
      '- `--format &lt;FORMAT&gt;`: format Default: `pretty`.'
    )
    expect(markdown).toContain(
      '- `-e, --env &lt;KEY=VALUE&gt;`: environment variables\n'
    )
  })
})
