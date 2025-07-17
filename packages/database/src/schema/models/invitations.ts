import { relations } from 'drizzle-orm'
import { pgTable, text, timestamp } from 'drizzle-orm/pg-core'

import type { InvitationEntity } from '@archesai/schemas'

import { INVITATION_ENTITY_KEY } from '@archesai/schemas'

import { roleEnum } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'
import { UserTable } from '#schema/models/user'

export const InvitationTable = pgTable(INVITATION_ENTITY_KEY, {
  ...baseFields,
  email: text('email').notNull(),
  expiresAt: timestamp({
    mode: 'string'
  }).notNull(),
  inviterId: text()
    .notNull()
    .references(() => UserTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  role: roleEnum().default('member').notNull(),
  status: text().notNull()
})

export const invitationRelations = relations(InvitationTable, ({ one }) => ({
  inviter: one(UserTable, {
    fields: [InvitationTable.inviterId],
    references: [UserTable.id]
  }),
  organization: one(OrganizationTable, {
    fields: [InvitationTable.organizationId],
    references: [OrganizationTable.id]
  })
}))

export type InvitationInsertModel = typeof InvitationTable.$inferInsert
export type InvitationSelectModel = typeof InvitationTable.$inferSelect

export type zInvitationCheck =
  InvitationEntity extends InvitationSelectModel ? true : false
