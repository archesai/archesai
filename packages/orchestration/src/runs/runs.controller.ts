import type { Controller } from '@archesai/core'
import type { RunEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreateRunDtoSchema,
  RUN_ENTITY_KEY,
  RunEntitySchema,
  UpdateRunDtoSchema
} from '@archesai/schemas'

import type { RunsService } from '#runs/runs.service'

/**
 * Controller for runs.
 */
export class RunsController
  extends BaseController<RunEntity>
  implements Controller
{
  constructor(runsService: RunsService) {
    super(
      RUN_ENTITY_KEY,
      RunEntitySchema,
      CreateRunDtoSchema,
      UpdateRunDtoSchema,
      runsService
    )
  }
}
