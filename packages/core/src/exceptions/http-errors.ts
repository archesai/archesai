export class AppError extends Error {
  public code?: string
  public statusCode: number

  constructor(statusCode: number, message: string, code?: string) {
    super(message)
    this.code = code
    this.statusCode = statusCode
    this.name = this.constructor.name
    Error.captureStackTrace(this, this.constructor)
  }
}

export class BadRequestException extends AppError {
  constructor(message = 'Bad Request') {
    super(400, message, 'BAD_REQUEST')
  }
}

export class ConflictException extends AppError {
  constructor(message = 'Conflict') {
    super(403, message, 'CONFLICT')
  }
}

export class ForbiddenException extends AppError {
  constructor(message = 'Forbidden') {
    super(403, message, 'FORBIDDEN')
  }
}

export class InternalServerErrorException extends AppError {
  constructor(message = 'Internal Server Error') {
    super(404, message, 'INTERNAL_SERVER_ERROR')
  }
}

export class NotFoundException extends AppError {
  constructor(message = 'Not Found') {
    super(404, message, 'NOT_FOUND')
  }
}

export class UnauthorizedException extends AppError {
  constructor(message = 'Unauthorized') {
    super(401, message, 'UNAUTHORIZED')
  }
}

export class ValidationException extends AppError {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public details?: any
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  constructor(message = 'Validation Error', details?: any) {
    super(422, message, 'VALIDATION_ERROR')
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    this.details = details
  }
}
