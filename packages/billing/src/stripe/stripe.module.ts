import type Stripe from 'stripe'

import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, createModule } from '@archesai/core'

import { StripeService } from '#stripe/stripe.service'

export const StripeModuleDefinition: ModuleMetadata = {
  exports: [StripeService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: StripeService,
      useFactory: (configService: ConfigService) => {
        if (!configService.get('billing.enabled')) {
          // Return a dummy or no-op service
          return {
            configService,
            createCheckoutSession: () => {
              throw new Error('Billing feature is disabled.')
            },
            createPortal: () => {
              throw new Error('Billing feature is disabled.')
            },
            createSetupIntent: () => {
              throw new Error('Billing feature is disabled.')
            },
            getPrice() {
              throw new Error('Billing feature is disabled.')
            },
            getProduct() {
              throw new Error('Billing feature is disabled.')
            },
            stripe: null as unknown as Stripe
          } satisfies StripeService
        } else {
          return new StripeService(configService)
        }
      }
    }
  ]
}

export const StripeModule = (() =>
  createModule(class StripeModule {}, StripeModuleDefinition))()
