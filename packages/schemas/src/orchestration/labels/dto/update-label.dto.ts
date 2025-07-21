import type { z } from 'zod'

import { CreateLabelDtoSchema } from '#orchestration/labels/dto/create-label.dto'

export const UpdateLabelDtoSchema: z.ZodObject<{
  name: z.ZodOptional<z.ZodString>
}> = CreateLabelDtoSchema.partial()

export type UpdateLabelDto = z.infer<typeof UpdateLabelDtoSchema>
