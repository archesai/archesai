import type { z } from 'zod'

import { CreateInvitationDtoSchema } from '#auth/invitations/dto/create-invitation.dto'

export const UpdateInvitationDtoSchema: z.ZodObject<{
  email: z.ZodOptional<z.ZodString>
  role: z.ZodOptional<
    z.ZodEnum<{
      admin: 'admin'
      member: 'member'
      owner: 'owner'
    }>
  >
}> = CreateInvitationDtoSchema.partial()

export type UpdateInvitationDto = z.infer<typeof UpdateInvitationDtoSchema>
