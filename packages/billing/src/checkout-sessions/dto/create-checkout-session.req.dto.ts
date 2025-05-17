import type { StaticDecode } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreateCheckoutSessionRequestSchema = Type.Object({
  priceId: Type.String({
    description: 'The ID of the price associated with the checkout session',
    maxLength: 255,
    minLength: 1
  })
})

export type CreateCheckoutSessionRequest = StaticDecode<
  typeof CreateCheckoutSessionRequestSchema
>
