import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { OrganizationEntitySchema } from '#auth/organizations/entities/organization.entity'

export const CreateOrganizationDtoSchema = Type.Object({
  billingEmail: OrganizationEntitySchema.properties.billingEmail,
  organizationId: OrganizationEntitySchema.properties.organizationId
})

export type CreateOrganizationDto = Static<typeof CreateOrganizationDtoSchema>
