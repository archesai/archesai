import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreateAccountDtoSchema = Type.Object({
  email: Type.String({
    description: 'The email address associated with the account'
  }),
  name: Type.String({
    description: 'The name of the user creating the account'
  }),
  password: Type.String({ description: 'The password for the account' })
})

export type CreateAccountDto = Static<typeof CreateAccountDtoSchema>
