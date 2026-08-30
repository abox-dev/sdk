import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'

const expected = (process.argv[2] ?? '').replace(/^v/, '')
if (!expected) throw new Error('usage: check-release-versions.mjs vX.Y.Z')

const manifests = [
  'packages/js-sdk/package.json',
  'packages/code-interpreter-js/package.json',
  'packages/cli/package.json',
]

for (const manifest of manifests) {
  const value = JSON.parse(fs.readFileSync(manifest, 'utf8')).version
  if (value !== expected)
    throw new Error(`${manifest}: expected ${expected}, got ${value}`)
}

for (const manifest of [
  'packages/python-sdk/pyproject.toml',
  'packages/code-interpreter-python/pyproject.toml',
]) {
  const content = fs.readFileSync(manifest, 'utf8')
  const value = content.match(/^version\s*=\s*"([^"]+)"/m)?.[1]
  if (value !== expected)
    throw new Error(`${manifest}: expected ${expected}, got ${value}`)
}

for (const directory of [
  'packages/js-sdk',
  'packages/code-interpreter-js',
  'packages/cli',
  'packages/python-sdk',
  'packages/code-interpreter-python',
]) {
  if (!fs.existsSync(path.join(directory, 'NOTICE'))) {
    throw new Error(`${directory}/NOTICE is missing`)
  }
}
