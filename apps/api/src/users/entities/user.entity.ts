import { ApiHideProperty } from '@nestjs/swagger'
import { Prisma, User } from '@prisma/client'
import { Expose } from 'class-transformer'
import {
  IsBoolean,
  IsEmail,
  IsString,
  MinLength,
  ValidateNested
} from 'class-validator'

import { BaseEntity } from '../../common/entities/base.entity'
import { MemberEntity } from '../../members/entities/member.entity'
import { AuthProviderEntity } from './auth-provider.entity'

export type UserWithMembershipsAndAuthProvidersModel = Prisma.UserGetPayload<{
  include: {
    authProviders: true
    memberships: true
  }
}>

export class UserEntity extends BaseEntity implements User {
  /**
   * The memberships of the currently signed-in user
   */
  @ValidateNested({ each: true })
  @Expose()
  authProviders: AuthProviderEntity[]

  /**
   * Whether or not the user is deactivated
   * @example false
   */
  @IsBoolean()
  @Expose()
  deactivated: boolean

  /**
   * The user's default organization name
   * @example 'my-organization'
   */
  @IsString()
  @Expose()
  defaultOrgname: string

  /**
   * The user's display name
   * @example 'John Smith'
   */
  @IsString()
  @Expose()
  displayName: string

  /**
   * The user's e-mail
   * @example 'example@archesai.com'
   */
  @IsEmail()
  @Expose()
  email: string

  /**
   * Whether or not the user's e-mail has been verified
   */
  @IsBoolean()
  @Expose()
  emailVerified: boolean

  /**
   * The user's first name
   * @example 'John'
   */
  @IsString()
  @Expose()
  firstName: string

  /**
   * The user's last name
   * @example 'Smith'
   */
  @IsString()
  @Expose()
  lastName: string

  /**
   * The memberships of the currently signed-in user
   */
  @ValidateNested({ each: true })
  @Expose()
  memberships: MemberEntity[]

  @ApiHideProperty()
  password: string

  /**
   * The user's photo URL
   * @example '/avatar.png'
   */
  @IsString()
  @Expose()
  photoUrl: string

  @ApiHideProperty()
  refreshToken: string

  /**
   * The user's username
   * @example 'jonathan'
   */
  @IsString()
  @MinLength(5)
  @Expose()
  username: string

  constructor(user: UserWithMembershipsAndAuthProvidersModel) {
    super()
    Object.assign(this, user)
    this.memberships = (this.memberships || []).map(
      (membership) => new MemberEntity(membership)
    )
    this.authProviders = (this.authProviders || []).map(
      (authProvider) => new AuthProviderEntity(authProvider)
    )
    this.displayName = this.firstName
      ? this.firstName + ' ' + this.lastName
      : this.username
  }
}
