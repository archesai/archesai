import { z } from 'zod'

import { ApiConfigSchema } from '#config/api.schema'
import { AuthConfigSchema } from '#config/auth.schema'
import { BillingConfigSchema } from '#config/billing.schema'
import { DatabaseConfigSchema } from '#config/database.schema'
import { InfrastructureConfigSchema } from '#config/infrastructure.schema'
import { IngressConfigSchema } from '#config/ingress.schema'
import { IntelligenceConfigSchema } from '#config/intelligence.schema'
import { LoggingConfigSchema } from '#config/logging.schema'
import { MonitoringConfigSchema } from '#config/monitoring.schema'
import { PlatformConfigSchema } from '#config/platform.schema'
import { RedisConfigSchema } from '#config/redis.schema'
import { StorageConfigSchema } from '#config/storage.schema'

export const ArchesConfigSchema: z.ZodObject<{
  api: typeof ApiConfigSchema
  auth: typeof AuthConfigSchema
  billing: typeof BillingConfigSchema
  database: typeof DatabaseConfigSchema
  infrastructure: typeof InfrastructureConfigSchema
  ingress: typeof IngressConfigSchema
  intelligence: typeof IntelligenceConfigSchema
  logging: typeof LoggingConfigSchema
  monitoring: typeof MonitoringConfigSchema
  platform: typeof PlatformConfigSchema
  redis: typeof RedisConfigSchema
  storage: typeof StorageConfigSchema
}> = z
  .object({
    api: ApiConfigSchema,
    auth: AuthConfigSchema,
    billing: BillingConfigSchema,
    database: DatabaseConfigSchema,
    infrastructure: InfrastructureConfigSchema,
    ingress: IngressConfigSchema,
    intelligence: IntelligenceConfigSchema,
    logging: LoggingConfigSchema,
    monitoring: MonitoringConfigSchema,
    platform: PlatformConfigSchema,
    redis: RedisConfigSchema,
    storage: StorageConfigSchema
  })
  .describe('Arches AI configuration schema')

export type ArchesConfig = z.infer<typeof ArchesConfigSchema>
