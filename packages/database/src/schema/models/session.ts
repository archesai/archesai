import { relations } from 'drizzle-orm'
import { pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import type { SessionEntity } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'
import { UserTable } from '#schema/models/user'

export const SessionTable = pgTable('session', {
  ...baseFields,
  activeOrganizationId: text(),
  expiresAt: timestamp({
    mode: 'string'
  }).notNull(),
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

export type SessionInsertModel = typeof SessionTable.$inferInsert
export type SessionSelectModel = typeof SessionTable.$inferSelect

export type zSessionCheck =
  SessionEntity extends SessionSelectModel ? true : false
