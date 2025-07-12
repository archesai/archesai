import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { SubscriptionDtoSchema } from '#billing/subscriptions/dto/subscription.dto'

export const CreateSubscriptionDtoSchema = Type.Object({
  planId: SubscriptionDtoSchema.properties.planId
})

export type CreateSubscriptionDto = Static<typeof CreateSubscriptionDtoSchema>
