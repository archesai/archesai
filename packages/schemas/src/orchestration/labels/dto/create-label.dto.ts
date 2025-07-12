import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LabelEntitySchema } from '#orchestration/labels/entities/label.entity'

export const CreateLabelDtoSchema = Type.Object({
  name: LabelEntitySchema.properties.name
})

export type CreateLabelDto = Static<typeof CreateLabelDtoSchema>
