import type { BaseService, WebsocketsService } from '@archesai/core'
import type { OrganizationEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { ORGANIZATION_ENTITY_KEY } from '@archesai/schemas'

import type { OrganizationRepository } from '#organizations/organization.repository'

export const createOrganizationsService = (
  organizationRepository: OrganizationRepository,
  websocketsService: WebsocketsService
): BaseService<OrganizationEntity> => {
  const emitOrganizationMutationEvent = (entity: OrganizationEntity): void => {
    websocketsService.broadcastEvent(entity.id, 'update', {
      queryKey: ['organizations', entity.id, ORGANIZATION_ENTITY_KEY]
    })
  }
  return createBaseService(
    organizationRepository,
    emitOrganizationMutationEvent
  )
}

export type OrganizationsService = ReturnType<typeof createOrganizationsService>
