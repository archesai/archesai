import { Controller } from '@nestjs/common'
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger'

import { BaseController } from '../common/base.controller'
import { CreateToolDto } from './dto/create-tool.dto'
import { UpdateToolDto } from './dto/update-tool.dto'
import { ToolEntity } from './entities/tool.entity'
import { ToolsService } from './tools.service'

@ApiBearerAuth()
@ApiTags('Tools')
@Controller('/organizations/:orgname/tools')
export class ToolsController extends BaseController<ToolEntity, CreateToolDto, UpdateToolDto, ToolsService> {
  constructor(private readonly toolsService: ToolsService) {
    super(toolsService)
  }
}
