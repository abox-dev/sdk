#!/usr/bin/env -S node --enable-source-maps

import simpleUpdateNotifier from 'simple-update-notifier'
import * as packageJSON from '../package.json'
import { program } from './commands'

export const pkg = packageJSON

const updateCheck = simpleUpdateNotifier({
  pkg,
  updateCheckInterval: 1000 * 60 * 60 * 8, // 8 hours
}).catch((e) => {
  if (process.env.DEBUG) {
    console.error('Update check failed:', e)
  }
})

const prog = program.version(
  packageJSON.version,
  undefined,
  'display AgentBox CLI version'
)

async function main() {
  await prog.parseAsync()
  await updateCheck
}

main()
