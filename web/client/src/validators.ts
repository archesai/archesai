import { z } from "zod";
import type { Problem } from "#generated/orval.schemas";

/**
 * Problem schema based on RFC 7807 Problem Details format
 */
export const problemSchema: z.ZodType<Problem> = z.object({
  detail: z.string(),
  errors: z.array(
    z.object({
      code: z.string().optional(),
      field: z.string(),
      message: z.string(),
    }),
  ),
  instance: z.string(),
  status: z.number().min(100).max(599),
  title: z.string(),
  type: z.string(),
});

/**
 * Type guards using Zod schemas
 */
export const isProblem = (
  obj: unknown,
): obj is z.infer<typeof problemSchema> => {
  return problemSchema.safeParse(obj).success;
};
