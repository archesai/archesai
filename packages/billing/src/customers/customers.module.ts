import type { ModuleMetadata } from '@archesai/core'

import { createModule, EventBus, EventBusModule } from '@archesai/core'

import { CustomersService } from '#customers/customers.service'
import { CustomersSubscriber } from '#customers/customers.subscriber'
import { StripeModule } from '#stripe/stripe.module'
import { StripeService } from '#stripe/stripe.service'

export const CustomersModuleDefinition: ModuleMetadata = {
  exports: [CustomersService],
  imports: [EventBusModule, StripeModule],
  providers: [
    {
      inject: [StripeService],
      provide: CustomersService,
      useFactory: (stripeService: StripeService) =>
        new CustomersService(stripeService)
    },
    {
      inject: [CustomersService, EventBus],
      provide: CustomersSubscriber,
      useFactory: (customersService: CustomersService, eventBus: EventBus) =>
        new CustomersSubscriber(customersService, eventBus)
    }
  ]
}

export const CustomersModule = (() =>
  createModule(class CustomersModule {}, CustomersModuleDefinition))()
