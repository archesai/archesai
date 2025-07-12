import { relations } from 'drizzle-orm'
import { boolean, pgTable, text } from 'drizzle-orm/pg-core'

import { INVITATION_ENTITY_KEY } from '@archesai/schemas'

import { roleEnum } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { OrganizationTable } from '#schema/models/organization'

export const InvitationTable = pgTable(INVITATION_ENTITY_KEY, {
  ...baseFields,
  accepted: boolean('accepted').default(false),
  email: text('email').notNull(),
  organizationId: text()
    .notNull()
    .references(() => OrganizationTable.id, {
      onDelete: 'cascade',
      onUpdate: 'cascade'
    }),
  role: roleEnum().default('USER').notNull()
})

export type InvitationInsertModel = typeof InvitationTable.$inferInsert
export type InvitationSelectModel = typeof InvitationTable.$inferSelect

export const invitationRelations = relations(InvitationTable, ({ one }) => ({
  organization: one(OrganizationTable, {
    fields: [InvitationTable.organizationId],
    references: [OrganizationTable.id]
  })
}))
