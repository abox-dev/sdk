import * as commander from 'commander'

import { ensureAPIKey } from 'src/api'
import { asBold } from 'src/utils/format'
import * as agentbox from '@abox-dev/sdk'
import { SandboxNotFoundError } from '@abox-dev/sdk'

async function pauseSandbox(sandboxID: string, apiKey: string) {
  try {
    const paused = await agentbox.Sandbox.pause(sandboxID, { apiKey })
    if (paused) {
      console.log(`Sandbox ${asBold(sandboxID)} has been paused`)
    } else {
      console.log(`Sandbox ${asBold(sandboxID)} is already paused`)
    }
  } catch (err: unknown) {
    if (err instanceof SandboxNotFoundError) {
      console.error(`Sandbox ${asBold(sandboxID)} wasn't found`)
      process.exit(1)
    }
    throw err
  }
}

export const pauseCommand = new commander.Command('pause')
  .description('pause sandbox')
  .argument(
    '<sandboxID>',
    `pause the sandbox specified by ${asBold('<sandboxID>')}`
  )
  .alias('ps')
  .action(async (sandboxID: string) => {
    try {
      const apiKey = ensureAPIKey()
      await pauseSandbox(sandboxID, apiKey)
    } catch (err: unknown) {
      console.error(err)
      process.exit(1)
    }
  })
