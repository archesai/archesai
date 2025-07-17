import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { ToolEntitySchema } from '#orchestration/tools/entities/tool.entity'

export const CreateToolDtoSchema: TObject<{
  description: TString
  name: TString
}> = Type.Object({
  description: ToolEntitySchema.properties.description,
  name: ToolEntitySchema.properties.name
})

export type CreateToolDto = Static<typeof CreateToolDtoSchema>
