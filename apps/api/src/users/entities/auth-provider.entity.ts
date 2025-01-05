import { BaseEntity } from '@/src/common/entities/base.entity'
import { AuthProvider as _PrismaAuthProvider } from '@prisma/client'
import { Expose } from 'class-transformer'
import { IsEnum, IsString } from 'class-validator'

export type AuthProviderModel = _PrismaAuthProvider

export enum AuthProviderTypeEnum {
  LOCAL = 'LOCAL',
  FIREBASE = 'FIREBASE',
  TWITTER = 'TWITTER'
}

export class AuthProviderEntity
  extends BaseEntity
  implements AuthProviderModel
{
  /**
   * The auth provider's provider
   * @example LOCAL
   */
  @IsEnum(AuthProviderTypeEnum)
  @Expose()
  provider: AuthProviderTypeEnum

  /**
   * The provider ID associated with the auth provider
   */
  @IsString()
  @Expose()
  providerId: string

  /**
   * The user ID associated with the auth provider
   * @example 123456
   */
  @IsString()
  @Expose()
  userId: string

  constructor(authProvider: AuthProviderModel) {
    super()
    Object.assign(this, authProvider)
  }
}
