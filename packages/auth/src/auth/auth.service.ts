import { betterAuth } from 'better-auth'
import { drizzleAdapter } from 'better-auth/adapters/drizzle'
import { organization } from 'better-auth/plugins'

import type { DrizzleDatabaseService } from '@archesai/database'

export const createAuthService = (databaseService: DrizzleDatabaseService) => {
  const auth = betterAuth({
    account: {
      modelName: 'AccountTable'
    },
    advanced: {
      crossSubDomainCookies: {
        domain: '.archesai.test', // your domain
        enabled: true
      },
      useSecureCookies: true
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
      // reactStartCookies()
    ],
    session: {
      cookieCache: {
        enabled: false
      },
      modelName: 'SessionTable'
    },
    trustedOrigins: [
      'https://platform.archesai.test',
      'http://platform.archesai.test',
      'http://api.archesai.test',
      'https://api.archesai.test'
    ],
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
