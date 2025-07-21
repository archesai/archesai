import { z } from 'zod'

export const SortSchema: z.ZodObject<{
  field: z.ZodString
  order: z.ZodEnum<{
    asc: 'asc'
    desc: 'desc'
  }>
}> = z
  .object({
    field: z.string(),
    order: z.enum(['asc', 'desc'])
  })
  .meta({
    description: 'Sorting configuration with field and order',
    id: 'Sort'
  })
