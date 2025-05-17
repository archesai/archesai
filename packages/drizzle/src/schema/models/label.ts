import { relations } from 'drizzle-orm'
import { pgTable, primaryKey, text, uniqueIndex } from 'drizzle-orm/pg-core'

import { LABEL_ENTITY_KEY } from '@archesai/domain'

import { baseFields } from '#schema/models/base'
import { ContentTable } from '#schema/models/content'
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
    contentId: text()
      .notNull()
      .references(() => ContentTable.id),
    labelId: text()
      .notNull()
      .references(() => LabelTable.id)
  },
  (t) => [primaryKey({ columns: [t.labelId, t.contentId] })]
)

export const _labelsToContentRelations = relations(
  _LabelsToContent,
  ({ one }) => ({
    content: one(ContentTable, {
      fields: [_LabelsToContent.contentId],
      references: [ContentTable.id]
    }),
    label: one(LabelTable, {
      fields: [_LabelsToContent.labelId],
      references: [LabelTable.id],
      relationName: 'content'
    })
  })
)

// SCHEMA
export type LabelModel = (typeof LabelTable)['$inferSelect']
