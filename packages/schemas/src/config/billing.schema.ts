import { z } from 'zod'

const BaseBillingConfigSchema = z.object({
  stripe: z.object({
    token: z
      .string()
      .describe(
        'Stripe secret API key (sk_live_... or sk_test_...) for payment processing'
      ),
    whsec: z
      .string()
      .describe(
        'Stripe webhook endpoint secret for verifying webhook signatures'
      )
  })
})

export const BillingConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          mode: z.ZodLiteral<'disabled'>
        }>,
        z.ZodObject<{
          mode: z.ZodLiteral<'enabled'>
          stripe: z.ZodObject<{
            token: z.ZodString
            whsec: z.ZodString
          }>
        }>
      ]
    >
  >
> = z
  .discriminatedUnion('mode', [
    z.object({
      mode: z.literal('disabled')
    }),
    BaseBillingConfigSchema.extend({
      mode: z.literal('enabled')
    })
  ])
  .optional()
  .default({ mode: 'disabled' })
  .describe(
    'Billing configuration for payment processing using Stripe. Includes API keys and webhook secrets.'
  )

export type BillingConfig = z.infer<typeof BillingConfigSchema>
