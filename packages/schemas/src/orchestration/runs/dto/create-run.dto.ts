import type { z } from 'zod'

import { RunEntitySchema } from '#orchestration/runs/entities/run.entity'

export const CreateRunDtoSchema: z.ZodObject<{
  pipelineId: z.ZodNullable<z.ZodString>
}> = RunEntitySchema.pick({
  pipelineId: true
})

export type CreateRunDto = z.infer<typeof CreateRunDtoSchema>
