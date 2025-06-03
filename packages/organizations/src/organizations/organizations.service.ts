import type { ConfigService, EventBus, WebsocketsService } from '@archesai/core'
import type {
  BaseInsertion,
  OrganizationEntity,
  PlanType
} from '@archesai/domain'

import { BaseService, ForbiddenException } from '@archesai/core'
import {
  ORGANIZATION_ENTITY_KEY,
  OrganizationCreatedEvent
} from '@archesai/domain'

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

  public async addOrRemoveCredits(orgname: string, numCredits: number) {
    this.logger.debug(`Adding ${numCredits.toString()} credits to ${orgname}`)
    const organization = await this.organizationRepository.addOrRemoveCredits(
      orgname,
      numCredits
    )
    this.emitMutationEvent(organization)
    return organization
  }

  public async checkCredits(orgname: string, numCredits: number) {
    this.logger.debug(
      `Checking ${numCredits.toString()} credits for ${orgname}`
    )
    const organization = await this.findByOrgname(orgname)
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
      plan: billingEnabled ? 'FREE' : 'UNLIMITED',
      type: ORGANIZATION_ENTITY_KEY
    })
    this.eventBus.emit(
      'organization.created',
      new OrganizationCreatedEvent({
        organization
      })
    )
    return organization
  }

  public async findByOrgname(orgname: string) {
    return this.organizationRepository.findByOrgname(orgname)
  }

  public async setPlan(orgname: string, plan: PlanType) {
    const organization = await this.findByOrgname(orgname)
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
    this.websocketsService.broadcastEvent(entity.orgname, 'update', {
      queryKey: ['organizations', entity.orgname]
    })
  }
}
