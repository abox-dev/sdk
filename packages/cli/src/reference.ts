import path from 'node:path'

import { program } from './commands'
import { commands2md } from './utils/commands2md'

const destination = process.argv[2]
if (!destination) {
  throw new Error('Usage: pnpm reference OUTPUT_DIR')
}

commands2md(program.commands, path.resolve(destination))
