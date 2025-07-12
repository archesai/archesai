import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateAccountDtoSchema } from '#accounts/dto/create-account.dto'

export const UpdateAccountDtoSchema = Type.Partial(CreateAccountDtoSchema)

export type UpdateAccountDtoDto = Static<typeof UpdateAccountDtoSchema>
