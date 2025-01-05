import { IntersectionType, PickType } from '@nestjs/swagger'

import { ContentEntity } from '../entities/content.entity'

export class CreateContentDto extends IntersectionType(
  PickType(ContentEntity, ['name', 'url', 'text', 'labels'] as const)
) {}
