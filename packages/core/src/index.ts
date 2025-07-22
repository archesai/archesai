export type { BaseService } from '#common/base-service'
export { createBaseService } from '#common/base-service'
export { crudPlugin } from '#common/crud.plugin'
export { configController } from '#config/config.controller'
export type { ConfigService } from '#config/config.service'
export { createConfigService } from '#config/config.service'
export type { EmailService } from '#email/email.service'
export { createEmailService } from '#email/email.service'
export {
  getEmailChangeConfirmationHtml,
  getEmailVerificationHtml,
  getPasswordResetHtml
} from '#email/templates'
export { errorHandlerPlugin } from '#exceptions/exception-handler'
export {
  AppError,
  BadRequestException,
  ConflictException,
  ForbiddenException,
  InternalServerErrorException,
  NotFoundException,
  UnauthorizedException,
  ValidationException
} from '#exceptions/http-errors'
export { PinoLoggerAdapter } from '#logging/adapters/pino-logger-adapter'
export type { Logger, LogLevel } from '#logging/logger'
export type { LoggerService } from '#logging/logger.service'
export { createLogger } from '#logging/logger.service'
export { createRedisService, RedisService } from '#redis/redis.service'
export { capitalize } from '#utils/capitalize'
export { catchError, catchErrorAsync } from '#utils/catch-error'
export { delay } from '#utils/delay'
export { generateExports, getAllFiles } from '#utils/generate-index'
export {
  isEmpty,
  isNil,
  isObject,
  isString,
  isUndefined
} from '#utils/is-primitive'
export { pluralize, singularize } from '#utils/pluralize'
export { retry } from '#utils/retry'
export { streamToBuffer } from '#utils/stream-to-buffer'
export {
  toCamelCase,
  toKebabCase,
  toSentenceCase,
  toSnakeCase,
  toTitleCase,
  toTitleCaseNoSpaces,
  vf
} from '#utils/strings'
export { toError } from '#utils/to-error'
export { RedisIoAdapter } from '#websockets/adapters/redis-io.adapter'
export {
  createWebsocketsService,
  WebsocketsService
} from '#websockets/websockets.service'
