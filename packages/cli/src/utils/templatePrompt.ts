import * as agentbox from '@abox-dev/sdk'
import * as chalk from 'chalk'

import { asFormattedSandboxTemplate } from 'src/utils/format'

export async function getPromptTemplates(
  templates: agentbox.components['schemas']['Template'][],
  text: string
) {
  const inquirer = await import('inquirer')
  const templatesAnswers = await inquirer.default.prompt([
    {
      name: 'templates',
      message: chalk.default.underline(text),
      type: 'checkbox',
      pageSize: 50,
      choices: templates.map((e) => ({
        name: asFormattedSandboxTemplate(e),
        value: e,
      })),
    },
  ])

  return templatesAnswers[
    'templates'
  ] as agentbox.components['schemas']['Template'][]
}
