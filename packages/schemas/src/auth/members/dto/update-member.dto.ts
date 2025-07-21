import type { z } from 'zod'

import { CreateMemberDtoSchema } from '#auth/members/dto/create-member.dto'

export const UpdateMemberDtoSchema: z.ZodObject<{
  role: z.ZodOptional<
    z.ZodEnum<{
      admin: 'admin'
      member: 'member'
      owner: 'owner'
    }>
  >
}> = CreateMemberDtoSchema.partial()

export type UpdateMemberDto = z.infer<typeof UpdateMemberDtoSchema>
