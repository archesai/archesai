export class AppError extends Error {
  public readonly code?: string
  public readonly detail: string
  public readonly meta?: Record<string, unknown>
  public readonly status: string = '500' // JSON:API wants string
  public readonly title: string

  constructor({
    code,
    message = 'An unexpected error occurred',
    detail = message,
    meta,
    status = '500',
    title = 'Internal Server Error'
  }: {
    code?: string
    detail?: string
    message?: string
    meta?: Record<string, unknown>
    status?: number | string
    title?: string
  } = {}) {
    super(message)
    this.name = 'AppError'
    this.title = title
    this.detail = detail
    this.status = String(status)
    if (code) {
      this.code = code
    }
    if (meta) {
      this.meta = meta
    }
  }
}

export class BadRequestException extends AppError {
  constructor(detail = 'Bad request', meta?: Record<string, unknown>) {
    super({
      code: 'BAD_REQUEST',
      detail,
      message: detail,
      ...(meta ? { meta } : {}),
      status: 400,
      title: 'Bad Request'
    })
    this.name = 'BadRequestException'
  }
}

export class ConflictException extends AppError {
  constructor(detail = 'Conflict', meta?: Record<string, unknown>) {
    super({
      code: 'CONFLICT',
      detail,
      message: detail,
      ...(meta ? { meta } : {}),
      status: 409,
      title: 'Conflict'
    })
    this.name = 'ConflictException'
  }
}

export class ForbiddenException extends AppError {
  constructor(detail = 'Forbidden', meta?: Record<string, unknown>) {
    super({
      code: 'FORBIDDEN',
      detail,
      message: detail,
      ...(meta ? { meta } : {}),
      status: 403,
      title: 'Forbidden'
    })
    this.name = 'ForbiddenException'
  }
}

export class InternalServerErrorException extends AppError {
  constructor(
    detail = 'Internal Server Error',
    meta?: Record<string, unknown>
  ) {
    super({
      code: 'INTERNAL_SERVER_ERROR',
      detail,
      message: detail,
      ...(meta ? { meta } : {}),
      status: 500,
      title: 'Internal Server Error'
    })
    this.name = 'InternalServerErrorException'
  }
}

export class NotFoundException extends AppError {
  constructor(detail = 'Not Found', meta?: Record<string, unknown>) {
    super({
      code: 'NOT_FOUND',
      detail,
      message: detail,
      ...(meta ? { meta } : {}),
      status: 404,
      title: 'Not Found'
    })
    this.name = 'NotFoundException'
  }
}

export class UnauthorizedException extends AppError {
  constructor(detail = 'Unauthorized', meta?: Record<string, unknown>) {
    super({
      code: 'UNAUTHORIZED',
      detail,
      message: detail,
      ...(meta ? { meta } : {}),
      status: 401,
      title: 'Unauthorized'
    })
    this.name = 'UnauthorizedException'
  }
}

export class UnprocessableEntityException extends AppError {
  constructor(detail = 'Unprocessable Entity', meta?: Record<string, unknown>) {
    super({
      code: 'UNPROCESSABLE_ENTITY',
      detail,
      message: detail,
      ...(meta ? { meta } : {}),
      status: 422,
      title: 'Unprocessable Entity'
    })
    this.name = 'UnprocessableEntityException'
  }
}
