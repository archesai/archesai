import type { ConfigService, EventBus, WebsocketsService } from '@archesai/core'
import type {
  BaseInsertion,
  OrganizationEntity,
  PlanType
} from '@archesai/schemas'

import { BaseService, ForbiddenException } from '@archesai/core'
import {
  ORGANIZATION_ENTITY_KEY,
  OrganizationCreatedEvent
} from '@archesai/schemas'

import type { OrganizationRepository } from '#organizations/organization.repository'

/**
 * Service for handling organizations.
 */
export class OrganizationsService extends BaseService<OrganizationEntity> {
  private readonly configService: ConfigService
  private readonly eventBus: EventBus
  private readonly organizationRepository: OrganizationRepository
  private readonly websocketsService: WebsocketsService

  constructor(
    configService: ConfigService,
    eventBus: EventBus,
    organizationRepository: OrganizationRepository,
    websocketsService: WebsocketsService
  ) {
    super(organizationRepository)
    this.configService = configService
    this.eventBus = eventBus
    this.organizationRepository = organizationRepository
    this.websocketsService = websocketsService
  }

  public async addOrRemoveCredits(organizationId: string, numCredits: number) {
    this.logger.debug(
      `Adding ${numCredits.toString()} credits to ${organizationId}`
    )

    let organization = await this.findByOrgname(organizationId)
    organization = await this.organizationRepository.update(
      ORGANIZATION_ENTITY_KEY,
      {
        billingEmail: organization.billingEmail,
        credits:
          numCredits < 0 ?
            organization.credits + numCredits
          : organization.credits - -1 * numCredits
      }
    )
    this.emitMutationEvent(organization)
    return organization
  }

  public async checkCredits(organizationId: string, numCredits: number) {
    this.logger.debug(
      `Checking ${numCredits.toString()} credits for ${organizationId}`
    )
    const organization = await this.findByOrgname(organizationId)
    if (organization.plan != 'PREMIUM' && organization.credits <= numCredits) {
      throw new ForbiddenException(
        'Sorry, you do not have enough credits. Please purchase more credits to continue' +
          (organization.credits < numCredits ?
            ` (estimated cost: ${numCredits.toString()})`
          : '')
      )
    }
  }

  public override async create(value: BaseInsertion<OrganizationEntity>) {
    const billingEnabled = this.configService.get('billing.enabled')
    const organization = await this.organizationRepository.create({
      ...value,
      billingEmail: value.billingEmail,
      organizationId: value.organizationId,
      plan: billingEnabled ? 'FREE' : 'UNLIMITED'
    })
    this.eventBus.emit(
      'organization.created',
      new OrganizationCreatedEvent({
        organization
      })
    )
    return organization
  }

  public async findByOrgname(organizationId: string) {
    return this.organizationRepository.findFirst({
      filter: {
        organizationId: {
          equals: organizationId
        }
      }
    })
  }

  public async setPlan(organizationId: string, plan: PlanType) {
    const organization = await this.findByOrgname(organizationId)
    const organizationEntity = await this.organizationRepository.update(
      organization.id,
      {
        plan
      }
    )

    this.emitMutationEvent(organizationEntity)
    return organizationEntity
  }

  protected emitMutationEvent(entity: OrganizationEntity): void {
    this.websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId]
    })
  }
}
