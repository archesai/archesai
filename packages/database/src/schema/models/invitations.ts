import { boolean, pgTable, text } from 'drizzle-orm/pg-core'

import { INVITATION_ENTITY_KEY } from '@archesai/schemas'

import { roleEnum } from '#schema/enums'
import { baseFields } from '#schema/models/base'
import { organizationFk } from '#schema/models/organization'

// TABLE
export const InvitationTable = pgTable(
  INVITATION_ENTITY_KEY,
  {
    ...baseFields,
    accepted: boolean('accepted').default(false),
    email: text('email').notNull(),
    role: roleEnum().default('USER').notNull()
  },
  (InvitationTable) => [organizationFk(InvitationTable)]
)

// SCHEMA
export type InvitationModel = (typeof InvitationTable)['$inferSelect']
