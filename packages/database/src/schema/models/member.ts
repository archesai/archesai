import { relations, sql } from 'drizzle-orm'
import { pgTable, text, uniqueIndex } from 'drizzle-orm/pg-core'

import { MEMBER_ENTITY_KEY } from '@archesai/schemas'

import { roleEnum } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { InvitationTable } from '#schema/models/invitations'
import { organizationFk } from '#schema/models/organization'
import { UserTable } from '#schema/models/user'

// TABLE
export const MemberTable = pgTable(
  MEMBER_ENTITY_KEY,
  {
    ...baseFields,
    invitationId: text('invitationId').references(() => InvitationTable.id, {
      onDelete: 'cascade'
    }),
    role: roleEnum().default('USER').notNull(),
    userId: text('userId')
      .notNull()
      .references(() => UserTable.id, { onDelete: 'cascade' })
  },
  (MemberTable) => [
    organizationFk(MemberTable),
    uniqueIndex()
      .on(MemberTable.userId, MemberTable.orgname)
      .where(sql`${MemberTable.userId} IS NOT NULL`)
  ]
)

// RELATIONS
export const memberRelations = relations(MemberTable, ({ one }) => ({
  user: one(UserTable, {
    fields: [MemberTable.userId],
    references: [UserTable.id],
    relationName: 'userMemberships'
  })
}))

// SCHEMA
export type MemberModel = (typeof MemberTable)['$inferSelect']
