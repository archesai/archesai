import type { z } from 'zod'

import { ToolEntitySchema } from '#orchestration/tools/entities/tool.entity'

export const CreateToolDtoSchema: z.ZodObject<{
  description: z.ZodString
  name: z.ZodString
}> = ToolEntitySchema.pick({
  description: true,
  name: true
})

export type CreateToolDto = z.infer<typeof CreateToolDtoSchema>
