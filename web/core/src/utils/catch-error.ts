import { toError } from '#utils/to-error'

export const catchErrorAsync = async <
  T,
  E extends new (message?: string) => Error
>(
  promise: Promise<T>,
  errorsToCatch?: E[]
): Promise<[Error] | [InstanceType<E>] | [undefined, T]> => {
  try {
    const data = await promise
    return [undefined, data]
  } catch (error) {
    if (errorsToCatch == undefined) {
      return [toError(error)]
    }

    for (const e of errorsToCatch) {
      if (error instanceof e) {
        return [error as InstanceType<E>]
      }
    }

    throw error
  }
}

export const catchError = <T, E extends new (message?: string) => Error>(
  fn: () => T,
  errorsToCatch?: E[]
): [Error] | [InstanceType<E>] | [undefined, T] => {
  try {
    const data = fn()
    return [undefined, data]
  } catch (error) {
    if (errorsToCatch == undefined) {
      return [toError(error)]
    }

    for (const e of errorsToCatch) {
      if (error instanceof e) {
        return [error as InstanceType<E>]
      }
    }

    throw error
  }
}
