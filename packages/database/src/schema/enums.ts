import { pgEnum } from 'drizzle-orm/pg-core'

import {
  AuthTypes,
  ContentBaseTypes,
  PlanTypes,
  ProviderTypes,
  RoleTypes,
  RunTypes,
  StatusTypes,
  VerificationTokenTypes
} from '@archesai/schemas'

export const roleEnum = pgEnum('role', RoleTypes)

export const verificationTokenType = pgEnum(
  'verificationTokenType',
  VerificationTokenTypes
)

export const toolIO = pgEnum('toolIO', ContentBaseTypes)

export const runStatus = pgEnum('runStatus', StatusTypes)

export const planType = pgEnum('planType', PlanTypes)

export const authType = pgEnum('authType', AuthTypes)

export const providerType = pgEnum('providerType', ProviderTypes)

export const runType = pgEnum('RunType', RunTypes)
