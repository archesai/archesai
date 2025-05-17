import { Type } from '@sinclair/typebox'

import { CreateRunRequestSchema } from '#runs/dto/create-run.req.dto'

export const UpdateRunRequestSchema = Type.Partial(CreateRunRequestSchema)
