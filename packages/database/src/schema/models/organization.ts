import { integer, pgTable, text } from 'drizzle-orm/pg-core'

import type { OrganizationEntity } from '@archesai/schemas'

import { ORGANIZATION_ENTITY_KEY } from '@archesai/schemas'

import { baseFields, planEnum } from '#schema/models/base'

export const OrganizationTable = pgTable(ORGANIZATION_ENTITY_KEY, {
  ...baseFields,
  billingEmail: text(),
  credits: integer().default(0).notNull(),
  logo: text(),
  metadata: text(),
  name: text().notNull(),
  plan: planEnum().default('FREE').notNull(),
  slug: text().notNull().unique(),
  stripeCustomerId: text().unique()
})

export type OrganizationInsertModel = typeof OrganizationTable.$inferInsert
export type OrganizationSelectModel = typeof OrganizationTable.$inferSelect

export type zOrganizationCheck =
  OrganizationEntity extends OrganizationSelectModel ? true : false
