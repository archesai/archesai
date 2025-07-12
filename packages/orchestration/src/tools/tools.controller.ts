import type { Controller } from '@archesai/core'
import type { ToolEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreateToolDtoSchema,
  TOOL_ENTITY_KEY,
  ToolEntitySchema,
  UpdateToolDtoSchema
} from '@archesai/schemas'

import type { ToolsService } from '#tools/tools.service'

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
      CreateToolDtoSchema,
      UpdateToolDtoSchema,
      toolsService
    )
  }
}
