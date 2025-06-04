import { relations } from 'drizzle-orm'
import {
  doublePrecision,
  pgTable,
  primaryKey,
  text,
  timestamp
} from 'drizzle-orm/pg-core'

import { RUN_ENTITY_KEY } from '@archesai/domain'

import { runStatus, runType } from '#schema/enums'
import { ArtifactTable } from '#schema/models/artifact'
import { baseFields } from '#schema/models/base'
import { organizationFk } from '#schema/models/organization'
import {
  PipelineTable
  // PipelineStepTable,
} from '#schema/models/pipeline'
import { ToolTable } from '#schema/models/tool'

// TABLE
export const RunTable = pgTable(
  RUN_ENTITY_KEY,
  {
    ...baseFields,
    completedAt: timestamp({ mode: 'string', precision: 3 }),
    error: text(),
    pipelineId: text().references(() => PipelineTable.id, {
      onDelete: 'set null',
      onUpdate: 'cascade'
    }),
    progress: doublePrecision().default(0).notNull(),
    runType: runType().notNull(),
    startedAt: timestamp({ mode: 'string', precision: 3 }),
    status: runStatus().default('QUEUED').notNull(),
    toolId: text().references(() => ToolTable.id, {
      onDelete: 'set null',
      onUpdate: 'cascade'
    })
    // pipelineRunId: text().notNull(),
    // pipelineStepId: text().notNull()
  },
  (RunTable) => [
    organizationFk(RunTable)
    // uniqueIndex().on(RunTable.pipelineRunId, RunTable.pipelineStepId)
  ]
)

// RELATIONS
export const runRelations = relations(RunTable, ({ many, one }) => ({
  inputs: many(_RunToArtifactTable),
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

// MANY TO MANY
export const _RunToArtifactTable = pgTable(
  '_runToContent',
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

export const _runToContentRelations = relations(
  _RunToArtifactTable,
  ({ one }) => ({
    consumer: one(RunTable, {
      fields: [_RunToArtifactTable.runId],
      references: [RunTable.id]
    }),
    input: one(ArtifactTable, {
      fields: [_RunToArtifactTable.artifactId],
      references: [ArtifactTable.id]
    })
  })
)

// SCHEMA
export type RunModel = (typeof RunTable)['$inferSelect']
