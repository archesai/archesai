import type { FastifyPluginAsync, FastifyRequest } from 'fastify'

import { AppError, ValidationException } from '#exceptions/http-errors'
import { ForbiddenResponseSchema } from '#exceptions/schemas/forbidden-response.schema'
import { NoContentResponseSchema } from '#exceptions/schemas/no-content-response.schema'
import { NotFoundResponseSchema } from '#exceptions/schemas/not-found-response.schema'
import { UnauthorizedResponseSchema } from '#exceptions/schemas/unauthorized-response.schema'

// Plugin options
interface ErrorHandlerOptions {
  includeStack?: boolean
  logLevel?: 'debug' | 'error' | 'info' | 'warn'
  sanitizeHeaders?: boolean
}

export const errorHandlerPlugin: FastifyPluginAsync<
  ErrorHandlerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, options = {}) => {
  const { includeStack = false, sanitizeHeaders = true } = options

  // Add error schema to app instance for reuse
  app.addSchema(ForbiddenResponseSchema)
  app.addSchema(NoContentResponseSchema)
  app.addSchema(NotFoundResponseSchema)
  app.addSchema(UnauthorizedResponseSchema)

  // Helper function to sanitize request data for logging
  const sanitizeRequest = (request: FastifyRequest) => {
    const { body, headers, method, params, query, url } = request

    const sanitizedHeaders =
      sanitizeHeaders ?
        {
          'content-type': headers['content-type'],
          'user-agent': headers['user-agent'],
          'x-forwarded-for': headers['x-forwarded-for']
        }
      : headers

    return {
      body: method !== 'GET' ? body : undefined,
      headers: sanitizedHeaders,
      method,
      params,
      query,
      url
    }
  }

  // Helper function to create error response
  const createErrorResponse = (
    error: {
      message: string
    },
    statusCode: number,
    code?: string,
    details?: unknown,
    path?: string
  ) => ({
    error: {
      code,
      details,
      message: error.message,
      path: path ?? 'unknown',
      statusCode,
      timestamp: new Date().toISOString()
    }
  })

  // Set up the global error handler
  app.setErrorHandler(async (error, request, reply) => {
    const requestContext = sanitizeRequest(request)
    const path = request.url

    // Handle custom application errors
    if (error instanceof AppError) {
      const logData = {
        code: error.code,
        error: error.message,
        request: requestContext,
        statusCode: error.statusCode,
        ...(includeStack && { stack: error.stack })
      }

      // Log based on error type
      if (error.statusCode >= 500) {
        request.log.error(logData)
      } else if (error.statusCode >= 400) {
        request.log.warn(logData)
      } else {
        request.log.info(logData)
      }

      const responseData = createErrorResponse(
        error,
        error.statusCode,
        error.code,
        error instanceof ValidationException ?
          (error as ValidationException).details
        : undefined,
        path
      )

      return reply.status(error.statusCode).send(responseData)
    }

    // Handle validation errors from TypeBox/Ajv
    if (error.validation) {
      const logData = {
        error: 'Validation failed',
        request: requestContext,
        validation: error.validation
      }

      request.log.warn(logData)

      const responseData = createErrorResponse(
        { message: 'Validation failed' },
        400,
        'VALIDATION_ERROR',
        error.validation,
        path
      )

      return reply.status(400).send(responseData)
    }

    // Handle other known Fastify errors
    if (error.statusCode) {
      const logData = {
        error: error.message,
        request: requestContext,
        statusCode: error.statusCode,
        ...(includeStack && { stack: error.stack })
      }

      if (error.statusCode >= 500) {
        request.log.error(logData)
      } else {
        request.log.warn(logData)
      }

      const responseData = createErrorResponse(
        error,
        error.statusCode,
        undefined,
        undefined,
        path
      )

      return reply.status(error.statusCode).send(responseData)
    }

    // Handle unexpected errors
    const logData = {
      error: error.message,
      request: requestContext,
      stack: error.stack,
      type: 'UNEXPECTED_ERROR'
    }

    request.log.error(logData)

    const responseData = createErrorResponse(
      { message: 'Internal Server Error' },
      500,
      'INTERNAL_ERROR',
      undefined,
      path
    )

    return reply.status(500).send(responseData)
  })

  // Add hook for handling async errors in routes
  app.addHook('onError', async (_request, _reply, _error) => {
    // This hook runs before the error handler
    // You can add additional processing here if needed
    return
  })
}
