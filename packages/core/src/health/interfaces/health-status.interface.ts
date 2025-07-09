import type { StatusType } from '@archesai/schemas'

import type { Errors } from '#http/schemas/errors.schema'

export interface HealthStatus {
  errors?: Errors
  status: StatusType
}
