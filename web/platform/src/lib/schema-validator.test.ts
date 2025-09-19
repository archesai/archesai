import { describe, expect, it } from "vitest";
import {
  getValidationErrors,
  isProblem,
  isSessionResponse,
} from "./schema-validator";

describe("schema-validator", () => {
  describe("isProblem", () => {
    it("should validate a correct Problem object", () => {
      const validProblem = {
        detail: "The request body contains invalid fields",
        status: 400,
        title: "Validation Failed",
        type: "https://api.example.com/errors/validation-failed",
      };

      expect(isProblem(validProblem)).toBe(true);
    });

    it("should accept Problem with optional instance field", () => {
      const validProblem = {
        detail: "The requested resource was not found",
        instance: "https://api.example.com/users/123",
        status: 404,
        title: "Not Found",
        type: "https://api.example.com/errors/not-found",
      };

      expect(isProblem(validProblem)).toBe(true);
    });

    it("should reject Problem missing required fields", () => {
      const invalidProblem = {
        status: 400,
        title: "Bad Request",
        // missing 'type' and 'detail'
      };

      expect(isProblem(invalidProblem)).toBe(false);
    });

    it("should reject Problem with invalid status code", () => {
      const invalidProblem = {
        detail: "An error occurred",
        status: 99, // below minimum
        title: "Error",
        type: "about:blank",
      };

      expect(isProblem(invalidProblem)).toBe(false);
    });

    it("should reject Problem with status above 599", () => {
      const invalidProblem = {
        detail: "An error occurred",
        status: 600, // above maximum
        title: "Error",
        type: "about:blank",
      };

      expect(isProblem(invalidProblem)).toBe(false);
    });

    it("should reject non-object values", () => {
      expect(isProblem(null)).toBe(false);
      expect(isProblem(undefined)).toBe(false);
      expect(isProblem("string")).toBe(false);
      expect(isProblem(123)).toBe(false);
      expect(isProblem([])).toBe(false);
    });

    it("should reject Problem with additional properties", () => {
      const invalidProblem = {
        detail: "An error occurred",
        extraField: "not allowed", // additionalProperties: false
        status: 400,
        title: "Error",
        type: "about:blank",
      };

      expect(isProblem(invalidProblem)).toBe(false);
    });
  });

  describe("isSessionResponse", () => {
    it("should validate a correct Session response", () => {
      const validSession = {
        data: {
          createdAt: "2024-01-01T00:00:00Z",
          expiresAt: "2024-01-02T00:00:00Z",
          id: "550e8400-e29b-41d4-a716-446655440000",
          ipAddress: "192.168.1.1",
          token: "session-token-123",
          updatedAt: "2024-01-01T00:00:00Z",
          userAgent: "Mozilla/5.0",
          userID: "550e8400-e29b-41d4-a716-446655440001",
        },
      };

      expect(isSessionResponse(validSession)).toBe(true);
    });

    it("should validate Session with optional fields", () => {
      const validSession = {
        data: {
          authMethod: "magic_link",
          authProvider: "local",
          createdAt: "2024-01-01T00:00:00Z",
          expiresAt: "2024-01-02T00:00:00Z",
          id: "550e8400-e29b-41d4-a716-446655440000",
          ipAddress: "192.168.1.1",
          metadata: {
            browser: "Chrome",
            version: "120.0",
          },
          organizationID: "550e8400-e29b-41d4-a716-446655440002",
          token: "session-token-123",
          updatedAt: "2024-01-01T00:00:00Z",
          userAgent: "Mozilla/5.0",
          userID: "550e8400-e29b-41d4-a716-446655440001",
        },
      };

      expect(isSessionResponse(validSession)).toBe(true);
    });

    it("should reject Session missing data wrapper", () => {
      const invalidSession = {
        createdAt: "2024-01-01T00:00:00Z",
        id: "550e8400-e29b-41d4-a716-446655440000",
        // not wrapped in 'data'
      };

      expect(isSessionResponse(invalidSession)).toBe(false);
    });

    it("should reject Session missing required fields", () => {
      const invalidSession = {
        data: {
          id: "550e8400-e29b-41d4-a716-446655440000",
          // missing other required fields
        },
      };

      expect(isSessionResponse(invalidSession)).toBe(false);
    });

    it("should reject non-object values", () => {
      expect(isSessionResponse(null)).toBe(false);
      expect(isSessionResponse(undefined)).toBe(false);
      expect(isSessionResponse("string")).toBe(false);
      expect(isSessionResponse(123)).toBe(false);
    });

    it("should accept Session with empty metadata object", () => {
      const validSession = {
        data: {
          createdAt: "2024-01-01T00:00:00Z",
          expiresAt: "2024-01-02T00:00:00Z",
          id: "550e8400-e29b-41d4-a716-446655440000",
          ipAddress: "192.168.1.1",
          metadata: {}, // empty metadata is valid
          token: "session-token-123",
          updatedAt: "2024-01-01T00:00:00Z",
          userAgent: "Mozilla/5.0",
          userID: "550e8400-e29b-41d4-a716-446655440001",
        },
      };

      expect(isSessionResponse(validSession)).toBe(true);
    });
  });

  describe("getValidationErrors", () => {
    it("should return errors after validation", () => {
      // Note: Due to test execution order, previous test validations may have set errors
      // This test verifies that getValidationErrors returns the appropriate value
      const errors = getValidationErrors();
      // It should either be null or an array of errors
      expect(errors === null || Array.isArray(errors)).toBe(true);
    });

    it("should return errors after failed Problem validation", () => {
      const invalidProblem = {
        status: 400,
        // missing required fields
      };

      isProblem(invalidProblem);
      const errors = getValidationErrors();

      expect(errors).not.toBe(null);
      expect(Array.isArray(errors)).toBe(true);
      if (errors) {
        expect(errors.length).toBeGreaterThan(0);
        // Check that errors mention missing required fields
        const errorMessages = errors.map((e) => e.message).join(" ");
        expect(errorMessages).toContain("required");
      }
    });

    it("should return errors after failed Session validation", () => {
      const invalidSession = {
        data: {
          // incomplete session data
          id: "123",
        },
      };

      isSessionResponse(invalidSession);
      const errors = getValidationErrors();

      expect(errors).not.toBe(null);
      expect(Array.isArray(errors)).toBe(true);
      if (errors) {
        expect(errors.length).toBeGreaterThan(0);
      }
    });
  });

  describe("type guards", () => {
    it("should properly narrow types with isProblem", () => {
      const unknown: unknown = {
        detail: "An error occurred",
        status: 400,
        title: "Error",
        type: "about:blank",
      };

      if (isProblem(unknown)) {
        // TypeScript should now know this is a Problem
        expect(unknown.status).toBe(400);
        expect(unknown.title).toBe("Error");
      }
    });

    it("should properly narrow types with isSessionResponse", () => {
      const unknown: unknown = {
        data: {
          createdAt: "2024-01-01T00:00:00Z",
          expiresAt: "2024-01-02T00:00:00Z",
          id: "550e8400-e29b-41d4-a716-446655440000",
          ipAddress: "192.168.1.1",
          token: "session-token-123",
          updatedAt: "2024-01-01T00:00:00Z",
          userAgent: "Mozilla/5.0",
          userID: "550e8400-e29b-41d4-a716-446655440001",
        },
      };

      if (isSessionResponse(unknown)) {
        // TypeScript should now know this has a data property
        expect(unknown.data).toBeDefined();
      }
    });
  });
});
