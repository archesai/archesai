import { pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import { UserTable } from '#schema/models/user'

export const SessionTable = pgTable('session', {
  expires: timestamp('expires', { mode: 'date' }).notNull(),
  sessionToken: text('sessionToken').primaryKey(),
  userId: text('userId')
    .notNull()
    .references(() => UserTable.id, { onDelete: 'cascade' })
})

export type SessionModel = (typeof SessionTable)['$inferSelect']
