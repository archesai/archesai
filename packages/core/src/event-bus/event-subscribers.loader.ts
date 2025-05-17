import type { EventBus } from '#event-bus/event-bus'
import type { EventSubscriber } from '#event-bus/interfaces/event-subscriber.interface'

/**
 * Loads event subscribers into the event bus.
 */
export class EventSubscribersLoader {
  private readonly eventBus: EventBus
  private readonly subscribers: EventSubscriber[]

  constructor(eventBus: EventBus, subscribers: EventSubscriber[]) {
    this.eventBus = eventBus
    this.subscribers = subscribers
  }

  public loadEventListeners() {
    this.subscribers.forEach((subscriber) => {
      subscriber.subscribe(this.eventBus)
    })
  }

  public removeAllListeners() {
    this.eventBus.removeAllListeners()
  }
}
