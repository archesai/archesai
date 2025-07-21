import { z } from 'zod'

export const PageSchema: z.ZodObject<{
  number: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
  size: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
}> = z
  .object({
    number: z.number().int().min(1).max(Number.MAX_VALUE).default(1).optional(),
    size: z.number().int().min(1).max(100).default(10).optional()
  })
  .meta({
    description: 'Pagination configuration with page number and size',
    id: 'Page'
  })
