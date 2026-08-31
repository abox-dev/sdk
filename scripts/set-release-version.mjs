import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'
import { fileURLToPath } from 'node:url'

const version = process.argv[2]
if (!version || !/^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$/.test(version)) {
  throw new Error('usage: pnpm release:version X.Y.Z')
}

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
const jsonManifests = [
  'packages/js-sdk/package.json',
  'packages/code-interpreter-js/package.json',
  'packages/cli/package.json',
  'packages/python-sdk/package.json',
  'packages/code-interpreter-python/package.json',
]
const pythonManifests = [
  'packages/python-sdk/pyproject.toml',
  'packages/code-interpreter-python/pyproject.toml',
]

const currentVersions = new Set()
for (const relative of jsonManifests) {
  const filename = path.join(root, relative)
  const manifest = JSON.parse(fs.readFileSync(filename, 'utf8'))
  currentVersions.add(manifest.version)
}
for (const relative of pythonManifests) {
  const content = fs.readFileSync(path.join(root, relative), 'utf8')
  const current = content.match(/^version\s*=\s*"([^"]+)"/m)?.[1]
  if (!current) throw new Error(`Cannot find version in ${relative}`)
  currentVersions.add(current)
}
if (currentVersions.size !== 1) {
  throw new Error(
    `Package versions are already inconsistent: ${[...currentVersions].join(', ')}`
  )
}

for (const relative of jsonManifests) {
  const filename = path.join(root, relative)
  const manifest = JSON.parse(fs.readFileSync(filename, 'utf8'))
  manifest.version = version
  fs.writeFileSync(filename, `${JSON.stringify(manifest, null, 2)}\n`)
}
for (const relative of pythonManifests) {
  const filename = path.join(root, relative)
  const content = fs.readFileSync(filename, 'utf8')
  fs.writeFileSync(
    filename,
    content.replace(/^version\s*=\s*"[^"]+"/m, `version = "${version}"`)
  )
}

process.stdout.write(`Updated all AgentBox SDK packages to ${version}\n`)
