import type { z } from 'zod'

import { LabelEntitySchema } from '#orchestration/labels/entities/label.entity'

export const CreateLabelDtoSchema: z.ZodObject<{
  name: z.ZodString
}> = LabelEntitySchema.pick({
  name: true
})

export type CreateLabelDto = z.infer<typeof CreateLabelDtoSchema>
