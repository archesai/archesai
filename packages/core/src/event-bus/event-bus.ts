import EventEmitter from 'node:events'

/**
 * A simple event bus that can be used to emit and listen to events.
 */
export class EventBus {
  private readonly emitter: EventEmitter

  constructor() {
    this.emitter = new EventEmitter()
  }

  public emit(event: string | symbol, ...payload: unknown[]): void {
    this.emitter.emit(event, ...payload)
  }

  // Register a listener
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public on(event: string | symbol, listener: (...args: any[]) => void): this {
    this.emitter.on(event, listener)
    return this
  }

  public removeAllListeners(): void {
    this.emitter.removeAllListeners()
  }
}
