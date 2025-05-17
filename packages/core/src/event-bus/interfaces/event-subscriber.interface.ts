import type { EventBus } from '#event-bus/event-bus'

export interface EventSubscriber {
  subscribe(eventBus: EventBus): void
}
