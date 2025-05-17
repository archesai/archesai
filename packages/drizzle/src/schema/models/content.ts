import type { AnyPgColumn } from 'drizzle-orm/pg-core'

import { relations } from 'drizzle-orm'
import { integer, pgTable, primaryKey, text, vector } from 'drizzle-orm/pg-core'

import { CONTENT_ENTITY_KEY } from '@archesai/domain'

import { baseFields } from '#schema/models/base'
import { _LabelsToContent } from '#schema/models/label'
import { organizationFk } from '#schema/models/organization'
import { _RunToContentTable, RunTable } from '#schema/models/run'

// TABLE
export const ContentTable = pgTable(
  CONTENT_ENTITY_KEY,
  {
    ...baseFields,
    credits: integer().default(0).notNull(),
    description: text(),
    embedding: vector({
      dimensions: 1536
    }),
    mimeType: text(),
    parentId: text().references((): AnyPgColumn => ContentTable.id, {
      onDelete: 'set null',
      onUpdate: 'cascade'
    }),
    previewImage: text(),
    producerId: text().references(() => RunTable.id, {
      onDelete: 'set null',
      onUpdate: 'cascade'
    }),
    text: text(),
    url: text()
  },
  (ContentTable) => [organizationFk(ContentTable)]
)

// RELATIONS
export const contentRelations = relations(ContentTable, ({ many, one }) => ({
  children: many(_ParentToChild, {
    relationName: 'children'
  }),
  consumers: many(_RunToContentTable),

  labels: many(_LabelsToContent),
  parent: one(_ParentToChild, {
    fields: [ContentTable.id],
    references: [_ParentToChild.parentId],
    relationName: 'parent'
  }),
  producer: one(RunTable, {
    fields: [ContentTable.producerId],
    references: [RunTable.id]
  })
}))

// SCHEMAS
export type ContentModel = (typeof ContentTable)['$inferSelect']

// MANY TO MANY
export const _ParentToChild = pgTable(
  '_parentToChild',
  {
    childId: text()
      .notNull()
      .references(() => ContentTable.id, {
        onDelete: 'cascade'
      }),
    parentId: text()
      .notNull()
      .references(() => ContentTable.id, {
        onDelete: 'cascade'
      })
  },
  (t) => [primaryKey({ columns: [t.parentId, t.childId] })]
)

export const _parentToChildRelations = relations(_ParentToChild, ({ one }) => ({
  child: one(ContentTable, {
    fields: [_ParentToChild.childId],
    references: [ContentTable.id]
  }),
  parent: one(ContentTable, {
    fields: [_ParentToChild.parentId],
    references: [ContentTable.id]
  })
}))
