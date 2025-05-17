import type { PgColumn } from 'drizzle-orm/pg-core'

import { foreignKey, integer, pgTable, text } from 'drizzle-orm/pg-core'

import { ORGANIZATION_ENTITY_KEY } from '@archesai/domain'

import { planType } from '#schema/enums'
import { baseFields } from '#schema/models/base'

// TABLE
export const OrganizationTable = pgTable(ORGANIZATION_ENTITY_KEY, {
  ...baseFields,
  billingEmail: text().notNull(),
  credits: integer().default(0).notNull(),
  orgname: text('orgname').notNull().unique(),
  plan: planType().default('FREE').notNull(),
  stripeCustomerId: text().unique()
})

// FOREIGN KEY HELPER
export const organizationFk = (table: { orgname: PgColumn }) =>
  foreignKey({
    columns: [table.orgname],
    foreignColumns: [OrganizationTable.id]
  })
    .onUpdate('cascade')
    .onDelete('cascade')

// SCHEMA
export type OrganizationModel = (typeof OrganizationTable)['$inferSelect']
