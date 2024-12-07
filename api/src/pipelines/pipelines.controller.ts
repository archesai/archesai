import { Controller, Logger } from '@nestjs/common'
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger'

import { BaseController } from '../common/base.controller'
import { CreatePipelineDto } from './dto/create-pipeline.dto'
import { UpdatePipelineDto } from './dto/update-pipeline.dto'
import { PipelineEntity } from './entities/pipeline.entity'
import { PipelinesService } from './pipelines.service'

@ApiBearerAuth()
@ApiTags('Pipelines')
@Controller('/organizations/:orgname/pipelines')
export class PipelinesController extends BaseController<
  PipelineEntity,
  CreatePipelineDto,
  UpdatePipelineDto,
  PipelinesService
> {
  private logger = new Logger(PipelinesController.name)

  constructor(private readonly pipelinesService: PipelinesService) {
    super(pipelinesService)
  }
}
