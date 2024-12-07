import { ApiProperty } from '@nestjs/swagger'
import { ApiToken as _PrismaApiToken, RoleType } from '@prisma/client'
import { IsEnum, IsString } from 'class-validator'

import { BaseEntity } from '../../common/entities/base.entity'

export type ApiTokenModel = _PrismaApiToken

export class ApiTokenEntity extends BaseEntity implements ApiTokenModel {
  /**
   *The domains that can access this API token
   * @example archesai.com,localhost:3000
   */
  @IsString()
  domains: string = '*'

  /**
   * The API token key. This will only be shown once
   * @example ********1234567890
   */
  key: string

  /**
   *  The name of the API token
   * @example My Token
   */
  @IsString()
  name: string

  /**
   * The organization name
   * @example my-organization
   */
  orgname: string

  @ApiProperty({
    description: 'The role of the API token',
    enum: RoleType,
    example: RoleType.ADMIN
  })
  @IsEnum(RoleType)
  role: RoleType

  /**
   * The username of the user who owns this API token
   * @example jonathan
   */
  username: string

  constructor(apiToken: ApiTokenModel) {
    super()
    Object.assign(this, apiToken)
  }
}
