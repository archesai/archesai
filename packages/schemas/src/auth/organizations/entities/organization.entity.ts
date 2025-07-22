import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { PlanTypes } from '#enums/role'

export const OrganizationEntitySchema: z.ZodObject<{
  billingEmail: z.ZodNullable<z.ZodString>
  createdAt: z.ZodString
  credits: z.ZodNumber
  id: z.ZodUUID
  logo: z.ZodNullable<z.ZodString>
  metadata: z.ZodNullable<z.ZodString>
  name: z.ZodString
  plan: z.ZodEnum<{
    BASIC: 'BASIC'
    FREE: 'FREE'
    PREMIUM: 'PREMIUM'
    STANDARD: 'STANDARD'
    UNLIMITED: 'UNLIMITED'
  }>
  slug: z.ZodString
  stripeCustomerId: z.ZodNullable<z.ZodString>
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  billingEmail: z
    .string()
    .nullable()
    .describe('The billing email to use for the organization'),
  credits: z
    .number()
    .describe('The number of credits you have remaining for this organization'),
  logo: z.string().nullable().describe('The URL of the organization logo'),
  metadata: z
    .string()
    .nullable()
    .describe('The metadata for the organization, used for custom data'),
  name: z.string().describe('The name of the organization'),
  plan: z
    .enum(PlanTypes)
    .describe('The plan that the organization is subscribed to'),
  slug: z
    .string()
    .describe('The unique slug for the organization, used in URLs'),
  stripeCustomerId: z.string().nullable().describe('The Stripe customer ID')
}).meta({
  description: 'Schema for Organization entity',
  id: 'OrganizationEntity'
})

export type OrganizationEntity = z.infer<typeof OrganizationEntitySchema>

export const ORGANIZATION_ENTITY_KEY = 'organizations'
