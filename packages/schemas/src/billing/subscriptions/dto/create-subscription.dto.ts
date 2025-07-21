import type { z } from 'zod'

import { SubscriptionDtoSchema } from '#billing/subscriptions/dto/subscription.dto'

export const CreateSubscriptionDtoSchema: z.ZodObject<{
  planId: z.ZodString
}> = SubscriptionDtoSchema.pick({
  planId: true
})

export type CreateSubscriptionDto = z.infer<typeof CreateSubscriptionDtoSchema>
