/**
 * Orval execution function for container executor.
 * Generates TypeScript/JavaScript clients from OpenAPI specifications.
 */

import fs from "node:fs";
import path from "node:path";
import { generate } from "orval";
import type { Executor } from "#runner";

interface ExecuteInput {
  openapi: string;
}

interface ExecuteOutput {
  files: Record<string, string>;
}

async function executeFunction(input: ExecuteInput): Promise<ExecuteOutput> {
  // Create a temporary working directory inside /tmp to avoid conflicts
  const workDir = "/app";
  console.error(`[orval] Using working directory: ${workDir}`);

  try {
    // Create working directory
    console.error(`[orval] Creating working directory: ${workDir}`);
    fs.mkdirSync(workDir, { recursive: true });

    // Write OpenAPI spec to file
    const specPath = path.join(workDir, "openapi.yaml");
    console.error(`[orval] Writing OpenAPI spec to: ${specPath}`);
    fs.writeFileSync(specPath, input.openapi);

    // Console log head of file for debugging
    const specHead = input.openapi.split("\n").slice(0, 10).join("\n");
    console.error(`[orval] OpenAPI spec head:\n${specHead}\n...`);

    // Redirect stdout to stderr during orval execution to prevent pollution
    const originalStdoutWrite = process.stdout.write.bind(process.stdout);
    process.stdout.write = (
      chunk: string | Uint8Array,
      ...args: any[]
    ): boolean => {
      return process.stderr.write(chunk, ...(args as [any]));
    };

    try {
      // Run orval using the installed package
      console.error("[orval] Running orval code generation");
      await generate({
        input: {
          target: specPath,
        },
        output: {
          biome: true,
          client: "zod",
          mode: "single",
          target: "./generated/zod.ts",
          workspace: "./src",
        },
      });
      console.error("[orval] Zod code generation completed");

      console.error(
        "[orval] Running orval code generation for React Query client",
      );
      await generate({
        input: {
          target: specPath,
        },
        output: {
          allParamsOptional: true,
          biome: true,
          client: "react-query",
          httpClient: "fetch",
          mode: "tags-split",
          override: {
            fetch: {
              includeHttpResponseReturnType: false,
            },
            mutator: {
              name: "customFetch",
              path: "./fetcher.ts",
            },
            query: {
              useInfinite: false,
              useQuery: true,
              useSuspenseInfiniteQuery: false,
              useSuspenseQuery: true,
              version: 5,
            },
          },
          // schemas: 'src/generated/models',
          target: "./generated/orval.ts",
          urlEncodeParameters: true,
          workspace: "./src",
          // propertySortOrder: 'Alphabetical'
        },
      });
      console.error("[orval] Code generation completed");
    } finally {
      // Restore original stdout
      process.stdout.write = originalStdoutWrite;
    }

    // Read generated files
    const generatedDir = path.join(workDir, "src/generated");
    const files: Record<string, string> = {};

    function readDirRecursive(dir: string, base = ""): void {
      const entries = fs.readdirSync(dir, { withFileTypes: true });
      for (const entry of entries) {
        const fullPath = path.join(dir, entry.name);
        const relativePath = path.join(base, entry.name);

        if (entry.isDirectory()) {
          readDirRecursive(fullPath, relativePath);
        } else {
          files[relativePath] = fs.readFileSync(fullPath, "utf8");
        }
      }
    }

    if (fs.existsSync(generatedDir)) {
      console.error(`[orval] Reading generated files from: ${generatedDir}`);
      readDirRecursive(generatedDir);
      console.error(`[orval] Read ${Object.keys(files).length} files`);
    } else {
      console.error(
        `[orval] Warning: Generated directory not found: ${generatedDir}`,
      );
    }

    console.error(`[orval] Successfully completed`);
    return {
      files: files,
    };
  } catch (error) {
    console.error(
      `[orval] Error during execution: ${(error as Error).message}`,
    );
    throw error;
  }
}

/**
 * Export the execute function as an Executor type
 */
const _executeFunction: Executor<ExecuteInput, ExecuteOutput> = executeFunction;

export { _executeFunction as executeFunction };
