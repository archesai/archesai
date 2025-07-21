import { z } from 'zod'

export const NoContentResponseSchema: z.ZodType<null> = z.null().meta({
  description: 'Schema for 204 No Content response',
  id: 'NoContentResponse'
})

export type NoContentResponse = z.infer<typeof NoContentResponseSchema>
