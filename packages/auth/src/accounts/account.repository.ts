import type { DatabaseService } from '@archesai/core'
import type { ProviderType } from '@archesai/domain'

import { BaseRepository } from '@archesai/core'
import { ACCOUNT_ENTITY_KEY, AccountEntity } from '@archesai/domain'

/**
 * Repository for managing accounts.
 */
export class AccountRepository extends BaseRepository<AccountEntity> {
  constructor(databaseService: DatabaseService<AccountEntity>) {
    super(databaseService, ACCOUNT_ENTITY_KEY, AccountEntity)
  }

  public async findByProviderAndProviderAccountId(
    provider: ProviderType,
    providerAccountId: string
  ): Promise<AccountEntity> {
    return this.findFirst({
      filter: {
        provider: {
          equals: provider
        },
        providerAccountId: {
          equals: providerAccountId
        }
      }
    })
  }

  public async updateRefreshToken(
    userId: string,
    newRefreshToken: string
  ): Promise<AccountEntity> {
    const refreshToken = await this.findFirst({
      filter: {
        provider: {
          equals: 'LOCAL'
        },
        userId: {
          equals: userId
        }
      }
    })
    return this.update(refreshToken.id, {
      refresh_token: newRefreshToken
    })
  }
}
