import type { Controller } from '@archesai/core'
import type { PipelineEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import { PIPELINE_ENTITY_KEY, PipelineEntitySchema } from '@archesai/domain'

import type { PipelinesService } from '#pipelines/pipelines.service'

import { CreatePipelineRequestSchema } from '#pipelines/dto/create-pipeline.req.dto'
import { UpdatePipelineRequestSchema } from '#pipelines/dto/update-pipeline.req.dto'

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
      CreatePipelineRequestSchema,
      UpdatePipelineRequestSchema,
      pipelinesService
    )
  }
}
