import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CheckoutSessionResponseSchema = Type.Object({
  url: Type.String({
    description: 'The URL that will bring you to the necessary Stripe page'
  })
})

export type CheckoutSessionResponse = Static<
  typeof CheckoutSessionResponseSchema
>
