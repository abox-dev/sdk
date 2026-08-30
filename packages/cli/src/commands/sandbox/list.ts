import * as commander from 'commander'
import { components, Sandbox, SandboxInfo } from '@abox-dev/sdk'

import { ensureAPIKey } from 'src/api'
import { renderTable } from 'src/utils/table'
import { parseMetadata } from './utils'

const DEFAULT_LIMIT = 1000
const PAGE_LIMIT = 100

export const listCommand = new commander.Command('list')
  .description('list all sandboxes, by default it list only running ones')
  .alias('ls')
  .option(
    '-s, --state <state>',
    'filter by state, eg. running, paused. Defaults to running',
    (value) => value.split(',')
  )
  .option('-m, --metadata <metadata>', 'filter by metadata, eg. key1=value1')
  .option(
    '-l, --limit <limit>',
    `limit the number of sandboxes returned (default: ${DEFAULT_LIMIT}, 0 for no limit)`,
    (value) => parseInt(value)
  )
  .option('-f, --format <format>', 'output format, eg. json, pretty')
  .action(async (options) => {
    try {
      const state = options.state || ['running']
      const format = options.format || 'pretty'
      const limit =
        options.limit === 0 ? undefined : (options.limit ?? DEFAULT_LIMIT)
      const { sandboxes, hasMore } = await listSandboxes({
        limit,
        state,
        metadataRaw: options.metadata,
      })

      if (format === 'pretty') {
        renderSandboxTable(sandboxes)
        if (hasMore) {
          console.log(
            `Showing first ${limit} sandboxes. Use --limit to change.`
          )
        }
      } else if (format === 'json') {
        console.log(JSON.stringify(sortSandboxes(sandboxes), null, 2))
      } else {
        console.error(`Unsupported output format: ${format}`)
        process.exit(1)
      }
    } catch (err: any) {
      console.error(err)
      process.exit(1)
    }
  })

export function sortSandboxes(sandboxes: SandboxInfo[]) {
  return sandboxes
    .slice()
    .sort(
      (a, b) =>
        new Date(b.startedAt).getTime() - new Date(a.startedAt).getTime() ||
        a.sandboxId.localeCompare(b.sandboxId)
    )
}

export function buildTableRows(sandboxes: SandboxInfo[]) {
  return sortSandboxes(sandboxes).map((sandbox) => ({
    ...sandbox,
    startedAt: new Date(sandbox.startedAt).toLocaleString(),
    endAt: new Date(sandbox.endAt).toLocaleString(),
    state: sandbox.state.charAt(0).toUpperCase() + sandbox.state.slice(1), // capitalize
    metadata: JSON.stringify(sandbox.metadata),
  }))
}

function renderSandboxTable(sandboxes: SandboxInfo[]) {
  if (!sandboxes?.length) {
    console.log('No sandboxes found')
    return
  }

  renderTable(buildTableRows(sandboxes), [
    { header: 'Sandbox ID', value: (row) => row.sandboxId },
    { header: 'Template ID', value: (row) => row.templateId },
    { header: 'Alias', value: (row) => row.name },
    { header: 'Started at', value: (row) => row.startedAt },
    { header: 'End at', value: (row) => row.endAt },
    { header: 'State', value: (row) => row.state },
    { header: 'vCPUs', value: (row) => String(row.cpuCount) },
    { header: 'RAM MiB', value: (row) => String(row.memoryMB) },
    { header: 'Envd version', value: (row) => row.envdVersion },
    { header: 'Metadata', value: (row) => row.metadata },
  ])
}

type ListSandboxesOptions = {
  limit?: number
  state?: components['schemas']['SandboxState'][]
  metadataRaw?: string
}

type ListSandboxesResult = {
  sandboxes: SandboxInfo[]
  hasMore: boolean
}

export async function listSandboxes({
  limit,
  state,
  metadataRaw,
}: ListSandboxesOptions = {}): Promise<ListSandboxesResult> {
  const apiKey = ensureAPIKey()
  const metadata = parseMetadata(metadataRaw)

  let pageLimit = limit
  if (!limit || limit > PAGE_LIMIT) {
    pageLimit = PAGE_LIMIT
  }

  const sandboxes: SandboxInfo[] = []
  const iterator = Sandbox.list({
    apiKey: apiKey,
    limit: pageLimit,
    query: { state, metadata },
  })

  while (iterator.hasNext && (!limit || sandboxes.length < limit)) {
    const batch = await iterator.nextItems()
    sandboxes.push(...batch)
  }

  return {
    sandboxes: limit ? sandboxes.slice(0, limit) : sandboxes,
    // We can't change the page size during iteration, so we may have to check if we have more sandboxes than the limit
    hasMore: iterator.hasNext || (limit ? sandboxes.length > limit : false),
  }
}
