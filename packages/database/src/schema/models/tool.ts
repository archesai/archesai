import { relations } from 'drizzle-orm'
import { pgTable, text } from 'drizzle-orm/pg-core'

import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import { toolIO } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'
import { RunTable } from '#schema/models/run'

export const ToolTable = pgTable(TOOL_ENTITY_KEY, {
  ...baseFields,
  description: text().notNull(),
  inputType: toolIO().notNull(),
  name: text().notNull(),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  outputType: toolIO().notNull(),
  toolBase: text().notNull()
})

export type ToolInsertModel = typeof ToolTable.$inferInsert
export type ToolSelectModel = typeof ToolTable.$inferSelect

export const toolRelations = relations(ToolTable, ({ many, one }) => ({
  organization: one(OrganizationTable, {
    fields: [ToolTable.organizationId],
    references: [OrganizationTable.id]
  }),
  runs: many(RunTable)
}))
