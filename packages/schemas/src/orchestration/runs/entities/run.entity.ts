import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { StatusTypes } from '#enums/role'

export const RunEntitySchema: z.ZodObject<{
  completedAt: z.ZodNullable<z.ZodString>
  createdAt: z.ZodString
  error: z.ZodNullable<z.ZodString>
  id: z.ZodString
  organizationId: z.ZodString
  pipelineId: z.ZodNullable<z.ZodString>
  progress: z.ZodNumber
  startedAt: z.ZodNullable<z.ZodString>
  status: z.ZodEnum<{
    COMPLETED: 'COMPLETED'
    FAILED: 'FAILED'
    PROCESSING: 'PROCESSING'
    QUEUED: 'QUEUED'
  }>
  toolId: z.ZodString
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  completedAt: z
    .string()
    .nullable()
    .describe('The timestamp when the run completed'),
  error: z.string().nullable().describe('The error message'),
  organizationId: z.string().describe('The organization name'),
  pipelineId: z
    .string()
    .nullable()
    .describe('The pipeline ID associated with the run'),
  progress: z.number().describe('The percent progress of the run'),
  startedAt: z
    .string()
    .nullable()
    .describe('The timestamp when the run started'),
  status: z.enum(StatusTypes),
  toolId: z.string().describe('The tool ID associated with the run')
}).meta({
  description: 'Schema for Run entity',
  id: 'RunEntity'
})

export type RunEntity = z.infer<typeof RunEntitySchema>

export const RUN_ENTITY_KEY = 'runs'
