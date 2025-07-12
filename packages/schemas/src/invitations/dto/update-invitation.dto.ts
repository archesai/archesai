import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateInvitationDtoSchema } from '#invitations/dto/create-invitation.dto'

export const UpdateInvitationDtoSchema = Type.Partial(CreateInvitationDtoSchema)

export type UpdateInvitationDto = Static<typeof UpdateInvitationDtoSchema>
