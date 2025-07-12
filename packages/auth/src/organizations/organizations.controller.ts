import type { Controller } from '@archesai/core'
import type { OrganizationEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreateOrganizationDtoSchema,
  ORGANIZATION_ENTITY_KEY,
  OrganizationEntitySchema,
  UpdateOrganizationDtoSchema
} from '@archesai/schemas'

import type { OrganizationsService } from '#organizations/organizations.service'

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
      CreateOrganizationDtoSchema,
      UpdateOrganizationDtoSchema,
      organizationsService
    )
  }
}
