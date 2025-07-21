import type { FastifyPluginAsync } from 'fastify'
import type { ZodTypeProvider } from 'fastify-type-provider-zod'

import type { EmailService, RedisService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { HealthCheckSchema } from '@archesai/schemas'

export const healthPlugin: FastifyPluginAsync<{
  databaseService: DrizzleDatabaseService
  emailService: EmailService
  redisService: RedisService
}> = async (app, { databaseService, emailService, redisService }) => {
  app.withTypeProvider<ZodTypeProvider>().get(
    '/health',
    {
      schema: {
        response: {
          200: HealthCheckSchema
        },
        summary: 'Health check endpoint',
        tags: ['System']
      }
    },
    async () => {
      // Use functional services for health checks
      const dbStatus = await checkServiceHealth('database', () =>
        databaseService.ping()
      )
      const redisStatus = await checkServiceHealth('redis', () =>
        redisService.ping()
      )
      const emailStatus = await checkServiceHealth('email', () =>
        emailService.ping()
      )

      return {
        services: {
          database: dbStatus,
          email: emailStatus,
          redis: redisStatus
        },
        timestamp: new Date().toISOString(),
        uptime: process.uptime()
      }
    }
  )

  await Promise.resolve()
}

async function checkServiceHealth(
  serviceName: string,
  healthCheck: () => Promise<boolean>
): Promise<string> {
  try {
    const isHealthy = await healthCheck()
    return isHealthy ? 'healthy' : 'unhealthy'
  } catch (error) {
    console.error(`Health check failed for ${serviceName}:`, error)
    return 'unhealthy'
  }
}
