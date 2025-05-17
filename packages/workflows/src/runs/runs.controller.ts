import type { Controller } from '@archesai/core'
import type { RunEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import { RUN_ENTITY_KEY, RunEntitySchema } from '@archesai/domain'

import type { RunsService } from '#runs/runs.service'

import { CreateRunRequestSchema } from '#runs/dto/create-run.req.dto'
import { UpdateRunRequestSchema } from '#runs/dto/update-run.req.dto'

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
      CreateRunRequestSchema,
      UpdateRunRequestSchema,
      runsService
    )
  }
}
