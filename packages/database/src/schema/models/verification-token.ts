import { pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import type { VerificationEntity } from '@archesai/schemas'

import { VERIFICATION_TOKEN_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'

export const VerificationTable = pgTable(VERIFICATION_TOKEN_ENTITY_KEY, {
  ...baseFields,
  expiresAt: timestamp({
    mode: 'string'
  }).notNull(),
  identifier: text().notNull(),
  value: text().notNull()
})

export type VerificationInsertModel = typeof VerificationTable.$inferInsert
export type VerificationSelectModel = typeof VerificationTable.$inferSelect

export type zVerificationCheck =
  VerificationEntity extends VerificationSelectModel ? true : false
