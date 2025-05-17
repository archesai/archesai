import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'

import { PlansController } from '#plans/plans.controller'
import { PlansService } from '#plans/plans.service'
import { StripeModule } from '#stripe/stripe.module'
import { StripeService } from '#stripe/stripe.service'

export const PlansModuleDefinition: ModuleMetadata = {
  exports: [PlansService],
  imports: [StripeModule],
  providers: [
    {
      inject: [StripeService],
      provide: PlansService,
      useFactory: (stripeService: StripeService) =>
        new PlansService(stripeService)
    },
    {
      inject: [PlansService],
      provide: PlansController,
      useFactory: (plansService: PlansService) =>
        new PlansController(plansService)
    }
  ]
}

@Module(PlansModuleDefinition)
export class PlansModule {}
