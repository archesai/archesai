/**
 * Default execute function that throws an error.
 * This file should be replaced by mounting a custom execute.ts/execute.js
 * or by building a derived image with a custom implementation.
 */

import type { Executor } from "#runner";

/**
 * Execute function - must be implemented by user
 * @param input - The input data
 * @returns The output data
 * @throws Error if no custom implementation is provided
 */
async function executeFunction(_input: unknown): Promise<unknown> {
  throw new Error(
    "No execution function provided. " +
      "Mount a custom execute.ts at /app/execute.ts ",
  );
}

/**
 * Export the execute function as an Executor type
 */
const _executeFunction: Executor<unknown, unknown> = executeFunction;

export { _executeFunction as executeFunction };
