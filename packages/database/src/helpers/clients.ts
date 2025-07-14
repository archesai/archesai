import type { LibSQLDatabase } from 'drizzle-orm/libsql'
import type pg from 'pg'

import { drizzle as libsqlDrizzle } from 'drizzle-orm/libsql/node'
import { drizzle as pgDrizzle } from 'drizzle-orm/node-postgres'

import * as schema from '#schema/index'

export const createClient = (url: string) => {
  const db = pgDrizzle({
    // casing: 'snake_case',
    connection: url,
    // logger: true,
    schema
  })

  return db
}

export const createPooledClient = (pool: pg.Pool) => {
  const db = pgDrizzle(pool, {
    // connection: databaseUrl,
    // casing: 'snake_case',
    // logger: true,
    schema
  })

  return db
}

export const createLibsqlClient = (
  url: string
): LibSQLDatabase<typeof schema> => {
  const db = libsqlDrizzle({
    // casing: 'snake_case',
    connection: url,
    schema
  })

  return db
}
