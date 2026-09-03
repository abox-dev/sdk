import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'

const expected = (process.argv[2] ?? '').replace(/^v/, '')
if (!expected) throw new Error('usage: check-release-versions.mjs vX.Y.Z')

const manifests = [
  'packages/js-sdk/package.json',
  'packages/code-interpreter-js/package.json',
  'packages/cli/package.json',
  'packages/python-sdk/package.json',
  'packages/code-interpreter-python/package.json',
]

for (const manifest of manifests) {
  const value = JSON.parse(fs.readFileSync(manifest, 'utf8')).version
  if (value !== expected)
    throw new Error(`${manifest}: expected ${expected}, got ${value}`)
}

const goVersion = fs
  .readFileSync('packages/go-sdk/version.go', 'utf8')
  .match(/^const Version = "([^"]+)"/m)?.[1]
if (goVersion !== expected)
  throw new Error(
    `packages/go-sdk/version.go: expected ${expected}, got ${goVersion}`
  )

for (const manifest of [
  'packages/python-sdk/pyproject.toml',
  'packages/code-interpreter-python/pyproject.toml',
]) {
  const content = fs.readFileSync(manifest, 'utf8')
  const value = content.match(/^version\s*=\s*"([^"]+)"/m)?.[1]
  if (value !== expected)
    throw new Error(`${manifest}: expected ${expected}, got ${value}`)
}

for (const [manifest, packages] of [
  ['packages/python-sdk/uv.lock', ['abox-sdk']],
  [
    'packages/code-interpreter-python/uv.lock',
    ['abox-code-interpreter', 'abox-sdk'],
  ],
]) {
  const blocks = fs.readFileSync(manifest, 'utf8').split('[[package]]')
  for (const packageName of packages) {
    const versions = blocks
      .filter((block) =>
        new RegExp(`^name\\s*=\\s*"${packageName}"$`, 'm').test(block)
      )
      .map((block) => block.match(/^version\s*=\s*"([^"]+)"/m)?.[1])
      .filter(Boolean)
    if (versions.length === 0 || versions.some((value) => value !== expected)) {
      throw new Error(
        `${manifest}: expected every ${packageName} entry to be ${expected}, got ${versions.join(', ') || 'none'}`
      )
    }
  }
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

if (!fs.existsSync('packages/go-sdk/LICENSE')) {
  throw new Error('packages/go-sdk/LICENSE is missing')
}
