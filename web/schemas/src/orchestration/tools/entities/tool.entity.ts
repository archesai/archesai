import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const ToolEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  description: z.ZodString
  id: z.ZodUUID
  inputMimeType: z.ZodString
  name: z.ZodString
  organizationId: z.ZodString
  outputMimeType: z.ZodString
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  description: z.string().describe('The tool description'),
  inputMimeType: z
    .string()
    .describe('The MIME type of the input for the tool, e.g. text/plain'),
  name: z.string().describe('The name of the tool'),
  organizationId: z.string().describe('The organization name'),
  outputMimeType: z
    .string()
    .describe('The MIME type of the output for the tool, e.g. text/plain')
}).meta({
  description: 'Schema for Tool entity',
  id: 'ToolEntity'
})

export type ToolEntity = z.infer<typeof ToolEntitySchema>

export const TOOL_ENTITY_KEY = 'tools'
