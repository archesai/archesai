import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const PipelineStepEntitySchema: z.ZodObject<{
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
}> = BaseEntitySchema.extend({
  dependents: z.array(
    z.object({
      pipelineStepId: z.string()
    })
  ),
  pipelineId: z.string(),
  prerequisites: z.array(
    z.object({
      pipelineStepId: z.string()
    })
  ),
  toolId: z.string()
}).meta({
  description: 'Schema for Pipeline Step entity',
  id: 'PipelineStepEntity'
})

export const PipelineEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  description: z.ZodNullable<z.ZodString>
  id: z.ZodString
  name: z.ZodNullable<z.ZodString>
  organizationId: z.ZodString
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
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  description: z.string().nullable().describe('The pipeline description'),
  name: z.string().nullable().describe('The pipeline name'),
  organizationId: z.string().describe('The organization id'),
  steps: z.array(PipelineStepEntitySchema).describe('The steps in the pipeline')
}).meta({
  description: 'Schema for Pipeline entity',
  id: 'PipelineEntity'
})

export type PipelineEntity = z.infer<typeof PipelineEntitySchema>

export type PipelineStepEntity = z.infer<typeof PipelineStepEntitySchema>

export const PIPELINE_ENTITY_KEY = 'pipelines'

export const PIPELINE_STEP_ENTITY_KEY = 'pipeline-steps'
