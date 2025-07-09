import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const ApiTokenDecodedJwtSchema = Type.Object({
  sub: Type.String({ description: 'The subject of the token' })
})

export type ApiTokenDecodedJwt = Static<typeof ApiTokenDecodedJwtSchema>
