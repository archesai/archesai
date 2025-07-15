import type { Static, TObject, TString } from '@sinclair/typebox'

import { CreateSubscriptionDtoSchema } from '#billing/subscriptions/dto/create-subscription.dto'

export const UpdateSubscriptionDtoSchema: TObject<{
  planId: TString
}> = CreateSubscriptionDtoSchema

export type UpdateSubscriptionDto = Static<typeof UpdateSubscriptionDtoSchema>
