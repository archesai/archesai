import type { Logger } from '#logging/logger'

/**
 * Retries a given asynchronous function a specified number of times with exponential backoff.
 *
 * @template T - The type of the value returned by the function.
 * @param {Logger} logger - The logger instance used to log warnings and errors.
 * @param {() => Promise<T>} fn - The asynchronous function to retry.
 * @param {number} maxAttempts - The maximum number of retry attempts.
 * @returns {Promise<T>} - A promise that resolves to the result of the function or rejects with an error after all retry attempts have been exhausted.
 */
export const retry = async <T>(
  logger: Logger,
  fn: () => Promise<T>,
  maxAttempts: number
): Promise<T> => {
  const execute = async (attempt: number): Promise<T> => {
    try {
      return await fn()
    } catch (err) {
      if (attempt <= maxAttempts) {
        const nextAttempt = attempt + 1
        const delayInSeconds = Math.max(
          Math.min(
            Math.pow(2, nextAttempt) + randInt(-nextAttempt, nextAttempt),
            10
          ),
          1
        )
        logger.warn(`retrying`, {
          attempt,
          err,
          maxAttempts
        })
        return await delay(() => execute(nextAttempt), delayInSeconds * 1000)
      } else {
        logger.error('error', { err })
        throw err
      }
    }
  }
  return execute(1)
}

const delay = async <T>(fn: () => Promise<T>, ms: number): Promise<T> =>
  new Promise((resolve) =>
    setTimeout(() => {
      resolve(fn())
    }, ms)
  )

const randInt = (min: number, max: number) =>
  Math.floor(Math.random() * (max - min + 1) + min)
