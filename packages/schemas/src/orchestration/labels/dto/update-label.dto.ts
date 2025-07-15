import type { Static, TObject, TOptional, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateLabelDtoSchema } from '#orchestration/labels/dto/create-label.dto'

export const UpdateLabelDtoSchema: TObject<{
  name: TOptional<TString>
}> = Type.Partial(CreateLabelDtoSchema)

export type UpdateLabelDto = Static<typeof UpdateLabelDtoSchema>
