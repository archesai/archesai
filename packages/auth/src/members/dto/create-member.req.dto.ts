import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { MemberEntitySchema } from '@archesai/schemas'

export const CreateMemberRequestSchema = Type.Object({
  name: MemberEntitySchema.properties.name,
  role: MemberEntitySchema.properties.role
})

export type CreateMemberRequest = Static<typeof CreateMemberRequestSchema>
