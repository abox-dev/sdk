import * as boxen from 'boxen'
import * as agentbox from '@abox-dev/sdk'

import * as packageJSON from '../package.json'
import { getUserConfig } from './user'
import { asBold, asPrimary } from './utils/format'

agentbox.ConnectionConfig.setIntegration(`agentbox-cli/${packageJSON.version}`)

export function resolveApiKey(): string | undefined {
  return process.env.AGENTBOX_API_KEY || getUserConfig()?.apiKey
}

export function resolveProjectId(cliProjectId?: string): string | undefined {
  return (
    cliProjectId ||
    process.env.AGENTBOX_PROJECT_ID ||
    getUserConfig()?.projectId
  )
}

const authErrorBox = () =>
  boxen.default(
    `AgentBox API key is not configured. Run ${asBold(
      'agentbox configure'
    )} or set ${asBold('AGENTBOX_API_KEY')}.\n\n${asPrimary(
      'https://docs.agentbox.ru/en/quickstart/api-key/'
    )}`,
    {
      width: 70,
      float: 'center',
      padding: 0.5,
      margin: 1,
      borderStyle: 'round',
      borderColor: 'redBright',
    }
  )

export function ensureAPIKey(): string {
  const apiKey = resolveApiKey()
  if (apiKey) return apiKey
  console.error(authErrorBox())
  return process.exit(1)
}

export const connectionConfig = new agentbox.ConnectionConfig({
  apiKey: resolveApiKey(),
})

export const client = new agentbox.ApiClient(connectionConfig, {
  requireApiKey: false,
})
