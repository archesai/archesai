import { PartialType, PickType } from '@nestjs/swagger'

import { CreateOrganizationDto } from './create-organization.dto'

export class UpdateOrganizationDto extends PartialType(
  PickType(CreateOrganizationDto, ['billingEmail'] as const)
) {}
