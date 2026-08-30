import * as chalk from 'chalk'
import * as agentbox from '@abox-dev/sdk'
import * as highlight from 'cli-highlight'
import * as boxen from 'boxen'

import { cwdRelative } from './filesystem'

export const primaryColor = '#FFB766'

export type SandboxTemplateRef = Pick<
  agentbox.components['schemas']['Template'],
  'templateID'
> & {
  names?: agentbox.components['schemas']['Template']['names']
}

export function asFormattedSandboxTemplate(template: SandboxTemplateRef) {
  const names = listNames(template.names)

  const name = names ? asBold(names) : ''
  const id = `${template.templateID} `

  return `${id}${name}`.trim()
}

export function asRed(text: string) {
  return chalk.default.redBright(text)
}

export function asFormattedError(text: string | undefined, err?: any) {
  return chalk.default.redBright(
    `${text ? `${text} \n` : ''}${err ? err.stack : ''}\n`
  )
}

export function asDim(content?: string) {
  return chalk.default.dim(content)
}

export function asBold(content: string) {
  return chalk.default.bold(content)
}

export function asPrimary(content: string) {
  return chalk.default.hex(primaryColor)(content)
}

export function asTimestamp(content: string) {
  return chalk.default.blue(content)
}

export function asLocal(pathInLocal?: string) {
  return chalk.default.blue(pathInLocal)
}

export function asLocalRelative(absolutePathInLocal?: string) {
  if (!absolutePathInLocal) return ''
  return asLocal('./' + cwdRelative(absolutePathInLocal))
}

export function asBuildLogs(content: string) {
  return chalk.default.blueBright(content)
}

export function withUnderline(content: string) {
  return chalk.default.underline(content)
}

export function listNames(names: string[] | undefined) {
  if (!names) return undefined
  return names.join(', ')
}

export function asTypescript(code: string) {
  return highlight.default(code, {
    language: 'typescript',
    ignoreIllegals: true,
  })
}

export function asPython(code: string) {
  return highlight.default(code, { language: 'python', ignoreIllegals: true })
}

export const borderStyle = {
  topLeft: '',
  topRight: '',
  bottomLeft: '',
  bottomRight: '',
  top: '',
  bottom: '',
  left: '',
  right: '',
} as const

const horizontalPadding = 2
const verticalPadding = 1

export function withDelimiter(
  content: string,
  title: string,
  isLast?: boolean
) {
  return boxen.default(content, {
    borderStyle: {
      ...borderStyle,
      top: '─',
      bottom: isLast ? '─' : '',
    },
    titleAlignment: 'center',
    float: 'left',
    title: title ? asBold(title) : undefined,
    margin: {
      top: 0,
      bottom: 0,
      left: 1,
      right: 0,
    },
    fullscreen: (w) => [w, 0],
    padding: {
      bottom: isLast ? verticalPadding : 0,
      left: horizontalPadding,
      right: horizontalPadding,
      top: verticalPadding,
    },
  })
}
