import type { z } from 'zod'

import { PipelineEntitySchema } from '#orchestration/pipelines/entities/pipeline.entity'

export const CreatePipelineDtoSchema: z.ZodObject<{
  description: z.ZodNullable<z.ZodString>
  name: z.ZodNullable<z.ZodString>
  steps: z.ZodArray<
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
}> = PipelineEntitySchema.pick({
  description: true,
  name: true,
  steps: true
})

export type CreatePipelineDto = z.infer<typeof CreatePipelineDtoSchema>
