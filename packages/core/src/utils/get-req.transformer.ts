import type {
  FastifyReply,
  FastifyRequest,
  RouteGenericInterface
} from 'fastify'

import type { UserEntity } from '@archesai/schemas'

export type ArchesApiRequest<
  T extends RouteGenericInterface = RouteGenericInterface
> = FastifyRequest<T> & {
  user?: UserEntity
}

export type ArchesApiResponse = FastifyReply
