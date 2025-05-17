import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const AccessTokenDecodedJwtSchema = Type.Object({
  sub: Type.String({ description: 'The subject of the token' })
})

export type AccessTokenDecodedJwt = Static<typeof AccessTokenDecodedJwtSchema>
