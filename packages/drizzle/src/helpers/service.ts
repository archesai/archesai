import pg from 'pg'

import type { DatabaseService } from '@archesai/core'

import { createPooledClient } from '#helpers/clients'
import { DrizzleDatabaseService } from '#helpers/drizzle-database.service'

export const createDrizzleDatabaseService = (
  connectionString: string
): DatabaseService => {
  const pool = new pg.Pool({ connectionString })
  const db = createPooledClient(pool)
  return new DrizzleDatabaseService(db)
}
