import type { Static, TObject, TOptional, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateToolDtoSchema } from '#orchestration/tools/dto/create-tool.dto'

export const UpdateToolDtoSchema: TObject<{
  description: TOptional<TString>
  name: TOptional<TString>
}> = Type.Partial(CreateToolDtoSchema)

export type UpdateToolDto = Static<typeof UpdateToolDtoSchema>
