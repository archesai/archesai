import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { InvitationEntitySchema } from '#invitations/entities/invitation.entity'

export const CreateInvitationDtoSchema = Type.Object({
  email: InvitationEntitySchema.properties.email,
  name: InvitationEntitySchema.properties.name,
  role: InvitationEntitySchema.properties.role
})

export type CreateInvitationDto = Static<typeof CreateInvitationDtoSchema>
