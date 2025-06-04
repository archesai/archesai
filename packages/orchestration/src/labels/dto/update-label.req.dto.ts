import { Type } from '@sinclair/typebox'

import { CreateLabelRequestSchema } from '#labels/dto/create-label.req.dto'

export const UpdateLabelRequestSchema = Type.Partial(CreateLabelRequestSchema)
