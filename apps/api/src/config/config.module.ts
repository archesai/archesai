import { Global, Module } from '@nestjs/common'
import { ConfigService } from './config.service'
import { ConfigController } from './config.controller'

@Global()
@Module({
  providers: [ConfigService],
  exports: [ConfigService],
  controllers: [ConfigController]
})
export class ConfigModule {}
