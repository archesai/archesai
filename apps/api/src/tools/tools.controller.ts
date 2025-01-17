import { Controller } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { BaseController } from '@/src/common/base.controller'
import { CreateToolDto } from '@/src/tools/dto/create-tool.dto'
import { UpdateToolDto } from '@/src/tools/dto/update-tool.dto'
import { ToolEntity } from '@/src/tools/entities/tool.entity'
import { ToolsService } from '@/src/tools/tools.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@ApiTags('Tools')
@Authenticated()
@Controller('/organizations/:orgname/tools')
export class ToolsController extends BaseController<
  ToolEntity,
  CreateToolDto,
  UpdateToolDto,
  ToolsService
>(ToolEntity, CreateToolDto, UpdateToolDto) {
  constructor(private readonly toolsService: ToolsService) {
    super(toolsService)
  }
}
