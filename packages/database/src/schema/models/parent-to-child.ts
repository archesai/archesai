import { relations } from 'drizzle-orm'
import { pgTable, primaryKey, text } from 'drizzle-orm/pg-core'

import { ArtifactTable } from '#schema/models/artifact'

export const ParentToChildTable = pgTable(
  'parentToChild',
  {
    childId: text()
      .notNull()
      .references(() => ArtifactTable.id, {
        onDelete: 'cascade'
      }),
    parentId: text()
      .notNull()
      .references(() => ArtifactTable.id, {
        onDelete: 'cascade'
      })
  },
  (t) => [primaryKey({ columns: [t.parentId, t.childId] })]
)

export const parentToChildRelations = relations(
  ParentToChildTable,
  ({ one }) => ({
    child: one(ArtifactTable, {
      fields: [ParentToChildTable.childId],
      references: [ArtifactTable.id]
    }),
    parent: one(ArtifactTable, {
      fields: [ParentToChildTable.parentId],
      references: [ArtifactTable.id]
    })
  })
)
