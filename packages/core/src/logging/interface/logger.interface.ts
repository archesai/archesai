/**
 * Logger interface
 */
export interface ILogger {
  debug(message: string, meta?: Record<string, unknown>): void
  error(message: string, meta?: Record<string, unknown>): void
  fatal(message: string, meta?: Record<string, unknown>): void
  log(message: string, meta?: Record<string, unknown>): void
  setLogLevels?(levels: LogLevel[]): void
  verbose(message: string, meta?: Record<string, unknown>): void
  warn(message: string, meta?: Record<string, unknown>): void
}

export type LogLevel = 'debug' | 'error' | 'fatal' | 'log' | 'verbose' | 'warn'
