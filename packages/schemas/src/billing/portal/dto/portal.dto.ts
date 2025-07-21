import { z } from 'zod'

export const PortalDtoSchema: z.ZodObject<{
  url: z.ZodString
}> = z.object({
  url: z
    .string()
    .describe('The URL that will bring you to the necessary Stripe page')
})

export type PortalDto = z.infer<typeof PortalDtoSchema>
