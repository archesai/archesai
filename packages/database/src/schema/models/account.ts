import { relations, sql } from 'drizzle-orm'
import { integer, pgTable, primaryKey, text } from 'drizzle-orm/pg-core'

import { ACCOUNT_ENTITY_KEY } from '@archesai/domain'

import { authType, providerType } from '#schema/enums'
import { UserTable } from '#schema/models/user'

export const AccountTable = pgTable(
  ACCOUNT_ENTITY_KEY,
  {
    access_token: text('access_token'),
    authType: authType('authType').notNull(),
    expires_at: integer('expires_at'),
    // my custom stuff
    hashed_password: text('hashed_password'),
    id: text('id')
      .default(sql`gen_random_uuid()`)
      .unique()
      .notNull(),
    id_token: text('id_token'),
    provider: providerType().notNull(),
    providerAccountId: text('providerAccountId').notNull(),
    refresh_token: text('refresh_token'),
    scope: text('scope'),
    session_state: text('session_state'),
    token_type: text('token_type'),
    userId: text('userId')
      .notNull()
      .references(() => UserTable.id, { onDelete: 'cascade' })
  },
  (account) => [
    primaryKey({
      columns: [account.provider, account.providerAccountId]
    })
  ]
)

export const accountRelations = relations(AccountTable, ({ one }) => ({
  user: one(UserTable, {
    fields: [AccountTable.userId],
    references: [UserTable.id],
    relationName: 'userAccounts'
  })
}))

export type AccountModel = (typeof AccountTable)['$inferSelect']
