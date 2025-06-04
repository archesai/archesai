import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'

import type { Errors } from '#http/schemas/errors.schema'
import type { ArgumentsHost } from '#types/context'
import type { ArchesApiResponse } from '#utils/get-req.transformer'

import { Logger } from '#logging/logger'

/**
 * Handles exceptions and returns a proper response.
 */
export class ExceptionsFilter {
  private readonly logger = new Logger(ExceptionsFilter.name)

  public catch(exception: Error, host: ArgumentsHost) {
    if (host.getType() != 'http') {
      const httpException = new Error(exception.message)
      return httpException
    }

    const ctx = host.switchToHttp()
    const response = ctx.getResponse<ArchesApiResponse>()

    let statusCode = 500
    let message = this.getErrorMessage(exception)

    if (
      Value.Check(
        Type.Object({
          response: Type.Object({
            error: Type.String(),
            message: Type.String(),
            statusCode: Type.Number()
          }),
          status: Type.Number()
        }),
        exception
      )
    ) {
      message = exception.response.message
      statusCode = exception.response.statusCode
    }

    if (this.isConflictEror(exception)) {
      statusCode = 409
    }

    if (this.isNotFoundError(exception)) {
      statusCode = 404
    }

    if (this.isBadRequestError(exception)) {
      statusCode = 400
    }

    if (this.isUnauthorizedError(exception)) {
      statusCode = 401
    }

    const errorResponse = {
      errors: [
        {
          detail: message,
          status: statusCode.toString(),
          title: exception.name
        }
      ]
    } satisfies { errors: Errors }

    if (statusCode >= 500) {
      this.logger.error(`server error`, {
        cause: exception.cause,
        errorResponse,
        stack: exception.stack
      })
    } else if (statusCode >= 400) {
      this.logger.warn(`client error`, { errorResponse })
    } else {
      this.logger.log(`unknown error`, { errorResponse })
    }

    return response.code(statusCode).send(errorResponse)
  }

  private getErrorMessage(exception: unknown): string {
    if (exception instanceof Error) {
      return exception.message
    } else if (typeof exception === 'string') {
      return exception
    } else {
      return JSON.stringify(exception)
    }
  }

  private isBadRequestError(exception: unknown): boolean {
    return this.isErrorOfType(exception, 'BadRequest')
  }

  private isConflictEror(exception: unknown): boolean {
    return this.isErrorOfType(exception, 'Conflict')
  }

  private isErrorOfType(exception: unknown, type: string): boolean {
    return (
      exception instanceof Error &&
      (exception.message.includes(type) || exception.name.includes(type))
    )
  }

  private isNotFoundError(exception: unknown): boolean {
    return this.isErrorOfType(exception, 'NotFound')
  }

  private isUnauthorizedError(exception: unknown): boolean {
    return this.isErrorOfType(exception, 'Unauthorized')
  }
}
