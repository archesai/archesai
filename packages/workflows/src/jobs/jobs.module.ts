import { readFileSync } from 'node:fs'

import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, Module } from '@archesai/core'

import { createQueue } from '#jobs/factories/queue.factory'

const QUEUE_SYMBOL = Symbol('QUEUE')

export const JobsModuleDefinition: ModuleMetadata = {
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: QUEUE_SYMBOL,
      useFactory: (configService: ConfigService) => {
        const redisCa = configService.get('redis.ca')
        return createQueue('my-queue', {
          host: configService.get('redis.host'),
          password: configService.get('redis.auth')!,
          port: configService.get('redis.port'),
          ...(redisCa ? { tls: { ca: readFileSync(redisCa) } } : {})
        })
      }
    }
  ]
}

@Module(JobsModuleDefinition)
export class JobsModule {}
