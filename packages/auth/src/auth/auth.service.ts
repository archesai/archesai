import { betterAuth } from 'better-auth'
import { drizzleAdapter } from 'better-auth/adapters/drizzle'
import { organization } from 'better-auth/plugins'

import type { DrizzleDatabaseService } from '@archesai/database'

export const createAuthService = (databaseService: DrizzleDatabaseService) => {
  const auth = betterAuth({
    account: {
      modelName: 'AccountTable'
    },
    database: drizzleAdapter(databaseService.db, {
      provider: 'pg'
    }),
    emailAndPassword: {
      enabled: true
    },
    emailVerification: {
      autoSignInAfterVerification: true
    },
    plugins: [
      organization({
        schema: {
          invitation: {
            modelName: 'invitations'
          },
          member: {
            modelName: 'members'
          },
          organization: {
            modelName: 'organizations'
          }
        }
      })
    ],
    session: {
      modelName: 'SessionTable'
    },

    trustedOrigins: ['https://platform.archesai.test'],
    user: {
      additionalFields: {
        deactivated: {
          defaultValue: false,
          type: 'boolean'
        }
      },
      modelName: 'UserTable'
    },
    verification: {
      modelName: 'VerificationTokenTable'
    }
  })

  return {
    handler: auth.handler
  }
}

export type AuthService = ReturnType<typeof createAuthService>
