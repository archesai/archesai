import type { z } from 'zod'

import { MemberEntitySchema } from '#auth/members/entities/member.entity'

export const CreateMemberDtoSchema: z.ZodObject<{
  role: z.ZodEnum<{
    admin: 'admin'
    member: 'member'
    owner: 'owner'
  }>
}> = MemberEntitySchema.pick({
  role: true
})

export type CreateMemberDto = z.infer<typeof CreateMemberDtoSchema>
