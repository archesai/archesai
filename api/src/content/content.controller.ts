import { Controller, UseGuards } from '@nestjs/common'
import { ApiBearerAuth } from '@nestjs/swagger'

import { EmailVerifiedGuard } from '../auth/guards/email-verified.guard'
import { BaseController } from '../common/base.controller'
import { ContentService } from './content.service'
import { CreateContentDto } from './dto/create-content.dto'
import { UpdateContentDto } from './dto/update-content.dto'
import { ContentEntity } from './entities/content.entity'

@ApiBearerAuth()
@Controller('/organizations/:orgname/content')
@UseGuards(EmailVerifiedGuard)
export class ContentController extends BaseController<
  ContentEntity,
  CreateContentDto,
  UpdateContentDto,
  ContentService
> {
  constructor(private readonly contentService: ContentService) {
    super(contentService)
  }
}
