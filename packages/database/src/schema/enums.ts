import type { PgEnum } from 'drizzle-orm/pg-core'

import { pgEnum } from 'drizzle-orm/pg-core'

import { PlanTypes, RoleTypes, StatusTypes } from '@archesai/schemas'

export const roleEnum: PgEnum<['admin', 'owner', 'member']> = pgEnum(
  'role',
  RoleTypes
)

export const runStatus: PgEnum<
  ['COMPLETED', 'FAILED', 'PROCESSING', 'QUEUED']
> = pgEnum('runStatus', StatusTypes)

export const planType: PgEnum<
  ['BASIC', 'FREE', 'PREMIUM', 'STANDARD', 'UNLIMITED']
> = pgEnum('planType', PlanTypes)
