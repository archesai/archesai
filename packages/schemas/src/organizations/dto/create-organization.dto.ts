import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { OrganizationEntitySchema } from '#organizations/entities/organization.entity'

export const CreateOrganizationDtoSchema = Type.Object({
  billingEmail: OrganizationEntitySchema.properties.billingEmail,
  orgname: OrganizationEntitySchema.properties.orgname
})

export type CreateOrganizationDto = Static<typeof CreateOrganizationDtoSchema>
