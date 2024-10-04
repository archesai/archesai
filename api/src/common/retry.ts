import { Logger } from "@nestjs/common";

export const retry = async <T>(
  logger: Logger,
  fn: () => Promise<T>,
  maxAttempts: number,
) => {
  const execute = async (attempt: any): Promise<T> => {
    try {
      return await fn();
    } catch (err) {
      if (attempt <= maxAttempts) {
        const nextAttempt = attempt + 1;
        const delayInSeconds = Math.max(
          Math.min(
            Math.pow(2, nextAttempt) + randInt(-nextAttempt, nextAttempt),
            30,
          ),
          1,
        );
        logger.warn(
          `Retrying after ${delayInSeconds} seconds due to error: ${err}`,
        );
        return await delay(() => execute(nextAttempt), delayInSeconds * 1000);
      } else {
        logger.error(err);
        throw err;
      }
    }
  };
  return execute(1);
};

const delay = async <T>(fn: () => Promise<T>, ms: number): Promise<T> =>
  new Promise((resolve) => setTimeout(() => resolve(fn()), ms));

const randInt = (min: number, max: number) =>
  Math.floor(Math.random() * (max - min + 1) + min);
