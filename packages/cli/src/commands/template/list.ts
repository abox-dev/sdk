import * as commander from 'commander'
import * as agentbox from '@abox-dev/sdk'

import { listNames } from '../../utils/format'
import { renderTable } from 'src/utils/table'
import { sortTemplateNames } from 'src/utils/templateSort'
import { client, ensureAPIKey, resolveProjectId } from 'src/api'
import { projectOption } from '../../options'
import { handleAgentBoxRequestError } from '../../utils/errors'

export const listCommand = new commander.Command('list')
  .description('list sandbox templates')
  .alias('ls')
  .addOption(projectOption)
  .option('-f, --format <format>', 'output format, eg. json, pretty')
  .action(async (opts: { project?: string; format: string }) => {
    try {
      const format = opts.format || 'pretty'
      ensureAPIKey()
      process.stdout.write('\n')

      const templates = await listSandboxTemplates({
        projectId: resolveProjectId(opts.project),
      })

      for (const template of templates) {
        sortTemplateNames(template.names)
      }

      if (format === 'pretty') {
        renderTemplateTable(templates)
      } else if (format === 'json') {
        console.log(JSON.stringify(templates, null, 2))
      } else {
        console.error(`Unsupported output format: ${format}`)
        process.exit(1)
      }
    } catch (err: any) {
      console.error(err)
      process.exit(1)
    }
  })

function renderTemplateTable(
  templates: agentbox.components['schemas']['Template'][]
) {
  if (!templates?.length) {
    console.log('No templates found.')
    return
  }

  renderTable(templates, [
    {
      header: 'Access',
      value: (template) => (template.public ? 'Public' : 'Private'),
    },
    { header: 'Template ID', value: (template) => template.templateID },
    {
      header: 'Template Name',
      value: (template) => listNames(template.names),
    },
    { header: 'vCPUs', value: (template) => String(template.cpuCount) },
    { header: 'RAM MiB', value: (template) => String(template.memoryMB) },
    {
      header: 'Created by',
      value: (template) => template.createdBy?.id,
    },
    {
      header: 'Created at',
      value: (template) => new Date(template.createdAt).toLocaleDateString(),
    },
    {
      header: 'Disk size MiB',
      value: (template) => String(template.diskSizeMB),
    },
    { header: 'Envd version', value: (template) => template.envdVersion },
  ])
}

export async function listSandboxTemplates({
  projectId,
}: {
  projectId?: string
}): Promise<agentbox.components['schemas']['Template'][]> {
  const templates = await client.api.GET('/v2/templates', {
    params: {
      // the backend API still calls the project ID "teamID"
      query: { teamID: projectId },
    },
  })

  handleAgentBoxRequestError(templates, 'Error getting templates')
  return templates.data
}
