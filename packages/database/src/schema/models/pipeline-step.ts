import { relations } from 'drizzle-orm'
import { pgTable, text } from 'drizzle-orm/pg-core'

import { PIPELINE_STEP_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { PipelineTable } from '#schema/models/pipeline'
import { PipelineStepToDependency } from '#schema/models/pipeline-step-to-dependency'
import { ToolTable } from '#schema/models/tool'

export const PipelineStepTable = pgTable(PIPELINE_STEP_ENTITY_KEY, {
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
})

export const pipelineStepRelations = relations(
  PipelineStepTable,
  ({ many, one }) => ({
    dependents: many(PipelineStepToDependency, {
      relationName: 'dependents'
    }),
    pipeline: one(PipelineTable, {
      fields: [PipelineStepTable.pipelineId],
      references: [PipelineTable.id]
    }),
    prerequisites: many(PipelineStepToDependency, {
      relationName: 'prerequisites'
    }),
    tool: one(ToolTable, {
      fields: [PipelineStepTable.toolId],
      references: [ToolTable.id]
    })
  })
)
