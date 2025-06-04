import { Type } from '@sinclair/typebox'

import { CreateToolRequestSchema } from '#tools/dto/create-tool.req.dto'

export const UpdateToolRequestSchema = Type.Partial(CreateToolRequestSchema)
