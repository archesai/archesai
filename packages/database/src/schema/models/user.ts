import { relations } from 'drizzle-orm'
import { boolean, pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import { USER_ENTITY_KEY } from '@archesai/schemas'

import { AccountTable } from '#schema/models/account'
import { baseFields } from '#schema/models/base'
import { MemberTable } from '#schema/models/member'

// TABLE
export const UserTable = pgTable(USER_ENTITY_KEY, {
  ...baseFields,
  deactivated: boolean('deactivated').default(false).notNull(),
  email: text('email').unique(),
  emailVerified: timestamp('emailVerified', { mode: 'date' }),
  image: text('image')
})

// RELATIONS
export const userRelations = relations(UserTable, ({ many }) => ({
  accounts: many(AccountTable, {
    relationName: 'userAccounts'
  }),
  memberships: many(MemberTable, {
    relationName: 'userMemberships'
  })
}))

// SCHEMA
export type UserModel = (typeof UserTable)['$inferSelect']
