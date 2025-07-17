import { relations } from 'drizzle-orm'
import { pgTable, text } from 'drizzle-orm/pg-core'

import type { PipelineEntity } from '@archesai/schemas'

import { PIPELINE_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'
import { PipelineStepTable } from '#schema/models/pipeline-step'

export const PipelineTable = pgTable(PIPELINE_ENTITY_KEY, {
  ...baseFields,
  description: text(),
  name: text(),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    })
})

export const pipelineRelations = relations(PipelineTable, ({ many, one }) => ({
  organization: one(OrganizationTable, {
    fields: [PipelineTable.organizationId],
    references: [OrganizationTable.id]
  }),
  steps: many(PipelineStepTable)
}))

export type PipelineInsertModel = typeof PipelineTable.$inferInsert
export type PipelineSelectModel = typeof PipelineTable.$inferSelect

export type zPipelineCheck =
  PipelineEntity extends PipelineSelectModel ? true : false
