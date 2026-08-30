import * as commander from 'commander'

import { createCommand } from './create'
import { deleteCommand } from './delete'
import { listCommand } from './list'

export const templateCommand = new commander.Command('template')
  .description('manage sandbox templates')
  .alias('tpl')
  .addCommand(createCommand)
  .addCommand(listCommand)
  .addCommand(deleteCommand)
