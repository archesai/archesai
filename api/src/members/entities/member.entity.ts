import { ApiProperty } from '@nestjs/swagger'
import { Member as _PrismaMember, RoleType } from '@prisma/client'
import { Expose } from 'class-transformer'
import { IsEnum, IsNotEmpty } from 'class-validator'

import { BaseEntity } from '../../common/entities/base.entity'

export type MemberModel = _PrismaMember

export class MemberEntity extends BaseEntity implements MemberModel {
  @ApiProperty({
    description: 'Whether the invite was accepted',
    example: false
  })
  @Expose()
  @IsNotEmpty()
  inviteAccepted: boolean

  // Exposed Properties
  @ApiProperty({
    description: 'The invited email of this member',
    example: 'invited-user@archesai.com'
  })
  @Expose()
  @IsNotEmpty()
  inviteEmail: string

  @ApiProperty({
    description: 'The organization name',
    example: 'my-organization'
  })
  @Expose()
  orgname: string

  @ApiProperty({ description: 'The role of the member', enum: RoleType })
  @Expose()
  @IsEnum(RoleType)
  role: RoleType

  @ApiProperty({
    description: 'The username of this member',
    example: 'jonathan',
    required: false,
    type: String
  })
  @Expose()
  username: null | string

  constructor(member: MemberModel) {
    super()
    Object.assign(this, member)
  }
}
