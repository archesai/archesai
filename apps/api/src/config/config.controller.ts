import { Controller, Get } from '@nestjs/common'
import { ConfigService } from './config.service'
import { ArchesConfig } from './schemas/config.schema'
import { ApiTags } from '@nestjs/swagger'

@ApiTags('Config')
@Controller('config')
export class ConfigController {
  constructor(private configService: ConfigService) {}

  @Get()
  async config(): Promise<ArchesConfig> {
    return this.configService.getConfig()
  }
}
