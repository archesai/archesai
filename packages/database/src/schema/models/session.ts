import { relations } from 'drizzle-orm'
import { date, pgTable, text } from 'drizzle-orm/pg-core'

import { baseFields } from '#schema/models/base'
import { UserTable } from '#schema/models/user'

export const SessionTable = pgTable('session', {
  ...baseFields,
  expiresAt: date().notNull(),
  ipAddress: text(),
  token: text().notNull(),
  userAgent: text(),
  userId: text('userId')
    .notNull()
    .references(() => UserTable.id, { onDelete: 'cascade' })
})

export const sessionRelations = relations(SessionTable, ({ one }) => ({
  user: one(UserTable, {
    fields: [SessionTable.userId],
    references: [UserTable.id]
  })
}))
