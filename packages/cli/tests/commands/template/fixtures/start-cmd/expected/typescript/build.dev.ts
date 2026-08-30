import { Template, defaultBuildLogger } from '@abox-dev/sdk'
import { template } from './template'

async function main() {
  await Template.build(template, 'start-cmd-dev', {
    cpuCount: 2,
    memoryMB: 1024,
    onBuildLogs: defaultBuildLogger(),
  })
}

main().catch(console.error)
