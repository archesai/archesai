import { relations } from 'drizzle-orm'
import { integer, pgTable, text } from 'drizzle-orm/pg-core'

import type { ArtifactEntity } from '@archesai/schemas'

import { ARTIFACT_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { LabelToArtifactTable } from '#schema/models/label-to-artifact'
import { OrganizationTable } from '#schema/models/organization'
import { RunTable } from '#schema/models/run'
import { RunToArtifactTable } from '#schema/models/run-to-artifact'

export const ArtifactTable = pgTable(ARTIFACT_ENTITY_KEY, {
  ...baseFields,
  credits: integer().default(0).notNull(),
  description: text(),
  // embedding: vector({
  //   dimensions: 1536
  // }),
  mimeType: text().default('application/octet-stream').notNull(),
  name: text(),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
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

export const artifactRelations = relations(ArtifactTable, ({ many, one }) => ({
  consumers: many(RunToArtifactTable),
  labels: many(LabelToArtifactTable),
  producer: one(RunTable, {
    fields: [ArtifactTable.producerId],
    references: [RunTable.id]
  })
}))

export type ArtifactInsertModel = typeof ArtifactTable.$inferInsert
export type ArtifactSelectModel = typeof ArtifactTable.$inferSelect

export type zArtifactCheck =
  ArtifactEntity extends ArtifactSelectModel ? true : false
