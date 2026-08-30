import { defineConfig } from 'vitest/config'
import { playwright } from '@vitest/browser-playwright'
import { config } from 'dotenv'

const env = config()
const unitTests = [
  'tests/api/{apiKey,handleApiError,http2,inflight}.test.ts',
  'tests/{client,cmdHelper,foreignPlatformObjects,is,logs,paginator,runtime,stripAnsi,undici,utils}.test.ts',
  'tests/envd/**/*.test.ts',
  'tests/sandbox/{configPropagation,iam,lifecycleRequest,networkTransform,rpcHeaders,urls}.test.ts',
  'tests/sandbox/commands/commandHandle.test.ts',
  'tests/sandbox/files/{entryInfo,watchHandle}.test.ts',
]

export default defineConfig({
  test: {
    projects: [
      {
        test: {
          name: 'unit',
          include: unitTests,
          // Isolation is required: several suites patch global fetch via msw
          // and rely on module mocks (vi.doMock / vi.resetModules). Under
          // vitest 4 a shared (non-isolated) context leaks this state across
          // files — e.g. aborted-request rejections and the cached undici
          // apiFetch singleton — causing cross-file failures.
          isolate: true,
          globals: false,
          testTimeout: 30_000,
          environment: 'node',
          bail: 0,
          setupFiles: ['tests/globalFetchFallback.setup.ts'],
          server: {},
          deps: {
            interopDefault: true,
          },
          env: {
            ...(process.env as Record<string, string>),
            ...env.parsed,
          },
        },
      },
      {
        test: {
          name: 'integration',
          include: [
            'tests/api/{info,kill,list,snapshot}.test.ts',
            'tests/sandbox/**/*.test.ts',
          ],
          exclude: unitTests,
          isolate: true,
          globals: false,
          testTimeout: 180_000,
          hookTimeout: 180_000,
          environment: 'node',
          setupFiles: ['tests/globalFetchFallback.setup.ts'],
          env: {
            ...(process.env as Record<string, string>),
            ...env.parsed,
          },
        },
      },
      {
        test: {
          name: 'browser',
          include: ['tests/runtimes/browser/**/*.{test,spec}.tsx'],
          browser: {
            enabled: true,
            headless: true,
            instances: [{ browser: 'chromium' }],
            provider: playwright(),
            // https://playwright.dev
          },
          provide: {
            AGENTBOX_API_KEY: process.env.AGENTBOX_API_KEY || env.parsed?.AGENTBOX_API_KEY,
            AGENTBOX_DOMAIN: process.env.AGENTBOX_DOMAIN || env.parsed?.AGENTBOX_DOMAIN,
          },
        },
      },
      {
        test: {
          name: 'template',
          include: ['tests/template/**/*.test.ts'],
          globals: false,
          testTimeout: 180_000,
          environment: 'node',
          setupFiles: ['tests/globalFetchFallback.setup.ts'],
        },
      },
      {
        test: {
          name: 'connectionConfig',
          include: ['tests/connectionConfig.test.ts'],
          globals: false,
          isolate: true,
          testTimeout: 10_000,
          environment: 'node',
        },
      },
    ],
  },
})
