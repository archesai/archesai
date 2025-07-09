import { LoggerService } from '#logging/logger.service'

/**
 * Logger class for logging messages with Pino-compatible child logger support.
 */
export class Logger {
  // Additional Pino-compatible methods
  public level = 'info' // Default log level
  private readonly bindings: Record<string, unknown>
  private readonly context: string

  constructor(context: string, bindings: Record<string, unknown> = {}) {
    this.context = context
    this.bindings = bindings
  }

  /**
   * Create a child logger with additional bindings
   */
  public child(bindings: Record<string, unknown>): Logger {
    return new Logger(this.context, {
      ...this.bindings,
      ...bindings
    })
  }

  public debug(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().debug(message, this.getLogObject(obj))
  }

  public error(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().error(message, this.getLogObject(obj))
  }

  public fatal(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().fatal(message, this.getLogObject(obj))
  }

  public info(message: string, obj?: Record<string, unknown>): void {
    this.log(message, obj)
  }

  /**
   * Check if a log level is enabled
   */
  public isLevelEnabled(level: string): boolean {
    const levels = ['trace', 'debug', 'info', 'warn', 'error', 'fatal']
    const currentLevelIndex = levels.indexOf(this.level)
    const checkLevelIndex = levels.indexOf(level)
    return checkLevelIndex >= currentLevelIndex
  }

  public log(message: string, obj?: Record<string, unknown>): void {
    // Handle the special NestApplication case
    if (this.context === 'NestApplication') {
      LoggerService.getInstance().log(message.toLowerCase(), {
        context: obj
      })
      return
    }

    LoggerService.getInstance().log(message, this.getLogObject(obj))
  }

  public silent(): void {
    // No-op for silent mode
  }

  public trace(message: string, obj?: Record<string, unknown>): void {
    this.verbose(message, obj)
  }

  public verbose(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().verbose(message, this.getLogObject(obj))
  }

  public warn(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().warn(message, this.getLogObject(obj))
  }

  private getLogObject(obj?: Record<string, unknown>): Record<string, unknown> {
    return {
      context: this.context,
      ...this.bindings,
      ...obj
    }
  }
}
