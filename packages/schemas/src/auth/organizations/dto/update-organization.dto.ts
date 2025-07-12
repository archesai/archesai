import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateOrganizationDtoSchema } from '#auth/organizations/dto/create-organization.dto'

export const UpdateOrganizationDtoSchema = Type.Partial(
  CreateOrganizationDtoSchema
)

export type UpdateOrganizationDto = Static<typeof UpdateOrganizationDtoSchema>
