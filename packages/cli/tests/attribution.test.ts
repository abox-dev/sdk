import { createServer, Server } from 'node:http'
import type { AddressInfo } from 'node:net'
import { afterAll, beforeAll, expect, test } from 'vitest'

import * as packageJSON from '../package.json'
import { runCliWithPipedStdin } from './setup'

let server: Server
let apiUrl: string
const userAgents: (string | undefined)[] = []

beforeAll(async () => {
  server = createServer((req, res) => {
    userAgents.push(req.headers['user-agent'])
    res.setHeader('Content-Type', 'application/json')
    res.end('[]')
  })
  await new Promise<void>((resolve) =>
    server.listen(0, '127.0.0.1', () => resolve())
  )
  const { port } = server.address() as AddressInfo
  apiUrl = `http://127.0.0.1:${port}`
})

afterAll(() => {
  server.close()
})

async function captureUserAgent(args: string[]): Promise<string | undefined> {
  const env: NodeJS.ProcessEnv = {
    ...process.env,
    AGENTBOX_API_URL: apiUrl,
    AGENTBOX_API_KEY: `ab_${'0'.repeat(40)}`,
  }
  delete env.AGENTBOX_DEBUG

  userAgents.length = 0
  const result = await runCliWithPipedStdin(args, Buffer.alloc(0), {
    timeoutMs: 30_000,
    env,
  })
  expect(result.error).toBeUndefined()
  return userAgents.find(Boolean)
}

test('CLI requests carry SDK and CLI attribution in the User-Agent', async () => {
  const userAgent = await captureUserAgent(['sandbox', 'list'])

  expect(userAgent).toMatch(/^agentbox-js-sdk\/\d/)
  expect(userAgent).toContain(`agentbox-cli/${packageJSON.version}`)
})
