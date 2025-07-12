import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateRunDtoSchema } from '#orchestration/runs/dto/create-run.dto'

export const UpdateRunDtoSchema = Type.Partial(CreateRunDtoSchema)

export type UpdateRunDto = Static<typeof UpdateRunDtoSchema>
