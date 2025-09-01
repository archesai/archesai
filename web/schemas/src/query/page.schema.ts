import { z } from 'zod'

export const PageSchema: z.ZodObject<{
  number: z.ZodOptional<z.ZodDefault<z.ZodCoercedNumber>>
  size: z.ZodOptional<z.ZodDefault<z.ZodCoercedNumber>>
}> = z
  .object({
    number: z.coerce
      .number()
      .int()
      .min(1)
      .max(Number.MAX_VALUE)
      .default(1)
      .optional(),
    size: z.coerce.number().int().min(1).max(100).default(10).optional()
  })
  .meta({
    description: 'Pagination configuration with page number and size',
    id: 'Page'
  })
