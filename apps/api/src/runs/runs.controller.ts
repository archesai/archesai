import { Controller } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { BaseController } from '@/src/common/base.controller'
import { CreateRunDto } from '@/src/runs/dto/create-run.dto'
import { RunEntity } from '@/src/runs/entities/run.entity'
import { RunsService } from '@/src/runs/runs.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@ApiTags('Runs')
@Authenticated()
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
