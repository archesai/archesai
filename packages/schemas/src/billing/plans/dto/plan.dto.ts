import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { PlanTypes } from '#enums/role'

export const PLAN_ENTITY_KEY = 'plans'

export const PlanDtoSchema: z.ZodObject<{
  createdAt: z.ZodString
  currency: z.ZodString
  description: z.ZodOptional<z.ZodNullable<z.ZodString>>
  id: z.ZodUUID
  metadata: z.ZodObject<{
    key: z.ZodOptional<
      z.ZodEnum<{
        BASIC: 'BASIC'
        FREE: 'FREE'
        PREMIUM: 'PREMIUM'
        STANDARD: 'STANDARD'
        UNLIMITED: 'UNLIMITED'
      }>
    >
  }>
  name: z.ZodString
  priceId: z.ZodString
  priceMetadata: z.ZodRecord<z.ZodString, z.ZodString>
  recurring: z.ZodOptional<
    z.ZodNullable<
      z.ZodObject<{
        interval: z.ZodString
        interval_count: z.ZodNumber
        trial_period_days: z.ZodOptional<z.ZodNullable<z.ZodNumber>>
      }>
    >
  >
  unitAmount: z.ZodOptional<z.ZodNullable<z.ZodNumber>>
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  currency: z.string().describe('The currency of the plan'),
  description: z
    .string()
    .nullable()
    .optional()
    .describe('The description of the plan'),
  metadata: z.object({
    key: z.enum(PlanTypes).optional().describe('The key of the metadata')
  }),
  name: z.string().describe('The name of the plan'),
  priceId: z.string().describe('The ID of the price associated with the plan'),
  priceMetadata: z
    .record(z.string(), z.string())
    .describe('The metadata of the price associated with the plan'),
  recurring: z
    .object({
      interval: z.string(),
      interval_count: z.number(),
      trial_period_days: z.number().nullable().optional()
    })
    .nullable()
    .optional()
    .describe('The interval of the plan'),
  unitAmount: z
    .number()
    .nullable()
    .optional()
    .describe('The amount in cents to be charged on the interval specified')
})

export type PlanDto = z.infer<typeof PlanDtoSchema>
