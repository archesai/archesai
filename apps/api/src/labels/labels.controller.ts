import { Controller } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { BaseController } from '@/src/common/base.controller'
import { CreateLabelDto } from '@/src/labels/dto/create-label.dto'
import { UpdateLabelDto } from '@/src/labels/dto/update-label.dto'
import { LabelEntity } from '@/src/labels/entities/label.entity'
import { LabelsService } from '@/src/labels/labels.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@Authenticated()
@ApiTags('Labels')
@Controller('/organizations/:orgname/labels')
export class LabelsController extends BaseController<
  LabelEntity,
  CreateLabelDto,
  UpdateLabelDto,
  LabelsService
>(LabelEntity, CreateLabelDto, UpdateLabelDto) {
  constructor(private readonly labelsService: LabelsService) {
    super(labelsService)
  }
}
