import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { EmailService } from '#email/email.service'
import { createModule } from '#utils/nest'

export const EmailModuleDefinition: ModuleMetadata = {
  exports: [EmailService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: EmailService,
      useFactory: (configService: ConfigService) =>
        new EmailService(configService)
    }
  ]
}

export const EmailModule = (() =>
  createModule(class EmailModule {}, EmailModuleDefinition))()
