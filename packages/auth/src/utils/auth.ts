import type { DB } from 'better-auth/adapters/drizzle'

import { betterAuth } from 'better-auth'
import { drizzleAdapter } from 'better-auth/adapters/drizzle'
import { organization } from 'better-auth/plugins'

export const createAuthService = (db: DB) => {
  const auth = betterAuth({
    database: drizzleAdapter(db, {
      provider: 'pg'
    }),
    plugins: [organization()]
  })

  return auth.api
}

export type AuthService = ReturnType<typeof createAuthService>
