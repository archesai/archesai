import { Type } from '@sinclair/typebox'

import { CreateContentRequestSchema } from '#content/dto/create-content.req.dto'

export const UpdateContentRequestSchema = Type.Partial(
  CreateContentRequestSchema
)
