import type { Controller } from '@archesai/core'
import type { OrganizationEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  ORGANIZATION_ENTITY_KEY,
  OrganizationEntitySchema
} from '@archesai/schemas'

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
  constructor(organizationsService: OrganizationsService) {
    super(
      ORGANIZATION_ENTITY_KEY,
      OrganizationEntitySchema,
      CreateOrganizationRequestSchema,
      UpdateOrganizationRequestSchema,
      organizationsService
    )
  }
}
