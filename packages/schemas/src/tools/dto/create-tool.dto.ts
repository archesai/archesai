import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { ToolEntitySchema } from '#tools/entities/tool.entity'

export const CreateToolDtoSchema = Type.Object({
  description: ToolEntitySchema.properties.description,
  name: ToolEntitySchema.properties.name
})

export type CreateToolDto = Static<typeof CreateToolDtoSchema>
