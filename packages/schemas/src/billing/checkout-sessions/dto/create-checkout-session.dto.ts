import { z } from 'zod'

export const CreateCheckoutSessionDtoSchema: z.ZodObject<{
  priceId: z.ZodString
}> = z.object({
  priceId: z
    .string()
    .min(1)
    .max(255)
    .describe('The ID of the price associated with the checkout session')
})

export type CreateCheckoutSessionDto = z.infer<
  typeof CreateCheckoutSessionDtoSchema
>
