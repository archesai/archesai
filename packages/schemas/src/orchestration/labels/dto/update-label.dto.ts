import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateLabelDtoSchema } from '#orchestration/labels/dto/create-label.dto'

export const UpdateLabelDtoSchema = Type.Partial(CreateLabelDtoSchema)

export type UpdateLabelDto = Static<typeof UpdateLabelDtoSchema>
