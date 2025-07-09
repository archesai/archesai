import { Type } from '@sinclair/typebox'

import { ToolEntitySchema } from '@archesai/schemas'

export const CreateToolRequestSchema = Type.Object({
  description: ToolEntitySchema.properties.description,
  name: ToolEntitySchema.properties.name
})
