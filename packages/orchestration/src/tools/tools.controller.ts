import type { Controller } from '@archesai/core'
import type { ToolEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import { TOOL_ENTITY_KEY, ToolEntitySchema } from '@archesai/domain'

import type { ToolsService } from '#tools/tools.service'

import { CreateToolRequestSchema } from '#tools/dto/create-tool.req.dto'
import { UpdateToolRequestSchema } from '#tools/dto/update-tool.req.dto'

/**
 * Controller for tools.
 */
export class ToolsController
  extends BaseController<ToolEntity>
  implements Controller
{
  constructor(toolsService: ToolsService) {
    super(
      TOOL_ENTITY_KEY,
      ToolEntitySchema,
      CreateToolRequestSchema,
      UpdateToolRequestSchema,
      toolsService
    )
  }
}
