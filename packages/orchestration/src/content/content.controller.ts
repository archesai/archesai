import type { Controller } from '@archesai/core'
import type { ContentEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import { CONTENT_ENTITY_KEY, ContentEntitySchema } from '@archesai/domain'

import type { ContentService } from '#content/content.service'

import { CreateContentRequestSchema } from '#content/dto/create-content.req.dto'
import { UpdateContentRequestSchema } from '#content/dto/update-content.req.dto'

/**
 * Controller for content.
 */
export class ContentController
  extends BaseController<ContentEntity>
  implements Controller
{
  constructor(contentService: ContentService) {
    super(
      CONTENT_ENTITY_KEY,
      ContentEntitySchema,
      CreateContentRequestSchema,
      UpdateContentRequestSchema,
      contentService
    )
  }
}
