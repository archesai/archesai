import { relations } from 'drizzle-orm'
import { pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import type { AccountEntity } from '@archesai/schemas'

import { ACCOUNT_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { UserTable } from '#schema/models/user'

export const AccountTable = pgTable(ACCOUNT_ENTITY_KEY, {
  ...baseFields,
  accessToken: text(),
  accessTokenExpiresAt: timestamp({
    mode: 'string'
  }),
  accountId: text().notNull(),
  idToken: text(),
  password: text(),
  providerId: text().notNull(),
  refreshToken: text(),
  refreshTokenExpiresAt: timestamp({
    mode: 'string'
  }),
  scope: text(),
  userId: text()
    .notNull()
    .references(() => UserTable.id, {
      onDelete: 'cascade'
    })
})

export type AccountInsertModel = typeof AccountTable.$inferInsert
export type AccountSelectModel = typeof AccountTable.$inferSelect

export const accountRelations = relations(AccountTable, ({ one }) => ({
  user: one(UserTable, {
    fields: [AccountTable.userId],
    references: [UserTable.id]
  })
}))

export type zAccountCheck =
  AccountEntity extends AccountSelectModel ? true : false
