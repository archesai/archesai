import { Injectable } from '@nestjs/common'
import { AuthProviderType } from '@prisma/client'

import { BaseService } from '../common/base.service'
import { OrganizationsService } from '../organizations/organizations.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { CreateUserDto } from './dto/create-user.dto'
import {
  UserEntity,
  UserWithMembershipsAndAuthProvidersModel
} from './entities/user.entity'
import { UserRepository } from './user.repository'

@Injectable()
export class UsersService extends BaseService<
  UserEntity,
  UserWithMembershipsAndAuthProvidersModel,
  UserRepository
> {
  constructor(
    private userRepository: UserRepository,
    private organizationsService: OrganizationsService,
    private websocketsService: WebsocketsService
  ) {
    super(userRepository)
  }

  async create(createUserDto: CreateUserDto) {
    let user = await this.userRepository.create({
      ...createUserDto
    })
    const organization = await this.organizationsService.create({
      billingEmail: user.email,
      orgname: user.username
    })
    await this.organizationsService.addUserToOrganization(
      organization.orgname,
      this.toEntity(user)
    )
    user = await this.userRepository.update(user.id, {
      defaultOrgname: organization.orgname
    })
    return this.toEntity(user)
  }

  async deactivate(id: string) {
    await this.userRepository.deactivate(id)
  }

  async findOneByEmail(email: string) {
    return this.toEntity(await this.userRepository.findOneByEmail(email))
  }

  async findOneByUsername(username: string) {
    return this.toEntity(await this.userRepository.findOneByUsername(username))
  }

  async setEmail(id: string, newEmail: string) {
    return this.toEntity(
      await this.userRepository.update(id, {
        email: newEmail
      })
    )
  }

  async setEmailVerified(id: string) {
    return this.toEntity(
      await this.userRepository.update(id, {
        emailVerified: true
      })
    )
  }

  async setRefreshToken(id: string, refreshToken: string) {
    return this.toEntity(
      await this.userRepository.update(id, {
        refreshToken
      })
    )
  }

  async syncAuthProvider(
    email: string,
    provider: AuthProviderType,
    providerId: string
  ): Promise<UserEntity> {
    const user = await this.userRepository.findOneByEmail(email)
    // if it does not have this provider, add it
    if (!user.authProviders.some((p) => p.provider === provider)) {
      return this.toEntity(
        await this.userRepository.addAuthProvider(email, provider, providerId)
      )
    }
    const userEntity = this.toEntity(user)
    this.emitMutationEvent(userEntity)
    return userEntity
  }

  protected emitMutationEvent(entity: UserEntity): void {
    this.websocketsService.socket?.to(entity.defaultOrgname).emit('update', {
      queryKey: ['user']
    })
  }

  protected toEntity(
    model: UserWithMembershipsAndAuthProvidersModel
  ): UserEntity {
    return new UserEntity(model)
  }
}
