import type { ModuleMetadata } from '@archesai/core'

import {
  createModule,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { CustomersModule } from '#customers/customers.module'
import { CustomersService } from '#customers/customers.service'
import { PaymentMethodRepository } from '#payment-methods/payment-method.repository'
import { PaymentMethodsController } from '#payment-methods/payment-methods.controller'
import { PaymentMethodsService } from '#payment-methods/payment-methods.service'
import { StripeModule } from '#stripe/stripe.module'
import { StripeService } from '#stripe/stripe.service'

export const PaymentMethodsModuleDefinition: ModuleMetadata = {
  exports: [PaymentMethodsService],
  imports: [CustomersModule, WebsocketsModule, StripeModule],
  providers: [
    {
      inject: [CustomersService, PaymentMethodRepository, WebsocketsService],
      provide: PaymentMethodsService,
      useFactory: (
        customersService: CustomersService,
        paymentMethodRepository: PaymentMethodRepository,
        websocketsService: WebsocketsService
      ) =>
        new PaymentMethodsService(
          customersService,
          paymentMethodRepository,
          websocketsService
        )
    },
    {
      inject: [StripeService],
      provide: PaymentMethodRepository,
      useFactory: (stripeService: StripeService) =>
        new PaymentMethodRepository(stripeService)
    },
    {
      inject: [PaymentMethodsService],
      provide: PaymentMethodsController,
      useFactory: (paymentMethodsService: PaymentMethodsService) =>
        new PaymentMethodsController(paymentMethodsService)
    }
  ]
}

export const PaymentMethodsModule = (() =>
  createModule(class PaymentMethodsModule {}, PaymentMethodsModuleDefinition))()
