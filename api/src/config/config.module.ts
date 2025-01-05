import { Global, Module } from '@nestjs/common'
import { ArchesConfigService } from './config.service'
import { loadConfiguration } from './configuration'
import { ConfigModule } from '@nestjs/config'
import { archesConfigSchema } from './schema'
import { ConfigController } from './config.controller'

@Global()
@Module({
  imports: [
    ConfigModule.forRoot({
      ignoreEnvFile: true,
      ignoreEnvVars: true,
      load: [() => loadConfiguration(archesConfigSchema)]
    })
  ],
  providers: [ArchesConfigService],
  exports: [ArchesConfigService],
  controllers: [ConfigController]
})
export class ArchesConfigModule {}
