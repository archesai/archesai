import type { Logger as PinoLogger } from 'pino'

/**
 * Adapter for the Pino logger.
 */
export class PinoLoggerAdapter {
  private readonly pinoLogger: PinoLogger

  constructor(pinoLogger: PinoLogger) {
    this.pinoLogger = pinoLogger
  }

  public debug(message: string, meta?: Record<string, unknown>): void {
    this.pinoLogger.debug(meta ?? {}, message)
  }

  public error(message: string, meta?: Record<string, unknown>): void {
    this.pinoLogger.error(meta ?? {}, message)
  }

  public fatal(message: string, meta?: Record<string, unknown>): void {
    this.pinoLogger.fatal(meta ?? {}, message)
  }

  public log(message: string, meta?: Record<string, unknown>): void {
    this.pinoLogger.info(meta ?? {}, message)
  }

  public verbose(message: string, meta?: Record<string, unknown>): void {
    this.pinoLogger.trace(meta ?? {}, message)
  }

  public warn(message: string, meta?: Record<string, unknown>): void {
    this.pinoLogger.warn(meta ?? {}, message)
  }
}
