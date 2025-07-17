import type { PgEnum } from 'drizzle-orm/pg-core'

import { pgEnum } from 'drizzle-orm/pg-core'

import {
  PlanTypes,
  RoleTypes,
  RunTypes,
  StatusTypes,
  VerificationTokenTypes
} from '@archesai/schemas'

export const roleEnum: PgEnum<['ADMIN', 'USER']> = pgEnum('role', RoleTypes)

export const verificationTokenType: PgEnum<
  ['EMAIL_CHANGE', 'EMAIL_VERIFICATION', 'PASSWORD_RESET']
> = pgEnum('verificationTokenType', VerificationTokenTypes)

export const runStatus: PgEnum<
  ['COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED']
> = pgEnum('runStatus', StatusTypes)

export const planType: PgEnum<
  ['BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED']
> = pgEnum('planType', PlanTypes)

export const runType: PgEnum<['PIPELINE_RUN', 'TOOL_RUN']> = pgEnum(
  'RunType',
  RunTypes
)
