import type { LibSQLDatabase } from 'drizzle-orm/libsql'

import { drizzle as libsqlDrizzle } from 'drizzle-orm/libsql/node'
import { drizzle as pgDrizzle } from 'drizzle-orm/node-postgres'
import pg from 'pg'

import type { DatabaseService } from '@archesai/core'

import { DrizzleDatabaseService } from '#adapters/drizzle-database.service'
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

export const createDrizzleDatabaseService = (
  connectionString: string
): DatabaseService => {
  const pool = new pg.Pool({ connectionString })
  const db = createPooledClient(pool)
  return new DrizzleDatabaseService(db)
}
