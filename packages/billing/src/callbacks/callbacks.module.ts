import type { ModuleMetadata } from '@archesai/core'

import {
  ConfigModule,
  ConfigService,
  createModule,
  EventBus,
  EventBusModule,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { CallbacksController } from '#callbacks/callbacks.controller'
import { CallbacksService } from '#callbacks/callbacks.service'
import { StripeModule } from '#stripe/stripe.module'
import { StripeService } from '#stripe/stripe.service'

export const CallbacksModuleDefinition: ModuleMetadata = {
  imports: [ConfigModule, EventBusModule, StripeModule, WebsocketsModule],
  providers: [
    {
      inject: [ConfigService, EventBus, StripeService, WebsocketsService],
      provide: CallbacksService,
      useFactory: (
        configService: ConfigService,
        eventBus: EventBus,
        stripeService: StripeService,
        websocketsService: WebsocketsService
      ) =>
        new CallbacksService(
          configService,
          eventBus,
          stripeService,
          websocketsService
        )
    },
    {
      inject: [CallbacksService],
      provide: CallbacksController,
      useFactory: (callbacksService: CallbacksService) =>
        new CallbacksController(callbacksService)
    }
  ]
}

export const CallbacksModule = (() =>
  createModule(class CallbacksModule {}, CallbacksModuleDefinition))()
