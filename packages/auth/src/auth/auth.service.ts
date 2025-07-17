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
            modelName: 'InvitationTable'
          },
          member: {
            modelName: 'MemberTable'
          },
          organization: {
            modelName: 'OrganizationTable'
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
      modelName: 'VerificationTable'
    }
  })

  return {
    createOrganization: auth.api.createOrganization,
    handler: auth.handler,
    setActiveOrganization: auth.api.setActiveOrganization
  }
}

export type AuthService = ReturnType<typeof createAuthService>
