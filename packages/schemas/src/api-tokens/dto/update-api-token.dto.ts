import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateApiTokenDtoSchema } from '#api-tokens/dto/create-api-token.dto'

export const UpdateApiTokenDtoSchema = Type.Partial(CreateApiTokenDtoSchema)

export type UpdateApiTokenDto = Static<typeof UpdateApiTokenDtoSchema>
