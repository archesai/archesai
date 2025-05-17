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
