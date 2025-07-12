import type { WebsocketsService } from '@archesai/core'
import type { UserInsertModel } from '@archesai/database'
import type { UserEntity } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { USER_ENTITY_KEY } from '@archesai/schemas'

import type { OrganizationsService } from '#organizations/organizations.service'
import type { UserRepository } from '#users/user.repository'

/**
 * Service for handling users.
 */
export class UsersService extends BaseService<UserEntity> {
  private readonly organizationsService: OrganizationsService
  private readonly userRepository: UserRepository
  private readonly websocketsService: WebsocketsService

  constructor(
    organizationsService: OrganizationsService,
    userRepository: UserRepository,
    websocketsService: WebsocketsService
  ) {
    super(userRepository)
    this.organizationsService = organizationsService
    this.userRepository = userRepository
    this.websocketsService = websocketsService
  }

  public async checkIfEmailExists(email: string): Promise<boolean> {
    try {
      await this.userRepository.findFirst({
        filter: {
          email: {
            equals: email
          }
        },
        page: {
          number: 1,
          size: 1
        },
        sort: '-createdAt'
      })
      return true
    } catch {
      return false
    }
  }

  public override async create(value: UserInsertModel): Promise<UserEntity> {
    // Create user
    const user = await this.userRepository.create({
      ...value
      // username:
      //   value.organizationId ||
      //   (value.email.split('@')[0] ?? value.email) +
      //     '-' +
      //     Math.random().toString(36).substring(2, 6)
    })

    // Create organization
    const usernamePrefix = user.email.split('@')[0]
    if (!usernamePrefix) {
      throw new Error('Could not extract username from email')
    }
    const organizationId = usernamePrefix + user.id.slice(0, 5)
    await this.organizationsService.create({
      billingEmail: user.email,
      createdAt: new Date().toISOString(),
      credits: 0,
      organizationId,
      plan: 'FREE',
      updatedAt: new Date().toISOString()
    })
    return user
  }

  public async deactivate(id: string): Promise<void> {
    await this.userRepository.update(id, {
      deactivated: true
    })
  }

  public async findOneByEmail(email: string): Promise<UserEntity> {
    return this.userRepository.findFirst({
      filter: {
        email: {
          equals: email
        }
      },
      page: {
        number: 1,
        size: 1
      },
      sort: '-createdAt'
    })
  }

  protected emitMutationEvent(user: UserEntity): void {
    this.websocketsService.broadcastEvent(user.id, 'update', {
      queryKey: [USER_ENTITY_KEY]
    })
  }
}
