import type { z } from 'zod'

import { PipelineEntitySchema } from '#orchestration/pipelines/entities/pipeline.entity'

export const CreatePipelineDtoSchema: z.ZodObject<{
  description: z.ZodNullable<z.ZodString>
  name: z.ZodNullable<z.ZodString>
}> = PipelineEntitySchema.pick({
  description: true,
  name: true
})

export type CreatePipelineDto = z.infer<typeof CreatePipelineDtoSchema>
