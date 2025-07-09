import { Type } from '@sinclair/typebox'

import { UserEntitySchema } from '@archesai/schemas'

export const CreateAccountRequestSchema = Type.Object({
  email: UserEntitySchema.properties.email,
  password: Type.String({ description: 'The password for the account' })
})
