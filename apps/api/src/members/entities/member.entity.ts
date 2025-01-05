import { Member as _PrismaMember } from '@prisma/client'
import {
  IsBoolean,
  IsEmail,
  IsEnum,
  IsOptional,
  IsString
} from 'class-validator'

import { BaseEntity } from '../../common/entities/base.entity'
import { Expose } from 'class-transformer'

export type MemberModel = _PrismaMember

export enum RoleTypeEnum {
  ADMIN = 'ADMIN',
  USER = 'USER'
}

export class MemberEntity extends BaseEntity implements MemberModel {
  /**
   * Whether the invite was accepted
   * @example false
   */
  @IsBoolean()
  @Expose()
  inviteAccepted: boolean

  /**
   * The email of the invited member
   * @example 'invited-user@archesai.com'
   */
  @IsEmail()
  @Expose()
  inviteEmail: string

  /**
   * The organization name
   * @example 'my-organization'
   */
  @IsString()
  @Expose()
  orgname: string

  /**
   * The role of the member
   * @example 'ADMIN'
   */
  @IsEnum(RoleTypeEnum)
  @Expose()
  role: RoleTypeEnum

  /**
   * The username of this member
   * @example 'jonathan'
   */
  @IsOptional()
  @IsString()
  @Expose()
  username: null | string

  constructor(member: MemberModel) {
    super()
    Object.assign(this, member)
  }
}
