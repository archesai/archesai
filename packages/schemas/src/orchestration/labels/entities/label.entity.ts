import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const LabelEntitySchema: TObject<{
  createdAt: TString
  id: TString
  name: TString
  organizationId: TString
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    name: Type.String({
      description: 'The name of the label'
    }),
    organizationId: Type.String({ description: 'The organization name' })
  },
  {
    $id: 'LabelEntity',
    description: 'The label entity',
    title: 'Label Entity'
  }
)

export type LabelEntity = Static<typeof LabelEntitySchema>

export const LABEL_ENTITY_KEY = 'labels'
