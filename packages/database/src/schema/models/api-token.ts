import { relations } from 'drizzle-orm'
import { pgTable, text } from 'drizzle-orm/pg-core'

import { API_TOKEN_ENTITY_KEY } from '@archesai/schemas'

import { roleEnum } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'

export const ApiTokenTable = pgTable(API_TOKEN_ENTITY_KEY, {
  ...baseFields,
  key: text().notNull(),
  name: text(),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  role: roleEnum().default('USER').notNull()
})

export type ApiTokenInsertModel = typeof ApiTokenTable.$inferInsert
export type ApiTokenSelectModel = typeof ApiTokenTable.$inferSelect

export const apiTokenRelations = relations(ApiTokenTable, ({ one }) => ({
  organization: one(OrganizationTable, {
    fields: [ApiTokenTable.organizationId],
    references: [OrganizationTable.id]
  })
}))
