import type { HealthStatus } from '#health/interfaces/health-status.interface'

export interface HealthCheck {
  getHealth: () => HealthStatus
}
