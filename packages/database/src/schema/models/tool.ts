import { relations } from 'drizzle-orm'
import { pgTable, text } from 'drizzle-orm/pg-core'

import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import { toolIO } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { organizationFk } from '#schema/models/organization'
import { RunTable } from '#schema/models/run'

// TABLE
export const ToolTable = pgTable(
  TOOL_ENTITY_KEY,
  {
    ...baseFields,
    description: text().notNull(),
    inputType: toolIO().notNull(),
    outputType: toolIO().notNull(),
    toolBase: text().notNull()
  },
  (ToolTable) => [organizationFk(ToolTable)]
)

// RELATIONS
export const toolRelations = relations(ToolTable, ({ many }) => ({
  runs: many(RunTable)
}))

// SCHEMA
export type ToolModel = (typeof ToolTable)['$inferSelect']
