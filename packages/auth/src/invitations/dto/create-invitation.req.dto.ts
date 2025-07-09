import { Type } from '@sinclair/typebox'

import { InvitationEntitySchema } from '@archesai/schemas'

export const CreateInvitationRequestSchema = Type.Object({
  email: InvitationEntitySchema.properties.email,
  name: InvitationEntitySchema.properties.name,
  role: InvitationEntitySchema.properties.role
})
