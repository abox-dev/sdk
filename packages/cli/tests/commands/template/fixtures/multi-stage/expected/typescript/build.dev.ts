import { Template, defaultBuildLogger } from '@abox-dev/sdk'
import { template } from './template'

async function main() {
  await Template.build(template, 'multi-stage-dev', {
    onBuildLogs: defaultBuildLogger(),
  })
}

main().catch(console.error)
