import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { OrganizationEntitySchema } from '#auth/organizations/entities/organization.entity'

export const CreateOrganizationDtoSchema: TObject<{
  billingEmail: TString
  organizationId: TString
}> = Type.Object({
  billingEmail: OrganizationEntitySchema.properties.billingEmail,
  organizationId: OrganizationEntitySchema.properties.id
})

export type CreateOrganizationDto = Static<typeof CreateOrganizationDtoSchema>
