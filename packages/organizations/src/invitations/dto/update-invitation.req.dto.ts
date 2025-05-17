import { Type } from '@sinclair/typebox'

import { CreateInvitationRequestSchema } from '#invitations/dto/create-invitation.req.dto'

export const UpdateInvitationRequestSchema = Type.Partial(
  CreateInvitationRequestSchema
)
