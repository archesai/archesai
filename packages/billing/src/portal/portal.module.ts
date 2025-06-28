import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, createModule } from '@archesai/core'

import { PortalController } from '#portal/portal.controller'
import { PortalService } from '#portal/portal.service'
import { StripeModule } from '#stripe/stripe.module'
import { StripeService } from '#stripe/stripe.service'

export const PortalModuleDefinition: ModuleMetadata = {
  exports: [PortalService],
  imports: [ConfigModule, StripeModule],
  providers: [
    {
      inject: [StripeService, ConfigService],
      provide: PortalService,
      useFactory: (
        stripeService: StripeService,
        configService: ConfigService
      ) => new PortalService(stripeService, configService)
    },
    {
      inject: [PortalService],
      provide: PortalController,
      useFactory: (portalService: PortalService) =>
        new PortalController(portalService)
    }
  ]
}

export const PortalModule = (() =>
  createModule(class PortalModule {}, PortalModuleDefinition))()
