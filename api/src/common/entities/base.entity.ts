import { ApiHideProperty } from '@nestjs/swagger'
import { Expose } from 'class-transformer'
import { IsDateString, IsUUID } from 'class-validator'

export class BaseEntity {
  /**
   * The date that this item was created
   * @example '2023-07-11T21:09:20.895Z'
   */
  @IsDateString()
  @Expose()
  createdAt: Date

  /**
   * The ID of the item
   * @example 'item-id'
   */
  @IsUUID()
  @Expose()
  id: string

  /**
   * The date that this item was last updated
   * @example '2023-07-11T21:09:20.895Z'
   */
  @ApiHideProperty()
  updatedAt: Date
}
