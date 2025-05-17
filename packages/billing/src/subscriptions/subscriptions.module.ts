import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'

import { StripeModule } from '#stripe/stripe.module'
import { StripeService } from '#stripe/stripe.service'
import { SubscriptionsController } from '#subscriptions/subscriptions.controller'
import { SubscriptionsService } from '#subscriptions/subscriptions.service'

export const SubscriptionModuleDefinition: ModuleMetadata = {
  exports: [SubscriptionsService],
  imports: [StripeModule],
  providers: [
    {
      inject: [StripeService],
      provide: SubscriptionsService,
      useFactory: (stripeService: StripeService) =>
        new SubscriptionsService(stripeService)
    },
    {
      inject: [SubscriptionsService],
      provide: SubscriptionsController,
      useFactory: (subscriptionsService: SubscriptionsService) =>
        new SubscriptionsController(subscriptionsService)
    }
  ]
}

@Module(SubscriptionModuleDefinition)
export class SubscriptionModule {}
