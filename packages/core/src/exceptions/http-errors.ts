export class AppError extends Error {
  public code: string
  public statusCode: number
  public title: string

  constructor(
    statusCode: number,
    message: string,
    code: string,
    title?: string
  ) {
    super(message)
    this.statusCode = statusCode
    this.code = code
    this.title = title ?? this.getDefaultTitle(statusCode)
    this.name = this.constructor.name
    Error.captureStackTrace(this, this.constructor)
  }

  private getDefaultTitle(statusCode: number): string {
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
}

export class BadRequestException extends AppError {
  constructor(message = 'Bad Request') {
    super(400, message, 'BAD_REQUEST')
  }
}

export class ConflictException extends AppError {
  constructor(message = 'Conflict') {
    super(409, message, 'CONFLICT')
  }
}

export class FileNotFoundException extends AppError {
  constructor(key: string) {
    super(404, `File not found: ${key}`, 'FILE_NOT_FOUND')
  }
}

export class ForbiddenException extends AppError {
  constructor(message = 'Forbidden') {
    super(403, message, 'FORBIDDEN')
  }
}

export class InternalServerErrorException extends AppError {
  constructor(message = 'Internal Server Error') {
    super(500, message, 'INTERNAL_SERVER_ERROR')
  }
}

export class NotFoundException extends AppError {
  constructor(message = 'Not Found') {
    super(404, message, 'NOT_FOUND')
  }
}

export class StorageException extends AppError {
  constructor(message = 'Storage Error') {
    super(500, message, 'STORAGE_ERROR')
  }
}

export class UnauthorizedException extends AppError {
  constructor(message = 'Unauthorized') {
    super(401, message, 'UNAUTHORIZED')
  }
}

export class ValidationException extends AppError {
  public details?: Record<string, unknown>
  constructor(message = 'Validation Error', details?: Record<string, unknown>) {
    super(422, message, 'VALIDATION_ERROR')
    if (details) {
      this.details = details
    }
  }
}
