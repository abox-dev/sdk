import * as fs from 'node:fs'
import * as os from 'node:os'
import * as path from 'node:path'

export interface UserConfig {
  apiKey: string
  projectId?: string
}

export const USER_CONFIG_PATH = path.join(
  os.homedir(),
  '.agentbox',
  'config.json'
)

function isUserConfig(value: unknown): value is UserConfig {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return false
  const config = value as Record<string, unknown>
  return (
    typeof config.apiKey === 'string' &&
    config.apiKey.length > 0 &&
    (config.projectId === undefined || typeof config.projectId === 'string')
  )
}

export function getUserConfig(
  configPath = USER_CONFIG_PATH
): UserConfig | null {
  if (!fs.existsSync(configPath)) return null
  try {
    const value: unknown = JSON.parse(fs.readFileSync(configPath, 'utf8'))
    return isUserConfig(value) ? value : null
  } catch {
    return null
  }
}

export function writeUserConfig(configPath: string, config: UserConfig): void {
  const directory = path.dirname(configPath)
  fs.mkdirSync(directory, { recursive: true, mode: 0o700 })
  fs.chmodSync(directory, 0o700)
  fs.writeFileSync(configPath, `${JSON.stringify(config, null, 2)}\n`, {
    mode: 0o600,
  })
  fs.chmodSync(configPath, 0o600)
}
