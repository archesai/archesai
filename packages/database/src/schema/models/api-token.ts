import { pgTable, text } from 'drizzle-orm/pg-core'

import { API_TOKEN_ENTITY_KEY } from '@archesai/schemas'

import { roleEnum } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { organizationFk } from '#schema/models/organization'

// TABLE
export const ApiTokenTable = pgTable(
  API_TOKEN_ENTITY_KEY,
  {
    ...baseFields,
    key: text().notNull(),
    role: roleEnum().default('USER').notNull()
  },
  (apiTokenTable) => [organizationFk(apiTokenTable)]
)

export type ApiTokenInsert = (typeof ApiTokenTable)['$inferInsert']
// SCHEMA
export type ApiTokenModel = (typeof ApiTokenTable)['$inferSelect']
