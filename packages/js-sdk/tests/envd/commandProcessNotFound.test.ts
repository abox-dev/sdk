import { Code, ConnectError } from '@connectrpc/connect'
import { describe, expect, it, vi } from 'vitest'

import { CommandNotFoundError } from '../../src/errors'
import { Commands } from '../../src/sandbox/commands'

function commandsWithMissingProcess(): Commands {
  const missing = () =>
    Promise.reject(new ConnectError('process not found', Code.NotFound))
  const commands = Object.create(Commands.prototype) as Commands
  Object.assign(commands, {
    rpc: { sendInput: vi.fn(missing), closeStdin: vi.fn(missing) },
    envdVersion: '0.6.15',
    connectionConfig: { getSignal: () => undefined },
    checkHealth: async () => true,
  })
  return commands
}

describe('command process disappearance', () => {
  it('identifies a missing PID while sending stdin', async () => {
    await expect(
      commandsWithMissingProcess().sendStdin(1234, 'input')
    ).rejects.toBeInstanceOf(CommandNotFoundError)
  })

  it('identifies a missing PID while closing stdin', async () => {
    await expect(
      commandsWithMissingProcess().closeStdin(1234)
    ).rejects.toBeInstanceOf(CommandNotFoundError)
  })
})
