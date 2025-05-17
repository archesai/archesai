import { sql } from 'drizzle-orm'
import { pgTable, primaryKey, text, timestamp } from 'drizzle-orm/pg-core'

import { VERIFICATION_TOKEN_ENTITY_KEY } from '@archesai/domain'

import { verificationTokenType } from '#schema/enums'

export const VerificationTokenTable = pgTable(
  VERIFICATION_TOKEN_ENTITY_KEY,
  {
    expires: timestamp('expires', { mode: 'date' }).notNull(),
    id: text('id')
      .default(sql`gen_random_uuid()`)
      .unique()
      .notNull(),
    identifier: text('identifier').notNull(),
    newEmail: text('newEmail'),
    token: text('token').notNull(),
    type: verificationTokenType('type').notNull()
  },
  (verificationToken) => [
    {
      compositePk: primaryKey({
        columns: [verificationToken.identifier, verificationToken.token]
      })
    }
  ]
)

export type VerificationTokenModel =
  (typeof VerificationTokenTable)['$inferSelect']
