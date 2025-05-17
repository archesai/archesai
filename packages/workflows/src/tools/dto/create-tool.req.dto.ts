import { Type } from '@sinclair/typebox'

import { ToolEntitySchema } from '@archesai/domain'

export const CreateToolRequestSchema = Type.Object({
  description: ToolEntitySchema.properties.description,
  name: ToolEntitySchema.properties.name
})
