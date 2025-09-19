import type { Problem } from "@archesai/client";
import Ajv, { type ValidateFunction } from "ajv";
import addFormats from "ajv-formats";
import type { JSONSchema7Definition } from "json-schema";

// Import the bundled OpenAPI spec
import openApiSpec from "../../../../api/openapi.bundled.json";

// Define types for OpenAPI schemas with explicit properties
interface OpenAPISchema {
  type?: string;
  properties?: Record<string, JSONSchema7Definition>;
  required?: string[];
  allOf?: JSONSchema7Definition[];
  $ref?: string;
  additionalProperties?: boolean | JSONSchema7Definition;
}

// Initialize AJV with formats support
const ajv = new Ajv({
  allErrors: true,
  removeAdditional: false, // Don't remove additional properties
  strict: false,
  verbose: true,
});
addFormats(ajv);

// Add all component schemas to AJV
if (openApiSpec.components?.schemas) {
  Object.entries(openApiSpec.components.schemas).forEach(([name, schema]) => {
    ajv.addSchema(schema, `#/components/schemas/${name}`);
  });
}

// Get the Problem schema directly from OpenAPI
const problemSchema = openApiSpec.components?.schemas?.Problem as
  | OpenAPISchema
  | undefined;
const sessionSchema = openApiSpec.components?.schemas?.Session as
  | OpenAPISchema
  | undefined;

// For Session, we need to handle the allOf properly
// The issue is that Session uses allOf with Base, and one schema has additionalProperties: false
// We need to merge the schemas for proper validation
let mergedSessionSchema: OpenAPISchema | null = null;
if (sessionSchema?.allOf) {
  // Create a merged schema that combines all properties
  mergedSessionSchema = {
    properties: {},
    required: [],
    type: "object",
  } as OpenAPISchema;

  // Merge all schemas in allOf
  for (const subSchema of sessionSchema.allOf) {
    if (typeof subSchema === "object" && subSchema !== null) {
      if ("$ref" in subSchema && subSchema.$ref) {
        // Resolve the reference
        const refName = subSchema.$ref.split("/").pop();
        if (refName && openApiSpec.components?.schemas) {
          const schemas = openApiSpec.components.schemas as Record<
            string,
            OpenAPISchema
          >;
          const refSchema = schemas[refName];
          if (refSchema?.properties) {
            Object.assign(
              mergedSessionSchema.properties || {},
              refSchema.properties,
            );
            if (refSchema.required && mergedSessionSchema.required) {
              (mergedSessionSchema.required as string[]).push(
                ...(refSchema.required as string[]),
              );
            }
          }
        }
      } else if ("properties" in subSchema) {
        const schema = subSchema as OpenAPISchema;
        if (schema.properties) {
          Object.assign(
            mergedSessionSchema.properties || {},
            schema.properties,
          );
        }
        if (schema.required && mergedSessionSchema.required) {
          (mergedSessionSchema.required as string[]).push(
            ...(schema.required as string[]),
          );
        }
      }
    }
  }

  // Remove duplicates from required array
  if (mergedSessionSchema.required) {
    mergedSessionSchema.required = [
      ...new Set(mergedSessionSchema.required as string[]),
    ];
  }
} else {
  mergedSessionSchema = sessionSchema || null;
}

// Compile validators from the actual OpenAPI schemas
const validateProblem: ValidateFunction | null = problemSchema
  ? ajv.compile(problemSchema)
  : null;
const validateSession: ValidateFunction | null = mergedSessionSchema
  ? ajv.compile(mergedSessionSchema)
  : null;

/**
 * Type guard for Problem using the actual OpenAPI schema.
 * Validates that an unknown object conforms to the RFC 7807 Problem Details format.
 *
 * @param obj - The object to validate
 * @returns True if the object is a valid Problem, false otherwise
 *
 * @example
 * ```typescript
 * const response = await fetch('/api/endpoint');
 * const data = await response.json();
 *
 * if (isProblem(data)) {
 *   console.error(`Error ${data.status}: ${data.detail}`);
 * }
 * ```
 */
export const isProblem = (obj: unknown): obj is Problem => {
  return validateProblem ? validateProblem(obj) : false;
};

/**
 * Type guard for Session response using the actual OpenAPI schema.
 * Validates that an unknown object contains a valid Session in its data property.
 *
 * @param obj - The object to validate
 * @returns True if the object contains a valid Session, false otherwise
 *
 * @example
 * ```typescript
 * const response = await fetch('/api/sessions/current');
 * const data = await response.json();
 *
 * if (isSessionResponse(data)) {
 *   console.log('User ID:', data.data.userID);
 * }
 * ```
 */
export const isSessionResponse = <T extends { data: unknown }>(
  obj: unknown,
): obj is T => {
  if (!obj || typeof obj !== "object") return false;
  const response = obj as { data: unknown };
  return !!response.data && !!validateSession && validateSession(response.data);
};

/**
 * Get validation errors from the last validation.
 * Returns the AJV validation errors from the most recent validation attempt.
 *
 * @returns Array of validation errors or null if no errors
 *
 * @example
 * ```typescript
 * if (!isProblem(data)) {
 *   const errors = getValidationErrors();
 *   console.error('Validation failed:', errors);
 * }
 * ```
 */
export const getValidationErrors = (): ValidateFunction["errors"] => {
  if (validateProblem?.errors) return validateProblem.errors;
  if (validateSession?.errors) return validateSession.errors;
  return null;
};

// Export validators for direct use if needed
export { validateProblem, validateSession };
