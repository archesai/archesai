import { relations } from 'drizzle-orm'
import { pgTable, primaryKey, text, uniqueIndex } from 'drizzle-orm/pg-core'

import { LABEL_ENTITY_KEY } from '@archesai/domain'

import { ArtifactTable } from '#schema/models/artifact'
import { baseFields } from '#schema/models/base'
import { organizationFk } from '#schema/models/organization'

// TABLE
export const LabelTable = pgTable(
  LABEL_ENTITY_KEY,
  {
    ...baseFields
  },
  (LabelTable) => [
    organizationFk(LabelTable),
    uniqueIndex().on(LabelTable.name, LabelTable.orgname)
  ]
)

// RELATIONS
export const labelRelations = relations(LabelTable, ({ many }) => ({
  content: many(_LabelsToContent)
}))

// MANY TO MANY
export const _LabelsToContent = pgTable(
  '_labelsToContent',
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

export const _labelsToContentRelations = relations(
  _LabelsToContent,
  ({ one }) => ({
    content: one(ArtifactTable, {
      fields: [_LabelsToContent.artifactId],
      references: [ArtifactTable.id]
    }),
    label: one(LabelTable, {
      fields: [_LabelsToContent.labelId],
      references: [LabelTable.id],
      relationName: 'artifact'
    })
  })
)

// SCHEMA
export type LabelModel = (typeof LabelTable)['$inferSelect']
