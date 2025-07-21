import { z } from 'zod'

export const SubscriptionDtoSchema: z.ZodObject<{
  planId: z.ZodString
}> = z.object({
  planId: z.string().describe('The ID of the plan')
})

export type SubscriptionDto = z.infer<typeof SubscriptionDtoSchema>
