import type { Static, TObject, TOptional, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateRunDtoSchema } from '#orchestration/runs/dto/create-run.dto'

export const UpdateRunDtoSchema: TObject<{
  pipelineId: TOptional<TString>
}> = Type.Partial(CreateRunDtoSchema)

export type UpdateRunDto = Static<typeof UpdateRunDtoSchema>
