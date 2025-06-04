import { relations } from 'drizzle-orm'
import { pgTable, primaryKey, text, uniqueIndex } from 'drizzle-orm/pg-core'

import { PIPELINE_ENTITY_KEY, PIPELINE_STEP_ENTITY_KEY } from '@archesai/domain'

import { baseFields } from '#schema/models/base'
import { organizationFk } from '#schema/models/organization'
import { ToolTable } from '#schema/models/tool'

// TABLE
export const PipelineTable = pgTable(
  PIPELINE_ENTITY_KEY,
  {
    ...baseFields,
    description: text()
  },
  (PipelineTable) => {
    return [organizationFk(PipelineTable)]
  }
)

// RELATIONS
export const pipelineRelations = relations(PipelineTable, ({ many }) => ({
  steps: many(PipelineStepTable)
}))

// TABLE
export const PipelineStepTable = pgTable(
  PIPELINE_STEP_ENTITY_KEY,
  {
    ...baseFields,
    pipelineId: text()
      .notNull()
      .references(() => PipelineTable.id, {
        onDelete: 'cascade',
        onUpdate: 'cascade'
      }),
    toolId: text()
      .notNull()
      .references(() => ToolTable.id, {
        onDelete: 'cascade',
        onUpdate: 'cascade'
      })
  },
  (PipelineStepTable) => [
    uniqueIndex().on(PipelineStepTable.name, PipelineStepTable.pipelineId)
  ]
)

// RELATIONS
export const pipelineStepRelations = relations(
  PipelineStepTable,
  ({ many, one }) => ({
    dependents: many(_PipelineStepDependencies, {
      relationName: 'dependents'
    }),
    pipeline: one(PipelineTable, {
      fields: [PipelineStepTable.pipelineId],
      references: [PipelineTable.id]
    }),
    prerequisites: many(_PipelineStepDependencies, {
      relationName: 'prerequisites'
    }),
    tool: one(ToolTable, {
      fields: [PipelineStepTable.toolId],
      references: [ToolTable.id]
    })
  })
)

// MANY TO MANY
export const _PipelineStepDependencies = pgTable(
  '_pipelineStepDependencies',
  {
    pipelineStepId: text()
      .notNull()
      .references(() => PipelineStepTable.id),
    prerequisiteStepId: text()
      .notNull()
      .references(() => PipelineStepTable.id)
  },
  (_PipelineStepDependencies) => [
    primaryKey({
      columns: [
        _PipelineStepDependencies.pipelineStepId,
        _PipelineStepDependencies.prerequisiteStepId
      ]
    })
  ]
)

export const _pipelineStepDependenciesRelations = relations(
  _PipelineStepDependencies,
  ({ one }) => ({
    pipelineStep: one(PipelineStepTable, {
      fields: [_PipelineStepDependencies.pipelineStepId],
      references: [PipelineStepTable.id],
      relationName: 'dependents'
    }),
    prerequisite: one(PipelineStepTable, {
      fields: [_PipelineStepDependencies.prerequisiteStepId],
      references: [PipelineStepTable.id],
      relationName: 'prerequisites'
    })
  })
)

// SCHEMA
export type PipelineModel = (typeof PipelineTable)['$inferSelect']
export type PipelineStepModel = (typeof PipelineStepTable)['$inferSelect']
