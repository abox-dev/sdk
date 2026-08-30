import { Template, defaultBuildLogger } from '@abox-dev/sdk'
import { template } from './template'

async function main() {
  await Template.build(template, 'env-test', {
    onBuildLogs: defaultBuildLogger(),
  })
}

main().catch(console.error)
