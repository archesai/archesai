import { LoggerService } from '#logging/logger.service'

/**
 * Logger class for logging messages.
 */
export class Logger {
  private readonly context: string

  constructor(context: string) {
    this.context = context
  }

  public debug(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().debug(message, {
      context: this.context,
      ...obj
    })
  }

  public error(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().error(message, {
      context: this.context,
      ...obj
    })
  }

  public fatal(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().fatal(message, {
      context: this.context,
      ...obj
    })
  }

  public log(message: string, obj?: Record<string, unknown>): void {
    // this is just a hacky fix to make sure these logs print correctly
    if (this.context === 'NestApplication') {
      LoggerService.getInstance().log(message.toLowerCase(), {
        context: obj
      })
      return
    }
    LoggerService.getInstance().log(message, { context: this.context, ...obj })
  }

  public verbose(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().verbose(message, {
      context: this.context,
      ...obj
    })
  }

  public warn(message: string, obj?: Record<string, unknown>): void {
    LoggerService.getInstance().warn(message, { context: this.context, ...obj })
  }
}
