import type { WebsocketsService } from '@archesai/core'
import type { AccountEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { ACCOUNT_ENTITY_KEY } from '@archesai/schemas'

import type { AccountRepository } from '#accounts/account.repository'

export const createAccountsService = (
  accountRepository: AccountRepository,
  websocketsService: WebsocketsService
) => {
  const emitAccountsMutationEvent = (entity: AccountEntity): void => {
    websocketsService.broadcastEvent(entity.userId, 'update', {
      queryKey: ['users', entity.userId, ACCOUNT_ENTITY_KEY]
    })
  }
  return createBaseService(accountRepository, emitAccountsMutationEvent)
}

export type AccountsService = ReturnType<typeof createAccountsService>
