import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { SubscriptionResponseSchema } from '#subscriptions/dto/subscription.res.dto'

export const CreateSubscriptionRequestSchema = Type.Object({
  planId: SubscriptionResponseSchema.properties.planId
})

export type CreateSubscriptionRequest = Static<
  typeof CreateSubscriptionRequestSchema
>
