/**
 * Test execute function that doubles a value
 */

/**
 * Input structure
 */
interface TestInput {
  value: number;
}

/**
 * Output structure
 */
interface TestOutput {
  doubled: number;
}

/**
 * Execute function for testing
 * @param input - The input data
 * @returns The output data
 */
export async function executeFunction(input: unknown): Promise<TestOutput> {
  // Validate input
  if (typeof input !== "object" || input === null) {
    throw new Error("Expected object input");
  }

  if (!("value" in input)) {
    throw new Error('Expected "value" field in input');
  }

  const typedInput = input as TestInput;

  if (typeof typedInput.value !== "number") {
    throw new Error('Expected "value" to be a number');
  }

  // Double the value
  return {
    doubled: typedInput.value * 2,
  };
}
