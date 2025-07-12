import type { WebsocketsService } from '@archesai/core'
import type { AccountInsertModel } from '@archesai/database'
import type { AccountEntity, ProviderType } from '@archesai/schemas'

import { BaseService } from '@archesai/core'

import type { AccountRepository } from '#accounts/account.repository'
import type { HashingService } from '#hashing/hashing.service'

/**
 * Service for managing accounts.
 */
export class AccountsService extends BaseService<AccountEntity> {
  private readonly accountRepository: AccountRepository
  private readonly hashingService: HashingService
  private readonly websocketsService: WebsocketsService

  constructor(
    accountRepository: AccountRepository,
    hashingService: HashingService,
    websocketsService: WebsocketsService
  ) {
    super(accountRepository)
    this.accountRepository = accountRepository
    this.hashingService = hashingService
    this.websocketsService = websocketsService
  }

  /**
   * Creates a new account entity with the provided data.
   * @param value The data to create the account entity with.
   * @returns The created account entity.
   */
  public override async create(
    value: AccountInsertModel
  ): Promise<AccountEntity> {
    if (value.providerId === 'LOCAL') {
      if (!value.password) {
        throw new Error('A hashed password is required for local')
      }
      value.password = await this.hashingService.hashPassword(value.password)
    }
    return this.accountRepository.create(value)
  }

  /**
   * Finds an account entity by the provider and provider account ID.
   * @param provider The provider to search by.
   * @param providerAccountId The provider account ID to search by.
   * @returns The account entity if found.
   */
  public async findByProviderAndProviderAccountId(
    provider: ProviderType,
    providerAccountId: string
  ): Promise<AccountEntity> {
    return this.accountRepository.findFirst({
      filter: {
        accountId: {
          equals: providerAccountId
        },
        providerId: {
          equals: provider
        }
      }
    })
  }

  protected override emitMutationEvent(entity: AccountEntity): void {
    this.websocketsService.broadcastEvent(entity.userId, 'update', {
      queryKey: ['auth', 'accounts']
    })
  }
}
