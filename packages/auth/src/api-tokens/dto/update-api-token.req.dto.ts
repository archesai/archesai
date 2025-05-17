import { Type } from '@sinclair/typebox'

import { CreateApiTokenRequestSchema } from '#api-tokens/dto/create-api-token.req.dto'

export const UpdateApiTokenRequestSchema = Type.Partial(
  CreateApiTokenRequestSchema
)
