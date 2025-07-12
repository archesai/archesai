import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

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

export type LabelEntity = Static<typeof LabelEntitySchema>

export const LABEL_ENTITY_KEY = 'labels'
