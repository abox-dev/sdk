import assert from 'node:assert/strict'

import {
  AgentBox,
  CommandExitError,
  Sandbox,
  TimeoutError,
} from '@abox-dev/sdk'

const marker = `sdk-runtime-js-${Date.now()}`
const sandboxes = new Map()

function track(sandbox) {
  sandboxes.set(sandbox.sandboxId, sandbox)
  return sandbox
}

async function cleanup() {
  await Promise.allSettled(
    [...sandboxes.values()].map((sandbox) => sandbox.kill())
  )
}

async function waitForHttp(sandbox, port, includeToken) {
  const headers = {
    'Agentbox-Sandbox-Id': sandbox.sandboxId,
    'Agentbox-Sandbox-Port': String(port),
  }
  if (includeToken) {
    headers['Agentbox-Traffic-Access-Token'] = sandbox.trafficAccessToken
  }

  let lastStatus = 0
  for (let attempt = 0; attempt < 30; attempt += 1) {
    const response = await fetch(process.env.AGENTBOX_SANDBOX_URL, { headers })
    lastStatus = response.status
    if (
      includeToken
        ? response.ok
        : response.status === 401 || response.status === 403
    ) {
      return response
    }
    await new Promise((resolve) => setTimeout(resolve, 250))
  }
  throw new Error(
    `sandbox HTTP proxy did not become ready, last status ${lastStatus}`
  )
}

try {
  assert.equal(typeof AgentBox, 'function')
  const client = new AgentBox()
  let sandbox = track(
    await client.Sandbox.create('base', {
      metadata: { sdkRuntime: marker },
      network: { allowPublicTraffic: false },
      timeoutMs: 180_000,
    })
  )

  assert.equal((await sandbox.getInfo()).sandboxId, sandbox.sandboxId)
  assert.match(
    sandbox.getHost(8123),
    new RegExp(`^8123-${sandbox.sandboxId}\\.`)
  )

  const paginator = Sandbox.list({
    query: { metadata: { sdkRuntime: marker } },
  })
  const listed = await paginator.nextItems()
  assert(listed.some((item) => item.sandboxId === sandbox.sandboxId))

  let streamedStdout = ''
  let streamedStderr = ''
  const command = await sandbox.commands.run(
    "printf 'agentbox-out'; printf 'agentbox-err' >&2",
    {
      onStdout: (chunk) => {
        streamedStdout += chunk
      },
      onStderr: (chunk) => {
        streamedStderr += chunk
      },
    }
  )
  assert.equal(command.stdout, 'agentbox-out')
  assert.equal(command.stderr, 'agentbox-err')
  assert.equal(streamedStdout, command.stdout)
  assert.equal(streamedStderr, command.stderr)

  await assert.rejects(
    sandbox.commands.run('sleep 2', { timeoutMs: 100 }),
    (error) => error instanceof TimeoutError
  )
  await assert.rejects(
    sandbox.commands.run("printf 'expected-error' >&2; exit 7"),
    (error) => error instanceof CommandExitError && error.exitCode === 7
  )

  await sandbox.files.write('/tmp/agentbox/text.txt', 'hello-agentbox')
  assert.equal(
    await sandbox.files.read('/tmp/agentbox/text.txt'),
    'hello-agentbox'
  )
  const bytes = new Uint8Array([0, 1, 2, 127, 255])
  await sandbox.files.write('/tmp/agentbox/data.bin', bytes.buffer)
  assert.deepEqual(
    [
      ...(await sandbox.files.read('/tmp/agentbox/data.bin', {
        format: 'bytes',
      })),
    ],
    [...bytes]
  )
  assert(
    (await sandbox.files.list('/tmp/agentbox')).some(
      (entry) => entry.name === 'data.bin'
    )
  )
  await sandbox.files.remove('/tmp/agentbox/data.bin')
  assert.equal(await sandbox.files.exists('/tmp/agentbox/data.bin'), false)

  const server = await sandbox.commands.run(
    'python3 -m http.server 8123 --bind 0.0.0.0',
    { background: true }
  )
  assert(sandbox.trafficAccessToken)
  const rejected = await waitForHttp(sandbox, 8123, false)
  assert([401, 403].includes(rejected.status))
  const accepted = await waitForHttp(sandbox, 8123, true)
  assert.equal(accepted.status, 200)
  await server.kill()

  await sandbox.setTimeout(180_000)
  const sandboxId = sandbox.sandboxId
  assert.equal(await sandbox.pause(), true)
  sandbox = track(await Sandbox.connect(sandboxId, { timeoutMs: 180_000 }))
  assert.equal((await sandbox.getInfo()).state, 'running')

  const [fork] = await sandbox.fork({ count: 1, timeoutMs: 180_000 })
  if (fork instanceof Error) throw fork
  track(fork)
  assert.equal(
    await fork.files.read('/tmp/agentbox/text.txt'),
    'hello-agentbox'
  )

  console.log('JavaScript core runtime smoke passed')
} finally {
  await cleanup()
}
