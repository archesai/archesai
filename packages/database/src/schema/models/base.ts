import { sql } from 'drizzle-orm'
import { text, timestamp } from 'drizzle-orm/pg-core'

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
