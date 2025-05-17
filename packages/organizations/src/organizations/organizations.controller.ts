import type { ArchesApiRequest, Controller } from '@archesai/core'
import type { OrganizationEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import {
  ORGANIZATION_ENTITY_KEY,
  OrganizationEntitySchema
} from '@archesai/domain'

import type { CreateOrganizationRequest } from '#organizations/dto/create-organization.req.dto'
import type { OrganizationsService } from '#organizations/organizations.service'

import { CreateOrganizationRequestSchema } from '#organizations/dto/create-organization.req.dto'
import { UpdateOrganizationRequestSchema } from '#organizations/dto/update-organization.req.dto'

/**
 * Controller for handling organizations.
 */
export class OrganizationsController
  extends BaseController<OrganizationEntity>
  implements Controller
{
  private readonly organizationsService: OrganizationsService
  constructor(organizationsService: OrganizationsService) {
    super(
      ORGANIZATION_ENTITY_KEY,
      OrganizationEntitySchema,
      CreateOrganizationRequestSchema,
      UpdateOrganizationRequestSchema,
      organizationsService
    )
    this.organizationsService = organizationsService
  }

  public override async create(
    request: ArchesApiRequest & { body: CreateOrganizationRequest }
  ) {
    return this.toIndividualResponse(
      await this.organizationsService.create({
        billingEmail: request.body.billingEmail,
        creatorId: request.user?.id ?? '',
        credits: 0,
        orgname: request.body.orgname,
        plan: 'FREE'
      })
    )
  }
}
