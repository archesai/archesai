import { BaseEntity } from '@/src/common/entities/base.entity'
import { Expose } from 'class-transformer'
import { IsBoolean, IsNumber, IsString } from 'class-validator'

export class StorageItemDto extends BaseEntity {
  /**
   * Whether or not this is a directory
   * @example true
   */
  @IsBoolean()
  @Expose()
  isDir: boolean

  /**
   * The path of the storage item
   * @example '/location/in/storage'
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The size of the item in bytes
   * @example 12341234
   */
  @IsNumber()
  @Expose()
  size: number

  constructor(storageItemDto: StorageItemDto) {
    super()
    Object.assign(this, storageItemDto)
  }
}
