import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const SubscriptionResponseSchema = Type.Object({
  planId: Type.String({
    description: 'The ID of the plan'
  })
})

export type SubscriptionResource = Static<typeof SubscriptionResponseSchema>
