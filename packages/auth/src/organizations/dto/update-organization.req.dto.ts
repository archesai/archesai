import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateOrganizationRequestSchema } from '#organizations/dto/create-organization.req.dto'

export const UpdateOrganizationRequestSchema = Type.Partial(
  CreateOrganizationRequestSchema
)

export type UpdateOrganizationRequest = Static<
  typeof UpdateOrganizationRequestSchema
>
