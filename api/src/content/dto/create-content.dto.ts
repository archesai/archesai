import { IntersectionType, PartialType, PickType } from '@nestjs/swagger'

import { ContentEntity } from '../entities/content.entity'
import { IsArray, IsString } from 'class-validator'
import { Expose } from 'class-transformer'

export class CreateContentDto extends IntersectionType(
  PickType(ContentEntity, ['name'] as const),
  PartialType(PickType(ContentEntity, ['url', 'text'] as const))
) {
  /**
   * The labels to associate with the content
   * @example ['label-1', 'label-2']
   */
  @IsArray()
  @IsString({ each: true })
  @Expose()
  labels: string[] = []
}
