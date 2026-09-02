import { assert } from 'vitest'

import { isDebug, sandboxTest, wait } from '../setup.js'
import { catchCmdExitErrorInBackground } from '../cmdHelper.js'

const sandboxProtocol = process.env.AGENTBOX_SANDBOX_URL
  ? new URL(process.env.AGENTBOX_SANDBOX_URL).protocol.slice(0, -1)
  : isDebug
    ? 'http'
    : 'https'
const sandboxWebSocketProtocol = sandboxProtocol === 'http' ? 'ws' : 'wss'

const websocketServer = `
import base64
import hashlib
import socket

server = socket.socket()
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(("0.0.0.0", 8765))
server.listen()

while True:
    connection, _ = server.accept()
    request = b""
    while b"\\r\\n\\r\\n" not in request:
        chunk = connection.recv(4096)
        if not chunk:
            break
        request += chunk

    headers = {}
    for line in request.decode("latin1").split("\\r\\n")[1:]:
        if ":" in line:
            name, value = line.split(":", 1)
            headers[name.lower()] = value.strip()

    key = headers.get("sec-websocket-key")
    if not key:
        connection.close()
        continue

    accept = base64.b64encode(
        hashlib.sha1((key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11").encode()).digest()
    ).decode()
    connection.sendall(
        (
            "HTTP/1.1 101 Switching Protocols\\r\\n"
            "Upgrade: websocket\\r\\n"
            "Connection: Upgrade\\r\\n"
            f"Sec-WebSocket-Accept: {accept}\\r\\n\\r\\n"
        ).encode()
    )
    message = b"agentbox-websocket-ok"
    connection.sendall(bytes([0x81, len(message)]) + message)
    connection.close()
`
sandboxTest(
  'ping server in running sandbox',
  async ({ sandbox }) => {
    const cmd = await sandbox.commands.run('python -m http.server 8000', {
      background: true,
    })

    const disable = catchCmdExitErrorInBackground(cmd)

    try {
      await wait(1000)

      const host = sandbox.getHost(8000)

      let res = await fetch(`${sandboxProtocol}://${host}`)

      for (let i = 0; i < 20; i++) {
        if (res.status === 200) {
          break
        }

        res = await fetch(`${sandboxProtocol}://${host}`)
        await wait(500)
      }
      assert.equal(res.status, 200)
      disable()
    } finally {
      try {
        await cmd.kill()
      } catch (e) {
        console.error(e)
      }
    }
  },
  60_000
)

sandboxTest(
  'routes WebSocket traffic through the public sandbox host',
  async ({ sandbox }) => {
    const encodedServer = Buffer.from(websocketServer).toString('base64')
    const cmd = await sandbox.commands.run(
      `python3 -c "import base64; exec(base64.b64decode('${encodedServer}'))"`,
      { background: true }
    )
    const disable = catchCmdExitErrorInBackground(cmd)

    try {
      await wait(500)
      const host = sandbox.getHost(8765)
      const message = await new Promise<string>((resolve, reject) => {
        const socket = new WebSocket(`${sandboxWebSocketProtocol}://${host}`)
        const timeout = setTimeout(
          () => reject(new Error('WebSocket response timed out')),
          10_000
        )
        socket.addEventListener('message', (event) => {
          clearTimeout(timeout)
          resolve(String(event.data))
          socket.close()
        })
        socket.addEventListener('error', () => {
          clearTimeout(timeout)
          reject(new Error('WebSocket connection failed'))
        })
      })

      assert.equal(message, 'agentbox-websocket-ok')
      disable()
    } finally {
      await cmd.kill().catch(() => undefined)
    }
  },
  60_000
)

sandboxTest.skipIf(isDebug)(
  'ping server in non-running sandbox',
  async ({ sandbox }) => {
    const host = sandbox.getHost(3000)
    const url = `https://${host}`

    await sandbox.kill()

    const res = await fetch(url)
    assert.equal(res.status, 502)

    const text = await res.text()
    const json = JSON.parse(text) as {
      message: string
      sandboxId: string
      code: number
    }
    assert.equal(json.message, 'The sandbox was not found')
    assert.isTrue(sandbox.sandboxId.startsWith(json.sandboxId))
    assert.equal(json.code, 502)
  }
)
