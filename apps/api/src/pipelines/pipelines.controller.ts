import { Controller } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { BaseController } from '@/src/common/base.controller'
import { CreatePipelineDto } from '@/src/pipelines/dto/create-pipeline.dto'
import { UpdatePipelineDto } from '@/src/pipelines/dto/update-pipeline.dto'
import { PipelineEntity } from '@/src/pipelines/entities/pipeline.entity'
import { PipelinesService } from '@/src/pipelines/pipelines.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@ApiTags('Pipelines')
@Authenticated()
@Controller('/organizations/:orgname/pipelines')
export class PipelinesController extends BaseController<
  PipelineEntity,
  CreatePipelineDto,
  UpdatePipelineDto,
  PipelinesService
>(PipelineEntity, CreatePipelineDto, UpdatePipelineDto) {
  constructor(private readonly pipelinesService: PipelinesService) {
    super(pipelinesService)
  }
}
