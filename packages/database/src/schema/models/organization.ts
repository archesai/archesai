import { integer, pgTable, text } from 'drizzle-orm/pg-core'

import { ORGANIZATION_ENTITY_KEY } from '@archesai/schemas'

import { planType } from '#schema/enums'
import { baseFields } from '#schema/models/base'

export const OrganizationTable = pgTable(ORGANIZATION_ENTITY_KEY, {
  ...baseFields,
  billingEmail: text().notNull(),
  credits: integer().default(0).notNull(),
  logo: text(),
  metadata: text(),
  name: text().notNull(),
  plan: planType().default('FREE').notNull(),
  slug: text().notNull().unique(),
  stripeCustomerId: text().unique()
})

export type OrganizationInsertModel = typeof OrganizationTable.$inferInsert
export type OrganizationSelectModel = typeof OrganizationTable.$inferSelect
