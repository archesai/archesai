import type { z } from 'zod'

import { InvitationEntitySchema } from '#auth/invitations/entities/invitation.entity'

export const CreateInvitationDtoSchema: z.ZodObject<{
  email: z.ZodString
  role: z.ZodEnum<{
    admin: 'admin'
    member: 'member'
    owner: 'owner'
  }>
}> = InvitationEntitySchema.pick({
  email: true,
  role: true
})

export type CreateInvitationDto = z.infer<typeof CreateInvitationDtoSchema>
