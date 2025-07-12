import { relations } from 'drizzle-orm'
import { pgTable, primaryKey, text } from 'drizzle-orm/pg-core'

import { ArtifactTable } from '#schema/models/artifact'
import { LabelTable } from '#schema/models/label'

export const LabelToArtifactTable = pgTable(
  'labelToArtifact',
  {
    artifactId: text()
      .notNull()
      .references(() => ArtifactTable.id),
    labelId: text()
      .notNull()
      .references(() => LabelTable.id)
  },
  (t) => [primaryKey({ columns: [t.labelId, t.artifactId] })]
)

export const labelToArtifactRelations = relations(
  LabelToArtifactTable,
  ({ one }) => ({
    artifact: one(ArtifactTable, {
      fields: [LabelToArtifactTable.artifactId],
      references: [ArtifactTable.id]
    }),
    label: one(LabelTable, {
      fields: [LabelToArtifactTable.labelId],
      references: [LabelTable.id]
    })
  })
)
