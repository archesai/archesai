import { Controller } from '@nestjs/common'

import { BaseController } from '@/src/common/base.controller'
import { ContentService } from '@/src/content/content.service'
import { CreateContentDto } from '@/src/content/dto/create-content.dto'
import { UpdateContentDto } from '@/src/content/dto/update-content.dto'
import { ContentEntity } from '@/src/content/entities/content.entity'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@Authenticated()
@Controller('organizations/:orgname/content')
export class ContentController extends BaseController<
  ContentEntity,
  CreateContentDto,
  UpdateContentDto,
  ContentService
>(ContentEntity, CreateContentDto, UpdateContentDto) {
  constructor(private readonly contentService: ContentService) {
    super(contentService)
  }
}
