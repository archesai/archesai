import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const SubscriptionDtoSchema = Type.Object({
  planId: Type.String({
    description: 'The ID of the plan'
  })
})

export type SubscriptionDto = Static<typeof SubscriptionDtoSchema>
