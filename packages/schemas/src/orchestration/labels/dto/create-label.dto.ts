import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LabelEntitySchema } from '#orchestration/labels/entities/label.entity'

export const CreateLabelDtoSchema: TObject<{
  name: TString
}> = Type.Object({
  name: LabelEntitySchema.properties.name
})

export type CreateLabelDto = Static<typeof CreateLabelDtoSchema>
