import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const SubscriptionDtoSchema: TObject<{
  planId: TString
}> = Type.Object({
  planId: Type.String({
    description: 'The ID of the plan'
  })
})

export type SubscriptionDto = Static<typeof SubscriptionDtoSchema>
