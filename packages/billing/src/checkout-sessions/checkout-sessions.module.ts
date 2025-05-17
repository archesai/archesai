import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, Module } from '@archesai/core'

import { CheckoutSessionsController } from '#checkout-sessions/checkout-sessions.controller'
import { CheckoutSessionsService } from '#checkout-sessions/checkout-sessions.service'
import { StripeModule } from '#stripe/stripe.module'
import { StripeService } from '#stripe/stripe.service'

export const CheckoutSessionsModuleDefinition: ModuleMetadata = {
  imports: [ConfigModule, StripeModule],
  providers: [
    {
      inject: [ConfigService, StripeService],
      provide: CheckoutSessionsService,
      useFactory: (
        configService: ConfigService,
        stripeService: StripeService
      ) => new CheckoutSessionsService(configService, stripeService)
    },
    {
      inject: [CheckoutSessionsService],
      provide: CheckoutSessionsController,
      useFactory: (checkoutSessionsService: CheckoutSessionsService) =>
        new CheckoutSessionsController(checkoutSessionsService)
    }
  ]
}

@Module(CheckoutSessionsModuleDefinition)
export class CheckoutSessionsModule {}
