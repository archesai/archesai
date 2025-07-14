import { relations } from 'drizzle-orm'
import { pgTable, text, timestamp } from 'drizzle-orm/pg-core'

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
  role: roleEnum().default('USER').notNull(),
  status: text().notNull()
})

export type InvitationInsertModel = typeof InvitationTable.$inferInsert
export type InvitationSelectModel = typeof InvitationTable.$inferSelect

export const invitationRelations = relations(InvitationTable, ({ one }) => ({
  organization: one(OrganizationTable, {
    fields: [InvitationTable.organizationId],
    references: [OrganizationTable.id]
  })
}))
