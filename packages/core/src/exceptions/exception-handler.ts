import type { FastifyPluginCallback, FastifyRequest } from 'fastify'

import type { ErrorDocument } from '@archesai/schemas'

import { AppError } from '#exceptions/http-errors'

// Plugin options
interface ErrorHandlerOptions {
  includeStack?: boolean
  logLevel?: 'debug' | 'error' | 'info' | 'warn'
  sanitizeHeaders?: boolean
}

export const errorHandlerPlugin: FastifyPluginCallback<ErrorHandlerOptions> = (
  app,
  options = {},
  done
) => {
  const { includeStack = true, sanitizeHeaders = true } = options

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
    message: string,
    statusCode: number,
    title?: string,
    validation: { field: string; message: string }[] = []
  ): ErrorDocument => ({
    error: {
      detail: message,
      status: statusCode.toString(),
      title: title ?? getDefaultTitle(statusCode),
      ...(validation.length > 0 && { validation })
    }
  })

  // Helper function to get default error titles
  const getDefaultTitle = (statusCode: number): string => {
    switch (statusCode) {
      case 400:
        return 'Bad Request'
      case 401:
        return 'Unauthorized'
      case 403:
        return 'Forbidden'
      case 404:
        return 'Not Found'
      case 409:
        return 'Conflict'
      case 422:
        return 'Validation Error'
      case 500:
        return 'Internal Server Error'
      default:
        return 'Error'
    }
  }

  // Helper function to log errors based on status code
  const logError = (
    request: FastifyRequest,
    error: { statusCode?: number },
    logData: Record<string, unknown>
  ) => {
    const statusCode = error.statusCode ?? 500

    if (statusCode >= 500) {
      request.log.error(logData)
    } else if (statusCode >= 400) {
      request.log.warn(logData)
    } else {
      request.log.info(logData)
    }
  }

  // Set up the global error handler
  app.setErrorHandler(async (error, request, reply) => {
    const requestContext = sanitizeRequest(request)

    // Handle validation errors from TypeBox/Ajv FIRST (highest priority)
    if (error.validation) {
      const logData = {
        error: error.message,
        request: requestContext,
        validation: error.validation
      }
      logError(request, { statusCode: 400 }, logData)
      const responseData = createErrorResponse(
        error.message,
        400,
        error.name,
        error.validation.map((v) => ({
          field: v.instancePath.replace(/^\//, ''),
          message: v.message ?? 'Invalid value'
        }))
      )
      return reply.status(400).send(responseData)
    }

    // Handle custom application errors (after validation)
    if (error instanceof AppError) {
      const logData = {
        code: error.code,
        error: error.message,
        request: requestContext,
        statusCode: error.statusCode,
        ...(includeStack && { stack: error.stack })
      }
      logError(request, error, logData)
      const responseData = createErrorResponse(
        error.message,
        error.statusCode,
        error.title
      )
      return reply.status(error.statusCode).send(responseData)
    }

    // Handle other known Fastify errors
    if (error.statusCode) {
      const logData = {
        error: error.message,
        request: requestContext,
        statusCode: error.statusCode,
        ...(includeStack && { stack: error.stack })
      }
      logError(request, error, logData)
      const responseData = createErrorResponse(error.message, error.statusCode)
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
    const responseData = createErrorResponse('Internal Server Error', 500)
    return reply.status(500).send(responseData)
  })

  // Add hook for handling async errors in routes
  app.addHook('onError', async (_request, _reply, _error) => {
    // This hook runs before the error handler
    // You can add additional processing here if needed
    return
  })

  done()
}
