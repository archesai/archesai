import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { SubscriptionDtoSchema } from '#billing/subscriptions/dto/subscription.dto'

export const CreateSubscriptionDtoSchema: TObject<{
  planId: TString
}> = Type.Object({
  planId: SubscriptionDtoSchema.properties.planId
})

export type CreateSubscriptionDto = Static<typeof CreateSubscriptionDtoSchema>
