import { z } from 'zod'

export const CreatePortalDtoSchema: z.ZodObject<{
  organizationId: z.ZodString
}> = z.object({
  organizationId: z.string()
})

export type CreatePortalDto = z.infer<typeof CreatePortalDtoSchema>
