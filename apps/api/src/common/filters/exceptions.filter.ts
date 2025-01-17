import { ConfigService } from '@/src/config/config.service'
import {
  ArgumentsHost,
  BadRequestException,
  Catch,
  ConflictException,
  HttpException,
  HttpStatus,
  NotFoundException
} from '@nestjs/common'
import { Logger } from '@nestjs/common'
import { ExceptionFilter } from '@nestjs/common'
import { Prisma } from '@prisma/client'
import { Response, Request } from 'express'

export interface ExtendedError<T = any> extends Error {
  cause?: T
  stack?: string
}
@Catch()
export class ExceptionsFilter implements ExceptionFilter {
  private readonly logger = new Logger(ExceptionsFilter.name)

  catch(exception: any, host: ArgumentsHost) {
    const ctx = host.switchToHttp()
    const response = ctx.getResponse<Response>()
    const request = ctx.getRequest<Request>()

    // Determine the status code
    let statusCode = this.getStatusCode(exception)

    // Determine the error message
    const message = this.getErrorMessage(exception)

    // Override status code based on specific error types
    if (this.isConflictEror(exception)) {
      statusCode = HttpStatus.CONFLICT
    }

    if (this.isNotFoundError(exception)) {
      statusCode = HttpStatus.NOT_FOUND
    }

    if (this.isBadRequestError(exception)) {
      statusCode = HttpStatus.BAD_REQUEST
    }

    // Handle non-HTTP contexts (e.g., microservices)
    if (host.getType() != 'http') {
      const httpException = new HttpException(exception.name, statusCode)
      return httpException
    }

    // Log the error with additional details
    this.logError(request, exception, statusCode)

    // Get loggineg level from environment
    const configService = new ConfigService()
    const logLevel = configService.get('logging.level')

    // Prepare the error response
    const errorResponse = {
      statusCode,
      message,
      // Include stack and cause only in development
      ...(logLevel === 'trace' && {
        stack: exception?.stack,
        cause: exception?.cause
      })
    }

    response.status(statusCode).json(errorResponse)
  }

  private logError(request: Request, error: ExtendedError, statusCode: number) {
    const logPayload: any = {
      timestamp: new Date().toISOString(),
      path: request.url,
      method: request.method,
      statusCode,
      message: error.message,
      ...(statusCode >= 500 && { stack: error.stack }),
      ...(statusCode >= 500 && error.cause && { cause: error.cause })
    }

    if (statusCode >= 500) {
      this.logger.error(logPayload, `server error`)
    } else if (statusCode >= 400) {
      this.logger.warn(logPayload, `client error`)
    } else {
      this.logger.log(logPayload, `unknown error`)
    }
  }

  private isNotFoundError(exception: any): boolean {
    return (
      exception instanceof NotFoundException ||
      (exception instanceof Prisma.PrismaClientKnownRequestError &&
        exception.message.toLowerCase().includes('not found'))
    )
  }

  private getStatusCode(exception: any): number {
    return exception instanceof HttpException
      ? exception.getStatus()
      : HttpStatus.INTERNAL_SERVER_ERROR
  }

  private getErrorMessage(exception: any): string {
    return exception.message || 'Internal server error'
  }

  private isConflictEror(exception: any): boolean {
    return (
      exception instanceof ConflictException ||
      (exception instanceof Prisma.PrismaClientKnownRequestError &&
        exception.code === 'P2002')
    )
  }

  private isBadRequestError(exception: any): boolean {
    return (
      exception instanceof BadRequestException ||
      exception instanceof Prisma.PrismaClientValidationError
    )
  }
}
