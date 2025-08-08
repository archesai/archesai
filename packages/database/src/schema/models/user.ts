import { relations } from 'drizzle-orm'
import { boolean, pgTable, text } from 'drizzle-orm/pg-core'

import type { UserEntity } from '@archesai/schemas'

import { USER_ENTITY_KEY } from '@archesai/schemas'

import { baseFields } from '#schema/models/base'

export const UserTable = pgTable(USER_ENTITY_KEY, {
  ...baseFields,
  deactivated: boolean().default(false).notNull(),
  email: text().notNull().unique(),
  emailVerified: boolean().default(false).notNull(),
  image: text(),
  name: text().notNull()
})

export const userRelations = relations(UserTable, () => ({}))

export type UserInsertModel = typeof UserTable.$inferInsert
export type UserSelectModel = typeof UserTable.$inferSelect

export type zUserCheck = UserEntity extends UserSelectModel ? true : false
