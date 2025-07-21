import type { z } from 'zod'

import { CreateToolDtoSchema } from '#orchestration/tools/dto/create-tool.dto'

export const UpdateToolDtoSchema: z.ZodObject<{
  description: z.ZodOptional<z.ZodString>
  name: z.ZodOptional<z.ZodString>
}> = CreateToolDtoSchema.partial()

export type UpdateToolDto = z.infer<typeof UpdateToolDtoSchema>
