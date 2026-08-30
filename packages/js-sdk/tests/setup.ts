import { test as base, onTestFailed } from 'vitest'
import {
  BuildInfo,
  LogEntry,
  Sandbox,
  SandboxOpts,
  Template,
  TemplateClass,
} from '../src'
import { template } from './template'

interface SandboxFixture {
  sandbox: Sandbox
  template: string
  sandboxTestId: string
  sandboxOpts: Partial<SandboxOpts>
}

interface BuildTemplateFixture {
  buildTemplate: (
    template: TemplateClass,
    options?: { name?: string; skipCache?: boolean },
    onBuildLogs?: (logEntry: LogEntry) => void
  ) => Promise<BuildInfo>
}

async function buildTemplate(
  template: TemplateClass,
  options?: { name?: string; skipCache?: boolean },
  onBuildLogs?: (logEntry: LogEntry) => void
): Promise<BuildInfo> {
  const buildName = options?.name || `agentbox-test-${generateRandomString()}`
  const buildInfo: { templateId?: string; buildId?: string } = {}

  const captureLogs = (log: LogEntry) => {
    if (log.message.includes('Template created with ID:')) {
      const match = log.message.match(
        /Template created with ID: ([^,]+), Build ID: (.+)/
      )
      if (match) {
        buildInfo.templateId = match[1]
        buildInfo.buildId = match[2]
      }
    }
    onBuildLogs?.(log)
  }

  try {
    return await Template.build(template, buildName, {
      cpuCount: 1,
      memoryMB: 1024,
      skipCache: options?.skipCache,
      onBuildLogs: captureLogs,
    })
  } catch (e) {
    console.error(
      `\n[BUILD FAILED] name=${buildName}, ` +
        `template_id=${buildInfo.templateId}, ` +
        `build_id=${buildInfo.buildId}, error=${e}`
    )
    throw e
  }
}

export const sandboxTest = base.extend<SandboxFixture>({
  template,
  sandboxTestId: [
    // eslint-disable-next-line no-empty-pattern
    async ({}, use) => {
      const id = `test-${generateRandomString()}`
      await use(id)
    },
    { auto: true },
  ],
  sandboxOpts: {},
  sandbox: [
    async ({ sandboxTestId, sandboxOpts }, use) => {
      const sandbox = await Sandbox.create(template, {
        metadata: { sandboxTestId },
        ...sandboxOpts,
      })
      onTestFailed(() => {
        console.error(`\n[TEST FAILED] Sandbox ID: ${sandbox.sandboxId}`)
      })
      try {
        await use(sandbox)
      } finally {
        try {
          await sandbox.kill()
        } catch (err) {
          if (!isDebug) {
            console.warn(
              'Failed to kill sandbox — this is expected if the test runs with local envd.'
            )
          }
        }
      }
    },
    { auto: false },
  ],
})

export const buildTemplateTest = base.extend<BuildTemplateFixture>({
  buildTemplate: [
    // eslint-disable-next-line no-empty-pattern
    async ({}, use) => {
      await use(buildTemplate)
    },
    { auto: true },
  ],
})

export const isDebug = process.env.AGENTBOX_DEBUG !== undefined

/** Placeholder API key with a valid format for tests that don't hit the API. */
export const TEST_API_KEY = `ab_${'0'.repeat(40)}`

function generateRandomString(length: number = 8): string {
  return Math.random()
    .toString(36)
    .substring(2, length + 2)
}

export async function wait(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

/**
 * Returns the API URL for the given path, using AGENTBOX_DOMAIN env var.
 * Supports msw path parameters like :templateID
 */
export function apiUrl(path: string): string {
  const domain = process.env.AGENTBOX_DOMAIN || 'agentbox-runtime.ru'
  return `https://api.${domain}${path}`
}

export { template }
