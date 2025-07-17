import { relations } from 'drizzle-orm'
import {
  boolean,
  integer,
  jsonb,
  pgTable,
  text,
  timestamp
} from 'drizzle-orm/pg-core'

import type { ApiTokenEntity } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'

const API_TOKEN_ENTITY_KEY = 'apiToken'

export const ApiTokenTable = pgTable(API_TOKEN_ENTITY_KEY, {
  ...baseFields,
  enabled: boolean().notNull(),
  expiresAt: timestamp({
    mode: 'string'
  }),
  key: text().notNull(),
  lastRefill: timestamp({
    mode: 'string'
  }),
  lastRequest: timestamp({
    mode: 'string'
  }),
  metadata: jsonb(),
  name: text(),
  permissions: text(),
  prefix: text(),
  rateLimitEnabled: boolean().notNull(),
  rateLimitMax: integer(),
  rateLimitTimeWindow: integer(),
  refillAmount: integer(),
  refillInterval: integer(),
  remaining: integer(),
  requestCount: integer().notNull().default(0),
  start: text(),
  userId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    })
})

export const apiTokenRelations = relations(ApiTokenTable, ({ one }) => ({
  user: one(OrganizationTable, {
    fields: [ApiTokenTable.userId],
    references: [OrganizationTable.id]
  })
}))

export type ApiTokenInsertModel = typeof ApiTokenTable.$inferInsert
export type ApiTokenSelectModel = typeof ApiTokenTable.$inferSelect

export type zApiTokenCheck =
  ApiTokenEntity extends ApiTokenSelectModel ? true : false
