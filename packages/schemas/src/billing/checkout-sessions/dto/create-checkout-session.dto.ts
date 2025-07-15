import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreateCheckoutSessionDtoSchema: TObject<{
  priceId: TString
}> = Type.Object({
  priceId: Type.String({
    description: 'The ID of the price associated with the checkout session',
    maxLength: 255,
    minLength: 1
  })
})

export type CreateCheckoutSessionDto = Static<
  typeof CreateCheckoutSessionDtoSchema
>
