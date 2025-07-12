import type { StatusType } from '@archesai/schemas'

import type { ErrorObject } from '#http/schemas/error-object.schema'

export interface HealthStatus {
  errors?: ErrorObject[]
  status: StatusType
}
