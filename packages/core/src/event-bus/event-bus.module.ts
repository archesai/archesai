import type { EventSubscriber } from '#event-bus/interfaces/event-subscriber.interface'
import type { ModuleMetadata } from '#utils/nest'

import { EventBus } from '#event-bus/event-bus'
import { EventSubscribersLoader } from '#event-bus/event-subscribers.loader'
import { Module } from '#utils/nest'

export const EventBusModuleDefinition: ModuleMetadata = {
  exports: [EventBus],
  providers: [
    {
      provide: EventBus,
      useFactory: () => new EventBus()
    },
    {
      inject: [EventBus],
      provide: EventSubscribersLoader,
      useFactory: (eventBus: EventBus, subscribers: EventSubscriber[]) =>
        new EventSubscribersLoader(eventBus, subscribers)
    }
  ]
}

@Module(EventBusModuleDefinition)
export class EventBusModule {}
