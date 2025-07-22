import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const PaymentMethodEntitySchema: z.ZodObject<{
  billing_details: z.ZodNullable<
    z.ZodObject<{
      address: z.ZodNullable<
        z.ZodObject<{
          city: z.ZodOptional<z.ZodNullable<z.ZodString>>
          country: z.ZodOptional<z.ZodNullable<z.ZodString>>
          line1: z.ZodOptional<z.ZodNullable<z.ZodString>>
          line2: z.ZodOptional<z.ZodNullable<z.ZodString>>
          postal_code: z.ZodOptional<z.ZodNullable<z.ZodString>>
          state: z.ZodOptional<z.ZodNullable<z.ZodString>>
        }>
      >
      email: z.ZodOptional<z.ZodNullable<z.ZodString>>
      name: z.ZodOptional<z.ZodNullable<z.ZodString>>
      phone: z.ZodOptional<z.ZodNullable<z.ZodString>>
    }>
  >
  card: z.ZodOptional<
    z.ZodNullable<
      z.ZodObject<{
        brand: z.ZodString
        country: z.ZodOptional<z.ZodNullable<z.ZodString>>
        exp_month: z.ZodNumber
        exp_year: z.ZodNumber
        fingerprint: z.ZodOptional<z.ZodNullable<z.ZodString>>
        funding: z.ZodString
        last4: z.ZodString
      }>
    >
  >
  createdAt: z.ZodString
  customer: z.ZodString
  id: z.ZodUUID
  type: z.ZodString
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  billing_details: z
    .object({
      address: z
        .object({
          city: z
            .string()
            .nullable()
            .optional()
            .describe('City/District/Suburb/Town/Village.'),
          country: z
            .string()
            .nullable()
            .optional()
            .describe('Two-letter country code (ISO 3166-1 alpha-2).'),
          line1: z
            .string()
            .nullable()
            .optional()
            .describe(
              'Address line 1 (e.g., street, PO Box, or company name).'
            ),
          line2: z
            .string()
            .nullable()
            .optional()
            .describe(
              'Address line 2 (e.g., apartment, suite, unit, or building).'
            ),
          postal_code: z
            .string()
            .nullable()
            .optional()
            .describe('ZIP or postal code.'),
          state: z
            .string()
            .nullable()
            .optional()
            .describe('State/County/Province/Region.')
        })
        .nullable(),
      email: z
        .string()
        .nullable()
        .optional()
        .describe('Email address associated with the payment method.'),
      name: z
        .string()
        .nullable()
        .optional()
        .describe('Full name associated with the payment method.'),
      phone: z
        .string()
        .nullable()
        .optional()
        .describe('Phone number associated with the payment method.')
    })
    .nullable(),
  card: z
    .object({
      brand: z.string().describe('Card brand (e.g., Visa, MasterCard).'),
      country: z
        .string()
        .nullable()
        .optional()
        .describe('Two-letter ISO code representing the country of the card.'),
      exp_month: z
        .number()
        .describe("Two-digit number representing the card's expiration month."),
      exp_year: z
        .number()
        .describe("Four-digit number representing the card's expiration year."),
      fingerprint: z
        .string()
        .nullable()
        .optional()
        .describe('Unencrypted PAN tokens (optional, sensitive).'),
      funding: z
        .string()
        .describe('Card funding type (credit, debit, prepaid, unknown).'),
      last4: z.string().describe('The last four digits of the card.')
    })
    .nullable()
    .optional()
    .describe(
      'If the PaymentMethod is a card, this contains the card details.'
    ),
  customer: z
    .string()
    .describe('ID of the customer this payment method is saved to.'),
  type: z
    .string()
    .describe('The type of the PaymentMethod. An example value is "card".')
})

export type PaymentMethodEntity = z.infer<typeof PaymentMethodEntitySchema>

export const PAYMENT_METHOD_ENTITY_KEY = 'payment-methods'
