import { relations } from 'drizzle-orm'
import { pgTable, text } from 'drizzle-orm/pg-core'

import { MEMBER_ENTITY_KEY } from '@archesai/schemas'

import { roleEnum } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'
import { UserTable } from '#schema/models/user'

export const MemberTable = pgTable(MEMBER_ENTITY_KEY, {
  ...baseFields,
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  role: roleEnum().default('USER').notNull(),
  userId: text()
    .notNull()
    .references(() => UserTable.id, { onDelete: 'cascade' })
})

export type MemberInsertModel = typeof MemberTable.$inferInsert
export type MemberSelectModel = typeof MemberTable.$inferSelect

export const memberRelations = relations(MemberTable, ({ one }) => ({
  organization: one(OrganizationTable, {
    fields: [MemberTable.organizationId],
    references: [OrganizationTable.id]
  }),
  user: one(UserTable, {
    fields: [MemberTable.userId],
    references: [UserTable.id]
  })
}))
