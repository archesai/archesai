import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { OrganizationEntitySchema } from '@archesai/schemas'

export const CreateOrganizationRequestSchema = Type.Object({
  billingEmail: OrganizationEntitySchema.properties.billingEmail,
  orgname: OrganizationEntitySchema.properties.orgname
})

export type CreateOrganizationRequest = Static<
  typeof CreateOrganizationRequestSchema
>
