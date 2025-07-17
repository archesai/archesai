import { createDrizzleDatabaseService } from '@archesai/database'

import { createAuthService } from '#auth/auth.service'

const url = process.env.DATABASE_URL
if (!url) {
  throw new Error('DATABASE_URL environment variable is not set')
}

const databaseService = createDrizzleDatabaseService(url)
export const auth = createAuthService(databaseService)
