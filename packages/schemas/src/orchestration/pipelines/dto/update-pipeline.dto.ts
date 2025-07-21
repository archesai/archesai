import type { z } from 'zod'

import { CreatePipelineDtoSchema } from '#orchestration/pipelines/dto/create-pipeline.dto'

export const UpdatePipelineDtoSchema: z.ZodObject<{
  description: z.ZodOptional<z.ZodNullable<z.ZodString>>
  name: z.ZodOptional<z.ZodNullable<z.ZodString>>
  steps: z.ZodOptional<
    z.ZodArray<
      z.ZodObject<{
        createdAt: z.ZodString
        dependents: z.ZodArray<
          z.ZodObject<{
            pipelineStepId: z.ZodString
          }>
        >
        id: z.ZodString
        pipelineId: z.ZodString
        prerequisites: z.ZodArray<
          z.ZodObject<{
            pipelineStepId: z.ZodString
          }>
        >
        toolId: z.ZodString
        updatedAt: z.ZodString
      }>
    >
  >
}> = CreatePipelineDtoSchema.partial()

export type UpdatePipelineDto = z.infer<typeof UpdatePipelineDtoSchema>
