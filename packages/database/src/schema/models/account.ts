import { relations } from 'drizzle-orm'
import { date, pgTable, text } from 'drizzle-orm/pg-core'

import { ACCOUNT_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { UserTable } from '#schema/models/user'

export const AccountTable = pgTable(ACCOUNT_ENTITY_KEY, {
  ...baseFields,
  accessToken: text(),
  accessTokenExpiresAt: date(),
  accountId: text().notNull(),
  idToken: text(),
  password: text(),
  providerId: text().notNull(),
  refreshToken: text(),
  refreshTokenExpiresAt: date(),
  scope: text(),
  userId: text()
    .notNull()
    .references(() => UserTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    })
})

export const accountRelations = relations(AccountTable, ({ one }) => ({
  user: one(UserTable, {
    fields: [AccountTable.userId],
    references: [UserTable.id]
  })
}))
