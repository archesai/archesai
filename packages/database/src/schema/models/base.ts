import { sql } from 'drizzle-orm'
import { pgEnum, text, timestamp } from 'drizzle-orm/pg-core'

import { PlanTypes, RoleTypes, StatusTypes } from '@archesai/schemas'

export const baseFields = {
  createdAt: timestamp({
    mode: 'string'
  })
    .defaultNow()
    .notNull(),
  id: text()
    .default(sql`gen_random_uuid()`)
    .primaryKey(),
  updatedAt: timestamp({
    mode: 'string'
  })
    .defaultNow()
    .notNull()
}

export const roleEnum = pgEnum('role', RoleTypes)

export const planEnum = pgEnum('planType', PlanTypes)

export const statusEnum = pgEnum('runStatus', StatusTypes)
