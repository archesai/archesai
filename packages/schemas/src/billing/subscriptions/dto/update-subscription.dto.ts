import type { z } from 'zod'

import { CreateSubscriptionDtoSchema } from '#billing/subscriptions/dto/create-subscription.dto'

export const UpdateSubscriptionDtoSchema: z.ZodObject<{
  planId: z.ZodString
}> = CreateSubscriptionDtoSchema

export type UpdateSubscriptionDto = z.infer<typeof UpdateSubscriptionDtoSchema>
