import { z } from 'zod'

export const RunpodResponseSchema: z.ZodObject<{
  id: z.ZodString
  output: z.ZodString
  status: z.ZodEnum<{
    COMPLETED: 'COMPLETED'
    FAILED: 'FAILED'
    IN_PROGRESS: 'IN_PROGRESS'
  }>
}> = z.object({
  id: z.string(),
  output: z.string(),
  status: z.enum(['COMPLETED', 'FAILED', 'IN_PROGRESS'])
})
