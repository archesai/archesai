import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const PaymentMethodEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    billing_details: Type.Union([
      Type.Object({
        address: Type.Union([
          Type.Object({
            city: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description: 'City/District/Suburb/Town/Village.'
              })
            ),
            country: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description: 'Two-letter country code (ISO 3166-1 alpha-2).'
              })
            ),
            line1: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description:
                  'Address line 1 (e.g., street, PO Box, or company name).'
              })
            ),
            line2: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description:
                  'Address line 2 (e.g., apartment, suite, unit, or building).'
              })
            ),
            postal_code: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description: 'ZIP or postal code.'
              })
            ),
            state: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description: 'State/County/Province/Region.'
              })
            )
          }),
          Type.Null()
        ]),
        email: Type.Optional(
          Type.Union([Type.String(), Type.Null()], {
            description: 'Email address associated with the payment method.'
          })
        ),
        name: Type.Optional(
          Type.Union([Type.String(), Type.Null()], {
            description: 'Full name associated with the payment method.'
          })
        ),
        phone: Type.Optional(
          Type.Union([Type.String(), Type.Null()], {
            description: 'Phone number associated with the payment method.'
          })
        )
      }),

      Type.Null()
    ]),
    card: Type.Optional(
      Type.Union(
        [
          Type.Object({
            brand: Type.String({
              description: 'Card brand (e.g., Visa, MasterCard).'
            }),
            country: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description:
                  'Two-letter ISO code representing the country of the card.'
              })
            ),
            exp_month: Type.Number({
              description:
                'Two-digit number representing the card’s expiration month.'
            }),
            exp_year: Type.Number({
              description:
                'Four-digit number representing the card’s expiration year.'
            }),
            fingerprint: Type.Optional(
              Type.Union([Type.String(), Type.Null()], {
                description: 'Unencrypted PAN tokens (optional, sensitive).'
              })
            ),
            funding: Type.String({
              description:
                'Card funding type (credit, debit, prepaid, unknown).'
            }),
            last4: Type.String({
              description: 'The last four digits of the card.'
            })
          }),
          Type.Null()
        ],
        {
          description:
            'If the PaymentMethod is a card, this contains the card details.'
        }
      )
    ),
    customer: Type.String({
      description: 'ID of the customer this payment method is saved to.'
    }),
    id: Type.String({
      description: 'Unique identifier for the payment method.'
    }),
    type: Type.String({
      description: 'The type of the PaymentMethod. An example value is "card".'
    })
  },
  {
    $id: 'PaymentMethodEntity',
    description: 'The payment method entity',
    title: 'Payment Method Entity'
  }
)

export type PaymentMethodEntity = Static<typeof PaymentMethodEntitySchema>

export const PAYMENT_METHOD_ENTITY_KEY = 'payment-methods'
