import type { Static, TObject, TOptional, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateOrganizationDtoSchema } from '#auth/organizations/dto/create-organization.dto'

export const UpdateOrganizationDtoSchema: TObject<{
  billingEmail: TOptional<TString>
  organizationId: TOptional<TString>
}> = Type.Partial(CreateOrganizationDtoSchema)

export type UpdateOrganizationDto = Static<typeof UpdateOrganizationDtoSchema>
