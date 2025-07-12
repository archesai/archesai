import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateToolDtoSchema } from '#tools/dto/create-tool.dto'

export const UpdateToolDtoSchema = Type.Partial(CreateToolDtoSchema)

export type UpdateToolDto = Static<typeof UpdateToolDtoSchema>
