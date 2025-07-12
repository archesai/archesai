import type { TypeBoxTypeProvider } from '@fastify/type-provider-typebox'
import type {
  FastifyBaseLogger,
  FastifyInstance,
  RawReplyDefaultExpression,
  RawRequestDefaultExpression,
  RawServerDefault
} from 'fastify'

import type { UserEntity } from '@archesai/schemas'

export type HttpInstance = FastifyInstance<
  RawServerDefault,
  RawRequestDefaultExpression,
  RawReplyDefaultExpression,
  FastifyBaseLogger,
  TypeBoxTypeProvider
>

declare module 'fastify' {
  type PassportUser = UserEntity
}
