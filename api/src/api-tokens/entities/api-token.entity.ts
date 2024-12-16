import { ApiToken as _PrismaApiToken } from '@prisma/client'

import { BaseEntity } from '@/src/common/entities/base.entity'
import { IsEnum, IsString } from 'class-validator'
import { Expose } from 'class-transformer'

export type ApiTokenModel = _PrismaApiToken

export enum RoleTypeEnum {
  ADMIN = 'ADMIN',
  USER = 'USER'
}

export class ApiTokenEntity extends BaseEntity implements ApiTokenModel {
  /**
   *The domains that can access this API token
   * @example archesai.com,localhost:3000
   */
  @IsString()
  @Expose()
  domains: string

  /**
   * The API token key. This will only be shown once
   * @example ********1234567890
   */
  @IsString()
  @Expose()
  key: string

  /**
   *  The name of the API token
   * @example My Token
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The organization name
   * @example my-organization
   */
  @IsString()
  @Expose()
  orgname: string

  /**
   * The role of the API token
   * @example ADMIN
   */
  @IsEnum(RoleTypeEnum)
  @Expose()
  role: RoleTypeEnum

  /**
   * The username of the user who owns this API token
   * @example jonathan
   */
  @IsString()
  @Expose()
  username: string

  constructor(apiToken: ApiTokenModel) {
    super()
    Object.assign(this, apiToken)
  }
}
