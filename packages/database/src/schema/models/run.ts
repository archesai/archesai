import { relations } from 'drizzle-orm'
import { doublePrecision, pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import type { RunEntity } from '@archesai/schemas'

import { RUN_ENTITY_KEY } from '@archesai/schemas'

import { runStatus } from '#schema/enums'
import { ArtifactTable } from '#schema/models/artifact'
import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'
import { PipelineTable } from '#schema/models/pipeline'
import { RunToArtifactTable } from '#schema/models/run-to-artifact'
import { ToolTable } from '#schema/models/tool'

export const RunTable = pgTable(RUN_ENTITY_KEY, {
  ...baseFields,
  completedAt: timestamp({ mode: 'string', precision: 3 }),
  error: text(),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  pipelineId: text().references(() => PipelineTable.id, {
    onDelete: 'set null',
    onUpdate: 'cascade'
  }),
  progress: doublePrecision().default(0).notNull(),
  startedAt: timestamp({ mode: 'string', precision: 3 }),
  status: runStatus().default('QUEUED').notNull(),
  toolId: text()
    .notNull()
    .references(() => ToolTable.id, {
      onDelete: 'set null',
      onUpdate: 'cascade'
    })
})

export const runRelations = relations(RunTable, ({ many, one }) => ({
  inputs: many(RunToArtifactTable),
  outputs: many(ArtifactTable),
  pipeline: one(PipelineTable, {
    fields: [RunTable.pipelineId],
    references: [PipelineTable.id]
  }),
  tool: one(ToolTable, {
    fields: [RunTable.toolId],
    references: [ToolTable.id]
  })
}))

export type RunInsertModel = typeof RunTable.$inferInsert
export type RunSelectModel = typeof RunTable.$inferSelect

export type zRunCheck = RunEntity extends RunSelectModel ? true : false
