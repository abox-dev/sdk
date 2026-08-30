import * as commander from 'commander'

import { USER_CONFIG_PATH, writeUserConfig } from 'src/user'

export const configureCommand = new commander.Command('configure')
  .description('store an AgentBox API key and optional project ID')
  .argument('[api-key]', 'AgentBox API key; prompted securely when omitted')
  .option('--project-id <project-id>', 'default AgentBox project ID')
  .action(
    async (apiKey: string | undefined, options: { projectId?: string }) => {
      let resolvedApiKey = apiKey
      let projectId = options.projectId

      if (!resolvedApiKey) {
        const { password, input } = await import('@inquirer/prompts')
        resolvedApiKey = await password({ message: 'AgentBox API key' })
        if (!projectId) {
          const value = await input({
            message: 'Project ID (optional)',
          })
          projectId = value || undefined
        }
      }

      if (!resolvedApiKey?.startsWith('ab_')) {
        throw new commander.InvalidArgumentError(
          'AgentBox API keys must start with ab_'
        )
      }

      writeUserConfig(USER_CONFIG_PATH, {
        apiKey: resolvedApiKey,
        ...(projectId ? { projectId } : {}),
      })
      console.log(`Configuration saved to ${USER_CONFIG_PATH}`)
    }
  )
