import { Command, Option } from 'commander'
import fs from 'node:fs'
import path from 'node:path'

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
}

function renderOption(option: Option): string {
  const fallback =
    option.defaultValue === undefined
      ? ''
      : ` Default: \`${String(option.defaultValue)}\`.`
  return `- \`${escapeHtml(option.flags)}\`: ${escapeHtml(option.description || 'No description.')}${fallback}`
}

function renderCommand(command: Command, parents: string[] = []): string {
  const fullName = [...parents, command.name()].join(' ')
  const usage = command.usage() || '[options]'
  const lines = [
    `## agentbox ${fullName}`,
    '',
    escapeHtml(command.description() || 'No description.'),
    '',
    '### Usage',
    '',
    '```bash',
    `agentbox ${fullName} ${usage}`.trimEnd(),
    '```',
    '',
  ]

  if (command.options.length) {
    lines.push('### Options', '', ...command.options.map(renderOption), '')
  }

  for (const subcommand of [...command.commands].sort((a, b) =>
    a.name().localeCompare(b.name())
  )) {
    lines.push(renderCommand(subcommand, [...parents, command.name()]))
  }
  return lines.join('\n')
}

export function commands2md(
  commands: readonly Command[],
  outputDir: string
): void {
  fs.rmSync(outputDir, { recursive: true, force: true })
  fs.mkdirSync(outputDir, { recursive: true })

  const groups = new Map<string, Command[]>([
    ['auth', commands.filter((command) => command.name() === 'configure')],
    ['sandbox', commands.filter((command) => command.name() === 'sandbox')],
    ['template', commands.filter((command) => command.name() === 'template')],
  ])

  for (const [group, groupCommands] of groups) {
    if (!groupCommands.length) {
      throw new Error(`CLI reference group ${group} has no commands`)
    }
    const title = group[0].toUpperCase() + group.slice(1)
    const markdown = [
      `# ${title} commands`,
      '',
      ...groupCommands.map((command) => renderCommand(command)),
    ].join('\n')
    fs.writeFileSync(path.join(outputDir, `${group}.md`), `${markdown}\n`)
  }
}
