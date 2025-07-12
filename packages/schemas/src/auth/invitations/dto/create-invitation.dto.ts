import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { InvitationEntitySchema } from '#auth/invitations/entities/invitation.entity'

export const CreateInvitationDtoSchema = Type.Object({
  email: InvitationEntitySchema.properties.email,
  role: InvitationEntitySchema.properties.role
})

export type CreateInvitationDto = Static<typeof CreateInvitationDtoSchema>
