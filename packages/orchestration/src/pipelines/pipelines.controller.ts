import type { Controller } from '@archesai/core'
import type { PipelineEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreatePipelineDtoSchema,
  PIPELINE_ENTITY_KEY,
  PipelineEntitySchema,
  UpdatePipelineDtoSchema
} from '@archesai/schemas'

import type { PipelinesService } from '#pipelines/pipelines.service'

/**
 * Controller for pipelines.
 */
export class PipelinesController
  extends BaseController<PipelineEntity>
  implements Controller
{
  constructor(pipelinesService: PipelinesService) {
    super(
      PIPELINE_ENTITY_KEY,
      PipelineEntitySchema,
      CreatePipelineDtoSchema,
      UpdatePipelineDtoSchema,
      pipelinesService
    )
  }
}
