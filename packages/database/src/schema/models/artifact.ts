import type { AnyPgColumn } from 'drizzle-orm/pg-core'

import { relations } from 'drizzle-orm'
import { integer, pgTable, primaryKey, text, vector } from 'drizzle-orm/pg-core'

import { ARTIFACT_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { LabelToArtifactTable } from '#schema/models/label-to-artifact'
import { OrganizationTable } from '#schema/models/organization'
import { _RunToArtifactTable, RunTable } from '#schema/models/run'

export const ArtifactTable = pgTable(ARTIFACT_ENTITY_KEY, {
  ...baseFields,
  credits: integer().default(0).notNull(),
  description: text(),
  embedding: vector({
    dimensions: 1536
  }),
  mimeType: text(),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  parentId: text().references((): AnyPgColumn => ArtifactTable.id, {
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
})

export const contentRelations = relations(ArtifactTable, ({ many, one }) => ({
  // children: many(_ParentToChild, {
  //   relationName: 'children'
  // }),
  consumers: many(_RunToArtifactTable),

  labels: many(LabelToArtifactTable),
  parent: one(_ParentToChild, {
    fields: [ArtifactTable.id],
    references: [_ParentToChild.parentId],
    relationName: 'parent'
  }),
  producer: one(RunTable, {
    fields: [ArtifactTable.producerId],
    references: [RunTable.id]
  })
}))

export const _ParentToChild = pgTable(
  '_parentToChild',
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

export const _parentToChildRelations = relations(_ParentToChild, ({ one }) => ({
  child: one(ArtifactTable, {
    fields: [_ParentToChild.childId],
    references: [ArtifactTable.id]
  }),
  parent: one(ArtifactTable, {
    fields: [_ParentToChild.parentId],
    references: [ArtifactTable.id]
  })
}))
