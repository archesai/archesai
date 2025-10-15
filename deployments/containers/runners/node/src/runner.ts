#!/usr/bin/env node

/**
 * Node.js runner for container executor.
 * Reads JSON from stdin, validates against schemas, executes function, and writes result to stdout.
 */

import type { JSONSchemaType } from "ajv";
import { Ajv } from "ajv";

// Dynamically import the execute module
// This can be mounted at runtime or copied in derived images
import { executeFunction } from "#execute";

// Initialize AJV for schema validation
const ajv = new Ajv({ allErrors: true, strict: true });

/**
 * Executor interface
 */
export type Executor<I, O> = (input: I) => Promise<O>;

/**
 * Request structure from stdin
 */
interface Request {
  schema_in?: JSONSchemaType<unknown>;
  schema_out?: JSONSchemaType<unknown>;
  input: unknown;
}

/**
 * Success response structure
 */
interface SuccessResponse {
  ok: true;
  output: unknown;
}

/**
 * Error response structure
 */
interface ErrorResponse {
  ok: false;
  error: {
    message: string;
    details?: string | null;
  };
}

/**
 * Response structure to stdout
 */
// type Response = SuccessResponse | ErrorResponse;

/**
 * Read all data from stdin
 */
async function readStdin(): Promise<string> {
  return new Promise((resolve, reject) => {
    const chunks: Buffer[] = [];

    process.stdin.on("data", (chunk: Buffer) => {
      chunks.push(chunk);
    });

    process.stdin.on("end", () => {
      resolve(Buffer.concat(chunks).toString("utf8"));
    });

    process.stdin.on("error", (err: Error) => {
      reject(err);
    });
  });
}

/**
 * Main entry point
 */
async function main(): Promise<void> {
  try {
    // Read input from stdin
    const rawInput = await readStdin();
    if (!rawInput) {
      throw new Error("No input provided");
    }

    // Parse the request
    const request = JSON.parse(rawInput) as Request;

    // Extract components
    const schemaIn = request.schema_in;
    const schemaOut = request.schema_out;
    const inputData = request.input;

    if (inputData === undefined || inputData === null) {
      throw new Error("Missing 'input' field in request");
    }

    // Validate input against schema if provided
    if (schemaIn) {
      const validateIn = ajv.compile(schemaIn);
      if (!validateIn(inputData)) {
        throw new Error(
          `Input validation failed: ${JSON.stringify(validateIn.errors)}`,
        );
      }
    }

    // Execute the function
    const outputData = await executeFunction(inputData);

    // Validate output against schema if provided
    if (schemaOut) {
      const validateOut = ajv.compile(schemaOut);
      if (!validateOut(outputData)) {
        throw new Error(
          `Output validation failed: ${JSON.stringify(validateOut.errors)}`,
        );
      }
    }

    // Return success response
    const response: SuccessResponse = {
      ok: true,
      output: outputData,
    };
    const responseStr = JSON.stringify(response);
    process.stdout.write(responseStr, () => {
      // Ensure all data is written before exiting
      process.exit(0);
    });
  } catch (error) {
    // Return error response
    const errorMessage = error instanceof Error ? error.message : String(error);
    const errorStack = error instanceof Error ? (error.stack ?? null) : null;

    const response: ErrorResponse = {
      error: {
        details: errorStack,
        message: errorMessage,
      },
      ok: false,
    };
    const responseStr = JSON.stringify(response);
    process.stdout.write(responseStr, () => {
      // Ensure all data is written before exiting
      process.exit(0);
    });
  }
}

main().catch((err: unknown) => {
  const errorMessage = err instanceof Error ? err.message : String(err);

  const response: ErrorResponse = {
    error: {
      message: `Unexpected error: ${errorMessage}`,
    },
    ok: false,
  };
  const responseStr = JSON.stringify(response);
  process.stdout.write(responseStr, () => {
    // Ensure all data is written before exiting
    process.exit(0);
  });
});
