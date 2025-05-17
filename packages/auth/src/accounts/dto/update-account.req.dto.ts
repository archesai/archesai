import { Type } from '@sinclair/typebox'

import { CreateAccountRequestSchema } from '#accounts/dto/create-account.req.dto'

export const UpdateAccountRequestSchema = Type.Partial(
  CreateAccountRequestSchema
)
