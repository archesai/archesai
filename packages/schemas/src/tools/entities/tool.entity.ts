import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { ContentBaseTypes } from '#enums/role'

export const ToolEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    description: Type.String({ description: 'The tool description' }),
    inputType: Type.Union(
      ContentBaseTypes.map((type) => Type.Literal(type)), // Using literals instead of enums
      { description: 'The input type of the tool' }
    ),
    name: Type.String({
      description: 'The name of the tool'
    }),
    orgname: Type.String({ description: 'The organization name' }),
    outputType: Type.Union(
      ContentBaseTypes.map((type) => Type.Literal(type)), // Using literals instead of enums
      { description: 'The output type of the tool' }
    ),
    toolBase: Type.String({ description: 'The base of the tool' })
  },
  {
    $id: 'ToolEntity',
    description: 'The tool entity',
    title: 'Tool Entity'
  }
)

export type ToolEntity = Static<typeof ToolEntitySchema>

export const TOOL_ENTITY_KEY = 'tools'
