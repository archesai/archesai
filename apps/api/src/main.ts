import type { FastifyInstance } from 'fastify'

import { bootstrap } from '#utils/bootstrap'

export const app: Promise<FastifyInstance> = bootstrap()
