import { sql } from 'drizzle-orm'
import { date, text } from 'drizzle-orm/pg-core'

export const baseFields = {
  createdAt: date()
    .default(sql`CURRENT_TIMESTAMP`)
    .notNull(),
  id: text()
    .default(sql`gen_random_uuid()`)
    .primaryKey(),
  updatedAt: date()
    .default(sql`CURRENT_TIMESTAMP`)
    .notNull()
}
