import { pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import { VERIFICATION_TOKEN_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'

export const VerificationTokenTable = pgTable(VERIFICATION_TOKEN_ENTITY_KEY, {
  ...baseFields,
  expiresAt: timestamp({
    mode: 'string'
  }).notNull(),
  identifier: text().notNull(),
  value: text().notNull()
})

export type VerificationTokenInsertModel =
  typeof VerificationTokenTable.$inferInsert
export type VerificationTokenSelectModel =
  typeof VerificationTokenTable.$inferSelect
