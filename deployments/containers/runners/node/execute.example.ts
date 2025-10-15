/**
 * Example custom execute function
 *
 * This file shows how to create a custom execute module for the node runner.
 *
 * To use:
 * 1. Copy this file and implement your logic
 * 2. Mount it to the container:
 *    docker run -i --rm -v ./my-execute.js:/app/execute.js archesai/runner-node:latest
 *
 * Or build a custom image:
 * 1. Create a Dockerfile:
 *    FROM archesai/runner-node:latest
 *    RUN npm install --omit=dev your-dependencies
 *    COPY execute.js ./execute.js
 * 2. Build: docker build -t my-generator:latest .
 */

/**
 * Example input structure
 */
interface ExampleInput {
  values: number[];
}

/**
 * Example output structure
 */
interface ExampleOutput {
  count: number;
  sum: number;
  mean: number;
  min: number;
  max: number;
}

/**
 * Execute function - implement your custom logic here
 *
 * @param input - The input data from the request
 * @returns The output data
 * @throws Error if execution fails
 */
export async function executeFunction(input: unknown): Promise<ExampleOutput> {
  // Type guard to validate input structure
  if (
    typeof input !== "object" ||
    input === null ||
    !("values" in input) ||
    !Array.isArray(input.values)
  ) {
    throw new Error("Expected 'values' array in input");
  }

  const typedInput = input as ExampleInput;
  const { values } = typedInput;

  // Validate that all values are numbers
  if (!values.every((v) => typeof v === "number")) {
    throw new Error("All values must be numbers");
  }

  // Calculate statistics
  const sum = values.reduce((a, b) => a + b, 0);
  const mean = sum / values.length;

  return {
    count: values.length,
    max: Math.max(...values),
    mean,
    min: Math.min(...values),
    sum,
  };
}

/**
 * Example usage:
 *
 * Input:
 * {
 *   "input": {
 *     "values": [1, 2, 3, 4, 5]
 *   }
 * }
 *
 * Output:
 * {
 *   "ok": true,
 *   "output": {
 *     "count": 5,
 *     "sum": 15,
 *     "mean": 3,
 *     "min": 1,
 *     "max": 5
 *   }
 * }
 */
