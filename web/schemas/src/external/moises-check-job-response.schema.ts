import { z } from 'zod'

export const MoisesCheckJobResponseSchema: z.ZodObject<{
  result: z.ZodObject<{
    Bass: z.ZodString
    Drums: z.ZodString
    Other: z.ZodString
    Vocals: z.ZodString
  }>
  status: z.ZodEnum<{
    FAILED: 'FAILED'
    PENDING: 'PENDING'
    PROCESSING: 'PROCESSING'
    SUCCEEDED: 'SUCCEEDED'
  }>
}> = z.object({
  result: z.object({
    Bass: z.string(),
    Drums: z.string(),
    Other: z.string(),
    Vocals: z.string()
  }),
  status: z.enum(['FAILED', 'PENDING', 'PROCESSING', 'SUCCEEDED'])
})

export type MoisesCheckJobResponse = z.infer<
  typeof MoisesCheckJobResponseSchema
>
