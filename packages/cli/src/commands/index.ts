import * as commander from 'commander'

import { asPrimary } from 'src/utils/format'
import { templateCommand } from './template'
import { sandboxCommand } from './sandbox'
import { configureCommand } from './auth/configure'

export const program = new commander.Command()
  .enablePositionalOptions()
  .description(
    `Create sandbox templates from Dockerfiles by running ${asPrimary(
      'agentbox template create'
    )} then use our SDKs to create sandboxes from these templates.

Visit ${asPrimary(
      'AgentBox docs (https://docs.agentbox.ru/en/cli/)'
    )} to learn how to create sandbox templates and start sandboxes.
`
  )
  .addCommand(configureCommand)
  .addCommand(templateCommand)
  .addCommand(sandboxCommand)
