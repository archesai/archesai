import { ApiHideProperty, ApiProperty } from '@nestjs/swagger'
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
  @ApiProperty({
    description: 'The memberships of the currently signed in user',
    type: [AuthProviderEntity]
  })
  @Expose()
  authProviders: AuthProviderEntity[]

  @ApiProperty({
    description: 'Whether or not the user is deactivated',
    example: false
  })
  @Expose()
  deactivated!: boolean

  @ApiProperty({
    description: "The user's default organization name",
    example: 'my-organization'
  })
  @Expose()
  defaultOrgname: string

  @ApiProperty({
    description: "The user's display name",
    example: 'John Smith'
  })
  @Expose()
  displayName: string

  @ApiProperty({
    description: "The user's e-mail",
    example: 'example@archesai.com'
  })
  @Expose()
  @IsEmail()
  email!: string

  @ApiProperty({
    description: "Whether or not the user's e-mail has been verified"
  })
  @Expose()
  emailVerified!: boolean

  @ApiProperty({
    description: "The user's first name",
    example: 'John'
  })
  @Expose()
  firstName: string

  @ApiProperty({
    description: "The user's last name",
    example: 'Smith'
  })
  @Expose()
  lastName: string

  @ApiProperty({
    description: 'The memberships of the currently signed in user',
    type: [MemberEntity]
  })
  @Expose()
  memberships: MemberEntity[]

  @ApiHideProperty()
  @Exclude()
  password: string

  @ApiProperty({
    description: "The user's photo url",
    example: '/avatar.png'
  })
  @Expose()
  @IsString()
  photoUrl!: string

  @ApiHideProperty()
  @Exclude()
  refreshToken: string

  // Exposed Properties
  @ApiProperty({
    description: "The user's username",
    example: 'jonathan',
    minLength: 5
  })
  @Expose()
  @MinLength(5)
  username!: string

  constructor(user: UserWithMembershipsAndAuthProvidersModel) {
    super()
    Object.assign(this, user)
    this.memberships = (this.memberships || []).map((membership) => new MemberEntity(membership))
    this.authProviders = (this.authProviders || []).map((authProvider) => new AuthProviderEntity(authProvider))
    this.displayName = this.firstName ? this.firstName + ' ' + this.lastName : this.username
  }
}
