import { betterAuth } from 'better-auth'
import { drizzleAdapter } from 'better-auth/adapters/drizzle'
import { organization } from 'better-auth/plugins'
import { reactStartCookies } from 'better-auth/react-start'

import type { ConfigService } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { SessionEntity, UserEntity } from '@archesai/schemas'

export interface AuthService {
  getSession(headers: Headers): Promise<null | {
    session: BetterAuthSession | null
    user: BetterAuthUser
  }>
  signInEmail(body: { email: string; password: string }): Promise<{
    headers: Headers
    response: {
      user: BetterAuthUser
    }
  }>
  signOut(headers: Headers): Promise<{
    headers: Headers
  }>
  signUpEmail(body: {
    email: string
    name: string
    password: string
  }): Promise<{
    headers: Headers
    response: {
      user: BetterAuthUser
    }
  }>
}

type BetterAuthSession = Omit<
  SessionEntity,
  | 'activeOrganizationId'
  | 'createdAt'
  | 'expiresAt'
  | 'ipAddress'
  | 'updatedAt'
  | 'userAgent'
> & {
  activeOrganizationId?: null | string | undefined
  createdAt: Date
  expiresAt: Date
  ipAddress?: null | string | undefined
  updatedAt: Date
  userAgent?: null | string | undefined
}

type BetterAuthUser = Omit<UserEntity, 'createdAt' | 'image' | 'updatedAt'> & {
  createdAt: Date
  image?: null | string | undefined
  updatedAt: Date
}

export const createAuthService = (
  databaseService: DatabaseService,
  configService: ConfigService
): AuthService => {
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
      }),
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
      modelName: 'UserTable'
    },
    verification: {
      modelName: 'VerificationTable'
    }
  })

  return {
    getSession: (headers: Headers) => auth.api.getSession({ headers }),
    signInEmail: (body: { email: string; password: string }) =>
      auth.api.signInEmail({
        body,
        returnHeaders: true
      }),
    signOut: (headers: Headers) =>
      auth.api.signOut({
        headers,
        returnHeaders: true
      }),
    signUpEmail: (body: { email: string; name: string; password: string }) =>
      auth.api.signUpEmail({
        body,
        returnHeaders: true
      })
  }
}
