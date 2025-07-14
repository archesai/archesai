import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

// export const SessionTable = pgTable('session', {
//   ...baseFields,
//   activeOrganizationId: text(),
//   expiresAt: timestamp({
//     mode: 'string'
//   }).notNull(),
//   ipAddress: text(),
//   token: text().notNull(),
//   userAgent: text(),
//   userId: text('userId')
//     .notNull()
//     .references(() => UserTable.id, { onDelete: 'cascade' })
// })

export const SessionEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    activeOrganizationId: Type.String({
      description: 'The active organization ID'
    }),
    expiresAt: Type.String({
      description: 'The expiration date of the session'
    }),
    ipAddress: Type.Optional(
      Type.String({ description: 'The IP address of the session' })
    ),
    token: Type.String({ description: 'The session token' }),
    userAgent: Type.Optional(
      Type.String({ description: 'The user agent of the session' })
    ),
    userId: Type.String({
      description: 'The ID of the user associated with the session'
    })
  },
  {
    $id: 'SessionEntity',
    description: 'The session entity',
    title: 'Session Entity'
  }
)

export type SessionEntity = Static<typeof SessionEntitySchema>
