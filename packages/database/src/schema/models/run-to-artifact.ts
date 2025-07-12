import { relations } from 'drizzle-orm'
import { pgTable, primaryKey, text } from 'drizzle-orm/pg-core'

import { ArtifactTable } from '#schema/models/artifact'
import { RunTable } from '#schema/models/run'

export const RunToArtifactTable = pgTable(
  'runToArtifact',
  {
    artifactId: text()
      .notNull()
      .references(() => ArtifactTable.id, {
        onDelete: 'cascade'
      }),
    runId: text()
      .notNull()
      .references(() => RunTable.id, {
        onDelete: 'cascade'
      })
  },
  (t) => [primaryKey({ columns: [t.runId, t.artifactId] })]
)

export const runToArtifactRelations = relations(
  RunToArtifactTable,
  ({ one }) => ({
    consumer: one(RunTable, {
      fields: [RunToArtifactTable.runId],
      references: [RunTable.id]
    }),
    input: one(ArtifactTable, {
      fields: [RunToArtifactTable.artifactId],
      references: [ArtifactTable.id]
    })
  })
)
