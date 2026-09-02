import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import process from 'node:process'
import { spawnSync } from 'node:child_process'

const artifactRoot = path.resolve(process.argv[2] ?? 'release')
if (!fs.existsSync(artifactRoot))
  throw new Error(`artifact directory not found: ${artifactRoot}`)

const temporary = fs.mkdtempSync(path.join(os.tmpdir(), 'agentbox-artifacts-'))
const forbidden = [
  /\bE2B[A-Za-z0-9_]*/,
  /\be2b[_-][A-Za-z0-9_-]*/,
  /@e2b\//,
  /\be2b\.(?:dev|app|com)\b/,
  /\b(?:E2B|e2b)_(?:API|DOMAIN|SANDBOX|PROJECT|DEBUG)/,
  /\b(?:class|interface|type|export|import)\s+(?:E2B|Desktop|Volume|Secret)\b/,
  /\b(?:mcp-gateway|addMcpServer|getMcpUrl|McpServer|add_mcp_server|get_mcp_url|get_mcp_token)\b/,
]
const textExtensions = new Set([
  '.cjs',
  '.cts',
  '.js',
  '.json',
  '.map',
  '.md',
  '.mjs',
  '.mts',
  '.py',
  '.toml',
  '.ts',
  '.txt',
  '.yaml',
  '.yml',
])

function run(command, args) {
  const result = spawnSync(command, args, { encoding: 'utf8' })
  if (result.status !== 0)
    throw new Error(`${command} failed: ${result.stderr}`)
}

function walk(directory) {
  return fs.readdirSync(directory, { withFileTypes: true }).flatMap((entry) => {
    const target = path.join(directory, entry.name)
    return entry.isDirectory() ? walk(target) : [target]
  })
}

try {
  for (const archive of walk(artifactRoot)) {
    const relative = path.relative(artifactRoot, archive)
    if (/e2b/i.test(relative))
      throw new Error(`legacy name in artifact path: ${relative}`)
    const destination = path.join(
      temporary,
      relative.replace(/[^a-zA-Z0-9_.-]/g, '_')
    )
    fs.mkdirSync(destination, { recursive: true })
    if (archive.endsWith('.tgz') || archive.endsWith('.tar.gz')) {
      run('tar', ['-xzf', archive, '-C', destination])
    } else if (archive.endsWith('.whl')) {
      run('unzip', ['-q', archive, '-d', destination])
    } else {
      continue
    }
  }

  for (const file of walk(temporary)) {
    const basename = path.basename(file)
    if (basename === 'LICENSE' || basename === 'NOTICE') continue
    if (!textExtensions.has(path.extname(file))) continue
    const content = fs.readFileSync(file, 'utf8')
    for (const pattern of forbidden) {
      if (pattern.test(content)) {
        throw new Error(
          `forbidden public identifier ${pattern} in ${path.relative(temporary, file)}`
        )
      }
    }
  }
} finally {
  fs.rmSync(temporary, { recursive: true, force: true })
}
