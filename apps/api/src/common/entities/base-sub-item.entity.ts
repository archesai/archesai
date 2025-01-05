import { Expose } from 'class-transformer'
import { IsString } from 'class-validator'

export type _PrismaSubItemModel = {
  id: string
  name: string
}

export class SubItemEntity implements _PrismaSubItemModel {
  /**
   * The id of the item
   * @example 'item-id'
   */
  @IsString()
  @Expose()
  id: string

  /**
   * The name of the item
   * @example 'item-name'
   */
  @IsString()
  @Expose()
  name: string

  constructor(subItem: _PrismaSubItemModel) {
    Object.assign(this, subItem)
  }
}
