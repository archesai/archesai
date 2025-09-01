import type { z } from 'zod'

import { CreateRunDtoSchema } from '#orchestration/runs/dto/create-run.dto'

export const UpdateRunDtoSchema: z.ZodObject<{
  pipelineId: z.ZodOptional<z.ZodNullable<z.ZodString>>
}> = CreateRunDtoSchema.partial()

export type UpdateRunDto = z.infer<typeof UpdateRunDtoSchema>
