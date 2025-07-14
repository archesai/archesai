import { createDrizzleDatabaseService } from '@archesai/database'

import { createAuthService } from '#auth/auth.service'

const databaseService = createDrizzleDatabaseService(process.env.DATABASE_URL!)
export const auth = createAuthService(databaseService)
