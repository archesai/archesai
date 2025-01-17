import { Injectable } from '@nestjs/common'
import { ForbiddenException } from '@nestjs/common'
import { PlanType } from '@prisma/client'

import { BillingService } from '../billing/billing.service'
import { BaseService } from '../common/base.service'
import { PipelinesService } from '../pipelines/pipelines.service'
import { ToolsService } from '../tools/tools.service'
import { WebsocketsService } from '../websockets/websockets.service'
import {
  OrganizationEntity,
  OrganizationModel,
  PlanTypeEnum
} from './entities/organization.entity'
import { OrganizationRepository } from './organization.repository'
import { RoleTypeEnum } from '../members/entities/member.entity'
import { UserEntity } from '../users/entities/user.entity'
import { ConfigService } from '../config/config.service'

@Injectable()
export class OrganizationsService extends BaseService<
  OrganizationEntity,
  OrganizationModel,
  OrganizationRepository
> {
  constructor(
    private billingService: BillingService,
    private configService: ConfigService,
    private organizationRepository: OrganizationRepository,
    private toolsService: ToolsService,
    private pipelinesService: PipelinesService,
    private websocketsService: WebsocketsService
  ) {
    super(organizationRepository)
  }

  async addOrRemoveCredits(orgname: string, numCredits: number) {
    this.logger.debug(`Adding ${numCredits} credits to ${orgname}`)
    const organization = await this.findByOrgname(orgname)
    const model = await this.organizationRepository.update(organization.id, {
      credits:
        numCredits < 0 ? { decrement: -numCredits } : { increment: numCredits },
      billingEmail: organization.billingEmail
    })
    const organizationEntity = this.toEntity(model)
    this.emitMutationEvent(organizationEntity)
    return this.toEntity(organizationEntity)
  }

  async checkCredits(orgname: string, numCredits: number) {
    this.logger.debug(`Checking ${numCredits} credits for ${orgname}`)
    const organization = await this.findByOrgname(orgname)
    if (organization.plan != 'PREMIUM' && organization.credits <= numCredits) {
      throw new ForbiddenException(
        'Sorry, you do not have enough credits. Please purchase more credits to continue' +
          (organization.credits < numCredits
            ? ` (estimated cost: ${numCredits})`
            : '')
      )
    }
  }

  async addUserToOrganization(orgname: string, user: UserEntity) {
    const organization = await this.findByOrgname(orgname)
    const model = await this.organizationRepository.update(organization.id, {
      members: {
        create: {
          inviteAccepted: true,
          inviteEmail: user.email,
          role: RoleTypeEnum.USER,
          username: user.username
        }
      }
    })
    const organizationEntity = this.toEntity(model)
    this.emitMutationEvent(organizationEntity)
    return organizationEntity
  }

  async create(data: Pick<OrganizationEntity, 'billingEmail' | 'orgname'>) {
    // If billing is enabled, create a stripe user, otherwsie set it to orgname
    const billingEnabled = this.configService.get('billing.enabled')
    let stripeCustomerId = data.orgname
    if (billingEnabled) {
      this.logger.debug(
        'BILLING ENABLED - Creating stripe customer for: ' + data.orgname
      )
      const stripeCustomer = await this.billingService.createCustomer(
        data.orgname,
        data.billingEmail
      )
      stripeCustomerId = stripeCustomer.id
    }

    // Create organization and tools
    const organization = await this.organizationRepository.create({
      ...data,
      credits: billingEnabled ? 0 : 100000000,
      stripeCustomerId,
      plan: billingEnabled ? PlanTypeEnum.FREE : PlanTypeEnum.UNLIMITED
    })

    await this.toolsService.createDefaultTools(organization.orgname)
    await this.pipelinesService.createDefaultPipeline(organization.orgname)

    return this.toEntity(organization)
  }

  async findByOrgname(orgname: string) {
    return this.toEntity(
      await this.organizationRepository.findByOrgname(orgname)
    )
  }

  async findByStripeCustomerId(stripeCustomerId: string) {
    return this.toEntity(
      await this.organizationRepository.findByStripeCustomerId(stripeCustomerId)
    )
  }

  async setPlan(orgname: string, plan: PlanType) {
    const organization = await this.findByOrgname(orgname)
    const organizationEntity = this.toEntity(
      await this.organizationRepository.update(organization.id, {
        plan
      })
    )
    this.emitMutationEvent(organizationEntity)
    return organizationEntity
  }

  protected emitMutationEvent(entity: OrganizationEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname]
    })
  }

  protected toEntity(model: OrganizationModel): OrganizationEntity {
    return new OrganizationEntity(model)
  }
}
