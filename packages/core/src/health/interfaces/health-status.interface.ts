import type { StatusType } from '@archesai/domain'

import type { Errors } from '#http/schemas/errors.schema'

export interface HealthStatus {
  errors?: Errors
  status: StatusType
}
