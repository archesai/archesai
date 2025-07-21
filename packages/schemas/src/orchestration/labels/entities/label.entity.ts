import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const LabelEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  id: z.ZodString
  name: z.ZodString
  organizationId: z.ZodString
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  name: z.string().describe('The name of the label'),
  organizationId: z.string().describe('The organization name')
}).meta({
  description: 'Schema for Label entity',
  id: 'LabelEntity'
})

export type LabelEntity = z.infer<typeof LabelEntitySchema>

export const LABEL_ENTITY_KEY = 'labels'
