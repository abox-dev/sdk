import assert from 'node:assert/strict'

import { Sandbox } from '@abox-dev/code-interpreter'

const sandbox = await Sandbox.create('code-interpreter', { timeoutMs: 180_000 })

try {
  const stdout = []
  const python = await sandbox.runCode("print('ci-stream')\n{'answer': 42}", {
    language: 'python',
    onStdout: (message) => stdout.push(message.line),
  })
  assert(stdout.join('').includes('ci-stream'))
  assert.equal(python.error, undefined)

  const javascript = await sandbox.runCode('console.log("js-ok"); 6 * 7', {
    language: 'javascript',
  })
  assert(javascript.logs.stdout.join('').includes('js-ok'))
  assert.equal(javascript.error, undefined)

  const typescript = await sandbox.runCode(
    'const value: number = 42; console.log(value)',
    { language: 'typescript' }
  )
  assert(typescript.logs.stdout.join('').includes('42'))

  const failed = await sandbox.runCode(
    'raise RuntimeError("expected-ci-error")'
  )
  assert.equal(failed.error?.name, 'RuntimeError')

  const context = await sandbox.createCodeContext({ language: 'python' })
  await sandbox.runCode('context_value = 41', { context })
  const contextual = await sandbox.runCode('context_value + 1', { context })
  assert.equal(contextual.error, undefined)
  assert(
    (await sandbox.listCodeContexts()).some((item) => item.id === context.id)
  )
  await sandbox.removeCodeContext(context)

  const data = await sandbox.runCode(
    "import pandas as pd\ndisplay(pd.DataFrame({'x': [1, 2], 'y': [3, 4]}))"
  )
  assert(data.results.some((result) => result.data || result.html))

  const chart = await sandbox.runCode(
    'import matplotlib.pyplot as plt\nplt.plot([1, 2], [3, 4])\nplt.show()'
  )
  assert(chart.results.some((result) => result.chart || result.png))

  await sandbox.files.write('/tmp/ci-js.txt', 'code-interpreter-file')
  assert.equal(
    await sandbox.files.read('/tmp/ci-js.txt'),
    'code-interpreter-file'
  )
  console.log('JavaScript Code Interpreter runtime smoke passed')
} finally {
  await sandbox.kill()
}
