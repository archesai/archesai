import {
  boolean,
  integer,
  pgTable,
  primaryKey,
  text
} from 'drizzle-orm/pg-core'

import { UserTable } from '#schema/models/user'

export const AuthenticatorTable = pgTable(
  'authenticator',
  {
    counter: integer('counter').notNull(),
    credentialBackedUp: boolean('credentialBackedUp').notNull(),
    credentialDeviceType: text('credentialDeviceType').notNull(),
    credentialID: text('credentialID').notNull().unique(),
    credentialPublicKey: text('credentialPublicKey').notNull(),
    providerAccountId: text('providerAccountId').notNull(),
    transports: text('transports'),
    userId: text('userId')
      .notNull()
      .references(() => UserTable.id, { onDelete: 'cascade' })
  },
  (authenticator) => [
    primaryKey({
      columns: [authenticator.userId, authenticator.credentialID]
    })
  ]
)

export type AuthenticatorModel = (typeof AuthenticatorTable)['$inferSelect']
