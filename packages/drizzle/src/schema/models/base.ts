import { sql } from 'drizzle-orm'
import { text, timestamp } from 'drizzle-orm/pg-core'

export const baseFields = {
  createdAt: timestamp({ mode: 'string', precision: 3 })
    .default(sql`CURRENT_TIMESTAMP`)
    .notNull(),
  id: text('id')
    .default(sql`gen_random_uuid()`)
    .primaryKey(),
  name: text('name'),
  orgname: text('orgname').notNull(),
  updatedAt: timestamp({ mode: 'string', precision: 3 })
    .default(sql`CURRENT_TIMESTAMP`)
    .notNull()
}
