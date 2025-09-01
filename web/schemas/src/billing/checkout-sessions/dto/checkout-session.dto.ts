import { z } from 'zod'

export const CheckoutSessionDtoSchema: z.ZodObject<{
  url: z.ZodString
}> = z.object({
  url: z
    .string()
    .describe('The URL that will bring you to the necessary Stripe page')
})

export type CheckoutSessionDto = z.infer<typeof CheckoutSessionDtoSchema>
