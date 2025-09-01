import { z } from 'zod'

export const IdParamsSchema: z.ZodObject<{
  id: z.ZodUUID
}> = z
  .object({
    id: z.uuid().describe('The unique identifier of the resource.')
  })
  .describe('ID Params')
