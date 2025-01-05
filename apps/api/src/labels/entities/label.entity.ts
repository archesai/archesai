import { Label as _PrismaLabel } from '@prisma/client'

import { BaseEntity } from '../../common/entities/base.entity'
import { IsString } from 'class-validator'
import { Expose } from 'class-transformer'

export type LabelModel = _PrismaLabel

export class LabelEntity extends BaseEntity implements LabelModel {
  /**
   * The chat label name
   * @example 'What are the morals of the story in Aesop's Fables?'
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The organization name
   * @example 'my-organization'
   */
  @IsString()
  @Expose()
  orgname: string

  constructor(label: LabelModel) {
    super()
    Object.assign(this, label)
  }
}
