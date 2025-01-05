import { Controller, Get } from '@nestjs/common'
import { IsPublic } from '../auth/decorators/is-public.decorator'
import { ArchesConfigService } from './config.service'
import { loadConfiguration } from './configuration'
import { ArchesConfig, archesConfigSchema } from './schema'
import { ApiTags } from '@nestjs/swagger'

@ApiTags('Config')
@IsPublic()
@Controller('config')
export class ConfigController {
  constructor(private configService: ArchesConfigService) {}

  @Get()
  async config(): Promise<ArchesConfig> {
    return loadConfiguration(archesConfigSchema)
  }
}
