import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreateAccountDtoSchema: TObject<{
  email: TString
  name: TString
  password: TString
}> = Type.Object({
  email: Type.String({
    description: 'The email address associated with the account'
  }),
  name: Type.String({
    description: 'The name of the user creating the account',
    minLength: 1
  }),
  password: Type.String({ description: 'The password for the account' })
})

export type CreateAccountDto = Static<typeof CreateAccountDtoSchema>
