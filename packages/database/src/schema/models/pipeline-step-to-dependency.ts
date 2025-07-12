import { relations } from 'drizzle-orm'
import { pgTable, primaryKey, text } from 'drizzle-orm/pg-core'

import { PipelineStepTable } from '#schema/models/pipeline-step'

export const PipelineStepToDependency = pgTable(
  'pipelineStepToDependency',
  {
    pipelineStepId: text()
      .notNull()
      .references(() => PipelineStepTable.id),
    prerequisiteId: text()
      .notNull()
      .references(() => PipelineStepTable.id)
  },
  (PipelineStepToDependency) => [
    primaryKey({
      columns: [
        PipelineStepToDependency.pipelineStepId,
        PipelineStepToDependency.prerequisiteId
      ]
    })
  ]
)

export const pipelineStepToDependencyRelations = relations(
  PipelineStepToDependency,
  ({ one }) => ({
    pipelineStep: one(PipelineStepTable, {
      fields: [PipelineStepToDependency.pipelineStepId],
      references: [PipelineStepTable.id],
      relationName: 'dependents'
    }),
    prerequisite: one(PipelineStepTable, {
      fields: [PipelineStepToDependency.prerequisiteId],
      references: [PipelineStepTable.id],
      relationName: 'prerequisites'
    })
  })
)
