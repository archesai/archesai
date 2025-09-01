export function toError(error: unknown): Error {
  if (error instanceof Error) {
    return error
  }

  const err = new Error('Unknown error', {
    cause: error
  })
  err.name = 'UnknownError'
  return err
}
