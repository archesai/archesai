import { Type } from '@sinclair/typebox'

import { UserEntitySchema } from '@archesai/domain'

export const CreateAccountRequestSchema = Type.Object({
  email: UserEntitySchema.properties.email,
  password: Type.String({ description: 'The password for the account' })
})
