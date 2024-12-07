import { ApiHideProperty } from '@nestjs/swagger'
import { AuthProvider, Member, User } from '@prisma/client'
import { Exclude, Expose } from 'class-transformer'
import { IsEmail, IsString, MinLength } from 'class-validator'

import { BaseEntity } from '../../common/entities/base.entity'
import { MemberEntity } from '../../members/entities/member.entity'
import { AuthProviderEntity } from './auth-provider.entity'

export type UserWithMembershipsAndAuthProvidersModel = User & {
  authProviders: AuthProvider[]
  memberships: Member[]
}

export class UserEntity extends BaseEntity implements User {
  /**
   * The memberships of the currently signed-in user
   */
  @Expose()
  authProviders: AuthProviderEntity[]

  /**
   * Whether or not the user is deactivated
   * @example false
   */
  @Expose()
  deactivated!: boolean

  /**
   * The user's default organization name
   * @example 'my-organization'
   */
  @Expose()
  defaultOrgname: string

  /**
   * The user's display name
   * @example 'John Smith'
   */
  @Expose()
  displayName: string

  /**
   * The user's e-mail
   * @example 'example@archesai.com'
   */
  @Expose()
  @IsEmail()
  email!: string

  /**
   * Whether or not the user's e-mail has been verified
   */
  @Expose()
  emailVerified!: boolean

  /**
   * The user's first name
   * @example 'John'
   */
  @Expose()
  firstName: string

  /**
   * The user's last name
   * @example 'Smith'
   */
  @Expose()
  lastName: string

  /**
   * The memberships of the currently signed-in user
   */
  @Expose()
  memberships: MemberEntity[]

  @ApiHideProperty()
  @Exclude()
  password: string

  /**
   * The user's photo URL
   * @example '/avatar.png'
   */
  @Expose()
  @IsString()
  photoUrl!: string
  @ApiHideProperty()
  @Exclude()
  refreshToken: string

  /**
   * The user's username
   * @example 'jonathan'
   */
  @Expose()
  @MinLength(5)
  username!: string

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
