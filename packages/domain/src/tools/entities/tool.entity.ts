import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { ContentBaseType } from '#enums/role'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'
import { ContentBaseTypes } from '#enums/role'

export const ToolEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    description: Type.String({ description: 'The tool description' }),
    inputType: Type.Union(
      ContentBaseTypes.map((type) => Type.Literal(type)), // Using literals instead of enums
      { description: 'The input type of the tool' }
    ),
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

export class ToolEntity
  extends BaseEntity
  implements Static<typeof ToolEntitySchema>
{
  public description: string
  public inputType: ContentBaseType
  public orgname: string
  public outputType: ContentBaseType
  public toolBase: string
  public type = TOOL_ENTITY_KEY

  constructor(props: ToolEntity) {
    super(props)
    this.description = props.description
    this.inputType = props.inputType
    this.orgname = props.orgname
    this.outputType = props.outputType
    this.toolBase = props.toolBase
  }
}

export const TOOL_ENTITY_KEY = 'tools'
