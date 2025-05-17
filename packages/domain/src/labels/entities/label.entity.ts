import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'

export const LabelEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    orgname: Type.String({ description: 'The organization name' })
  },
  {
    $id: 'LabelEntity',
    description: 'The label entity',
    title: 'Label Entity'
  }
)

export class LabelEntity
  extends BaseEntity
  implements Static<typeof LabelEntitySchema>
{
  public orgname: string
  public type = LABEL_ENTITY_KEY

  constructor(props: LabelEntity) {
    super(props)
    this.orgname = props.orgname
  }
}

export const LABEL_ENTITY_KEY = 'labels'
