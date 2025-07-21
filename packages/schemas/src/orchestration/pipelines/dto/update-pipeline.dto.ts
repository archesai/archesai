import type { z } from 'zod'

import { CreatePipelineDtoSchema } from '#orchestration/pipelines/dto/create-pipeline.dto'

export const UpdatePipelineDtoSchema: z.ZodObject<{
  description: z.ZodOptional<z.ZodNullable<z.ZodString>>
  name: z.ZodOptional<z.ZodNullable<z.ZodString>>
}> = CreatePipelineDtoSchema.partial()

export type UpdatePipelineDto = z.infer<typeof UpdatePipelineDtoSchema>
