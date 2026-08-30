import crypto from 'node:crypto'
import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'
import { spawnSync } from 'node:child_process'

const directory = path.resolve(process.argv[2] ?? 'release/npm')

function npm(args, options = {}) {
  return spawnSync('npm', args, { encoding: 'utf8', ...options })
}

function manifestFromTarball(filename) {
  const result = spawnSync('tar', ['-xOzf', filename, 'package/package.json'], {
    encoding: 'utf8',
  })
  if (result.status !== 0) {
    throw new Error(
      result.stderr || `Cannot read package.json from ${filename}`
    )
  }
  return JSON.parse(result.stdout)
}

for (const tarball of fs
  .readdirSync(directory)
  .filter((name) => name.endsWith('.tgz'))
  .sort()) {
  const filename = path.join(directory, tarball)
  const manifest = manifestFromTarball(filename)
  const spec = `${manifest.name}@${manifest.version}`
  const expected = `sha512-${crypto.createHash('sha512').update(fs.readFileSync(filename)).digest('base64')}`
  const current = npm(['view', spec, 'dist.integrity', '--json'])

  if (current.status === 0 && current.stdout.trim()) {
    const actual = JSON.parse(current.stdout)
    if (actual !== expected)
      throw new Error(`${spec} exists with a different digest`)
    process.stdout.write(`${spec} already exists with the expected digest\n`)
    continue
  }

  const published = npm(['publish', filename, '--access', 'public'], {
    stdio: 'inherit',
  })
  if (published.status !== 0) process.exit(published.status ?? 1)
}
