import type { BetterAuthPlugin } from 'better-auth'

import { betterAuth } from 'better-auth'
import { drizzleAdapter } from 'better-auth/adapters/drizzle'
import { organization } from 'better-auth/plugins'
import { reactStartCookies } from 'better-auth/react-start'

import type { ConfigService } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'

export const createAuthService = (
  databaseService: DatabaseService,
  configService: ConfigService
) => {
  const auth = betterAuth({
    account: {
      modelName: 'AccountTable'
    },
    advanced: {
      crossSubDomainCookies: {
        domain: '.' + configService.get('ingress.domain'), // your domain
        enabled: true
      },
      generateId: false,
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
      }) satisfies BetterAuthPlugin,
      reactStartCookies()
    ],
    session: {
      cookieCache: {
        enabled: false
      },
      modelName: 'SessionTable'
    },
    trustedOrigins: configService.get('api.cors.origins').split(','),
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
    getSession: auth.api.getSession,
    handler: auth.handler,
    signInEmail: auth.api.signInEmail,
    signOut: auth.api.signOut,
    signUpEmail: auth.api.signUpEmail
  }
}

export type AuthService = ReturnType<typeof createAuthService>
