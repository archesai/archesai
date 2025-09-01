import { z } from 'zod'

export const RunpodResponseDtoSchema: z.ZodObject<{
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
  status: z.enum(['IN_PROGRESS', 'COMPLETED', 'FAILED'])
})

export type RunpodResponseDto = z.infer<typeof RunpodResponseDtoSchema>
