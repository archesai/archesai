import { relations } from 'drizzle-orm'
import { boolean, pgTable, text } from 'drizzle-orm/pg-core'

import { USER_ENTITY_KEY } from '@archesai/schemas'

import { AccountTable } from '#schema/models/account'
import { baseFields } from '#schema/models/base'
import { MemberTable } from '#schema/models/member'

export const UserTable = pgTable(USER_ENTITY_KEY, {
  ...baseFields,
  deactivated: boolean().default(false).notNull(),
  email: text().unique().notNull(),
  emailVerified: boolean().default(false).notNull(),
  image: text(),
  name: text().notNull()
})

export const userRelations = relations(UserTable, ({ many }) => ({
  accounts: many(AccountTable),
  memberships: many(MemberTable)
}))
