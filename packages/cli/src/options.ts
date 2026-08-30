import * as commander from 'commander'

import { asBold, asLocal } from './utils/format'

/**
 * Parse a CLI option as a positive integer, rejecting non-numeric values so
 * they don't silently become NaN.
 */
export function parsePositiveInt(label: string): (value: string) => number {
  return (value) => {
    const parsed = Number(value)
    if (!Number.isInteger(parsed) || parsed < 1) {
      throw new commander.InvalidArgumentError(
        `${label} must be a positive integer. You provided ${asLocal(value)}.`
      )
    }
    return parsed
  }
}

export const pathOption = new commander.Option(
  '-p, --path <path>',
  `change root directory where command is executed to ${asBold(
    '<path>'
  )} directory`
)

export const selectMultipleOption = new commander.Option(
  '-s, --select',
  'select sandbox template from interactive list'
)

export const projectOption = new commander.Option(
  '-t, --project <project-id>',
  'specify the project ID for the operation. See https://docs.agentbox.ru/en/cli/configuration/.'
)
