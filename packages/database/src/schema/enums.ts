import type { PgEnum } from 'drizzle-orm/pg-core'

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

export const roleEnum: PgEnum<['ADMIN', 'USER']> = pgEnum('role', RoleTypes)

export const verificationTokenType: PgEnum<
  ['EMAIL_CHANGE', 'EMAIL_VERIFICATION', 'PASSWORD_RESET']
> = pgEnum('verificationTokenType', VerificationTokenTypes)

export const toolIO: PgEnum<['AUDIO', 'IMAGE', 'TEXT', 'VIDEO']> = pgEnum(
  'toolIO',
  ContentBaseTypes
)

export const runStatus: PgEnum<
  ['COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED']
> = pgEnum('runStatus', StatusTypes)

export const planType: PgEnum<
  ['BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED']
> = pgEnum('planType', PlanTypes)

export const authType: PgEnum<['email', 'oauth', 'oidc', 'webauthn']> = pgEnum(
  'authType',
  AuthTypes
)

export const providerType: PgEnum<['API_KEY', 'FIREBASE', 'LOCAL', 'TWITTER']> =
  pgEnum('providerType', ProviderTypes)

export const runType: PgEnum<['PIPELINE_RUN', 'TOOL_RUN']> = pgEnum(
  'RunType',
  RunTypes
)
