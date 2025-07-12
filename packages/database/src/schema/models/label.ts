import { relations } from 'drizzle-orm'
import { pgTable, text, uniqueIndex } from 'drizzle-orm/pg-core'

import { LABEL_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { LabelToArtifactTable } from '#schema/models/label-to-artifact'
import { OrganizationTable } from '#schema/models/organization'

export const LabelTable = pgTable(
  LABEL_ENTITY_KEY,
  {
    ...baseFields,
    name: text('name').notNull(),
    organizationId: text()
      .notNull()
      .references(() => OrganizationTable.id, {
        onDelete: 'cascade',
        onUpdate: 'cascade'
      })
  },
  (LabelTable) => [uniqueIndex().on(LabelTable.name, LabelTable.organizationId)]
)

export const labelRelations = relations(LabelTable, ({ many, one }) => ({
  artifacts: many(LabelToArtifactTable),
  organization: one(OrganizationTable, {
    fields: [LabelTable.organizationId],
    references: [OrganizationTable.id]
  })
}))
