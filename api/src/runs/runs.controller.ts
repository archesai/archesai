import { Controller } from '@nestjs/common'
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger'

import { BaseController } from '../common/base.controller'
import { CreateRunDto } from './dto/create-run.dto'
import { RunEntity } from './entities/run.entity'
import { RunsService } from './runs.service'

@ApiBearerAuth()
@ApiTags('Runs')
@Controller('/organizations/:orgname/runs')
export class RunsController extends BaseController<
  RunEntity,
  CreateRunDto,
  any,
  RunsService
>(RunEntity, CreateRunDto, String) {
  constructor(private readonly runsService: RunsService) {
    super(runsService)
  }
}
