import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { UserEntitySchema } from '@archesai/schemas'

export const CreateAccountDtoSchema = Type.Object({
  email: UserEntitySchema.properties.email,
  password: Type.String({ description: 'The password for the account' })
})

export type CreateAccountDto = Static<typeof CreateAccountDtoSchema>
